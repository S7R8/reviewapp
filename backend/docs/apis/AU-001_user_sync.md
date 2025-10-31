# AU-001: ユーザー同期API

## 📋 基本情報

| 項目 | 内容 |
|------|------|
| API Code | AU-001 |
| Method | POST |
| Endpoint | /api/v1/auth/sync |
| 認証 | 必須（JWT Bearer Token） |
| Phase | Phase 1（MVP） |

---

## 🎯 存在意義

### 目的
Auth0でログインしたユーザーをバックエンドのデータベースと同期する。初回ログイン時にユーザーレコードを自動作成し、既存ユーザーの場合は情報を取得する。

### ユースケース
- Auth0でのログイン直後、フロントエンドから自動的に呼び出される
- 初回ログインユーザーのアカウント作成
- 既存ユーザーの認証状態確認とユーザー情報取得
- 後続のAPIリクエストで使用するユーザーIDの取得

---

## 📥 リクエスト

### Headers
```
Content-Type: application/json
Authorization: Bearer {jwt_token}
```

**重要**: このAPIは`id_token`（JWT形式）を使用します。Auth0から取得した`id_token`を`Authorization`ヘッダーに設定してください。

### Body
このAPIはリクエストボディを必要としません（空のJSONオブジェクト`{}`を送信）。

### リクエスト例

```http
POST /api/v1/auth/sync
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlJQVWE4...

{}
```

---

## 📤 レスポンス

### 成功（200 OK）

#### 既存ユーザーの場合
```json
{
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "auth0_user_id": "auth0|69038b5428c69abdd48d4d65",
    "email": "user@example.com",
    "name": "John Doe",
    "avatar_url": "https://example.com/avatar.jpg",
    "preferences": "{}",
    "created_at": "2024-11-01T10:00:00Z",
    "updated_at": "2024-11-01T10:00:00Z"
  }
}
```

#### 新規ユーザーの場合（初回ログイン）
```json
{
  "user": {
    "id": "987e6543-e21b-12d3-a456-426614174999",
    "auth0_user_id": "auth0|69038b5428c69abdd48d4d65",
    "email": "user-69038b5428c69abdd48d4d65@temp.local",
    "name": "New User",
    "avatar_url": null,
    "preferences": "{}",
    "created_at": "2024-11-01T15:30:00Z",
    "updated_at": "2024-11-01T15:30:00Z"
  }
}
```

### エラーレスポンス

#### 401 Unauthorized（認証エラー）
```json
{
  "error": "unauthorized",
  "message": "認証情報が見つかりません"
}
```

**発生条件**:
- `Authorization`ヘッダーが存在しない
- JWTトークンが無効
- トークンに`sub`（subject）クレームが存在しない

#### 500 Internal Server Error（サーバーエラー）
```json
{
  "error": "sync_failed",
  "message": "ユーザー同期に失敗しました"
}
```

**発生条件**:
- データベース接続エラー
- ユーザー作成時のエラー

---

## 🔧 ビジネスロジック

### トークンからの情報抽出

JWTトークンの`sub`クレームから`auth0_user_id`を取得します：

```
Token claims:
{
  "sub": "auth0|69038b5428c69abdd48d4d65",  ← この値を使用
  "iss": "https://dev-863amkaw1kj03j7y.us.auth0.com/",
  "aud": "gewOjqNbMHgsCcpGh9Ch2IYuI4FxiE62",
  "email": "user@example.com",
  "name": "John Doe"
}
```

### 処理フロー

```
1. リクエスト受信
   ↓
2. JWT検証（Middleware）
   - トークンの署名検証（簡易版）
   - Issuerの確認
   - 有効期限の確認
   - auth0_sub（subject）の取得
   ↓
3. コンテキストに auth0_sub を保存
   ↓
4. Handler: auth0_sub の取得
   ↓
5. UseCase: ユーザー検索または作成
   ├─ 既存ユーザーの場合
   │  └→ FindByAuth0UserID() でユーザー取得
   │
   └─ 新規ユーザーの場合
      ├→ 新規Userモデル作成
      │  - id: UUID生成
      │  - auth0_user_id: auth0_sub
      │  - email: 仮メール生成
      │  - name: "New User"
      │  - preferences: "{}"
      └→ Create() でDB保存
   ↓
6. ユーザーIDをコンテキストに保存
   - 後続のAPIリクエストで使用
   ↓
7. レスポンス返却
```

### 仮メールアドレスの生成

Auth0の`sub`から仮のメールアドレスを生成します：

