import React from 'react';
import { ReviewStatus, StatusBadgeConfig } from '../types/review';

interface StatusBadgeProps {
  status: ReviewStatus;
  className?: string;
}

/**
 * ステータスバッジコンポーネント
 * レビューのステータスを視覚的に表示
 */
const StatusBadge: React.FC<StatusBadgeProps> = ({ status, className = '' }) => {
  const getStatusConfig = (status: ReviewStatus): StatusBadgeConfig => {
    switch (status) {
      case 'success':
        return {
          label: '成功',
          bgColor: 'bg-green-100',
          textColor: 'text-green-800',
        };
      case 'warning':
        return {
          label: '改善点あり',
          bgColor: 'bg-yellow-100',
          textColor: 'text-yellow-800',
        };
      case 'error':
        return {
          label: 'エラー',
          bgColor: 'bg-red-100',
          textColor: 'text-red-800',
        };
      case 'pending':
        return {
          label: '処理中',
          bgColor: 'bg-blue-100',
          textColor: 'text-blue-800',
        };
      default:
        return {
          label: '不明',
          bgColor: 'bg-gray-100',
          textColor: 'text-gray-800',
        };
    }
  };

  const config = getStatusConfig(status);

  return (
    <span
      className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${config.bgColor} ${config.textColor} ${className}`}
    >
      {config.label}
    </span>
  );
};

export default StatusBadge;
