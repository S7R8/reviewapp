# ReviewApp

RAGベースのAIコードレビューツール。ユーザーが登録したナレッジを参照し、一貫性のあるレビューを提供します。

![Demo](C:\output-palette.gif)

## 概要

Claude APIを利用したコードレビューアプリケーションです。

※このプロジェクトは東京都デジタル人材育成支援事業「若手エンジニアコース」の単位取得の一環として開発しました。

### 開発の経緯

Claude等のAIツールを用いたコーディング支援において、メモリ容量の制限により前回のレビュー内容が失われ、AIのフィードバックに一貫性が保てないという悩みがありました。この悩みに対し、RAGを用いてナレッジベースを構築し、過去の判断基準を参照することで一貫性のあるコードレビューを実現することを目的としています。

### 機能概要

- **ナレッジベースのレビュー**: ユーザーが登録した判断基準を参照
- **構造化された結果**: Summary、Good Points、Improvements を分類表示
- **フィードバック機能**: レビュー結果への評価による学習
- **ナレッジ管理**: カテゴリ・優先度別の整理

## 技術スタック

**Backend**
- Go 1.25.2 + Echo v4
- PostgreSQL 15 (pgvector)
- Claude 3.5 Sonnet API

**Frontend**
- React 19 + TypeScript
- Vite 7
- Zustand (状態管理)
- Tailwind CSS 4
- Monaco Editor

**Architecture**
- Clean Architecture (Backend)
- Docker (開発環境)

## セットアップ

### 前提条件

- Docker Desktop
- Go 1.25.2+
- Node.js 18.0+

### 環境変数の設定

`backend/.env.example` を `backend/.env` にコピーして以下を設定:
```env
# Anthropic API Key (https://console.anthropic.com/ で取得)
ANTHROPIC_API_KEY=your_api_key_here

# OpenAI API Key (https://platform.openai.com/ で取得)
OPENAI_API_KEY=your_openai_api_key_here

# Auth0 設定 (https://auth0.com/ でアプリケーション作成)
AUTH0_DOMAIN=your_domain.auth0.com
AUTH0_CLIENT_ID=your_client_id

# Database (デフォルト値でOK)
DATABASE_URL=postgres://dev_user:dev_password@postgres:5432/reviewapp?sslmode=disable
REDIS_URL=redis://redis:6379
```

### 起動
```bash
# 初回セットアップ (コンテナ起動 + DB migration)
make setup

# バックエンド起動 (ホットリロード)
cd backend && air

# フロントエンド起動 (別ターミナル)
cd frontend && npm install && npm run dev
```

**アクセス先:**
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080

## プロジェクト構造
```
review-app/
├── backend/              # Go (Clean Architecture)
│   ├── cmd/              # エントリーポイント
│   ├── internal/
│   │   ├── application/  # ユースケース
│   │   ├── domain/       # ドメインロジック
│   │   ├── infrastructure/ # DB, API実装
│   │   └── interfaces/   # HTTPハンドラー
│   └── migrations/       # DBマイグレーション
│
├── frontend/             # React + TypeScript
│   └── src/
│       ├── components/   # UIコンポーネント
│       ├── pages/        # ページ
│       └── api/          # APIクライアント
│
└── docker-compose.yaml   # 開発環境
```

## 開発コマンド
```bash
make help           # コマンド一覧
make up             # コンテナ起動
make down           # コンテナ停止
make migrate        # マイグレーション実行
make test           # テスト実行
make db-shell       # PostgreSQL接続
```

## システム構成図
``` mermaid
graph TB
    subgraph "フロントエンド"
        A[React + TypeScript]
    end
    
    subgraph "バックエンド (Go + Echo)"
        C[認証ミドルウェア<br/>Auth0 JWT]
        C --> D[ハンドラー層]
        D --> E[ユースケース層]
        E --> F[ドメイン層]
        E --> G[インフラ層]
    end
    
    subgraph "外部API"
        H[Claude API<br/>コードレビュー生成]
        I[OpenAI API<br/>Embedding生成]
    end
    
    subgraph "データストア"
        J[(PostgreSQL + pgvector)]
    end
    
    subgraph "認証"
        L[Auth0<br/>認証・認可]
    end
    
    A -->|REST API| C
    A -->|OAuth 2.0| L
    E -->|レビュー生成| H
    E -->|Embedding生成| I
    G -->|ベクトル検索| J
    C -->|JWT検証| L
    
    style H fill:#CC9B7A,stroke:#333,stroke-width:2px
    style I fill:#10a37f,stroke:#333,stroke-width:2px
    style J fill:#336791,stroke:#333,stroke-width:2px
```
