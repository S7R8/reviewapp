package knowledge

import (
	"context"
	"fmt"

	"github.com/s7r8/reviewapp/internal/domain/repository"
)

// DeleteKnowledgeInput - 入力
type DeleteKnowledgeInput struct {
	UserID      string
	KnowledgeID string
}

// DeleteKnowledgeOutput - 出力
type DeleteKnowledgeOutput struct {
	Success bool
}

// DeleteKnowledgeUseCase - ナレッジ削除UseCase
type DeleteKnowledgeUseCase struct {
	knowledgeRepo repository.KnowledgeRepository
}

// NewDeleteKnowledgeUseCase - コンストラクタ
func NewDeleteKnowledgeUseCase(knowledgeRepo repository.KnowledgeRepository) *DeleteKnowledgeUseCase {
	return &DeleteKnowledgeUseCase{
		knowledgeRepo: knowledgeRepo,
	}
}

// Execute - ナレッジを削除
func (uc *DeleteKnowledgeUseCase) Execute(ctx context.Context, input DeleteKnowledgeInput) (*DeleteKnowledgeOutput, error) {
	// 1. ナレッジの存在確認と権限チェック
	knowledge, err := uc.knowledgeRepo.FindByID(ctx, input.KnowledgeID)
	if err != nil {
		return nil, fmt.Errorf("ナレッジが見つかりません: %w", err)
	}

	// 2. 所有者チェック
	if knowledge.UserID != input.UserID {
		return nil, fmt.Errorf("このナレッジを削除する権限がありません")
	}

	// 3. 論理削除
	if err := uc.knowledgeRepo.Delete(ctx, input.KnowledgeID); err != nil {
		return nil, fmt.Errorf("ナレッジの削除に失敗しました: %w", err)
	}

	return &DeleteKnowledgeOutput{
		Success: true,
	}, nil
}
