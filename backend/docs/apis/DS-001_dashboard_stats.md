# DS-001: ダッシュボード統計取得API

## 📋 基本情報

| 項目 | 内容 |
|------|------|
| API Code | DS-001 |
| Method | GET |
| Endpoint | /api/v1/dashboard/stats |
| 認証 | 必須（JWT Bearer Token） |
| Phase | Phase 1（MVP） |

---

## 🎯 存在意義

### 目的
ダッシュボード画面に表示する統計情報を一括で取得する。

### ユースケース
- ダッシュボード画面の表示
- ユーザーの成長状況を可視化
- AIクローンの学習状況を確認

---

## 📥 リクエスト

### Headers
```
Authorization: Bearer {jwt_token}
```

### リクエスト例

```http
GET /api/v1/dashboard/stats
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 📤 レスポンス

### 成功（200 OK）

```json
{
  "stats": {
    "total_reviews": 127,
    "knowledge_count": 89,
    "consistency_score": 87,
    "weekly_reviews": 23
  },
  "recent_reviews": [
    {
      "id": "01J5XXXXXXXXXXXXXXXXXX",
      "file_name": "auth.go",
      "language": "Go",
      "created_at": "2024-01-20T15:30:00Z",
      "improvements_count": 3,
      "status": "warning"
    },
    {
      "id": "01J5YYYYYYYYYYYYYYYYYY",
      "file_name": "user_service.js",
      "language": "JavaScript",
      "created_at": "2024-01-19T10:15:00Z",
      "improvements_count": 5,
      "status": "warning"
    },
    {
      "id": "01J5ZZZZZZZZZZZZZZZZZZ",
      "file_name": "api_handler.py",
      "language": "Python",
      "created_at": "2024-01-18T14:45:00Z",
      "improvements_count": 2,
      "status": "success"
    }
  ],
  "skill_analysis": {
    "error_handling": 35,
    "testing": 25,
    "performance": 15,
    "security": 10,
    "clean_code": 10,
    "architecture": 5,
    "other": 0
  }
}
```

### レスポンスフィールド説明

#### stats（統計情報）

| フィールド | 型 | 説明 |
|-----------|-----|------|
| total_reviews | int | 総レビュー回数 |
| knowledge_count | int | 総ナレッジ数（有効なもののみ） |
| consistency_score | int | 一貫性スコア（0-100、フィードバックスコアの平均） |
| weekly_reviews | int | 今週のレビュー回数 |

#### recent_reviews（最近のレビュー、最大5件）

| フィールド | 型 | 説明 |
|-----------|-----|------|
| id | string | レビューID |
| file_name | string | ファイル名（nullの場合は"Untitled"） |
| language | string | プログラミング言語 |
| created_at | string | 作成日時（ISO 8601） |
| improvements_count | int | 改善点の数 |
| status | string | ステータス（success, warning, error） |

#### skill_analysis（スキル分析、カテゴリ別ナレッジ数の割合）

| フィールド | 型 | 説明 |
|-----------|-----|------|
| error_handling | int | エラーハンドリングのナレッジ割合（%） |
| testing | int | テストのナレッジ割合（%） |
| performance | int | パフォーマンスのナレッジ割合（%） |
| security | int | セキュリティのナレッジ割合（%） |
| clean_code | int | クリーンコードのナレッジ割合（%） |
| architecture | int | アーキテクチャのナレッジ割合（%） |
| other | int | その他のナレッジ割合（%） |

### エラーレスポンス

#### 401 Unauthorized（認証エラー）
```json
{
  "error": "unauthorized",
  "message": "認証が必要です"
}
```

#### 500 Internal Server Error
```json
{
  "error": "internal_error",
  "message": "サーバーエラーが発生しました"
}
```

---

## 🔧 ビジネスロジック

### 処理フロー

```
1. リクエスト受信（Handler）
   ↓
2. JWT検証 → user_id 取得
   ↓
3. UseCase実行
   a) 総レビュー回数を取得
      - ReviewRepository.CountByUserID()
   
   b) 総ナレッジ数を取得（有効なもののみ）
      - KnowledgeRepository.CountByUserID()
   
   c) 一貫性スコアを計算
      - ReviewRepository.GetAverageFeedbackScore()
      - フィードバックスコア（1-3）の平均を0-100に変換
      - 計算式: ((average - 1) / 2) * 100
      - スコアなしの場合は0
   
   d) 今週のレビュー回数を取得
      - ReviewRepository.CountByUserIDAndDateRange()
      - 今週の月曜 00:00:00 から現在まで
   
   e) 最近のレビューを取得（最大5件）
      - ReviewRepository.FindRecentByUserID(limit=5)
      - structured_result.improvements の長さをカウント
   
   f) カテゴリ別ナレッジ数を取得
      - KnowledgeRepository.CountByCategory()
      - 各カテゴリの割合を計算（%）
      - 合計が100%になるように調整
   ↓
4. レスポンスを構築してクライアントに返す
   ↓
