import { X, Trash2, Calendar, Tag } from 'lucide-react';
import { KnowledgeListItem } from '../types/knowledge';
import { getCategoryLabel } from '../api/knowledgeApi';
import { getPriorityStars } from '../types/knowledge';
import { useState } from 'react';

interface KnowledgeDetailModalProps {
  knowledge: KnowledgeListItem;
  isOpen: boolean;
  onClose: () => void;
  onDelete: (id: string) => Promise<void>;
}

export default function KnowledgeDetailModal({
  knowledge,
  isOpen,
  onClose,
  onDelete,
}: KnowledgeDetailModalProps) {
  const [isDeleting, setIsDeleting] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  if (!isOpen) return null;

  const formatDate = (dateString: string | null) => {
    if (!dateString) return '-';
    const date = new Date(dateString);
    return date.toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getPriorityColor = (priority: number): string => {
    switch (priority) {
      case 3: return 'text-red-600 bg-red-50';
      case 2: return 'text-yellow-600 bg-yellow-50';
      case 1: return 'text-blue-600 bg-blue-50';
      default: return 'text-gray-600 bg-gray-50';
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

  const handleDelete = async () => {
    if (!showDeleteConfirm) {
      setShowDeleteConfirm(true);
      return;
    }

    setIsDeleting(true);
    try {
      await onDelete(knowledge.id);
      onClose();
    } catch (error) {
      console.error('削除に失敗しました:', error);
      alert('削除に失敗しました。もう一度お試しください。');
    } finally {
      setIsDeleting(false);
      setShowDeleteConfirm(false);
    }
  };

  const handleClose = () => {
    setShowDeleteConfirm(false);
    onClose();
  };

  return (
    <>
      {/* オーバーレイ */}
      <div
        className="fixed inset-0 bg-black bg-opacity-50 z-40 transition-opacity"
        onClick={handleClose}
      />

      {/* モーダル */}
      <div className="fixed inset-0 z-50 flex items-center justify-center p-4 pointer-events-none">
        <div
          className="bg-white rounded-2xl shadow-2xl max-w-3xl w-full max-h-[90vh] overflow-hidden pointer-events-auto"
          onClick={(e) => e.stopPropagation()}
        >
          {/* ヘッダー */}
          <div className="flex items-start justify-between p-6 border-b border-gray-200">
            <div className="flex-1">
              <h2 className="text-2xl font-bold text-[#111827] mb-2">
                {knowledge.title}
              </h2>
              <div className="flex items-center gap-3 flex-wrap">
                {/* カテゴリ */}
                <span
                  className={`inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm font-medium ${getCategoryColor(
                    knowledge.category
                  )}`}
                >
                  <Tag size={14} />
                  {getCategoryLabel(knowledge.category)}
                </span>

                {/* 重要度 */}
                <span
                  className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-semibold ${getPriorityColor(
                    knowledge.priority
                  )}`}
                >
                  {getPriorityStars(knowledge.priority)}
                </span>

                {/* 使用回数 */}
                <span className="text-sm text-[#6B7280]">
                  使用回数: <strong>{knowledge.usage_count}回</strong>
                </span>
              </div>
            </div>

            {/* 閉じるボタン */}
            <button
              onClick={handleClose}
              className="text-gray-400 hover:text-gray-600 transition-colors ml-4"
            >
              <X size={24} />
            </button>
          </div>

          {/* コンテンツ */}
          <div className="p-6 overflow-y-auto max-h-[calc(90vh-200px)]">
            {/* 本文 */}
            <div className="mb-6">
              <h3 className="text-sm font-semibold text-[#6B7280] mb-2 uppercase tracking-wider">
                内容
              </h3>
              <div className="bg-gray-50 rounded-lg p-4 text-[#111827] whitespace-pre-wrap leading-relaxed">
                {knowledge.content}
              </div>
            </div>

            {/* メタ情報 */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
              <div className="flex items-center gap-2 text-[#6B7280]">
                <Calendar size={16} />
                <span>作成日: {formatDate(knowledge.created_at)}</span>
              </div>
              {knowledge.last_used_at && (
                <div className="flex items-center gap-2 text-[#6B7280]">
                  <Calendar size={16} />
                  <span>最終使用: {formatDate(knowledge.last_used_at)}</span>
                </div>
              )}
            </div>
          </div>

          {/* フッター（アクション） */}
          <div className="flex items-center justify-between p-6 border-t border-gray-200 bg-gray-50">
            {!showDeleteConfirm ? (
              <>
                <button
                  onClick={handleDelete}
                  disabled={isDeleting}
                  className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-red-600 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors disabled:opacity-50"
                >
                  <Trash2 size={16} />
                  削除
                </button>

                <button
                  onClick={handleClose}
                  className="px-6 py-2 text-sm font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
                >
                  閉じる
                </button>
              </>
            ) : (
              <>
                <div className="text-sm text-red-600 font-medium">
                  本当に削除しますか？この操作は取り消せません。
                </div>
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => setShowDeleteConfirm(false)}
                    disabled={isDeleting}
                    className="px-4 py-2 text-sm font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors disabled:opacity-50"
                  >
                    キャンセル
                  </button>
                  <button
                    onClick={handleDelete}
                    disabled={isDeleting}
                    className="px-4 py-2 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {isDeleting ? '削除中...' : '削除する'}
                  </button>
                </div>
              </>
            )}
          </div>
        </div>
      </div>
    </>
  );
}
