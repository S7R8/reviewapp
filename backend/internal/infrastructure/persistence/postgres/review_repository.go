package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/s7r8/reviewapp/internal/domain/model"
)

// ReviewRepository - PostgreSQL実装
type ReviewRepository struct {
	db *sql.DB
}

// NewReviewRepository - コンストラクタ
func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// Create - レビューを作成
func (r *ReviewRepository) Create(ctx context.Context, review *model.Review) error {
	// トランザクション開始
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. reviewsテーブルにINSERT
	query := `
		INSERT INTO reviews (
			id, user_id, code, language, context,
			review_result, llm_provider, llm_model, tokens_used,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		review.ID,
		review.UserID,
		review.Code,
		review.Language,
		review.Context,
		review.ReviewResult,
		review.LLMProvider,
		review.LLMModel,
		review.TokensUsed,
		review.CreatedAt,
		review.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	// 2. review_knowledgeテーブルにINSERT（参照されたナレッジ）
	if len(review.ReferencedKnowledge) > 0 {
		stmt, err := tx.PrepareContext(ctx, `
			INSERT INTO review_knowledge (review_id, knowledge_id, created_at)
			VALUES ($1, $2, NOW())
		`)
		if err != nil {
			return fmt.Errorf("failed to prepare review_knowledge insert: %w", err)
		}
		defer stmt.Close()

		for _, knowledgeID := range review.ReferencedKnowledge {
			_, err = stmt.ExecContext(ctx, review.ID, knowledgeID)
			if err != nil {
				return fmt.Errorf("failed to insert review_knowledge: %w", err)
			}
		}
	}

	// トランザクションコミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// FindByID - IDでレビューを取得（中間テーブルからナレッジIDを取得）
func (r *ReviewRepository) FindByID(ctx context.Context, id string) (*model.Review, error) {
	// 1. reviewsテーブルから基本情報を取得
	query := `
		SELECT 
			id, user_id, code, language, context,
			review_result, llm_provider, llm_model, tokens_used,
			feedback_score, feedback_comment, created_at, updated_at, deleted_at
		FROM reviews
		WHERE id = $1 AND deleted_at IS NULL
	`

	review := &model.Review{}
	var context, llmProvider, llmModel, feedbackComment sql.NullString
	var feedbackScore sql.NullInt32
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&review.ID,
		&review.UserID,
		&review.Code,
		&review.Language,
		&context,
		&review.ReviewResult,
		&llmProvider,
		&llmModel,
		&review.TokensUsed,
		&feedbackScore,
		&feedbackComment,
		&review.CreatedAt,
		&review.UpdatedAt,
		&deletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("review not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find review: %w", err)
	}

	// Nullable フィールドの処理
	if context.Valid {
		review.Context = context.String
	}
	if llmProvider.Valid {
		review.LLMProvider = llmProvider.String
	}
	if llmModel.Valid {
		review.LLMModel = llmModel.String
	}
	if feedbackScore.Valid {
		score := int(feedbackScore.Int32)
		review.FeedbackScore = &score
	}
	if feedbackComment.Valid {
		review.FeedbackComment = feedbackComment.String
	}
	if deletedAt.Valid {
		review.DeletedAt = &deletedAt.Time
	}

	// 2. review_knowledgeテーブルから参照されたナレッジIDを取得
	knowledgeQuery := `
		SELECT knowledge_id
		FROM review_knowledge
		WHERE review_id = $1
		ORDER BY created_at
	`

	rows, err := r.db.QueryContext(ctx, knowledgeQuery, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query review_knowledge: %w", err)
	}
	defer rows.Close()

	var referencedKnowledge []string
	for rows.Next() {
		var knowledgeID string
		if err := rows.Scan(&knowledgeID); err != nil {
			return nil, fmt.Errorf("failed to scan knowledge_id: %w", err)
		}
		referencedKnowledge = append(referencedKnowledge, knowledgeID)
	}

	review.ReferencedKnowledge = referencedKnowledge
	return review, nil
}

// FindByUserID - ユーザーIDで全レビューを取得
func (r *ReviewRepository) FindByUserID(ctx context.Context, userID string, limit int) ([]*model.Review, error) {
	query := `
		SELECT 
			id, user_id, code, language, context,
			review_result, llm_provider, llm_model, tokens_used,
			feedback_score, feedback_comment, created_at, updated_at, deleted_at
		FROM reviews
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find reviews: %w", err)
	}
	defer rows.Close()

	var reviews []*model.Review
	for rows.Next() {
		review := &model.Review{}
		var fileName, context, llmProvider, llmModel, feedbackComment sql.NullString
		var feedbackScore sql.NullInt32

		err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.Code,
			&review.Language,
			&fileName,
			&context,
			&review.ReviewResult,
			&llmProvider,
			&llmModel,
			&review.TokensUsed,
			&feedbackScore,
			&feedbackComment,
			&review.CreatedAt,
			&review.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}

		// Nullable フィールドの処理
		if context.Valid {
			review.Context = context.String
		}
		if llmProvider.Valid {
			review.LLMProvider = llmProvider.String
		}
		if llmModel.Valid {
			review.LLMModel = llmModel.String
		}
		if feedbackScore.Valid {
			score := int(feedbackScore.Int32)
			review.FeedbackScore = &score
		}
		if feedbackComment.Valid {
			review.FeedbackComment = feedbackComment.String
		}

		// 各レビューの参照ナレッジを取得
		knowledgeQuery := `
			SELECT knowledge_id
			FROM review_knowledge
			WHERE review_id = $1
			ORDER BY created_at
		`
		knowledgeRows, err := r.db.QueryContext(ctx, knowledgeQuery, review.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to query review_knowledge: %w", err)
		}

		var referencedKnowledge []string
		for knowledgeRows.Next() {
			var knowledgeID string
			if err := knowledgeRows.Scan(&knowledgeID); err != nil {
				knowledgeRows.Close()
				return nil, fmt.Errorf("failed to scan knowledge_id: %w", err)
			}
			referencedKnowledge = append(referencedKnowledge, knowledgeID)
		}
		knowledgeRows.Close()

		review.ReferencedKnowledge = referencedKnowledge
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// Update - レビューを更新
func (r *ReviewRepository) Update(ctx context.Context, review *model.Review) error {
	query := `
		UPDATE reviews
		SET 
			review_result = $1,
			llm_provider = $2,
			llm_model = $3,
			tokens_used = $4,
			feedback_score = $5,
			feedback_comment = $6,
			updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		review.ReviewResult,
		review.LLMProvider,
		review.LLMModel,
		review.TokensUsed,
		review.FeedbackScore,
		review.FeedbackComment,
		review.UpdatedAt,
		review.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}

	return nil
}

// Delete - レビューを削除（物理削除、中間テーブルはCASCADE削除）
func (r *ReviewRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM reviews WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	return nil
}

// FindRecentByUserID - 最近のレビューを取得
func (r *ReviewRepository) FindRecentByUserID(ctx context.Context, userID string, limit int) ([]*model.Review, error) {
	return r.FindByUserID(ctx, userID, limit)
}
