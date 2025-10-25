# =====================================================
# ReviewApp - Makefile
# =====================================================

.PHONY: help setup up down restart logs clean test lint format migrate db-shell redis-shell

# デフォルトターゲット
.DEFAULT_GOAL := help

# カラー定義
BLUE := \033[0;34m
GREEN := \033[0;32m
YELLOW := \033[1;33m
NC := \033[0m # No Color

# 環境検出（Docker内 or ホスト）
# Docker内では /proc/1/cgroup にdockerの文字列がある
IN_DOCKER := $(shell test -f /.dockerenv && echo "yes" || echo "no")

# データベース接続設定
ifeq ($(IN_DOCKER),yes)
    # Docker内からの接続
    DB_HOST := postgres
    REDIS_HOST := redis
else
    # ホストマシンからの接続
    DB_HOST := localhost
    REDIS_HOST := localhost
endif

DB_PORT := 5432
DB_USER := dev_user
DB_PASSWORD := dev_password
DB_NAME := reviewapp

## =====================================================
## ヘルプ
## =====================================================

help: ## このヘルプを表示
	@echo "$(BLUE)ReviewApp - Available Commands$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)Environment: $(IN_DOCKER)$(NC)"
	@echo "$(YELLOW)DB Host: $(DB_HOST)$(NC)"

## =====================================================
## セットアップ & 起動
## =====================================================

setup: ## 初期セットアップを実行
	@echo "$(BLUE)🚀 Running setup...$(NC)"
	@bash scripts/setup-dev.sh

up: ## Docker コンテナを起動
	@echo "$(BLUE)🐳 Starting containers...$(NC)"
	@docker-compose up -d postgres redis
	@echo "$(GREEN)✓ Containers started$(NC)"
	@echo "$(YELLOW)Waiting for PostgreSQL to be ready...$(NC)"
	@sleep 5

up-tools: ## DBツール込みで起動（pgAdmin + Redis Commander）
	@echo "$(BLUE)🐳 Starting containers with tools...$(NC)"
	@docker-compose --profile tools up -d
	@echo "$(GREEN)✓ Containers with tools started$(NC)"
	@echo "$(YELLOW)pgAdmin: http://localhost:5050 (admin@example.com / admin)$(NC)"
	@echo "$(YELLOW)Redis Commander: http://localhost:8081$(NC)"

down: ## Docker コンテナを停止
	@echo "$(BLUE)🛑 Stopping containers...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✓ Containers stopped$(NC)"

restart: down up ## コンテナを再起動

logs: ## コンテナのログを表示
	@docker-compose logs -f

logs-api: ## APIサーバーのログを表示
	@docker-compose logs -f dev

logs-db: ## PostgreSQLのログを表示
	@docker-compose logs -f postgres

logs-redis: ## Redisのログを表示
	@docker-compose logs -f redis

clean: down ## コンテナとボリュームを完全削除
	@echo "$(YELLOW)⚠️  Removing all containers and volumes...$(NC)"
	@docker-compose down -v
	@echo "$(GREEN)✓ Cleanup complete$(NC)"

## =====================================================
## データベース
## =====================================================

migrate: ## マイグレーションを実行
	@echo "$(BLUE)🗄️  Running migrations...$(NC)"
	@echo "$(YELLOW)Connecting to $(DB_HOST):$(DB_PORT)$(NC)"
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f backend/migrations/001_init.sql
	@echo "$(GREEN)✓ Migrations complete$(NC)"

migrate-reset: ## データベースをリセットして再マイグレーション
	@echo "$(YELLOW)⚠️  Resetting database...$(NC)"
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d postgres -c "CREATE DATABASE $(DB_NAME);"
	@make migrate
	@echo "$(GREEN)✓ Database reset complete$(NC)"

db-shell: ## PostgreSQLに接続
	@echo "$(BLUE)🐘 Connecting to PostgreSQL...$(NC)"
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME)

db-status: ## データベースの状態を確認
	@echo "$(BLUE)📊 Database Status$(NC)"
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -c "\dt"

redis-shell: ## Redisに接続
	@redis-cli -h $(REDIS_HOST) -p 6379

## =====================================================
## 開発
## =====================================================

run: ## APIサーバーを起動（ホットリロードなし）
	@echo "$(BLUE)🚀 Starting API server...$(NC)"
	@cd backend && go run cmd/api/main.go

dev: ## APIサーバーを起動（Air でホットリロード）
	@echo "$(BLUE)🔥 Starting API server with hot reload...$(NC)"
	@cd backend && air

build: ## バイナリをビルド
	@echo "$(BLUE)🔨 Building binary...$(NC)"
	@cd backend && go build -o ../bin/api cmd/api/main.go
	@echo "$(GREEN)✓ Binary built: bin/api$(NC)"

## =====================================================
## テスト & 品質
## =====================================================

test: ## テストを実行
	@echo "$(BLUE)🧪 Running tests...$(NC)"
	@cd backend && go test -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)✓ Tests complete$(NC)"

test-coverage: test ## テストカバレッジを表示
	@cd backend && go tool cover -html=coverage.out

lint: ## Lintを実行
	@echo "$(BLUE)🔍 Running linter...$(NC)"
	@cd backend && golangci-lint run ./...
	@echo "$(GREEN)✓ Lint complete$(NC)"

format: ## コードをフォーマット
	@echo "$(BLUE)✨ Formatting code...$(NC)"
	@cd backend && gofmt -w -s .
	@cd backend && goimports -w .
	@echo "$(GREEN)✓ Format complete$(NC)"

## =====================================================
## Go Modules
## =====================================================

mod-download: ## Go modulesをダウンロード
	@cd backend && go mod download

mod-tidy: ## Go modulesを整理
	@cd backend && go mod tidy

mod-verify: ## Go modulesを検証
	@cd backend && go mod verify

## =====================================================
## 便利コマンド
## =====================================================

ps: ## コンテナの状態を確認
	@docker-compose ps

stats: ## リソース使用状況を表示
	@docker stats --no-stream

health: ## サービスのヘルスチェック
	@echo "$(BLUE)🏥 Health Check$(NC)"
	@echo ""
	@echo "PostgreSQL:"
	@PGPASSWORD=$(DB_PASSWORD) pg_isready -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) && echo "$(GREEN)✓ OK$(NC)" || echo "$(YELLOW)✗ Not Ready$(NC)"
	@echo ""
	@echo "Redis:"
	@redis-cli -h $(REDIS_HOST) -p 6379 ping > /dev/null 2>&1 && echo "$(GREEN)✓ OK$(NC)" || echo "$(YELLOW)✗ Not Ready$(NC)"
	@echo ""

seed: ## サンプルデータを投入
	@echo "$(BLUE)🌱 Seeding database...$(NC)"
	@cd backend && go run scripts/seed/main.go
	@echo "$(GREEN)✓ Seed complete$(NC)"

## =====================================================
## ドキュメント
## =====================================================

docs: ## API ドキュメントを生成
	@echo "$(BLUE)📚 Generating API docs...$(NC)"
	@cd backend && swag init -g cmd/api/main.go -o docs
	@echo "$(GREEN)✓ Docs generated$(NC)"

## =====================================================
## リリース
## =====================================================

docker-build: ## Dockerイメージをビルド
	@echo "$(BLUE)🐳 Building Docker image...$(NC)"
	@docker build -t reviewapp:latest -f backend/Dockerfile backend/
	@echo "$(GREEN)✓ Docker image built$(NC)"

## =====================================================
## 情報表示
## =====================================================

info: ## プロジェクト情報を表示
	@echo "$(BLUE)ReviewApp - Project Info$(NC)"
	@echo ""
	@echo "Environment:     $(IN_DOCKER)"
	@echo "Go version:      $(shell go version 2>/dev/null || echo 'Not installed')"
	@echo "Project path:    $(shell pwd)"
	@echo "Docker version:  $(shell docker --version 2>/dev/null || echo 'Not installed')"
	@echo ""
	@echo "Connection Info:"
	@echo "  - DB Host:     $(DB_HOST):$(DB_PORT)"
	@echo "  - Redis Host:  $(REDIS_HOST):6379"
	@echo ""
	@echo "Services:"
	@echo "  - API:         http://localhost:8080"
	@echo "  - PostgreSQL:  $(DB_HOST):$(DB_PORT)"
	@echo "  - Redis:       $(REDIS_HOST):6379"
	@echo ""

## =====================================================
## Docker開発環境
## =====================================================

shell: ## 開発コンテナに入る
	@docker-compose exec dev /bin/sh

exec: ## 開発コンテナでコマンドを実行（使用: make exec CMD="go version"）
	@docker-compose exec dev $(CMD)