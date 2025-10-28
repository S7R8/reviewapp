package parser

import (
	"regexp"
	"strings"

	"github.com/s7r8/reviewapp/internal/domain/model"
)

// ParseReviewMarkdown - マークダウン形式のレビュー結果を構造化データに変換
func ParseReviewMarkdown(markdown string) *model.StructuredReviewResult {
	result := &model.StructuredReviewResult{
		GoodPoints:   []string{},
		Improvements: []model.Improvement{},
	}

	// サマリーを抽出
	result.Summary = extractSummary(markdown)

	// 良い点を抽出
	result.GoodPoints = extractGoodPoints(markdown)

	// 改善点を抽出
	result.Improvements = extractImprovements(markdown)

	return result
}

// extractSummary - 総合評価セクションから抽出
func extractSummary(text string) string {
	// ### 総合評価 または ## 総合評価
	re := regexp.MustCompile(`###+?\s*総合評価\s*\n([\s\S]*?)(?:\n##|$)`)
	if match := re.FindStringSubmatch(text); len(match) > 1 {
		lines := strings.Split(match[1], "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				return trimmed
			}
		}
	}
	return "詳細は下記をご確認ください。"
}

// extractGoodPoints - 良い点セクションから箇条書きを抽出
func extractGoodPoints(text string) []string {
	points := []string{}

	// ### 良い点 または ## 良い点
	re := regexp.MustCompile(`###+?\s*良い点\s*\n([\s\S]*?)(?:\n##|$)`)
	if match := re.FindStringSubmatch(text); len(match) > 1 {
		content := match[1]
		lines := strings.Split(content, "\n")

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			// - で始まる行を抽出
			if strings.HasPrefix(trimmed, "-") {
				point := strings.TrimPrefix(trimmed, "-")
				point = strings.TrimSpace(point)
				if point != "" {
					points = append(points, point)
				}
			}
		}
	}

	if len(points) == 0 {
		points = append(points, "コードの基本的な構造は良好です")
	}

	return points
}

// extractImprovements - 改善点セクションを抽出
func extractImprovements(text string) []model.Improvement {
	improvements := []model.Improvement{}

	// ### 1. タイトル または ## 1. タイトル の形式でセクションを探す
	sectionRe := regexp.MustCompile(`###+?\s*(\d+)\.\s+(.+)\n`)
	matches := sectionRe.FindAllStringSubmatchIndex(text, -1)

	for i, match := range matches {
		// タイトルを取得
		title := text[match[4]:match[5]]

		// セクションの内容範囲を決定
		contentStart := match[1]
		contentEnd := len(text)

		// 次のセクションまたは ## 総合評価 で終了
		if i+1 < len(matches) {
			contentEnd = matches[i+1][0]
		} else {
			// 総合評価セクションを探す
			summaryRe := regexp.MustCompile(`\n###+?\s*総合評価`)
			if loc := summaryRe.FindStringIndex(text[contentStart:]); loc != nil {
				contentEnd = contentStart + loc[0]
			}
		}

		content := text[contentStart:contentEnd]

		// 説明とコードを抽出
		description := extractDescription(content)
		codeAfter := extractCodeBlock(content)

		// 重要度を判定
		severity := determineSeverity(title, description)

		improvements = append(improvements, model.Improvement{
			Title:       strings.TrimSpace(title),
			Description: description,
			CodeAfter:   codeAfter,
			Severity:    severity,
		})
	}

	return improvements
}

// extractDescription - 説明文を抽出（箇条書き + 通常テキスト）
func extractDescription(content string) string {
	lines := strings.Split(content, "\n")
	descriptions := []string{}

	inCodeBlock := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// コードブロックをスキップ
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}

		// "改善例：" "改善案：" をスキップ
		if regexp.MustCompile(`^改善[例案][:：]`).MatchString(trimmed) {
			continue
		}

		// 空行をスキップ
		if trimmed == "" {
			continue
		}

		// - で始まる箇条書きを抽出
		if strings.HasPrefix(trimmed, "-") {
			desc := strings.TrimPrefix(trimmed, "-")
			desc = strings.TrimSpace(desc)
			descriptions = append(descriptions, desc)
		}
	}

	if len(descriptions) == 0 {
		return "改善が推奨されます"
	}

	return strings.Join(descriptions, "\n")
}

// extractCodeBlock - コードブロックを抽出
func extractCodeBlock(content string) string {
	codeRe := regexp.MustCompile("```[a-z]*\\n([\\s\\S]*?)```")
	if match := codeRe.FindStringSubmatch(content); len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

// determineSeverity - 重要度を自動判定
func determineSeverity(title, description string) string {
	text := strings.ToLower(title + " " + description)

	// 高優先度キーワード
	highKeywords := []string{
		"重大", "脆弱性", "エラーハンドリング", "エラー処理",
		"セキュリティ", "危険", "バグ", "クリティカル",
	}
	for _, keyword := range highKeywords {
		if strings.Contains(text, keyword) {
			return "high"
		}
	}

	// 中優先度キーワード
	mediumKeywords := []string{
		"パフォーマンス", "効率", "最適化",
		"クリーンコード", "保守性", "可読性",
		"テスト", "ドキュメント",
	}
	for _, keyword := range mediumKeywords {
		if strings.Contains(text, keyword) {
			return "medium"
		}
	}

	return "low"
}
