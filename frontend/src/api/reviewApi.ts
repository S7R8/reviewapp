import { ReviewResult, ReviewRequest } from '../types/review';
import { apiPost, apiPut } from './client';

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
  feedback_score: number | null;
  feedback_comment: string | null;
  created_at: string;
  updated_at: string;
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
}

export const reviewApiClient = new ReviewApiClient();
