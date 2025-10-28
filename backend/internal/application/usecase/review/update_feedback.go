package review

import (
	"context"
	"fmt"
	"time"

	"github.com/s7r8/reviewapp/internal/domain/repository"
)

// UpdateFeedbackUseCase - フィードバック更新のユースケース
type UpdateFeedbackUseCase struct {
	reviewRepo repository.ReviewRepository
}

// NewUpdateFeedbackUseCase - コンストラクタ
func NewUpdateFeedbackUseCase(reviewRepo repository.ReviewRepository) *UpdateFeedbackUseCase {
	return &UpdateFeedbackUseCase{
		reviewRepo: reviewRepo,
	}
}

// UpdateFeedbackInput - 入力
type UpdateFeedbackInput struct {
	ReviewID string
	UserID   string // 権限チェック用
	Score    int
	Comment  string
}

// UpdateFeedbackOutput - 出力
type UpdateFeedbackOutput struct {
	ReviewID        string
	FeedbackScore   int
	FeedbackComment string
	UpdatedAt       string
}

// Execute - フィードバック更新を実行
func (uc *UpdateFeedbackUseCase) Execute(ctx context.Context, input UpdateFeedbackInput) (*UpdateFeedbackOutput, error) {
	// 1. バリデーション
	if err := uc.validate(input); err != nil {
		return nil, err
	}

	// 2. レビューの存在確認
	review, err := uc.reviewRepo.FindByID(ctx, input.ReviewID)
	if err != nil {
		return nil, fmt.Errorf("レビューが見つかりません: %w", err)
	}

	// 3. 権限チェック
	if review.UserID != input.UserID {
		return nil, fmt.Errorf("このレビューを更新する権限がありません")
	}

	// 4. フィードバック更新
	if err := uc.reviewRepo.UpdateFeedback(ctx, input.ReviewID, input.Score, input.Comment); err != nil {
		return nil, fmt.Errorf("フィードバックの更新に失敗しました: %w", err)
	}

	// 5. 出力を生成（入力値を使用）
	output := &UpdateFeedbackOutput{
		ReviewID:        input.ReviewID,
		FeedbackScore:   input.Score,
		FeedbackComment: input.Comment,
		UpdatedAt:       time.Now().Format("2006-01-02T15:04:05Z07:00"),
	}

	return output, nil
}

// validate - バリデーション
func (uc *UpdateFeedbackUseCase) validate(input UpdateFeedbackInput) error {
	if input.ReviewID == "" {
		return fmt.Errorf("レビューIDは必須です")
	}

	if input.UserID == "" {
		return fmt.Errorf("ユーザーIDは必須です")
	}

	if input.Score < 1 || input.Score > 3 {
		return fmt.Errorf("スコアは1-3の整数で指定してください")
	}

	if len(input.Comment) > 500 {
		return fmt.Errorf("コメントは500文字以内にしてください")
	}

	return nil
}
