package repository

import (
	"context"

	"github.com/s7r8/reviewapp/internal/domain/model"
)

// ReviewRepository - レビューリポジトリのインターフェース
type ReviewRepository interface {
	// Create - レビューを作成
	Create(ctx context.Context, review *model.Review) error

	// FindByID - IDでレビューを取得
	FindByID(ctx context.Context, id string) (*model.Review, error)

	// FindByUserID - ユーザーIDで全レビューを取得
	FindByUserID(ctx context.Context, userID string, limit int) ([]*model.Review, error)

	// Update - レビューを更新
	Update(ctx context.Context, review *model.Review) error

	// Delete - レビューを削除
	Delete(ctx context.Context, id string) error

	// FindRecentByUserID - 最近のレビューを取得
	FindRecentByUserID(ctx context.Context, userID string, limit int) ([]*model.Review, error)
}
