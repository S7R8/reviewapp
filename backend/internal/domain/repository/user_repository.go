package repository

import (
	"context"

	"github.com/s7r8/reviewapp/internal/domain/model"
)

// UserRepository - ユーザーリポジトリのインターフェース
type UserRepository interface {
	// Create - ユーザーを作成
	Create(ctx context.Context, user *model.User) error

	// FindByID - IDでユーザーを取得
	FindByID(ctx context.Context, id string) (*model.User, error)

	// FindByEmail - メールアドレスでユーザーを取得
	FindByEmail(ctx context.Context, email string) (*model.User, error)

	// FindByAuth0UserID - Auth0ユーザーIDでユーザーを取得
	FindByAuth0UserID(ctx context.Context, auth0UserID string) (*model.User, error)

	// Update - ユーザーを更新
	Update(ctx context.Context, user *model.User) error

	// Delete - ユーザーを削除
	Delete(ctx context.Context, id string) error
}
