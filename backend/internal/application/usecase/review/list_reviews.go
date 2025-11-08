package review

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/domain/repository"
)

// ListReviewsUseCase - レビュー履歴一覧取得UseCase
type ListReviewsUseCase struct {
	reviewRepo repository.ReviewRepository
}

// NewListReviewsUseCase - コンストラクタ
func NewListReviewsUseCase(reviewRepo repository.ReviewRepository) *ListReviewsUseCase {
	return &ListReviewsUseCase{
		reviewRepo: reviewRepo,
	}
}

// ListReviewsInput - 入力
type ListReviewsInput struct {
	UserID    string
	Page      int
	PageSize  int
	Language  string
	Status    string
	SortBy    string
	SortOrder string
	DateFrom  *time.Time
	DateTo    *time.Time
}

// ListReviewsOutput - 出力
type ListReviewsOutput struct {
	Items      []*model.Review
	Total      int
	Page       int
	PageSize   int
	TotalPages int
}

// Execute - 実行
func (u *ListReviewsUseCase) Execute(ctx context.Context, input ListReviewsInput) (*ListReviewsOutput, error) {
	// デフォルト値の設定
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PageSize <= 0 {
		input.PageSize = 10
	}
	if input.PageSize > 100 {
		input.PageSize = 100 // 最大100件まで
	}
	if input.SortBy == "" {
		input.SortBy = "created_at"
	}
	if input.SortOrder == "" {
		input.SortOrder = "desc"
	}

	// バリデーション
	if err := u.validateInput(input); err != nil {
		return nil, err
	}

	// フィルター条件を構築
	filters := make(map[string]interface{})
	if input.Language != "" {
		filters["language"] = input.Language
	}
	if input.Status != "" {
		filters["status"] = input.Status
	}
	if input.DateFrom != nil {
		filters["date_from"] = *input.DateFrom
	}
	if input.DateTo != nil {
		filters["date_to"] = *input.DateTo
	}

	// LIMIT/OFFSETの計算
	offset := (input.Page - 1) * input.PageSize

	// レビュー一覧を取得
	reviews, err := u.reviewRepo.ListWithFilters(
		ctx,
		input.UserID,
		filters,
		input.SortBy,
		input.SortOrder,
		input.PageSize,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list reviews: %w", err)
	}

	// 総件数を取得
	total, err := u.reviewRepo.CountWithFilters(ctx, input.UserID, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to count reviews: %w", err)
	}

	// 総ページ数を計算
	totalPages := int(math.Ceil(float64(total) / float64(input.PageSize)))

	return &ListReviewsOutput{
		Items:      reviews,
		Total:      total,
		Page:       input.Page,
		PageSize:   input.PageSize,
		TotalPages: totalPages,
	}, nil
}

// validateInput - 入力バリデーション
func (u *ListReviewsUseCase) validateInput(input ListReviewsInput) error {
	// 言語のバリデーション
	if input.Language != "" {
		validLanguages := []string{
			"TypeScript", "JavaScript", "Python", "Go", "Java", "C++", "C#",
			"Ruby", "PHP", "Rust", "Swift", "Kotlin", "Other",
		}
		if !stringInSlice(input.Language, validLanguages) {
			return fmt.Errorf("サポートされていない言語です: %s", input.Language)
		}
	}

	// ステータスのバリデーション
	if input.Status != "" {
		validStatuses := []string{"success", "warning", "error", "pending"}
		if !stringInSlice(input.Status, validStatuses) {
			return fmt.Errorf("無効なステータスです: %s", input.Status)
		}
	}

	// ソート対象のバリデーション
	validSortBy := []string{"created_at", "language", "status"}
	if !stringInSlice(input.SortBy, validSortBy) {
		return fmt.Errorf("無効なソート対象です: %s", input.SortBy)
	}

	// ソート順のバリデーション
	if input.SortOrder != "asc" && input.SortOrder != "desc" {
		return fmt.Errorf("無効なソート順です: %s", input.SortOrder)
	}

	// 期間のバリデーション
	if input.DateFrom != nil && input.DateTo != nil {
		if input.DateFrom.After(*input.DateTo) {
			return fmt.Errorf("開始日は終了日より前である必要があります")
		}
	}

	return nil
}

// stringInSlice - スライスに要素が含まれているかチェック
func stringInSlice(item string, slice []string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
