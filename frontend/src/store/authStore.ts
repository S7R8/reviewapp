import { create } from 'zustand';

interface User {
  email: string;
  name: string;
}

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,

  // 仮実装：何でも通す
  login: async (email: string, password: string) => {
    // TODO: 後でバックエンドAPI呼び出しに置き換える
    console.log('仮ログイン:', { email, password });
    
    // 1秒待つ（ローディング状態を確認するため）
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    // ログイン成功
    set({
      user: {
        email,
        name: '山田 太郎', // 仮のユーザー名
      },
      isAuthenticated: true,
    });

    // localStorageに保存（ページリロード対応）
    localStorage.setItem('auth_token', 'dummy_token');
    localStorage.setItem('user_email', email);
  },

  logout: () => {
    set({
      user: null,
      isAuthenticated: false,
    });
    
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user_email');
  },
}));
