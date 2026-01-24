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
	"github.com/stretchr/testify/assert"
)

func TestReviewCodeUseCase_Execute(t *testing.T) {
	tests := []struct {
		name          string
		input         review.ReviewCodeInput
		knowledges    []*model.Knowledge
		claudeResponse *external.ReviewCodeOutput
		embeddingError error
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
		{
			name: "Embeddingエラーでもフォールバック",
			input: review.ReviewCodeInput{
				UserID:   "test-user-id",
				Code:     "func test() {}",
				Language: "go",
			},
			knowledges: []*model.Knowledge{},
			claudeResponse: &external.ReviewCodeOutput{
				ReviewResult: "レビュー結果",
				TokensUsed:   100,
			},
			embeddingError: errors.New("embedding generation failed"),
			expectedError:  false, // フォールバックするのでエラーにならない
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックを準備
			mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
			mockReviewRepo := testutil.NewMockReviewRepository()
			mockClaudeClient := testutil.NewMockClaudeClient()
			mockEmbeddingClient := testutil.NewMockEmbeddingClient()
			reviewService := service.NewReviewService()

			// テストデータを設定
			mockKnowledgeRepo.SetKnowledges(tt.knowledges)
			mockClaudeClient.SetResponse(tt.claudeResponse)

			if tt.embeddingError != nil {
				mockEmbeddingClient.SetError(tt.embeddingError)
			}

			// UseCaseを初期化
			uc := review.NewReviewCodeUseCase(
				mockReviewRepo,
				mockKnowledgeRepo,
				reviewService,
				mockClaudeClient,
				mockEmbeddingClient,
			)

			// 実行
			output, err := uc.Execute(context.Background(), tt.input)

			// 検証
			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.Equal(t, tt.input.UserID, output.Review.UserID)
			assert.Equal(t, tt.claudeResponse.ReviewResult, output.Review.ReviewResult)
			assert.Equal(t, tt.claudeResponse.TokensUsed, output.Review.TokensUsed)
		})
	}
}

func TestReviewCodeUseCase_Execute_Error_Cases(t *testing.T) {
	t.Run("ナレッジリポジトリエラー", func(t *testing.T) {
		mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
		mockReviewRepo := testutil.NewMockReviewRepository()
		mockClaudeClient := testutil.NewMockClaudeClient()
		mockEmbeddingClient := testutil.NewMockEmbeddingClient()
		reviewService := service.NewReviewService()

		mockKnowledgeRepo.SetError(errors.New("database error"))

		uc := review.NewReviewCodeUseCase(
			mockReviewRepo,
			mockKnowledgeRepo,
			reviewService,
			mockClaudeClient,
			mockEmbeddingClient,
		)

		input := review.ReviewCodeInput{
			UserID:   "test-user-id",
			Code:     "test code",
			Language: "go",
		}

		_, err := uc.Execute(context.Background(), input)
		assert.Error(t, err)
	})

	t.Run("Claude APIエラー", func(t *testing.T) {
		mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
		mockReviewRepo := testutil.NewMockReviewRepository()
		mockClaudeClient := testutil.NewMockClaudeClient()
		mockEmbeddingClient := testutil.NewMockEmbeddingClient()
		reviewService := service.NewReviewService()

		mockClaudeClient.SetError(errors.New("API error"))

		uc := review.NewReviewCodeUseCase(
			mockReviewRepo,
			mockKnowledgeRepo,
			reviewService,
			mockClaudeClient,
			mockEmbeddingClient,
		)

		input := review.ReviewCodeInput{
			UserID:   "test-user-id",
			Code:     "test code",
			Language: "go",
		}

		_, err := uc.Execute(context.Background(), input)
		assert.Error(t, err)
	})

	t.Run("レビュー保存エラー", func(t *testing.T) {
		mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
		mockReviewRepo := testutil.NewMockReviewRepository()
		mockClaudeClient := testutil.NewMockClaudeClient()
		mockEmbeddingClient := testutil.NewMockEmbeddingClient()
		reviewService := service.NewReviewService()

		mockClaudeClient.SetResponse(&external.ReviewCodeOutput{
			ReviewResult: "Good code",
			TokensUsed:   100,
		})
		mockReviewRepo.SetError(errors.New("save error"))

		uc := review.NewReviewCodeUseCase(
			mockReviewRepo,
			mockKnowledgeRepo,
			reviewService,
			mockClaudeClient,
			mockEmbeddingClient,
		)

		input := review.ReviewCodeInput{
			UserID:   "test-user-id",
			Code:     "test code",
			Language: "go",
		}

		_, err := uc.Execute(context.Background(), input)
		assert.Error(t, err)
	})
}

