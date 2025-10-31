package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/domain/repository"
)

// UseCase はユーザー関連のユースケース
type UseCase struct {
	userRepo repository.UserRepository
}

// NewUseCase は新しいUseCaseを作成します
func NewUseCase(userRepo repository.UserRepository) *UseCase {
	return &UseCase{
		userRepo: userRepo,
	}
}

// SyncUserByAuth0Sub はAuth0のSubjectからユーザーを同期します
// ユーザーが存在しない場合は自動作成します
func (uc *UseCase) SyncUserByAuth0Sub(ctx context.Context, auth0Sub string) (*model.User, error) {
	// 1. auth0_user_idでユーザー検索
	user, err := uc.userRepo.FindByAuth0UserID(ctx, auth0Sub)
	if err == nil {
		// 既存ユーザーが見つかった
		return user, nil
	}

	// 2. ユーザーが存在しない場合、新規作成
	// （初回ログイン時）
	newUser := &model.User{
		ID:          uuid.New().String(),
		Auth0UserID: auth0Sub,
		Email:       extractEmailFromAuth0Sub(auth0Sub), // 仮のメール
		Name:        "New User",                         // デフォルト名
		AvatarURL:   nil,
		Preferences: "{}",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return newUser, nil
}

// GetUserByID はIDでユーザーを取得します
func (uc *UseCase) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return uc.userRepo.FindByID(ctx, id)
}

// UpdateUser はユーザー情報を更新します
func (uc *UseCase) UpdateUser(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()
	return uc.userRepo.Update(ctx, user)
}

// extractEmailFromAuth0Sub はAuth0 Subから仮のメールアドレスを生成します
// 実際のメールアドレスはJWTのclaimsから取得すべきですが、
// ここでは簡易的に実装しています
func extractEmailFromAuth0Sub(auth0Sub string) string {
	// auth0|123456 → user-123456@temp.local
	parts := strings.Split(auth0Sub, "|")
	if len(parts) == 2 {
		return fmt.Sprintf("user-%s@temp.local", parts[1])
	}
	return "unknown@temp.local"
}
