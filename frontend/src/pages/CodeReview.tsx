import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';
import Editor from '@monaco-editor/react';
import Sidebar, { SidebarToggle } from '../components/Sidebar';
import { ReviewResultSkeleton } from '../components/Skeleton';
import { FeedbackButtons } from '../components/FeedbackButtons';
import { Toast } from '../components/Toast';
import { useSidebar } from '../hooks/useSidebar';
import { useReviewStore } from '../store/reviewStore';
import { detectLanguage } from '../utils/languageDetector';
import { createKnowledgeFromReview } from '../utils/knowledgeHelper';
import { knowledgeApiClient, getCategoryLabel } from '../api/knowledgeApi';
import { reviewApiClient } from '../api/reviewApi';
import {
  Bookmark,
  ChevronRight,
  AlertTriangle,
  CheckCircle2,
  X,
  AlertCircle,
  Loader,
  ArrowLeft
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

type ToastType = 'success' | 'error';

export default function CodeReview() {
  const { id: reviewId } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { isOpen: sidebarOpen, toggle: toggleSidebar } = useSidebar();
  
  // 閲覧モード判定
  const isViewMode = !!reviewId;
  
  const { 
    currentReview, 
    relatedKnowledge, 
    isLoading, 
    error, 
    executeReview,
    currentCode,
    currentLanguage,
    setCode,
    reset,
    loadReviewById
  } = useReviewStore();
  
  const [code, setCodeLocal] = useState(currentCode);
  const [language, setLanguageLocal] = useState(currentLanguage);
  const [showKnowledge, setShowKnowledge] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  
  // トースト通知
  const [toast, setToast] = useState<{ type: ToastType; message: string } | null>(null);

  // 閲覧モード: レビュー履歴をロード
  useEffect(() => {
    if (reviewId && loadReviewById) {
      loadReviewById(reviewId);
    }
  }, [reviewId, loadReviewById]);

  // 閲覧モード・通常モード両方でグローバル状態を復元
  useEffect(() => {
    // currentCodeが空でも、閲覧モードの場合は反映（初回ロード後に更新される）
    setCodeLocal(currentCode);
    setLanguageLocal(currentLanguage);
  }, [currentCode, currentLanguage]);

  // 閲覧モード → 通常モード への遷移を検知してリセット
  useEffect(() => {
    // reviewIdがない（通常モード）かつ、以前のレビューデータが残っている場合
    if (!reviewId && currentReview) {
      reset();
      setCodeLocal('');
      setLanguageLocal('python');
    }
  }, [reviewId]); // reviewIdの変化のみを監視

  const handleCodeChange = (value: string | undefined) => {
    if (isViewMode) return; // 閲覧モードでは編集不可
    
    const newCode = value || '';
    setCodeLocal(newCode);
    
    if (newCode.trim().length > 10) {
      const detected = detectLanguage(newCode);
      if (detected !== language) {
        setLanguageLocal(detected);
        setCode(newCode, detected);
        return;
      }
    }
    
    setCode(newCode, language);
  };

  const handleLanguageChange = (newLanguage: string) => {
    if (isViewMode) return; // 閲覧モードでは変更不可
    
    setLanguageLocal(newLanguage);
    setCode(code, newLanguage);
  };

  const handleNewReview = () => {
    if (code.trim() || currentReview) {
      const confirmed = window.confirm(
        '現在の作業をクリアして新規レビューを開始しますか？'
      );
      if (!confirmed) return;
    }
    
    reset();
    setCodeLocal('');
    setLanguageLocal('python');
  };

  const handleReview = async () => {
    await executeReview(code, language, `code.${getFileExtension(language)}`);
  };

  const handleSaveAsKnowledge = async () => {
    if (!currentReview) return;
    
    setIsSaving(true);
    
    try {
      const knowledgeData = createKnowledgeFromReview(currentReview, language);
      await knowledgeApiClient.createKnowledge(knowledgeData);
      
      setToast({
        type: 'success',
        message: 'ナレッジを保存しました！'
      });
      
    } catch (error) {
      console.error('ナレッジ保存エラー:', error);
      
      setToast({
        type: 'error',
        message: 'ナレッジの保存に失敗しました'
      });
    } finally {
      setIsSaving(false);
    }
  };

  const handleBackToHistory = () => {
    navigate('/history');
  };

  const getFileExtension = (lang: string): string => {
    const extensions: Record<string, string> = {
      python: 'py',
      javascript: 'js',
      typescript: 'ts',
      go: 'go',
      java: 'java',
      html: 'html',
      css: 'css',
    };
    return extensions[lang] || 'txt';
  };

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

        <div className="flex-1 flex flex-col p-6 lg:p-8 overflow-hidden animate-fade-in">
          <div className="flex flex-col gap-4 h-full">
            <header className="ml-12 flex-shrink-0">
              <div className="flex items-center gap-4 mb-2">
                {isViewMode && (
                  <button
                    onClick={handleBackToHistory}
                    className="flex items-center gap-2 text-[#6B7280] hover:text-[#111827] transition-colors"
                  >
                    <ArrowLeft size={20} />
                    <span className="text-sm font-medium">履歴に戻る</span>
                  </button>
                )}
                <h1 className="text-[#111827] text-4xl font-black">
                  {isViewMode ? 'レビュー詳細' : 'AIコードレビュー'}
                </h1>
              </div>
              <p className="text-[#6B7280] text-base">
                {isViewMode 
                  ? '過去のレビュー結果を表示しています（読み取り専用）'
                  : 'コードを貼り付けて、AIによるレビューを開始します。'
                }
              </p>
            </header>

            <div className="flex-1 flex flex-col lg:grid lg:grid-cols-10 gap-6 min-h-0 overflow-y-auto lg:overflow-y-hidden">
              {/* コード入力エリア */}
              <div className="lg:col-span-4 flex flex-col gap-4 h-auto lg:h-full">
                <div className="flex items-center gap-2 flex-shrink-0">
                  <label className="flex flex-col" style={{ width: '200px' }}>
                    <p className="text-[#111827] text-sm font-medium mb-2">言語</p>
                    <select
                      value={language}
                      onChange={(e) => handleLanguageChange(e.target.value)}
                      disabled={isViewMode}
                      className="w-full h-12 px-4 rounded-lg border border-gray-300 bg-white text-[#111827] focus:border-[#FBBF24] focus:ring-[#FBBF24] disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {LANGUAGE_OPTIONS.map((option) => (
                        <option key={option.value} value={option.value}>
                          {option.label}
                        </option>
                      ))}
                    </select>
                  </label>
                  
                  {!isViewMode && (
                    <button
                      onClick={handleNewReview}
                      className="px-4 h-12 border border-gray-300 text-[#111827] rounded-lg text-sm font-medium hover:bg-gray-50 transition-colors mt-6 flex-shrink-0"
                      title="新規レビューを開始"
                    >
                      🆕 新規
                    </button>
                  )}
                </div>

                <div className="flex-1 min-h-[500px] lg:min-h-0 rounded-xl overflow-hidden border border-gray-300 flex-shrink-0">
                  <Editor
                    height="100%"
                    defaultLanguage={language}
                    language={language}
                    value={code}
                    onChange={handleCodeChange}
                    theme="vs-dark"
                    options={{
                      minimap: { enabled: false },
                      fontSize: 14,
                      lineNumbers: 'on',
                      scrollBeyondLastLine: false,
                      automaticLayout: true,
                      readOnly: isViewMode, // 閲覧モードでは読み取り専用
                    }}
                  />
                </div>

                {/* LLM情報表示（閲覧モードのみ） */}
                {isViewMode && currentReview && (currentReview.llmProvider || currentReview.llmModel || currentReview.tokensUsed) && (
                  <div className="bg-gray-50 border border-gray-200 rounded-lg p-4 flex-shrink-0">
                    <h3 className="text-sm font-semibold text-[#111827] mb-3 flex items-center gap-2">
                      <AlertCircle size={16} className="text-[#F4C753]" />
                      LLM情報
                    </h3>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-3 text-sm">
                      {currentReview.llmProvider && (
                        <div>
                          <div className="text-[#6B7280] text-xs mb-1">プロバイダー</div>
                          <div className="font-medium text-[#111827]">{currentReview.llmProvider}</div>
                        </div>
                      )}
                      {currentReview.llmModel && (
                        <div>
                          <div className="text-[#6B7280] text-xs mb-1">モデル</div>
                          <div className="font-medium text-[#111827]">{currentReview.llmModel}</div>
                        </div>
                      )}
                      {currentReview.tokensUsed !== undefined && currentReview.tokensUsed !== null && (
                        <div>
                          <div className="text-[#6B7280] text-xs mb-1">使用トークン数</div>
                          <div className="font-medium text-[#111827]">{currentReview.tokensUsed.toLocaleString()} tokens</div>
                        </div>
                      )}
                    </div>
                  </div>
                )}

                {!isViewMode && (
                  <button
                    onClick={handleReview}
                    disabled={isLoading || !code.trim()}
                    className="w-full h-12 px-4 bg-[#FBBF24] text-[#111827] rounded-lg font-bold text-base hover:bg-amber-400 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 flex-shrink-0"
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
                )}
              </div>

              {/* レビュー結果エリア */}
              <div className="lg:col-span-6 flex flex-col bg-white rounded-xl border border-gray-200 h-auto lg:h-full min-h-[600px] overflow-hidden">
                <div className="p-6 border-b border-gray-200 flex items-center justify-between">
                  <h2 className="text-lg font-bold text-[#111827]">レビュー結果</h2>
                  {currentReview && relatedKnowledge.length > 0 && (
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
                        {isViewMode 
                          ? 'レビュー結果を読み込んでいます...'
                          : 'コードを入力して「レビュー実行」をクリックしてください'
                        }
                      </p>
                    </div>
                  )}
                </div>

                {currentReview && !isViewMode && (
                  <div className="p-4 flex items-center justify-between border-t border-gray-200">
                    <FeedbackButtons />
                    
                    <button 
                      onClick={handleSaveAsKnowledge}
                      disabled={isSaving}
                      className="flex items-center gap-2 px-4 h-10 bg-[#F4C753]/20 text-[#111827] rounded-lg text-sm font-bold hover:bg-[#F4C753]/30 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {isSaving ? (
                        <>
                          <Loader size={16} className="animate-spin" />
                          保存中...
                        </>
                      ) : (
                        <>
                          <Bookmark size={16} />
                          ナレッジとして保存
                        </>
                      )}
                    </button>
                  </div>
                )}

                {currentReview && isViewMode && (
                  <div className="p-4 flex items-center justify-between border-t border-gray-200">
                    <FeedbackButtons />
                    
                    <button 
                      onClick={handleSaveAsKnowledge}
                      disabled={isSaving}
                      className="flex items-center gap-2 px-4 h-10 bg-[#F4C753]/20 text-[#111827] rounded-lg text-sm font-bold hover:bg-[#F4C753]/30 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {isSaving ? (
                        <>
                          <Loader size={16} className="animate-spin" />
                          保存中...
                        </>
                      ) : (
                        <>
                          <Bookmark size={16} />
                          ナレッジとして保存
                        </>
                      )}
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* 関連ナレッジサイドパネル */}
        {showKnowledge && relatedKnowledge.length > 0 && (
          <>
            <div
              className="fixed inset-0 bg-black/30 z-40"
              onClick={() => setShowKnowledge(false)}
            />
            <div className="fixed top-0 right-0 h-full w-full lg:w-[500px] bg-white border-l border-gray-200 shadow-2xl z-50 overflow-y-auto animate-slide-in-right">
              <div className="p-6">
                <div className="flex items-center justify-between mb-6">
                  <h3 className="text-xl font-bold text-[#111827]">
                    📚 このレビューで参照されたナレッジ
                  </h3>
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
                      className="p-4 rounded-xl bg-gray-50 border border-gray-200 hover:border-[#F4C753]/50 hover:shadow-md transition-all"
                    >
                      <h4 className="font-bold text-sm text-[#111827] mb-2">
                        {knowledge.title}
                      </h4>
                      
                      <div className="flex gap-2 mb-3">
                        <span className="text-xs font-medium px-2 py-1 rounded-full bg-blue-100 text-blue-700">
                          {getCategoryLabel(knowledge.category as any)}
                        </span>
                        {knowledge.tags.map((tag, i) => (
                          tag !== knowledge.category && (
                            <span
                              key={i}
                              className="text-xs font-medium px-2 py-1 rounded-full bg-gray-100 text-gray-700"
                            >
                              {tag}
                            </span>
                          )
                        ))}
                      </div>
                      
                      <p className="text-xs text-[#6B7280]">
                        {knowledge.description}
                      </p>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </>
        )}

        {/* トースト通知 */}
        {toast && (
          <Toast
            type={toast.type}
            message={toast.message}
            onClose={() => setToast(null)}
          />
        )}
      </main>
    </div>
  );
}
