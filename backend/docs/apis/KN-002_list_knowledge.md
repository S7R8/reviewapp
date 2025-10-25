# KN-002: ナレッジ一覧取得API

## 📋 基本情報

| 項目 | 内容 |
|------|------|
| API Code | KN-002 |
| Method | GET |
| Endpoint | /api/v1/knowledge |
| 認証 | 必須（JWT Bearer Token） |
| Phase | Phase 1（MVP） |

---

## 🎯 存在意義

### 目的
ユーザーが登録したナレッジの一覧を取得する。

### ユースケース
- Dashboard で「あなたのクローンの成長度」を表示
- ナレッジ管理画面で一覧表示・編集
- カテゴリ別のナレッジ分析
- レビュー前の事前確認

---

## 📥 リクエスト

### Headers
```
Authorization: Bearer {jwt_token}
```

### Query Parameters

| パラメータ | 型 | 必須 | 説明 | 例 |
|-----------|-----|------|------|-----|
| category | string | ❌ | カテゴリでフィルタ（指定なし=全件） | `error_handling` |

### Category 許可値（KN-001と同じ）

| 値 | 説明 |
|----|------|
| error_handling | エラーハンドリング |
| testing | テスト |
| performance | パフォーマンス |
| security | セキュリティ |
| clean_code | クリーンコード |
| architecture | アーキテクチャ |
| other | その他 |

### リクエスト例

#### 全件取得
```http
GET /api/v1/knowledge
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### カテゴリでフィルタ
```http
GET /api/v1/knowledge?category=error_handling
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 📤 レスポンス

### 成功（200 OK）

#### 全件取得の例
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "user_id": "00000000-0000-0000-0000-000000000001",
    "title": "エラーハンドリングの原則",
    "content": "エラーは必ずログに出力し、ユーザー向けメッセージと開発者向け詳細を分ける。contextを使ってエラーチェーンを保持する。",
    "category": "error_handling",
    "priority": 5,
    "source_type": "manual",
    "source_id": null,
    "usage_count": 12,
    "last_used_at": "2024-01-20T15:30:00Z",
    "is_active": true,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-20T15:30:00Z"
  },
  {
    "id": "234e5678-e89b-12d3-a456-426614174001",
    "user_id": "00000000-0000-0000-0000-000000000001",
    "title": "関数は50行以内",
    "content": "関数は1つのことだけをする。50行を超える場合は分割を検討する。",
    "category": "clean_code",
    "priority": 4,
    "source_type": "manual",
    "source_id": null,
    "usage_count": 8,
    "last_used_at": "2024-01-18T09:00:00Z",
    "is_active": true,
    "created_at": "2024-01-16T14:00:00Z",
    "updated_at": "2024-01-18T09:00:00Z"
  }
]
```

#### 該当なし（空配列）
```json
[]
```

### エラーレスポンス

#### 400 Bad Request（無効なカテゴリ）
```json
{
  "error": "validation_error",
  "message": "無効なカテゴリです"
}
```

#### 401 Unauthorized（認証エラー）
```json
{
  "error": "unauthorized",
  "message": "認証が必要です"
}
```

#### 500 Internal Server Error（サーバーエラー）
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
| category | オプショナル | - |
| category | 許可値のみ | "無効なカテゴリです" |

---

## 🔧 ビジネスロジック

### ソート順
```
ORDER BY priority DESC, created_at DESC
```
- 重要度の高いナレッジが上位
- 同じ重要度なら新しいものが上位

### フィルタ条件
- `is_active = true` のみ取得
- `deleted_at IS NULL` のみ取得（論理削除済みは除外）

### 処理フロー

```
1. リクエスト受信（Handler）
   ↓
2. JWT検証 → user_id 取得
   ↓
3. クエリパラメータのバリデーション
   - category があれば許可値チェック
   ↓
4. UseCase実行
   - category あり → FindByCategory()
   - category なし → FindByUserID()
   ↓
5. ナレッジ一覧を返す（Handler）
   - シンプルな配列形式
   ↓
