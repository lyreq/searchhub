
import React, { useState, useEffect } from 'react';
import { SearchParams, ContentType, SortOption } from '../../types';

interface SearchBarProps {
  searchParams: SearchParams;
  onParamsChange: (params: Partial<SearchParams>) => void;
  onFetchNow: () => void;
  isFetchingProviders: boolean;
}

export const SearchBar: React.FC<SearchBarProps> = ({ searchParams, onParamsChange, onFetchNow, isFetchingProviders }) => {
  const [query, setQuery] = useState(searchParams.query);

  useEffect(() => {
    setQuery(searchParams.query);
  }, [searchParams.query]);

  const handleQueryChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setQuery(e.target.value);
    onParamsChange({ query: e.target.value });
  };
  
  const handleTypeChange = (type: ContentType | '') => {
    onParamsChange({ type: searchParams.type === type ? '' : type });
  };
  
  return (
    <div className="p-4 bg-white rounded-lg shadow-md mb-6 space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 items-end">
        {/* Search Input */}
        <div className="lg:col-span-2">
          <label htmlFor="search" className="block text-sm font-medium text-gray-700">Search</label>
          <input
            id="search"
            type="text"
            value={query}
            onChange={handleQueryChange}
            placeholder="Search by keyword..."
            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-brand-blue focus:border-brand-blue"
          />
        </div>
        
        {/* Sort Select */}
        <div>
          <label htmlFor="sort" className="block text-sm font-medium text-gray-700">Sort By</label>
          <select
            id="sort"
            value={searchParams.sort}
            onChange={(e) => onParamsChange({ sort: e.target.value as SortOption })}
            className="mt-1 block w-full px-3 py-2 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-brand-blue focus:border-brand-blue"
          >
            {Object.values(SortOption).map(option => (
              <option key={option} value={option} className="capitalize">{option}</option>
            ))}
          </select>
        </div>
        
        {/* Per Page Select */}
        <div>
          <label htmlFor="per_page" className="block text-sm font-medium text-gray-700">Per Page</label>
          <select
            id="per_page"
            value={searchParams.per_page}
            onChange={(e) => onParamsChange({ per_page: parseInt(e.target.value, 10) })}
            className="mt-1 block w-full px-3 py-2 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-brand-blue focus:border-brand-blue"
          >
            {[10, 20, 50].map(val => (
              <option key={val} value={val}>{val}</option>
            ))}
          </select>
        </div>
      </div>
      
      {/* Type Filter & Fetch Button */}
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-2">
            <span className="text-sm font-medium text-gray-700">Type:</span>
            <button
                onClick={() => handleTypeChange(ContentType.VIDEO)}
                className={`px-3 py-1 text-sm font-semibold rounded-full ${searchParams.type === ContentType.VIDEO ? 'bg-brand-blue text-white' : 'bg-gray-200 text-gray-700'}`}
            >
                Video
            </button>
            <button
                onClick={() => handleTypeChange(ContentType.ARTICLE)}
                className={`px-3 py-1 text-sm font-semibold rounded-full ${searchParams.type === ContentType.ARTICLE ? 'bg-brand-blue text-white' : 'bg-gray-200 text-gray-700'}`}
            >
                Article
            </button>
        </div>
        <button
          onClick={onFetchNow}
          disabled={isFetchingProviders}
          className="px-4 py-2 bg-brand-blue text-white font-semibold rounded-md shadow-sm hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-blue disabled:bg-gray-400"
        >
          {isFetchingProviders ? 'Fetching...' : 'Fetch Now'}
        </button>
      </div>
      
      <p className="text-xs text-gray-500 text-center pt-2">List responses are cached server-side for ~60 seconds.</p>
    </div>
  );
};
