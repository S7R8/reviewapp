package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pgvector/pgvector-go"
	"github.com/s7r8/reviewapp/internal/domain/model"
)

// KnowledgeRepositoryImpl - PostgreSQL実装
type KnowledgeRepository struct {
	db *sql.DB
}

// NewKnowledgeRepository - コンストラクタ
func NewKnowledgeRepository(db *sql.DB) *KnowledgeRepository {
	return &KnowledgeRepository{db: db}
}

// Create - ナレッジを作成
func (r *KnowledgeRepository) Create(ctx context.Context, knowledge *model.Knowledge) error {
	query := `
		INSERT INTO knowledge (
			id, user_id, title, content, category, priority,
			source_type, source_id, usage_count, last_used_at,
			is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	var embeddingVector interface{}
	if knowledge.HasEmbedding() {
		embeddingVector = pgvector.NewVector(knowledge.GetEmbedding())
	} else {
		embeddingVector = nil
	}
	_, err := r.db.ExecContext(
		ctx,
		query,
		knowledge.ID,
		knowledge.UserID,
		knowledge.Title,
		knowledge.Content,
		knowledge.Category,
		knowledge.Priority,
		knowledge.SourceType,
		knowledge.SourceID,
		knowledge.UsageCount,
		knowledge.LastUsedAt,
		embeddingVector,
		knowledge.IsActive,
		knowledge.CreatedAt,
		knowledge.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create knowledge: %w", err)
	}

	return nil
}

// SearchByKeyword - キーワードで検索（フルテキスト検索）
func (r *KnowledgeRepository) SearchByKeyword(ctx context.Context, userID, keyword string, limit int) ([]*model.Knowledge, error) {
	query := `
		SELECT 
			id, user_id, title, content, category, priority,
			source_type, source_id, usage_count, last_used_at,
			is_active, created_at, updated_at
		FROM knowledge
		WHERE user_id = $1 
			AND is_active = true 
			AND deleted_at IS NULL
			AND (title ILIKE $2 OR content ILIKE $2)
		ORDER BY priority DESC, created_at DESC
		LIMIT $3
	`

	keywordPattern := "%" + keyword + "%"
	rows, err := r.db.QueryContext(ctx, query, userID, keywordPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search knowledge by keyword: %w", err)
	}
	defer rows.Close()

	var knowledges []*model.Knowledge
	for rows.Next() {
		k := &model.Knowledge{}
		var embeddingVector pgvector.Vector
		err := rows.Scan(
			&k.ID,
			&k.UserID,
			&k.Title,
			&k.Content,
			&k.Category,
			&k.Priority,
			&k.SourceType,
			&k.SourceID,
			&k.UsageCount,
			&k.LastUsedAt,
			&embeddingVector,
			&k.IsActive,
			&k.CreatedAt,
			&k.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan knowledge: %w", err)
		}
		if len(embeddingVector.Slice()) > 0 {
			k.Embedding = embeddingVector.Slice()
		}
		knowledges = append(knowledges, k)
	}

	return knowledges, nil
}

// FindByID - IDでナレッジを取得
func (r *KnowledgeRepository) FindByID(ctx context.Context, id string) (*model.Knowledge, error) {
	query := `
		SELECT 
			id, user_id, title, content, category, priority,
			source_type, source_id, usage_count, last_used_at,
			is_active, created_at, updated_at
		FROM knowledge
		WHERE id = $1 AND deleted_at IS NULL
	`

	k := &model.Knowledge{}
	var embeddingVector pgvector.Vector
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&k.ID,
		&k.UserID,
		&k.Title,
		&k.Content,
		&k.Category,
		&k.Priority,
		&k.SourceType,
		&k.SourceID,
		&k.UsageCount,
		&k.LastUsedAt,
		&embeddingVector,
		&k.IsActive,
		&k.CreatedAt,
		&k.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("knowledge not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find knowledge: %w", err)
	}

	if len(embeddingVector.Slice()) > 0 {
		k.Embedding = embeddingVector.Slice()
	}

	return k, nil
}

// FindByUserID - ユーザーIDで全ナレッジを取得
func (r *KnowledgeRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Knowledge, error) {
	query := `
		SELECT 
			id, user_id, title, content, category, priority,
			source_type, source_id, usage_count, last_used_at,
			is_active, created_at, updated_at
		FROM knowledge
		WHERE user_id = $1 AND is_active = true AND deleted_at IS NULL
		ORDER BY priority DESC, created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find knowledge: %w", err)
	}
	defer rows.Close()

	var knowledges []*model.Knowledge
	for rows.Next() {
		k := &model.Knowledge{}
		var embeddingVector pgvector.Vector
		err := rows.Scan(
			&k.ID,
			&k.UserID,
			&k.Title,
			&k.Content,
			&k.Category,
			&k.Priority,
			&k.SourceType,
			&k.SourceID,
			&k.UsageCount,
			&k.LastUsedAt,
			&embeddingVector,
			&k.IsActive,
			&k.CreatedAt,
			&k.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan knowledge: %w", err)
		}
		if len(embeddingVector.Slice()) > 0 {
			k.Embedding = embeddingVector.Slice()
		}

		knowledges = append(knowledges, k)
	}

	return knowledges, nil
}

