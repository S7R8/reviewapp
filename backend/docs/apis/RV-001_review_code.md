# RV-001: コードレビュー実行API

## 📋 基本情報

| 項目 | 内容 |
|------|------|
| API Code | RV-001 |
| Method | POST |
| Endpoint | /api/v1/reviews |
| 認証 | 必須（JWT Bearer Token） |
| Phase | Phase 1（MVP） |

---

## 🎯 存在意義

### 目的
ユーザーのコードを、そのユーザー独自のナレッジに基づいてレビューする。

### ユースケース
- 開発中のコードをレビュー
- Pull Request前のセルフレビュー
- コーディング学習の補助

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
| code | string | ✅ | - | レビュー対象のコード |
| language | string | ✅ | - | プログラミング言語 |
| file_name | string | ❌ | max 255文字 | ファイル名（オプション） |
| context | string | ❌ | - | 追加のコンテキスト（オプション） |

### Language 推奨値

主要な言語をサポート：
```
go, python, javascript, typescript, java, c, cpp, csharp, 
rust, ruby, php, swift, kotlin, scala, html, css, sql
```

（上記以外も受け付け可能）

### リクエスト例

```json
POST /api/v1/reviews
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "code": "func HandleError(err error) {\n    if err != nil {\n        log.Println(err)\n    }\n}",
  "language": "go",
  "file_name": "handler.go",
  "context": "HTTPハンドラのエラー処理です。ユーザーにエラーメッセージを返す必要があります。"
}
```

---

## 📤 レスポンス

### 成功（201 Created）

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174001",
  "user_id": "00000000-0000-0000-0000-000000000001",
  "code": "func HandleError(err error) {\n    if err != nil {\n        log.Println(err)\n    }\n}",
  "language": "go",
  "file_name": "handler.go",
  "review_result": "## 総評\nエラーハンドリングが不十分です。以下の点を改善してください。\n\n## 改善点\n\n### 1. ユーザー向けメッセージがない\nあなたのナレッジ「エラーハンドリングの原則」によると、エラーはログ出力だけでなく、ユーザー向けメッセージと開発者向け詳細を分ける必要があります。\n\n```go\nfunc HandleError(w http.ResponseWriter, err error) {\n    if err != nil {\n        log.Printf(\"Error occurred: %+v\", err) // 開発者向け\n        http.Error(w, \"サーバーエラーが発生しました\", http.StatusInternalServerError) // ユーザー向け\n    }\n}\n```\n\n### 2. contextを使ったエラーチェーン\ncontextを使ってエラーチェーンを保持すると、デバッグが容易になります。\n\n## 参考にしたナレッジ\n- [エラーハンドリング] エラーハンドリングの原則（Priority: 5）",
  "llm_provider": "claude",
  "llm_model": "claude-3-5-sonnet-20241022",
  "tokens_used": 1250,
  "referenced_knowledge": [
    {
      "id": "knowledge-123",
      "title": "エラーハンドリングの原則",
      "category": "error_handling",
      "priority": 5
    }
  ],
  "feedback_score": null,
  "feedback_comment": null,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### エラーレスポンス

