
export interface ApiSuccess<T> {
  status: "success" | "error";
  data: T;
  message: string;
}

export enum ContentType {
  VIDEO = 'video',
  ARTICLE = 'article',
}

export enum SortOption {
  SCORE = 'score',
  RELEVANCE = 'relevance',
  RECENT = 'recent',
  POPULARITY = 'popularity',
}

export interface ListContentItem {
  id: number;
  title: string;
  type: ContentType;
  score: number;
  published_at: string;
  tags: string[];
  provider: string;
}

export interface Pagination {
  total: number;
  page: number;
  per_page: number;
}

export interface ListEnvelope {
  data: ListContentItem[];
  pagination: Pagination;
  sort: string;
  filters: Record<string, any>;
}

export interface ContentMetrics {
  views?: number;
  likes?: number;
  duration_seconds?: number;
  reading_time_minutes?: number;
  reactions?: number;
  comments?: number;
}

export interface ContentScores {
  base: number;
  type_coefficient: number;
  freshness: number;
  engagement: number;
  final: number;
}

export interface ContentDetail {
  id: number;
  title: string;
  type: ContentType;
  provider: string;
  published_at: string;
  tags: string[];
  metrics: ContentMetrics;
  scores: ContentScores;
  raw_payload?: Record<string, any>;
}

export interface ProviderRunResult {
  provider: string;
  fetched: number;
  inserted: number;
  updated: number;
  skipped: number;
  errors: string[] | null;
}

export interface FetchResult {
  run_id: string;
  total: number;
  inserted: number;
  updated: number;
  skipped: number;
  providers: ProviderRunResult[];
}

export interface SearchParams {
  query: string;
  type: ContentType | '';
  sort: SortOption;
  page: number;
  per_page: number;
}

export type ToastMessage = {
  id: number;
  message: string;
  type: 'success' | 'error';
};
