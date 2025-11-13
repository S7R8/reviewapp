import { useNavigate } from 'react-router-dom';
import { AlertCircle } from 'lucide-react';
import Sidebar, { SidebarToggle } from '../components/Sidebar';
import { useSidebar } from '../hooks/useSidebar';
import { useKnowledgeList } from '../hooks/useKnowledgeList';
import KnowledgeTable from '../components/KnowledgeTable';
import Pagination from '../components/Pagination';
import TableSkeleton from '../components/TableSkeleton';
import { KnowledgeListItem } from '../types/knowledge';
import { KnowledgeCategory, getCategoryLabel } from '../api/knowledgeApi';

const PRIORITY_OPTIONS = [
  { value: undefined, label: '全て' },
  { value: 3, label: '高 (★★★)' },
  { value: 2, label: '中 (★★☆)' },
  { value: 1, label: '低 (★☆☆)' },
];

const CATEGORY_OPTIONS: Array<{ value: KnowledgeCategory | undefined; label: string }> = [
  { value: undefined, label: '全カテゴリ' },
  { value: 'error_handling', label: getCategoryLabel('error_handling') },
  { value: 'testing', label: getCategoryLabel('testing') },
  { value: 'performance', label: getCategoryLabel('performance') },
  { value: 'security', label: getCategoryLabel('security') },
  { value: 'clean_code', label: getCategoryLabel('clean_code') },
  { value: 'architecture', label: getCategoryLabel('architecture') },
  { value: 'other', label: getCategoryLabel('other') },
];

export default function KnowledgePage() {
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
  } = useKnowledgeList();

  const handleCategoryChange = (category: KnowledgeCategory | undefined) => {
    updateFilter({ category });
  };

  const handlePriorityChange = (priority: number | undefined) => {
    updateFilter({ priority });
  };

  const handleItemClick = (item: KnowledgeListItem) => {
    // TODO: ナレッジ詳細ページへ遷移（Phase 2）
    console.log('ナレッジ詳細:', item);
  };

  return (
    <div className="flex h-screen overflow-hidden bg-[#f8f7f6]">
      {/* サイドバー */}
      <Sidebar currentPage="knowledge" isOpen={sidebarOpen} onToggle={toggleSidebar} />

      {/* メインコンテンツ */}
      <main className="flex-1 overflow-y-auto p-8 relative">
        {/* サイドバートグルボタン */}
        <div className="absolute top-8 left-8 z-10">
          <SidebarToggle isOpen={sidebarOpen} onToggle={toggleSidebar} />
        </div>

        <div className="max-w-7xl mx-auto">
          {/* ヘッダー */}
          <header className="mb-8 ml-16">
            <h1 className="text-[#111827] text-4xl font-black mb-2">
              📚 ナレッジ一覧
            </h1>
            <p className="text-[#6B7280] text-base">
              あなたの学びとルールを管理します
            </p>
          </header>

          {/* コンテンツ */}
          <div className="space-y-6">
            {/* フィルター */}
            <div className="flex items-center justify-between gap-4">
              <div className="flex items-center gap-3">
                {/* カテゴリフィルタ */}
                <select
                  value={filter.category || ''}
                  onChange={(e) =>
                    handleCategoryChange(
                      e.target.value ? (e.target.value as KnowledgeCategory) : undefined
                    )
                  }
                  className="h-10 px-4 rounded-lg border border-gray-300 bg-white text-[#111827] text-sm focus:border-[#FBBF24] focus:ring-[#FBBF24]"
                >
                  {CATEGORY_OPTIONS.map((option) => (
                    <option key={option.value || 'all'} value={option.value || ''}>
                      {option.label}
                    </option>
                  ))}
                </select>

                {/* 重要度フィルタ */}
                <select
                  value={filter.priority || ''}
                  onChange={(e) =>
                    handlePriorityChange(
                      e.target.value ? parseInt(e.target.value) : undefined
                    )
                  }
                  className="h-10 px-4 rounded-lg border border-gray-300 bg-white text-[#111827] text-sm focus:border-[#FBBF24] focus:ring-[#FBBF24]"
                >
                  {PRIORITY_OPTIONS.map((option) => (
                    <option key={option.value || 'all'} value={option.value || ''}>
                      {option.label}
                    </option>
                  ))}
                </select>
              </div>

              {/* 件数表示 */}
              <div className="text-sm text-[#6B7280]">
                全 <span className="font-semibold text-[#111827]">{totalItems}</span> 件
              </div>
            </div>

            {/* ローディング状態 */}
            {loading && <TableSkeleton rows={10} />}

            {/* エラー状態 */}
            {error && !loading && (
              <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                <div className="flex items-start gap-3">
                  <AlertCircle className="text-red-600 flex-shrink-0 mt-0.5" size={20} />
                  <div className="flex-1">
                    <p className="text-sm font-medium text-red-800 mb-1">
                      エラーが発生しました
                    </p>
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
                <KnowledgeTable
                  items={items}
                  onItemClick={handleItemClick}
                  sortBy={filter.sortBy}
                  sortOrder={filter.sortOrder}
                  onSort={changeSort}
                />

                {/* ページネーション */}
                {items.length > 0 && (
                  <Pagination
                    currentPage={filter.page}
                    totalPages={totalPages}
                    totalItems={totalItems}
                    pageSize={filter.pageSize}
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
