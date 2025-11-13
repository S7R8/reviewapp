/**
 * ナレッジの型定義（仮）
 */
export interface Knowledge {
  id: string;
  title: string;
  description: string;
  category: string;
  tags: string[];
}

/**
 * レビューステータスの型定義
 */
export type ReviewStatus = 'success' | 'warning' | 'error' | 'pending';

/**
 * プログラミング言語の型定義
 */
export type ProgrammingLanguage =
  | 'TypeScript'
  | 'JavaScript'
  | 'Python'
  | 'Go'
  | 'Java'
  | 'C++'
  | 'C#'
  | 'Ruby'
  | 'PHP'
  | 'Rust'
  | 'Swift'
  | 'Kotlin'
  | 'SCSS'
  | 'CSS'
  | 'HTML'
  | 'Markdown'
  | 'Other';

/**
 * レビュー履歴アイテムの型定義
 */
export interface ReviewHistoryItem {
  id: string;
  createdAt: string; // ISO 8601 format
  language: ProgrammingLanguage;
  status: ReviewStatus;
  reviewContent?: string;
  knowledgeReferences?: string[]; // 参照されたナレッジのID
}

/**
 * レビュー履歴リストのレスポンス型
 */
export interface ReviewHistoryListResponse {
  items: ReviewHistoryItem[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}

/**
 * レビュー履歴の検索・フィルター条件
 */
export interface ReviewHistoryFilter {
  language?: ProgrammingLanguage;
  status?: ReviewStatus;
  dateFrom?: string;
  dateTo?: string;
  page?: number;
  pageSize?: number;
  sortBy?: 'createdAt' | 'language' | 'status';
  sortOrder?: 'asc' | 'desc';
}

/**
 * ステータスバッジの設定
 */
export interface StatusBadgeConfig {
  label: string;
  bgColor: string;
  textColor: string;
}

/**
 * レビューリクエストの型定義
 */
export interface ReviewRequest {
  code: string;
  language: string;
  filename?: string;
}

/**
 * レビュー結果の型定義
 */
export interface ReviewResult {
  id: string;
  summary: string;
  goodPoints: string[];
  improvements: Array<{
    title: string;
    description: string;
    codeAfter?: string;
    severity: 'low' | 'medium' | 'high';
  }>;
  references: Array<{
    source: string;
    description: string;
  }>;
  referencedKnowledgeIds?: string[]; // ★ 追加
  createdAt: string;
  rawMarkdown: string;
}
