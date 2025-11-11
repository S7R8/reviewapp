import React, { useState, useEffect, useRef } from 'react';
import { ChevronDown } from 'lucide-react';
import { ProgrammingLanguage } from '../types/review';

interface SearchFilterProps {
  onLanguageChange: (language: ProgrammingLanguage | undefined) => void;
}

const SearchFilter: React.FC<SearchFilterProps> = ({
  onLanguageChange,
}) => {
  const [showLanguageMenu, setShowLanguageMenu] = useState(false);
  const [selectedLanguage, setSelectedLanguage] = useState<ProgrammingLanguage | undefined>();

  const languageRef = useRef<HTMLDivElement>(null);

  const languages: (ProgrammingLanguage | 'all')[] = [
    'all', 'TypeScript', 'JavaScript', 'Python', 'Go', 'Java', 'C++', 'C#', 
    'Ruby', 'PHP', 'Rust', 'Swift', 'Kotlin', 'Other',
  ];

  // 外側クリックでメニューを閉じる
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (languageRef.current && !languageRef.current.contains(event.target as Node)) {
        setShowLanguageMenu(false);
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

  return (
    <div className="flex gap-3 items-center">
      {/* 言語フィルター */}
      <div className="relative" ref={languageRef}>
        <button
          onClick={() => setShowLanguageMenu(!showLanguageMenu)}
          className="flex h-10 items-center gap-x-2 rounded-lg bg-white border border-gray-300 px-4 hover:bg-gray-50 transition-colors"
        >
          <p className="text-sm font-medium text-[#111827]">{selectedLanguage || '言語でフィルター'}</p>
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
    </div>
  );
};

export default SearchFilter;
