import { ChevronDown, ChevronUp } from 'lucide-react';
import { KnowledgeListItem } from '../types/knowledge';
import { getCategoryLabel } from '../api/knowledgeApi';
import { getPriorityStars } from '../types/knowledge';

interface KnowledgeTableProps {
  items: KnowledgeListItem[];
  onItemClick: (item: KnowledgeListItem) => void;
  sortBy: string;
  sortOrder: 'asc' | 'desc';
  onSort: (field: string) => void;
}

export default function KnowledgeTable({
  items,
  onItemClick,
  sortBy,
  sortOrder,
  onSort,
}: KnowledgeTableProps) {
  const formatDate = (dateString: string | null) => {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    });
  };

  const SortIcon = ({ field }: { field: string }) => {
    if (sortBy !== field) return null;
    return sortOrder === 'asc' ? (
      <ChevronUp size={16} className="inline" />
    ) : (
      <ChevronDown size={16} className="inline" />
    );
  };

  const getPriorityColor = (priority: number): string => {
    switch (priority) {
      case 3: return 'text-red-600';
      case 2: return 'text-yellow-600';
      case 1: return 'text-blue-600';
      default: return 'text-gray-600';
    }
  };

  const getCategoryColor = (category: string): string => {
    const colors: Record<string, string> = {
      error_handling: 'bg-red-100 text-red-700',
      testing: 'bg-green-100 text-green-700',
      performance: 'bg-blue-100 text-blue-700',
      security: 'bg-purple-100 text-purple-700',
      clean_code: 'bg-yellow-100 text-yellow-700',
      architecture: 'bg-indigo-100 text-indigo-700',
      other: 'bg-gray-100 text-gray-700',
    };
    return colors[category] || 'bg-gray-100 text-gray-700';
  };

  if (items.length === 0) {
    return (
      <div className="bg-white rounded-xl border border-gray-200 p-12 text-center">
        <p className="text-[#6B7280] text-sm">
          ナレッジが見つかりませんでした
        </p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-gray-50 border-b border-gray-200">
            <tr>
              <th
                className="px-6 py-4 text-left text-xs font-semibold text-[#6B7280] uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                onClick={() => onSort('created_at')}
              >
                タイトル <SortIcon field="created_at" />
              </th>
              <th className="px-6 py-4 text-left text-xs font-semibold text-[#6B7280] uppercase tracking-wider">
                カテゴリ
              </th>
              <th
                className="px-6 py-4 text-center text-xs font-semibold text-[#6B7280] uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                onClick={() => onSort('priority')}
              >
                重要度 <SortIcon field="priority" />
              </th>
              <th
                className="px-6 py-4 text-center text-xs font-semibold text-[#6B7280] uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                onClick={() => onSort('usage_count')}
              >
                使用回数 <SortIcon field="usage_count" />
              </th>
              <th className="px-6 py-4 text-left text-xs font-semibold text-[#6B7280] uppercase tracking-wider">
                最終使用
              </th>
              <th className="px-6 py-4 text-left text-xs font-semibold text-[#6B7280] uppercase tracking-wider">
                作成日
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {items.map((item) => (
              <tr
                key={item.id}
                onClick={() => onItemClick(item)}
                className="hover:bg-gray-50 cursor-pointer transition-colors"
              >
                {/* タイトル */}
                <td className="px-6 py-4">
                  <div className="text-sm font-medium text-[#111827] line-clamp-2">
                    {item.title}
                  </div>
                </td>

                {/* カテゴリ */}
                <td className="px-6 py-4">
                  <span
                    className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getCategoryColor(
                      item.category
                    )}`}
                  >
                    {getCategoryLabel(item.category)}
                  </span>
                </td>

                {/* 重要度 */}
                <td className="px-6 py-4 text-center">
                  <span className={`text-sm font-semibold ${getPriorityColor(item.priority)}`}>
                    {getPriorityStars(item.priority)}
                  </span>
                </td>

                {/* 使用回数 */}
                <td className="px-6 py-4 text-center">
                  <span className="text-sm text-[#6B7280]">
                    {item.usage_count}回
                  </span>
                </td>

                {/* 最終使用 */}
                <td className="px-6 py-4">
                  <span className="text-sm text-[#6B7280]">
                    {formatDate(item.last_used_at)}
                  </span>
                </td>

                {/* 作成日 */}
                <td className="px-6 py-4">
                  <span className="text-sm text-[#6B7280]">
                    {formatDate(item.created_at)}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