// FindByCategory - カテゴリでナレッジを取得
func (r *KnowledgeRepository) FindByCategory(ctx context.Context, userID, category string) ([]*model.Knowledge, error) {
	query := `
		SELECT 
			id, user_id, title, content, category, priority,
			source_type, source_id, usage_count, last_used_at,
			is_active, created_at, updated_at
		FROM knowledge
		WHERE user_id = $1 AND category = $2 AND is_active = true AND deleted_at IS NULL
		ORDER BY priority DESC, created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID, category)
	if err != nil {
		return nil, fmt.Errorf("failed to find knowledge by category: %w", err)
	}
	defer rows.Close()

	var knowledges []*model.Knowledge
	for rows.Next() {
		k := &model.Knowledge{}
		var embeddingVector pgvector.Vector
		err := rows.Scan(
			&k.ID,
			&k.UserID,
			&k.Title,
			&k.Content,
			&k.Category,
			&k.Priority,
			&k.SourceType,
			&k.SourceID,
			&k.UsageCount,
			&k.LastUsedAt,
			&embeddingVector,
			&k.IsActive,
			&k.CreatedAt,
			&k.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan knowledge: %w", err)
		}
		if len(embeddingVector.Slice()) > 0 {
			k.Embedding = embeddingVector.Slice()
		}
		knowledges = append(knowledges, k)
	}

	return knowledges, nil
}

// Update - ナレッジを更新
func (r *KnowledgeRepository) Update(ctx context.Context, knowledge *model.Knowledge) error {
	query := `
		UPDATE knowledge
		SET 
			title = $1,
			content = $2,
			category = $3,
			priority = $4,
			usage_count = $5,
			last_used_at = $6,
			is_active = $7,
			updated_at = $8
		WHERE id = $9 AND deleted_at IS NULL
	`

	// EmbeddingをPostgreSQL vector型に変換
	var embeddingVector interface{}
	if knowledge.HasEmbedding() {
		embeddingVector = pgvector.NewVector(knowledge.GetEmbedding())
	} else {
		embeddingVector = nil
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		knowledge.Title,
		knowledge.Content,
		knowledge.Category,
		knowledge.Priority,
		knowledge.UsageCount,
		knowledge.LastUsedAt,
		embeddingVector,
		knowledge.IsActive,
		knowledge.UpdatedAt,
		knowledge.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update knowledge: %w", err)
	}

	return nil
}

// Delete - ナレッジを削除（論理削除）
func (r *KnowledgeRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE knowledge
		SET deleted_at = NOW(), is_active = false, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete knowledge: %w", err)
	}

	// 削除対象が存在しない場合のチェック
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("knowledge not found or already deleted: %s", id)
	}

	return nil
}

// CountByUserID - ユーザーIDでナレッジ総数を取得（有効なもののみ）
func (r *KnowledgeRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM knowledge
		WHERE user_id = $1 AND is_active = true AND deleted_at IS NULL
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count knowledge: %w", err)
	}

	return count, nil
}

// CountByCategory - カテゴリ別のナレッジ数を取得
func (r *KnowledgeRepository) CountByCategory(ctx context.Context, userID string) (map[string]int, error) {
	query := `
		SELECT category, COUNT(*) as count
		FROM knowledge
		WHERE user_id = $1 AND is_active = true AND deleted_at IS NULL
		GROUP BY category
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count by category: %w", err)
	}
	defer rows.Close()

	categoryCounts := make(map[string]int)
	for rows.Next() {
		var category string
		var count int
		err := rows.Scan(&category, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category count: %w", err)
		}
		categoryCounts[category] = count
	}

	return categoryCounts, nil
}

