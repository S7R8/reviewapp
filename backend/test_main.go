package main

import (
	"os"
	"testing"
)

// TestMain - テスト実行前後の共通処理
func TestMain(m *testing.M) {
	// テスト用環境変数を設定
	os.Setenv("ENV", "test")
	os.Setenv("PORT", "8080")
	os.Setenv("DATABASE_URL", "postgres://test_user:test_password@localhost:5432/reviewapp_test?sslmode=disable")
	os.Setenv("CLAUDE_API_KEY", "test-key")
	os.Setenv("CLAUDE_MODEL", "claude-3-5-sonnet-20241022")

	// テストを実行
	code := m.Run()

	// クリーンアップ処理があればここに記述

	os.Exit(code)
}
