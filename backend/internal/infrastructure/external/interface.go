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
