package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/s7r8/reviewapp/internal/application/usecase/review"
	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/domain/service"
	"github.com/s7r8/reviewapp/internal/infrastructure/external"
	"github.com/s7r8/reviewapp/internal/interfaces/http/handler"
	"github.com/s7r8/reviewapp/internal/interfaces/http/middleware"
	"github.com/s7r8/reviewapp/test/testutil"
	"github.com/stretchr/testify/assert"
)

func TestReviewHandler_ReviewCode(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setUserID      bool
		claudeError    error
		expectedStatus int
		expectedError  string
	}{
		{
			name: "正常なリクエスト",
			requestBody: map[string]string{
				"code":     "func main() { fmt.Println(\"Hello\") }",
				"language": "go",
				"context":  "テスト用コード",
			},
			setUserID:      true,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "codeが空の場合",
			requestBody: map[string]string{
				"code":     "",
				"language": "go",
			},
			setUserID:      true,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name: "languageが空の場合",
			requestBody: map[string]string{
				"code":     "func main() {}",
				"language": "",
			},
			setUserID:      true,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name:           "不正なJSON",
			requestBody:    "invalid json",
			setUserID:      true,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name: "認証情報なし",
			requestBody: map[string]string{
				"code":     "func main() {}",
				"language": "go",
			},
			setUserID:      false,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "unauthorized",
		},
		{
			name: "Claude APIエラー",
			requestBody: map[string]string{
				"code":     "func main() {}",
				"language": "go",
			},
			setUserID:      true,
			claudeError:    errors.New("API error"),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "internal_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックのUseCaseを準備
			mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
			mockReviewRepo := testutil.NewMockReviewRepository()
			mockClaudeClient := testutil.NewMockClaudeClient()
			mockEmbeddingClient := testutil.NewMockEmbeddingClient()
			reviewService := service.NewReviewService()

			// Claudeエラーを設定
			if tt.claudeError != nil {
				mockClaudeClient.SetError(tt.claudeError)
			}

			reviewUseCase := review.NewReviewCodeUseCase(
				mockReviewRepo,
				mockKnowledgeRepo,
				reviewService,
				mockClaudeClient,
				mockEmbeddingClient,
			)

			feedbackUseCase := review.NewUpdateFeedbackUseCase(
				mockReviewRepo,
			)

			listReviewsUseCase := review.NewListReviewsUseCase(
				mockReviewRepo,
			)

			getReviewUseCase := review.NewGetReviewUseCase(
				mockReviewRepo,
			)

			// ハンドラーを初期化
			h := handler.NewReviewHandler(reviewUseCase, feedbackUseCase, listReviewsUseCase, getReviewUseCase)

			// リクエストを準備
			var reqBody []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/review", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Echoコンテキストを作成
			e := echo.New()
			c := e.NewContext(req, rec)

			// 認証情報を設定（テスト用）
			if tt.setUserID {
				middleware.SetUserID(c, "test-user-id")
			}

			// ハンドラーを実行
			err = h.ReviewCode(c)

			// ステータスコードを検証
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tt.expectedStatus, he.Code)
				} else {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			// エラーレスポンスを検証
			if tt.expectedError != "" {
				var response map[string]string
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err == nil {
					assert.Equal(t, tt.expectedError, response["error"])
				}
			}

			// 正常レスポンスを検証
			if tt.expectedStatus == http.StatusCreated {
				var response handler.ReviewCodeResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.ID)
				assert.NotEmpty(t, response.ReviewResult)
				assert.Equal(t, "RV-001", rec.Header().Get("X-API-Code"))
			}
		})
	}
}

func TestReviewHandler_ReviewCode_WithStructuredResult(t *testing.T) {
	mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
	mockReviewRepo := testutil.NewMockReviewRepository()
	mockClaudeClient := testutil.NewMockClaudeClient()
	mockEmbeddingClient := testutil.NewMockEmbeddingClient()
	reviewService := service.NewReviewService()

	// 構造化結果を含むレスポンスを設定
	mockClaudeClient.SetResponse(&external.ReviewCodeOutput{
		ReviewResult: "Good code with structured feedback",
		TokensUsed:   200,
	})

	reviewUseCase := review.NewReviewCodeUseCase(
		mockReviewRepo,
		mockKnowledgeRepo,
		reviewService,
		mockClaudeClient,
		mockEmbeddingClient,
	)

	h := handler.NewReviewHandler(reviewUseCase, nil, nil, nil)

	reqBody, _ := json.Marshal(map[string]string{
		"code":     "func test() {}",
		"language": "go",
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/review", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	e := echo.New()
	c := e.NewContext(req, rec)
	middleware.SetUserID(c, "test-user-id")

	err := h.ReviewCode(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestReviewHandler_UpdateFeedback(t *testing.T) {
	tests := []struct {
		name           string
		reviewID       string
		requestBody    interface{}
		setUserID      bool
		mockReview     *model.Review
		expectedStatus int
		expectedError  string
	}{
		{
			name:     "正常なリクエスト（Good）",
			reviewID: "test-review-id",
			requestBody: map[string]interface{}{
				"score":   3,
				"comment": "とても役に立ちました",
			},
			setUserID: true,
			mockReview: &model.Review{
				ID:           "test-review-id",
				UserID:       "test-user-id",
				Code:         "func test() {}",
				Language:     "go",
				ReviewResult: "Good",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "scoreが0",
			reviewID: "test-review-id",
			requestBody: map[string]interface{}{
				"score": 0,
			},
			setUserID:      true,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name:     "scoreが4",
			reviewID: "test-review-id",
			requestBody: map[string]interface{}{
				"score": 4,
			},
			setUserID:      true,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name:     "reviewIDが空",
			reviewID: "",
			requestBody: map[string]interface{}{
				"score": 3,
			},
			setUserID:      true,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name:           "不正なJSON",
			reviewID:       "test-review-id",
			requestBody:    "invalid json",
			setUserID:      true,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid_request",
		},
		{
			name:     "認証情報なし",
			reviewID: "test-review-id",
			requestBody: map[string]interface{}{
				"score": 3,
			},
			setUserID:      false,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "unauthorized",
		},
		{
			name:     "レビューが見つからない",
			reviewID: "non-existent-id",
			requestBody: map[string]interface{}{
				"score": 3,
			},
			setUserID:      true,
			mockReview:     nil,
			expectedStatus: http.StatusNotFound,
			expectedError:  "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReviewRepo := testutil.NewMockReviewRepository()

			if tt.mockReview != nil {
				mockReviewRepo.Create(nil, tt.mockReview)
			}

			feedbackUseCase := review.NewUpdateFeedbackUseCase(mockReviewRepo)
			h := handler.NewReviewHandler(nil, feedbackUseCase, nil, nil)

			var reqBody []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPut, "/api/v1/reviews/"+tt.reviewID+"/feedback", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.reviewID)

			if tt.setUserID {
				middleware.SetUserID(c, "test-user-id")
			}

			err = h.UpdateFeedback(c)

			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tt.expectedStatus, he.Code)
				}
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			if tt.expectedError != "" {
				var response map[string]string
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err == nil {
					assert.Equal(t, tt.expectedError, response["error"])
				}
			}

			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, "RV-004", rec.Header().Get("X-API-Code"))
			}
		})
	}
}

