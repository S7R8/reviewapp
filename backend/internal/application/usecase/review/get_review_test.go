package review

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReviewRepositoryForGet - モックリポジトリ
type MockReviewRepositoryForGet struct {
	mock.Mock
}

func (m *MockReviewRepositoryForGet) FindByID(ctx context.Context, reviewID string) (*model.Review, error) {
	args := m.Called(ctx, reviewID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Review), args.Error(1)
}

func (m *MockReviewRepositoryForGet) Create(ctx context.Context, review *model.Review) error {
	return nil
}

func (m *MockReviewRepositoryForGet) Update(ctx context.Context, review *model.Review) error {
	return nil
}

func (m *MockReviewRepositoryForGet) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *MockReviewRepositoryForGet) FindByUserID(ctx context.Context, userID string, limit int) ([]*model.Review, error) {
	return nil, nil
}

func (m *MockReviewRepositoryForGet) FindRecentByUserID(ctx context.Context, userID string, limit int) ([]*model.Review, error) {
	return nil, nil
}

func (m *MockReviewRepositoryForGet) UpdateFeedback(ctx context.Context, reviewID string, score int, comment string) error {
	return nil
}

func (m *MockReviewRepositoryForGet) ListWithFilters(ctx context.Context, userID string, filters map[string]interface{}, sortBy, sortOrder string, limit, offset int) ([]*model.Review, error) {
	return nil, nil
}

func (m *MockReviewRepositoryForGet) CountWithFilters(ctx context.Context, userID string, filters map[string]interface{}) (int, error) {
	return 0, nil
}

func (m *MockReviewRepositoryForGet) CountByUserID(ctx context.Context, userID string) (int, error) {
	return 0, nil
}

func (m *MockReviewRepositoryForGet) CountByUserIDAndDateRange(ctx context.Context, userID string, from, to time.Time) (int, error) {
	return 0, nil
}

func (m *MockReviewRepositoryForGet) GetAverageFeedbackScore(ctx context.Context, userID string) (float64, error) {
	return 0.0, nil
}

// TestGetReviewUseCase_Execute - 正常系テスト
func TestGetReviewUseCase_Execute(t *testing.T) {
	// Arrange
	mockRepo := new(MockReviewRepositoryForGet)
	usecase := NewGetReviewUseCase(mockRepo)

	reviewID := "review-123"
	userID := "user-456"

	expectedReview := &model.Review{
		ID:           reviewID,
		UserID:       userID,
		Code:         "test code",
		Language:     "go",
		ReviewResult: "test review",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.On("FindByID", mock.Anything, reviewID).Return(expectedReview, nil)

	// Act
	input := GetReviewInput{
		ReviewID: reviewID,
		UserID:   userID,
	}
	output, err := usecase.Execute(context.Background(), input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, expectedReview.ID, output.Review.ID)
	assert.Equal(t, expectedReview.UserID, output.Review.UserID)
	assert.Equal(t, expectedReview.Code, output.Review.Code)

	mockRepo.AssertExpectations(t)
}

// TestGetReviewUseCase_Execute_NotFound - レビューが見つからない場合
func TestGetReviewUseCase_Execute_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockReviewRepositoryForGet)
	usecase := NewGetReviewUseCase(mockRepo)

	reviewID := "non-existent-review"
	userID := "user-456"

	mockRepo.On("FindByID", mock.Anything, reviewID).Return(nil, errors.New("not found"))

	// Act
	input := GetReviewInput{
		ReviewID: reviewID,
		UserID:   userID,
	}
	output, err := usecase.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "failed to find review")

	mockRepo.AssertExpectations(t)
}