// SearchBySimilarity - ベクトル類似度検索
func (r *KnowledgeRepository) SearchBySimilarity(
	ctx context.Context,
	userID string,
	embedding []float32,
	limit int,
	threshold float64,
) ([]*model.Knowledge, error) {
	// PostgreSQL + pgvectorでコサイン類似度検索
	// <=> はコサイン距離（0に近いほど類似）
	// 1 - (embedding <=> $2) でコサイン類似度に変換（1に近いほど類似）
	query := `
		SELECT 
			id, user_id, title, content, category, priority,
			source_type, source_id, usage_count, last_used_at,
			embedding, is_active, created_at, updated_at,
			1 - (embedding <=> $2) AS similarity
		FROM knowledge
		WHERE 
			user_id = $1 
			AND deleted_at IS NULL
			AND is_active = true
			AND embedding IS NOT NULL
			AND (1 - (embedding <=> $2)) >= $3
		ORDER BY embedding <=> $2
		LIMIT $4
	`

	embeddingVector := pgvector.NewVector(embedding)

	rows, err := r.db.QueryContext(ctx, query, userID, embeddingVector, threshold, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search by similarity: %w", err)
	}
	defer rows.Close()

	var knowledges []*model.Knowledge
	for rows.Next() {
		k := &model.Knowledge{}
		var embeddingVec pgvector.Vector
		var similarity float64

		err := rows.Scan(
			&k.ID,
			&k.UserID,
			&k.Title,
			&k.Content,
			&k.Category,
			&k.Priority,
			&k.SourceType,
			&k.SourceID,
			&k.UsageCount,
			&k.LastUsedAt,
			&embeddingVec,
			&k.IsActive,
			&k.CreatedAt,
			&k.UpdatedAt,
			&similarity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan knowledge: %w", err)
		}

		if len(embeddingVec.Slice()) > 0 {
			k.Embedding = embeddingVec.Slice()
		}

		knowledges = append(knowledges, k)
	}

	return knowledges, nil
}

// FindWithoutEmbedding - Embeddingが未設定のナレッジを取得
func (r *KnowledgeRepository) FindWithoutEmbedding(ctx context.Context, limit int) ([]*model.Knowledge, error) {
	query := `
		SELECT 
			id, user_id, title, content, category, priority,
			source_type, source_id, usage_count, last_used_at,
			embedding, is_active, created_at, updated_at
		FROM knowledge
		WHERE 
			deleted_at IS NULL
			AND embedding IS NULL
		ORDER BY created_at ASC
	`

	// limitが指定されている場合はLIMIT句を追加
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find knowledge without embedding: %w", err)
	}
	defer rows.Close()

	var knowledges []*model.Knowledge
	for rows.Next() {
		k := &model.Knowledge{}
		var embeddingVector pgvector.Vector
		err := rows.Scan(
			&k.ID,
			&k.UserID,
			&k.Title,
			&k.Content,
			&k.Category,
			&k.Priority,
			&k.SourceType,
			&k.SourceID,
			&k.UsageCount,
			&k.LastUsedAt,
			&embeddingVector,
			&k.IsActive,
			&k.CreatedAt,
			&k.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan knowledge: %w", err)
		}

		if len(embeddingVector.Slice()) > 0 {
			k.Embedding = embeddingVector.Slice()
		}

		knowledges = append(knowledges, k)
	}

	return knowledges, nil
}
