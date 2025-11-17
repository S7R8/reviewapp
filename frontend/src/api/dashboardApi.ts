import { apiGet } from './client';

// ダッシュボード統計情報の型定義
export interface DashboardStats {
  total_reviews: number;
  knowledge_count: number;
  consistency_score: number;
  weekly_reviews: number;
}

export interface RecentReviewItem {
  id: string;
  file_name: string;
  language: string;
  created_at: string;
  improvements_count: number;
  status: 'success' | 'warning' | 'error';
}

export interface SkillAnalysis {
  error_handling: number;
  testing: number;
  performance: number;
  security: number;
  clean_code: number;
  architecture: number;
  other: number;
}

export interface DashboardStatsResponse {
  stats: DashboardStats;
  recent_reviews: RecentReviewItem[];
  skill_analysis: SkillAnalysis;
}

class DashboardApiClient {
  /**
   * ダッシュボード統計情報を取得
   */
  async getStats(): Promise<DashboardStatsResponse> {
    return apiGet<DashboardStatsResponse>('/api/v1/dashboard/stats');
  }
}

export const dashboardApiClient = new DashboardApiClient();
