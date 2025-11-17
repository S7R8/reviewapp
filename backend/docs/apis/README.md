# API一覧

## 命名規則

```
[カテゴリ]-[連番]_[機能名].md

カテゴリ:
- AU: Auth（認証）
- DS: Dashboard（ダッシュボード）
- KN: Knowledge（ナレッジ）
- RV: Review（レビュー）
- US: User（ユーザー）
- TG: Tag（タグ）
```

## Auth APIs

| API Code | Method | Endpoint | 概要 | Status | ドキュメント |
|----------|--------|----------|------|--------|-------------|
| AU-001 | POST | /api/v1/auth/sync | ユーザー同期（初回ログイン） | ✅ 完了 | [AU-001](./AU-001_user_sync.md) |

## Dashboard APIs

| API Code | Method | Endpoint | 概要 | Status | ドキュメント |
|----------|--------|----------|------|--------|-----------|
| DS-001 | GET | /api/v1/dashboard/stats | ダッシュボード統計取得 | ✅ 完了 | [DS-001](./DS-001_dashboard_stats.md) |

---

## Knowledge APIs

| API Code | Method | Endpoint | 概要 | Status | ドキュメント |
|----------|--------|----------|------|--------|-------------|
| KN-001 | POST | /api/v1/knowledge | ナレッジ作成 | ✅ 完了 | [KN-001](./KN-001_create_knowledge.md) |
| KN-002 | GET | /api/v1/knowledge | ナレッジ一覧取得 | ✅ 完了 | [KN-002](./KN-002_list_knowledge.md) |
| KN-003 | GET | /api/v1/knowledge/:id | ナレッジ詳細取得 | ⏳ 未着手 | - |
| KN-004 | PUT | /api/v1/knowledge/:id | ナレッジ更新 | ⏳ Phase 2 | - |
| KN-005 | DELETE | /api/v1/knowledge/:id | ナレッジ削除 | ⏳ Phase 2 | - |
| KN-006 | GET | /api/v1/knowledge/search | ナレッジ検索（RAG） | ⏳ Phase 2 | - |

---

## Review APIs

| API Code | Method | Endpoint | 概要 | Status | ドキュメント |
|----------|--------|----------|------|--------|-------------|
| RV-001 | POST | /api/v1/reviews | コードレビュー実行 | ✅ 完了 | [RV-001](./RV-001_review_code.md) |
| RV-002 | GET | /api/v1/reviews | レビュー履歴一覧 | 📝 仕様作成完了 | [RV-002](./RV-002_list_reviews.md) |
| RV-003 | GET | /api/v1/reviews/:id | レビュー詳細取得 | ⏳ Phase 1 | - |
| RV-004 | PUT | /api/v1/reviews/:id/feedback | レビューフィードバック | ✅ 完了 | [RV-004](./RV-004_update_feedback.md) |

---

## User APIs

| API Code | Method | Endpoint | 概要 | Status | ドキュメント |
|----------|--------|----------|------|--------|-------------|
| US-001 | GET | /api/v1/users/me | 現在のユーザー情報取得 | ⏳ Phase 1 | - |

---

## Tag APIs

| API Code | Method | Endpoint | 概要 | Status | ドキュメント |
|----------|--------|----------|------|--------|-------------|
| TG-001 | GET | /api/v1/tags | タグ一覧取得 | ⏳ Phase 2 | - |
| TG-002 | POST | /api/v1/tags | タグ作成 | ⏳ Phase 2 | - |

---

## レスポンスヘッダー

全てのAPIレスポンスに以下のヘッダーを含める：

```
X-API-Code: KN-001
X-Request-ID: uuid
```

---

## ログフォーマット

```
[{request_id}] [{api_code}] {method} {path} - {message}

例:
[req-123456] [KN-001] POST /api/v1/knowledge - Handler: 開始
[req-123456] [KN-001] POST /api/v1/knowledge - UseCase: バリデーションOK
[req-123456] [KN-001] POST /api/v1/knowledge - Repository: INSERT成功
[req-123456] [KN-001] POST /api/v1/knowledge - Handler: 201レスポンス
```

---

## 最近の更新

- 2025-01-XX: DS-001 ダッシュボード統計APIを追加
- 2025-01-XX: RV-002 レビュー履歴一覧APIの仕様作成完了
- 2025-01-XX: RV-004 フィードバック更新API完了
