import { apiPost } from './client';

interface User {
  id: string;
  auth0_user_id: string;
  email: string;
  name: string;
  avatar_url: string | null;
  preferences: string;
  created_at: string;
  updated_at: string;
}

interface SyncUserResponse {
  user: User;
}

/**
 * Auth0ログイン後にバックエンドのユーザーと同期
 */
export const syncUser = async (): Promise<User> => {
  const response = await apiPost<SyncUserResponse>('/api/v1/auth/sync', {});
  return response.user;
};
