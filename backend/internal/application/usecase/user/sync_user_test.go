package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/test/testutil"
	"github.com/stretchr/testify/assert"
)

func TestSyncUserByAuth0Sub_ExistingUser(t *testing.T) {
	// Arrange
	mockRepo := testutil.NewMockUserRepository()
	uc := NewUseCase(mockRepo)
	ctx := context.Background()

	auth0Sub := "auth0|123456789"
	existingUser := &model.User{
		ID:          uuid.New().String(),
		Auth0UserID: auth0Sub,
		Email:       "existing@example.com",
		Name:        "Existing User",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now().Add(-24 * time.Hour),
	}

	// Mock設定: ユーザーが既に存在する
	mockRepo.SetUser(existingUser)

	// Act
	user, err := uc.SyncUserByAuth0Sub(ctx, auth0Sub)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, existingUser.ID, user.ID)
	assert.Equal(t, existingUser.Auth0UserID, user.Auth0UserID)
	assert.Equal(t, existingUser.Email, user.Email)
	assert.Equal(t, existingUser.Name, user.Name)
}

func TestSyncUserByAuth0Sub_NewUser(t *testing.T) {
	// Arrange
	mockRepo := testutil.NewMockUserRepository()
	uc := NewUseCase(mockRepo)
	ctx := context.Background()

	auth0Sub := "auth0|987654321"

	// Mock設定: ユーザーが存在しない（何も設定しない）

	// Act
	user, err := uc.SyncUserByAuth0Sub(ctx, auth0Sub)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, auth0Sub, user.Auth0UserID)
	assert.Equal(t, "New User", user.Name)
	assert.Equal(t, "user-987654321@temp.local", user.Email)
	assert.Equal(t, "{}", user.Preferences)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestSyncUserByAuth0Sub_CreateError(t *testing.T) {
	// Arrange
	mockRepo := testutil.NewMockUserRepository()
	uc := NewUseCase(mockRepo)
	ctx := context.Background()

	auth0Sub := "auth0|error_case"
	createError := errors.New("database error")

	// Mock設定: エラーを設定
	mockRepo.SetError(createError)

	// Act
	user, err := uc.SyncUserByAuth0Sub(ctx, auth0Sub)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to create user")
}

func TestGetUserByID_Success(t *testing.T) {
	// Arrange
	mockRepo := testutil.NewMockUserRepository()
	uc := NewUseCase(mockRepo)
	ctx := context.Background()

	userID := uuid.New().String()
	expectedUser := &model.User{
		ID:          userID,
		Auth0UserID: "auth0|123456",
		Email:       "test@example.com",
		Name:        "Test User",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.SetUser(expectedUser)

	// Act
	user, err := uc.GetUserByID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
}

func TestGetUserByID_NotFound(t *testing.T) {
	// Arrange
	mockRepo := testutil.NewMockUserRepository()
	uc := NewUseCase(mockRepo)
	ctx := context.Background()

	userID := uuid.New().String()
	// ユーザーを設定しない（存在しない状態）

	// Act
	user, err := uc.GetUserByID(ctx, userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUpdateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := testutil.NewMockUserRepository()
	uc := NewUseCase(mockRepo)
	ctx := context.Background()

	oldTime := time.Now().Add(-1 * time.Hour)
	user := &model.User{
		ID:          uuid.New().String(),
		Auth0UserID: "auth0|123456",
		Email:       "updated@example.com",
		Name:        "Updated User",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   oldTime,
	}

	// 既存ユーザーとして設定
	mockRepo.SetUser(user)

	// Act
	err := uc.UpdateUser(ctx, user)

	// Assert
	assert.NoError(t, err)
	assert.True(t, user.UpdatedAt.After(oldTime))
}

func TestUpdateUser_Error(t *testing.T) {
	// Arrange
	mockRepo := testutil.NewMockUserRepository()
	uc := NewUseCase(mockRepo)
	ctx := context.Background()

	user := &model.User{
		ID:          uuid.New().String(),
		Auth0UserID: "auth0|123456",
		Email:       "test@example.com",
	}
	updateError := errors.New("update failed")

	mockRepo.SetError(updateError)

	// Act
	err := uc.UpdateUser(ctx, user)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, updateError, err)
}

func TestExtractEmailFromAuth0Sub(t *testing.T) {
	tests := []struct {
		name     string
		auth0Sub string
		expected string
	}{
		{
			name:     "standard auth0 format",
			auth0Sub: "auth0|123456789",
			expected: "user-123456789@temp.local",
		},
		{
			name:     "google oauth format",
			auth0Sub: "google-oauth2|987654321",
			expected: "user-987654321@temp.local",
		},
		{
			name:     "invalid format without pipe",
			auth0Sub: "invalid_format",
			expected: "unknown@temp.local",
		},
		{
			name:     "empty string",
			auth0Sub: "",
			expected: "unknown@temp.local",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractEmailFromAuth0Sub(tt.auth0Sub)
			assert.Equal(t, tt.expected, result)
		})
	}
}
