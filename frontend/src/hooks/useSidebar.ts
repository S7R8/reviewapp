import { useState, useEffect } from 'react';

export function useSidebar() {
  const [isOpen, setIsOpen] = useState(true);

  useEffect(() => {
    // 画面幅の変更を監視
    const handleResize = () => {
      const isMobile = window.innerWidth < 1024; // lg breakpoint
      
      if (isMobile) {
        setIsOpen(false); // モバイルでは閉じる
      } else {
        setIsOpen(true); // デスクトップでは開く
      }
    };

    // 初回実行
    handleResize();

    // リスナー登録
    window.addEventListener('resize', handleResize);

    // クリーンアップ
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  const toggle = () => setIsOpen(!isOpen);

  return { isOpen, toggle };
}
