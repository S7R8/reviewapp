package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/s7r8/reviewapp/internal/infrastructure/config"
	"github.com/s7r8/reviewapp/internal/infrastructure/persistence/postgres"
)

// TestDatabase - テスト用データベース
type TestDatabase struct {
	DB     *sql.DB
	Config *config.DatabaseConfig
}

// NewTestDatabase - テスト用データベースを作成
func NewTestDatabase(t *testing.T) *TestDatabase {
	t.Helper()

	// テスト用の設定（環境変数で上書き可能）
	cfg := &config.DatabaseConfig{
		Host:     getEnvOrDefault("TEST_DB_HOST", "postgres"),
		Port:     getEnvOrDefault("TEST_DB_PORT", "5432"),
		User:     getEnvOrDefault("TEST_DB_USER", "dev_user"),
		Password: getEnvOrDefault("TEST_DB_PASSWORD", "dev_password"),
		Name:     getEnvOrDefault("TEST_DB_NAME", "reviewapp"),
		SSLMode:  getEnvOrDefault("TEST_DB_SSLMODE", "disable"),
	}

	// テスト用データベースに接続
	db, err := postgres.NewDB(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v\nConfig: %+v", err, cfg)
	}

	return &TestDatabase{
		DB:     db.DB,
		Config: cfg,
	}
}

// Close - データベース接続を閉じる
func (td *TestDatabase) Close() {
	if td.DB != nil {
		td.DB.Close()
	}
}

// Cleanup - テストデータをクリーンアップ
func (td *TestDatabase) Cleanup(t *testing.T) {
	t.Helper()

	tables := []string{
		"review_knowledge",
		"reviews",
		"knowledge_tags",
		"tags",
		"knowledge",
		"users",
	}

	for _, table := range tables {
		_, err := td.DB.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			t.Logf("Failed to cleanup table %s: %v", table, err)
		}
	}
}

// SeedTestData - テスト用のシードデータを投入
func (td *TestDatabase) SeedTestData(t *testing.T) {
	t.Helper()

	// テスト用ユーザー
	_, err := td.DB.Exec(`
		INSERT INTO users (id, auth0_user_id, email, name) VALUES
		('00000000-0000-0000-0000-000000000001', 'auth0|test-user-001', 'test@example.com', 'Test User')
	`)
	if err != nil {
		t.Fatalf("Failed to seed test user: %v", err)
	}

	// テスト用ナレッジ
	_, err = td.DB.Exec(`
		INSERT INTO knowledge (id, user_id, title, content, category, priority, source_type) VALUES
		('11111111-1111-1111-1111-111111111111', '00000000-0000-0000-0000-000000000001', 
		 'Test Knowledge 1', 'Test content 1', 'testing', 5, 'manual'),
		('22222222-2222-2222-2222-222222222222', '00000000-0000-0000-0000-000000000001',
		 'Test Knowledge 2', 'Test content 2', 'clean_code', 4, 'manual')
	`)
	if err != nil {
		t.Fatalf("Failed to seed test knowledge: %v", err)
	}
}

// getEnvOrDefault - 環境変数かデフォルト値を取得
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
