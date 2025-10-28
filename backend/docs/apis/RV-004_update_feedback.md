# RV-004: レビューフィードバックAPI

## 📋 基本情報

| 項目     | 内容                         |
| -------- | ---------------------------- |
| API Code | RV-004                       |
| Method   | PUT                          |
| Endpoint | /api/v1/reviews/:id/feedback |
| 認証     | 必須（JWT Bearer Token）     |
| Phase    | Phase 1（MVP）               |

---

## 🎯 存在意義

### 目的
レビュー結果に対するユーザーの評価（スコア・コメント）を記録し、将来的なAI学習データとして活用する。

### ユースケース
- レビューの品質評価
- 役に立った/立たなかった指摘の記録
- AIの一貫性スコア算出に活用
- 将来的なファインチューニングデータとして蓄積

---

## 📥 リクエスト

### Headers
```
Content-Type: application/json
Authorization: Bearer {jwt_token}
```

### Path Parameters

| パラメータ | 型            | 説明       |
| ---------- | ------------- | ---------- |
| id         | string (UUID) | レビューID |

### Body Schema

| フィールド | 型     | 必須 | 制約        | 説明                                               |
| ---------- | ------ | ---- | ----------- | -------------------------------------------------- |
| score      | int    | ✅    | 1-3         | フィードバックスコア（1: Bad, 2: Normal, 3: Good） |
| comment    | string | ❌    | max 500文字 | フィードバックコメント（オプション）               |

### スコアの意味

| スコア | 評価     | 意味               |
| ------ | -------- | ------------------ |
| 1      | 👎 Bad    | 役に立たなかった   |
| 2      | - Normal | 普通（デフォルト） |
| 3      | 👍 Good   | 役に立った         |

### リクエスト例

```json
PUT /api/v1/reviews/123e4567-e89b-12d3-a456-426614174001/feedback
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "score": 3,
  "comment": "エラーハンドリングの指摘が的確でした。特に開発者向け詳細とユーザー向けメッセージを分ける点が参考になりました。"
}
```

```json
PUT /api/v1/reviews/123e4567-e89b-12d3-a456-426614174001/feedback
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "score": 1
}
```

---

## 📤 レスポンス

