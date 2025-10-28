package review

import (
	"context"
	"errors"
	"testing"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/test/testutil"
)

// テスト用のレビューを作成
func createTestReview(userID string) *model.Review {
	review := model.NewReview(userID, "test code", "go", "test context")
	review.SetReviewResult("test result", nil, []string{}, "claude", "claude-3-5-sonnet", 100)
	return review
}

// TC-RV-004-01: 全フィールド正常（score=3, commentあり）
func TestUpdateFeedbackUseCase_Execute_Success_WithComment(t *testing.T) {
	ctx := context.Background()
	testUserID := "user-123"
	testReviewID := "review-123"
	testReview := createTestReview(testUserID)
	testReview.ID = testReviewID

	// モックRepositoryを準備
	mockRepo := testutil.NewMockReviewRepository()
	mockRepo.Create(ctx, testReview) // レビューを追加

	useCase := NewUpdateFeedbackUseCase(mockRepo)

	input := UpdateFeedbackInput{
		ReviewID: testReviewID,
		UserID:   testUserID,
		Score:    3,
		Comment:  "とても役に立ちました",
	}

	output, err := useCase.Execute(ctx, input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if output == nil {
		t.Fatal("Expected output, got nil")
	}

	if output.ReviewID != testReviewID {
		t.Errorf("Expected review ID %s, got %s", testReviewID, output.ReviewID)
	}

	if output.FeedbackScore != 3 {
		t.Errorf("Expected score 3, got %d", output.FeedbackScore)
	}

	if output.FeedbackComment != "とても役に立ちました" {
		t.Errorf("Expected comment 'とても役に立ちました', got %s", output.FeedbackComment)
	}
}

// TC-RV-004-02: スコアのみ（commentなし）
func TestUpdateFeedbackUseCase_Execute_Success_WithoutComment(t *testing.T) {
	ctx := context.Background()
	testUserID := "user-123"
	testReviewID := "review-123"
	testReview := createTestReview(testUserID)
	testReview.ID = testReviewID

	mockRepo := testutil.NewMockReviewRepository()
	mockRepo.Create(ctx, testReview)

	useCase := NewUpdateFeedbackUseCase(mockRepo)

	input := UpdateFeedbackInput{
		ReviewID: testReviewID,
		UserID:   testUserID,
		Score:    1,
		Comment:  "",
	}

	output, err := useCase.Execute(ctx, input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if output.FeedbackScore != 1 {
		t.Errorf("Expected score 1, got %d", output.FeedbackScore)
	}

	if output.FeedbackComment != "" {
		t.Errorf("Expected empty comment, got %s", output.FeedbackComment)
	}
}

// TC-RV-004-06: scoreが空（0）
func TestUpdateFeedbackUseCase_Execute_InvalidScore_Zero(t *testing.T) {
	ctx := context.Background()
	testUserID := "user-123"
	testReviewID := "review-123"

	mockRepo := testutil.NewMockReviewRepository()
	useCase := NewUpdateFeedbackUseCase(mockRepo)

	input := UpdateFeedbackInput{
		ReviewID: testReviewID,
		UserID:   testUserID,
		Score:    0,
		Comment:  "",
	}

	_, err := useCase.Execute(ctx, input)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "スコアは1-3の整数で指定してください"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TC-RV-004-08: scoreが4
func TestUpdateFeedbackUseCase_Execute_InvalidScore_Four(t *testing.T) {
	ctx := context.Background()
	testUserID := "user-123"
	testReviewID := "review-123"

	mockRepo := testutil.NewMockReviewRepository()
	useCase := NewUpdateFeedbackUseCase(mockRepo)

	input := UpdateFeedbackInput{
		ReviewID: testReviewID,
		UserID:   testUserID,
		Score:    4,
		Comment:  "",
	}

	_, err := useCase.Execute(ctx, input)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "スコアは1-3の整数で指定してください" {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TC-RV-004-09: commentが501文字
func TestUpdateFeedbackUseCase_Execute_InvalidComment_TooLong(t *testing.T) {
	ctx := context.Background()
	testUserID := "user-123"
	testReviewID := "review-123"

	mockRepo := testutil.NewMockReviewRepository()
	useCase := NewUpdateFeedbackUseCase(mockRepo)

	longComment := ""
	for i := 0; i < 501; i++ {
		longComment += "a"
	}

	input := UpdateFeedbackInput{
		ReviewID: testReviewID,
		UserID:   testUserID,
		Score:    3,
		Comment:  longComment,
	}

	_, err := useCase.Execute(ctx, input)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "コメントは500文字以内にしてください"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TC-RV-004-13: 他人のレビューにフィードバック
func TestUpdateFeedbackUseCase_Execute_Forbidden_OtherUserReview(t *testing.T) {
	ctx := context.Background()
	testUserID := "user-123"
	otherUserID := "user-456"
	testReviewID := "review-123"
	
	testReview := createTestReview(otherUserID) // 他のユーザーのレビュー
	testReview.ID = testReviewID

	mockRepo := testutil.NewMockReviewRepository()
	mockRepo.Create(ctx, testReview)

	useCase := NewUpdateFeedbackUseCase(mockRepo)

	input := UpdateFeedbackInput{
		ReviewID: testReviewID,
		UserID:   testUserID, // 別のユーザー
		Score:    3,
		Comment:  "",
	}

	_, err := useCase.Execute(ctx, input)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "このレビューを更新する権限がありません"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

// TC-RV-004-14: 存在しないreview_id
func TestUpdateFeedbackUseCase_Execute_NotFound_InvalidReviewID(t *testing.T) {
	ctx := context.Background()
	testUserID := "user-123"
	testReviewID := "invalid-review-id"

	mockRepo := testutil.NewMockReviewRepository()
	mockRepo.SetError(errors.New("review not found"))

	useCase := NewUpdateFeedbackUseCase(mockRepo)

	input := UpdateFeedbackInput{
		ReviewID: testReviewID,
		UserID:   testUserID,
		Score:    3,
		Comment:  "",
	}

	_, err := useCase.Execute(ctx, input)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	// エラーメッセージに "レビューが見つかりません" が含まれることを確認
	if !contains(err.Error(), "レビューが見つかりません") {
		t.Errorf("Expected error to contain 'レビューが見つかりません', got '%s'", err.Error())
	}
}

// TC-RV-004-04: スコア=1（Bad）
func TestUpdateFeedbackUseCase_Execute_Success_ScoreBad(t *testing.T) {
	ctx := context.Background()
	testUserID := "user-123"
	testReviewID := "review-123"
	testReview := createTestReview(testUserID)
	testReview.ID = testReviewID

	mockRepo := testutil.NewMockReviewRepository()
	mockRepo.Create(ctx, testReview)

	useCase := NewUpdateFeedbackUseCase(mockRepo)

	input := UpdateFeedbackInput{
		ReviewID: testReviewID,
		UserID:   testUserID,
		Score:    1,
		Comment:  "",
	}

	output, err := useCase.Execute(ctx, input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if output.FeedbackScore != 1 {
		t.Errorf("Expected score 1, got %d", output.FeedbackScore)
	}
}

// TC-RV-004-05: スコア=2（Normal）
func TestUpdateFeedbackUseCase_Execute_Success_ScoreNormal(t *testing.T) {
	ctx := context.Background()
	testUserID := "user-123"
	testReviewID := "review-123"
	testReview := createTestReview(testUserID)
	testReview.ID = testReviewID

	mockRepo := testutil.NewMockReviewRepository()
	mockRepo.Create(ctx, testReview)

	useCase := NewUpdateFeedbackUseCase(mockRepo)

	input := UpdateFeedbackInput{
		ReviewID: testReviewID,
		UserID:   testUserID,
		Score:    2,
		Comment:  "",
	}

	output, err := useCase.Execute(ctx, input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if output.FeedbackScore != 2 {
		t.Errorf("Expected score 2, got %d", output.FeedbackScore)
	}
}

// ヘルパー関数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
