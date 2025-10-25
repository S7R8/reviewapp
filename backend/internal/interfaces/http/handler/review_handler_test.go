package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/s7r8/reviewapp/internal/application/usecase/review"
	"github.com/s7r8/reviewapp/internal/domain/service"
	"github.com/s7r8/reviewapp/internal/interfaces/http/handler"
	"github.com/s7r8/reviewapp/test/testutil"
)

func TestReviewHandler_ReviewCode(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
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
			expectedStatus: http.StatusCreated,
		},
		{
			name: "codeが空の場合",
			requestBody: map[string]string{
				"code":     "",
				"language": "go",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name: "languageが空の場合",
			requestBody: map[string]string{
				"code":     "func main() {}",
				"language": "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
		{
			name:           "不正なJSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "validation_error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックのUseCaseを準備
			mockKnowledgeRepo := testutil.NewMockKnowledgeRepository()
			mockReviewRepo := testutil.NewMockReviewRepository()
			mockClaudeClient := testutil.NewMockClaudeClient()
			reviewService := service.NewReviewService()

			reviewUseCase := review.NewReviewCodeUseCase(
				mockReviewRepo,
				mockKnowledgeRepo,
				reviewService,
				mockClaudeClient,
			)

			// ハンドラーを初期化
			h := handler.NewReviewHandler(reviewUseCase)

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

			// ハンドラーを実行
			err = h.ReviewCode(c)

			// ステータスコードを検証
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					if he.Code != tt.expectedStatus {
						t.Errorf("Expected status %d, got %d", tt.expectedStatus, he.Code)
					}
				} else {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				if rec.Code != tt.expectedStatus {
					t.Errorf("Expected status %d, got %d", tt.expectedStatus, rec.Code)
				}
			}

			// エラーレスポンスを検証
			if tt.expectedError != "" {
				var response map[string]string
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err == nil {
					if response["error"] != tt.expectedError {
						t.Errorf("Expected error %s, got %s", tt.expectedError, response["error"])
					}
				}
			}

			// 正常レスポンスを検証
			if tt.expectedStatus == http.StatusCreated {
				var response handler.ReviewCodeResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				if response.ID == "" {
					t.Errorf("Expected non-empty ID in response")
				}

				if response.ReviewResult == "" {
					t.Errorf("Expected non-empty ReviewResult in response")
				}

				// ヘッダーを検証
				apiCode := rec.Header().Get("X-API-Code")
				if apiCode != "RV-001" {
					t.Errorf("Expected X-API-Code to be RV-001, got %s", apiCode)
				}
			}
		})
	}
}
