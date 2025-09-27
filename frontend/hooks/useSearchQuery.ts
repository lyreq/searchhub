
import { useState, useCallback, useMemo } from 'react';
import { SearchParams, ContentType, SortOption } from '../types';

const getParamsFromUrl = (): Partial<SearchParams> => {
  const params = new URLSearchParams(window.location.search);
  return {
    query: params.get('query') || '',
    type: (params.get('type') as ContentType) || '',
    sort: (params.get('sort') as SortOption) || SortOption.SCORE,
    page: parseInt(params.get('page') || '1', 10),
    per_page: parseInt(params.get('per_page') || '10', 10),
  };
};

export function useSearchQuery(): [SearchParams, (newParams: Partial<SearchParams>) => void] {
  const [searchParams, setSearchParams] = useState<SearchParams>(() => {
    const params = getParamsFromUrl();
    return {
      query: params.query || '',
      type: params.type || '',
      sort: params.sort || SortOption.SCORE,
      page: params.page || 1,
      per_page: params.per_page || 10,
    };
  });

  const updateSearchParams = useCallback((newParams: Partial<SearchParams>) => {
    const oldParams = getParamsFromUrl();
    const updated: SearchParams = {
      query: 'query' in newParams ? newParams.query! : oldParams.query || '',
      type: 'type' in newParams ? newParams.type! : (oldParams.type || ''),
      sort: 'sort' in newParams ? newParams.sort! : oldParams.sort || SortOption.SCORE,
      page: 'page' in newParams ? newParams.page! : oldParams.page || 1,
      per_page: 'per_page' in newParams ? newParams.per_page! : oldParams.per_page || 10,
    };

    // Reset page to 1 if filters change
    const hasFilterChanged =
      ('query' in newParams && newParams.query !== oldParams.query) ||
      ('type' in newParams && newParams.type !== oldParams.type) ||
      ('sort' in newParams && newParams.sort !== oldParams.sort) ||
      ('per_page' in newParams && newParams.per_page !== oldParams.per_page);

    if (hasFilterChanged) {
      updated.page = 1;
    }
    
    setSearchParams(updated);

    const urlParams = new URLSearchParams();
    if (updated.query) urlParams.set('query', updated.query);
    if (updated.type) urlParams.set('type', updated.type);
    if (updated.sort) urlParams.set('sort', updated.sort);
    if (updated.page) urlParams.set('page', updated.page.toString());
    if (updated.per_page) urlParams.set('per_page', updated.per_page.toString());
    
    const newUrl = `${window.location.pathname}?${urlParams.toString()}`;
    window.history.pushState({}, '', newUrl);
  }, []);
  
  const memoizedSearchParams = useMemo(() => searchParams, [searchParams]);

  return [memoizedSearchParams, updateSearchParams];
}
