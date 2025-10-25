-- =====================================================
-- ReviewApp - Initial Schema Migration
-- =====================================================
-- Phase 1: キーワード検索ベース + Auth0認証
-- Phase 2: ベクトル検索対応（pgvector）
-- =====================================================

-- pgvector拡張を有効化（Phase 2用）
CREATE EXTENSION IF NOT EXISTS vector;

-- =====================================================
-- 1. users テーブル
-- =====================================================
-- Auth0で認証を行うため、パスワードは不要
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Auth0連携
    auth0_user_id VARCHAR(255) NOT NULL UNIQUE,  -- Auth0のユーザーID（例: auth0|123456）
    
    -- 基本情報（Auth0から取得）
    email VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    avatar_url TEXT,  -- Auth0のpicture
    
    -- アプリ独自設定
    preferences JSONB DEFAULT '{}',  -- ユーザー設定
    
    -- メタデータ
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- インデックス
CREATE INDEX idx_users_auth0_user_id ON users(auth0_user_id);
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at);

-- =====================================================
-- 2. knowledge テーブル（ナレッジ - コア）
-- =====================================================
-- ユーザーのコーディング哲学・ルール・学習履歴
CREATE TABLE knowledge (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- コンテンツ
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    
    -- 分類
    category VARCHAR(50),  -- 'error_handling', 'testing', 'performance', etc.
    priority INTEGER NOT NULL DEFAULT 3 CHECK (priority BETWEEN 1 AND 5),
    
    -- メタデータ
    source_type VARCHAR(50) NOT NULL,  -- 'review', 'conversation', 'manual'
    source_id UUID,  -- 元のreview_id等
    
    -- 使用状況（追加）
    usage_count INTEGER NOT NULL DEFAULT 0,  -- 何回参照されたか
    last_used_at TIMESTAMP WITH TIME ZONE,    -- 最後に使われた日時
    
    -- ベクトル検索用（Phase 2）
    embedding vector(1536),  -- OpenAI text-embedding-3-small の次元数
    
    -- 状態
    is_active BOOLEAN NOT NULL DEFAULT true,
    
    -- タイムスタンプ
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- インデックス
CREATE INDEX idx_knowledge_user_id ON knowledge(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_knowledge_category ON knowledge(category) WHERE is_active = true AND deleted_at IS NULL;
CREATE INDEX idx_knowledge_priority ON knowledge(priority DESC, created_at DESC);
CREATE INDEX idx_knowledge_source ON knowledge(source_type, source_id);
CREATE INDEX idx_knowledge_usage ON knowledge(usage_count DESC, last_used_at DESC);  -- 使用頻度順

-- 全文検索用インデックス（Phase 1: キーワード検索）
CREATE INDEX idx_knowledge_title_content ON knowledge USING gin(to_tsvector('english', title || ' ' || content)) 
    WHERE is_active = true AND deleted_at IS NULL;

-- ベクトル検索用インデックス（Phase 2）
-- HNSW: 高速な近似最近傍探索
CREATE INDEX idx_knowledge_embedding ON knowledge USING hnsw (embedding vector_cosine_ops)
    WHERE embedding IS NOT NULL AND deleted_at IS NULL;

-- =====================================================
-- 3. tags テーブル（タグマスタ）
-- =====================================================
-- Phase 1ではテーブルのみ作成、API実装はPhase 2
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    color VARCHAR(7),  -- Hex color: #FF5733
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, name)
);

CREATE INDEX idx_tags_user_id ON tags(user_id);

