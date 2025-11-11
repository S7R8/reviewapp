import React from 'react';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  totalItems: number;
  pageSize: number;
  onPageChange: (page: number) => void;
}

/**
 * ページネーションコンポーネント
 * レビュー履歴のページ切り替えを管理
 */
const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  totalPages,
  totalItems,
  pageSize,
  onPageChange,
}) => {
  const startItem = (currentPage - 1) * pageSize + 1;
  const endItem = Math.min(currentPage * pageSize, totalItems);

  const handlePrevious = () => {
    if (currentPage > 1) {
      onPageChange(currentPage - 1);
    }
  };

  const handleNext = () => {
    if (currentPage < totalPages) {
      onPageChange(currentPage + 1);
    }
  };

  return (
    <nav
      aria-label="Pagination"
      className="flex items-center justify-between border-t border-gray-200 px-4 py-3 sm:px-6"
    >
      <div className="hidden sm:block">
        <p className="text-sm text-[#6B7280]">
          <span className="font-medium">{totalItems}</span>件中
          <span className="font-medium"> {startItem}</span> -{' '}
          <span className="font-medium">{endItem}</span>件を表示
        </p>
      </div>
      <div className="flex flex-1 justify-between sm:justify-end gap-2">
        <button
          onClick={handlePrevious}
          disabled={currentPage <= 1}
          className="relative inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-[#111827] hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          前へ
        </button>
        <button
          onClick={handleNext}
          disabled={currentPage >= totalPages}
          className="relative inline-flex items-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-[#111827] hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          次へ
        </button>
      </div>
    </nav>
  );
};

export default Pagination;
