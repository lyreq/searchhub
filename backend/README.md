
# SearchHub Backend — Setup & Run

This project is a **Golang (Echo)** backend that ingests content from two providers (JSON/XML), normalizes & scores them, stores in **Postgres**, caches list responses in **Redis**, and exposes a small API.


## Prerequisites

-   **Docker & Docker Compose v2** (recommended path)
    
-   (Optional for host-run) **Go 1.22+**, **Make**

## 1) Environment

Create a `.env` at the repo root. Minimal vars for this project:
```bash
#### DOCKER HOST PORT MAPPINGS ####

APP_HOST_PORT="8080"
APP_GREEN_HOST_PORT="8005"
APP_BLUE_HOST_PORT="8006"
POSTGRES_HOST_PORT="5433"
MEILISEARCH_HOST_PORT="7701"
REDIS_HOST_PORT="6380"
REDIS_EXPORTER_HOST_PORT="9122"
PROMETHEUS_HOST_PORT="9091"
GRAFANA_HOST_PORT="3001"

#### APP ENVIRONMENT ####

APP_URL="http://localhost:8082"
APP_PORT="8082"
APP_JWT_SECRET="cZ08EMbLEbj3REp7fU"
APP_ADMIN_JWT_SECRET="H12EcCoGvNpKYrp6Vw"
APP_AUTH_EXPIRE_HOURS=24
OTP_EXPIRE_SECONDS=45
APP_ENVIRONMENT=dev

DB_NAME="lytemp"
DB_HOST="127.0.0.1"
DB_PASSWORD="db_password"
DB_USER="root"
DB_PORT="5432"
DB_DEBUG="false"
DB_MIGRATE="true"

MEILISEARCH_KEY="ms"
MEILISEARCH_PORT="7700"
MEILISEARCH_HOST="localhost"
MEILISEARCH_MASTER_KEY="meilisearch_master_key"
MEILISEARCH_API_KEY="meilisearch_master_key"
  

REDIS_HOST="localhost"
REDIS_PORT="6379"
REDIS_PASSWORD="redis_password"
REDIS_EXPORTER_PORT="9121"

PROMETHEUS_PORT="9090"

GRAFANA_SECURITY_ADMIN_USER="grafana_security_admin_user"
GRAFANA_SECURITY_ADMIN_PASSWORD="grafana_security_admin_password"
GRAFANA_PORT="3000"

LOGAR_ADMIN_USERNAME="lytemp"
LOGAR_ADMIN_PASSWORD="81tLsVj4uSMtudTr8k"

MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_USERNAME=no-reply@lytemp.com
MAIL_PASSWORD=super_secure_password
MAIL_FROM_NAME=lytemp
MAIL_FROM_EMAIL=no-reply@lytemp.com
MAIL_SSL=false

TELEGRAM_BOT_TOKEN="asdasd"

IYZICO_SANDBOX_API_KEY="sandbox-abc"
IYZICO_SANDBOX_SECRET_KEY="def"
IYZICO_SANDBOX_BASE_URL="https://sandbox-api.iyzipay.com"

ONESIGNAL_APP_ID=""
ONESIGNAL_REST_API_KEY=""

PROVIDER1_URL="https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/v2/provider1"

PROVIDER2_URL="https://raw.githubusercontent.com/WEG-Technology/mock/refs/heads/main/v2/provider2"
```

> **Note:**
> 
> -   If you run the backend **inside Docker**, keep `DB_HOST=postgres_db` and `REDIS_HOST=redis` (service names).
>     
> -   If you run the backend **on your host** with only Postgres/Redis in Docker, set `DB_HOST=127.0.0.1`, `REDIS_HOST=127.0.0.1`.
>

## 2) Local Dev (Run Go on host)

Bring up Postgres, Redis, and the API in containers.
```bash
# build & start (dev compose)
make services-start
go run .
```
-   Backend API: **[http://localhost:8082](http://localhost:8082)**
-   Postgres: mapped to host (see your compose file)
-   Redis: mapped to host (see your compose file)
## 3) Running Tests

Run all unit tests (race detector, no cache):
```bash
make test
```

## 4) Notes

-   **Ports:** ensure `APP_PORT=8082` so the frontend can call `http://localhost:8082/api/v1/...`.
    
-   **CORS:** if you serve the frontend from a different origin, enable/configure CORS in the Echo server as needed.
    
-   **Rate limiting & timeouts:** provider HTTP calls are guarded by timeouts, retries, and an in-memory rate limit per host.
    
-   **No migrations required:** the app uses GORM auto-migrate (based on current models).
    
-   **Logging:** application logs to stdout; Docker `logs -f` shows live output.

## 5) Troubleshooting

-   **DB connection errors:** verify `DB_HOST` matches your run mode (container vs host). Ensure Postgres is up.
    
-   **Empty list:** run `GET /api/v1/fetch-provider` first to ingest mock data.
    
-   **CORS issues:** configure allowed origins for the frontend host.
    
-   **Cache confusion:** list results are cached per `(query,type,sort,page,per_page)`. Wait TTL or vary params to bypass.