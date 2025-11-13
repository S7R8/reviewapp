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
          // エラー時は空配列
          knowledgeList = [];
        }
      }
      // referenced_knowledgeがない場合は空配列（モックは使わない）

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