#### 400 Bad Request（バリデーションエラー）
```json
{
  "error": "validation_error",
  "message": "コードは必須です",
  "details": {
    "field": "code",
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

#### 429 Too Many Requests（レート制限）
```json
{
  "error": "rate_limit_exceeded",
  "message": "リクエスト数が上限に達しました。しばらくお待ちください。",
  "details": {
    "retry_after": 60
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

#### 503 Service Unavailable（外部APIエラー）
```json
{
  "error": "llm_api_error",
  "message": "AI APIが一時的に利用できません"
}
```

---

## ✅ バリデーションルール

| 項目 | ルール | エラーメッセージ |
|------|--------|-----------------|
| code | 必須 | "コードは必須です" |
| code | 10,000文字以内 | "コードは10,000文字以内にしてください" |
| language | 必須 | "プログラミング言語は必須です" |
| file_name | 255文字以内 | "ファイル名は255文字以内にしてください" |

---

## 🔧 ビジネスロジック

### 処理フロー（RAG統合）

```
1. リクエスト受信（Handler）
   ↓
2. JWT検証 → user_id 取得
   ↓
3. リクエストボディのバリデーション
   ↓
4. 関連ナレッジを検索（RAG: Retrieval）
   - KnowledgeRepository.FindByUserID()
   - ユーザーの全アクティブナレッジ取得
   ↓
5. プロンプト生成（ReviewService）
   - BuildPromptFromKnowledge()
   - Priority順でソート
   - Top 10に絞り込み
   - カテゴリ名付きでフォーマット
   ↓
6. LLM APIでレビュー生成（RAG: Augmented Generation）
   - ClaudeClient.ReviewCode()
   - システムプロンプト + ナレッジ + コード
   - レビュー結果を取得
   ↓
7. レビュー結果を保存（Repository）
   - INSERT INTO reviews ...
   - INSERT INTO review_knowledge ... (参照されたナレッジ)
   ↓
8. ナレッジのusage_countを更新
   - Knowledge.IncrementUsage()
   - UPDATE knowledge SET usage_count = usage_count + 1, last_used_at = NOW()
   ↓
9. レビュー結果を返す（Handler）
   ↓
10. レスポンスヘッダーに X-API-Code: RV-001 を追加
```

### 自動設定される値

| フィールド | 値 | 説明 |
|-----------|-----|------|
| id | UUID v4 | 自動生成 |
| user_id | JWT から取得 | 認証情報から取得 |
| llm_provider | "claude" | Phase 1は固定 |
| llm_model | config.CLAUDE_MODEL | 環境変数から取得 |
| tokens_used | LLMレスポンスから | Claude APIのレスポンス |
| feedback_score | null | 初期値はnull |
| feedback_comment | null | 初期値はnull |
| created_at | 現在時刻 | 自動設定 |
| updated_at | 現在時刻 | 自動設定 |

### プロンプト構造

```
あなたは {USER_NAME} のクローンとして、コードレビューを行ってください。

以下は、{USER_NAME} が重視しているコーディング原則です：

{KNOWLEDGE_PROMPT}
  ### [エラーハンドリング] エラーハンドリングの原則
  エラーは必ずログに出力し、ユーザー向けメッセージと開発者向け詳細を分ける。
  
  ### [クリーンコード] 関数は1つのことだけをする
  関数は50行以内に抑え、1つの責務のみを持つ。

## レビュー対象のコード

言語: {LANGUAGE}
ファイル名: {FILE_NAME}

```
{CODE}
```

## 追加コンテキスト
{CONTEXT}

## レビュー方針
1. 上記のナレッジに基づいて、一貫性のあるレビューを行う
2. 具体的な改善案を提示する
3. どのナレッジに基づいた指摘かを明記する
4. ポジティブな点も指摘する

## 出力形式
以下の形式でレビュー結果を出力してください：

## 総評
[全体的な評価]

## 良い点
[良い点を列挙]

## 改善点
[改善点を列挙。各項目でナレッジを引用]

## 参考にしたナレッジ
[使用したナレッジのリスト]
```

---

## 📁 実装ファイル

| 層 | ファイルパス | 役割 |
|----|-------------|------|
| Handler | `internal/interfaces/http/handler/review_handler.go` | HTTPリクエスト処理 |
| UseCase | `internal/application/usecase/review/review_code.go` | ビジネスロジック |
| Service | `internal/domain/service/review_service.go` | プロンプト生成 |
| Repository | `internal/infrastructure/persistence/postgres/review_repository.go` | DB操作 |
| Repository | `internal/infrastructure/persistence/postgres/knowledge_repository.go` | ナレッジ検索 |
| External | `internal/infrastructure/external/claude_client.go` | Claude API |
| Domain | `internal/domain/model/review.go` | エンティティ定義 |

---

## 🧪 テストケース

### 正常系

- [ ] **TC-RV-001-01**: 全フィールド正常（ナレッジあり）
  - 期待結果: 201 Created、ナレッジに基づいたレビュー

- [ ] **TC-RV-001-02**: オプションフィールドなし
  - 期待結果: 201 Created、レビュー実行

- [ ] **TC-RV-001-03**: ナレッジがない場合
  - 期待結果: 201 Created、一般的なベストプラクティスでレビュー

### 異常系（バリデーション）

- [ ] **TC-RV-001-04**: code が空
  - 期待結果: 400 Bad Request

- [ ] **TC-RV-001-05**: code が10,001文字
  - 期待結果: 400 Bad Request

- [ ] **TC-RV-001-06**: language が空
  - 期待結果: 400 Bad Request

- [ ] **TC-RV-001-07**: file_name が256文字
  - 期待結果: 400 Bad Request

### 異常系（認証）

- [ ] **TC-RV-001-08**: JWT トークンなし
  - 期待結果: 401 Unauthorized

- [ ] **TC-RV-001-09**: JWT トークンが無効
  - 期待結果: 401 Unauthorized

### 異常系（外部API）

- [ ] **TC-RV-001-10**: Claude APIエラー
  - 期待結果: 503 Service Unavailable

- [ ] **TC-RV-001-11**: Claude APIタイムアウト
  - 期待結果: 504 Gateway Timeout

- [ ] **TC-RV-001-12**: レート制限超過
  - 期待結果: 429 Too Many Requests

### 統合テスト

- [ ] **TC-RV-001-13**: RAG統合テスト
  - ナレッジ作成 → レビュー実行 → ナレッジが参照されることを確認

- [ ] **TC-RV-001-14**: usage_count更新テスト
  - レビュー実行 → ナレッジのusage_countとlast_used_atが更新されることを確認

---

## 📊 実装状況

- [x] ドキュメント作成
- [x] 設計レビュー
- [x] ReviewService実装（プロンプト生成）
- [x] **実装完了** 🎉
- [x] **Claude API統合** 🎉
- [x] **RAG統合** 🎉
- [x] **Handler実装** 🎉
- [ ] 単体テスト
- [ ] 統合テスト
- [ ] コードレビュー
- [ ] デプロイ

### 実装済み機能 🎉
✅ Review エンティティ  
✅ ReviewRepository インターフェース & 実装  
✅ ReviewService（プロンプト生成）  
✅ **ReviewCodeUseCase（完全実装）**  
✅ **ClaudeClient 実装** (anthropic-sdk-go)  
✅ **ReviewHandler 実装**  
✅ **RAG統合** (Knowledge検索 → プロンプト生成)  
✅ **Knowledge usage 更新**  
✅ Wire DI 統合  

### 未実装・制限事項
❌ JWT認証（現在は固定ユーザーID）  
❌ レート制限  
❌ リトライロジック  
❌ テストケース  

---

## 📝 変更履歴

| 日付 | バージョン | 変更内容 | 担当 |
|------|-----------|---------|------|
| 2024-XX-XX | 1.0 | 初版作成（設計） | - |
| 2025-01-XX | 1.1 | ReviewService実装完了 | - |
| 2025-01-24 | 2.0 | **実装完了** 🎉 | - |
|  |  | - Claude API統合 | - |
|  |  | - RAG統合 | - |
|  |  | - Handler実装 | - |
|  |  | - UseCase完全実装 | - |

---

## 💡 実装時の注意点

### パフォーマンス
- **ナレッジ検索:** Phase 1では全件取得だが、ReviewServiceで Top 10 に絞るため問題なし
- **Claude API:** レスポンスタイムは 3-5秒程度を想定
- **非同期処理:** Phase 1では不要（将来的には検討）

### セキュリティ
- user_id は必ずJWTから取得
- コードは保存するが、暗号化は不要（ユーザー自身のコード）
- Claude APIキーは環境変数で管理

### エラーハンドリング
- **Claude APIエラー:** リトライ3回まで、その後503を返す
- **DB エラー:** 詳細をログに記録、ユーザーには汎用メッセージ
- **タイムアウト:** 30秒でタイムアウト

### コスト管理
- **Claude API:** 1リクエスト約$0.005-0.01
- **Phase 1:** 1日100リクエスト想定 → $0.5-1.0/日
- **レート制限:** ユーザーあたり10リクエスト/分

### ログ
レビュー実行時は以下をログ出力：
```
[RV-001] Review started - user_id: xxx, language: go
[RV-001] Found 5 relevant knowledge
[RV-001] Calling Claude API - tokens: ~1000
[RV-001] Claude API response - tokens_used: 1250, duration: 3.2s
[RV-001] Review saved - review_id: xxx
[RV-001] Updated usage_count for 5 knowledge
```

---

## 🔗 関連API

- [KN-001: ナレッジ作成](./KN-001_create_knowledge.md)（ナレッジを追加）
- [KN-002: ナレッジ一覧取得](./KN-002_list_knowledge.md)（参照されるナレッジ）
- [RV-002: レビュー履歴一覧](./RV-002_list_reviews.md)（過去のレビュー確認）
- [RV-003: レビュー詳細取得](./RV-003_get_review.md)（レビュー詳細）
- [RV-004: レビューフィードバック](./RV-004_feedback.md)（精度向上）

---

## 🎨 プロンプトエンジニアリング

### システムプロンプトのポイント
1. **役割定義:** 「あなたは {USER_NAME} のクローンです」
2. **ナレッジ参照:** Priority順で整理されたナレッジを提示
3. **一貫性:** 過去の判断基準に基づくことを強調
4. **根拠明示:** どのナレッジに基づいた指摘かを明記

### 出力形式
- **構造化:** 総評、良い点、改善点、参考ナレッジ
- **具体性:** 改善案はコード例を含める
- **引用:** ナレッジのタイトルとPriorityを明記

### Phase 2での改善案
- ベクトル検索で関連度の高いナレッジのみを選択
- コードの類似度スコアを活用
- ユーザーフィードバックでプロンプト最適化

---

## 📚 参考資料

- [prompt-design.md](../../backend/docs/prompt-design.md) - プロンプト設計詳細
- [DESIGN.md](../../backend/docs/DESIGN.md) - RAG設計
- [Claude API Documentation](https://docs.anthropic.com/)

---
