import { ApiSuccess, ListEnvelope, ContentDetail, FetchResult, SearchParams } from '../types';

// Fix: The triple-slash directive for Vite client types was causing an error because the
// type definition file could not be found. To resolve the error on `import.meta.env`
// without these types, a type assertion `(import.meta as any)` is used.
const BASE_URL = (import.meta as any).env?.VITE_API_BASE_URL;
const API_BASE_URL = `${BASE_URL}/api/v1`;

async function handleResponse<T,>(response: Response): Promise<ApiSuccess<T>> {
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.message || 'An API error occurred.');
  }
  return data;
}

export async function getContents(params: Partial<SearchParams>): Promise<ApiSuccess<ListEnvelope>> {
  const urlParams = new URLSearchParams();
  if (params.query) urlParams.append('query', params.query);
  if (params.type) urlParams.append('type', params.type);
  if (params.sort) urlParams.append('sort', params.sort);
  if (params.page) urlParams.append('page', params.page.toString());
  if (params.per_page) urlParams.append('per_page', params.per_page.toString());

  const response = await fetch(`${API_BASE_URL}/contents?${urlParams.toString()}`);
  return handleResponse<ListEnvelope>(response);
}

export async function getContentDetail(id: number, includeRaw: boolean = true): Promise<ApiSuccess<ContentDetail>> {
  const urlParams = new URLSearchParams({ include_raw: includeRaw ? 'true' : 'false' });
  const response = await fetch(`${API_BASE_URL}/contents/${id}?${urlParams.toString()}`);
  return handleResponse<ContentDetail>(response);
}

export async function fetchProviders(): Promise<ApiSuccess<FetchResult>> {
  const response = await fetch(`${API_BASE_URL}/fetch-provider`);
  return handleResponse<FetchResult>(response);
}