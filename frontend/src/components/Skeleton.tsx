// シンプルなスケルトンコンポーネント

export function SkeletonLine({ className = '' }: { className?: string }) {
  return (
    <div className={`h-4 bg-gray-200 rounded animate-pulse ${className}`} />
  );
}

export function SkeletonBox({ className = '' }: { className?: string }) {
  return (
    <div className={`bg-gray-200 rounded animate-pulse ${className}`} />
  );
}

export function SkeletonCard({ className = '' }: { className?: string }) {
  return (
    <div className={`p-6 bg-white rounded-xl border border-gray-200 ${className}`}>
      <div className="space-y-3">
        <SkeletonLine className="w-3/4" />
        <SkeletonLine className="w-full" />
        <SkeletonLine className="w-5/6" />
      </div>
    </div>
  );
}

// レビュー結果用のスケルトン
export function ReviewResultSkeleton() {
  return (
    <div className="space-y-6 p-6">
      {/* サマリー */}
      <div className="space-y-2">
        <SkeletonLine className="w-full" />
        <SkeletonLine className="w-5/6" />
      </div>

      {/* 良い点 */}
      <div className="p-4 rounded-lg bg-gray-50 border border-gray-200">
        <SkeletonBox className="h-5 w-32 mb-3" />
        <div className="space-y-2">
          <SkeletonLine className="w-full" />
          <SkeletonLine className="w-4/5" />
        </div>
      </div>

      {/* 改善点 */}
      <div className="p-4 rounded-lg bg-gray-50 border border-gray-200">
        <SkeletonBox className="h-5 w-32 mb-3" />
        <div className="space-y-4">
          <div>
            <SkeletonLine className="w-2/3 mb-2" />
            <SkeletonLine className="w-full" />
            <SkeletonLine className="w-5/6" />
            <SkeletonBox className="h-20 w-full mt-2" />
          </div>
          <div>
            <SkeletonLine className="w-2/3 mb-2" />
            <SkeletonLine className="w-full" />
            <SkeletonLine className="w-4/5" />
          </div>
        </div>
      </div>

      {/* 根拠 */}
      <div className="space-y-2">
        <SkeletonLine className="w-48" />
      </div>
    </div>
  );
}