6. レスポンスヘッダーに X-API-Code: KN-002 を追加
```

---

## 📊 レスポンス仕様

### Phase 1（今回）
- **シンプルな配列**で返す
- ページングなし（全件返す）
- 総件数情報なし

```json
[
  { ...knowledge1 },
  { ...knowledge2 }
]
```

### Phase 2（将来）
- ページング追加（limit/offset）
- 総件数追加
- has_more フラグ追加

```json
{
  "knowledges": [...],
  "total": 127,
  "limit": 20,
  "offset": 0,
  "has_more": true
}
```

---

## 📁 実装ファイル

| 層 | ファイルパス | 役割 |
|----|-------------|------|
| Handler | `internal/interfaces/http/handler/knowledge_handler.go` | HTTPリクエスト処理 |
| UseCase | `internal/application/usecase/knowledge/list_knowledge.go` | ビジネスロジック |
| Repository | `internal/infrastructure/persistence/postgres/knowledge_repository.go` | DB操作（既存利用） |
| Domain | `internal/domain/model/knowledge.go` | エンティティ定義（既存） |

### 使用するRepository メソッド
- `FindByUserID()` - 既に実装済み
- `FindByCategory()` - 既に実装済み

---

## 🧪 テストケース

### 正常系

- [x] **TC-KN-002-01**: category なし（全件取得）
  - 期待結果: 200 OK、全ナレッジが返る

- [x] **TC-KN-002-02**: category あり（フィルタ）
  - 期待結果: 200 OK、該当カテゴリのみ返る

- [x] **TC-KN-002-03**: 該当なし
  - 期待結果: 200 OK、空配列 `[]`

- [x] **TC-KN-002-04**: ソート順確認
  - 期待結果: priority DESC, created_at DESC

### 異常系（バリデーション）

- [x] **TC-KN-002-05**: category が許可値外
  - 期待結果: 400 Bad Request

### 異常系（認証）

- [ ] **TC-KN-002-06**: JWT トークンなし
  - 期待結果: 401 Unauthorized

- [ ] **TC-KN-002-07**: JWT トークンが無効
  - 期待結果: 401 Unauthorized

### 異常系（サーバーエラー）

- [ ] **TC-KN-002-08**: DB接続エラー
  - 期待結果: 500 Internal Server Error

---

## 📊 実装状況

- [x] ドキュメント作成
- [x] 設計レビュー
- [x] 実装完了
  - [x] UseCase実装
  - [x] Handler実装
  - [x] DI設定（Wire）
  - [x] ルーティング追加
- [ ] 単体テスト
- [ ] 統合テスト
- [ ] コードレビュー
- [ ] デプロイ

### 実装済み機能
✅ ナレッジ一覧取得（全件）  
✅ カテゴリフィルタ機能  
✅ ソート機能（priority DESC, created_at DESC）  
✅ 空配列レスポンス（正常系）  
✅ エラーハンドリング  
✅ バリデーション（カテゴリ許可値チェック）  

### 未実装・制限事項
❌ JWT認証（現在は固定ユーザーID: `00000000-0000-0000-0000-000000000001`）  
❌ ページング機能（limit/offset）← Phase 2で実装予定  
❌ 総件数情報  
❌ 自動テスト  
❌ API Rate Limiting  
❌ ロギング（構造化ログ）  

---

## 💡 実装時の注意点

### セキュリティ
- user_id は必ずJWTから取得（リクエストから受け取らない）
- 他のユーザーのナレッジは絶対に返さない

### パフォーマンス
- Phase 1 では全件取得だが、通常のユーザーなら数十〜数百件程度
- ナレッジが1000件を超える場合は Phase 2 でページング実装
- インデックスが適切に効いているか確認（user_id, category）

### エラーハンドリング
- DB エラーは詳細をログに記録、ユーザーには汎用メッセージ
- 空配列は正常系として扱う（エラーではない）

---

## 🔗 関連API

- [KN-001: ナレッジ作成](./KN-001_create_knowledge.md)
- [KN-003: ナレッジ詳細取得](./KN-003_get_knowledge.md)（未実装）
- [KN-004: ナレッジ更新](./KN-004_update_knowledge.md)（未実装）
- [KN-005: ナレッジ削除](./KN-005_delete_knowledge.md)（未実装）

---

## 📝 変更履歴

| 日付 | バージョン | 変更内容 | 担当 |
|------|-----------|---------|------|
| 2025-01-XX | 1.0 | 初版作成 | - |
| 2025-01-XX | 1.1 | 基本機能実装完了（JWT認証・ページング除く） | - |
