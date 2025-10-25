package review_test

import (
	"context"
	"errors"
	"testing"

	"github.com/s7r8/reviewapp/internal/application/usecase/review"
	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/domain/service"
	"github.com/s7r8/reviewapp/internal/infrastructure/external"
	"github.com/s7r8/reviewapp/test/testutil"
)

func TestReviewCodeUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		input         review.ReviewCodeInput
		knowledges    []*model.Knowledge
		claudeResponse *external.ReviewCodeOutput
		expectedError bool
	}{
		{
			name: "正常なコードレビュー",
			input: review.ReviewCodeInput{
				UserID:   "test-user-id",
				Code:     "func main() { fmt.Println(\"Hello\") }",
				Language: "go",
				Context:  "テスト用コード",
			},
			knowledges: []*model.Knowledge{
				{
					ID:       "knowledge-1",
					UserID:   "test-user-id",
					Title:    "関数は1つのことだけをする",
					Content:  "関数は50行以内に抑える",
					Category: "clean_code",
					Priority: 5,
				},
			},
			claudeResponse: &external.ReviewCodeOutput{
				ReviewResult: "コードは適切に書かれています。",
				TokensUsed:   150,
			},
			expectedError: false,
		},
		{
			name: "ナレッジが存在しない場合",
			input: review.ReviewCodeInput{
				UserID:   "test-user-id",
				Code:     "func main() {}",
				Language: "go",
			},
			knowledges: []*model.Knowledge{},
			claudeResponse: &external.ReviewCodeOutput{
				ReviewResult: "一般的なベストプラクティスに基づいてレビューしました。",
				TokensUsed:   100,
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックを準備
			mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
			mockReviewRepo := testutil.NewMockReviewRepository()
			mockClaudeClient := testutil.NewMockClaudeClient()
			reviewService := service.NewReviewService()

			// テストデータを設定
			mockKnowledgeRepo.SetKnowledges(tt.knowledges)
			mockClaudeClient.SetResponse(tt.claudeResponse)

			// UseCaseを初期化
			uc := review.NewReviewCodeUseCase(
				mockReviewRepo,
				mockKnowledgeRepo,
				reviewService,
				mockClaudeClient,
			)

			// 実行
			output, err := uc.Execute(context.Background(), tt.input)

			// 検証
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if output == nil {
				t.Errorf("Expected output, but got nil")
				return
			}

			if output.Review.UserID != tt.input.UserID {
				t.Errorf("Expected UserID %s, got %s", tt.input.UserID, output.Review.UserID)
			}

			if output.Review.ReviewResult != tt.claudeResponse.ReviewResult {
				t.Errorf("Expected ReviewResult %s, got %s", tt.claudeResponse.ReviewResult, output.Review.ReviewResult)
			}

			if output.Review.TokensUsed != tt.claudeResponse.TokensUsed {
				t.Errorf("Expected TokensUsed %d, got %d", tt.claudeResponse.TokensUsed, output.Review.TokensUsed)
			}
		})
	}
}

func TestReviewCodeUseCase_Execute_Error_Cases(t *testing.T) {
	t.Run("ナレッジリポジトリエラー", func(t *testing.T) {
		mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
		mockReviewRepo := testutil.NewMockReviewRepository()
		mockClaudeClient := testutil.NewMockClaudeClient()
		reviewService := service.NewReviewService()

		mockKnowledgeRepo.SetError(errors.New("database error"))

		uc := review.NewReviewCodeUseCase(
			mockReviewRepo,
			mockKnowledgeRepo,
			reviewService,
			mockClaudeClient,
		)

		input := review.ReviewCodeInput{
			UserID:   "test-user-id",
			Code:     "test code",
			Language: "go",
		}

		_, err := uc.Execute(context.Background(), input)
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	})

	t.Run("Claude APIエラー", func(t *testing.T) {
		mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
		mockReviewRepo := testutil.NewMockReviewRepository()
		mockClaudeClient := testutil.NewMockClaudeClient()
		reviewService := service.NewReviewService()

		mockClaudeClient.SetError(errors.New("API error"))

		uc := review.NewReviewCodeUseCase(
			mockReviewRepo,
			mockKnowledgeRepo,
			reviewService,
			mockClaudeClient,
		)

		input := review.ReviewCodeInput{
			UserID:   "test-user-id",
			Code:     "test code",
			Language: "go",
		}

		_, err := uc.Execute(context.Background(), input)
		if err == nil {
			t.Errorf("Expected error, but got nil")
		}
	})
}
