import { create } from 'zustand';
import axios from 'axios';
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
  login: (email: string, password: string) => Promise<void>;
  signup: (email: string, password: string, name: string) => Promise<void>;
  logout: () => void;
  clearError: () => void;
  initializeAuth: () => Promise<void>;
}

// Auth0設定
const AUTH0_DOMAIN = import.meta.env.VITE_AUTH0_DOMAIN;
const AUTH0_CLIENT_ID = import.meta.env.VITE_AUTH0_CLIENT_ID;
const AUTH0_CLIENT_SECRET = import.meta.env.VITE_AUTH0_CLIENT_SECRET;
const AUTH0_AUDIENCE = import.meta.env.VITE_AUTH0_AUDIENCE;

// JWTからユーザー情報を抽出
const parseJWT = (token: string): User | null => {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split('')
        .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    );

    const payload = JSON.parse(jsonPayload);

    return {
      email: payload.email || '',
      name: payload.name || payload.email || '',
      sub: payload.sub || '',
    };
  } catch (error) {
    console.error('JWT解析エラー:', error);
    return null;
  }
};

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,

  // アプリ起動時の認証状態復元
  initializeAuth: async () => {
    const accessToken = localStorage.getItem('access_token');
    const idToken = localStorage.getItem('id_token');

    if (accessToken && idToken) {
      const user = parseJWT(idToken);
      if (user) {
        set({
          user,
          isAuthenticated: true,
        });

        // バックエンドとユーザー同期
        try {
          await syncUser();
        } catch (error) {
          console.error('ユーザー同期エラー:', error);
          // 同期失敗しても認証状態は維持
        }
      }
    }
  },

  // ログイン（Auth0 Password Grant）
  login: async (email: string, password: string) => {
    set({ isLoading: true, error: null });

    try {
      // Auth0 Token Endpoint
      const response = await axios.post(
        `https://${AUTH0_DOMAIN}/oauth/token`,
        {
          grant_type: 'http://auth0.com/oauth/grant-type/password-realm',
          username: email,
          password: password,
          client_id: AUTH0_CLIENT_ID,
          client_secret: AUTH0_CLIENT_SECRET,
          // audience: AUTH0_AUDIENCE, // Auth0の設定でPassword Grantが許可されていないためコメントアウト
          realm: 'Username-Password-Authentication',
          scope: 'openid profile email',
        }
      );

      const { access_token, id_token } = response.data;

      // access_tokenがJWT形式かチェック
      let user = null;
      if (access_token && access_token.split('.').length === 3) {
        // access_tokenがJWT形式の場合
        user = parseJWT(access_token);
      }
      
      // access_tokenから取得できない場合はid_tokenから
      if (!user && id_token) {
        user = parseJWT(id_token);
      }

      if (!user) {
        throw new Error('ユーザー情報の取得に失敗しました');
      }

      // トークンを保存
      localStorage.setItem('access_token', access_token);
      localStorage.setItem('id_token', id_token);

      // 状態更新
      set({
        user,
        isAuthenticated: true,
        isLoading: false,
      });

      // バックエンドとユーザー同期
      try {
      await syncUser();
      } catch (error) {
        console.error('ユーザー同期エラー:', error);
      // 同期失敗してもログインは成功扱い
      }
    } catch (error: any) {
      console.error('ログインエラー:', error);

      let errorMessage = 'ログインに失敗しました';
      if (error.response?.data?.error_description) {
        errorMessage = error.response.data.error_description;
      } else if (error.response?.data?.error) {
        errorMessage = error.response.data.error;
      } else if (error.response?.status === 403 || error.response?.status === 401) {
        errorMessage = 'メールアドレスまたはパスワードが間違っています';
      } else if (error.response?.status === 400) {
        errorMessage = 'リクエストが無効です。Auth0の設定を確認してください。';
      }

      set({
        error: errorMessage,
        isLoading: false,
      });

      throw new Error(errorMessage);
    }
  },

  // サインアップ（Auth0 Signup API）
  signup: async (email: string, password: string, name: string) => {
    set({ isLoading: true, error: null });

    try {
      // 1. Auth0でユーザー作成
      await axios.post(
        `https://${AUTH0_DOMAIN}/dbconnections/signup`,
        {
          client_id: AUTH0_CLIENT_ID,
          email: email,
          password: password,
          connection: 'Username-Password-Authentication',
          name: name,
        }
      );

      // 2. 自動ログイン
      await useAuthStore.getState().login(email, password);
    } catch (error: any) {
      console.error('サインアップエラー:', error);

      let errorMessage = 'サインアップに失敗しました';
      if (error.response?.data?.description) {
        errorMessage = error.response.data.description;
      } else if (error.response?.data?.message) {
        errorMessage = error.response.data.message;
      }

      set({
        error: errorMessage,
        isLoading: false,
      });

      throw new Error(errorMessage);
    }
  },

  // ログアウト
  logout: () => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('id_token');

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
