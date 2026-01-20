package external

import (
	"context"
)

// ClaudeClientInterface - Claude API クライアントのインターフェース
type ClaudeClientInterface interface {
	ReviewCode(ctx context.Context, input ReviewCodeInput) (*ReviewCodeOutput, error)
}

// 元のClaudeClientがインターフェースを実装していることを保証
var _ ClaudeClientInterface = (*ClaudeClient)(nil)

type EmbeddingClientInterface interface {
	// GenerateEmbedding - テキストからEmbeddingを生成する
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
	// GenerateEmbeddings - 複数のテキストからEmbeddingを生成する
	GenerateEmbeddings(ctx context.Context, text []string) ([][]float32, error)
}

// OpenAIClientがインターフェースを実装していることを保証
var _ EmbeddingClientInterface = (*OpenAIClient)(nil)