5. レスポンスヘッダーに X-API-Code: DS-001 を追加
```

### 一貫性スコアの計算

```go
// フィードバックスコア（1-3）を0-100に変換
// 1 → 0%, 2 → 50%, 3 → 100%
func calculateConsistencyScore(averageScore float64) int {
    if averageScore == 0 {
        return 0
    }
    score := ((averageScore - 1) / 2) * 100
    return int(math.Round(score))
}
```

### スキル分析の計算

```go
// カテゴリ別ナレッジ数を％に変換
func calculateSkillPercentages(categoryCounts map[string]int) map[string]int {
    total := 0
    for _, count := range categoryCounts {
        total += count
    }
    
    if total == 0 {
        return map[string]int{
            "error_handling": 0,
            "testing": 0,
            "performance": 0,
            "security": 0,
            "clean_code": 0,
            "architecture": 0,
            "other": 0,
        }
    }
    
    percentages := make(map[string]int)
    for category, count := range categoryCounts {
        percentages[category] = int(math.Round(float64(count) / float64(total) * 100))
    }
    
    return percentages
}
```

### 今週の計算

```go
// 今週の月曜 00:00:00 から現在まで
func getThisWeekRange() (time.Time, time.Time) {
    now := time.Now()
    weekday := now.Weekday()
    
    // 月曜日を週の開始とする
    daysFromMonday := int(weekday) - 1
    if daysFromMonday < 0 {
        daysFromMonday = 6 // 日曜日の場合
    }
    
    monday := now.AddDate(0, 0, -daysFromMonday)
    mondayStart := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
    
    return mondayStart, now
}
```

---

## 📁 実装ファイル

| 層 | ファイルパス | 役割 |
|----|-------------|------|
| Handler | `internal/interfaces/http/handler/dashboard_handler.go` | HTTPリクエスト処理 |
| UseCase | `internal/application/usecase/dashboard/get_stats.go` | ビジネスロジック |
| Repository | `internal/infrastructure/persistence/postgres/review_repository.go` | DB操作（拡張） |
| Repository | `internal/infrastructure/persistence/postgres/knowledge_repository.go` | DB操作（拡張） |
| Domain | `internal/domain/repository/review_repository.go` | インターフェース拡張 |
| Domain | `internal/domain/repository/knowledge_repository.go` | インターフェース拡張 |

---

## 🗄️ 必要なリポジトリメソッド

### ReviewRepository に追加

```go
// CountByUserID - ユーザーIDでレビュー総数を取得
CountByUserID(ctx context.Context, userID string) (int, error)

// CountByUserIDAndDateRange - 期間内のレビュー数を取得
CountByUserIDAndDateRange(ctx context.Context, userID string, from, to time.Time) (int, error)

// GetAverageFeedbackScore - フィードバックスコアの平均を取得
GetAverageFeedbackScore(ctx context.Context, userID string) (float64, error)
```

### KnowledgeRepository に追加

```go
// CountByUserID - ユーザーIDでナレッジ総数を取得（有効なもののみ）
CountByUserID(ctx context.Context, userID string) (int, error)

// CountByCategory - カテゴリ別のナレッジ数を取得
CountByCategory(ctx context.Context, userID string) (map[string]int, error)
```

---

## 🧪 テストケース

### 正常系

- [ ] **TC-DS-001-01**: データが存在する場合の取得
  - 前提: レビュー10件、ナレッジ5件が存在
  - 期待結果: 200 OK、正しい統計情報

- [ ] **TC-DS-001-02**: データが存在しない場合（新規ユーザー）
  - 前提: レビュー0件、ナレッジ0件
  - 期待結果: 200 OK、全て0

- [ ] **TC-DS-001-03**: フィードバックスコアがない場合
  - 前提: レビューはあるがフィードバックなし
  - 期待結果: 200 OK、consistency_score=0

- [ ] **TC-DS-001-04**: 今週のレビューがない場合
  - 前提: レビューはあるが今週はなし
  - 期待結果: 200 OK、weekly_reviews=0

- [ ] **TC-DS-001-05**: カテゴリ別ナレッジが偏っている場合
  - 前提: error_handling のみ100%
  - 期待結果: 200 OK、error_handling=100、他は0

- [ ] **TC-DS-001-06**: 一貫性スコアの計算確認
  - 前提: フィードバックスコア平均 = 2.5
  - 期待結果: consistency_score = 75

### 異常系（認証）

- [ ] **TC-DS-001-07**: JWT トークンなし
  - 期待結果: 401 Unauthorized

- [ ] **TC-DS-001-08**: JWT トークンが無効
  - 期待結果: 401 Unauthorized

---

## 📊 実装状況

- [x] ドキュメント作成
- [ ] 設計レビュー
- [ ] 実装
  - [ ] UseCase実装
  - [ ] Handler実装
  - [ ] ReviewRepository拡張
  - [ ] KnowledgeRepository拡張
  - [ ] DI設定（Wire）
  - [ ] ルーティング追加
- [ ] 単体テスト
- [ ] 統合テスト
- [ ] コードレビュー
- [ ] デプロイ

---

## 💡 実装時の注意点

### パフォーマンス
- 統計情報は複数のクエリを並行実行可能（goroutine）
- キャッシュ可能（5分程度）
- インデックス必須:
  - `idx_reviews_user_created (user_id, created_at)`
  - `idx_reviews_user_feedback (user_id, feedback_score)`
  - `idx_knowledge_user_active (user_id, is_active)`
  - `idx_knowledge_user_category (user_id, category)`

### セキュリティ
- user_id は必ずJWTから取得
- 他ユーザーのデータにアクセス不可

### エラーハンドリング
- 個別のクエリエラーは0として扱う（全体をエラーにしない）
- ログに詳細を記録

### ログ
```
[DS-001] GetDashboardStats started - user_id: xxx
[DS-001] Total reviews: 127
[DS-001] Knowledge count: 89
[DS-001] Consistency score: 87 (avg feedback: 2.74)
[DS-001] Weekly reviews: 23
[DS-001] Recent reviews: 5 items
[DS-001] Skill analysis: error_handling=35%, testing=25%...
[DS-001] Response: 200 OK
```

---

## 🔗 関連API

- [RV-002: レビュー履歴一覧取得](./RV-002_list_reviews.md)
- [KN-002: ナレッジ一覧取得](./KN-002_list_knowledge.md)

---

## 📝 変更履歴

| 日付 | バージョン | 変更内容 | 担当 |
|------|-----------|---------|------|
| 2025-01-XX | 1.0 | 初版作成 | - |

---
