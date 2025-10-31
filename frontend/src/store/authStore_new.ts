import { create } from 'zustand';
import { syncUser } from '../api/auth';

interface User {
  email: string;
  name: string;
  sub: string; // Auth0 User ID
}

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  setUser: (user: User | null) => void;
  setAuthenticated: (isAuthenticated: boolean) => void;
  setLoading: (isLoading: boolean) => void;
  setError: (error: string | null) => void;
  logout: () => void;
  clearError: () => void;
  syncUserWithBackend: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,

  setUser: (user) => set({ user }),
  setAuthenticated: (isAuthenticated) => set({ isAuthenticated }),
  setLoading: (isLoading) => set({ isLoading }),
  setError: (error) => set({ error }),

  // バックエンドとユーザー同期
  syncUserWithBackend: async () => {
    try {
      await syncUser();
      console.log('✅ ユーザー同期成功');
    } catch (error) {
      console.error('⚠️ ユーザー同期失敗:', error);
      throw error;
    }
  },

  // ログアウト
  logout: () => {
    set({
      user: null,
      isAuthenticated: false,
      error: null,
    });
  },

  // エラークリア
  clearError: () => {
    set({ error: null });
  },
}));
