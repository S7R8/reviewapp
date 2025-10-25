package testutil_test

import (
	"testing"

	"github.com/s7r8/reviewapp/test/testutil"
)

// TestDatabaseConnection - データベース接続テスト
func TestDatabaseConnection(t *testing.T) {
	t.Log("Testing database connection...")

	// データベース接続テスト
	testDB := testutil.NewTestDatabase(t)
	defer testDB.Close()

	// 接続確認
	err := testDB.DB.Ping()
	if err != nil {
		t.Fatalf("Database ping failed: %v", err)
	}

	t.Log("Database connection successful")

	// 基本的なクエリテスト
	var count int
	err = testDB.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		t.Logf("Users table query failed (expected if table doesn't exist): %v", err)
		
		// テーブルが存在するか確認
		var exists bool
		err = testDB.DB.QueryRow(`
			SELECT EXISTS (
				SELECT 1 FROM information_schema.tables 
				WHERE table_name = 'users'
			)
		`).Scan(&exists)
		
		if err != nil {
			t.Fatalf("Failed to check if users table exists: %v", err)
		}
		
		if !exists {
			t.Error("Users table does not exist. Did you run migrations?")
		}
	} else {
		t.Logf("Users table has %d records", count)
	}

	// シードデータテスト
	t.Log("Testing seed data...")
	testDB.SeedTestData(t)
	t.Log("Seed data successful")

	// クリーンアップテスト
	t.Log("Testing cleanup...")
	testDB.Cleanup(t)
	t.Log("Cleanup successful")
}
