package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/s7r8/reviewapp/internal/infrastructure/config"
)

// DB - データベース接続のラッパー
type DB struct {
	*sql.DB
}

// NewDB - データベース接続を確立
func NewDB(cfg *config.DatabaseConfig) (*DB, error) {
	dsn := cfg.GetDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 接続プールの設定
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 接続確認（タイムアウト付き）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("✅ Database connected successfully")

	return &DB{db}, nil
}

// Close - データベース接続を閉じる
func (db *DB) Close() error {
	if err := db.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	fmt.Println("🔌 Database connection closed")
	return nil
}

// HealthCheck - ヘルスチェック
func (db *DB) HealthCheck(ctx context.Context) error {
	return db.PingContext(ctx)
}