```go
// auth0|69038b5428c69abdd48d4d65
// ↓
// user-69038b5428c69abdd48d4d65@temp.local
```

**注意**: 将来的にはJWTの`email`クレームから実際のメールアドレスを取得する必要があります。

---

## 📁 実装ファイル

| 層 | ファイルパス | 役割 |
|----|-------------|------|
| Handler | `internal/interfaces/http/handler/auth_handler.go` | HTTPリクエスト処理 |
| Middleware | `internal/interfaces/http/middleware/auth.go` | JWT検証 |
| Middleware | `internal/interfaces/http/middleware/context.go` | コンテキスト管理 |
| UseCase | `internal/application/usecase/user/sync_user.go` | ビジネスロジック |
| Repository | `internal/infrastructure/persistence/postgres/user_repository.go` | DB操作 |
| Domain | `internal/domain/model/user.go` | エンティティ定義 |
| Auth | `internal/infrastructure/auth/validator.go` | JWT検証ロジック |
| Auth | `internal/infrastructure/auth/jwks.go` | JWKSキャッシュ管理 |

---

## 🔐 認証フロー詳細

### 1. フロントエンド側の処理

```typescript
// 1. Auth0でログイン
const { access_token, id_token } = await auth0.login(email, password);

// 2. id_tokenをlocalStorageに保存
localStorage.setItem('id_token', id_token);
localStorage.setItem('access_token', access_token);

// 3. /auth/sync を呼び出し
const response = await fetch('/api/v1/auth/sync', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${id_token}`,  // id_tokenを使用
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({})
});

const { user } = await response.json();
```

### 2. バックエンド側の処理

#### ミドルウェアでの認証
```go
// 1. トークン抽出
token := extractBearerToken(request)

// 2. トークン検証（簡易版）
parsedToken, err := jwt.Parse(token,
    jwt.WithVerify(false),      // 署名検証を無効化
    jwt.WithValidate(false))    // 厳密な検証を無効化

// 3. 基本検証
- Issuerチェック: https://dev-863amkaw1kj03j7y.us.auth0.com/
- 有効期限チェック
- Subjectの存在確認

// 4. auth0_subをコンテキストに保存
context.Set("auth0_sub", token.Subject())

// 5. 既存ユーザーの場合、user_idもコンテキストに保存
user, _ := userRepo.FindByAuth0UserID(auth0_sub)
if user != nil {
    context.Set("user_id", user.ID)
}
```

#### ハンドラーでの処理
```go
// 1. auth0_sub取得
auth0Sub := middleware.GetAuth0SubFromContext(c)

// 2. ユーザー同期
user := userUseCase.SyncUserByAuth0Sub(ctx, auth0Sub)

// 3. user_idをコンテキストに保存
middleware.SetUserID(c, user.ID)

