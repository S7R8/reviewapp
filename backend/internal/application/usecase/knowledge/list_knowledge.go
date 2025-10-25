package knowledge

import (
	"context"
	"fmt"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/domain/repository"
)

// ListKnowledgeUseCase - ナレッジ一覧取得のユースケース
type ListKnowledgeUseCase struct {
	knowledgeRepo repository.KnowledgeRepository
}

// NewListKnowledgeUseCase - コンストラクタ
func NewListKnowledgeUseCase(knowledgeRepo repository.KnowledgeRepository) *ListKnowledgeUseCase {
	return &ListKnowledgeUseCase{
		knowledgeRepo: knowledgeRepo,
	}
}

// ListKnowledgeInput - 入力
type ListKnowledgeInput struct {
	UserID   string
	Category string
}

// ListKnowledgeOutput - 出力
type ListKnowledgeOutput struct {
	Knowledges []*model.Knowledge
}

// Execute - ナレッジ一覧取得を実行
func (uc *ListKnowledgeUseCase) Execute(ctx context.Context, input ListKnowledgeInput) (*ListKnowledgeOutput, error) {
	var knowledges []*model.Knowledge
	var err error

	// カテゴリ指定の有無で分岐
	if input.Category != "" {
		// カテゴリでフィルタ
		knowledges, err = uc.knowledgeRepo.FindByCategory(ctx, input.UserID, input.Category)
		if err != nil {
			return nil, fmt.Errorf("failed to find knowledge by category: %w", err)
		}
	} else {
		// 全件取得
		knowledges, err = uc.knowledgeRepo.FindByUserID(ctx, input.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to find knowledge by user_id: %w", err)
		}
	}

	// 空配列でもエラーではない（正常系）
	return &ListKnowledgeOutput{
		Knowledges: knowledges,
	}, nil
}
