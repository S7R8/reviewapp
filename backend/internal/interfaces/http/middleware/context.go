package middleware

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
)

// GetAuth0SubFromContext はコンテキストからAuth0のSubjectを取得します
func GetAuth0SubFromContext(c echo.Context) (string, error) {
	auth0Sub, ok := c.Get("auth0_sub").(string)
	if !ok || auth0Sub == "" {
		return "", fmt.Errorf("auth0_sub not found in context")
	}
	return auth0Sub, nil
}

// ContextKey はコンテキストに値を保存するためのキー型
type ContextKey string

const (
	// UserIDKey はユーザーIDを保存するためのキー
	UserIDKey ContextKey = "user_id"
	// Auth0SubKey はAuth0のSubjectを保存するためのキー
	Auth0SubKey ContextKey = "auth0_sub"
)

// SetUserID はコンテキストにユーザーIDを保存します
func SetUserID(c echo.Context, userID string) {
	c.Set(string(UserIDKey), userID)
}

// GetUserID はコンテキストからユーザーIDを取得します
func GetUserID(c echo.Context) (string, error) {
	userID, ok := c.Get(string(UserIDKey)).(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user_id not found in context")
	}
	return userID, nil
}

// GetUserIDFromStdContext は標準のcontextからユーザーIDを取得します
func GetUserIDFromStdContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user_id not found in context")
	}
	return userID, nil
}
