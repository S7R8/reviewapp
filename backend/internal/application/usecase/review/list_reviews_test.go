package review

import (
	"context"
	"testing"
	"time"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/test/testutil"
	"github.com/stretchr/testify/assert"
)

func TestListReviewsUseCase_Execute(t *testing.T) {
	tests := []struct {
		name           string
		input          ListReviewsInput
		mockReviews    []*model.Review
		expectedCount  int
		expectedTotal  int
		expectedError  bool
	}{
		{
			name: "正常系: レビュー一覧取得",
			input: ListReviewsInput{
				UserID:   "test-user-id",
				Page:     1,
				PageSize: 10,
			},
			mockReviews: []*model.Review{
				{
					ID:           "review-1",
					UserID:       "test-user-id",
					Code:         "func test() {}",
					Language:     "go",
					ReviewResult: "Good code",
					CreatedAt:    time.Now(),
				},
				{
					ID:           "review-2",
					UserID:       "test-user-id",
					Code:         "func test2() {}",
					Language:     "go",
					ReviewResult: "Needs improvement",
					CreatedAt:    time.Now().Add(-1 * time.Hour),
				},
			},
			expectedCount: 2,
			expectedTotal: 2,
			expectedError: false,
		},
		{
			name: "正常系: 空の結果",
			input: ListReviewsInput{
				UserID:   "test-user-id",
				Page:     1,
				PageSize: 10,
			},
			mockReviews:   []*model.Review{},
			expectedCount: 0,
			expectedTotal: 0,
			expectedError: false,
		},
		{
			name: "正常系: 言語フィルター付き",
			input: ListReviewsInput{
				UserID:   "test-user-id",
				Page:     1,
				PageSize: 10,
				Language: "Go",
			},
			mockReviews: []*model.Review{
				{
					ID:           "review-1",
					UserID:       "test-user-id",
					Code:         "func test() {}",
					Language:     "go",
					ReviewResult: "Good code",
					CreatedAt:    time.Now(),
				},
			},
			expectedCount: 1,
			expectedTotal: 1,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックを準備
			mockRepo := testutil.NewMockReviewRepository()

			// テストデータを設定
			for _, review := range tt.mockReviews {
				mockRepo.Create(context.Background(), review)
			}

			// UseCaseを初期化
			uc := NewListReviewsUseCase(mockRepo)

			// 実行
			output, err := uc.Execute(context.Background(), tt.input)

			// 検証
			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.Equal(t, tt.expectedCount, len(output.Items))
			assert.Equal(t, tt.expectedTotal, output.Total)
		})
	}
}

func TestListReviewsUseCase_Execute_Pagination(t *testing.T) {
	mockRepo := testutil.NewMockReviewRepository()

	// 15件のレビューを作成
	for i := 1; i <= 15; i++ {
		review := &model.Review{
			ID:           "review-" + string(rune(i)),
			UserID:       "test-user-id",
			Code:         "func test() {}",
			Language:     "go",
			ReviewResult: "Test review",
			CreatedAt:    time.Now().Add(-time.Duration(i) * time.Hour),
		}
		mockRepo.Create(context.Background(), review)
	}

	uc := NewListReviewsUseCase(mockRepo)

	// Page 1
	output1, err := uc.Execute(context.Background(), ListReviewsInput{
		UserID:   "test-user-id",
		Page:     1,
		PageSize: 10,
	})

	assert.NoError(t, err)
	assert.Equal(t, 10, len(output1.Items))
	assert.Equal(t, 15, output1.Total)

	// Page 2
	output2, err := uc.Execute(context.Background(), ListReviewsInput{
		UserID:   "test-user-id",
		Page:     2,
		PageSize: 10,
	})

	assert.NoError(t, err)
	assert.Equal(t, 5, len(output2.Items))
	assert.Equal(t, 15, output2.Total)
}

func TestListReviewsUseCase_Execute_Sorting(t *testing.T) {
	mockRepo := testutil.NewMockReviewRepository()

	uc := NewListReviewsUseCase(mockRepo)

	tests := []struct {
		name      string
		sortBy    string
		sortOrder string
	}{
		{
			name:      "ソート: created_at DESC",
			sortBy:    "created_at",
			sortOrder: "desc",
		},
		{
			name:      "ソート: created_at ASC",
			sortBy:    "created_at",
			sortOrder: "asc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := uc.Execute(context.Background(), ListReviewsInput{
				UserID:    "test-user-id",
				Page:      1,
				PageSize:  10,
				SortBy:    tt.sortBy,
				SortOrder: tt.sortOrder,
			})

			assert.NoError(t, err)
			assert.NotNil(t, output)
		})
	}
}

func TestListReviewsUseCase_Execute_InvalidLanguage(t *testing.T) {
	mockRepo := testutil.NewMockReviewRepository()
	uc := NewListReviewsUseCase(mockRepo)

	_, err := uc.Execute(context.Background(), ListReviewsInput{
		UserID:   "test-user-id",
		Page:     1,
		PageSize: 10,
		Language: "InvalidLanguage",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "サポートされていない言語です")
}

func TestListReviewsUseCase_Execute_InvalidStatus(t *testing.T) {
	mockRepo := testutil.NewMockReviewRepository()
	uc := NewListReviewsUseCase(mockRepo)

	_, err := uc.Execute(context.Background(), ListReviewsInput{
		UserID:   "test-user-id",
		Page:     1,
		PageSize: 10,
		Status:   "invalid_status",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "無効なステータスです")
}

func TestListReviewsUseCase_Execute_InvalidDateRange(t *testing.T) {
	mockRepo := testutil.NewMockReviewRepository()
	uc := NewListReviewsUseCase(mockRepo)

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	_, err := uc.Execute(context.Background(), ListReviewsInput{
		UserID:   "test-user-id",
		Page:     1,
		PageSize: 10,
		DateFrom: &now,
		DateTo:   &yesterday,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "開始日は終了日より前である必要があります")
}

func TestListReviewsUseCase_Execute_DefaultValues(t *testing.T) {
	mockRepo := testutil.NewMockReviewRepository()
	uc := NewListReviewsUseCase(mockRepo)

	// PageとPageSizeを0で送信（デフォルト値が使われるはず）
	output, err := uc.Execute(context.Background(), ListReviewsInput{
		UserID:   "test-user-id",
		Page:     0,
		PageSize: 0,
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, output.Page)      // デフォルトは1
	assert.Equal(t, 10, output.PageSize) // デフォルトは10
}

func TestListReviewsUseCase_Execute_MaxPageSize(t *testing.T) {
	mockRepo := testutil.NewMockReviewRepository()
	uc := NewListReviewsUseCase(mockRepo)

	// PageSizeを200で送信（100に制限されるはず）
	output, err := uc.Execute(context.Background(), ListReviewsInput{
		UserID:   "test-user-id",
		Page:     1,
		PageSize: 200,
	})

	assert.NoError(t, err)
	assert.Equal(t, 100, output.PageSize) // 最大100
}
