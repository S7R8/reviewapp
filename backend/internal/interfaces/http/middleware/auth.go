package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/s7r8/reviewapp/internal/domain/repository"
	"github.com/s7r8/reviewapp/internal/infrastructure/auth"
)

// AuthMiddleware はJWT認証を行うミドルウェア
type AuthMiddleware struct {
	validator  *auth.Validator
	userRepo   repository.UserRepository
}

// NewAuthMiddleware は新しいAuthMiddlewareを作成します
func NewAuthMiddleware(validator *auth.Validator, userRepo repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		validator: validator,
		userRepo:  userRepo,
	}
}

// Authenticate はJWTトークンを検証し、ユーザー情報をコンテキストに保存します
func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Authorizationヘッダーを取得
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
		}

		// "Bearer "プレフィックスを確認
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
		}

		// トークンを抽出
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// トークンを検証
		token, err := m.validator.ValidateToken(context.Background(), tokenString)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token: "+err.Error())
		}

		// Auth0のSubject（ユーザーID）を取得
		auth0Sub := auth.GetAuth0Sub(token)
		if auth0Sub == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing subject in token")
		}

		// コンテキストに保存
		c.Set("auth0_sub", auth0Sub)
		c.Set("token", token)

		// ユーザーIDを取得してコンテキストに保存
		user, err := m.userRepo.FindByAuth0UserID(context.Background(), auth0Sub)
		if err == nil && user != nil {
			// ユーザーが見つかった場合、IDをコンテキストに保存
			c.Set(string(UserIDKey), user.ID)
		}
		// エラーの場合は無視（/auth/syncエンドポイントで作成される）

		return next(c)
	}
}

// GetAuth0Sub はコンテキストからAuth0 Subjectを取得します
func GetAuth0Sub(c echo.Context) string {
	if sub, ok := c.Get("auth0_sub").(string); ok {
		return sub
	}
	return ""
}
