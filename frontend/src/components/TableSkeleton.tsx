import React from 'react';

interface TableSkeletonProps {
  rows?: number;
}

/**
 * テーブルローディング用スケルトンコンポーネント
 */
const TableSkeleton: React.FC<TableSkeletonProps> = ({ rows = 10 }) => {
  return (
    <div className="overflow-x-auto bg-white border border-gray-200 rounded-xl">
      <table className="w-full text-left">
        <thead className="border-b border-gray-200">
          <tr>
            <th className="px-6 py-4 text-xs font-semibold text-[#6B7280] uppercase tracking-wider">
              日付
            </th>
            <th className="px-6 py-4 text-xs font-semibold text-[#6B7280] uppercase tracking-wider">
              言語
            </th>
            <th className="px-6 py-4 text-xs font-semibold text-[#6B7280] uppercase tracking-wider">
              ステータス
            </th>
            <th className="relative px-6 py-4"></th>
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-200">
          {Array.from({ length: rows }).map((_, index) => (
            <tr key={index}>
              <td className="px-6 py-4 whitespace-nowrap">
                <div className="h-4 bg-gray-200 rounded animate-pulse w-32"></div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <div className="h-4 bg-gray-200 rounded animate-pulse w-24"></div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap">
                <div className="h-6 bg-gray-200 rounded-full animate-pulse w-20"></div>
              </td>
              <td className="px-6 py-4 whitespace-nowrap text-right">
                <div className="h-4 bg-gray-200 rounded animate-pulse w-4 ml-auto"></div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default TableSkeleton;
