import { useState, useEffect, useCallback } from 'react';
import {
  ReviewHistoryItem,
  ReviewHistoryListResponse,
  ReviewHistoryFilter,
} from '../types/review';
import { reviewApiClient } from '../api/reviewApi';

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
      // 実際のAPIを呼び出す
      const data: ReviewHistoryListResponse = await reviewApiClient.getReviewHistory(filter);

      setItems(data.items);
      setTotalItems(data.total);
      setTotalPages(data.totalPages);
    } catch (err) {
      console.error('Failed to fetch review history:', err);
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
