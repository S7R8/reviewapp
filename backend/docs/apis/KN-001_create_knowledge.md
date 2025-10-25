# KN-001: ナレッジ作成API

## 📋 基本情報

| 項目 | 内容 |
|------|------|
| API Code | KN-001 |
| Method | POST |
| Endpoint | /api/v1/knowledge |
| 認証 | 必須（JWT Bearer Token） |
| Phase | Phase 1（MVP） |

---

## 🎯 存在意義

### 目的
ユーザーのコーディング哲学・ルール・学びを手動で登録する。

### ユースケース
- ユーザーが重視するコーディングルールを登録
- レビューから学んだことを記録
- チーム内の暗黙知を明文化

---

## 📥 リクエスト

### Headers
```
Content-Type: application/json
Authorization: Bearer {jwt_token}
```

### Body Schema

| フィールド | 型 | 必須 | 制約 | 説明 |
|-----------|-----|------|------|------|
| title | string | ✅ | max 200文字 | ナレッジのタイトル |
| content | string | ✅ | - | ナレッジの内容 |
| category | string | ✅ | - | カテゴリ（後述） |
| priority | integer | ✅ | 1-5 | 重要度（1=低、5=高） |

### Category 許可値

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

```json
POST /api/v1/knowledge
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "title": "エラーハンドリングの原則",
  "content": "エラーは必ずログに出力し、ユーザー向けメッセージと開発者向け詳細を分ける。contextを使ってエラーチェーンを保持する。",
  "category": "error_handling",
  "priority": 5
}
```

---

## 📤 レスポンス

### 成功（201 Created）

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "user_id": "00000000-0000-0000-0000-000000000001",
  "title": "エラーハンドリングの原則",
  "content": "エラーは必ずログに出力し、ユーザー向けメッセージと開発者向け詳細を分ける。contextを使ってエラーチェーンを保持する。",
  "category": "error_handling",
  "priority": 5,
  "source_type": "manual",
  "source_id": null,
  "usage_count": 0,
  "last_used_at": null,
  "is_active": true,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

### エラーレスポンス

#### 400 Bad Request（バリデーションエラー）
```json
{
  "error": "validation_error",
  "message": "タイトルは必須です",
  "details": {
    "field": "title",
    "reason": "required"
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
| title | 必須 | "タイトルは必須です" |
| title | 200文字以内 | "タイトルは200文字以内にしてください" |
| content | 必須 | "内容は必須です" |
| category | 必須 | "カテゴリは必須です" |
| category | 許可値 | "無効なカテゴリです" |
| priority | 必須 | "重要度は必須です" |
| priority | 1-5の整数 | "重要度は1-5の整数で指定してください" |

---

## 🔧 ビジネスロジック

### 自動設定される値

| フィールド | 値 | 説明 |
|-----------|-----|------|
| id | UUID v4 | 自動生成 |
| user_id | JWT から取得 | 認証情報から取得 |
| source_type | "manual" | 手動作成の場合は固定 |
| source_id | null | 手動作成の場合はnull |
| usage_count | 0 | 初期値は0 |
| last_used_at | null | 初期値はnull |
| is_active | true | 初期値はtrue |
| created_at | 現在時刻 | 自動設定 |
| updated_at | 現在時刻 | 自動設定 |

### 処理フロー

```
1. リクエスト受信（Handler）
   ↓
2. JWT検証 → user_id 取得
   ↓
3. リクエストボディのバリデーション
   ↓
4. Domainモデルの作成（UseCase）
   - NewKnowledge() を呼び出し
   - 必須フィールドを設定
   ↓
5. データベースに保存（Repository）
   - INSERT INTO knowledge ...
   ↓
6. 作成されたナレッジを返す（Handler）
   ↓
7. レスポンスヘッダーに X-API-Code: KN-001 を追加
```

---

## 📁 実装ファイル

| 層 | ファイルパス | 役割 |
|----|-------------|------|
| Handler | `internal/interfaces/http/handler/knowledge_handler.go` | HTTPリクエスト処理 |
| UseCase | `internal/application/usecase/knowledge/create_knowledge.go` | ビジネスロジック |
| Repository | `internal/infrastructure/persistence/postgres/knowledge_repository.go` | DB操作 |
| Domain | `internal/domain/model/knowledge.go` | エンティティ定義 |

---

## 🧪 テストケース

### 正常系

- [x] **TC-KN-001-01**: 全フィールド正常
  - 期待結果: 201 Created、ナレッジが作成される

### 異常系（バリデーション）

- [ ] **TC-KN-001-02**: title が空
  - 期待結果: 400 Bad Request

- [ ] **TC-KN-001-03**: title が201文字
  - 期待結果: 400 Bad Request

- [ ] **TC-KN-001-04**: content が空
  - 期待結果: 400 Bad Request

- [ ] **TC-KN-001-05**: category が空
  - 期待結果: 400 Bad Request

- [ ] **TC-KN-001-06**: category が許可値外
  - 期待結果: 400 Bad Request

- [ ] **TC-KN-001-07**: priority が 0
  - 期待結果: 400 Bad Request

- [ ] **TC-KN-001-08**: priority が 6
  - 期待結果: 400 Bad Request

### 異常系（認証）

- [ ] **TC-KN-001-09**: JWT トークンなし
  - 期待結果: 401 Unauthorized

- [ ] **TC-KN-001-10**: JWT トークンが無効
  - 期待結果: 401 Unauthorized

### 異常系（サーバーエラー）

- [ ] **TC-KN-001-11**: DB接続エラー
  - 期待結果: 500 Internal Server Error

---

## 📊 実装状況

- [x] ドキュメント作成
- [x] 設計レビュー
- [x] 実装開始
- [x] 実装完了（基本機能）
- [ ] JWT認証実装（現在は開発用固定ID）
- [ ] 単体テスト
- [ ] 統合テスト
- [ ] コードレビュー
- [ ] デプロイ

### 実装済み機能
✅ DB接続（PostgreSQL）  
✅ Repositoryパターン（DIP対応）  
✅ UseCase実装  
✅ Handler実装  
✅ バリデーション（ドメイン層・ハンドラ層）  
✅ Wire（DI自動化）  
✅ エラーハンドリング  
✅ グレースフルシャットダウン  

### 未実装・制限事項
❌ JWT認証（現在は固定ユーザーID: `00000000-0000-0000-0000-000000000001`）  
❌ 自動テスト  
❌ API Rate Limiting  
❌ ロギング（構造化ログ）  

---

## 📝 変更履歴

| 日付 | バージョン | 変更内容 | 担当 |
|------|-----------|---------|------|
| 2024-XX-XX | 1.0 | 初版作成 | - |
| 2025-01-XX | 1.1 | 基本機能実装完了（JWT認証除く） | - |

---

## 💡 実装時の注意点

### セキュリティ
- user_id は必ずJWTから取得（リクエストボディから受け取らない）
- SQLインジェクション対策（prepared statement使用）

### パフォーマンス
- INSERT は軽い処理なので特別な最適化は不要
- ただし、将来的に Embedding 生成を追加する場合は非同期処理を検討

### エラーハンドリング
- DB エラーは詳細をログに記録、ユーザーには汎用メッセージ
- バリデーションエラーは具体的に返す

---

## 🔗 関連API

- [KN-002: ナレッジ一覧取得](./KN-002_list_knowledge.md)
- [KN-003: ナレッジ詳細取得](./KN-003_get_knowledge.md)
- [RV-001: コードレビュー実行](./RV-001_review_code.md)（このナレッジを参照する）
