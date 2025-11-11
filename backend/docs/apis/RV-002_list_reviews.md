# RV-002: レビュー履歴一覧取得API

## 📋 基本情報

| 項目 | 内容 |
|------|------|
| API Code | RV-002 |
| Method | GET |
| Endpoint | /api/v1/reviews |
| 認証 | 必須（JWT Bearer Token） |
| Phase | Phase 1（MVP） |

---

## 🎯 存在意義

### 目的
ユーザーの過去のコードレビュー履歴を一覧取得し、フィルタリング・ソート・ページネーション機能を提供する。

### ユースケース
- Dashboard で最近のレビューを表示
- レビュー履歴画面で過去のレビューを確認
- 特定の言語やステータスでフィルタリング
- 成長の記録として振り返り

---

## 📥 リクエスト

### Headers
```
Authorization: Bearer {jwt_token}
```

### Query Parameters

| パラメータ | 型 | 必須 | デフォルト | 説明 |
|-----------|-----|------|-----------|------|
| page | int | ❌ | 1 | ページ番号（1から開始） |
| page_size | int | ❌ | 10 | 1ページあたりの件数（最大100） |
| language | string | ❌ | - | プログラミング言語でフィルター |
| status | string | ❌ | - | ステータスでフィルター |
| sort_by | string | ❌ | created_at | ソート対象 |
| sort_order | string | ❌ | desc | ソート順 |
| date_from | string | ❌ | - | 開始日（ISO 8601） |
| date_to | string | ❌ | - | 終了日（ISO 8601） |

### Language 許可値
```
TypeScript, JavaScript, Python, Go, Java, C++, C#, 
Ruby, PHP, Rust, Swift, Kotlin, Other
```

### Status 許可値

| 値 | 説明 |
|----|------|
| success | レビュー成功、問題なし |
| warning | レビュー成功、改善点あり |
| error | レビュー失敗、エラー発生 |
| pending | レビュー処理中 |

### SortBy 許可値

| 値 | 説明 |
|----|------|
| created_at | 作成日時でソート |
| language | 言語でソート |
| status | ステータスでソート |

### SortOrder 許可値

| 値 | 説明 |
|----|------|
| asc | 昇順 |
| desc | 降順 |

### リクエスト例

#### 基本的な取得
```http
GET /api/v1/reviews?page=1&page_size=10
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### フィルター付き
```http
GET /api/v1/reviews?language=TypeScript&status=warning
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### ソート
```http
GET /api/v1/reviews?sort_by=created_at&sort_order=desc
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### 期間指定
```http
GET /api/v1/reviews?date_from=2024-01-01T00:00:00Z&date_to=2024-01-31T23:59:59Z
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 📤 レスポンス

### 成功（200 OK）

```json
{
  "items": [
    {
      "id": "01J5XXXXXXXXXXXXXXXXXX",
      "user_id": "auth0|123456789",
      "code": "function hello() { console.log('Hello'); }",
      "language": "JavaScript",
      "status": "warning",
      "review_result": "改善点: console.logは本番環境で使用しない...",
      "knowledge_references": [
        "01J5YYYYYYYYYYYYYYYYYY"
      ],
      "created_at": "2024-01-27T15:30:00Z",
      "updated_at": "2024-01-27T15:30:05Z"
    }
  ],
  "total": 127,
  "page": 1,
  "page_size": 10,
  "total_pages": 13
}
```

### エラーレスポンス

#### 400 Bad Request（無効なパラメータ）
```json
{
  "error": "validation_error",
  "message": "無効なパラメータです",
  "details": {
    "page": "1以上の整数を指定してください"
  }
}
```

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

## ✅ バリデーションルール

| 項目 | ルール | エラーメッセージ |
|------|--------|-----------------|
| page | 1以上の整数 | "pageは1以上の整数を指定してください" |
| page_size | 1〜100の整数 | "page_sizeは1〜100の整数を指定してください" |
| language | 許可値のみ | "サポートされていない言語です" |
| status | 許可値のみ | "無効なステータスです" |
| sort_by | 許可値のみ | "無効なソート対象です" |
| sort_order | asc または desc | "無効なソート順です" |
| date_from | ISO 8601形式 | "無効な日付形式です" |
| date_to | ISO 8601形式 | "無効な日付形式です" |

---

## 🔧 ビジネスロジック

### 処理フロー

```
1. リクエスト受信（Handler）
   ↓
2. JWT検証 → user_id 取得
   ↓
3. クエリパラメータのバリデーション
   - デフォルト値設定
   - 許可値チェック
   ↓
4. UseCase実行
   - ReviewRepository.ListByUserID()
   - フィルター、ソート、ページネーション適用
   ↓
5. 総件数を取得
   - ReviewRepository.CountByUserID()
   ↓
6. ページネーション情報を計算
   - total_pages = ceil(total / page_size)
   ↓
7. レビュー履歴一覧を返す（Handler）
   ↓
8. レスポンスヘッダーに X-API-Code: RV-002 を追加
```

### デフォルト値の設定
```go
if query.Page == 0 {
    query.Page = 1
}
if query.PageSize == 0 {
    query.PageSize = 10
}
if query.SortBy == "" {
    query.SortBy = "created_at"
}
if query.SortOrder == "" {
    query.SortOrder = "desc"
}
```

