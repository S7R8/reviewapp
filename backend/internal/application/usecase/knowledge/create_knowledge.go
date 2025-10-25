package knowledge

import (
	"context"
	"fmt"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/domain/repository"
)

// CreateKnowledgeUseCase - ナレッジ作成のユースケース
type CreateKnowledgeUseCase struct {
	knowledgeRepo repository.KnowledgeRepository
}

// NewCreateKnowledgeUseCase - コンストラクタ
func NewCreateKnowledgeUseCase(knowledgeRepo repository.KnowledgeRepository) *CreateKnowledgeUseCase {
	return &CreateKnowledgeUseCase{
		knowledgeRepo: knowledgeRepo,
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
	
	// 2. リポジトリに保存
	if err := uc.knowledgeRepo.Create(ctx, knowledge); err != nil {
		return nil, fmt.Errorf("failed to create knowledge: %w", err)
	}
	
	// 3. 作成されたナレッジを返す
	return &CreateKnowledgeOutput{
		Knowledge: knowledge,
	}, nil
}
