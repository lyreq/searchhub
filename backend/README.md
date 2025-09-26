
# # LyTemp Boilerplate (Echo Framework)

This repository provides a production-ready **Golang Echo** boilerplate with a clean modular structure, service integrations (**Postgres**, **Redis**, **Meilisearch**), and monitoring stack (**Prometheus** + **Grafana**).


## 🚀 Features
-   Modular project structure (`internal`, `pkg`, `server`) with clear separation of concerns
    
-   Postgres, Redis, Meilisearch containers via Docker Compose
    
-   Prometheus + Grafana monitoring stack
    
-   Air hot-reload support for fast development
    
-   Configurable via `.env` + `config.yml`
    
-   Includes common integrations: Mailer, Redis, Iyzico, OneSignal, Telegram, RevenueCat, etc.


## 📦 Prerequisites

Make sure you have the following installed:
-   Docker & Docker Compose v2
  
-   Go 1.22+
    
-   Make
## ⚙️ Configuration

Copy `.env` and update values as needed:
```bash
cp .env.example .env
```
-   The `.env` file defines container ports, database credentials, API keys, and environment variables.
    
-   Application-level configuration is managed in `config.yml`.  
    All values are mapped from `.env`.

## 🛠️ Running in Development

You have two main workflows:

### 1) Host-based (manual run)

Bring up dependencies (Postgres, Redis, Meilisearch, etc.) with Docker, then run the app locally:
```bash
make services-start
go run .
```
-   `make services-start` → starts core services
    
-   `go run .` → runs the Echo app from your machine using `.env` values
    

> In this mode, `DB_HOST=127.0.0.1` is required in `.env`.

### 2. **Containerized (hot-reload with Air)**

Run the app inside Docker with hot-reload:
```bash
make up PROFILE=dev ENV_FILE=.env.docker
```

-   Uses **Air** for live reload on file changes
    
-   Container links use service names (`postgres_db`, `redis`, `meilisearch`) instead of `localhost`
    

> In this mode, your `.env.docker` should set `DB_HOST=postgres_db`.


## Useful Makefile Commands

-   `make up PROFILE=dev` → start full stack in dev mode
    
-   `make down PROFILE=dev` → stop stack
    
-   `make logs PROFILE=dev` → follow logs
    
-   `make restart PROFILE=dev` → restart containers
    

----------

## 📊 Monitoring

-   Prometheus → [http://localhost:9090](http://localhost:9090)
    
-   Grafana → [http://localhost:3000](http://localhost:3000) (default admin user/password from `.env`)
    
-   Redis Exporter → [http://localhost:9121/metrics](http://localhost:9121/metrics)

## 📂 Project Structure (highlights)
```bash
internal/       # Domain logic, handlers, routers, middleware
pkg/            # Shared utilities & integrations (db, redis, mailer, etc.)
server/         # Server setup, validation
deployments/    # Dockerfiles (dev/stage/prod)
scripts/        # Helper scripts for deploy/restart
config.yml      # App configuration
.env            # Environment variables
```

## 📝 Notes

-   Do not commit real secrets inside `.env`. Use `.env.example` as a template.
    
-   For staging and production, use `PROFILE=stage` / `PROFILE=prod` with their respective env files.
    
-   Healthchecks are enabled so the app waits for Postgres, Redis, and Meilisearch to be ready.

✅ That’s it! After following the steps, you’ll have **LyTemp running locally with full dependencies and monitoring.**
