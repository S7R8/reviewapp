import React from 'react';
import { ChevronRight, ArrowUp, ArrowDown } from 'lucide-react';
import { ReviewHistoryItem } from '../types/review';
import StatusBadge from './StatusBadge';

interface ReviewHistoryTableProps {
  items: ReviewHistoryItem[];
  onItemClick: (item: ReviewHistoryItem) => void;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
  onSort?: (field: string) => void;
}

const ReviewHistoryTable: React.FC<ReviewHistoryTableProps> = ({
  items,
  onItemClick,
  sortBy,
  sortOrder,
  onSort,
}) => {
  const formatDate = (dateString: string): string => {
    const date = new Date(dateString);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${year}/${month}/${day} ${hours}:${minutes}`;
  };

  const handleSort = (field: string) => {
    if (onSort) {
      onSort(field);
    }
  };

  const renderSortIcon = (field: string) => {
    if (sortBy !== field) return null;
    return sortOrder === 'asc' ? (
      <ArrowUp className="w-4 h-4" />
    ) : (
      <ArrowDown className="w-4 h-4" />
    );
  };

  if (items.length === 0) {
    return (
      <div className="bg-white border border-gray-200 rounded-xl p-16 text-center">
        <div className="flex flex-col items-center gap-4">
          <div className="w-16 h-16 rounded-full bg-gray-100 flex items-center justify-center">
            <ChevronRight className="w-8 h-8 text-gray-400" />
          </div>
          <div>
            <h3 className="text-lg font-semibold text-[#111827] mb-1">
              まだレビュー履歴がありません
            </h3>
            <p className="text-sm text-[#6B7280]">
              最初のレビューを実行してみましょう
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto bg-white border border-gray-200 rounded-xl">
      <table className="w-full text-left">
        <thead className="border-b border-gray-200 bg-gray-50">
          <tr>
            <th
              className="px-6 py-4 text-xs font-semibold text-[#6B7280] uppercase tracking-wider cursor-pointer hover:text-[#111827] transition-colors"
              scope="col"
              onClick={() => handleSort('createdAt')}
            >
              <div className="flex items-center gap-2">
                <span>日付</span>
                {renderSortIcon('createdAt')}
              </div>
            </th>
            <th
              className="px-6 py-4 text-xs font-semibold text-[#6B7280] uppercase tracking-wider cursor-pointer hover:text-[#111827] transition-colors"
              scope="col"
              onClick={() => handleSort('language')}
            >
              <div className="flex items-center gap-2">
                <span>言語</span>
                {renderSortIcon('language')}
              </div>
            </th>
            <th
              className="px-6 py-4 text-xs font-semibold text-[#6B7280] uppercase tracking-wider cursor-pointer hover:text-[#111827] transition-colors"
              scope="col"
              onClick={() => handleSort('status')}
            >
              <div className="flex items-center gap-2">
                <span>ステータス</span>
                {renderSortIcon('status')}
              </div>
            </th>
            <th className="relative px-6 py-4" scope="col">
              <span className="sr-only">Actions</span>
            </th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-200">
          {items.map((item) => (
            <tr
              key={item.id}
              onClick={() => onItemClick(item)}
              className="hover:bg-gray-50 cursor-pointer transition-colors group"
            >
              <td className="px-6 py-4 whitespace-nowrap text-sm text-[#111827]">
                {formatDate(item.createdAt)}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm text-[#111827] font-medium">
                {item.language}
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-sm">
                <StatusBadge status={item.status} />
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <ChevronRight className="w-5 h-5 text-gray-400 group-hover:text-[#F4C753] transition-colors inline-block" />
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default ReviewHistoryTable;
