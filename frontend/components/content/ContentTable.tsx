
import React from 'react';
import { ListContentItem, Pagination as PaginationType, SearchParams } from '../../types';
import { formatDateTime } from '../../utils/dateFormatter';
import { Badge } from '../ui/Badge';

interface ContentTableProps {
  data?: ListContentItem[];
  pagination?: PaginationType;
  isLoading: boolean;
  isError: boolean;
  onRowClick: (id: number) => void;
  onPageChange: (page: number) => void;
  searchParams: SearchParams;
}

const TableSkeleton: React.FC = () => (
  <>
    {[...Array(5)].map((_, i) => (
      <tr key={i} className="animate-pulse">
        <td className="p-4"><div className="h-4 bg-gray-200 rounded w-3/4"></div></td>
        <td className="p-4"><div className="h-4 bg-gray-200 rounded w-1/2"></div></td>
        <td className="p-4"><div className="h-4 bg-gray-200 rounded w-1/4"></div></td>
        <td className="p-4"><div className="h-4 bg-gray-200 rounded w-1/2"></div></td>
        <td className="p-4"><div className="h-4 bg-gray-200 rounded w-3/4"></div></td>
        <td className="p-4"><div className="h-4 bg-gray-200 rounded w-full"></div></td>
      </tr>
    ))}
  </>
);

const Pagination: React.FC<{ pagination: PaginationType, onPageChange: (page: number) => void }> = ({ pagination, onPageChange }) => {
    const { total, page, per_page } = pagination;
    const totalPages = Math.ceil(total / per_page);

    if (totalPages <= 1) return null;

    return (
        <div className="flex items-center justify-between mt-4">
            <span className="text-sm text-gray-700">
                Page <span className="font-semibold">{page}</span> of <span className="font-semibold">{totalPages}</span>
            </span>
            <div className="inline-flex -space-x-px">
                <button onClick={() => onPageChange(page - 1)} disabled={page === 1} className="px-3 py-2 leading-tight text-gray-500 bg-white border border-gray-300 rounded-l-lg hover:bg-gray-100 disabled:opacity-50">
                    Previous
                </button>
                <button onClick={() => onPageChange(page + 1)} disabled={page === totalPages} className="px-3 py-2 leading-tight text-gray-500 bg-white border border-gray-300 rounded-r-lg hover:bg-gray-100 disabled:opacity-50">
                    Next
                </button>
            </div>
        </div>
    );
}

export const ContentTable: React.FC<ContentTableProps> = ({ data, pagination, isLoading, isError, onRowClick, onPageChange }) => {
  if (isError) {
    return <div className="text-center p-8 bg-red-50 text-red-700 rounded-lg">Failed to load content. Please try again.</div>;
  }
  
  return (
    <div className="bg-white shadow-md rounded-lg overflow-hidden">
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50 sticky top-0">
            <tr>
              <th className="p-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Title</th>
              <th className="p-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
              <th className="p-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Score</th>
              <th className="p-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Provider</th>
              <th className="p-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Published At</th>
              <th className="p-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tags</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {isLoading ? <TableSkeleton /> : (
              data?.map(item => (
                <tr key={item.id} onClick={() => onRowClick(item.id)} className="hover:bg-gray-50 cursor-pointer">
                  <td className="p-4 whitespace-nowrap text-sm font-medium text-gray-900">{item.title}</td>
                  <td className="p-4 whitespace-nowrap text-sm text-gray-500 capitalize">{item.type}</td>
                  <td className="p-4 whitespace-nowrap text-sm text-gray-500">{item.score.toFixed(2)}</td>
                  <td className="p-4 whitespace-nowrap text-sm text-gray-500">{item.provider}</td>
                  <td className="p-4 whitespace-nowrap text-sm text-gray-500" title={item.published_at}>
                    {formatDateTime(item.published_at)}
                  </td>
                  <td className="p-4 whitespace-nowrap">
                    {item.tags.map(tag => <Badge key={tag}>{tag}</Badge>)}
                  </td>
                </tr>
              ))
            )}
            {!isLoading && data?.length === 0 && (
                <tr>
                    <td colSpan={6} className="text-center p-8 text-gray-500">No content found.</td>
                </tr>
            )}
          </tbody>
        </table>
      </div>
      {pagination && <div className="p-4"><Pagination pagination={pagination} onPageChange={onPageChange} /></div>}
    </div>
  );
};
