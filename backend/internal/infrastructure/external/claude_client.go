package external

import (
	"context"
	"fmt"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// ClaudeClient - Claude API クライアント
type ClaudeClient struct {
	client      anthropic.Client
	model       string
	maxTokens   int
	temperature float64
}

// NewClaudeClient - コンストラクタ
func NewClaudeClient(apiKey, model string, maxTokens int) *ClaudeClient {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &ClaudeClient{
		client:      client,
		model:       model,
		maxTokens:   maxTokens,
		temperature: 0.7, // デフォルト
	}
}

// ReviewCodeInput - レビュー入力
type ReviewCodeInput struct {
	Code            string
	Language        string
	Context         string
	KnowledgePrompt string // ナレッジから生成したプロンプト
}

// ReviewCodeOutput - レビュー結果
type ReviewCodeOutput struct {
	ReviewResult string
	TokensUsed   int
}

// ReviewCode - コードをレビュー
func (c *ClaudeClient) ReviewCode(ctx context.Context, input ReviewCodeInput) (*ReviewCodeOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// プロンプト生成
	systemPrompt := c.buildSystemPrompt(input.KnowledgePrompt)
	userPrompt := c.buildUserPrompt(input.Code, input.Language, input.Context)

	// Claude API呼び出し
	message, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(c.model),
		MaxTokens: int64(c.maxTokens),
		System: []anthropic.TextBlockParam{
			{
				Type: "text",
				Text: systemPrompt,
			},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock(userPrompt),
			),
		},
		Temperature: anthropic.Float(c.temperature),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call Claude API: %w", err)
	}

	// レスポンスからテキストを抽出
	var reviewText string
	for _, block := range message.Content {
		if block.Type == "text" {
			reviewText += block.Text
		}
	}

	// トークン使用量を計算
	tokensUsed := int(message.Usage.InputTokens + message.Usage.OutputTokens)

	return &ReviewCodeOutput{
		ReviewResult: reviewText,
		TokensUsed:   tokensUsed,
	}, nil
}

// buildSystemPrompt - システムプロンプト生成
func (c *ClaudeClient) buildSystemPrompt(knowledgePrompt string) string {
	return fmt.Sprintf(`あなたはコードレビュアーです。
以下のルールと過去の判断基準に基づいてレビューしてください。

## ユーザーのコーディング哲学・ルール
%s

## レビュー指示
1. 上記のルールに違反している箇所を指摘
2. 改善案を具体的に提示
3. なぜそのルールが重要か説明
4. 良い点も必ず指摘する

**重要**: ユーザーの哲学・ルールを最優先してください。

## 出力フォーマット（この形式を厳密に守ること）

**必ず以下の構造で出力してください:**

### 良い点
- 良い点1
- 良い点2

### 1. 改善点のタイトル

- 問題点の説明
- 理由の説明

改善例：
`+"```python"+`
# 改善後のコード
`+"```"+`

### 2. 改善点のタイトル

- 問題点の説明
- 理由の説明

改善例：
`+"```python"+`
# 改善後のコード
`+"```"+`

### 総合評価
総合的な評価を1-2文で記述

**絶対に守るべきルール:**
1. 各セクションは必ず「### 」で始める（###の後にスペース）
2. 改善点は「### 数字. タイトル」の形式
3. コードブロックは`+"```言語名"+`で囲む
4. この順序を必ず守る: 良い点 → 改善点 → 総合評価`, knowledgePrompt)
}

// buildUserPrompt - ユーザープロンプト生成
func (c *ClaudeClient) buildUserPrompt(code, language, context string) string {
	prompt := fmt.Sprintf(`## レビュー対象コード
言語: %s

`, language)

	if context != "" {
		prompt += fmt.Sprintf(`コンテキスト: %s

`, context)
	}

	prompt += fmt.Sprintf("```%s\n%s\n```", language, code)

	return prompt
}
