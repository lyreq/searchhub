# SearchHub — Content Search & Scoring (Monorepo)

A small, production-ready demo that **fetches content from multiple providers (JSON/XML)**, **normalizes & scores** items, **stores** them in Postgres, **caches** list responses in Redis, and serves a **React + TypeScript dashboard**.

----------

## Table of Contents

-   [Project Structure](#project-structure)
    
-   [High-Level Architecture](#high-level-architecture)
    
-   [Scoring Model](#scoring-model)
    
-   [APIs (Summary)](#apis-summary)
    
-   [API Docs (Swagger)](#api-docs-swagger)
    
-   [Backend (Overview)](#backend-overview)
    
-   [Frontend (Overview)](#frontend-overview)
    
-   [Local Development (Quick Start)](#local-development-quick-start)
    
-   [Operational Notes](#operational-notes)


## Project Structure
```bash
.
├── backend/              # Golang (Echo) API: providers, scoring, storage, cache
│   └── README.md         # Backend-specific setup & details
├── frontend/             # React + TypeScript dashboard (Vite)
│   └── README.md         # Frontend-specific setup & details
└── README.md             # This file (root overview)
```

## High-Level Architecture
```bash
+-------------------+          +--------------------+
|  Provider: JSON   |          |  Provider: XML     |
|  (mock endpoint)  |          |  (mock endpoint)   |
+---------+---------+          +----------+---------+
          \                            /
           \   Adapters (pkg/provider)/
            \      Normalize         /
             +----------+-----------+
                        |
                        v
                Scoring Engine
         (base + typeCoef + freshness + engagement)
                        |
                        v
                 Persistence (DB)
              (Postgres via GORM models)
                        ^
                        |
        List API <-- Cache (Redis, short TTL)
                        |
                        v
                  Frontend Dashboard
```

**Key goals**

-   Merge heterogeneous provider formats into a **single, standard content model**.
    
-   **Score** items for ranking (relevance/popularity/freshness).
    
-   Provide a **clean, queryable API** + **simple dashboard**.

## Scoring Model

> **FinalScore = (BaseScore × TypeCoefficient) + Freshness + Engagement**

-   **BaseScore**
    
    -   Video: `views / 1000 + likes / 100`
        
    -   Article: `reading_time + reactions / 50`
        
-   **TypeCoefficient**
    
    -   Video: `1.5`
        
    -   Article: `1.0`
        
-   **Freshness (published_at)**
    
    -   within 1 week: `+5`
        
    -   within 1 month: `+3`
        
    -   within 3 months: `+1`
        
    -   older: `+0`
        
-   **Engagement**
    
    -   Video: `(likes / views) * 10`
        
    -   Article: `(reactions / reading_time) * 5`
        

Scoring is **computed server-side** when ingesting content and persisted for fast sorting.

----------

## APIs (Summary)

**Base URL**: `http://localhost:8082/api/v1`

1.  **Fetch Providers (ingest & score)**
    
    -   `GET /fetch-provider`
        
    -   Pulls from both providers, normalizes, scores, upserts into DB.
        
    -   Returns a run summary per provider (fetched/inserted/updated/errors).
        
2.  **List Contents (search/filter/sort/pagination)**
    
    -   `GET /contents`
        
    -   **Query params** (all optional):
        
        -   `query`: keyword
            
        -   `type`: `video` | `article`
            
        -   `sort`: `score` | `relevance` | `recent` | `popularity` (default: `score`)
            
        -   `page`: integer (default: `1`)
            
        -   `per_page`: integer 1..100 (default: `20`)
            
    -   Response includes `data[]`, `pagination`, effective `sort` and `filters`.
        
3.  **Content Detail**
    
    -   `GET /contents/{id}`
        
    -   **Query**:
        
        -   `include_raw`: `true|1` to include provider raw payload
            
    -   Returns metrics and score breakdown.
        

> List responses are **cached in Redis** per `(query,type,sort,page,per_page)` with a **short TTL (~60s)**.

----------

## API Docs (Swagger)

When the backend is running, open:

**`http://localhost:8082/api-docs`**

You’ll see interactive **Swagger/OpenAPI** documentation for all endpoints and models.

----------

## Backend (Overview)

-   **Language/Framework**: Go (Echo)
    
-   **Data Store**: Postgres (GORM)
    
-   **Cache**: Redis (short TTL for list results)
    
-   **Providers**: Two mock providers (JSON, XML) via dedicated **adapters** under `pkg/provider`.
    
-   **Rate Limiting & Reliability**:
    
    -   Tuned HTTP client (timeouts)
        
    -   Retries with backoff (for provider errors)
        
    -   In-memory rate limit (protect provider quotas)
        
-   **Extensibility**:
    
    -   New provider = implement `Provider` interface (normalize to the shared `Item`), register in a factory.
        
-   **Docs**:
    
    -   See **`backend/README.md`** for environment variables, Docker Compose usage, and local run instructions.
        

----------

## Frontend (Overview)

-   **Stack**: React + TypeScript (Vite)
    
-   **Data fetching**: TanStack Query (react-query)
    
-   **UI**:
    
    -   **Contents** page: search, type filter, sort, pagination
        
    -   **Table** with columns: Title, Type, Score, Provider, PublishedAt, Tags
        
    -   Row click → **Detail drawer** with metrics & score breakdown (and optional raw JSON)
        
    -   **“Fetch Now”** button to trigger ingestion; shows summarized result per provider
        
-   **Config**:
    
    -   `.env.local` → `VITE_API_BASE_URL=http://localhost:8082`
        
-   **Docs**:
    
    -   See **`frontend/README.md`** for setup, scripts, and dev workflow.
        

----------

## Local Development (Quick Start)

### 1) Start the backend

-   Follow **`backend/README.md`**.
    
-   Ensure API is reachable at **`http://localhost:8082/api/v1/...`**.
    
-   Verify Swagger at **`/api-docs`**.
    

### 2) Start the frontend

-   Follow **`frontend/README.md`**.
    
-   Create `.env.local` with:
```bash
VITE_API_BASE_URL=http://localhost:8082
```
-   Run `npm install` then `npm run dev`.
    
-   Open the printed dev URL (e.g., `http://localhost:5173`).

### 3) Ingest some data

-   Click **“Fetch Now”** in the UI or call:
```bash
GET http://localhost:8082/api/v1/fetch-provider
```

-   Then use the **Contents** page to search/sort/browse items.
    

----------

## Operational Notes

-   **Caching**: List results cached (Redis) ~60s to reduce DB load and improve snappiness.
    
-   **Resilience**: Provider calls use timeouts, limited retries, and in-memory rate limiting.
    
-   **Consistency**: Items are **upserted** (new → insert, existing → update scores/fields).
    
-   **Observability**: The boilerplate includes hooks for logging/metrics; see backend docs for details.