
# SearchHub Dashboard — Frontend Setup

> Minimal steps to run the React + TypeScript dashboard locally.

## Prerequisites

-   **Node.js 18+**
    
-   **npm** (or yarn/pnpm if you prefer)
    

## 1) Clone & enter the project
```bash
git clone <your-repo-url>
cd frontend
```
## 2) Create environment file

Create a `.env.local` at the project root with your API base URL:
```bash
# .env.local
VITE_API_BASE_URL=http://localhost:8082
```
## 3) Install dependencies
```bash
npm install
```
## 4) Start the dev server
```bash
npm run dev
```
Open the URL printed in your terminal (usually `http://localhost:5173`).

## Troubleshooting

-   **CORS / Network errors**: Confirm `VITE_API_BASE_URL` is correct and backend is reachable.
    
-   **Empty results**: Run a fetch cycle via **Fetch Now** in the UI or call `GET /api/v1/fetch-provider` on the backend first.