### 成功（200 OK）

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174001",
  "feedback_score": 3,
  "feedback_comment": "エラーハンドリングの指摘が的確でした。特に開発者向け詳細とユーザー向けメッセージを分ける点が参考になりました。",
  "updated_at": "2025-10-28T12:34:56Z"
}
```

### エラーレスポンス

#### 400 Bad Request（バリデーションエラー）
```json
{
  "error": "validation_error",
  "message": "スコアは1-3の整数で指定してください",
  "details": {
    "field": "score",
    "value": 5,
    "constraint": "min=1,max=3"
  }
}
```

```json
{
  "error": "validation_error",
  "message": "コメントは500文字以内にしてください",
  "details": {
    "field": "comment",
    "value_length": 523,
    "max_length": 500
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

#### 403 Forbidden（権限エラー）
```json
{
  "error": "forbidden",
  "message": "このレビューを更新する権限がありません"
}
```

#### 404 Not Found（レビューが存在しない）
```json
{
  "error": "not_found",
  "message": "レビューが見つかりません",
  "details": {
    "review_id": "123e4567-e89b-12d3-a456-426614174001"
  }
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

| 項目    | ルール      | エラーメッセージ                      |
| ------- | ----------- | ------------------------------------- |
| score   | 必須        | "スコアは必須です"                    |
| score   | 1-3の整数   | "スコアは1-3の整数で指定してください" |
| comment | 500文字以内 | "コメントは500文字以内にしてください" |

---

## 🔧 ビジネスロジック

### 処理フロー

```
1. リクエスト受信（Handler）
   ↓
2. JWT検証 → user_id 取得
   ↓
3. パスパラメータからreview_id取得
   ↓
4. リクエストボディのバリデーション
   - score: 1-3の範囲チェック
   - comment: 500文字以内チェック
   ↓
5. レビューの存在確認（Repository）
   - ReviewRepository.FindByID()
   ↓
6. 権限チェック
   - review.user_id == 認証user_id
   ↓
7. フィードバック更新（Repository）
   - ReviewRepository.UpdateFeedback()
   - UPDATE reviews SET feedback_score = ?, feedback_comment = ?, updated_at = NOW()
   ↓
8. 更新されたレビュー情報を返す（Handler）
   ↓
9. レスポンスヘッダーに X-API-Code: RV-004 を追加
```

### 上書き仕様

- ユーザーは何度でもフィードバックを変更可能
- 既存のフィードバックがある場合は上書き
- `updated_at` は更新時刻で上書き

### デフォルト値

| フィールド       | デフォルト値 | 説明                 |
| ---------------- | ------------ | -------------------- |
| feedback_score   | null         | レビュー作成時はnull |
| feedback_comment | null         | レビュー作成時はnull |

---

## 📁 実装ファイル

| 層         | ファイルパス                                                        | 役割                                |
| ---------- | ------------------------------------------------------------------- | ----------------------------------- |
| Handler    | `internal/interfaces/http/handler/review_handler.go`                | HTTPリクエスト処理                  |
| UseCase    | `internal/application/usecase/review/update_feedback.go`            | ビジネスロジック                    |
| Repository | `internal/infrastructure/persistence/postgres/review_repository.go` | DB操作（UpdateFeedback追加）        |
| Domain     | `internal/domain/model/review.go`                                   | エンティティ定義（SetFeedback更新） |

---

## 🧪 テストケース

### 正常系

- [ ] **TC-RV-004-01**: 全フィールド正常（score=3, commentあり）
  - 期待結果: 200 OK、フィードバック更新成功

- [ ] **TC-RV-004-02**: スコアのみ（commentなし）
  - 期待結果: 200 OK、commentはnullのまま

- [ ] **TC-RV-004-03**: 既存フィードバックの上書き
  - 期待結果: 200 OK、スコアとコメントが上書きされる

- [ ] **TC-RV-004-04**: スコア=1（Bad）
  - 期待結果: 200 OK

- [ ] **TC-RV-004-05**: スコア=2（Normal）
  - 期待結果: 200 OK

### 異常系（バリデーション）

- [ ] **TC-RV-004-06**: scoreが空
  - 期待結果: 400 Bad Request

- [ ] **TC-RV-004-07**: scoreが0
  - 期待結果: 400 Bad Request

- [ ] **TC-RV-004-08**: scoreが4
  - 期待結果: 400 Bad Request

- [ ] **TC-RV-004-09**: commentが501文字
  - 期待結果: 400 Bad Request

- [ ] **TC-RV-004-10**: review_idがUUID形式でない
  - 期待結果: 400 Bad Request

### 異常系（認証・権限）

- [ ] **TC-RV-004-11**: JWTトークンなし
  - 期待結果: 401 Unauthorized

- [ ] **TC-RV-004-12**: JWTトークンが無効
  - 期待結果: 401 Unauthorized

- [ ] **TC-RV-004-13**: 他人のレビューにフィードバック
  - 期待結果: 403 Forbidden

### 異常系（データ）

- [ ] **TC-RV-004-14**: 存在しないreview_id
  - 期待結果: 404 Not Found

- [ ] **TC-RV-004-15**: 削除済みレビュー（deleted_at != null）
  - 期待結果: 404 Not Found

### 統合テスト

- [ ] **TC-RV-004-16**: レビュー作成 → フィードバック送信 → 再取得で確認
  - 期待結果: フィードバックが正しく保存されている

- [ ] **TC-RV-004-17**: フィードバック送信 → 複数回上書き
  - 期待結果: 最後のフィードバックが保存されている

---

## 📊 実装状況

- [x] ドキュメント作成
- [x] 設計レビュー
- [x] Domain モデル更新（SetFeedback拡張）
- [x] Repository インターフェース追加
- [x] Repository 実装
- [x] UseCase 作成
- [x] Handler 実装
- [x] ルーティング追加
- [x] **実装完了** 🎉
- [ ] 単体テスト
- [ ] 統合テスト
- [ ] コードレビュー
- [ ] デプロイ

---

## 📝 変更履歴

| 日付       | バージョン | 変更内容           | 担当 |
| ---------- | ---------- | ------------------ | ---- |
| 2025-10-28 | 1.0        | 初版作成           | -    |
| 2025-10-29 | 2.0        | **実装完了** 🎉     | -    |
|            |            | - Domain層実装     | -    |
|            |            | - Repository層実装 | -    |
|            |            | - UseCase層実装    | -    |
|            |            | - Handler層実装    | -    |
|            |            | - ルーティング追加 | -    |
|            |            | - DI設定完了       | -    |

---

## 💡 実装時の注意点

### パフォーマンス
- **単純なUPDATE:** 高速に完了（数ms）
- **インデックス:** review_id（PK）とuser_id にインデックスが必要

### セキュリティ
- **権限チェック必須:** review.user_id == JWT.user_id を必ず確認
- **SQLインジェクション対策:** プレースホルダーを使用

### エラーハンドリング
- **レビューが存在しない:** 404 Not Found（詳細はログに記録）
- **権限エラー:** 403 Forbidden（ユーザーIDはログに記録しない）
- **DB エラー:** 500 Internal Server Error（詳細をログに記録、ユーザーには汎用メッセージ）

### ログ
フィードバック送信時は以下をログ出力：
```
[RV-004] Feedback started - review_id: xxx, user_id: xxx
[RV-004] Validation OK - score: 3, comment_length: 45
[RV-004] Review found - review_id: xxx
[RV-004] Permission check OK
[RV-004] Feedback updated - review_id: xxx, score: 3
[RV-004] Response 200 OK
```

---

## 🔗 関連API

- [RV-001: コードレビュー実行](./RV-001_review_code.md)（レビュー作成）
- [RV-002: レビュー履歴一覧](./RV-002_list_reviews.md)（過去のレビュー確認）
- [RV-003: レビュー詳細取得](./RV-003_get_review.md)（フィードバック含む詳細取得）

---

## 📚 参考資料

- [DESIGN.md](../../backend/docs/DESIGN.md) - システム設計
- [review.go](../../backend/internal/domain/model/review.go) - Reviewエンティティ

---
