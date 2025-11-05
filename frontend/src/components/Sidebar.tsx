import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import {
  Code,
  Book,
  History,
  Settings,
  LogOut,
  BarChart3,
  Plus,
  Menu,
  X
} from 'lucide-react';

interface SidebarProps {
  currentPage: 'dashboard' | 'review' | 'knowledge' | 'history' | 'settings';
  isOpen: boolean;
  onToggle: () => void;
}

export default function Sidebar({ currentPage, isOpen, onToggle }: SidebarProps) {
  const navigate = useNavigate();
  const { user, logout } = useAuthStore();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const menuItems = [
    { id: 'dashboard', label: 'ダッシュボード', icon: BarChart3, path: '/dashboard' },
    { id: 'review', label: 'コードレビュー', icon: Code, path: '/review' },
    { id: 'knowledge', label: 'ナレッジベース', icon: Book, path: '/knowledge' },
    { id: 'history', label: '履歴', icon: History, path: '/history' },
    { id: 'settings', label: '設定', icon: Settings, path: '/settings' },
  ];

  if (!isOpen) {
    return null;
  }

  return (
    <aside className="flex w-64 flex-col bg-white border-r border-gray-200 p-4 relative">
      {/* 閉じるボタン（モバイル用） */}
      <button
        onClick={onToggle}
        className="absolute top-4 right-4 lg:hidden p-2 rounded-lg hover:bg-gray-100 transition-colors"
        aria-label="サイドバーを閉じる"
      >
        <X size={20} />
      </button>

      <div className="flex flex-col gap-4">
        {/* ユーザー情報 */}
        <div className="flex items-center gap-3 pb-4 border-b border-gray-200">
          <div className="bg-gradient-to-br from-[#F4C753] to-[#FBBF24] rounded-full w-10 h-10 flex items-center justify-center text-white font-bold text-lg">
            {user?.name?.charAt(0) || 'U'}
          </div>
          <div className="flex flex-col">
            <h1 className="text-[#111827] text-base font-medium">
              {user?.name || 'ユーザー'}
            </h1>
            <p className="text-[#6B7280] text-sm">ようこそ</p>
          </div>
        </div>

        {/* ナビゲーション */}
        <nav className="flex flex-col gap-2">
          {menuItems.map((item) => {
            const Icon = item.icon;
            const isActive = currentPage === item.id;

            return (
              <button
                key={item.id}
                onClick={() => navigate(item.path)}
                className={`flex items-center gap-3 px-3 py-2 rounded-lg transition-colors ${isActive
                    ? 'bg-[#F4C753]/20 text-[#111827]'
                    : 'text-[#6B7280] hover:bg-gray-50'
                  }`}
              >
                <Icon size={20} />
                <p className="text-sm font-medium">{item.label}</p>
              </button>
            );
          })}
        </nav>
      </div>

      {/* 下部ボタン */}
      <div className="mt-auto flex flex-col gap-4">
        <button
          onClick={() => navigate('/review')}
          className="flex items-center justify-center gap-2 w-full h-10 px-4 bg-[#FBBF24] text-[#111827] rounded-lg font-bold text-sm hover:bg-amber-400 transition-colors"
        >
          <Plus size={20} />
          新規レビュー
        </button>

        <button
          onClick={handleLogout}
          className="flex items-center gap-3 px-3 py-2 rounded-lg text-[#6B7280] hover:bg-gray-50 transition-colors"
        >
          <LogOut size={20} />
          <p className="text-sm font-medium">ログアウト</p>
        </button>
      </div>
    </aside>
  );
}

// トグルボタンコンポーネント（メインコンテンツ側で使用）
export function SidebarToggle({ isOpen, onToggle }: { isOpen: boolean; onToggle: () => void }) {
  return (
    <button
      onClick={onToggle}
      className="p-2 bg-white rounded-lg border border-gray-200 hover:bg-gray-50 transition-colors shadow-sm"
      aria-label={isOpen ? 'サイドバーを閉じる' : 'サイドバーを開く'}
    >
      {isOpen ? <X size={20} /> : <Menu size={20} />}
    </button>
  );
}