func TestReviewCodeUseCase_Execute_EmptyInput(t *testing.T) {
	mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
	mockReviewRepo := testutil.NewMockReviewRepository()
	mockClaudeClient := testutil.NewMockClaudeClient()
	mockEmbeddingClient := testutil.NewMockEmbeddingClient()
	reviewService := service.NewReviewService()

	uc := review.NewReviewCodeUseCase(
		mockReviewRepo,
		mockKnowledgeRepo,
		reviewService,
		mockClaudeClient,
		mockEmbeddingClient,
	)

	tests := []struct {
		name  string
		input review.ReviewCodeInput
	}{
		{
			name: "UserIDが空",
			input: review.ReviewCodeInput{
				UserID:   "",
				Code:     "func test() {}",
				Language: "go",
			},
		},
		{
			name: "Codeが空",
			input: review.ReviewCodeInput{
				UserID:   "test-user-id",
				Code:     "",
				Language: "go",
			},
		},
		{
			name: "Languageが空",
			input: review.ReviewCodeInput{
				UserID:   "test-user-id",
				Code:     "func test() {}",
				Language: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.Execute(context.Background(), tt.input)
			assert.Error(t, err)
		})
	}
}

func TestReviewCodeUseCase_Execute_WithMultipleKnowledge(t *testing.T) {
	mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
	mockReviewRepo := testutil.NewMockReviewRepository()
	mockClaudeClient := testutil.NewMockClaudeClient()
	mockEmbeddingClient := testutil.NewMockEmbeddingClient()
	reviewService := service.NewReviewService()

	// 複数のナレッジを設定
	knowledges := []*model.Knowledge{
		{
			ID:       "knowledge-1",
			UserID:   "test-user-id",
			Title:    "Rule 1",
			Content:  "Content 1",
			Category: "clean_code",
			Priority: 5,
		},
		{
			ID:       "knowledge-2",
			UserID:   "test-user-id",
			Title:    "Rule 2",
			Content:  "Content 2",
			Category: "security",
			Priority: 4,
		},
		{
			ID:       "knowledge-3",
			UserID:   "test-user-id",
			Title:    "Rule 3",
			Content:  "Content 3",
			Category: "performance",
			Priority: 3,
		},
	}

	mockKnowledgeRepo.SetKnowledges(knowledges)
	mockClaudeClient.SetResponse(&external.ReviewCodeOutput{
		ReviewResult: "Review based on multiple knowledge",
		TokensUsed:   200,
	})

	uc := review.NewReviewCodeUseCase(
		mockReviewRepo,
		mockKnowledgeRepo,
		reviewService,
		mockClaudeClient,
		mockEmbeddingClient,
	)

	input := review.ReviewCodeInput{
		UserID:   "test-user-id",
		Code:     "func test() {}",
		Language: "go",
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	// 複数のナレッジが参照されていることを確認
	assert.NotEmpty(t, output.Review.ReferencedKnowledge)
}

func TestReviewCodeUseCase_Execute_LongCode(t *testing.T) {
	mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
	mockReviewRepo := testutil.NewMockReviewRepository()
	mockClaudeClient := testutil.NewMockClaudeClient()
	mockEmbeddingClient := testutil.NewMockEmbeddingClient()
	reviewService := service.NewReviewService()

	mockClaudeClient.SetResponse(&external.ReviewCodeOutput{
		ReviewResult: "Long code review",
		TokensUsed:   500,
	})

	uc := review.NewReviewCodeUseCase(
		mockReviewRepo,
		mockKnowledgeRepo,
		reviewService,
		mockClaudeClient,
		mockEmbeddingClient,
	)

	// 長いコード（1000行以上を想定）
	longCode := ""
	for i := 0; i < 1000; i++ {
		longCode += "func test" + string(rune(i)) + "() {}\n"
	}

	input := review.ReviewCodeInput{
		UserID:   "test-user-id",
		Code:     longCode,
		Language: "go",
		Context:  "Long code test",
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, longCode, output.Review.Code)
}

func TestReviewCodeUseCase_Execute_DifferentLanguages(t *testing.T) {
	languages := []string{"go", "python", "javascript", "typescript", "java", "rust"}

	for _, lang := range languages {
		t.Run("Language: "+lang, func(t *testing.T) {
			mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
			mockReviewRepo := testutil.NewMockReviewRepository()
			mockClaudeClient := testutil.NewMockClaudeClient()
			mockEmbeddingClient := testutil.NewMockEmbeddingClient()
			reviewService := service.NewReviewService()

			mockClaudeClient.SetResponse(&external.ReviewCodeOutput{
				ReviewResult: "Review for " + lang,
				TokensUsed:   100,
			})

			uc := review.NewReviewCodeUseCase(
				mockReviewRepo,
				mockKnowledgeRepo,
				reviewService,
				mockClaudeClient,
				mockEmbeddingClient,
			)

			input := review.ReviewCodeInput{
				UserID:   "test-user-id",
				Code:     "test code",
				Language: lang,
			}

			output, err := uc.Execute(context.Background(), input)

			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.Equal(t, lang, output.Review.Language)
		})
	}
}