func TestReviewHandler_ListReviews(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		setUserID      bool
		mockReviews    []*model.Review
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "正常系: レビュー一覧取得",
			queryParams: map[string]string{
				"page":      "1",
				"page_size": "10",
			},
			setUserID: true,
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
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name: "正常系: 空の結果",
			queryParams: map[string]string{
				"page":      "1",
				"page_size": "10",
			},
			setUserID:      true,
			mockReviews:    []*model.Review{},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "認証情報なし",
			queryParams: map[string]string{
				"page":      "1",
				"page_size": "10",
			},
			setUserID:      false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "不正なページ番号",
			queryParams: map[string]string{
				"page":      "-1",
				"page_size": "10",
			},
			setUserID:      true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "不正なページサイズ",
			queryParams: map[string]string{
				"page":      "1",
				"page_size": "200",
			},
			setUserID:      true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReviewRepo := testutil.NewMockReviewRepository()
			for _, r := range tt.mockReviews {
				mockReviewRepo.Create(nil, r)
			}

			listReviewsUseCase := review.NewListReviewsUseCase(mockReviewRepo)
			h := handler.NewReviewHandler(nil, nil, listReviewsUseCase, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/reviews", nil)
			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			if tt.setUserID {
				middleware.SetUserID(c, "test-user-id")
			}

			err := h.ListReviews(c)

			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tt.expectedStatus, he.Code)
				}
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response handler.ListReviewsResponse
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(response.Items))
				assert.Equal(t, "RV-002", rec.Header().Get("X-API-Code"))
			}
		})
	}
}

func TestReviewHandler_GetReviewByID(t *testing.T) {
	tests := []struct {
		name           string
		reviewID       string
		userID         string
		setUserID      bool
		mockReview     *model.Review
		expectedStatus int
	}{
		{
			name:      "正常系: レビュー取得成功",
			reviewID:  "review-123",
			userID:    "test-user-id",
			setUserID: true,
			mockReview: &model.Review{
				ID:           "review-123",
				UserID:       "test-user-id",
				Code:         "func test() {}",
				Language:     "go",
				ReviewResult: "Good code",
				CreatedAt:    time.Now(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "異常系: レビューが見つからない",
			reviewID:       "non-existent",
			userID:         "test-user-id",
			setUserID:      true,
			mockReview:     nil,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:      "異常系: 他人のレビューにアクセス",
			reviewID:  "review-123",
			userID:    "other-user-id",
			setUserID: true,
			mockReview: &model.Review{
				ID:           "review-123",
				UserID:       "test-user-id",
				Code:         "func test() {}",
				Language:     "go",
				ReviewResult: "Good code",
				CreatedAt:    time.Now(),
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "異常系: 認証情報なし",
			reviewID:       "review-123",
			userID:         "",
			setUserID:      false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "異常系: reviewIDが空",
			reviewID:       "",
			userID:         "test-user-id",
			setUserID:      true,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReviewRepo := testutil.NewMockReviewRepository()
			if tt.mockReview != nil {
				mockReviewRepo.Create(nil, tt.mockReview)
			}

			getReviewUseCase := review.NewGetReviewUseCase(mockReviewRepo)
			h := handler.NewReviewHandler(nil, nil, nil, getReviewUseCase)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/reviews/"+tt.reviewID, nil)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.reviewID)

			if tt.setUserID {
				middleware.SetUserID(c, tt.userID)
			}

			err := h.GetReviewByID(c)

			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tt.expectedStatus, he.Code)
				}
			} else {
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, "RV-003", rec.Header().Get("X-API-Code"))
			}
		})
	}
}
