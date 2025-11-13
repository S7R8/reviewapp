import { ReviewResult, ReviewRequest, ReviewHistoryFilter, ReviewHistoryListResponse } from '../types/review';
import { apiPost, apiPut, apiGet } from './client';

interface ApiReviewResponse {
  id: string;
  user_id: string;
  code: string;
  language: string;
  file_name: string | null;
  review_result: string;
  structured_result?: {
    summary: string;
    good_points: string[];
    improvements: Array<{
      title: string;
      description: string;
      code_after: string;
      severity: string;
    }>;
  };
  referenced_knowledge?: Array<{  // ★ 追加
    id: string;
    title: string;
    category: string;
    priority: number;
  }>;
  feedback_score: number | null;
  feedback_comment: string | null;
  created_at: string;
  updated_at: string;
}

interface ApiReviewHistoryResponse {
  items: Array<{
    id: string;
    user_id: string;
    code: string;
    language: string;
    status: string;
    review_result: string;
    knowledge_references: string[];
    created_at: string;
    updated_at: string;
  }>;
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

class ReviewApiClient {
  async reviewCode(request: ReviewRequest): Promise<ReviewResult> {
    const data = await apiPost<ApiReviewResponse>('/api/v1/reviews', {
      code: request.code,
      language: request.language,
      file_name: request.filename || null,
      context: null,
    });

    // バックエンドから構造化データを受け取る
    const structured = data.structured_result || {
      summary: 'レビュー結果を確認してください',
      good_points: [],
      improvements: [],
    };

    return {
      id: data.id,
      summary: structured.summary,
      goodPoints: structured.good_points || [],
      improvements: (structured.improvements || []).map((imp: any) => ({
        title: imp.title,
        description: imp.description,
        codeAfter: imp.code_after,
        severity: imp.severity || 'medium',
      })),
      references: [],
      referencedKnowledgeIds: (data.referenced_knowledge || []).map(k => k.id), // ★ 追加
      createdAt: data.created_at,
      rawMarkdown: data.review_result,
    };
  }

  async updateFeedback(reviewId: string, score: number, comment?: string): Promise<void> {
    await apiPut(`/api/v1/reviews/${reviewId}/feedback`, {
      score,
      comment: comment || '',
    });
  }

  /**
   * レビュー履歴一覧を取得
   */
  async getReviewHistory(filter: ReviewHistoryFilter): Promise<ReviewHistoryListResponse> {
    // クエリパラメータを構築
    const params = new URLSearchParams();
    
    if (filter.page) params.append('page', filter.page.toString());
    if (filter.pageSize) params.append('page_size', filter.pageSize.toString());
    if (filter.language) params.append('language', filter.language);
    if (filter.status) params.append('status', filter.status);
    if (filter.sortBy) params.append('sort_by', filter.sortBy === 'createdAt' ? 'created_at' : filter.sortBy);
    if (filter.sortOrder) params.append('sort_order', filter.sortOrder);
    if (filter.dateFrom) params.append('date_from', filter.dateFrom);
    if (filter.dateTo) params.append('date_to', filter.dateTo);

    const queryString = params.toString();
    const endpoint = `/api/v1/reviews${queryString ? `?${queryString}` : ''}`;

    // スケルトン表示のためにわざと遅延を追加（開発用）
    const [data] = await Promise.all([
      apiGet<ApiReviewHistoryResponse>(endpoint),
      new Promise(resolve => setTimeout(resolve, 800)), // 800ms遅延
    ]);

    // バックエンドのレスポンスをフロントエンドの型に変換
    return {
      items: data.items.map((item) => ({
        id: item.id,
        createdAt: item.created_at,
        language: item.language as any, // ProgrammingLanguageにキャスト
        status: item.status as any, // ReviewStatusにキャスト
        reviewContent: item.review_result,
        knowledgeReferences: item.knowledge_references,
      })),
      total: data.total,
      page: data.page,
      pageSize: data.page_size,
      totalPages: data.total_pages,
    };
  }
}

export const reviewApiClient = new ReviewApiClient();
