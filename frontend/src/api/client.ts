/**
 * 共通APIクライアント
 * すべてのAPI呼び出しに認証ヘッダーを自動付与
 */

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export interface ApiError {
  error: string;
  message: string;
  details?: Record<string, unknown>;
}

/**
 * 認証ヘッダーを取得
 * 注: audience なしでAuth0を使用する場合、access_tokenはJWE形式（暗号化）で返されるため、
 * id_tokenを使用します。id_tokenはJWT形式で、バックエンドで検証可能です。
 */
export const getAuthHeaders = (): HeadersInit => {
  // id_tokenを優先的に使用（JWT形式）
  const idToken = localStorage.getItem('id_token');
  const accessToken = localStorage.getItem('access_token');
  
  // id_tokenがJWT形式（3パーツ）かチェック
  const token = idToken && idToken.split('.').length === 3 ? idToken : accessToken;
  
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  return headers;
};

/**
 * 共通fetchラッパー
 */
export const apiFetch = async <T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> => {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    ...options,
    headers: {
      ...getAuthHeaders(),
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error: ApiError = await response.json().catch(() => ({
      error: 'UnknownError',
      message: `HTTP ${response.status}: ${response.statusText}`,
    }));
    
    throw new Error(error.message || 'APIリクエストに失敗しました');
  }

  return response.json();
};

/**
 * GET リクエスト
 */
export const apiGet = async <T>(endpoint: string): Promise<T> => {
  return apiFetch<T>(endpoint, { method: 'GET' });
};

/**
 * POST リクエスト
 */
export const apiPost = async <T>(endpoint: string, body: any): Promise<T> => {
  return apiFetch<T>(endpoint, {
    method: 'POST',
    body: JSON.stringify(body),
  });
};

/**
 * PUT リクエスト
 */
export const apiPut = async <T>(endpoint: string, body: any): Promise<T> => {
  return apiFetch<T>(endpoint, {
    method: 'PUT',
    body: JSON.stringify(body),
  });
};

/**
 * DELETE リクエスト
 */
export const apiDelete = async <T>(endpoint: string): Promise<T> => {
  return apiFetch<T>(endpoint, { method: 'DELETE' });
};

export { API_BASE_URL };