-- =====================================================
-- 4. knowledge_tags テーブル（中間テーブル）
-- =====================================================
-- Phase 1ではテーブルのみ作成、API実装はPhase 2
CREATE TABLE knowledge_tags (
    knowledge_id UUID NOT NULL REFERENCES knowledge(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (knowledge_id, tag_id)
);

CREATE INDEX idx_knowledge_tags_tag_id ON knowledge_tags(tag_id);

-- =====================================================
-- 5. reviews テーブル（レビュー履歴）
-- =====================================================
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- レビュー対象
    code TEXT NOT NULL,
    language VARCHAR(50) NOT NULL,  -- 'go', 'python', 'javascript', etc.
    context TEXT,
    
    -- レビュー結果
    review_result TEXT NOT NULL,
    
    -- LLM情報
    llm_provider VARCHAR(50) NOT NULL,  -- 'claude', 'openai'
    llm_model VARCHAR(100) NOT NULL,    -- 'claude-3-5-sonnet-20241022'
    tokens_used INTEGER,
    
    -- フィードバック
    feedback_score INTEGER CHECK (feedback_score BETWEEN 1 AND 5),
    feedback_comment TEXT,
    
    -- タイムスタンプ
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- インデックス
CREATE INDEX idx_reviews_user_id ON reviews(user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_reviews_created_at ON reviews(created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_reviews_language ON reviews(language) WHERE deleted_at IS NULL;
CREATE INDEX idx_reviews_feedback ON reviews(feedback_score) WHERE feedback_score IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_reviews_deleted_at ON reviews(deleted_at) WHERE deleted_at IS NULL;

-- =====================================================
-- 6. review_knowledge テーブル（中間テーブル）
-- =====================================================
-- レビュー時に参照されたナレッジを記録
CREATE TABLE review_knowledge (
    review_id UUID NOT NULL REFERENCES reviews(id) ON DELETE CASCADE,
    knowledge_id UUID NOT NULL REFERENCES knowledge(id) ON DELETE CASCADE,
    
    -- 関連度スコア（Phase 2で使用）
    relevance_score FLOAT,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (review_id, knowledge_id)
);

CREATE INDEX idx_review_knowledge_knowledge_id ON review_knowledge(knowledge_id);
CREATE INDEX idx_review_knowledge_relevance ON review_knowledge(relevance_score DESC) 
    WHERE relevance_score IS NOT NULL;

-- =====================================================
-- トリガー: updated_at自動更新
-- =====================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- トリガーを各テーブルに適用
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_knowledge_updated_at BEFORE UPDATE ON knowledge
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reviews_updated_at BEFORE UPDATE ON reviews
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================
-- シードデータ（開発用）
-- =====================================================
-- テスト用ユーザー（Auth0 ID はダミー）
INSERT INTO users (id, auth0_user_id, email, name) VALUES
    ('00000000-0000-0000-0000-000000000001', 'auth0|dev-user-001', 'dev@example.com', 'Dev User');

-- サンプルナレッジ
INSERT INTO knowledge (user_id, title, content, category, priority, source_type, usage_count) VALUES
    ('00000000-0000-0000-0000-000000000001', 
     'エラーハンドリングの原則', 
     'エラーは必ずログに出力し、ユーザー向けメッセージと開発者向け詳細を分ける。contextを使ってエラーチェーンを保持する。',
     'error_handling',
     5,
     'manual',
     0),
    ('00000000-0000-0000-0000-000000000001',
     '関数は1つのことだけをする',
     '関数は50行以内に抑え、1つの責務のみを持つ。複雑な処理は小さな関数に分割する。',
     'clean_code',
     4,
     'manual',
     0),
    ('00000000-0000-0000-0000-000000000001',
     'テストは必須',
     '新機能には必ずユニットテストを書く。カバレッジは80%以上を目標とする。',
     'testing',
     4,
     'manual',
     0);

-- =====================================================
-- コメント
-- =====================================================
COMMENT ON TABLE users IS 'ユーザー情報（Auth0で認証）';
COMMENT ON TABLE knowledge IS 'ナレッジ（ユーザーのコーディング哲学・ルール）';
COMMENT ON TABLE tags IS 'タグマスタ（Phase 2で実装）';
COMMENT ON TABLE knowledge_tags IS 'ナレッジとタグの紐付け（Phase 2で実装）';
COMMENT ON TABLE reviews IS 'コードレビュー履歴';
COMMENT ON TABLE review_knowledge IS 'レビュー時に参照されたナレッジ';

COMMENT ON COLUMN users.auth0_user_id IS 'Auth0のユーザーID（例: auth0|123456）';
COMMENT ON COLUMN knowledge.embedding IS 'OpenAI Embedding vector (1536次元) - Phase 2で使用';
COMMENT ON COLUMN knowledge.priority IS '重要度 (1=低, 5=高)';
COMMENT ON COLUMN knowledge.source_type IS 'ナレッジの出所 (review/conversation/manual)';
COMMENT ON COLUMN knowledge.usage_count IS '参照された回数';
COMMENT ON COLUMN knowledge.last_used_at IS '最後に使われた日時';
