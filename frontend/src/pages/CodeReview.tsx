import { useState } from 'react';
import ReactMarkdown from 'react-markdown';
import Editor from '@monaco-editor/react';
import Sidebar, { SidebarToggle } from '../components/Sidebar';
import { ReviewResultSkeleton } from '../components/Skeleton';
import { FeedbackButtons } from '../components/FeedbackButtons';
import { useSidebar } from '../hooks/useSidebar';
import { useReviewStore } from '../store/reviewStore';
import {
  Bookmark,
  ChevronRight,
  AlertTriangle,
  CheckCircle2,
  X,
  AlertCircle
} from 'lucide-react';

const LANGUAGE_OPTIONS = [
  { value: 'python', label: 'Python' },
  { value: 'javascript', label: 'JavaScript' },
  { value: 'typescript', label: 'TypeScript' },
  { value: 'go', label: 'Go' },
  { value: 'java', label: 'Java' },
  { value: 'html', label: 'HTML' },
  { value: 'css', label: 'CSS' },
];

const DEFAULT_CODE = `def binary_search(arr, target):
    low = 0
    high = len(arr) - 1
    
    while low <= high:
        # This is a potential bug, can cause overflow
        mid = (low + high) // 2
        guess = arr[mid]
        
        if guess == target:
            return mid
        if guess > target:
            high = mid - 1
        else:
            low = mid + 1
    return None

# Test the function
my_list = [1, 3, 5, 7, 9]
print(binary_search(my_list, 3))  # Expected: 1
print(binary_search(my_list, -1)) # Expected: None`;

