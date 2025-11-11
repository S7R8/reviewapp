import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Plus, AlertCircle } from 'lucide-react';
import Sidebar, { SidebarToggle } from '../components/Sidebar';
import { useSidebar } from '../hooks/useSidebar';
import { useReviewHistory } from '../hooks/useReviewHistory';
import ReviewHistoryTable from '../components/ReviewHistoryTable';
import SearchFilter from '../components/SearchFilter';
import Pagination from '../components/Pagination';
import TableSkeleton from '../components/TableSkeleton';
import { ReviewHistoryItem, ProgrammingLanguage } from '../types/review';

export default function HistoryPage() {
  const navigate = useNavigate();
  const { isOpen: sidebarOpen, toggle: toggleSidebar } = useSidebar();
  
  const {
    items,
    loading,
    error,
    filter,
    totalItems,
    totalPages,
    updateFilter,
    changePage,
    changeSort,
    refetch,
  } = useReviewHistory();

  const handleLanguageChange = (language: ProgrammingLanguage | undefined) => {
    updateFilter({ language });
  };

  const handleItemClick = (item: ReviewHistoryItem) => {
    // レビュー詳細ページへ遷移
    navigate(`/review/${item.id}`);
  };

  const handleNewReview = () => {
    navigate('/review');
  };

  return (
    <div className="flex h-screen overflow-hidden bg-[#f8f7f6]">
      {/* サイドバー */}
      <Sidebar currentPage="history" isOpen={sidebarOpen} onToggle={toggleSidebar} />

      {/* メインコンテンツ */}
      <main className="flex-1 overflow-y-auto p-8 relative">
        {/* サイドバートグルボタン */}
        <div className="absolute top-8 left-8 z-10">
          <SidebarToggle isOpen={sidebarOpen} onToggle={toggleSidebar} />
        </div>

        <div className="max-w-7xl mx-auto">
          {/* ヘッダー */}
          <header className="mb-8 ml-16">
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-[#111827] text-4xl font-black mb-2">
                  レビュー履歴
                </h1>
                <p className="text-[#6B7280] text-base">
                  過去のAIコードレビュー履歴を確認します
                </p>
              </div>
              <button
                onClick={handleNewReview}
                className="flex items-center justify-center gap-2 h-10 px-4 rounded-lg bg-[#FBBF24] text-[#111827] text-sm font-bold shadow-sm hover:bg-amber-400 transition-colors"
              >
                <Plus size={20} />
                <span>新規レビュー</span>
              </button>
            </div>
          </header>

          {/* コンテンツ */}
          <div className="space-y-6">
            {/* 検索・フィルター */}
            <div className="flex items-center justify-between">
              <SearchFilter
                onLanguageChange={handleLanguageChange}
              />
              
              {/* 件数表示 */}
              <div className="text-sm text-[#6B7280]">
                全 <span className="font-semibold text-[#111827]">{totalItems}</span> 件
              </div>
            </div>

            {/* ローディング状態（スケルトン表示） */}
            {loading && <TableSkeleton rows={10} />}

            {/* エラー状態 */}
            {error && !loading && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                <div className="flex items-start gap-3">
                  <AlertCircle className="text-red-600 flex-shrink-0 mt-0.5" size={20} />
                  <div className="flex-1">
                    <p className="text-sm font-medium text-red-800 mb-1">エラーが発生しました</p>
                    <p className="text-sm text-red-700">{error.message}</p>
                    <button
                      onClick={refetch}
                      className="mt-3 text-sm text-red-600 hover:text-red-800 underline"
                    >
                      再読み込み
                    </button>
                  </div>
                </div>
              </div>
            )}

            {/* テーブル */}
            {!loading && !error && (
              <>
                <ReviewHistoryTable
                  items={items}
                  onItemClick={handleItemClick}
                  sortBy={filter.sortBy}
                  sortOrder={filter.sortOrder}
                  onSort={changeSort}
                />

                {/* ページネーション */}
                {items.length > 0 && (
                  <Pagination
                    currentPage={filter.page || 1}
                    totalPages={totalPages}
                    totalItems={totalItems}
                    pageSize={filter.pageSize || 10}
                    onPageChange={changePage}
                  />
                )}
              </>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
