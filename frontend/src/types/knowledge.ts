import { KnowledgeCategory } from '../api/knowledgeApi';

// ナレッジ一覧アイテム
export interface KnowledgeListItem {
  id: string;
  title: string;
  content: string;
  category: KnowledgeCategory;
  priority: number;  // 1-3
  usage_count: number;
  last_used_at: string | null;
  created_at: string;
  updated_at: string;
}

// ナレッジフィルタ
export interface KnowledgeFilter {
  category?: KnowledgeCategory;
  priority?: number;  // 1, 2, 3
  page: number;
  pageSize: number;
  sortBy: 'created_at' | 'priority' | 'usage_count';
  sortOrder: 'asc' | 'desc';
}

// 優先度ラベル
export const getPriorityLabel = (priority: number): string => {
  const labels: Record<number, string> = {
    1: '低',
    2: '中',
    3: '高'
  };
  return labels[priority] || '-';
};

// 優先度を星で表示
export const getPriorityStars = (priority: number): string => {
  const stars: Record<number, string> = {
    1: '★☆☆',
    2: '★★☆',
    3: '★★★'
  };
  return stars[priority] || '☆☆☆';
};
