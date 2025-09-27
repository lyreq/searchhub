
import React from 'react';
import { FetchResult, ProviderRunResult } from '../types';

interface FetchProvidersResultProps {
  data: FetchResult;
}

const ProviderRow: React.FC<{ providerResult: ProviderRunResult }> = ({ providerResult }) => {
    const hasErrors = providerResult.errors && providerResult.errors.length > 0;
    const statusColor = hasErrors ? 'bg-yellow-100 text-yellow-800' : 'bg-green-100 text-green-800';
    const statusText = hasErrors ? 'Degraded' : 'OK';

    return (
        <tr className={hasErrors ? "bg-yellow-50" : ""}>
            <td className="p-3 text-sm font-medium text-gray-900">{providerResult.provider}</td>
            <td className="p-3 text-sm text-gray-700">{providerResult.fetched}</td>
            <td className="p-3 text-sm text-gray-700">{providerResult.inserted}</td>
            <td className="p-3 text-sm text-gray-700">{providerResult.updated}</td>
            <td className="p-3 text-sm text-gray-700">{providerResult.skipped}</td>
            <td className="p-3">
                <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${statusColor}`}>
                    {statusText}
                </span>
            </td>
        </tr>
    );
};


export const FetchProvidersResult: React.FC<FetchProvidersResultProps> = ({ data }) => {
    const anyProviderHasErrors = data.providers.some(p => p.errors && p.errors.length > 0);
    return (
        <div className="space-y-6">
            {anyProviderHasErrors && (
                <div className="p-4 bg-yellow-50 border-l-4 border-yellow-400">
                    <div className="flex">
                        <div className="flex-shrink-0">
                             <svg className="h-5 w-5 text-yellow-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                                <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.21 3.03-1.742 3.03H4.42c-1.532 0-2.492-1.696-1.742-3.03l5.58-9.92zM10 13a1 1 0 110-2 1 1 0 010 2zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                            </svg>
                        </div>
                        <div className="ml-3">
                            <p className="text-sm text-yellow-700">
                                One or more providers encountered errors.
                            </p>
                        </div>
                    </div>
                </div>
            )}
            
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
                <div><span className="block text-2xl font-bold text-brand-dark">{data.inserted}</span><span className="text-sm text-gray-500">Inserted</span></div>
                <div><span className="block text-2xl font-bold text-brand-dark">{data.updated}</span><span className="text-sm text-gray-500">Updated</span></div>
                <div><span className="block text-2xl font-bold text-brand-dark">{data.skipped}</span><span className="text-sm text-gray-500">Skipped</span></div>
                <div><span className="block text-2xl font-bold text-brand-dark">{data.total}</span><span className="text-sm text-gray-500">Total</span></div>
            </div>

            <div>
                <h3 className="text-lg font-semibold text-gray-800 mb-2">Provider Breakdown</h3>
                <div className="overflow-x-auto border rounded-lg">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="p-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Provider</th>
                                <th className="p-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Fetched</th>
                                <th className="p-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Inserted</th>
                                <th className="p-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Updated</th>
                                <th className="p-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Skipped</th>
                                <th className="p-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {data.providers.map(p => <ProviderRow key={p.provider} providerResult={p} />)}
                        </tbody>
                    </table>
                </div>
            </div>

            <div className="text-xs text-gray-500">
                Run ID: <span className="font-mono">{data.run_id}</span>
            </div>
        </div>
    );
};
