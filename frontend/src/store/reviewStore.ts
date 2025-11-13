import { create } from 'zustand';
import { ReviewResult, Knowledge } from '../types/review';
import { reviewApiClient } from '../api/reviewApi';
import { knowledgeApiClient, Knowledge as ApiKnowledge } from '../api/knowledgeApi';

interface ReviewState {
  currentReview: ReviewResult | null;
  relatedKnowledge: Knowledge[];
  isLoading: boolean;
  error: string | null;
  feedbackScore: number | null;
  isSubmittingFeedback: boolean;
  
  // ★ コードと言語の状態を追加
  currentCode: string;
  currentLanguage: string;
  
  // レビュー実行（実際のAPI呼び出し）
  executeReview: (code: string, language: string, filename: string) => Promise<void>;
  
  // フィードバック送信
  submitFeedback: (score: number, comment?: string) => Promise<void>;
  
  // ★ コードと言語を保存
  setCode: (code: string, language: string) => void;
  
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

export const useReviewStore = create<ReviewState>((set, get) => ({
  currentReview: null,
  relatedKnowledge: [],
  isLoading: false,
  error: null,
  feedbackScore: null,
  isSubmittingFeedback: false,
  currentCode: '',           // ★ 初期値
  currentLanguage: 'python', // ★ 初期値

  // ★ コードと言語を保存するメソッド
  setCode: (code: string, language: string) => {
    set({ currentCode: code, currentLanguage: language });
  },

  executeReview: async (code: string, language: string, filename: string) => {
    set({ 
      isLoading: true, 
      error: null,
      currentCode: code,      // ★ レビュー実行時に保存
      currentLanguage: language
    });

    try {
      // 実際のAPI呼び出し
      const review = await reviewApiClient.reviewCode({
        code,
        language,
        filename,
      });

      // ★ 関連ナレッジを取得（referenced_knowledgeがあれば）
      let knowledgeList: Knowledge[] = [];
      if (review.referencedKnowledgeIds && review.referencedKnowledgeIds.length > 0) {
        try {
          // 複数のIDからナレッジを取得
          const knowledgeDetails = await Promise.all(
            review.referencedKnowledgeIds.map(id => knowledgeApiClient.getKnowledgeById(id))
          );
          
          // APIのKnowledge型をフロントの型に変換
          knowledgeList = knowledgeDetails.map((k: ApiKnowledge) => ({
            id: k.id,
            title: k.title,
            description: k.content,
            category: k.category,
            tags: [k.category], // カテゴリをタグとして使用
          }));
        } catch (error) {
          console.error('ナレッジ取得エラー:', error);
          // エラー時はモックデータを使用
          knowledgeList = generateMockKnowledge();
        }
      } else {
        // referenced_knowledgeがない場合はモック
        knowledgeList = generateMockKnowledge();
      }

      set({
        currentReview: review,
        relatedKnowledge: knowledgeList,
        isLoading: false,
        feedbackScore: null, // リセット
      });
    } catch (error) {
      console.error('レビュー実行エラー:', error);
      
      set({
        error: error instanceof Error ? error.message : 'レビューの実行に失敗しました',
        isLoading: false,
      });
    }
  },

  submitFeedback: async (score: number, comment?: string) => {
    const { currentReview } = get();
    if (!currentReview) {
      console.error('レビューが存在しません');
      return;
    }

    set({ isSubmittingFeedback: true });

    try {
      await reviewApiClient.updateFeedback(currentReview.id, score, comment);
      set({ 
        feedbackScore: score,
        isSubmittingFeedback: false,
      });
      console.log('フィードバックを送信しました:', { score, comment });
    } catch (error) {
      console.error('フィードバックの送信に失敗しました:', error);
      set({ isSubmittingFeedback: false });
    }
  },

  reset: () => {
    set({
      currentReview: null,
      relatedKnowledge: [],
      isLoading: false,
      error: null,
      feedbackScore: null,
      isSubmittingFeedback: false,
      currentCode: '',           // ★ リセット
      currentLanguage: 'python', // ★ リセット
    });
  },
}));
