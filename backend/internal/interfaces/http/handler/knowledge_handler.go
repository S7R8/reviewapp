package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/s7r8/reviewapp/internal/application/usecase/knowledge"
	"github.com/s7r8/reviewapp/internal/interfaces/http/response"
)

// KnowledgeHandler - ナレッジ関連のHTTPハンドラ
type KnowledgeHandler struct {
	createKnowledgeUC *knowledge.CreateKnowledgeUseCase
	listKnowledgeUC   *knowledge.ListKnowledgeUseCase
}

// NewKnowledgeHandler - コンストラクタ
func NewKnowledgeHandler(
	createKnowledgeUC *knowledge.CreateKnowledgeUseCase,
	listKnowledgeUC *knowledge.ListKnowledgeUseCase,
) *KnowledgeHandler {
	return &KnowledgeHandler{
		createKnowledgeUC: createKnowledgeUC,
		listKnowledgeUC:   listKnowledgeUC,
	}
}

// CreateKnowledgeRequest - リクエストボディ
type CreateKnowledgeRequest struct {
	Title    string `json:"title" validate:"required,max=200"`
	Content  string `json:"content" validate:"required"`
	Category string `json:"category" validate:"required"`
	Priority int    `json:"priority" validate:"required,min=1,max=5"`
}

// CreateKnowledge - ナレッジ作成エンドポイント
// POST /api/v1/knowledge
func (h *KnowledgeHandler) CreateKnowledge(c echo.Context) error {
	// 1. リクエストボディをパース
	var req CreateKnowledgeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "invalid_request",
			Message: "リクエストボディが不正です",
		})
	}

	// 2. バリデーション
	if err := validateCreateKnowledgeRequest(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	// 3. ユーザーIDを取得（TODO: JWTから取得、現在は固定）
	// 本来はミドルウェアでJWTを検証し、c.Get("user_id")で取得
	userID := "00000000-0000-0000-0000-000000000001" // 開発用固定値

	// 4. UseCase実行
	output, err := h.createKnowledgeUC.Execute(c.Request().Context(), knowledge.CreateKnowledgeInput{
		UserID:   userID,
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Priority: req.Priority,
	})
	if err != nil {
		// ドメインエラー（バリデーションエラー）
		return c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
	}

	// 5. レスポンスヘッダーに API Code を追加
	c.Response().Header().Set("X-API-Code", "KN-001")

	// 6. 成功レスポンス
	return c.JSON(http.StatusCreated, output.Knowledge)
}

// validateCreateKnowledgeRequest - リクエストバリデーション
func validateCreateKnowledgeRequest(req *CreateKnowledgeRequest) error {
	if req.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "タイトルは必須です")
	}
	if len(req.Title) > 200 {
		return echo.NewHTTPError(http.StatusBadRequest, "タイトルは200文字以内にしてください")
	}
	if req.Content == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "内容は必須です")
	}
	if req.Category == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "カテゴリは必須です")
	}
	if req.Priority < 1 || req.Priority > 5 {
		return echo.NewHTTPError(http.StatusBadRequest, "重要度は1-5の整数で指定してください")
	}
	return nil
}

func (h *KnowledgeHandler) ListKnowledge(c echo.Context) error {
	// 1．クエリパラメータを取得
	category := c.QueryParam("category")

	// 2. カテゴリのバリデーション
	if category != "" {
		if err := validateCategory(category); err != nil {
			return c.JSON(http.StatusBadRequest, response.ErrorResponse{
				Error:   "validation_error",
				Message: err.Error(),
			})
		}
	}

	// 3. ユーザーIDを取得（TODO: JWTから取得、現在は固定）
	userID := "00000000-0000-0000-0000-000000000001" // 開発用固定値

	// 4. UseCase実行
	output, err := h.listKnowledgeUC.Execute(c.Request().Context(), knowledge.ListKnowledgeInput{
		UserID:   userID,
		Category: category,
	})
	if err != nil {
		// DBエラーなど
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error:   "internal_error",
			Message: "サーバーエラーが発生しました",
		})
	}

	// 5. レスポンスヘッダーに API Code を追加
	c.Response().Header().Set("X-API-Code", "KN-002")

	// 6. 成功レスポンス（シンプルな配列）
	return c.JSON(http.StatusOK, output.Knowledges)
}

// validateCategory - カテゴリのバリデーション
func validateCategory(category string) error {
	validCategories := map[string]bool{
		"error_handling": true,
		"testing":        true,
		"performance":    true,
		"security":       true,
		"clean_code":     true,
		"architecture":   true,
		"other":          true,
	}

	if !validCategories[category] {
		return echo.NewHTTPError(http.StatusBadRequest, "無効なカテゴリです")
	}
	return nil
}
