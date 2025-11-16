package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/s7r8/reviewapp/internal/domain/model"
)

// ReviewService - レビュードメインサービス
type ReviewService struct{}

// NewReviewService - コンストラクタ
func NewReviewService() *ReviewService {
	return &ReviewService{}
}

// BuildPromptFromKnowledge - ナレッジからプロンプトを生成し、実際に使用したナレッジを返す
func (s *ReviewService) BuildPromptFromKnowledge(knowledges []*model.Knowledge) (string, []*model.Knowledge) {
	// ナレッジが存在しない場合
	if len(knowledges) == 0 {
		return "一般的なベストプラクティスに基づいてレビューしてください。", []*model.Knowledge{}
	}

	sortedKnowledges := make([]*model.Knowledge, len(knowledges))
	copy(sortedKnowledges, knowledges)
	sort.Slice(sortedKnowledges, func(i, j int) bool {
		if sortedKnowledges[i].Priority == sortedKnowledges[j].Priority {
			return sortedKnowledges[i].CreatedAt.After(sortedKnowledges[j].CreatedAt)
		}
		return sortedKnowledges[i].Priority > sortedKnowledges[j].Priority
	})

	var sb strings.Builder

	limit := 10
	if len(sortedKnowledges) < limit {
		limit = len(sortedKnowledges)
	}

	// 実際に使用したナレッジを記録
	usedKnowledges := make([]*model.Knowledge, limit)

	for i := 0; i < limit; i++ {
		k := sortedKnowledges[i]
		usedKnowledges[i] = k

		categoryName := s.getCategoryName(k.Category)
		sb.WriteString(fmt.Sprintf("### [%s] %s\n", categoryName, k.Title))
		sb.WriteString(fmt.Sprintf("%s\n\n", k.Content))
	}

	return sb.String(), usedKnowledges
}

// getCategoryName - カテゴリIDを日本語名に変換
// TODO DBからマッピングを取得するよう修正
func (s *ReviewService) getCategoryName(category string) string {
	categoryMap := map[string]string{
		"error_handling": "エラーハンドリング",
		"testing":        "テスト",
		"performance":    "パフォーマンス",
		"security":       "セキュリティ",
		"clean_code":     "クリーンコード",
		"architecture":   "アーキテクチャ",
		"other":          "その他",
	}

	if name, ok := categoryMap[category]; ok {
		return name
	}
	return category
}
