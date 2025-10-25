#!/bin/bash
set -e

echo "🚀 ReviewApp - Development Environment Setup"
echo "=================================================="

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 関数定義
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# .envファイルの確認と作成
echo ""
echo "📝 Step 1: Environment Variables"
if [ ! -f .env ]; then
    print_info "Creating .env file from .env.example..."
    cp .env.example .env
    print_success ".env file created"
    print_warning "Please update .env with your API keys!"
else
    print_success ".env file already exists"
fi

# Go modules の初期化
echo ""
echo "📦 Step 2: Go Modules"
if [ ! -f backend/go.mod ]; then
    print_info "Initializing Go modules..."
    cd backend
    go mod init github.com/s7r8/reviewapp
    cd ..
    print_success "Go modules initialized"
else
    print_success "Go modules already initialized"
fi

# データベースの起動確認
echo ""
echo "🐘 Step 3: Database Check"
print_info "Checking PostgreSQL connection..."

# 最大30秒待機
for i in {1..30}; do
    if pg_isready -h localhost -p 5432 -U dev_user > /dev/null 2>&1; then
        print_success "PostgreSQL is ready"
        break
    fi
    if [ $i -eq 30 ]; then
        print_error "PostgreSQL is not ready after 30 seconds"
        exit 1
    fi
    sleep 1
done

# pgvectorの確認
echo ""
echo "🔍 Step 4: pgvector Extension"
print_info "Checking pgvector extension..."

PGPASSWORD=dev_password psql -h localhost -p 5432 -U dev_user -d reviewapp -c "CREATE EXTENSION IF NOT EXISTS vector;" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    print_success "pgvector extension enabled"
else
    print_error "Failed to enable pgvector extension"
    exit 1
fi

# Redisの起動確認
echo ""
echo "🔴 Step 5: Redis Check"
print_info "Checking Redis connection..."

for i in {1..10}; do
    if redis-cli -h localhost -p 6379 ping > /dev/null 2>&1; then
        print_success "Redis is ready"
        break
    fi
    if [ $i -eq 10 ]; then
        print_error "Redis is not ready after 10 seconds"
        exit 1
    fi
    sleep 1
done

# マイグレーションファイルの確認
echo ""
echo "🗄️  Step 6: Database Migrations"
if [ -f backend/migrations/001_init.sql ]; then
    print_info "Running migrations..."
    PGPASSWORD=dev_password psql -h localhost -p 5432 -U dev_user -d reviewapp -f backend/migrations/001_init.sql > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        print_success "Migrations completed"
    else
        print_warning "Migrations may have already been applied"
    fi
else
    print_warning "No migration files found"
fi

# テーブルの確認
echo ""
echo "📊 Step 7: Verify Tables"
PGPASSWORD=dev_password psql -h localhost -p 5432 -U dev_user -d reviewapp -c "\dt" 2>&1 | grep -q "users"
if [ $? -eq 0 ]; then
    print_success "Database tables created successfully"
    
    # テーブル一覧を表示
    print_info "Created tables:"
    PGPASSWORD=dev_password psql -h localhost -p 5432 -U dev_user -d reviewapp -c "\dt" | grep -E "users|knowledge|reviews|conversations|messages" | awk '{print "  - " $3}'
else
    print_warning "Tables not found - migrations may need to be run manually"
fi

# Go依存関係のインストール
echo ""
echo "📥 Step 8: Installing Go Dependencies"
if [ -f backend/go.mod ]; then
    print_info "Running go mod tidy..."
    cd backend
    go mod tidy
    cd ..
    print_success "Go dependencies installed"
fi

# 完了メッセージ
echo ""
echo "=================================================="
print_success "Setup completed successfully!"
echo ""
echo "🎯 Next Steps:"
echo ""
echo "1. Update your .env file with API keys:"
echo "   - AUTH0_DOMAIN, AUTH0_AUDIENCE, etc."
echo "   - CLAUDE_API_KEY"
echo "   - OPENAI_API_KEY"
echo ""
echo "2. Start developing:"
echo "   cd backend"
echo "   go run cmd/api/main.go"
echo ""
echo "   Or use Air for hot reload:"
echo "   air"
echo ""
echo "3. Access the services:"
echo "   - API: http://localhost:8080"
echo "   - PostgreSQL: localhost:5432"
echo "   - Redis: localhost:6379"
echo "   - pgAdmin: http://localhost:5050 (optional, run with --profile tools)"
echo ""
echo "4. Run tests:"
echo "   make test"
echo ""
echo "=================================================="