# ReviewApp Backend - Testing Guide

## 🧪 テスト戦略

### テストの種類
1. **Unit Tests** - 個別コンポーネントのテスト
2. **Integration Tests** - データベース連携を含むテスト
3. **E2E Tests** - エンドツーエンドのテスト

## 📁 テストディレクトリ構成

```
backend/
├── internal/
│   ├── application/usecase/
│   │   ├── knowledge/
│   │   │   ├── create_knowledge_test.go    ✅ 作成済み
│   │   │   └── list_knowledge_test.go      📝 作成予定
│   │   └── review/
│   │       └── review_code_test.go         ✅ 作成済み
│   ├── domain/
│   │   ├── model/
│   │   │   ├── knowledge_test.go           📝 作成予定
│   │   │   ├── review_test.go              📝 作成予定
│   │   │   └── user_test.go                📝 作成予定
│   │   └── service/
│   │       └── review_service_test.go      📝 作成予定
│   ├── infrastructure/
│   │   ├── config/
│   │   │   └── config_test.go              📝 作成予定
│   │   ├── external/
│   │   │   └── claude_client_test.go       📝 作成予定
│   │   └── persistence/postgres/
│   │       ├── knowledge_repository_test.go 📝 作成予定
│   │       └── review_repository_test.go    📝 作成予定
│   └── interfaces/http/handler/
│       ├── knowledge_handler_test.go       📝 作成予定
│       └── review_handler_test.go          ✅ 作成済み
├── test/
│   ├── testutil/
│   │   ├── database.go                     ✅ 作成済み
│   │   ├── fixtures.go                     📝 作成予定
│   │   └── mocks.go                        ✅ 作成済み
│   ├── integration/
│   │   ├── api_test.go                     ✅ 作成済み
│   │   └── database_test.go                📝 作成予定
│   └── e2e/
│       └── review_flow_test.go             📝 作成予定
└── test_main.go                            ✅ 作成済み
```

## 🚀 テストの実行方法

### 1. Unit Tests
```bash
# すべてのユニットテストを実行
go test ./internal/...

# 特定のパッケージのテスト
go test ./internal/application/usecase/review

# カバレッジ付きで実行
go test -cover ./internal/...

# 詳細出力
go test -v ./internal/...
```

### 2. Integration Tests
```bash
# データベースが必要なため、事前に起動
docker-compose up -d postgres

# 統合テストを実行
go test ./test/integration/...

# 統合テストをスキップしてユニットテストのみ
go test -short ./...
```

### 3. 全テスト実行
```bash
# すべてのテストを実行
go test ./...

# カバレッジレポート生成
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 🔧 テスト用データベース設定

### PostgreSQL設定
```bash
# テスト用データベースを作成
createdb reviewapp_test

# マイグレーション実行
psql reviewapp_test < migrations/001_init.sql
```

### Docker Compose（推奨）
```yaml
# docker-compose.test.yml
version: '3.8'
services:
  postgres-test:
    image: postgres:15
    environment:
      POSTGRES_DB: reviewapp_test
      POSTGRES_USER: test_user
      POSTGRES_PASSWORD: test_password
    ports:
      - "5433:5432"
    volumes:
      - ./migrations:/docker-entrypoint-initdb.d
```

## 📝 テストのベストプラクティス

### 1. ファイル命名規則
- テストファイル: `*_test.go`
- テスト関数: `TestXxx(*testing.T)`
- ベンチマーク: `BenchmarkXxx(*testing.B)`

### 2. テストケース構造
```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name           string
        input          SomeInput
        expected       SomeOutput
        expectedError  bool
    }{
        {
            name:     "正常ケース",
            input:    SomeInput{},
            expected: SomeOutput{},
        },
        // ...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // テスト実行
        })
    }
}
```

### 3. モックの使用
```go
// testutil.MockXxx を使用
mockRepo := testutil.NewMockKnowledgeRepository()
mockRepo.SetKnowledges(testKnowledges)
```

### 4. データベーステスト
```go
func TestWithDatabase(t *testing.T) {
    testDB := testutil.NewTestDatabase(t)
    defer testDB.Close()
    defer testDB.Cleanup(t)
    testDB.SeedTestData(t)
    
    // テスト実行
}
```

## 🎯 優先度の高いテスト

### Phase 1 (MVP用)
1. ✅ `review_code_test.go` - コアビジネスロジック
2. ✅ `review_handler_test.go` - API エンドポイント
3. ✅ `api_test.go` - 統合テスト
4. 📝 `review_service_test.go` - プロンプト生成ロジック
5. 📝 `knowledge_repository_test.go` - データベース操作

### Phase 2 (拡張機能用)
1. 📝 `claude_client_test.go` - 外部API連携
2. 📝 `config_test.go` - 設定管理
3. 📝 `model_test.go` - ドメインモデル
4. 📝 `e2e_test.go` - エンドツーエンド

## 🚦 CI/CD での実行

### GitHub Actions 例
```yaml
name: Test
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15
        env:
          POSTGRES_PASSWORD: test_password
          POSTGRES_DB: reviewapp_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: 1.21
    
    - name: Run tests
      env:
        DATABASE_URL: postgres://postgres:test_password@localhost:5432/reviewapp_test?sslmode=disable
      run: |
        go test -v -cover ./...
```

## 📊 テストカバレッジ目標

- **Unit Tests**: 80% 以上
- **Integration Tests**: 主要なAPIエンドポイント
- **E2E Tests**: 重要なユーザーフロー

## 🛠️ 次のステップ

1. `review_service_test.go` の実装
2. `knowledge_repository_test.go` の実装
3. E2Eテストの実装
4. CIパイプラインの設定
