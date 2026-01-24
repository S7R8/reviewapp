package knowledge

import (
	"context"
	"fmt"
	"log"

	"github.com/s7r8/reviewapp/internal/domain/repository"
	"github.com/s7r8/reviewapp/internal/infrastructure/external"
)

// UpdateKnowledgeUseCase - ナレッジ更新のユースケース
type UpdateKnowledgeUseCase struct {
	knowledgeRepo   repository.KnowledgeRepository
	embeddingClient external.EmbeddingClientInterface
}

// NewUpdateKnowledgeUseCase - コンストラクタ
func NewUpdateKnowledgeUseCase(
	knowledgeRepo repository.KnowledgeRepository,
	embeddingClient external.EmbeddingClientInterface,
) *UpdateKnowledgeUseCase {
	return &UpdateKnowledgeUseCase{
		knowledgeRepo:   knowledgeRepo,
		embeddingClient: embeddingClient,
	}
}

// UpdateKnowledgeInput - 入力
type UpdateKnowledgeInput struct {
	UserID      string
	KnowledgeID string
	Title       string
	Content     string
	Category    string
	Priority    int
}

// UpdateKnowledgeOutput - 出力
type UpdateKnowledgeOutput struct {
	Success bool
	Message string
}

// Execute - ナレッジ更新を実行
func (uc *UpdateKnowledgeUseCase) Execute(ctx context.Context, input UpdateKnowledgeInput) (*UpdateKnowledgeOutput, error) {
	// 1. 既存のナレッジを取得
	knowledge, err := uc.knowledgeRepo.FindByID(ctx, input.KnowledgeID)
	if err != nil {
		return nil, fmt.Errorf("knowledge not found: %w", err)
	}

	// 2. 権限チェック（自分のナレッジのみ更新可能）
	if knowledge.UserID != input.UserID {
		return nil, fmt.Errorf("permission denied: you can only update your own knowledge")
	}

	// 3. 内容を更新
	if err := knowledge.UpdateContent(input.Title, input.Content, input.Category, input.Priority); err != nil {
		return nil, fmt.Errorf("invalid update data: %w", err)
	}

	// 4. タイトルまたはコンテンツが変更された場合、Embeddingを再生成
	embeddingText := input.Title + "\n\n" + input.Content
	embedding, err := uc.embeddingClient.GenerateEmbedding(ctx, embeddingText)
	if err != nil {
		log.Printf("Warning: failed to regenerate embedding for knowledge %s: %v", input.KnowledgeID, err)
	} else {
		knowledge.SetEmbedding(embedding)
	}

	// 5. リポジトリで更新
	if err := uc.knowledgeRepo.Update(ctx, knowledge); err != nil {
		return nil, fmt.Errorf("failed to update knowledge: %w", err)
	}

	return &UpdateKnowledgeOutput{
		Success: true,
		Message: "ナレッジを更新しました",
	}, nil
}
