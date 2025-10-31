package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/domain/repository"
)

// userRepository はUserRepositoryのPostgreSQL実装
type userRepository struct {
	db *sql.DB
}

// NewUserRepository は新しいUserRepositoryを作成します
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create - ユーザーを作成
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (id, auth0_user_id, email, name, avatar_url, preferences, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Auth0UserID,
		user.Email,
		user.Name,
		user.AvatarURL,
		user.Preferences,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindByID - IDでユーザーを取得
func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, auth0_user_id, email, name, avatar_url, preferences, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Auth0UserID,
		&user.Email,
		&user.Name,
		&user.AvatarURL,
		&user.Preferences,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return user, nil
}

// FindByEmail - メールアドレスでユーザーを取得
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, auth0_user_id, email, name, avatar_url, preferences, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Auth0UserID,
		&user.Email,
		&user.Name,
		&user.AvatarURL,
		&user.Preferences,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return user, nil
}

// FindByAuth0UserID - Auth0ユーザーIDでユーザーを取得
func (r *userRepository) FindByAuth0UserID(ctx context.Context, auth0UserID string) (*model.User, error) {
	query := `
		SELECT id, auth0_user_id, email, name, avatar_url, preferences, created_at, updated_at, deleted_at
		FROM users
		WHERE auth0_user_id = $1 AND deleted_at IS NULL
	`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, auth0UserID).Scan(
		&user.ID,
		&user.Auth0UserID,
		&user.Email,
		&user.Name,
		&user.AvatarURL,
		&user.Preferences,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user by auth0_user_id: %w", err)
	}

	return user, nil
}

// Update - ユーザーを更新
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET email = $2, name = $3, avatar_url = $4, preferences = $5, updated_at = $6
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Name,
		user.AvatarURL,
		user.Preferences,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete - ユーザーを削除（論理削除）
func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