// TestGetReviewUseCase_Execute_Forbidden - 他人のレビューへのアクセス
func TestGetReviewUseCase_Execute_Forbidden(t *testing.T) {
	// Arrange
	mockRepo := new(MockReviewRepositoryForGet)
	usecase := NewGetReviewUseCase(mockRepo)

	reviewID := "review-123"
	ownerUserID := "user-owner"
	requestUserID := "user-other"

	otherUsersReview := &model.Review{
		ID:           reviewID,
		UserID:       ownerUserID, // 別のユーザーのレビュー
		Code:         "test code",
		Language:     "go",
		ReviewResult: "test review",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo.On("FindByID", mock.Anything, reviewID).Return(otherUsersReview, nil)

	// Act
	input := GetReviewInput{
		ReviewID: reviewID,
		UserID:   requestUserID, // 別のユーザーがアクセス
	}
	output, err := usecase.Execute(context.Background(), input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "このレビューにアクセスする権限がありません")

	mockRepo.AssertExpectations(t)
}

// TestGetReviewUseCase_Execute_EmptyReviewID - ReviewIDが空の場合
func TestGetReviewUseCase_Execute_EmptyReviewID(t *testing.T) {
	mockRepo := new(MockReviewRepositoryForGet)
	usecase := NewGetReviewUseCase(mockRepo)

	input := GetReviewInput{
		ReviewID: "",
		UserID:   "user-456",
	}

	output, err := usecase.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "review ID is required")
}

// TestGetReviewUseCase_Execute_EmptyUserID - UserIDが空の場合
func TestGetReviewUseCase_Execute_EmptyUserID(t *testing.T) {
	mockRepo := new(MockReviewRepositoryForGet)
	usecase := NewGetReviewUseCase(mockRepo)

	input := GetReviewInput{
		ReviewID: "review-123",
		UserID:   "",
	}

	output, err := usecase.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "user ID is required")
}

// TestGetReviewUseCase_Execute_WithStructuredResult - 構造化結果を含む場合
func TestGetReviewUseCase_Execute_WithStructuredResult(t *testing.T) {
	mockRepo := new(MockReviewRepositoryForGet)
	usecase := NewGetReviewUseCase(mockRepo)

	reviewID := "review-123"
	userID := "user-456"

	review := &model.Review{
		ID:           reviewID,
		UserID:       userID,
		Code:         "test code",
		Language:     "go",
		ReviewResult: "test review",
		StructuredResult: &model.StructuredReviewResult{
			Summary:    "Good code overall",
			GoodPoints: []string{"Clean", "Well documented"},
			Improvements: []model.Improvement{
				{
					Title:       "Optimization needed",
					Description: "Can be optimized",
					Severity:    "medium",
				},
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", mock.Anything, reviewID).Return(review, nil)

	input := GetReviewInput{
		ReviewID: reviewID,
		UserID:   userID,
	}

	output, err := usecase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.NotNil(t, output.Review.StructuredResult)
	assert.Equal(t, "Good code overall", output.Review.StructuredResult.Summary)
	assert.Len(t, output.Review.StructuredResult.Improvements, 1)

	mockRepo.AssertExpectations(t)
}

// TestGetReviewUseCase_Execute_WithReferencedKnowledge - 参照ナレッジを含む場合
func TestGetReviewUseCase_Execute_WithReferencedKnowledge(t *testing.T) {
	mockRepo := new(MockReviewRepositoryForGet)
	usecase := NewGetReviewUseCase(mockRepo)

	reviewID := "review-123"
	userID := "user-456"

	review := &model.Review{
		ID:                  reviewID,
		UserID:              userID,
		Code:                "test code",
		Language:            "go",
		ReviewResult:        "test review",
		ReferencedKnowledge: []string{"knowledge-1", "knowledge-2"},
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	mockRepo.On("FindByID", mock.Anything, reviewID).Return(review, nil)

	input := GetReviewInput{
		ReviewID: reviewID,
		UserID:   userID,
	}

	output, err := usecase.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output.Review.ReferencedKnowledge, 2)
	assert.Equal(t, "knowledge-1", output.Review.ReferencedKnowledge[0])

	mockRepo.AssertExpectations(t)
}
