
import React, { useState, useCallback } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getContents, fetchProviders } from './services/apiClient';
import { useSearchQuery } from './hooks/useSearchQuery';
import { useDebounce } from './hooks/useDebounce';
import { ContentTable } from './components/content/ContentTable';
import { SearchBar } from './components/filters/SearchBar';
import { ContentDetailDrawer } from './components/content/ContentDetailDrawer';
import { ToastContainer } from './components/ui/Toast';
import { Modal } from './components/ui/Modal';
import { FetchProvidersResult } from './components/FetchProvidersResult';
import { ToastMessage, FetchResult, SearchParams } from './types';

const App: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchQuery();
  const debouncedQuery = useDebounce(searchParams.query, 400);

  const queryClient = useQueryClient();

  const [selectedContentId, setSelectedContentId] = useState<number | null>(null);
  const [toasts, setToasts] = useState<ToastMessage[]>([]);
  const [fetchResult, setFetchResult] = useState<FetchResult | null>(null);

  const addToast = (message: string, type: 'success' | 'error') => {
    setToasts((prev) => [...prev, { id: Date.now(), message, type }]);
  };

  const removeToast = (id: number) => {
    setToasts((prev) => prev.filter((toast) => toast.id !== id));
  };

  const { data: contentsData, isLoading, isError } = useQuery({
    queryKey: ['contents', debouncedQuery, searchParams.type, searchParams.sort, searchParams.page, searchParams.per_page],
    queryFn: () => getContents({
      ...searchParams,
      query: debouncedQuery,
    }),
  });

  const fetchProvidersMutation = useMutation({
    mutationFn: fetchProviders,
    onSuccess: (data) => {
      addToast(data.message, 'success');
      setFetchResult(data.data);
      queryClient.invalidateQueries({ queryKey: ['contents'] });
    },
    onError: (error: Error) => {
      addToast(error.message, 'error');
    },
  });
  
  const handleParamsChange = useCallback((params: Partial<SearchParams>) => {
    setSearchParams(params);
  }, [setSearchParams]);

  const handlePageChange = (page: number) => {
    setSearchParams({ page });
  };
  
  return (
    <div className="min-h-screen bg-gray-50 text-brand-dark">
      <ToastContainer toasts={toasts} onDismiss={removeToast} />
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto py-4 px-4 sm:px-6 lg:px-8">
            <h1 className="text-3xl font-bold leading-tight text-brand-dark">
                SearchHub
            </h1>
            <p className="text-sm text-gray-500">Content Search Service Dashboard</p>
        </div>
      </header>
      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <SearchBar
          searchParams={searchParams}
          onParamsChange={handleParamsChange}
          onFetchNow={() => fetchProvidersMutation.mutate()}
          isFetchingProviders={fetchProvidersMutation.isPending}
        />
        <ContentTable
          data={contentsData?.data.data}
          pagination={contentsData?.data.pagination}
          isLoading={isLoading}
          isError={isError}
          onRowClick={setSelectedContentId}
          onPageChange={handlePageChange}
          searchParams={searchParams}
        />
      </main>
      <ContentDetailDrawer
        contentId={selectedContentId}
        onClose={() => setSelectedContentId(null)}
      />
      <Modal 
        isOpen={!!fetchResult} 
        onClose={() => setFetchResult(null)}
        title="Provider Fetch Summary"
      >
        {fetchResult && <FetchProvidersResult data={fetchResult} />}
      </Modal>
    </div>
  );
};

export default App;
