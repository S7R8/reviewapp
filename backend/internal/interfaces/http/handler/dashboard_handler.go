package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/s7r8/reviewapp/internal/application/usecase/dashboard"
	"github.com/s7r8/reviewapp/internal/interfaces/http/response"
)

// DashboardHandler - ダッシュボードハンドラー
type DashboardHandler struct {
	getStatsUseCase *dashboard.GetStatsUseCase
}

// NewDashboardHandler - コンストラクタ
func NewDashboardHandler(
	getStatsUseCase *dashboard.GetStatsUseCase,
) *DashboardHandler {
	return &DashboardHandler{
		getStatsUseCase: getStatsUseCase,
	}
}

// GetStats - ダッシュボード統計取得
// @Summary ダッシュボード統計取得
// @Description ダッシュボード画面に表示する統計情報を一括で取得
// @Tags dashboard
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} dashboard.DashboardStatsResponse
// @Failure 401 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/dashboard/stats [get]
func (h *DashboardHandler) GetStats(c echo.Context) error {
	// JWTからユーザーIDを取得
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証が必要です",
		})
	}

	// UseCaseを実行
	stats, err := h.getStatsUseCase.Execute(c.Request().Context(), userID.(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error:   "internal_error",
			Message: "統計情報の取得に失敗しました",
		})
	}

	// レスポンスヘッダーにAPI Codeを追加
	c.Response().Header().Set("X-API-Code", "DS-001")

	// 成功レスポンス
	return c.JSON(http.StatusOK, stats)
}
