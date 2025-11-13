import { ReviewResult } from '../types/review';
import { KnowledgeCreateRequest, KnowledgeCategory } from '../api/knowledgeApi';

/**
 * カテゴリ推定用キーワードマッピング
 */
const CATEGORY_KEYWORDS: Record<KnowledgeCategory, string[]> = {
  error_handling: [
    'エラー', 'error', '例外', 'exception', 'エラーハンドリング',
    'try', 'catch', 'panic', 'recover', 'throw', 'ハンドリング'
  ],
  testing: [
    'テスト', 'test', 'testing', 'カバレッジ', 'coverage',
    'ユニットテスト', 'unit test', 'モック', 'mock', 'テストケース'
  ],
  performance: [
    'パフォーマンス', 'performance', '最適化', 'optimization',
    '速度', 'speed', 'メモリ', 'memory', '効率', 'efficiency',
    '高速化', 'キャッシュ', 'cache'
  ],
  security: [
    'セキュリティ', 'security', '脆弱性', 'vulnerability',
    'XSS', 'SQL injection', '認証', 'authentication', '暗号化', 'encryption',
    'インジェクション', 'サニタイズ', 'sanitize'
  ],
  clean_code: [
    'クリーンコード', 'clean code', '可読性', 'readability',
    'リファクタリング', 'refactoring', '命名', 'naming', 'コメント', 'comment',
    '変数名', '関数名', 'メソッド名', '整理'
  ],
  architecture: [
    'アーキテクチャ', 'architecture', '設計', 'design',
    '構造', 'structure', 'パターン', 'pattern', '依存', 'dependency',
    'SOLID', 'DDD', 'クリーンアーキテクチャ'
  ],
  other: []
};

/**
 * レビュー結果からカテゴリを推定
 */
function estimateCategory(review: ReviewResult): KnowledgeCategory {
  // 改善点のタイトルと説明を結合
  const text = [
    review.summary,
    ...review.improvements.map(i => `${i.title} ${i.description}`)
  ].join(' ').toLowerCase();

  // 各カテゴリのマッチ数をカウント
  const scores: Record<string, number> = {};
  
  for (const [category, keywords] of Object.entries(CATEGORY_KEYWORDS)) {
    if (category === 'other') continue;
    
    scores[category] = keywords.filter(keyword => 
      text.includes(keyword.toLowerCase())
    ).length;
  }

  // 最もマッチ数が多いカテゴリを返す
  const bestMatch = Object.entries(scores)
    .sort((a, b) => b[1] - a[1])[0];

  return (bestMatch && bestMatch[1] > 0) 
    ? bestMatch[0] as KnowledgeCategory
    : 'other';
}

/**
 * レビュー結果から優先度を推定
 * severity → priority のマッピング
 * high → 3, medium → 2, low → 1
 */
function estimatePriority(review: ReviewResult): number {
  if (review.improvements.length === 0) {
    return 2; // デフォルト: medium
  }

  const severities = review.improvements.map(i => i.severity);
  
  if (severities.includes('high')) return 3;
  if (severities.includes('medium')) return 2;
  return 1; // low
}

/**
 * タイトルを生成
 * 形式: [言語] 最初の改善点のタイトル
 */
function generateTitle(review: ReviewResult, language: string): string {
  const lang = language.charAt(0).toUpperCase() + language.slice(1);
  
  // 最初の改善点のタイトルを使用
  if (review.improvements.length > 0) {
    const firstImprovement = review.improvements[0].title;
    return `[${lang}] ${firstImprovement}`;
  }
  
  // 改善点がない場合はsummaryの最初の一文
  const summaryFirstLine = review.summary
    .split('\n')[0]
    .replace(/^【.*?】/g, '')  // 【】を除去
    .trim()
    .slice(0, 50);
  
  return `[${lang}] ${summaryFirstLine}`;
}

/**
 * コンテンツを生成
 */
function generateContent(review: ReviewResult, language: string): string {
  const now = new Date();
  const dateStr = now.toLocaleString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  });

  const lang = language.charAt(0).toUpperCase() + language.slice(1);

  let content = `【言語】${lang}\n`;
  content += `【レビュー日時】${dateStr}\n\n`;
  
  // 総評
  content += `【総評】\n${review.summary}\n\n`;
  
  // 改善点
  if (review.improvements.length > 0) {
    content += `【改善点】\n`;
    review.improvements.forEach((improvement, index) => {
      const severity = improvement.severity.toUpperCase();
      content += `${index + 1}. ${improvement.title} (${severity})\n`;
      content += `   ${improvement.description}\n`;
      
      if (improvement.codeAfter) {
        content += `\n   改善例:\n\`\`\`${language}\n${improvement.codeAfter}\n\`\`\`\n`;
      }
      content += '\n';
    });
  }
  
  // 良い点
  if (review.goodPoints.length > 0) {
    content += `【良い点】\n`;
    review.goodPoints.forEach(point => {
      content += `- ${point}\n`;
    });
    content += '\n';
  }
  
  // 参照情報
  if (review.references.length > 0) {
    content += `【参照情報】\n`;
    review.references.forEach(ref => {
      content += `- ${ref.source}: ${ref.description}\n`;
    });
  }

  return content.trim();
}

/**
 * レビュー結果からナレッジ作成リクエストを生成
 */
export function createKnowledgeFromReview(
  review: ReviewResult,
  language: string
): KnowledgeCreateRequest {
  const category = estimateCategory(review);
  
  return {
    title: generateTitle(review, language),
    content: generateContent(review, language),
    category,
    priority: estimatePriority(review)
  };
}
