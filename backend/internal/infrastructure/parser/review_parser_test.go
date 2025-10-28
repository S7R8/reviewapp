package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseReviewMarkdown_ActualResponse(t *testing.T) {
	// 実際にClaudeから返ってきたレスポンス
	markdown := `## コードレビュー結果

### 良い点
- 二分探索のアルゴリズムが基本的に正しく実装されている
- シンプルで読みやすいコード構造
- テストケースが含まれており、異なる入力に対する動作を確認している

### 1. オーバーフロー脆弱性の修正

- 現在` + "`mid`" + `の計算方法は大きな整数配列でオーバーフローのリスクがある
- Pythonでは自動的に大きな整数に対応するが、他の言語では問題になる可能性がある

改善例：
` + "```python" + `
def binary_search(arr, target):
    low = 0
    high = len(arr) - 1
    
    while low <= high:
        # より安全な中央値の計算方法
        mid = low + (high - low) // 2
        guess = arr[mid]
        
        if guess == target:
            return mid
        if guess > target:
            high = mid - 1
        else:
            low = mid + 1
    return None
` + "```" + `

### 2. エラーハンドリングの改善

- 現在のコードはエラーハンドリングが不足している
- 入力の型チェックや空配列のケースを考慮していない

改善例：
` + "```python" + `
import logging

def binary_search(arr, target):
    # 入力バリデーション
    if not isinstance(arr, list):
        logging.error(f"Invalid input type: {type(arr)}")
        raise TypeError("Input must be a list")
    
    if not arr:
        logging.warning("Empty array provided")
        return None
    
    low = 0
    high = len(arr) - 1
    
    try:
        while low <= high:
            mid = low + (high - low) // 2
            guess = arr[mid]
            
            if guess == target:
                return mid
            if guess > target:
                high = mid - 1
            else:
                low = mid + 1
        return None
    except Exception as e:
        logging.error(f"Unexpected error in binary search: {e}", exc_info=True)
        raise
` + "```" + `

### 3. テストカバレッジの向上

- 現在のテストは限定的
- より多様なケースをテストする必要がある

改善例：
` + "```python" + `
import unittest

class TestBinarySearch(unittest.TestCase):
    def test_normal_case(self):
        arr = [1, 3, 5, 7, 9]
        self.assertEqual(binary_search(arr, 3), 1)
    
    def test_first_element(self):
        arr = [1, 3, 5, 7, 9]
        self.assertEqual(binary_search(arr, 1), 0)
    
    def test_last_element(self):
        arr = [1, 3, 5, 7, 9]
        self.assertEqual(binary_search(arr, 9), 4)
    
    def test_not_found(self):
        arr = [1, 3, 5, 7, 9]
        self.assertIsNone(binary_search(arr, -1))
    
    def test_empty_array(self):
        self.assertIsNone(binary_search([], 5))

if __name__ == '__main__':
    unittest.main()
` + "```" + `

### 総合評価
基本的な二分探索アルゴリズムは正しく実装されていますが、エラーハンドリング、入力検証、テストカバレッジの観点で改善の余地があります。`

	result := ParseReviewMarkdown(markdown)

	// デバッグ出力
	t.Logf("\n=== Parse Result ===")
	t.Logf("Summary: %s", result.Summary)
	t.Logf("\nGood Points (%d):", len(result.GoodPoints))
	for i, point := range result.GoodPoints {
		t.Logf("  [%d] %s", i+1, point)
	}
	t.Logf("\nImprovements (%d):", len(result.Improvements))
	for i, imp := range result.Improvements {
		t.Logf("\n  === Improvement %d ===", i+1)
		t.Logf("  Title: %s", imp.Title)
		t.Logf("  Description: %s", imp.Description)
		t.Logf("  Severity: %s", imp.Severity)
		t.Logf("  Has code: %v", imp.CodeAfter != "")
		if imp.CodeAfter != "" {
			t.Logf("  Code lines: %d", len(splitLines(imp.CodeAfter)))
		}
	}

	// アサーション
	assert.NotEmpty(t, result.Summary, "Summary should not be empty")
	assert.Contains(t, result.Summary, "正しく実装", "Summary should contain expected text")

	// 良い点のチェック
	assert.Equal(t, 3, len(result.GoodPoints), "Should have 3 good points")
	assert.Contains(t, result.GoodPoints[0], "二分探索")
	assert.Contains(t, result.GoodPoints[1], "シンプル")
	assert.Contains(t, result.GoodPoints[2], "テストケース")

	// 改善点のチェック
	assert.Equal(t, 3, len(result.Improvements), "Should have 3 improvements")

	// 1つ目: オーバーフロー脆弱性
	assert.Contains(t, result.Improvements[0].Title, "オーバーフロー")
	assert.Contains(t, result.Improvements[0].Description, "計算方法")
	assert.Equal(t, "high", result.Improvements[0].Severity, "オーバーフロー'脆弱性'なので high")
	assert.NotEmpty(t, result.Improvements[0].CodeAfter, "Should have code for improvement 1")
	assert.Contains(t, result.Improvements[0].CodeAfter, "low + (high - low)")

	// 2つ目: エラーハンドリング
	assert.Contains(t, result.Improvements[1].Title, "エラーハンドリング")
	assert.Contains(t, result.Improvements[1].Description, "不足")
	assert.Equal(t, "high", result.Improvements[1].Severity)
	assert.NotEmpty(t, result.Improvements[1].CodeAfter, "Should have code for improvement 2")

	// 3つ目: テスト
	assert.Contains(t, result.Improvements[2].Title, "テスト")
	assert.Contains(t, result.Improvements[2].Description, "限定的")
	assert.Equal(t, "medium", result.Improvements[2].Severity)
	assert.NotEmpty(t, result.Improvements[2].CodeAfter, "Should have code for improvement 3")
}

func splitLines(s string) []string {
	lines := []string{}
	current := ""
	for _, ch := range s {
		if ch == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}
