package postgres

import (
	"context"
	"database/sql"
	"encoding/json" // ★ 追加
	"fmt"
	"strings"
	"time"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/infrastructure/parser"
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

	// ★ StructuredResult を JSONB に変換
	var reviewResultJSON []byte
	if review.StructuredResult != nil {
		reviewResultJSON, err = json.Marshal(review.StructuredResult)
		if err != nil {
			return fmt.Errorf("failed to marshal review result: %w", err)
		}
	} else {
		// フォールバック: 空の構造化データ
		reviewResultJSON = []byte(`{"summary":"","good_points":[],"improvements":[]}`)
	}

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
		reviewResultJSON,
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

// FindByID - IDでレビューを取得
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
	var reviewResultJSON []byte
	var feedbackScore sql.NullInt32
	var deletedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&review.ID,
		&review.UserID,
		&review.Code,
		&review.Language,
		&context,
		&reviewResultJSON,
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

	// ★ JSONBから構造化データを復元
	if len(reviewResultJSON) > 0 {
		var structured model.StructuredReviewResult
		if err := json.Unmarshal(reviewResultJSON, &structured); err != nil {
			return nil, fmt.Errorf("failed to unmarshal review result: %w", err)
		}
		review.StructuredResult = &structured
		// ReviewResult には空文字列を設定（互換性のため）
		review.ReviewResult = ""
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
			feedback_score, feedback_comment, created_at, updated_at
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
		var context, llmProvider, llmModel, feedbackComment sql.NullString
		var reviewResultJSON []byte
		var feedbackScore sql.NullInt32

		err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.Code,
			&review.Language,
			&context,
			&reviewResultJSON,
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

		// ★ JSONBから構造化データを復元
		if len(reviewResultJSON) > 0 {
			var structured model.StructuredReviewResult
			if err := json.Unmarshal(reviewResultJSON, &structured); err == nil {
				review.StructuredResult = &structured
			}
			review.ReviewResult = ""
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

// ListWithFilters - フィルター、ソート、ページネーション付きでレビュー一覧を取得
func (r *ReviewRepository) ListWithFilters(ctx context.Context, userID string, filters map[string]interface{}, sortBy, sortOrder string, limit, offset int) ([]*model.Review, error) {
	// WHERE句を動的に構築
	where := "user_id = $1 AND deleted_at IS NULL"
	params := []interface{}{userID}
	paramIndex := 2

	// 言語フィルター
	if language, ok := filters["language"].(string); ok && language != "" {
		where += fmt.Sprintf(" AND language = $%d", paramIndex)
		params = append(params, language)
		paramIndex++
	}

	// ステータスフィルター
	if status, ok := filters["status"].(string); ok && status != "" {
		where += fmt.Sprintf(" AND status = $%d", paramIndex)
		params = append(params, status)
		paramIndex++
	}

	// 期間フィルター（開始日）
	if dateFrom, ok := filters["date_from"].(time.Time); ok && !dateFrom.IsZero() {
		where += fmt.Sprintf(" AND created_at >= $%d", paramIndex)
		params = append(params, dateFrom)
		paramIndex++
	}

	// 期間フィルター（終了日）
	if dateTo, ok := filters["date_to"].(time.Time); ok && !dateTo.IsZero() {
		where += fmt.Sprintf(" AND created_at <= $%d", paramIndex)
		params = append(params, dateTo)
		paramIndex++
	}

	// ORDER BY句
	orderBy := fmt.Sprintf("%s %s", sortBy, strings.ToUpper(sortOrder))

	// クエリ構築
	query := fmt.Sprintf(`
		SELECT 
			id, user_id, code, language, context,
			review_result, llm_provider, llm_model, tokens_used,
			feedback_score, feedback_comment, created_at, updated_at
		FROM reviews
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d
	`, where, orderBy, paramIndex, paramIndex+1)

	params = append(params, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to list reviews: %w", err)
	}
	defer rows.Close()

	var reviews []*model.Review
	for rows.Next() {
		review := &model.Review{}
		var context, llmProvider, llmModel, feedbackComment sql.NullString
		var reviewResultJSON []byte
		var feedbackScore sql.NullInt32

		err := rows.Scan(
			&review.ID,
			&review.UserID,
			&review.Code,
			&review.Language,
			&context,
			&reviewResultJSON,
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

		// ★ JSONBから構造化データを復元
		if len(reviewResultJSON) > 0 {
			var structured model.StructuredReviewResult
			json.Unmarshal(reviewResultJSON, &structured)
			review.StructuredResult = &structured
			review.ReviewResult = ""
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

		// ★ markdownから構造化データをパース
		if review.ReviewResult != "" {
			structured := parser.ParseReviewMarkdown(review.ReviewResult)
			review.StructuredResult = structured
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

// CountWithFilters - フィルター条件に合致するレビューの総数を取得
func (r *ReviewRepository) CountWithFilters(ctx context.Context, userID string, filters map[string]interface{}) (int, error) {
	// WHERE句を動的に構築
	where := "user_id = $1 AND deleted_at IS NULL"
	params := []interface{}{userID}
	paramIndex := 2

	// 言語フィルター
	if language, ok := filters["language"].(string); ok && language != "" {
		where += fmt.Sprintf(" AND language = $%d", paramIndex)
		params = append(params, language)
		paramIndex++
	}

	// ステータスフィルター
	if status, ok := filters["status"].(string); ok && status != "" {
		where += fmt.Sprintf(" AND status = $%d", paramIndex)
		params = append(params, status)
		paramIndex++
	}

	// 期間フィルター（開始日）
	if dateFrom, ok := filters["date_from"].(time.Time); ok && !dateFrom.IsZero() {
		where += fmt.Sprintf(" AND created_at >= $%d", paramIndex)
		params = append(params, dateFrom)
		paramIndex++
	}

	// 期間フィルター（終了日）
	if dateTo, ok := filters["date_to"].(time.Time); ok && !dateTo.IsZero() {
		where += fmt.Sprintf(" AND created_at <= $%d", paramIndex)
		params = append(params, dateTo)
		paramIndex++
	}

	// クエリ構築
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM reviews
		WHERE %s
	`, where)

	var count int
	err := r.db.QueryRowContext(ctx, query, params...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count reviews: %w", err)
	}

	return count, nil
}

// Update - レビューを更新
func (r *ReviewRepository) Update(ctx context.Context, review *model.Review) error {
	var reviewResultJSON []byte
	var err error
	if review.StructuredResult != nil {
		reviewResultJSON, err = json.Marshal(review.StructuredResult)
		if err != nil {
			return fmt.Errorf("failed to marshal review result: %w", err)
		}
	} else {
		reviewResultJSON = []byte(`{"summary":"","good_points":[],"improvements":[]}`)
	}

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

	_, err = r.db.ExecContext(
		ctx,
		query,
		reviewResultJSON, // ★ JSONB
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

// UpdateFeedback - フィードバックを更新
func (r *ReviewRepository) UpdateFeedback(ctx context.Context, reviewID string, score int, comment string) error {
query := `
UPDATE reviews
SET feedback_score = $1,
feedback_comment = $2,
updated_at = NOW()
WHERE id = $3 AND deleted_at IS NULL
`

result, err := r.db.ExecContext(ctx, query, score, comment, reviewID)
if err != nil {
return fmt.Errorf("failed to update feedback: %w", err)
}

rowsAffected, err := result.RowsAffected()
if err != nil {
return fmt.Errorf("failed to get rows affected: %w", err)
}

if rowsAffected == 0 {
return fmt.Errorf("レビューが見つかりません: %s", reviewID)
}

return nil
}

// CountByUserID - ユーザーIDでレビュー総数を取得
func (r *ReviewRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM reviews
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count reviews: %w", err)
	}

	return count, nil
}

// CountByUserIDAndDateRange - 期間内のレビュー数を取得
func (r *ReviewRepository) CountByUserIDAndDateRange(ctx context.Context, userID string, from, to time.Time) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM reviews
		WHERE user_id = $1 
		  AND created_at >= $2 
		  AND created_at <= $3
		  AND deleted_at IS NULL
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID, from, to).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count reviews by date range: %w", err)
	}

	return count, nil
}

// GetAverageFeedbackScore - フィードバックスコアの平均を取得
func (r *ReviewRepository) GetAverageFeedbackScore(ctx context.Context, userID string) (float64, error) {
	query := `
		SELECT COALESCE(AVG(feedback_score), 0)
		FROM reviews
		WHERE user_id = $1 
		  AND feedback_score IS NOT NULL
		  AND deleted_at IS NULL
	`

	var average float64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&average)
	if err != nil {
		return 0, fmt.Errorf("failed to get average feedback score: %w", err)
	}

	return average, nil
}
