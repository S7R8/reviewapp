import { useNavigate } from 'react-router-dom';
import Sidebar, { SidebarToggle } from '../components/Sidebar';
import { useSidebar } from '../hooks/useSidebar';
import { useAuthStore } from '../store/authStore';
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

  // 仮データ
  const stats = {
    totalReviews: 127,
    knowledgeCount: 89,
    consistency: 87,
    weeklyReviews: 23,
  };

  const recentReviews = [
    {
      id: '1',
      fileName: 'auth.go',
      date: '2024-01-20',
      status: 'completed',
      improvements: 3,
    },
    {
      id: '2',
      fileName: 'user_service.js',
      date: '2024-01-19',
      status: 'completed',
      improvements: 5,
    },
    {
      id: '3',
      fileName: 'api_handler.py',
      date: '2024-01-18',
      status: 'completed',
      improvements: 2,
    },
  ];

  const skillRadar = {
    errorHandling: 95,
    testing: 80,
    performance: 65,
    security: 50,
    documentation: 70,
  };

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
                <div className="space-y-4">
                  {recentReviews.map((review) => (
                    <div
                      key={review.id}
                      className="flex items-center justify-between p-4 rounded-lg border border-gray-200 hover:border-[#F4C753] transition-colors cursor-pointer"
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
                        <CheckCircle2 className="text-[#10B981]" size={20} />
                      </div>
                    </div>
                  ))}
                </div>
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

          {/* 開発中メッセージ */}
          <div className="mt-8 p-4 bg-green-50 border border-green-200 rounded-lg">
            <p className="text-sm text-green-800">
              ✅ Dashboard画面（仮実装）が動作しています<br />
              統計データは仮のデータです。バックエンド実装後に実際のデータを表示します。
            </p>
          </div>
        </div>
      </main>
    </div>
  );
}
