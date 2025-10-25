package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/s7r8/reviewapp/internal/application/usecase/review"
)

// ReviewHandler - レビューハンドラー
type ReviewHandler struct {
	reviewCodeUsecase *review.ReviewCodeUseCase
}

// NewReviewHandler - コンストラクタ
func NewReviewHandler(reviewCodeUsecase *review.ReviewCodeUseCase) *ReviewHandler {
	return &ReviewHandler{
		reviewCodeUsecase: reviewCodeUsecase,
	}
}

// ReviewCode - POST /api/v1/reviews
func (h *ReviewHandler) ReviewCode(c echo.Context) error {
	var req ReviewCodeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   "validation_error",
			"message": "Invalid request body",
		})
	}

	// 2. バリデーション
	if err := validateReviewCodeRequest(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	// TODO: 認証からUserIDを取得（Phase 1は固定）
	userID := "00000000-0000-0000-0000-000000000001"

	// Usecaseを実行
	input := review.ReviewCodeInput{
		UserID:   userID,
		Code:     req.Code,
		Language: req.Language,
		Context:  req.Context,
	}

	output, err := h.reviewCodeUsecase.Execute(c.Request().Context(), input)
	if err != nil {
		// エラーの詳細をログに出力
		c.Logger().Errorf("ReviewCode failed: %v", err)

		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "サーバーエラーが発生しました",
		})
	}

	// レスポンスを構築
	response := ReviewCodeResponse{
		ID:               output.Review.ID,
		UserID:           output.Review.UserID,
		Code:             output.Review.Code,
		Language:         output.Review.Language,
		Context:          output.Review.Context,
		ReviewResult:     output.Review.ReviewResult,
		UsedKnowledgeIDs: output.Review.ReferencedKnowledge,
		LLMProvider:      output.Review.LLMProvider,
		LLMModel:         output.Review.LLMModel,
		TokensUsed:       output.Review.TokensUsed,
		CreatedAt:        output.Review.CreatedAt,
	}

	// ヘッダーにAPI Codeを追加
	c.Response().Header().Set("X-API-Code", "RV-001")

	return c.JSON(http.StatusCreated, response)
}

// バリデーション
func validateReviewCodeRequest(req *ReviewCodeRequest) error {
	if req.Code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "コードは必須です")
	}
	if req.Language == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "プログラミング言語は必須です")
	}
	// Context はオプション（空でもOK）
	return nil
}

// ReviewCodeRequest - リクエスト
type ReviewCodeRequest struct {
	Code     string `json:"code" validate:"required"`
	Language string `json:"language" validate:"required"`
	Context  string `json:"context"`
}

// ReviewCodeResponse - レスポンス
type ReviewCodeResponse struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	Code             string    `json:"code"`
	Language         string    `json:"language"`
	FileName         string    `json:"file_name,omitempty"`
	Context          string    `json:"context,omitempty"`
	ReviewResult     string    `json:"review_result"`
	UsedKnowledgeIDs []string  `json:"used_knowledge_ids"`
	LLMProvider      string    `json:"llm_provider"`
	LLMModel         string    `json:"llm_model"`
	TokensUsed       int       `json:"tokens_used"`
	CreatedAt        time.Time `json:"created_at"`
}
