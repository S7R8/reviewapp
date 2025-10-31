package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/s7r8/reviewapp/internal/application/usecase/review"
	"github.com/s7r8/reviewapp/internal/interfaces/http/middleware"
	"github.com/s7r8/reviewapp/internal/interfaces/http/response"
)

// ReviewHandler - レビューハンドラー
type ReviewHandler struct {
	reviewCodeUsecase     *review.ReviewCodeUseCase
	updateFeedbackUsecase *review.UpdateFeedbackUseCase
}

// NewReviewHandler - コンストラクタ
func NewReviewHandler(
	reviewCodeUsecase *review.ReviewCodeUseCase,
	updateFeedbackUsecase *review.UpdateFeedbackUseCase,
) *ReviewHandler {
	return &ReviewHandler{
		reviewCodeUsecase:     reviewCodeUsecase,
		updateFeedbackUsecase: updateFeedbackUsecase,
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
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	// 3. ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error:   "unauthorized",
			Message: "ユーザー情報が見つかりません。/auth/syncを先に呼び出してください。",
		})
	}

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
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error:   "internal_error",
			Message: "サーバーエラーが発生しました",
		})
	}

	// 構造化データを含めてレスポンス
	var structuredResult *StructuredReviewResult
	if output.Review.StructuredResult != nil {
		structuredResult = &StructuredReviewResult{
			Summary:    output.Review.StructuredResult.Summary,
			GoodPoints: output.Review.StructuredResult.GoodPoints,
			Improvements: func() []Improvement {
				improvements := make([]Improvement, len(output.Review.StructuredResult.Improvements))
				for i, imp := range output.Review.StructuredResult.Improvements {
					improvements[i] = Improvement{
						Title:       imp.Title,
						Description: imp.Description,
						CodeAfter:   imp.CodeAfter,
						Severity:    imp.Severity,
					}
				}
				return improvements
			}(),
		}
	}

	responseData := ReviewCodeResponse{
		ID:               output.Review.ID,
		UserID:           output.Review.UserID,
		Code:             output.Review.Code,
		Language:         output.Review.Language,
		Context:          output.Review.Context,
		ReviewResult:     output.Review.ReviewResult,
		StructuredResult: structuredResult,
		UsedKnowledgeIDs: output.Review.ReferencedKnowledge,
		LLMProvider:      output.Review.LLMProvider,
		LLMModel:         output.Review.LLMModel,
		TokensUsed:       output.Review.TokensUsed,
		CreatedAt:        output.Review.CreatedAt,
	}

	// ヘッダーにAPI Codeを追加
	c.Response().Header().Set("X-API-Code", "RV-001")
	return c.JSON(http.StatusCreated, responseData)
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
	ID               string                  `json:"id"`
	UserID           string                  `json:"user_id"`
	Code             string                  `json:"code"`
	Language         string                  `json:"language"`
	FileName         string                  `json:"file_name,omitempty"`
	Context          string                  `json:"context,omitempty"`
	ReviewResult     string                  `json:"review_result"`
	StructuredResult *StructuredReviewResult `json:"structured_result,omitempty"`
	UsedKnowledgeIDs []string                `json:"used_knowledge_ids"`
	LLMProvider      string                  `json:"llm_provider"`
	LLMModel         string                  `json:"llm_model"`
	TokensUsed       int                     `json:"tokens_used"`
	CreatedAt        time.Time               `json:"created_at"`
}

// StructuredReviewResult - 構造化されたレビュー結果（レスポンス用）
type StructuredReviewResult struct {
	Summary      string        `json:"summary"`
	GoodPoints   []string      `json:"good_points"`
	Improvements []Improvement `json:"improvements"`
}

// Improvement - 改善点（レスポンス用）
type Improvement struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CodeAfter   string `json:"code_after,omitempty"`
	Severity    string `json:"severity"`
}

// UpdateFeedback - PUT /api/v1/reviews/:id/feedback
func (h *ReviewHandler) UpdateFeedback(c echo.Context) error {
	// 1. パスパラメータからReviewIDを取得
	reviewID := c.Param("id")
	if reviewID == "" {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "validation_error",
			Message: "レビューIDは必須です",
		})
	}

	// 2. リクエストボディをパース
	var req UpdateFeedbackRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "invalid_request",
			Message: "リクエストボディが不正です",
		})
	}

	// 3. バリデーション
	if req.Score < 1 || req.Score > 3 {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "validation_error",
			Message: "スコアは1-3の整数で指定してください",
		})
	}

	if len(req.Comment) > 500 {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "validation_error",
			Message: "コメントは500文字以内にしてください",
		})
	}

	// 4. ユーザーIDを取得
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error:   "unauthorized",
			Message: "ユーザー情報が見つかりません。/auth/syncを先に呼び出してください。",
		})
	}

	// 5. UseCase実行
	input := review.UpdateFeedbackInput{
		ReviewID: reviewID,
		UserID:   userID,
		Score:    req.Score,
		Comment:  req.Comment,
	}

	output, err := h.updateFeedbackUsecase.Execute(c.Request().Context(), input)
	if err != nil {
		// エラーの種類を判定
		c.Logger().Errorf("UpdateFeedback failed: %v", err)
		
		// レビューが見つからない
		if err.Error() == "レビューが見つかりません" {
			return c.JSON(http.StatusNotFound, response.ErrorResponse{
				Error:   "not_found",
				Message: "レビューが見つかりません",
			})
		}
		
		// 権限エラー
		if err.Error() == "このレビューを更新する権限がありません" {
			return c.JSON(http.StatusForbidden, response.ErrorResponse{
				Error:   "forbidden",
				Message: "このレビューを更新する権限がありません",
			})
		}
		
		// その他のエラー
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error:   "internal_error",
			Message: "サーバーエラーが発生しました",
		})
	}

	// 6. 成功レスポンス
	responseData := UpdateFeedbackResponse{
		ID:              output.ReviewID,
		FeedbackScore:   output.FeedbackScore,
		FeedbackComment: output.FeedbackComment,
		UpdatedAt:       output.UpdatedAt,
	}

	// ヘッダーにAPI Codeを追加
	c.Response().Header().Set("X-API-Code", "RV-004")
	return c.JSON(http.StatusOK, responseData)
}

// UpdateFeedbackRequest - フィードバック更新リクエスト
type UpdateFeedbackRequest struct {
	Score   int    `json:"score" validate:"required,min=1,max=3"`
	Comment string `json:"comment" validate:"max=500"`
}

// UpdateFeedbackResponse - フィードバック更新レスポンス
type UpdateFeedbackResponse struct {
	ID              string `json:"id"`
	FeedbackScore   int    `json:"feedback_score"`
	FeedbackComment string `json:"feedback_comment,omitempty"`
	UpdatedAt       string `json:"updated_at"`
}
