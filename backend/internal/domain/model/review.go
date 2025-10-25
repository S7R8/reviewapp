package model

import (
	"time"

	"github.com/google/uuid"
)

// Review - コードレビューエンティティ
type Review struct {
	ID                  string     `json:"id"`
	UserID              string     `json:"user_id"`
	Code                string     `json:"code"`     // レビュー対象のコード
	Language            string     `json:"language"` // プログラミング言語
	Context             string     `json:"context,omitempty"`  // 追加コンテキスト（オプション）
	ReviewResult        string     `json:"review_result"`        // AIのレビュー結果
	ReferencedKnowledge []string   `json:"referenced_knowledge"` // 参照したナレッジID（中間テーブルから取得）
	LLMProvider         string     `json:"llm_provider"`         // LLMプロバイダー（claude, openai）
	LLMModel            string     `json:"llm_model"`            // モデル名
	TokensUsed          int        `json:"tokens_used"`          // 使用トークン数
	FeedbackScore       *int       `json:"feedback_score,omitempty"`       // ユーザーのフィードバック（1-5）
	FeedbackComment     string     `json:"feedback_comment,omitempty"`     // フィードバックコメント
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"` // 論理削除
}

// NewReview - レビューを生成
func NewReview(userID, code, language, context string) *Review {
	return &Review{
		ID:        uuid.New().String(),
		UserID:    userID,
		Code:      code,
		Language:  language,
		Context:   context,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// SetReviewResult - レビュー結果を設定
func (r *Review) SetReviewResult(result string, knowledgeIDs []string, llmProvider, llmModel string, tokensUsed int) {
	r.ReviewResult = result
	r.ReferencedKnowledge = knowledgeIDs
	r.LLMProvider = llmProvider
	r.LLMModel = llmModel
	r.TokensUsed = tokensUsed
	r.UpdatedAt = time.Now()
}

// SetFeedback - ユーザーフィードバックを設定
func (r *Review) SetFeedback(score int) {
	r.FeedbackScore = &score
	r.UpdatedAt = time.Now()
}