export default function CodeReview() {
  const [language, setLanguage] = useState('python');
  const [code, setCode] = useState(DEFAULT_CODE);
  const [showKnowledge, setShowKnowledge] = useState(false);
  const { isOpen: sidebarOpen, toggle: toggleSidebar } = useSidebar();

  const { currentReview, relatedKnowledge, isLoading, error, executeReview } = useReviewStore();

  const handleReview = async () => {
    await executeReview(code, language, 'code.py');
  };

  // パース失敗判定
  const isParseFailure = currentReview && 
    currentReview.goodPoints.length === 0 && 
    currentReview.improvements.length === 0;

  return (
    <div className="flex h-screen overflow-hidden bg-[#f8f7f6]">
      <Sidebar currentPage="review" isOpen={sidebarOpen} onToggle={toggleSidebar} />

      <main className="flex-1 flex flex-col overflow-hidden relative">
        <div className="absolute top-6 left-6 z-10">
          <SidebarToggle isOpen={sidebarOpen} onToggle={toggleSidebar} />
        </div>

        <div className="flex-1 flex flex-col p-6 lg:p-8 overflow-y-auto animate-fade-in">
          <div className="flex flex-col gap-8 h-full">
            <header className="ml-12">
              <h1 className="text-[#111827] text-4xl font-black mb-2">
                AIコードレビュー
              </h1>
              <p className="text-[#6B7280] text-base">
                コードを貼り付けて、AIによるレビューを開始します。
              </p>
            </header>

            <div className="flex-1 grid grid-cols-10 gap-6 h-full min-h-0">
              <div className="col-span-10 lg:col-span-4 flex flex-col gap-4 h-full">
                <div className="flex items-center">
                  <label className="flex flex-col" style={{ width: '200px' }}>
                    <p className="text-[#111827] text-sm font-medium mb-2">言語</p>
                    <select
                      value={language}
                      onChange={(e) => setLanguage(e.target.value)}
                      className="w-full h-12 px-4 rounded-lg border border-gray-300 bg-white text-[#111827] focus:border-[#FBBF24] focus:ring-[#FBBF24]"
                    >
                      {LANGUAGE_OPTIONS.map((option) => (
                        <option key={option.value} value={option.value}>
                          {option.label}
                        </option>
                      ))}
                    </select>
                  </label>
                </div>

                <div className="flex-1 min-h-[400px] rounded-xl overflow-hidden border border-gray-300">
                  <Editor
                    height="100%"
                    defaultLanguage={language}
                    language={language}
                    value={code}
                    onChange={(value) => setCode(value || '')}
                    theme="vs-dark"
                    options={{
                      minimap: { enabled: false },
                      fontSize: 14,
                      lineNumbers: 'on',
                      scrollBeyondLastLine: false,
                      automaticLayout: true,
                    }}
                  />
                </div>

                <button
                  onClick={handleReview}
                  disabled={isLoading || !code.trim()}
                  className="w-full h-12 px-4 bg-[#FBBF24] text-[#111827] rounded-lg font-bold text-base hover:bg-amber-400 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                >
                  {isLoading ? (
                    <>
                      <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-[#111827]" />
                      レビュー実行中...
                    </>
                  ) : (
                    'レビュー実行'
                  )}
                </button>
              </div>

              <div className="col-span-10 lg:col-span-6 flex flex-col bg-white rounded-xl border border-gray-200 h-full overflow-hidden">
                <div className="p-6 border-b border-gray-200 flex items-center justify-between">
                  <h2 className="text-lg font-bold text-[#111827]">レビュー結果</h2>
                  {currentReview && (
                    <button
                      onClick={() => setShowKnowledge(!showKnowledge)}
                      className="flex items-center gap-2 px-4 h-9 bg-[#F4C753]/20 text-[#111827] rounded-lg text-sm font-medium hover:bg-[#F4C753]/30 transition-colors"
                    >
                      関連ナレッジ {showKnowledge ? '非表示' : '表示'}
                    </button>
                  )}
                </div>

                <div className="flex-1 overflow-y-auto">
                  {isLoading ? (
                    <ReviewResultSkeleton />
                  ) : currentReview ? (
                    isParseFailure ? (
                      /* パース失敗時：生のマークダウンを表示 */
                      <div className="p-6">
                        <div className="mb-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                          <p className="text-sm text-yellow-800 flex items-center gap-2">
                            <AlertCircle size={16} />
                            レビュー結果のパースに失敗しました。生のレビュー結果を表示しています。
                          </p>
                        </div>
                        <div className="prose prose-sm max-w-none">
                          <ReactMarkdown
                            components={{
                              code(props) {
                                const { className, children, ...rest } = props;
                                const match = /language-(\w+)/.exec(className || '');
                                
                                return match ? (
                                  <pre className="bg-gray-900 text-gray-100 p-4 rounded-lg overflow-x-auto">
                                    <code className={className} {...rest}>
                                      {children}
                                    </code>
                                  </pre>
                                ) : (
                                  <code className="bg-gray-100 text-red-600 px-1 py-0.5 rounded text-sm" {...rest}>
                                    {children}
                                  </code>
                                );
                              },
                            }}
                          >
                            {currentReview.rawMarkdown}
                          </ReactMarkdown>
                        </div>
                      </div>
                    ) : (
                      /* 通常表示 */
                      <div className="space-y-6 p-6">
                        <p className="text-[#111827]">{currentReview.summary}</p>

                        {currentReview.goodPoints.length > 0 && (
                          <div className="p-4 rounded-lg bg-green-50 border border-green-200">
                            <h4 className="font-semibold text-green-700 mb-2 flex items-center gap-2">
                              <CheckCircle2 size={20} />
                              良い点
                            </h4>
                            <ul className="space-y-1">
                              {currentReview.goodPoints.map((point, i) => (
                                <li key={i} className="text-sm text-green-800">
                                  • {point}
                                </li>
                              ))}
                            </ul>
                          </div>
                        )}

                        {currentReview.improvements.length > 0 && (
                          <div className="p-4 rounded-lg bg-yellow-50 border border-yellow-200">
                            <h4 className="font-semibold text-amber-700 mb-3 flex items-center gap-2">
                              <AlertTriangle size={20} />
                              改善点
                            </h4>
                            <ol className="space-y-4">
                              {currentReview.improvements.map((improvement, i) => (
                                <li key={i} className="text-sm">
                                  <div className="flex items-start gap-2 mb-2">
                                    <strong className="text-amber-900">
                                      {i + 1}. {improvement.title}
                                    </strong>
                                    {improvement.severity === 'high' && (
                                      <span className="px-2 py-0.5 text-xs font-bold bg-red-100 text-red-700 rounded">
                                        HIGH
                                      </span>
                                    )}
                                    {improvement.severity === 'medium' && (
                                      <span className="px-2 py-0.5 text-xs font-bold bg-yellow-100 text-yellow-700 rounded">
                                        MEDIUM
                                      </span>
                                    )}
                                    {improvement.severity === 'low' && (
                                      <span className="px-2 py-0.5 text-xs font-bold bg-blue-100 text-blue-700 rounded">
                                        LOW
                                      </span>
                                    )}
                                  </div>
                                  <p className="text-amber-800 mt-1 whitespace-pre-line">
                                    {improvement.description}
                                  </p>
                                  {improvement.codeAfter && (
                                    <div className="mt-3">
                                      <p className="text-xs font-semibold text-amber-700 mb-1">
                                        改善例：
                                      </p>
                                      <pre className="p-3 bg-gray-900 text-gray-100 rounded border border-amber-200 text-xs overflow-x-auto">
                                        <code>{improvement.codeAfter}</code>
                                      </pre>
                                    </div>
                                  )}
                                </li>
                              ))}
                            </ol>
                          </div>
                        )}

                        {currentReview.references.length > 0 && (
                          <details className="group">
                            <summary className="flex items-center gap-2 cursor-pointer text-sm font-medium text-[#6B7280] hover:text-[#F4C753] transition-colors">
                              <ChevronRight className="transition-transform group-open:rotate-90" size={16} />
                              レビューの根拠を表示
                            </summary>
                            <div className="mt-3 ml-6 pl-4 border-l-2 border-[#F4C753]/30 text-xs text-[#6B7280]">
                              {currentReview.references.map((ref, i) => (
                                <div key={i} className="mb-2">
                                  <p className="font-medium text-[#111827]">{ref.source}</p>
                                  <p>{ref.description}</p>
                                </div>
                              ))}
                            </div>
                          </details>
                        )}
                      </div>
                    )
                  ) : (
                    <div className="flex items-center justify-center h-full text-[#6B7280] p-6">
                      <p className="text-sm">
                        コードを入力して「レビュー実行」をクリックしてください
                      </p>
                    </div>
                  )}
                </div>

                {currentReview && (
                  <div className="p-4 flex items-center justify-between border-t border-gray-200">
                    <FeedbackButtons />
                    
                    <button className="flex items-center gap-2 px-4 h-10 bg-[#F4C753]/20 text-[#111827] rounded-lg text-sm font-bold hover:bg-[#F4C753]/30 transition-colors">
                      <Bookmark size={16} />
                      ナレッジとして保存
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>

        {showKnowledge && relatedKnowledge.length > 0 && (
          <>
            <div
              className="fixed inset-0 bg-black/30 z-40"
              onClick={() => setShowKnowledge(false)}
            />
            <div className="fixed top-0 right-0 h-full w-full lg:w-[500px] bg-white border-l border-gray-200 shadow-2xl z-50 overflow-y-auto animate-slide-in-right">
              <div className="p-6">
                <div className="flex items-center justify-between mb-6">
                  <h3 className="text-xl font-bold text-[#111827]">関連ナレッジ</h3>
                  <button
                    onClick={() => setShowKnowledge(false)}
                    className="p-2 rounded-lg hover:bg-gray-100 transition-colors"
                  >
                    <X size={20} />
                  </button>
                </div>
                <div className="space-y-4">
                  {relatedKnowledge.map((knowledge) => (
                    <div
                      key={knowledge.id}
                      className="p-4 rounded-xl bg-gray-50 border border-gray-200 hover:border-[#F4C753]/50 hover:shadow-md transition-all cursor-pointer"
                    >
                      <h4 className="font-bold text-sm text-[#111827] mb-2">
                        {knowledge.title}
                      </h4>
                      <p className="text-xs text-[#6B7280] mb-3">
                        {knowledge.description}
                      </p>
                      <div className="flex gap-2 flex-wrap">
                        {knowledge.tags.map((tag) => (
                          <span
                            key={tag}
                            className="text-xs font-medium px-2 py-1 rounded-full bg-blue-100 text-blue-700"
                          >
                            {tag}
                          </span>
                        ))}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </>
        )}
      </main>
    </div>
  );
}