### SQLクエリの動的構築
```go
// WHERE句の動的構築
where := "user_id = $1"
params := []interface{}{userID}
paramIndex := 2

if query.Language != "" {
    where += fmt.Sprintf(" AND language = $%d", paramIndex)
    params = append(params, query.Language)
    paramIndex++
}

// ORDER BY句
orderBy := fmt.Sprintf("%s %s", query.SortBy, strings.ToUpper(query.SortOrder))

// LIMIT/OFFSET
offset := (query.Page - 1) * query.PageSize
```

---

## 📁 実装ファイル

| 層 | ファイルパス | 役割 |
|----|-------------|------|
| Handler | `internal/interfaces/http/handler/review_handler.go` | HTTPリクエスト処理 |
| UseCase | `internal/application/usecase/review/list_reviews.go` | ビジネスロジック |
| Repository | `internal/infrastructure/persistence/postgres/review_repository.go` | DB操作 |
| Domain | `internal/domain/model/review.go` | エンティティ定義（既存） |

---

## 🧪 テストケース

### 正常系

- [ ] **TC-RV-002-01**: デフォルトパラメータで取得
  - 期待結果: 200 OK、10件取得

- [ ] **TC-RV-002-02**: ページネーション動作確認
  - 期待結果: 200 OK、指定ページの件数

- [ ] **TC-RV-002-03**: 言語フィルター動作確認
  - 期待結果: 200 OK、該当言語のみ

- [ ] **TC-RV-002-04**: ステータスフィルター動作確認
  - 期待結果: 200 OK、該当ステータスのみ

- [ ] **TC-RV-002-05**: ソート動作確認（昇順・降順）
  - 期待結果: 200 OK、ソート順が正しい

- [ ] **TC-RV-002-06**: 期間フィルター動作確認
  - 期待結果: 200 OK、期間内のみ

- [ ] **TC-RV-002-07**: 複数フィルター組み合わせ
  - 期待結果: 200 OK、全条件を満たすもののみ

- [ ] **TC-RV-002-08**: 該当なし（空配列）
  - 期待結果: 200 OK、空のitemsと total=0

### 異常系（バリデーション）

- [ ] **TC-RV-002-09**: 無効なページ番号（0以下）
  - 期待結果: 400 Bad Request

- [ ] **TC-RV-002-10**: 無効なページサイズ（101以上）
  - 期待結果: 400 Bad Request

- [ ] **TC-RV-002-11**: 無効な言語
  - 期待結果: 400 Bad Request

- [ ] **TC-RV-002-12**: 無効なステータス
  - 期待結果: 400 Bad Request

- [ ] **TC-RV-002-13**: 無効な日付フォーマット
  - 期待結果: 400 Bad Request

### 異常系（認証）

- [ ] **TC-RV-002-14**: JWT トークンなし
  - 期待結果: 401 Unauthorized

- [ ] **TC-RV-002-15**: JWT トークンが無効
  - 期待結果: 401 Unauthorized

---

## 📊 実装状況

- [x] ドキュメント作成
- [ ] 設計レビュー
- [ ] 実装
  - [ ] UseCase実装
  - [ ] Handler実装
  - [ ] Repository拡張
  - [ ] DI設定（Wire）
  - [ ] ルーティング追加
- [ ] 単体テスト
- [ ] 統合テスト
- [ ] コードレビュー
- [ ] デプロイ

### 実装予定機能
⏳ レビュー履歴一覧取得  
⏳ ページネーション機能  
⏳ 言語フィルター  
⏳ ステータスフィルター  
⏳ ソート機能  
⏳ 期間フィルター  
⏳ 総件数取得  

### 未実装・制限事項
❌ JWT認証（Phase 1では固定ユーザーID）  
❌ 全文検索機能（Phase 2）  
❌ キャッシング（Phase 2）  
❌ レート制限  

---

## 💡 実装時の注意点

### パフォーマンス
- **インデックス:** `idx_user_created (user_id, created_at DESC)` が必須
- **COUNT クエリ:** 総件数は別途取得（キャッシュ可能）
- **ページサイズ上限:** 100件まで（過負荷防止）

### セキュリティ
- user_id は必ずJWTから取得（クエリパラメータでは指定させない）
- SQLインジェクション対策（プリペアドステートメント使用）
- 他ユーザーのデータにアクセス不可

### エラーハンドリング
- DB エラー: 詳細をログに記録、ユーザーには汎用メッセージ
- 空配列は正常系として扱う（エラーではない）

### ログ
```
[RV-002] ListReviews started - user_id: xxx
[RV-002] Query params: page=1, page_size=10, language=TypeScript
[RV-002] Found 10 reviews (total: 127)
[RV-002] Response: 200 OK
```

---

## 🔗 関連API

- [RV-001: コードレビュー実行](./RV-001_review_code.md)（レビュー作成）
- [RV-003: レビュー詳細取得](./RV-003_get_review.md)（詳細表示）
- [RV-004: レビューフィードバック](./RV-004_update_feedback.md)（フィードバック更新）

---

## 📝 変更履歴

| 日付 | バージョン | 変更内容 | 担当 |
|------|-----------|---------|------|
| 2025-01-XX | 1.0 | 初版作成 | - |

---
