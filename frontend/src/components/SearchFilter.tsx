import React, { useState, useEffect, useRef } from 'react';
import { ChevronDown } from 'lucide-react';
import { ProgrammingLanguage, ReviewStatus } from '../types/review';

interface SearchFilterProps {
  onLanguageChange: (language: ProgrammingLanguage | undefined) => void;
  onStatusChange: (status: ReviewStatus | undefined) => void;
}

const SearchFilter: React.FC<SearchFilterProps> = ({
  onLanguageChange,
  onStatusChange,
}) => {
  const [showLanguageMenu, setShowLanguageMenu] = useState(false);
  const [showStatusMenu, setShowStatusMenu] = useState(false);
  const [selectedLanguage, setSelectedLanguage] = useState<ProgrammingLanguage | undefined>();
  const [selectedStatus, setSelectedStatus] = useState<ReviewStatus | undefined>();

  const languageRef = useRef<HTMLDivElement>(null);
  const statusRef = useRef<HTMLDivElement>(null);

  const languages: (ProgrammingLanguage | 'all')[] = [
    'all', 'TypeScript', 'JavaScript', 'Python', 'Go', 'Java',
  ];

  const statuses: ({ value: ReviewStatus | 'all'; label: string })[] = [
    { value: 'all', label: 'すべて' },
    { value: 'success', label: '成功' },
    { value: 'warning', label: '改善点あり' },
    { value: 'error', label: 'エラー' },
  ];

  // 外側クリックでメニューを閉じる
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (languageRef.current && !languageRef.current.contains(event.target as Node)) {
        setShowLanguageMenu(false);
      }
      if (statusRef.current && !statusRef.current.contains(event.target as Node)) {
        setShowStatusMenu(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  const handleLanguageSelect = (language: ProgrammingLanguage | 'all') => {
    const selectedLang = language === 'all' ? undefined : language;
    setSelectedLanguage(selectedLang);
    onLanguageChange(selectedLang);
    setShowLanguageMenu(false);
  };

  const handleStatusSelect = (status: ReviewStatus | 'all') => {
    const selectedStat = status === 'all' ? undefined : status;
    setSelectedStatus(selectedStat);
    onStatusChange(selectedStat);
    setShowStatusMenu(false);
  };

  return (
    <div className="flex gap-3 items-center">
      {/* 言語フィルター */}
      <div className="relative" ref={languageRef}>
        <button
          onClick={() => setShowLanguageMenu(!showLanguageMenu)}
          className="flex h-10 items-center gap-x-2 rounded-lg bg-white border border-gray-300 px-4 hover:bg-gray-50 transition-colors"
        >
          <p className="text-sm font-medium text-[#111827]">{selectedLanguage || '言語'}</p>
          <ChevronDown className="w-4 h-4 text-[#6B7280]" />
        </button>
        {showLanguageMenu && (
          <div className="absolute z-10 mt-2 w-48 rounded-lg bg-white border border-gray-200 shadow-lg">
            <div className="py-1 max-h-60 overflow-y-auto">
              {languages.map((lang) => (
                <button
                  key={lang}
                  onClick={() => handleLanguageSelect(lang as ProgrammingLanguage | 'all')}
                  className="block w-full text-left px-4 py-2 text-sm text-[#111827] hover:bg-gray-50 transition-colors"
                >
                  {lang === 'all' ? 'すべて' : lang}
                </button>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* ステータスフィルター */}
      <div className="relative" ref={statusRef}>
        <button
          onClick={() => setShowStatusMenu(!showStatusMenu)}
          className="flex h-10 items-center gap-x-2 rounded-lg bg-white border border-gray-300 px-4 hover:bg-gray-50 transition-colors"
        >
          <p className="text-sm font-medium text-[#111827]">
            {selectedStatus ? statuses.find((s) => s.value === selectedStatus)?.label : 'ステータス'}
          </p>
          <ChevronDown className="w-4 h-4 text-[#6B7280]" />
        </button>
        {showStatusMenu && (
          <div className="absolute z-10 mt-2 w-48 rounded-lg bg-white border border-gray-200 shadow-lg">
            <div className="py-1">
              {statuses.map((status) => (
                <button
                  key={status.value}
                  onClick={() => handleStatusSelect(status.value as ReviewStatus | 'all')}
                  className="block w-full text-left px-4 py-2 text-sm text-[#111827] hover:bg-gray-50 transition-colors"
                >
                  {status.label}
                </button>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default SearchFilter;
