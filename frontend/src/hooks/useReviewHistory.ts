import { useState, useEffect, useCallback } from 'react';
import {
  ReviewHistoryItem,
  ReviewHistoryListResponse,
  ReviewHistoryFilter,
} from '../types/review';

export const useReviewHistory = (initialFilter?: ReviewHistoryFilter) => {
  const [items, setItems] = useState<ReviewHistoryItem[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);
  const [filter, setFilter] = useState<ReviewHistoryFilter>(
    initialFilter || {
      page: 1,
      pageSize: 10,
      sortBy: 'createdAt',
      sortOrder: 'desc',
    }
  );
  const [totalItems, setTotalItems] = useState<number>(0);
  const [totalPages, setTotalPages] = useState<number>(0);

  const fetchReviewHistory = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      // TODO: 実際のAPIコールに置き換える
      // const response = await fetch('/api/reviews/history', {
      //   method: 'POST',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify(filter),
      // });
      // const data: ReviewHistoryListResponse = await response.json();

      // モックデータ（ローディング遅延をシミュレート）
      await new Promise((resolve) => setTimeout(resolve, 800));

      const mockData: ReviewHistoryListResponse = {
        items: [
          {
            id: '1',
            createdAt: '2024-01-27T15:30:00Z',
            language: 'TypeScript',
            status: 'warning',
          },
          {
            id: '2',
            createdAt: '2024-01-27T14:20:00Z',
            language: 'Python',
            status: 'success',
          },
          {
            id: '3',
            createdAt: '2024-01-27T11:15:00Z',
            language: 'JavaScript',
            status: 'success',
          },
          {
            id: '4',
            createdAt: '2024-01-26T18:15:00Z',
            language: 'Go',
            status: 'error',
          },
          {
            id: '5',
            createdAt: '2024-01-26T16:45:00Z',
            language: 'TypeScript',
            status: 'warning',
          },
          {
            id: '6',
            createdAt: '2024-01-26T09:45:00Z',
            language: 'Python',
            status: 'success',
          },
          {
            id: '7',
            createdAt: '2024-01-25T17:30:00Z',
            language: 'Java',
            status: 'success',
          },
          {
            id: '8',
            createdAt: '2024-01-25T11:02:00Z',
            language: 'JavaScript',
            status: 'error',
          },
          {
            id: '9',
            createdAt: '2024-01-24T14:20:00Z',
            language: 'TypeScript',
            status: 'success',
          },
          {
            id: '10',
            createdAt: '2024-01-24T10:10:00Z',
            language: 'Python',
            status: 'warning',
          },
        ],
        total: 53,
        page: filter.page || 1,
        pageSize: filter.pageSize || 10,
        totalPages: Math.ceil(53 / (filter.pageSize || 10)),
      };

      let filteredItems = mockData.items;

      if (filter.language) {
        filteredItems = filteredItems.filter((item) => item.language === filter.language);
      }

      if (filter.status) {
        filteredItems = filteredItems.filter((item) => item.status === filter.status);
      }

      setItems(filteredItems);
      setTotalItems(mockData.total);
      setTotalPages(mockData.totalPages);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('データの取得に失敗しました'));
    } finally {
      setLoading(false);
    }
  }, [filter]);

  useEffect(() => {
    fetchReviewHistory();
  }, [fetchReviewHistory]);

  const updateFilter = useCallback((newFilter: Partial<ReviewHistoryFilter>) => {
    setFilter((prev) => ({
      ...prev,
      ...newFilter,
      page: newFilter.page !== undefined ? newFilter.page : 1,
    }));
  }, []);

  const changePage = useCallback((page: number) => {
    setFilter((prev) => ({ ...prev, page }));
  }, []);

  const changeSort = useCallback((sortBy: string) => {
    setFilter((prev) => ({
      ...prev,
      sortBy: sortBy as 'createdAt' | 'language' | 'status',
      sortOrder: prev.sortBy === sortBy && prev.sortOrder === 'asc' ? 'desc' : 'asc',
    }));
  }, []);

  const refetch = useCallback(() => {
    fetchReviewHistory();
  }, [fetchReviewHistory]);

  return {
    items,
    loading,
    error,
    filter,
    totalItems,
    totalPages,
    updateFilter,
    changePage,
    changeSort,
    refetch,
  };
};
