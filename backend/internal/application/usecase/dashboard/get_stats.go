package dashboard

import (
	"context"
	"math"
	"time"

	"github.com/s7r8/reviewapp/internal/domain/model"
	"github.com/s7r8/reviewapp/internal/domain/repository"
)

// GetStatsUseCase - ダッシュボード統計取得UseCase
type GetStatsUseCase struct {
	reviewRepo    repository.ReviewRepository
	knowledgeRepo repository.KnowledgeRepository
}

// NewGetStatsUseCase - コンストラクタ
func NewGetStatsUseCase(
	reviewRepo repository.ReviewRepository,
	knowledgeRepo repository.KnowledgeRepository,
) *GetStatsUseCase {
	return &GetStatsUseCase{
		reviewRepo:    reviewRepo,
		knowledgeRepo: knowledgeRepo,
	}
}

// DashboardStats - ダッシュボード統計情報
type DashboardStats struct {
	TotalReviews     int               `json:"total_reviews"`
	KnowledgeCount   int               `json:"knowledge_count"`
	ConsistencyScore int               `json:"consistency_score"`
	WeeklyReviews    int               `json:"weekly_reviews"`
}

// RecentReviewItem - 最近のレビューアイテム
type RecentReviewItem struct {
	ID                string `json:"id"`
	FileName          string `json:"file_name"`
	Language          string `json:"language"`
	CreatedAt         string `json:"created_at"`
	ImprovementsCount int    `json:"improvements_count"`
	Status            string `json:"status"`
}

// SkillAnalysis - スキル分析（カテゴリ別のナレッジ割合）
type SkillAnalysis struct {
	ErrorHandling int `json:"error_handling"`
	Testing       int `json:"testing"`
	Performance   int `json:"performance"`
	Security      int `json:"security"`
	CleanCode     int `json:"clean_code"`
	Architecture  int `json:"architecture"`
	Other         int `json:"other"`
}

// DashboardStatsResponse - レスポンス全体
type DashboardStatsResponse struct {
	Stats          DashboardStats     `json:"stats"`
	RecentReviews  []RecentReviewItem `json:"recent_reviews"`
	SkillAnalysis  SkillAnalysis      `json:"skill_analysis"`
}

// Execute - UseCase実行
func (u *GetStatsUseCase) Execute(ctx context.Context, userID string) (*DashboardStatsResponse, error) {
	// 1. 総レビュー回数を取得
	totalReviews, err := u.reviewRepo.CountByUserID(ctx, userID)
	if err != nil {
		totalReviews = 0 // エラーでも0として扱う
	}

	// 2. 総ナレッジ数を取得（有効なもののみ）
	knowledgeCount, err := u.knowledgeRepo.CountByUserID(ctx, userID)
	if err != nil {
		knowledgeCount = 0
	}

	// 3. 一貫性スコアを計算
	averageScore, err := u.reviewRepo.GetAverageFeedbackScore(ctx, userID)
	if err != nil {
		averageScore = 0
	}
	consistencyScore := calculateConsistencyScore(averageScore)

	// 4. 今週のレビュー回数を取得
	from, to := getThisWeekRange()
	weeklyReviews, err := u.reviewRepo.CountByUserIDAndDateRange(ctx, userID, from, to)
	if err != nil {
		weeklyReviews = 0
	}

	// 5. 最近のレビューを取得（最大5件）
	recentReviewsData, err := u.reviewRepo.FindRecentByUserID(ctx, userID, 5)
	if err != nil {
		recentReviewsData = []*model.Review{}
	}
	recentReviews := convertToRecentReviewItems(recentReviewsData)

	// 6. カテゴリ別ナレッジ数を取得
	categoryCounts, err := u.knowledgeRepo.CountByCategory(ctx, userID)
	if err != nil {
		categoryCounts = make(map[string]int)
	}
	skillAnalysis := calculateSkillPercentages(categoryCounts)

	// レスポンスを構築
	response := &DashboardStatsResponse{
		Stats: DashboardStats{
			TotalReviews:     totalReviews,
			KnowledgeCount:   knowledgeCount,
			ConsistencyScore: consistencyScore,
			WeeklyReviews:    weeklyReviews,
		},
		RecentReviews: recentReviews,
		SkillAnalysis: skillAnalysis,
	}

	return response, nil
}

// calculateConsistencyScore - 一貫性スコアを計算
// フィードバックスコア（1-3）を0-100%に変換
// 1 → 0%, 2 → 50%, 3 → 100%
func calculateConsistencyScore(averageScore float64) int {
	if averageScore == 0 {
		return 0
	}
	score := ((averageScore - 1) / 2) * 100
	return int(math.Round(score))
}

// getThisWeekRange - 今週の範囲を取得（月曜 00:00:00 から現在まで）
func getThisWeekRange() (time.Time, time.Time) {
	now := time.Now()
	weekday := now.Weekday()

	// 月曜日を週の開始とする
	daysFromMonday := int(weekday) - 1
	if daysFromMonday < 0 {
		daysFromMonday = 6 // 日曜日の場合
	}

	monday := now.AddDate(0, 0, -daysFromMonday)
	mondayStart := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())

	return mondayStart, now
}

// calculateSkillPercentages - カテゴリ別ナレッジ数を％に変換
func calculateSkillPercentages(categoryCounts map[string]int) SkillAnalysis {
	total := 0
	for _, count := range categoryCounts {
		total += count
	}

	if total == 0 {
		return SkillAnalysis{
			ErrorHandling: 0,
			Testing:       0,
			Performance:   0,
			Security:      0,
			CleanCode:     0,
			Architecture:  0,
			Other:         0,
		}
	}

	return SkillAnalysis{
		ErrorHandling: calculatePercentage(categoryCounts[model.CategoryErrorHandling], total),
		Testing:       calculatePercentage(categoryCounts[model.CategoryTesting], total),
		Performance:   calculatePercentage(categoryCounts[model.CategoryPerformance], total),
		Security:      calculatePercentage(categoryCounts[model.CategorySecurity], total),
		CleanCode:     calculatePercentage(categoryCounts[model.CategoryCleanCode], total),
		Architecture:  calculatePercentage(categoryCounts[model.CategoryArchitecture], total),
		Other:         calculatePercentage(categoryCounts[model.CategoryOther], total),
	}
}

// calculatePercentage - 割合を計算
func calculatePercentage(count, total int) int {
	if total == 0 {
		return 0
	}
	return int(math.Round(float64(count) / float64(total) * 100))
}

// convertToRecentReviewItems - Review配列をRecentReviewItem配列に変換
func convertToRecentReviewItems(reviews []*model.Review) []RecentReviewItem {
	items := make([]RecentReviewItem, 0, len(reviews))
	
	for _, review := range reviews {
		// ファイル名を決定（CodeからファイルっぽいものをパースするかUntitledにする）
		fileName := "Untitled"
		
		// 改善点の数をカウント
		improvementsCount := 0
		if review.StructuredResult != nil {
			improvementsCount = len(review.StructuredResult.Improvements)
		}
		
		// ステータスを決定（改善点の数に基づく）
		status := "success"
		if improvementsCount > 0 {
			status = "warning"
		}
		
		items = append(items, RecentReviewItem{
			ID:                review.ID,
			FileName:          fileName,
			Language:          review.Language,
			CreatedAt:         review.CreatedAt.Format(time.RFC3339),
			ImprovementsCount: improvementsCount,
			Status:            status,
		})
	}
	
	return items
}
