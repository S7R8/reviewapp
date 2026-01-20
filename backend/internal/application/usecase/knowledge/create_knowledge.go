package knowledge

import (
	"context"
	"fmt"
	"log"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/domain/repository"
	"github.com/s7r8/reviewapp/internal/infrastructure/external"
)

// CreateKnowledgeUseCase - ナレッジ作成のユースケース
type CreateKnowledgeUseCase struct {
	knowledgeRepo   repository.KnowledgeRepository
	embeddingClient external.EmbeddingClientInterface
}

// NewCreateKnowledgeUseCase - コンストラクタ
func NewCreateKnowledgeUseCase(
	knowledgeRepo repository.KnowledgeRepository,
	embeddingClient external.EmbeddingClientInterface,
) *CreateKnowledgeUseCase {
	return &CreateKnowledgeUseCase{
		knowledgeRepo:   knowledgeRepo,
		embeddingClient: embeddingClient,
	}
}

// CreateKnowledgeInput - 入力
type CreateKnowledgeInput struct {
	UserID   string
	Title    string
	Content  string
	Category string
	Priority int
}

// CreateKnowledgeOutput - 出力
type CreateKnowledgeOutput struct {
	Knowledge *model.Knowledge
}

// Execute - ナレッジ作成を実行
func (uc *CreateKnowledgeUseCase) Execute(ctx context.Context, input CreateKnowledgeInput) (*CreateKnowledgeOutput, error) {
	// 1. ドメインモデルを作成（バリデーション含む）
	knowledge, err := model.NewKnowledge(
		input.UserID,
		input.Title,
		input.Content,
		input.Category,
		input.Priority,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid knowledge data: %w", err)
	}

	// 2. Embeddingベクトルを生成して設定
	embeddingText := input.Title + "\n\n" + input.Content
	embedding, err := uc.embeddingClient.GenerateEmbedding(ctx, embeddingText)
	if err != nil {
		log.Printf("Warning: failed to generate embedding for knowledge: %v", err)
	} else {
		knowledge.SetEmbedding(embedding)
	}

	// 3. リポジトリに保存
	if err := uc.knowledgeRepo.Create(ctx, knowledge); err != nil {
		return nil, fmt.Errorf("failed to create knowledge: %w", err)
	}

	// 4. 作成されたナレッジを返す
	return &CreateKnowledgeOutput{
		Knowledge: knowledge,
	}, nil
}
