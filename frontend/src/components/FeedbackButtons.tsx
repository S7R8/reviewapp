import { useState } from 'react';
import { ThumbsUp, ThumbsDown } from 'lucide-react';
import { useReviewStore } from '../store/reviewStore';

export const FeedbackButtons = () => {
  const { feedbackScore, isSubmittingFeedback, submitFeedback } = useReviewStore();
  const [localScore, setLocalScore] = useState<number | null>(feedbackScore);

  const handleFeedback = async (score: number) => {
    setLocalScore(score);
    await submitFeedback(score);
  };

  return (
    <div className="flex items-center gap-4">
      <p className="text-sm text-[#6B7280]">
        このレビューは役に立ちましたか？
      </p>
      
      <button
        onClick={() => handleFeedback(1)}
        disabled={isSubmittingFeedback}
        className={`p-2 rounded-lg transition-all ${
          localScore === 1
            ? 'bg-red-100 text-red-600 shadow-sm'
            : 'text-[#6B7280] hover:bg-gray-100 hover:text-red-500'
        } ${isSubmittingFeedback ? 'opacity-50 cursor-not-allowed' : ''}`}
        title="役に立たなかった"
      >
        <ThumbsDown size={18} />
      </button>
      
      <button
        onClick={() => handleFeedback(3)}
        disabled={isSubmittingFeedback}
        className={`p-2 rounded-lg transition-all ${
          localScore === 3
            ? 'bg-[#F4C753] text-[#111827] shadow-sm'
            : 'text-[#6B7280] hover:bg-gray-100 hover:text-[#F4C753]'
        } ${isSubmittingFeedback ? 'opacity-50 cursor-not-allowed' : ''}`}
        title="役に立った"
      >
        <ThumbsUp size={18} />
      </button>

      {isSubmittingFeedback && (
        <div className="flex items-center gap-2 text-xs text-[#6B7280]">
          <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-[#6B7280]" />
          送信中...
        </div>
      )}

      {localScore && !isSubmittingFeedback && (
        <span className="text-xs text-green-600 font-medium">
          ✓ 送信完了
        </span>
      )}
    </div>
  );
};
