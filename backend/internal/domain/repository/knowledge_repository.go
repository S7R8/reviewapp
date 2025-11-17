package repository

import (
	"context"

	"github.com/s7r8/reviewapp/internal/domain/model"
)

// KnowledgeRepository - ナレッジリポジトリのインターフェース
type KnowledgeRepository interface {
	// Create - ナレッジを作成
	Create(ctx context.Context, knowledge *model.Knowledge) error

	// FindByID - IDでナレッジを取得
	FindByID(ctx context.Context, id string) (*model.Knowledge, error)

	// FindByUserID - ユーザーIDで全ナレッジを取得
	FindByUserID(ctx context.Context, userID string) ([]*model.Knowledge, error)

	// FindByCategory - カテゴリで検索
	FindByCategory(ctx context.Context, userID, category string) ([]*model.Knowledge, error)

	// Update - ナレッジを更新
	Update(ctx context.Context, knowledge *model.Knowledge) error

	// Delete - ナレッジを削除
	Delete(ctx context.Context, id string) error

	// SearchByKeyword - キーワードで検索（フルテキスト検索）
	// SearchByKeyword(ctx context.Context, userID, keyword string, limit int) ([]*model.Knowledge, error)

	// CountByUserID - ユーザーIDでナレッジ総数を取得（有効なもののみ）
	CountByUserID(ctx context.Context, userID string) (int, error)

	// CountByCategory - カテゴリ別のナレッジ数を取得
	CountByCategory(ctx context.Context, userID string) (map[string]int, error)
}
