import { apiPost, apiGet } from './client';

// ナレッジカテゴリ
export type KnowledgeCategory = 
  | 'error_handling'
  | 'testing'
  | 'performance'
  | 'security'
  | 'clean_code'
  | 'architecture'
  | 'other';

// ナレッジ作成リクエスト
export interface KnowledgeCreateRequest {
  title: string;
  content: string;
  category: KnowledgeCategory;
  priority: number; // 1-3
}

// ナレッジ一覧フィルタ
export interface KnowledgeListFilter {
  category?: KnowledgeCategory;
}

// ナレッジエンティティ
export interface Knowledge {
  id: string;
  user_id: string;
  title: string;
  content: string;
  category: KnowledgeCategory;
  priority: number;
  source_type: string;
  source_id: string | null;
  usage_count: number;
  last_used_at: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

// カテゴリラベル取得
export const getCategoryLabel = (category: KnowledgeCategory): string => {
  const labels: Record<KnowledgeCategory, string> = {
    error_handling: 'エラーハンドリング',
    testing: 'テスト',
    performance: 'パフォーマンス',
    security: 'セキュリティ',
    clean_code: 'クリーンコード',
    architecture: 'アーキテクチャ',
    other: 'その他',
  };
  return labels[category];
};

// ナレッジAPIクライアント
class KnowledgeApiClient {
  /**
   * ナレッジを作成
   */
  async createKnowledge(request: KnowledgeCreateRequest): Promise<Knowledge> {
    return await apiPost<Knowledge>('/api/v1/knowledge', request);
  }

  /**
   * ナレッジ一覧を取得
   */
  async listKnowledge(category?: KnowledgeCategory): Promise<Knowledge[]> {
    const query = category ? `?category=${category}` : '';
    return await apiGet<Knowledge[]>(`/api/v1/knowledge${query}`);
  }

  /**
   * ナレッジ詳細を取得（ID指定）
   */
  async getKnowledgeById(id: string): Promise<Knowledge> {
    return await apiGet<Knowledge>(`/api/v1/knowledge/${id}`);
  }
}

export const knowledgeApiClient = new KnowledgeApiClient();