// 4. レスポンス
return user
```

---

## 🧪 テストケース

### 正常系

- [x] **TC-AU-001-01**: 新規ユーザー作成
  - 前提: `auth0_user_id`がDBに存在しない
  - 期待結果: 200 OK、新規ユーザーが作成される

- [x] **TC-AU-001-02**: 既存ユーザー取得
  - 前提: `auth0_user_id`がDBに存在する
  - 期待結果: 200 OK、既存ユーザー情報が返される

### 異常系（認証）

- [x] **TC-AU-001-03**: JWTトークンなし
  - 期待結果: 401 Unauthorized

- [x] **TC-AU-001-04**: JWTトークンが無効
  - 期待結果: 401 Unauthorized

- [x] **TC-AU-001-05**: トークンにsubクレームがない
  - 期待結果: 401 Unauthorized

- [x] **TC-AU-001-06**: トークンの有効期限切れ
  - 期待結果: 401 Unauthorized

### 異常系（サーバーエラー）

- [ ] **TC-AU-001-07**: DB接続エラー
  - 期待結果: 500 Internal Server Error

- [ ] **TC-AU-001-08**: ユーザー作成失敗（DB制約違反）
  - 期待結果: 500 Internal Server Error

---

## 📊 実装状況

- [x] ドキュメント作成
- [x] 設計レビュー
- [x] 実装完了
- [x] JWT認証実装（簡易版）
- [x] Auth0統合
- [x] ミドルウェア実装
- [x] 動作確認（手動テスト）
- [ ] 単体テスト
- [ ] 統合テスト
- [ ] コードレビュー
- [ ] 本番デプロイ

### 実装済み機能
✅ Auth0統合（`audience`なし構成）  
✅ JWT検証（簡易版：署名検証無効化）  
✅ id_token使用（JWT形式）  
✅ ユーザー自動作成  
✅ 既存ユーザー取得  
✅ コンテキストへのユーザーID保存  
✅ エラーハンドリング  

### 現在の制限事項・今後の改善点

#### 🔴 セキュリティ関連
- **署名検証が無効化されている**: 開発環境で`audience`を使用しない構成のため、JWT署名検証を無効化しています。本番環境では以下のいずれかが必要：
  - Auth0でAPIを作成し`audience`を有効化 → 署名検証を有効化
  - または、id_tokenの署名をJWKSで検証する実装を追加

#### 🟡 機能追加予定
- **実際のメールアドレス取得**: 現在は仮メール。JWTの`email`クレームから実際のメールを取得
- **プロフィール情報の同期**: `name`, `picture`などのクレームを反映
- **ユーザー更新機能**: 既存ユーザーの情報を最新に更新

#### 🟢 パフォーマンス最適化
- **キャッシュ**: ユーザー情報のキャッシュ（Redis等）
- **非同期処理**: ユーザー作成時の追加処理を非同期化

---

## 💡 実装時の注意点

### セキュリティ

1. **トークンの取り扱い**
   - フロントエンドでは`id_token`を使用（JWT形式）
   - `access_token`はJWE形式のため使用しない
   - トークンは必ずHTTPS経由で送信

2. **auth0_user_idの一意性**
   - DBの`users.auth0_user_id`カラムにUNIQUE制約
   - 重複作成を防ぐ

3. **コンテキストの使用**
   - `user_id`は必ずコンテキスト経由で取得
   - リクエストボディから受け取らない

### エラーハンドリング

1. **認証エラー**
   - 詳細なエラーメッセージはログのみ
   - ユーザーには汎用的なメッセージ

2. **DB エラー**
   - スタックトレースをログに記録
   - ユーザーには「同期失敗」のみ通知

### パフォーマンス

1. **ミドルウェアでのユーザー検索**
   - 認証の度にDBアクセスが発生
   - 将来的にはキャッシュの導入を検討

2. **新規ユーザー作成**
   - トランザクション不要（単一INSERT）
   - UUID生成はメモリ上で完結

---

## 🔗 関連ドキュメント

### API
- すべての保護されたAPIエンドポイント（このAPIで取得したuser_idを使用）
- [KN-001: ナレッジ作成](./KN-001_create_knowledge.md)
- [RV-001: コードレビュー実行](./RV-001_review_code.md)

### 認証関連
- [Auth0統合ガイド](../AUTH0_INTEGRATION.md)
- [認証フロー設計](../DESIGN.md#認証)

---

## 📝 変更履歴

| 日付 | バージョン | 変更内容 | 担当 |
|------|-----------|---------|------|
| 2025-01-XX | 1.0 | 初版作成 | - |
| 2025-01-XX | 1.1 | 実装完了（簡易JWT検証版） | - |

---

## 🎯 フロントエンド実装ガイド

### ログインフローの実装

```typescript
// 1. ログイン
async function login(email: string, password: string) {
  // Auth0でログイン
  const response = await axios.post(
    `https://${AUTH0_DOMAIN}/oauth/token`,
    {
      grant_type: 'http://auth0.com/oauth/grant-type/password-realm',
      username: email,
      password: password,
      client_id: AUTH0_CLIENT_ID,
      client_secret: AUTH0_CLIENT_SECRET,
      realm: 'Username-Password-Authentication',
      scope: 'openid profile email',
    }
  );

  const { access_token, id_token } = response.data;

  // トークンを保存
  localStorage.setItem('access_token', access_token);
  localStorage.setItem('id_token', id_token);

  // 2. バックエンドと同期
  await syncUser();
}

// 2. ユーザー同期
async function syncUser() {
  const idToken = localStorage.getItem('id_token');

  const response = await fetch('/api/v1/auth/sync', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${idToken}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({}),
  });

  if (!response.ok) {
    throw new Error('ユーザー同期に失敗しました');
  }

  const { user } = await response.json();
  return user;
}
```

### すべてのAPIリクエストでid_tokenを使用

```typescript
// APIクライアント
export const getAuthHeaders = (): HeadersInit => {
  // id_tokenを優先的に使用
  const idToken = localStorage.getItem('id_token');
  const accessToken = localStorage.getItem('access_token');
  
  // id_tokenがJWT形式（3パーツ）かチェック
  const token = idToken && idToken.split('.').length === 3 
    ? idToken 
    : accessToken;
  
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  return headers;
};
```
