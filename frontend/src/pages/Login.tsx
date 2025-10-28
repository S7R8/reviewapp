import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { Mail, Lock, Eye, EyeOff, Hexagon } from 'lucide-react';

export default function Login() {
  const navigate = useNavigate();
  const login = useAuthStore((state) => state.login);

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      await login(email, password);
      navigate('/dashboard');
    } catch (error) {
      console.error('ログインエラー:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="relative flex min-h-screen w-full flex-col items-center justify-center overflow-x-hidden p-4 sm:p-6 lg:p-8">
      {/* 背景グラデーション */}
      <div className="absolute top-0 left-0 w-full h-full bg-gradient-to-br from-[#F4C753]/10 via-transparent to-[#F4C753]/5 -z-10"></div>

      {/* ヘッダー（ロゴ） */}
      <header className="absolute top-0 left-0 w-full p-6">
        <a className="flex items-center gap-2 text-xl font-bold text-[#111827]" href="#">
          <Hexagon className="text-[#F4C753]" size={32} />
          <span>CodeRev</span>
        </a>
      </header>

      {/* メインカード */}
      <main className="w-full max-w-md">
        <div className="rounded-xl border border-gray-200 bg-white p-8 shadow-sm">
          {/* タイトル */}
          <div className="flex flex-col items-center">
            <h1 className="text-[#111827] tracking-tight text-3xl font-bold leading-tight text-center pb-2">
              おかえりなさい
            </h1>
            <p className="text-[#6B7280] text-base font-normal leading-normal pb-8 text-center">
              アカウントにログインしてください
            </p>
          </div>

          {/* フォーム */}
          <form onSubmit={handleSubmit} className="flex flex-col gap-6">
            {/* メールアドレス */}
            <label className="flex flex-col w-full">
              <p className="text-[#111827] text-sm font-medium leading-normal pb-2">
                メールアドレス
              </p>
              <div className="relative flex w-full items-center">
                <Mail className="absolute left-3 text-[#6B7280]" size={20} />
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="flex w-full min-w-0 flex-1 resize-none overflow-hidden rounded-lg border border-gray-300 bg-white placeholder:text-gray-400 focus:border-[#FBBF24] focus:ring-[#FBBF24] h-12 py-2 pl-10 pr-4 text-base font-normal leading-normal text-[#111827]"
                  placeholder="your@email.com"
                  required
                />
              </div>
            </label>

            {/* パスワード */}
            <label className="flex flex-col w-full">
              <p className="text-[#111827] text-sm font-medium leading-normal pb-2">
                パスワード
              </p>
              <div className="relative flex w-full items-center">
                <Lock className="absolute left-3 text-[#6B7280]" size={20} />
                <input
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="flex w-full min-w-0 flex-1 resize-none overflow-hidden rounded-lg border border-gray-300 bg-white placeholder:text-gray-400 focus:border-[#FBBF24] focus:ring-[#FBBF24] h-12 py-2 pl-10 pr-10 text-base font-normal leading-normal text-[#111827]"
                  placeholder="パスワードを入力"
                  required
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  aria-label="パスワードを表示"
                  className="absolute right-3 text-[#6B7280] hover:text-[#111827] transition-colors"
                >
                  {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
                </button>
              </div>
            </label>

            {/* パスワード忘れリンク */}
            <div className="flex justify-end">
              <a className="text-[#F4C753] text-sm font-medium hover:underline" href="#">
                パスワードをお忘れですか？
              </a>
            </div>

            {/* ログインボタン */}
            <button
              type="submit"
              disabled={isLoading}
              className="flex min-w-[84px] cursor-pointer items-center justify-center overflow-hidden rounded-lg h-12 px-5 w-full bg-[#FBBF24] text-[#111827] text-base font-bold leading-normal tracking-[0.015em] transition-colors hover:bg-amber-400 focus:outline-none focus:ring-2 focus:ring-amber-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isLoading ? (
                <span className="flex items-center gap-2">
                  <svg className="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  ログイン中...
                </span>
              ) : (
                <span className="truncate">ログイン</span>
              )}
            </button>
          </form>

          {/* または区切り */}
          <div className="relative flex items-center py-6">
            <div className="flex-grow border-t border-gray-300"></div>
            <span className="mx-4 flex-shrink text-sm text-[#6B7280]">または</span>
            <div className="flex-grow border-t border-gray-300"></div>
          </div>

          {/* OAuth ボタン */}
          <div className="flex flex-col gap-4">
            {/* Google */}
            <button className="flex h-12 w-full cursor-pointer items-center justify-center gap-3 overflow-hidden rounded-lg border border-gray-300 bg-white px-5 text-base font-medium leading-normal text-[#111827] transition-colors hover:bg-gray-50">
              <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <g clipPath="url(#clip0_105_2393)">
                  <path d="M22.56 12.25C22.56 11.45 22.49 10.68 22.36 9.92H12V14.28H18.19C17.92 15.77 17.06 17.08 15.82 17.94V20.59H19.66C21.56 18.83 22.56 15.83 22.56 12.25Z" fill="#4285F4"></path>
                  <path d="M12 23C15.06 23 17.67 22.01 19.66 20.59L15.82 17.94C14.75 18.66 13.45 19.12 12 19.12C9.11 19.12 6.63 17.24 5.72 14.66H1.73V17.41C3.67 20.8 7.55 23 12 23Z" fill="#34A853"></path>
                  <path d="M5.72 14.66C5.46 13.92 5.32 13.13 5.32 12.33C5.32 11.53 5.46 10.74 5.72 10H1.73C0.94 11.57 0.49 13.31 0.49 15.17C0.49 17.03 0.94 18.77 1.73 20.34L5.72 17.59V14.66Z" fill="#FBBC05"></path>
                  <path d="M12 5.48C13.67 5.48 15.14 6.08 16.29 7.15L20.03 3.41C17.67 1.3 15.06 0 12 0C7.55 0 3.67 2.2 1.73 5.59L5.72 8.34C6.63 5.76 9.11 3.88 12 3.88V5.48Z" fill="#EA4335"></path>
                </g>
                <defs>
                  <clipPath id="clip0_105_2393">
                    <rect fill="white" height="24" width="24"></rect>
                  </clipPath>
                </defs>
              </svg>
              <span className="truncate">Googleでログイン</span>
            </button>

            {/* GitHub */}
            <button className="flex h-12 w-full cursor-pointer items-center justify-center gap-3 overflow-hidden rounded-lg border border-gray-300 bg-white px-5 text-base font-medium leading-normal text-[#111827] transition-colors hover:bg-gray-50">
              <svg className="h-6 w-6" fill="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path clipRule="evenodd" d="M12 2C6.477 2 2 6.477 2 12C2 16.417 4.86 20.165 8.84 21.49C9.48 21.6 9.68 21.21 9.68 20.88C9.68 20.58 9.67 19.55 9.67 18.45C6.98 19.01 6.13 17.09 6.13 17.09C5.55 15.63 4.56 15.22 4.56 15.22C3.41 14.43 4.65 14.44 4.65 14.44C5.92 14.53 6.64 15.73 6.64 15.73C7.79 17.7 9.79 17.09 10.51 16.78C10.62 15.93 10.96 15.35 11.33 15.03C8.61 14.72 5.72 13.68 5.72 9.7C5.72 8.57 6.13 7.65 6.78 6.94C6.66 6.63 6.27 5.53 6.9 4.38C6.9 4.38 7.97 4.04 9.65 5.21C10.67 4.93 11.75 4.79 12.83 4.79C13.9 4.79 14.98 4.93 16 5.21C17.68 4.04 18.75 4.38 18.75 4.38C19.38 5.53 18.99 6.63 18.87 6.94C19.52 7.65 19.93 8.57 19.93 9.7C19.93 13.69 17.03 14.72 14.31 15.03C14.77 15.42 15.15 16.19 15.15 17.41C15.15 19.14 15.14 20.53 15.14 20.88C15.14 21.21 15.33 21.61 15.98 21.49C19.96 20.16 22.82 16.41 22.82 12C22.82 6.477 18.347 2 12.82 2H12Z" fillRule="evenodd"></path>
              </svg>
              <span className="truncate">GitHubでログイン</span>
            </button>
          </div>
        </div>

        {/* サインアップリンク */}
        <div className="pt-6">
          <p className="text-[#6B7280] text-sm font-normal leading-normal text-center">
            アカウントをお持ちでないですか？{' '}
            <Link to="/signup" className="font-bold text-[#F4C753] hover:underline">
              新規登録はこちら
            </Link>
          </p>
        </div>
      </main>
    </div>
  );
}
