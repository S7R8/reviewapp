package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/s7r8/reviewapp/internal/application/usecase/user"
	"github.com/s7r8/reviewapp/internal/interfaces/http/middleware"
	"github.com/s7r8/reviewapp/internal/interfaces/http/response"
)

// AuthHandler - 認証関連のHTTPハンドラ
type AuthHandler struct {
	userUC *user.UseCase
}

// NewAuthHandler - コンストラクタ
func NewAuthHandler(userUC *user.UseCase) *AuthHandler {
	return &AuthHandler{
		userUC: userUC,
	}
}

// SyncUser - ユーザー同期エンドポイント
// POST /api/v1/auth/sync
// Auth0でログイン後、初回アクセス時にバックエンドのユーザーと同期する
func (h *AuthHandler) SyncUser(c echo.Context) error {
	// 1. Auth0のSubjectを取得
	auth0Sub, err := middleware.GetAuth0SubFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error:   "unauthorized",
			Message: "認証情報が見つかりません",
		})
	}

	// 2. UseCase実行（ユーザーを取得または作成）
	user, err := h.userUC.SyncUserByAuth0Sub(c.Request().Context(), auth0Sub)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error:   "sync_failed",
			Message: "ユーザー同期に失敗しました",
		})
	}

	// 3. ユーザーIDをコンテキストに保存（後続のリクエストで使用）
	middleware.SetUserID(c, user.ID)

	// 4. 成功レスポンス
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user": user,
	})
}
