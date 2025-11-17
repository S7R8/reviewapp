import { X } from 'lucide-react';
import { useState } from 'react';
import { KnowledgeCategory, KnowledgeCreateRequest, getCategoryLabel } from '../api/knowledgeApi';

interface KnowledgeCreateModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreate: (request: KnowledgeCreateRequest) => Promise<void>;
}

const CATEGORY_OPTIONS: Array<{ value: KnowledgeCategory; label: string }> = [
  { value: 'error_handling', label: getCategoryLabel('error_handling') },
  { value: 'testing', label: getCategoryLabel('testing') },
  { value: 'performance', label: getCategoryLabel('performance') },
  { value: 'security', label: getCategoryLabel('security') },
  { value: 'clean_code', label: getCategoryLabel('clean_code') },
  { value: 'architecture', label: getCategoryLabel('architecture') },
  { value: 'other', label: getCategoryLabel('other') },
];

const PRIORITY_OPTIONS = [
  { value: 1, label: '低 (★☆☆)', description: '参考程度の情報' },
  { value: 2, label: '中 (★★☆)', description: '重要な情報' },
  { value: 3, label: '高 (★★★)', description: '必須の情報' },
];

export default function KnowledgeCreateModal({
  isOpen,
  onClose,
  onCreate,
}: KnowledgeCreateModalProps) {
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [category, setCategory] = useState<KnowledgeCategory>('other');
  const [priority, setPriority] = useState(2);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  if (!isOpen) return null;

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!title.trim()) {
      newErrors.title = 'タイトルは必須です';
    } else if (title.length > 200) {
      newErrors.title = 'タイトルは200文字以内にしてください';
    }

    if (!content.trim()) {
      newErrors.content = '内容は必須です';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) {
      return;
    }

    setIsSubmitting(true);
    try {
      await onCreate({
        title: title.trim(),
        content: content.trim(),
        category,
        priority,
      });

      // 成功したらフォームをリセットして閉じる
      setTitle('');
      setContent('');
      setCategory('other');
      setPriority(2);
      setErrors({});
      onClose();
    } catch (error) {
      console.error('ナレッジの作成に失敗しました:', error);
      setErrors({ submit: 'ナレッジの作成に失敗しました。もう一度お試しください。' });
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleClose = () => {
    setTitle('');
    setContent('');
    setCategory('other');
    setPriority(2);
    setErrors({});
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
          className="bg-white rounded-2xl shadow-2xl max-w-2xl w-full max-h-[90vh] overflow-hidden pointer-events-auto"
          onClick={(e) => e.stopPropagation()}
        >
          {/* ヘッダー */}
          <div className="flex items-center justify-between p-6 border-b border-gray-200">
            <h2 className="text-2xl font-bold text-[#111827]">
              ナレッジを追加
            </h2>
            <button
              onClick={handleClose}
              className="text-gray-400 hover:text-gray-600 transition-colors"
            >
              <X size={24} />
            </button>
          </div>

          {/* フォーム */}
          <form onSubmit={handleSubmit} className="p-6 overflow-y-auto max-h-[calc(90vh-140px)]">
            <div className="space-y-6">
              {/* タイトル */}
              <div>
                <label
                  htmlFor="title"
                  className="block text-sm font-semibold text-[#111827] mb-2"
                >
                  タイトル <span className="text-red-500">*</span>
                </label>
                <input
                  id="title"
                  type="text"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  placeholder="例: エラーハンドリングの原則"
                  className={`w-full px-4 py-2 rounded-lg border ${
                    errors.title ? 'border-red-500' : 'border-gray-300'
                  } focus:border-[#FBBF24] focus:ring-1 focus:ring-[#FBBF24] outline-none transition-colors`}
                  maxLength={200}
                />
                {errors.title && (
                  <p className="mt-1 text-sm text-red-600">{errors.title}</p>
                )}
                <p className="mt-1 text-xs text-gray-500">
                  {title.length}/200文字
                </p>
              </div>

              {/* 内容 */}
              <div>
                <label
                  htmlFor="content"
                  className="block text-sm font-semibold text-[#111827] mb-2"
                >
                  内容 <span className="text-red-500">*</span>
                </label>
                <textarea
                  id="content"
                  value={content}
                  onChange={(e) => setContent(e.target.value)}
                  placeholder="例: エラーは必ずログに出力し、ユーザー向けメッセージと開発者向け詳細を分ける。contextを使ってエラーチェーンを保持する。"
                  rows={6}
                  className={`w-full px-4 py-2 rounded-lg border ${
                    errors.content ? 'border-red-500' : 'border-gray-300'
                  } focus:border-[#FBBF24] focus:ring-1 focus:ring-[#FBBF24] outline-none transition-colors resize-vertical`}
                />
                {errors.content && (
                  <p className="mt-1 text-sm text-red-600">{errors.content}</p>
                )}
              </div>

              {/* カテゴリと重要度（横並び） */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* カテゴリ */}
                <div>
                  <label
                    htmlFor="category"
                    className="block text-sm font-semibold text-[#111827] mb-2"
                  >
                    カテゴリ <span className="text-red-500">*</span>
                  </label>
                  <select
                    id="category"
                    value={category}
                    onChange={(e) => setCategory(e.target.value as KnowledgeCategory)}
                    className="w-full px-4 py-2 rounded-lg border border-gray-300 focus:border-[#FBBF24] focus:ring-1 focus:ring-[#FBBF24] outline-none transition-colors"
                  >
                    {CATEGORY_OPTIONS.map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                  </select>
                </div>

                {/* 重要度 */}
                <div>
                  <label
                    htmlFor="priority"
                    className="block text-sm font-semibold text-[#111827] mb-2"
                  >
                    重要度 <span className="text-red-500">*</span>
                  </label>
                  <select
                    id="priority"
                    value={priority}
                    onChange={(e) => setPriority(Number(e.target.value))}
                    className="w-full px-4 py-2 rounded-lg border border-gray-300 focus:border-[#FBBF24] focus:ring-1 focus:ring-[#FBBF24] outline-none transition-colors"
                  >
                    {PRIORITY_OPTIONS.map((option) => (
                      <option key={option.value} value={option.value}>
                        {option.label}
                      </option>
                    ))}
                  </select>
                  <p className="mt-1 text-xs text-gray-500">
                    {PRIORITY_OPTIONS.find((o) => o.value === priority)?.description}
                  </p>
                </div>
              </div>

              {/* エラーメッセージ */}
              {errors.submit && (
                <div className="bg-red-50 border border-red-200 rounded-lg p-3">
                  <p className="text-sm text-red-800">{errors.submit}</p>
                </div>
              )}
            </div>
          </form>

          {/* フッター */}
          <div className="flex items-center justify-end gap-3 p-6 border-t border-gray-200 bg-gray-50">
            <button
              type="button"
              onClick={handleClose}
              disabled={isSubmitting}
              className="px-6 py-2 text-sm font-medium text-gray-700 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors disabled:opacity-50"
            >
              キャンセル
            </button>
            <button
              onClick={handleSubmit}
              disabled={isSubmitting}
              className="px-6 py-2 text-sm font-medium text-white bg-[#FBBF24] hover:bg-[#F59E0B] rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isSubmitting ? '作成中...' : '作成する'}
            </button>
          </div>
        </div>
      </div>
    </>
  );
}
