package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Review - コードレビューエンティティ
type Review struct {
	ID                  string                  `json:"id"`
	UserID              string                  `json:"user_id"`
	Code                string                  `json:"code"`
	Language            string                  `json:"language"`
	Context             string                  `json:"context,omitempty"`
	ReviewResult        string                  `json:"review_result"`                  // マークダウン（元データ）
	StructuredResult    *StructuredReviewResult `json:"structured_result,omitempty"`    // 構造化データ
	ReferencedKnowledge []string                `json:"referenced_knowledge"`
	LLMProvider         string                  `json:"llm_provider"`
	LLMModel            string                  `json:"llm_model"`
	TokensUsed          int                     `json:"tokens_used"`
	FeedbackScore       *int                    `json:"feedback_score,omitempty"`
	FeedbackComment     string                  `json:"feedback_comment,omitempty"`
	CreatedAt           time.Time               `json:"created_at"`
	UpdatedAt           time.Time               `json:"updated_at"`
	DeletedAt           *time.Time              `json:"deleted_at,omitempty"`
}

// StructuredReviewResult - 構造化されたレビュー結果
type StructuredReviewResult struct {
	Summary      string        `json:"summary"`
	GoodPoints   []string      `json:"good_points"`
	Improvements []Improvement `json:"improvements"`
}

// Improvement - 改善点
type Improvement struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CodeAfter   string `json:"code_after,omitempty"`
	Severity    string `json:"severity"` // low, medium, high
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
func (r *Review) SetReviewResult(result string, structuredResult *StructuredReviewResult, knowledgeIDs []string, llmProvider, llmModel string, tokensUsed int) {
	r.ReviewResult = result
	r.StructuredResult = structuredResult
	r.ReferencedKnowledge = knowledgeIDs
	r.LLMProvider = llmProvider
	r.LLMModel = llmModel
	r.TokensUsed = tokensUsed
	r.UpdatedAt = time.Now()
}

// SetFeedback - ユーザーフィードバックを設定
func (r *Review) SetFeedback(score int, comment string) error {
	// スコアのバリデーション
	if score < 1 || score > 3 {
		return fmt.Errorf("スコアは1-3の整数で指定してください")
	}
	
	// コメントの長さチェック
	if len(comment) > 500 {
		return fmt.Errorf("コメントは500文字以内にしてください")
	}
	
	r.FeedbackScore = &score
	r.FeedbackComment = comment
	r.UpdatedAt = time.Now()
	return nil
}
