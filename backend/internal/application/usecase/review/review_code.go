package review

import (
	"context"
	"fmt"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/domain/repository"
	"github.com/s7r8/reviewapp/internal/domain/service"
	"github.com/s7r8/reviewapp/internal/infrastructure/external"
	"github.com/s7r8/reviewapp/internal/infrastructure/parser"
)

// ReviewCodeUseCase - コードレビューのユースケース
type ReviewCodeUseCase struct {
	reviewRepo    repository.ReviewRepository
	knowledgeRepo repository.KnowledgeRepository
	reviewService *service.ReviewService
	claudeClient  external.ClaudeClientInterface
}

// NewReviewCodeUsecase - コンストラクタ
func NewReviewCodeUseCase(
	reviewRepo repository.ReviewRepository,
	knowledgeRepo repository.KnowledgeRepository,
	reviewService *service.ReviewService,
	claudeClient external.ClaudeClientInterface,
) *ReviewCodeUseCase {
	return &ReviewCodeUseCase{
		reviewRepo:    reviewRepo,
		knowledgeRepo: knowledgeRepo,
		reviewService: reviewService,
		claudeClient:  claudeClient,
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
	// 1. 関連ナレッジを検索（RAG: Retrieval）
	knowledges, err := uc.knowledgeRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find knowledge: %w", err)
	}

	// 2. LLMでレビュー生成（RAG: Augmented Generation）
	knowledgePrompt := uc.reviewService.BuildPromptFromKnowledge(knowledges)

	// 3. レビューを保存
	reviewResult, err := uc.claudeClient.ReviewCode(ctx, external.ReviewCodeInput{
		Code:            input.Code,
		Language:        input.Language,
		Context:         input.Context,
		KnowledgePrompt: knowledgePrompt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to review code: %w", err)
	}

	// 4. マークダウンを構造化データに変換
	structuredResult := parser.ParseReviewMarkdown(reviewResult.ReviewResult)

	// 5. レビューエンティティを作成
	review := model.NewReview(
		input.UserID,
		input.Code,
		input.Language,
		input.Context,
	)

	// 6. レビュー結果を設定
	knowledgeIDs := extractKnowledgeIDs(knowledges)
	review.SetReviewResult(
		reviewResult.ReviewResult,
		structuredResult,
		knowledgeIDs,
		"claude",                     // LLM Provider
		"claude-3-5-sonnet-20241022", // LLM Model (TODO: configから取得)
		reviewResult.TokensUsed,
	)

	// 7. レビュー結果を保存
	if err := uc.reviewRepo.Create(ctx, review); err != nil {
		return nil, fmt.Errorf("failed to save review: %w", err)
	}

	// 8. ナレッジの使用カウントを更新
	if err := uc.updateKnowledgeUsage(ctx, knowledges); err != nil {
		// 更新失敗してもレビュー結果は返す
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

// extractKnowledgeIDs - ナレッジIDのリストを抽出
func extractKnowledgeIDs(knowledges []*model.Knowledge) []string {
	ids := make([]string, len(knowledges))
	for i, k := range knowledges {
		ids[i] = k.ID
	}
	return ids
}
