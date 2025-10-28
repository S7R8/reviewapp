export interface ReviewResult {
  id: string;
  summary: string;
  goodPoints: string[];
  improvements: Improvement[];
  references: Reference[];
  createdAt: string;
  rawMarkdown: string;
}

export interface Improvement {
  title: string;
  description: string;
  codeAfter?: string;
  severity: 'low' | 'medium' | 'high';
}

export interface Reference {
  source: string;
  description: string;
}

export interface Knowledge {
  id: string;
  title: string;
  description: string;
  category: string;
  tags: string[];
}

export interface ReviewRequest {
  code: string;
  language: string;
  filename: string;
}
