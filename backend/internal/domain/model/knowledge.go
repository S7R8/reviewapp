package model

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Knowledge - ナレッジエンティティ
type Knowledge struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Category   string     `json:"category"`
	Priority   int        `json:"priority"`
	SourceType string     `json:"source_type"`
	SourceID   *string    `json:"source_id"`
	UsageCount int        `json:"usage_count"`
	LastUsedAt *time.Time `json:"last_used_at"`
	Embedding  []float32  `json:"-"`
	IsActive   bool       `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// カテゴリの定数
const (
	CategoryErrorHandling = "error_handling"
	CategoryTesting       = "testing"
	CategoryPerformance   = "performance"
	CategorySecurity      = "security"
	CategoryCleanCode     = "clean_code"
	CategoryArchitecture  = "architecture"
	CategoryOther         = "other"
)

// ソースタイプの定数
const (
	SourceTypeManual       = "manual"
	SourceTypeReview       = "review"
	SourceTypeConversation = "conversation"
)

// バリデーションエラー
var (
	ErrTitleRequired      = errors.New("タイトルは必須です")
	ErrTitleTooLong       = errors.New("タイトルは200文字以内にしてください")
	ErrContentRequired    = errors.New("内容は必須です")
	ErrCategoryRequired   = errors.New("カテゴリは必須です")
	ErrCategoryInvalid    = errors.New("無効なカテゴリです")
	ErrPriorityRequired   = errors.New("重要度は必須です")
	ErrPriorityOutOfRange = errors.New("重要度は1-5の整数で指定してください")
	ErrSourceTypeRequired = errors.New("ソースタイプは必須です")
	ErrSourceTypeInvalid  = errors.New("無効なソースタイプです")
)

// 許可されたカテゴリ
var validCategories = map[string]bool{
	CategoryErrorHandling: true,
	CategoryTesting:       true,
	CategoryPerformance:   true,
	CategorySecurity:      true,
	CategoryCleanCode:     true,
	CategoryArchitecture:  true,
	CategoryOther:         true,
}

// 許可されたソースタイプ
var validSourceTypes = map[string]bool{
	SourceTypeManual:       true,
	SourceTypeReview:       true,
	SourceTypeConversation: true,
}

// NewKnowledge - ナレッジを生成
func NewKnowledge(userID, title, content, category string, priority int) (*Knowledge, error) {
	now := time.Now()
	k := &Knowledge{
		ID:         uuid.New().String(),
		UserID:     userID,
		Title:      strings.TrimSpace(title),
		Content:    strings.TrimSpace(content),
		Category:   category,
		Priority:   priority,
		SourceType: SourceTypeManual,
		SourceID:   nil,
		UsageCount: 0,
		LastUsedAt: nil,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := k.Validate(); err != nil {
		return nil, err
	}

	return k, nil
}

// NewKnowledgeFromReview - レビューから生成されたナレッジ
func NewKnowledgeFromReview(userID, reviewID, title, content, category string, priority int) (*Knowledge, error) {
	now := time.Now()
	k := &Knowledge{
		ID:         uuid.New().String(),
		UserID:     userID,
		Title:      strings.TrimSpace(title),
		Content:    strings.TrimSpace(content),
		Category:   category,
		Priority:   priority,
		SourceType: SourceTypeReview,
		SourceID:   &reviewID,
		UsageCount: 0,
		LastUsedAt: nil,
		IsActive:   true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := k.Validate(); err != nil {
		return nil, err
	}

	return k, nil
}

// Validate - ナレッジのバリデーション
func (k *Knowledge) Validate() error {
	// タイトルチェック
	if k.Title == "" {
		return ErrTitleRequired
	}
	if len(k.Title) > 200 {
		return ErrTitleTooLong
	}

	// 内容チェック
	if k.Content == "" {
		return ErrContentRequired
	}

	// カテゴリチェック
	if k.Category == "" {
		return ErrCategoryRequired
	}
	if !validCategories[k.Category] {
		return ErrCategoryInvalid
	}

	// 優先度チェック
	if k.Priority < 1 || k.Priority > 5 {
		return ErrPriorityOutOfRange
	}

	// ソースタイプチェック
	if k.SourceType == "" {
		return ErrSourceTypeRequired
	}
	if !validSourceTypes[k.SourceType] {
		return ErrSourceTypeInvalid
	}

	return nil
}

// IncrementUsage - 使用カウントを増やす
func (k *Knowledge) IncrementUsage() {
	k.UsageCount++
	now := time.Now()
	k.LastUsedAt = &now
	k.UpdatedAt = now
}

// Deactivate - 無効化
func (k *Knowledge) Deactivate() {
	k.IsActive = false
	k.UpdatedAt = time.Now()
}

// Activate - 有効化
func (k *Knowledge) Activate() {
	k.IsActive = true
	k.UpdatedAt = time.Now()
}

// UpdateContent - 内容を更新
func (k *Knowledge) UpdateContent(title, content, category string, priority int) error {
	k.Title = strings.TrimSpace(title)
	k.Content = strings.TrimSpace(content)
	k.Category = category
	k.Priority = priority
	k.UpdatedAt = time.Now()

	return k.Validate()
}

// SetEmbedding - Embeddingベクトルを設定
func (k *Knowledge) SetEmbedding(embedding []float32) {
	k.Embedding = embedding
	k.UpdatedAt = time.Now()
}

// HasEmbedding - Embeddingが設定されているか確認
func (k *Knowledge) HasEmbedding() bool {
	return k.Embedding != nil && len(k.Embedding) > 0
}

// GetEmbedding - Embeddingを取得
func (k *Knowledge) GetEmbedding() []float32 {
	if !k.HasEmbedding() {
		return nil
	}
	embedding := make([]float32, len(k.Embedding))
	copy(embedding, k.Embedding)
	return embedding
}

// ClearEmbedding - Embeddingをクリア
func (k *Knowledge) ClearEmbedding() {
	k.Embedding = nil
	k.UpdatedAt = time.Now()
}
