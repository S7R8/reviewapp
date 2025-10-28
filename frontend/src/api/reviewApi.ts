import { ReviewResult, ReviewRequest } from '../types/review';

const API_BASE_URL = 'http://localhost:8080';

interface ApiError {
  error: string;
  message: string;
  details?: Record<string, unknown>;
}

class ReviewApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  async reviewCode(request: ReviewRequest): Promise<ReviewResult> {
    const response = await fetch(`${this.baseUrl}/api/v1/reviews`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        code: request.code,
        language: request.language,
        file_name: request.filename || null,
        context: null,
      }),
    });

    if (!response.ok) {
      const error: ApiError = await response.json();
      throw new Error(error.message || 'レビューの実行に失敗しました');
    }

    const data = await response.json();

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
    const response = await fetch(`${this.baseUrl}/api/v1/reviews/${reviewId}/feedback`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        score,
        comment: comment || '',
      }),
    });

    if (!response.ok) {
      const error: ApiError = await response.json();
      throw new Error(error.message || 'フィードバックの送信に失敗しました');
    }
  }
}

export const reviewApiClient = new ReviewApiClient();
