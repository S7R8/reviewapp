import React from 'react';
import ReactMarkdown from 'react-markdown';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { ReviewResult as ReviewResultType } from '../types/review';

interface ReviewResultProps {
  result: ReviewResultType;
}

export const ReviewResult: React.FC<ReviewResultProps> = ({ result }) => {
  return (
    <div className="bg-white rounded-lg shadow-lg p-6">
      <h2 className="text-2xl font-bold mb-6 text-gray-800">
        レビュー結果
      </h2>

      {/* マークダウンをそのまま表示 */}
      <div className="prose prose-lg max-w-none">
        <ReactMarkdown
          components={{
            code({ node, inline, className, children, ...props }) {
              const match = /language-(\w+)/.exec(className || '');
              return !inline && match ? (
                <SyntaxHighlighter
                  style={vscDarkPlus}
                  language={match[1]}
                  PreTag="div"
                  {...props}
                >
                  {String(children).replace(/\n$/, '')}
                </SyntaxHighlighter>
              ) : (
                <code className={className} {...props}>
                  {children}
                </code>
              );
            },
          }}
        >
          {result.rawMarkdown || ''}
        </ReactMarkdown>
      </div>

      {/* メタ情報 */}
      <div className="mt-8 pt-6 border-t border-gray-200">
        <p className="text-sm text-gray-500">
          レビュー日時: {new Date(result.createdAt).toLocaleString('ja-JP')}
        </p>
      </div>
    </div>
  );
};
