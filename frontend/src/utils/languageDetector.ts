/**
 * コードから言語を自動判定する
 * 
 * @param code - 判定対象のコード
 * @returns 検出された言語 (python, javascript, go, など)
 */
export const detectLanguage = (code: string): string => {
  const trimmed = code.trim();
  
  if (!trimmed) return 'python';
  
  // 各言語の特徴的なパターンを定義
  const patterns = [
    {
      language: 'python',
      patterns: [
        /^def\s+\w+\s*\(/m,
        /^class\s+\w+/m,
        /^import\s+\w+/m,
        /^from\s+\w+\s+import/m,
        /if\s+__name__\s*==\s*['"]__main__['"]/,
        /print\s*\(/,
        /:\s*$/m, // Pythonの行末コロン
      ],
      weight: 1,
    },
    {
      language: 'javascript',
      patterns: [
        /^(const|let|var)\s+\w+/m,
        /^function\s+\w+/m,
        /=>\s*\{/,
        /console\.(log|error|warn)/,
        /^import\s+.*from\s+['"]/m,
        /^export\s+(default|const|function)/m,
        /\.then\s*\(/,
        /async\s+function/,
      ],
      weight: 1,
    },
    {
      language: 'typescript',
      patterns: [
        /:\s*(string|number|boolean|any|void|unknown)/,
        /^interface\s+\w+/m,
        /^type\s+\w+\s*=/m,
        /<[\w\s,]+>\s*\(/,
        /as\s+(string|number|boolean)/,
      ],
      weight: 2, // TypeScriptパターンは重み付け
    },
    {
      language: 'go',
      patterns: [
        /^package\s+\w+/m,
        /^func\s+\w+/m,
        /^type\s+\w+\s+struct/m,
        /fmt\.(Print|Sprintf|Println)/,
        /^import\s+\(/m,
        /:=\s*/,
        /^func\s+\(\w+\s+\*?\w+\)/m, // メソッド定義
      ],
      weight: 1,
    },
    {
      language: 'java',
      patterns: [
        /^public\s+class\s+\w+/m,
        /^private\s+(static\s+)?[\w<>]+\s+\w+/m,
        /System\.(out|err)\.(print|println)/,
        /^import\s+java\./m,
        /^@\w+/m, // アノテーション
        /public\s+static\s+void\s+main/,
      ],
      weight: 1,
    },
    {
      language: 'html',
      patterns: [
        /^<!DOCTYPE\s+html>/i,
        /^<html/i,
        /<(div|span|p|body|head|h\d)[\s>]/i,
        /<\/\w+>/,
      ],
      weight: 1,
    },
    {
      language: 'css',
      patterns: [
        /^(\.|#|[\w-]+)\s*\{/m,
        /@media\s+/,
        /@import\s+/,
        /:\s*[\w-]+\s*;/,
      ],
      weight: 1,
    },
  ];

  // 各言語のスコアを計算
  const scores: Record<string, number> = {};
  
  for (const { language, patterns: langPatterns, weight } of patterns) {
    let score = 0;
    
    for (const pattern of langPatterns) {
      if (pattern.test(trimmed)) {
        score += weight;
      }
    }
    
    if (score > 0) {
      scores[language] = score;
    }
  }

  // TypeScriptとJavaScriptの特殊処理
  if (scores['typescript'] && scores['javascript']) {
    // TypeScript特有のパターンがあれば優先
    if (scores['typescript'] >= 2) {
      delete scores['javascript'];
    } else {
      // そうでなければJavaScriptとして扱う
      scores['javascript'] += scores['typescript'];
      delete scores['typescript'];
    }
  }

  // 最高スコアの言語を返す
  let maxScore = 0;
  let detectedLanguage = 'python'; // デフォルト

  for (const [lang, score] of Object.entries(scores)) {
    if (score > maxScore) {
      maxScore = score;
      detectedLanguage = lang;
    }
  }

  // スコアが低すぎる場合はデフォルトに戻す
  if (maxScore < 2) {
    return 'python';
  }

  return detectedLanguage;
};

/**
 * 検出された言語名を表示用ラベルに変換
 */
export const getLanguageLabel = (language: string): string => {
  const labels: Record<string, string> = {
    python: 'Python',
    javascript: 'JavaScript',
    typescript: 'TypeScript',
    go: 'Go',
    java: 'Java',
    html: 'HTML',
    css: 'CSS',
  };
  
  return labels[language] || language.charAt(0).toUpperCase() + language.slice(1);
};

/**
 * デバッグ用: 検出ロジックの詳細を取得
 */
export const detectLanguageWithDetails = (code: string): {
  language: string;
  confidence: number;
  matches: Array<{ language: string; score: number }>;
} => {
  const trimmed = code.trim();
  
  if (!trimmed) {
    return {
      language: 'python',
      confidence: 0,
      matches: [],
    };
  }

  // 簡易版: 実際の検出ロジックを再利用
  const detected = detectLanguage(code);
  
  return {
    language: detected,
    confidence: 0.8, // 仮の信頼度
    matches: [{ language: detected, score: 1 }],
  };
};
