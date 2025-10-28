import { create } from 'zustand';
import { ReviewResult, Knowledge } from '../types/review';
import { reviewApiClient } from '../api/reviewApi';

interface ReviewState {
  currentReview: ReviewResult | null;
  relatedKnowledge: Knowledge[];
  isLoading: boolean;
  error: string | null;
  
  // レビュー実行（実際のAPI呼び出し）
  executeReview: (code: string, language: string, filename: string) => Promise<void>;
  
  // リセット
  reset: () => void;
}

// 仮の関連ナレッジを生成（将来的にはAPIから取得）
const generateMockKnowledge = (): Knowledge[] => {
  return [
    {
      id: '1',
      title: 'エラーハンドリングのベストプラクティス',
      description: 'try-catchを使用した適切なエラー処理の方法。ユーザー向けメッセージと開発者向け詳細を分ける。',
      category: 'Error Handling',
      tags: ['JavaScript', 'Best Practice'],
    },
    {
      id: '2',
      title: '効果的なコメントの書き方',
      description: 'コードの「なぜ」を説明するコメントを書く。「何を」しているかは書かない。',
      category: 'Documentation',
      tags: ['Documentation', 'Clean Code'],
    },
    {
      id: '3',
      title: 'JavaScriptのNullチェック',
      description: 'Optional Chaining (?.) を用いた安全なプロパティアクセスの実践方法。',
      category: 'JavaScript',
      tags: ['JavaScript', 'Safety'],
    },
  ];
};

export const useReviewStore = create<ReviewState>((set) => ({
  currentReview: null,
  relatedKnowledge: [],
  isLoading: false,
  error: null,

  executeReview: async (code: string, language: string, filename: string) => {
    set({ isLoading: true, error: null });

    try {
      // 実際のAPI呼び出し
      const review = await reviewApiClient.reviewCode({
        code,
        language,
        filename,
      });

      // TODO: 将来的には関連ナレッジもAPIから取得
      const mockKnowledge = generateMockKnowledge();

      set({
        currentReview: review,
        relatedKnowledge: mockKnowledge,
        isLoading: false,
      });
    } catch (error) {
      console.error('レビュー実行エラー:', error);
      
      set({
        error: error instanceof Error ? error.message : 'レビューの実行に失敗しました',
        isLoading: false,
      });
    }
  },

  reset: () => {
    set({
      currentReview: null,
      relatedKnowledge: [],
      isLoading: false,
      error: null,
    });
  },
}));
