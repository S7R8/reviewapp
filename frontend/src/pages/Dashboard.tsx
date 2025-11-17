import { useNavigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import Sidebar, { SidebarToggle } from '../components/Sidebar';
import { useSidebar } from '../hooks/useSidebar';
import { useAuthStore } from '../store/authStore';
import { DashboardSkeleton } from '../components/DashboardSkeleton';
import { dashboardApiClient, DashboardStatsResponse } from '../api/dashboardApi';
import {
  Code,
  Book,
  TrendingUp,
  CheckCircle2,
  AlertCircle,
  Zap
} from 'lucide-react';

export default function Dashboard() {
  const navigate = useNavigate();
  const { user } = useAuthStore();
  const { isOpen: sidebarOpen, toggle: toggleSidebar } = useSidebar();

  const [loading, setLoading] = useState(true);
  const [data, setData] = useState<DashboardStatsResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        setLoading(true);
        console.log('[Dashboard] Fetching stats...');
        const response = await dashboardApiClient.getStats();
        console.log('[Dashboard] Stats received:', response);
        setData(response);
        setError(null);
      } catch (err) {
        console.error('Failed to fetch dashboard stats:', err);
        setError('統計情報の取得に失敗しました');
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  // データがない場合のデフォルト値
  const stats = data ? {
    totalReviews: data.stats.total_reviews,
    knowledgeCount: data.stats.knowledge_count,
    consistency: data.stats.consistency_score,
    weeklyReviews: data.stats.weekly_reviews,
  } : {
    totalReviews: 0,
    knowledgeCount: 0,
    consistency: 0,
    weeklyReviews: 0,
  };

  const recentReviews = data ? data.recent_reviews.map(review => {
    console.log('[Dashboard] Processing review:', review);

    // 言語名を日本語表示用に変換
    const languageMap: { [key: string]: string } = {
      'typescript': 'TypeScript',
      'javascript': 'JavaScript',
      'python': 'Python',
      'go': 'Go',
      'java': 'Java',
      'rust': 'Rust',
      'ruby': 'Ruby',
      'php': 'PHP',
    };

    const displayName = languageMap[review.language.toLowerCase()] || review.language;

    return {
      id: review.id,
      fileName: `${displayName} レビュー`,
      date: new Date(review.created_at).toLocaleDateString('ja-JP'),
      status: review.status,
      improvements: review.improvements_count,
    };
  }) : [];

  console.log('[Dashboard] Recent reviews:', recentReviews);

  const skillRadar = data ? {
    errorHandling: data.skill_analysis.error_handling,
    testing: data.skill_analysis.testing,
    performance: data.skill_analysis.performance,
    security: data.skill_analysis.security,
    documentation: data.skill_analysis.clean_code,
  } : {
    errorHandling: 0,
    testing: 0,
    performance: 0,
    security: 0,
    documentation: 0,
  };

  console.log('[Dashboard] Skill analysis data:', data?.skill_analysis);
  console.log('[Dashboard] Skill radar:', skillRadar);

  return (
    <div className="flex h-screen overflow-hidden bg-[#f8f7f6]">
      {/* サイドバー */}
      <Sidebar currentPage="dashboard" isOpen={sidebarOpen} onToggle={toggleSidebar} />

      {/* メインコンテンツ */}
      <main className="flex-1 overflow-y-auto p-8 relative">
        {/* サイドバートグルボタン */}
        <div className="absolute top-8 left-8 z-10">
          <SidebarToggle isOpen={sidebarOpen} onToggle={toggleSidebar} />
        </div>

        <div className="max-w-6xl mx-auto">
          {/* ヘッダー */}
          <header className="mb-8 ml-16">
            <h1 className="text-[#111827] text-4xl font-black mb-2">
              ダッシュボード
            </h1>
            <p className="text-[#6B7280] text-base">
              あなたのコーディング習慣とAIクローンの成長を確認しましょう
            </p>
          </header>

          {/* エラー表示 */}
          {error && (
            <div className="mb-8 p-4 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-sm text-red-800">{error}</p>
            </div>
          )}

          {/* コンテンツ */}
          {loading ? (
            <DashboardSkeleton />
          ) : (
            <>
              {/* 統計カード */}
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
                {/* レビュー回数 */}
                <div className="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
                  <div className="flex items-center justify-between mb-2">
                    <p className="text-[#6B7280] text-sm font-medium">レビュー回数</p>
                    <Code className="text-[#F4C753]" size={20} />
                  </div>
                  <p className="text-[#111827] text-3xl font-bold">{stats.totalReviews}</p>
                  <p className="text-[#10B981] text-sm font-medium mt-1">
                    今週 +{stats.weeklyReviews}
                  </p>
                </div>

                {/* ナレッジ数 */}
                <div className="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
                  <div className="flex items-center justify-between mb-2">
                    <p className="text-[#6B7280] text-sm font-medium">ナレッジ数</p>
                    <Book className="text-[#F4C753]" size={20} />
                  </div>
                  <p className="text-[#111827] text-3xl font-bold">{stats.knowledgeCount}</p>
                  <p className="text-[#6B7280] text-sm mt-1">蓄積済み</p>
                </div>

                {/* 一貫性スコア */}
                <div className="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
                  <div className="flex items-center justify-between mb-2">
                    <p className="text-[#6B7280] text-sm font-medium">一貫性スコア</p>
                    <CheckCircle2 className="text-[#F4C753]" size={20} />
                  </div>
                  <p className="text-[#111827] text-3xl font-bold">{stats.consistency}%</p>
                  <p className="text-[#10B981] text-sm font-medium mt-1">+5% 先月比</p>
                </div>

                {/* 成長率 */}
                <div className="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
                  <div className="flex items-center justify-between mb-2">
                    <p className="text-[#6B7280] text-sm font-medium">成長率</p>
                    <TrendingUp className="text-[#F4C753]" size={20} />
                  </div>
                  <p className="text-[#111827] text-3xl font-bold">+23%</p>
                  <p className="text-[#6B7280] text-sm mt-1">先月比</p>
                </div>
              </div>

              <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* 最近のレビュー */}
                <div className="lg:col-span-2">
                  <div className="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
                    <h2 className="text-[#111827] text-xl font-bold mb-4">
                      最近のレビュー
                    </h2>
                    {recentReviews.length === 0 ? (
                      <div className="text-center py-8">
                        <Code className="mx-auto text-gray-400 mb-2" size={48} />
                        <p className="text-gray-500 text-sm">レビュー履歴がありません</p>
                        <button
                          onClick={() => navigate('/review')}
                          className="mt-4 px-4 py-2 bg-[#F4C753] text-white rounded-lg hover:bg-[#E5B84C] transition-colors"
                        >
                          コードレビューを開始
                        </button>
                      </div>
                    ) : (
                      <div className="space-y-4">
                        {recentReviews.map((review) => (
                          <div
                            key={review.id}
                            className="flex items-center justify-between p-4 rounded-lg border border-gray-200 hover:border-[#F4C753] transition-colors cursor-pointer"
                            onClick={() => navigate(`/review/${review.id}`)}
                          >
                            <div className="flex items-center gap-4">
                              <div className="w-10 h-10 rounded-lg bg-[#F4C753]/20 flex items-center justify-center">
                                <Code className="text-[#F4C753]" size={20} />
                              </div>
                              <div>
                                <p className="text-[#111827] font-medium">
                                  {review.fileName}
                                </p>
                                <p className="text-[#6B7280] text-sm">{review.date}</p>
                              </div>
                            </div>
                            <div className="flex items-center gap-3">
                              <span className="flex items-center gap-1 text-[#6B7280] text-sm">
                                <AlertCircle size={16} />
                                {review.improvements}件の改善
                              </span>
                              {review.status === 'success' ? (
                                <CheckCircle2 className="text-[#10B981]" size={20} />
                              ) : review.status === 'warning' ? (
                                <AlertCircle className="text-[#F4C753]" size={20} />
                              ) : (
                                <AlertCircle className="text-[#EF4444]" size={20} />
                              )}
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>

                {/* スキル分析 */}
                <div>
                  <div className="bg-white rounded-xl p-6 border border-gray-200 shadow-sm">
                    <h2 className="text-[#111827] text-xl font-bold mb-4">
                      あなたの重視ポイント
                    </h2>
                    <div className="space-y-4">
                      {/* エラーハンドリング */}
                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-[#111827] text-sm font-medium">
                            エラーハンドリング
                          </span>
                          <span className="text-[#6B7280] text-sm">
                            {skillRadar.errorHandling}%
                          </span>
                        </div>
                        <div className="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
                          <div
                            className="h-full bg-[#10B981] rounded-full"
                            style={{ width: `${skillRadar.errorHandling}%` }}
                          ></div>
                        </div>
                      </div>

                      {/* テスト */}
                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-[#111827] text-sm font-medium">
                            テスト
                          </span>
                          <span className="text-[#6B7280] text-sm">
                            {skillRadar.testing}%
                          </span>
                        </div>
                        <div className="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
                          <div
                            className="h-full bg-[#3B82F6] rounded-full"
                            style={{ width: `${skillRadar.testing}%` }}
                          ></div>
                        </div>
                      </div>

                      {/* ドキュメント */}
                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-[#111827] text-sm font-medium">
                            ドキュメント
                          </span>
                          <span className="text-[#6B7280] text-sm">
                            {skillRadar.documentation}%
                          </span>
                        </div>
                        <div className="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
                          <div
                            className="h-full bg-[#F4C753] rounded-full"
                            style={{ width: `${skillRadar.documentation}%` }}
                          ></div>
                        </div>
                      </div>

                      {/* パフォーマンス */}
                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-[#111827] text-sm font-medium">
                            パフォーマンス
                          </span>
                          <span className="text-[#6B7280] text-sm">
                            {skillRadar.performance}%
                          </span>
                        </div>
                        <div className="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
                          <div
                            className="h-full bg-[#8B5CF6] rounded-full"
                            style={{ width: `${skillRadar.performance}%` }}
                          ></div>
                        </div>
                      </div>

                      {/* セキュリティ */}
                      <div>
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-[#111827] text-sm font-medium">
                            セキュリティ
                          </span>
                          <span className="text-[#6B7280] text-sm">
                            {skillRadar.security}%
                          </span>
                        </div>
                        <div className="w-full h-2 bg-gray-200 rounded-full overflow-hidden">
                          <div
                            className="h-full bg-[#EF4444] rounded-full"
                            style={{ width: `${skillRadar.security}%` }}
                          ></div>
                        </div>
                      </div>
                    </div>

                    <div className="mt-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
                      <div className="flex items-start gap-2">
                        <Zap className="text-blue-600 flex-shrink-0" size={20} />
                        <div>
                          <p className="text-blue-900 text-sm font-medium mb-1">
                            改善のヒント
                          </p>
                          <p className="text-blue-700 text-xs">
                            セキュリティのナレッジを増やすと、さらに精度が向上します
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </>
          )}
        </div>
      </main>
    </div>
  );
}
