package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config - アプリケーション全体の設定
type Config struct {
	Env      string
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	LLM      LLMConfig
	Redis    RedisConfig
	Features FeatureFlags
}

// ServerConfig - サーバー設定
type ServerConfig struct {
	Port     string
	LogLevel string
}

// DatabaseConfig - データベース設定
type DatabaseConfig struct {
	URL             string
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// AuthConfig - 認証設定（Auth0）
type AuthConfig struct {
	Domain       string
	Audience     string
	ClientID     string
	ClientSecret string
}

// LLMConfig - LLM設定
type LLMConfig struct {
	ClaudeAPIKey    string
	ClaudeModel     string
	ClaudeMaxTokens int
	OpenAIAPIKey    string
	OpenAIEmbedding string
	EmbeddingDim    int
	OpenAITimeout   time.Duration
}

// RedisConfig - Redis設定
type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

// FeatureFlags - 機能フラグ
type FeatureFlags struct {
	VectorSearch         bool
	HybridSearch         bool
	AutoKnowledgeExtract bool
	ConversationMode     bool
}

// Load - 環境変数から設定を読み込み
func Load() (*Config, error) {
	cfg := &Config{
		Env: getEnv("ENV", "development"),
		Server: ServerConfig{
			Port:     getEnv("PORT", "8080"),
			LogLevel: getEnv("LOG_LEVEL", "debug"),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", ""),
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "dev_user"),
			Password:        getEnv("DB_PASSWORD", "dev_password"),
			Name:            getEnv("DB_NAME", "reviewapp"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", "5m"),
		},
		Auth: AuthConfig{
			Domain:       getEnv("AUTH0_DOMAIN", ""),
			Audience:     getEnv("AUTH0_AUDIENCE", ""),
			ClientID:     getEnv("AUTH0_CLIENT_ID", ""),
			ClientSecret: getEnv("AUTH0_CLIENT_SECRET", ""),
		},
		LLM: LLMConfig{
			ClaudeAPIKey:    getEnv("CLAUDE_API_KEY", ""),
			ClaudeModel:     getEnv("CLAUDE_MODEL", "claude-3-5-haiku-latest"),
			ClaudeMaxTokens: getEnvAsInt("CLAUDE_MAX_TOKENS", 4096),
			OpenAIAPIKey:    getEnv("OPENAI_API_KEY", ""),
			OpenAIEmbedding: getEnv("OPENAI_EMBEDDING_MODEL", "text-embedding-3-small"),
			EmbeddingDim:    getEnvAsInt("OPENAI_EMBEDDING_DIMENSIONS", 1536),
			OpenAITimeout:   getEnvAsDuration("OPENAI_API_TIMEOUT", "30s"),
		},
		Redis: RedisConfig{
			URL:      getEnv("REDIS_URL", "redis://localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		Features: FeatureFlags{
			VectorSearch:         getEnvAsBool("FEATURE_VECTOR_SEARCH", false),
			HybridSearch:         getEnvAsBool("FEATURE_HYBRID_SEARCH", false),
			AutoKnowledgeExtract: getEnvAsBool("FEATURE_AUTO_KNOWLEDGE_EXTRACT", true),
			ConversationMode:     getEnvAsBool("FEATURE_CONVERSATION_MODE", true),
		},
	}

	return cfg, nil
}

// GetDSN - PostgreSQL接続文字列を生成
func (c *DatabaseConfig) GetDSN() string {
	// DATABASE_URLが設定されていればそれを使用
	if c.URL != "" {
		return c.URL
	}

	// 個別設定から生成
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

// ヘルパー関数

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	valueStr := getEnv(key, defaultValue)
	if duration, err := time.ParseDuration(valueStr); err == nil {
		return duration
	}
	// パースに失敗したらデフォルト値をパース
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}
