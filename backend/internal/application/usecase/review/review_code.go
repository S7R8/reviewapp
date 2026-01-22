package review

import (
	"context"
	"fmt"
	"log"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/domain/repository"
	"github.com/s7r8/reviewapp/internal/domain/service"
	"github.com/s7r8/reviewapp/internal/infrastructure/external"
	"github.com/s7r8/reviewapp/internal/infrastructure/parser"
)

// ReviewCodeUseCase - コードレビューのユースケース
type ReviewCodeUseCase struct {
	reviewRepo      repository.ReviewRepository
	knowledgeRepo   repository.KnowledgeRepository
	reviewService   *service.ReviewService
	claudeClient    external.ClaudeClientInterface
	embeddingClient external.EmbeddingClientInterface
}

// NewReviewCodeUsecase - コンストラクタ
func NewReviewCodeUseCase(
	reviewRepo repository.ReviewRepository,
	knowledgeRepo repository.KnowledgeRepository,
	reviewService *service.ReviewService,
	claudeClient external.ClaudeClientInterface,
	embeddingClient external.EmbeddingClientInterface,
) *ReviewCodeUseCase {
	return &ReviewCodeUseCase{
		reviewRepo:      reviewRepo,
		knowledgeRepo:   knowledgeRepo,
		reviewService:   reviewService,
		claudeClient:    claudeClient,
		embeddingClient: embeddingClient,
	}
}

// ReviewCodeInput - 入力
type ReviewCodeInput struct {
	UserID   string
	Code     string
	Language string
	Context  string // オプショナル
}

// ReviewCodeOutput - 出力
type ReviewCodeOutput struct {
	Review *model.Review
}

// Execute - コードレビューを実行
func (uc *ReviewCodeUseCase) Execute(ctx context.Context, input ReviewCodeInput) (*ReviewCodeOutput, error) {
	// 1. コードからEmbeddingを生成
	embeddingText := fmt.Sprintf("Language: %s\n\n%s", input.Language, input.Code)
	if input.Context != "" {
		embeddingText += fmt.Sprintf("\n\nContext: %s", input.Context)
	}

	embedding, err := uc.embeddingClient.GenerateEmbedding(ctx, embeddingText)
	if err != nil {
		// Embeddingエラーの場合、全ナレッジ取得にフォールバック
		log.Printf("Warning: failed to generate embedding, falling back to all knowledge: %v", err)
		return uc.executeWithAllKnowledge(ctx, input)
	}

	// 2. ベクトル類似度検索で関連ナレッジを取得（RAG: Retrieval）
	knowledges, err := uc.knowledgeRepo.SearchBySimilarity(
		ctx,
		input.UserID,
		embedding,
		10,   // Top-10ナレッジを取得
		0.35, // 類似度35%以上のものだけ
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search knowledge by similarity: %w", err)
	}

	// ナレッジが見つからない場合のログ
	if len(knowledges) == 0 {
		log.Printf("No knowledge found with similarity >= 0.35 for user %s", input.UserID)
	} else {
		log.Printf("Found %d relevant knowledge items for review", len(knowledges))
	}

	// 3. プロンプト生成（RAG: Augmented）
	knowledgePrompt, usedKnowledges := uc.reviewService.BuildPromptFromKnowledge(knowledges)

	// 4. LLMでレビュー生成（RAG: Generation）
	reviewResult, err := uc.claudeClient.ReviewCode(ctx, external.ReviewCodeInput{
		Code:            input.Code,
		Language:        input.Language,
		Context:         input.Context,
		KnowledgePrompt: knowledgePrompt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to review code: %w", err)
	}

	// 5. マークダウンを構造化データに変換
	structuredResult := parser.ParseReviewMarkdown(reviewResult.ReviewResult)

	// 6. レビューエンティティを作成
	review := model.NewReview(
		input.UserID,
		input.Code,
		input.Language,
		input.Context,
	)

	// 7. レビュー結果を設定（実際に使用したナレッジIDのみ記録）
	knowledgeIDs := extractKnowledgeIDs(usedKnowledges)
	review.SetReviewResult(
		reviewResult.ReviewResult,
		structuredResult,
		knowledgeIDs,
		"claude",                     // LLM Provider
		"claude-3-5-sonnet-20241022", // LLM Model (TODO: configから取得)
		reviewResult.TokensUsed,
	)

	// 8. レビュー結果を保存
	if err := uc.reviewRepo.Create(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to save review: %w", err)
	}

	// 9. ナレッジの使用カウントを更新（実際に使用したナレッジのみ）
	if err := uc.updateKnowledgeUsage(ctx, usedKnowledges); err != nil {
		// 更新失敗してもレビュー結果は返す
		log.Printf("Warning: failed to update knowledge usage: %v", err)
	}

	return &ReviewCodeOutput{
		Review: review,
	}, nil
}

// updateKnowledgeUsage - ナレッジの使用カウントと最終使用日時を更新
func (uc *ReviewCodeUseCase) updateKnowledgeUsage(ctx context.Context, knowledges []*model.Knowledge) error {
	for _, k := range knowledges {
		k.IncrementUsage()
		if err := uc.knowledgeRepo.Update(ctx, k); err != nil {
			return fmt.Errorf("failed to update knowledge %s: %w", k.ID, err)
		}
	}
	return nil
}

// executeWithAllKnowledge - Embedding生成失敗時のフォールバック処理
func (uc *ReviewCodeUseCase) executeWithAllKnowledge(ctx context.Context, input ReviewCodeInput) (*ReviewCodeOutput, error) {
	// 全ナレッジを取得
	knowledges, err := uc.knowledgeRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find knowledge: %w", err)
	}

	// プロンプト生成
	knowledgePrompt, usedKnowledges := uc.reviewService.BuildPromptFromKnowledge(knowledges)

	// LLMでレビュー生成
	reviewResult, err := uc.claudeClient.ReviewCode(ctx, external.ReviewCodeInput{
		Code:            input.Code,
		Language:        input.Language,
		Context:         input.Context,
		KnowledgePrompt: knowledgePrompt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to review code: %w", err)
	}

	// マークダウンを構造化データに変換
	structuredResult := parser.ParseReviewMarkdown(reviewResult.ReviewResult)

	// レビューエンティティを作成
	review := model.NewReview(
		input.UserID,
		input.Code,
		input.Language,
		input.Context,
	)

	// レビュー結果を設定
	knowledgeIDs := extractKnowledgeIDs(usedKnowledges)
	review.SetReviewResult(
		reviewResult.ReviewResult,
		structuredResult,
		knowledgeIDs,
		"claude",
		"claude-3-5-sonnet-20241022",
		reviewResult.TokensUsed,
	)

	// レビュー結果を保存
	if err := uc.reviewRepo.Create(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to save review: %w", err)
	}

	// ナレッジの使用カウントを更新
	if err := uc.updateKnowledgeUsage(ctx, usedKnowledges); err != nil {
		log.Printf("Warning: failed to update knowledge usage: %v", err)
	}

	return &ReviewCodeOutput{
		Review: review,
	}, nil
}

// extractKnowledgeIDs - ナレッジIDのリストを抽出
func extractKnowledgeIDs(knowledges []*model.Knowledge) []string {
	ids := make([]string, len(knowledges))
	for i, k := range knowledges {
		ids[i] = k.ID
	}
	return ids
}
