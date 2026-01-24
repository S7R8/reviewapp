package review

import (
	"context"
	"fmt"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/domain/repository"
)

// GetReviewUseCase - レビュー詳細取得のユースケース
type GetReviewUseCase struct {
	reviewRepo repository.ReviewRepository
}

// NewGetReviewUseCase - コンストラクタ
func NewGetReviewUseCase(
	reviewRepo repository.ReviewRepository,
) *GetReviewUseCase {
	return &GetReviewUseCase{
		reviewRepo: reviewRepo,
	}
}

// GetReviewInput - 入力
type GetReviewInput struct {
	ReviewID string
	UserID   string
}

// GetReviewOutput - 出力
type GetReviewOutput struct {
	Review *model.Review
}

// Execute - レビュー詳細を取得
func (uc *GetReviewUseCase) Execute(ctx context.Context, input GetReviewInput) (*GetReviewOutput, error) {
	// 0. バリデーション
	if input.ReviewID == "" {
		return nil, fmt.Errorf("review ID is required")
	}
	if input.UserID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	// 1. レビューを取得
	review, err := uc.reviewRepo.FindByID(ctx, input.ReviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to find review: %w", err)
	}

	// 2. 権限チェック（自分のレビューのみ取得可能）
	if review.UserID != input.UserID {
		return nil, fmt.Errorf("このレビューにアクセスする権限がありません")
	}

	return &GetReviewOutput{
		Review: review,
	}, nil
}
