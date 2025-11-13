import { useState, useEffect, useCallback } from 'react';
import { KnowledgeListItem, KnowledgeFilter } from '../types/knowledge';
import { knowledgeApiClient, Knowledge } from '../api/knowledgeApi';

export const useKnowledgeList = (initialFilter?: Partial<KnowledgeFilter>) => {
  const [items, setItems] = useState<KnowledgeListItem[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);
  const [filter, setFilter] = useState<KnowledgeFilter>({
    page: 1,
    pageSize: 10,
    sortBy: 'created_at',
    sortOrder: 'desc',
    ...initialFilter,
  });

  // クライアントサイドでのソート処理
  const sortItems = useCallback((data: Knowledge[]): KnowledgeListItem[] => {
    const sorted = [...data].sort((a, b) => {
      let compareValue = 0;

      switch (filter.sortBy) {
        case 'created_at':
          compareValue = new Date(a.created_at).getTime() - new Date(b.created_at).getTime();
          break;
        case 'priority':
          compareValue = a.priority - b.priority;
          break;
        case 'usage_count':
          compareValue = a.usage_count - b.usage_count;
          break;
      }

      return filter.sortOrder === 'asc' ? compareValue : -compareValue;
    });

    return sorted.map(item => ({
      id: item.id,
      title: item.title,
      content: item.content,
      category: item.category,
      priority: item.priority,
      usage_count: item.usage_count,
      last_used_at: item.last_used_at,
      created_at: item.created_at,
      updated_at: item.updated_at,
    }));
  }, [filter.sortBy, filter.sortOrder]);

  // クライアントサイドでのページング処理
  const paginateItems = useCallback((data: KnowledgeListItem[]): KnowledgeListItem[] => {
    const startIndex = (filter.page - 1) * filter.pageSize;
    const endIndex = startIndex + filter.pageSize;
    return data.slice(startIndex, endIndex);
  }, [filter.page, filter.pageSize]);

  const fetchKnowledge = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      // APIからナレッジ一覧を取得
      const data = await knowledgeApiClient.listKnowledge(filter.category);

      // ソート
      const sortedData = sortItems(data);

      // 優先度フィルタ（クライアントサイド）
      const filteredData = filter.priority
        ? sortedData.filter(item => item.priority === filter.priority)
        : sortedData;

      // 全体の件数を保存
      const totalItems = filteredData.length;
      const totalPages = Math.ceil(totalItems / filter.pageSize);

      // ページング
      const paginatedData = paginateItems(filteredData);

      setItems(paginatedData);
      // totalとpagesを返すために、stateに保存
      (fetchKnowledge as any).totalItems = totalItems;
      (fetchKnowledge as any).totalPages = totalPages;
    } catch (err) {
      console.error('Failed to fetch knowledge:', err);
      setError(err instanceof Error ? err : new Error('データの取得に失敗しました'));
    } finally {
      setLoading(false);
    }
  }, [filter, sortItems, paginateItems]);

  useEffect(() => {
    fetchKnowledge();
  }, [fetchKnowledge]);

  const updateFilter = useCallback((newFilter: Partial<KnowledgeFilter>) => {
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
      sortBy: sortBy as 'created_at' | 'priority' | 'usage_count',
      sortOrder: prev.sortBy === sortBy && prev.sortOrder === 'asc' ? 'desc' : 'asc',
    }));
  }, []);

  const refetch = useCallback(() => {
    fetchKnowledge();
  }, [fetchKnowledge]);

  // 総件数と総ページ数を取得（クロージャから）
  const totalItems = (fetchKnowledge as any).totalItems || 0;
  const totalPages = (fetchKnowledge as any).totalPages || 0;

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
