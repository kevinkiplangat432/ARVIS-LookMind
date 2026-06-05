# ARVIS — Technical Documentation

> **Version**: 0.1.0 (Alpha)
> **Last Updated**: 2025
> **Maintained By**: Kevin — kiplangatkevin335@gmail.com

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Problem Definition](#problem-definition)
3. [System Architecture](#system-architecture)
4. [Repository Structure](#repository-structure)
5. [Go Backend](#go-backend)
6. [Anomaly Detection](#anomaly-detection)
7. [Database Schema](#database-schema)
8. [Management API](#management-api)
9. [React Dashboard](#react-dashboard)
10. [Configuration](#configuration)
11. [Testing Strategy](#testing-strategy)
12. [Deployment](#deployment)
13. [Roadmap](#roadmap)
14. [Engineering Principles](#engineering-principles)
15. [Tech Stack Decisions](#tech-stack-decisions)

---

## Executive Summary

ARVIS is a Go-based reverse proxy that sits between any application and any LLM API. It intercepts every request, records it to Postgres, runs anomaly detection rules, and exposes the data through a REST API consumed by a React dashboard.

The entire backend is Go. There is no Python, no FastAPI, no separate management service. The proxy and the API are two HTTP listeners in one binary. React is the only frontend technology.

---

## Problem Definition

Organisations deploying LLM-powered applications face these risks:

- No visibility into what is being sent to external AI APIs
- No mechanism to detect unusual or expensive request patterns
- No audit trail for compliance or incident investigation
- Vendor lock-in through direct API coupling in application code

ARVIS solves all of this at the infrastructure level. Applications point their LLM base URL at the proxy — nothing else changes.

---

## System Architecture

```
Client Application
       │
       ▼
┌──────────────────────────────────────────────┐
│          Go Process                          │
│                                              │
│  ┌─────────────────────────────────────┐    │
│  │  Proxy Server (:8080)               │    │
│  │  httputil.ReverseProxy              │    │
│  │  → intercept → forward → log async  │    │
│  └─────────────────────────────────────┘    │
│                                              │
│  ┌─────────────────────────────────────┐    │
│  │  API Server (:8081)                 │    │
│  │  chi router                         │    │
│  │  GET /health                        │    │
│  │  GET /requests                      │    │
│  │  GET /anomalies                     │    │
│  └─────────────────────────────────────┘    │
└──────────────────────────────────────────────┘
       │
       ▼
  PostgreSQL
  (requests, anomalies)
       │
       ▼
  React Dashboard (:3000)
  polls Go API every 5s
```

### Data Flow

1. Client sends request to `:8080`
2. Middleware logs method + path + latency
3. Forwarder proxies to upstream LLM API
4. After response, a goroutine writes a `requests` row to Postgres
5. Detector runs rules against the request — any triggered rule writes an `anomalies` row
6. Dashboard polls `GET /requests` and `GET /anomalies` every 5s, renders live feed

---

## Repository Structure

```
arvis/
├── cmd/arvis/main.go            ← starts proxy goroutine + API server
├── internal/
│   ├── config/config.go         ← env-based config, sensible defaults
│   ├── proxy/
│   │   ├── proxy.go             ← New() wires everything together
│   │   ├── forwarder.go         ← ServeHTTP, async logging, detector call
│   │   └── middleware.go        ← withMiddleware wraps handler with logger
│   ├── api/
│   │   ├── router.go            ← NewRouter() returns chi.Mux
│   │   ├── middleware.go        ← withLogging (unused by chi, available)
│   │   └── handlers/
│   │       ├── health.go        ← GET /health → {"status":"ok"}
│   │       ├── requests.go      ← GET /requests?limit=N
│   │       └── anomalies.go     ← GET /anomalies?limit=N
│   ├── store/
│   │   ├── db.go                ← pgxpool.New
│   │   ├── requests.go          ← Request struct, InsertRequest, ListRequests
│   │   └── anomalies.go         ← Anomaly struct, InsertAnomaly, ListAnomalies
│   └── detector/
│       └── rules.go             ← Check(Request) []Flag
├── dashboard/                   ← React + TypeScript
│   ├── src/
│   │   ├── api/client.ts        ← getRequests(), getAnomalies()
│   │   ├── components/
│   │   │   ├── RequestTable.tsx
│   │   │   ├── AnomalyFeed.tsx
│   │   │   └── StatCards.tsx
│   │   ├── pages/
│   │   │   ├── Dashboard.tsx    ← polls, computes stats, renders both tables
│   │   │   └── Requests.tsx     ← full request log (limit 200)
│   │   ├── App.tsx              ← BrowserRouter, nav, routes
│   │   └── main.tsx
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts           ← /api/* proxied to localhost:8081
├── migrations/
│   ├── 001_create_requests.sql
│   └── 002_create_anomalies.sql
├── scripts/seed.go
├── docker-compose.yml
├── .env.example
├── Makefile
├── go.mod
└── go.sum
```

---

## Go Backend

### Entry Point — `cmd/arvis/main.go`

Loads config, opens the database pool, starts the proxy in a goroutine, starts the API in the main goroutine. Both servers block on `http.ListenAndServe`.

```go
go func() {
    log.Fatal(http.ListenAndServe(cfg.ProxyAddr, proxy.New(cfg, db)))
}()
log.Fatal(http.ListenAndServe(cfg.APIAddr, api.NewRouter(db)))
```

### Config — `internal/config/config.go`

All config from environment variables. Defaults allow the binary to start without any `.env` file.

| Variable | Default | Purpose |
|---|---|---|
| `PROXY_ADDR` | `:8080` | Proxy listen address |
| `API_ADDR` | `:8081` | API listen address |
| `DATABASE_URL` | `postgres://arvis:arvis@localhost:5432/arvis?sslmode=disable` | Postgres connection |
| `TARGET_URL` | `https://api.openai.com` | Upstream LLM API |
| `API_KEY` | `""` | LLM API key (forwarded in requests) |
| `MAX_TOKENS` | `4096` | Max tokens per request |

### Proxy — `internal/proxy/`

**proxy.go** — `New(cfg, db)` parses the target URL, creates an `httputil.ReverseProxy`, wraps it with middleware, and returns the handler.

**forwarder.go** — the core `ServeHTTP` implementation:

1. Records start time
2. Wraps `ResponseWriter` in a `statusRecorder` to capture the status code
3. Calls the upstream `ReverseProxy`
4. Spawns a goroutine to write the `Request` row and run the detector

The goroutine never blocks the response path. If the database is slow or unavailable, the proxy continues serving requests — log writes will simply fail silently and log to stderr.

**middleware.go** — `withMiddleware` wraps any handler with a logger that prints method, path, and latency.

### Store — `internal/store/`

**db.go** — `Connect(url string)` returns a `*pgxpool.Pool`. The pool is created with `context.Background()` and shared across both the proxy and API.

**requests.go** — defines the `Request` struct and two functions:

- `InsertRequest(ctx, db, r)` — single row insert
- `ListRequests(ctx, db, limit)` — `ORDER BY created_at DESC LIMIT $1`

**anomalies.go** — same pattern for `Anomaly`:

- `InsertAnomaly(ctx, db, a)`
- `ListAnomalies(ctx, db, limit)`

All database functions accept a `context.Context` and return errors to the caller.

---

## Anomaly Detection

`internal/detector/rules.go` — `Check(r store.Request) []Flag`

Three rules run on every logged request:

| Rule | Condition | Detail |
|---|---|---|
| `high_latency` | `latency_ms > 10000` | latency exceeded 10s |
| `high_token_usage` | `prompt_tokens + completion_tokens > 3000` | total tokens exceeded 3000 |
| `upstream_error` | `status_code >= 500` | upstream returned 5xx |

Each triggered rule produces a `Flag{Rule, Detail}`. The forwarder converts each flag into an `Anomaly` row with a new UUID, linked to the request by `request_id`.

Adding new rules means adding a condition block in `Check`. No other files need to change.

---

## Database Schema

### `requests`

```sql
CREATE TABLE IF NOT EXISTS requests (
    id                  TEXT PRIMARY KEY,
    model               TEXT,
    prompt_tokens       INTEGER DEFAULT 0,
    completion_tokens   INTEGER DEFAULT 0,
    latency_ms          INTEGER DEFAULT 0,
    status_code         INTEGER DEFAULT 200,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_requests_created ON requests(created_at DESC);
```

### `anomalies`

```sql
CREATE TABLE IF NOT EXISTS anomalies (
    id          TEXT PRIMARY KEY,
    request_id  TEXT REFERENCES requests(id),
    rule        TEXT NOT NULL,
    detail      TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_anomalies_created ON anomalies(created_at DESC);
```

Both tables are append-only. No update or delete operations are ever performed.

Migration files live in `migrations/` and are mounted into Postgres via Docker's `docker-entrypoint-initdb.d` — they run automatically on first container start.

---

## Management API

`internal/api/router.go` — `NewRouter(db)` returns a `chi.Mux` with three routes:

### `GET /health`

```json
{"status": "ok"}
```

### `GET /requests?limit=N`

Returns up to N requests ordered by `created_at DESC`. Default limit: 50, max enforced by the query.

```json
[
  {
    "id": "uuid",
    "model": "gpt-4o",
    "prompt_tokens": 120,
    "completion_tokens": 80,
    "latency_ms": 340,
    "status_code": 200,
    "created_at": "2025-01-01T12:00:00Z"
  }
]
```

### `GET /anomalies?limit=N`

Returns up to N anomalies ordered by `created_at DESC`. Default limit: 50.

```json
[
  {
    "id": "uuid",
    "request_id": "uuid",
    "rule": "high_latency",
    "detail": "latency exceeded 10s",
    "created_at": "2025-01-01T12:00:00Z"
  }
]
```

All handlers follow the same pattern: parse optional `limit` from query string, call store function, encode JSON response.

---

## React Dashboard

### Tech Stack

| Layer | Technology |
|---|---|
| Framework | React 18 |
| Build tool | Vite |
| Language | TypeScript |
| Routing | React Router v6 |
| API calls | native `fetch` |

No external component library, no state management library. Keeps the dependency surface minimal.

### API Client — `src/api/client.ts`

Two functions: `getRequests(limit?)` and `getAnomalies(limit?)`. Both call `/api/*` which Vite proxies to `localhost:8081` in development. In production, Nginx handles the routing.

### Pages

**Dashboard.tsx** — polls both endpoints every 5 seconds via `setInterval`. Computes stats (total requests, total tokens, average latency) client-side from the request array. Renders `StatCards`, `RequestTable`, and `AnomalyFeed`.

**Requests.tsx** — loads 200 requests on mount. Full paginated log view.

### Components

**StatCards.tsx** — three cards: requests count, total tokens, average latency.

**RequestTable.tsx** — table with columns: time, model, tokens, latency, status.

**AnomalyFeed.tsx** — table with columns: time, rule, detail, request ID (first 8 chars).

### Vite Proxy Config

```typescript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8081',
      rewrite: path => path.replace(/^\/api/, '')
    }
  }
}
```

This means the dashboard makes calls to `/api/requests` in development, which are transparently forwarded to `http://localhost:8081/requests`.

---

## Configuration

### `.env.example`

```bash
PROXY_ADDR=:8080
API_ADDR=:8081
DATABASE_URL=postgres://arvis:arvis@localhost:5432/arvis?sslmode=disable
TARGET_URL=https://api.openai.com
API_KEY=your_llm_api_key_here
MAX_TOKENS=4096
```

Copy to `.env` before running. Never commit `.env`.

### `docker-compose.yml`

Spins up Postgres and the Go binary. Postgres mounts `migrations/` into `docker-entrypoint-initdb.d` for automatic schema creation. The Go service depends on Postgres.

---

## Testing Strategy

### Unit Tests

Test detector rules in isolation — no database required:

```go
func TestHighLatency(t *testing.T) {
    r := store.Request{LatencyMs: 15000, StatusCode: 200}
    flags := detector.Check(r)
    if len(flags) != 1 || flags[0].Rule != "high_latency" {
        t.Fatalf("expected high_latency flag, got %v", flags)
    }
}
```

### Integration Tests

Spin up a real Postgres via Docker, run migrations, insert rows via `InsertRequest`, query via `ListRequests`, assert correctness.

### End-to-End

`docker compose up`, make a request through the proxy, query `GET /requests` on the API, assert the row is present and latency is recorded.

Run all tests:

```bash
make test
# go test ./...
```

---

## Deployment

### Local

```bash
make dev        # docker compose up -d
make logs       # docker compose logs -f
make down       # docker compose down
```

### VPS (Phase 2)

1. Copy `docker-compose.yml` and `.env` to the server
2. Add an Nginx config to terminate HTTPS and reverse proxy `:8080` (proxy) and `:8081` (API)
3. Bind Postgres to `127.0.0.1` in `docker-compose.yml`
4. Run `docker compose up -d`

Nginx routing example:

```nginx
server {
    server_name proxy.yourcompany.co.ke;
    location / { proxy_pass http://127.0.0.1:8080; }
}

server {
    server_name api.yourcompany.co.ke;
    location / { proxy_pass http://127.0.0.1:8081; }
}
```

### Production (Phase 3)

Multiple Go instances behind a load balancer. Postgres moves to a managed database (RDS or equivalent). Container orchestration (ECS, Kubernetes) handles restarts and zero-downtime deploys.

---

## Roadmap

### Near Term
- **PII detection** — regex-based Kenyan National ID, M-PESA, KRA PIN masking before forwarding
- **Redis policy engine** — per-key token budgets enforced before proxying
- **Kill switch** — terminate response stream on policy violation mid-stream
- **Request body parsing** — extract model and token counts from actual OpenAI request/response JSON

### Medium Term
- **Per-key identity** — attribute requests to specific employees or agents
- **Webhook alerts** — POST to Slack/PagerDuty when anomaly rules fire
- **Semantic rules** — embedding-based topic blocking beyond keyword matching
- **Shadow mode** — log everything, block nothing, for zero-risk enterprise onboarding

### Long Term
- **Cost attribution** — token spend breakdown by team or agent
- **ML-based anomaly detection** — baseline normal behaviour, alert on deviations
- **Forensic replay** — replay any session exactly as it occurred
- **Compliance packs** — pre-built rule sets for specific regulatory frameworks

---

## Engineering Principles

### Go for everything backend
One language, one binary, one process. The proxy and the API are not separate services — they are two `http.ListenAndServe` calls in the same `main()`. This eliminates inter-service latency, simplifies deployment, and reduces operational surface area.

### Async on the critical path
The proxy's critical path is: receive request → forward to LLM → return response. Everything else — database writes, anomaly detection — happens in goroutines after the response is sent. The proxy adds microseconds of overhead, not milliseconds.

### Append-only storage
No row in `requests` or `anomalies` is ever updated or deleted. This makes the audit trail tamper-evident and simplifies the data model.

### Minimal dependencies
Go standard library + chi (routing) + pgx (Postgres). No ORM, no framework, no middleware forest. Every dependency has a clear justification.

### Fail open
If the database is unavailable, the proxy keeps proxying. Log write errors are emitted to stderr but do not interrupt request handling.

---

## Tech Stack Decisions

| Layer | Technology | Reason |
|---|---|---|
| Proxy engine | Go `httputil.ReverseProxy` | Standard library, zero-copy streaming, minimal latency |
| Management API | Go chi | Lightweight router, composable middleware, no magic |
| Database driver | pgx v5 | Native Postgres protocol, connection pooling built in |
| Database | PostgreSQL 16 | ACID, JSONB for future flexibility, mature tooling |
| Frontend | React + Vite + TypeScript | Component model, fast builds, type safety |
| Containers | Docker + Docker Compose | Simple local dev, same config on VPS |
| Production proxy | Nginx | HTTPS termination, subdomain routing |
| Observability | Prometheus + Grafana (Phase 3) | Industry standard, open source |

---

**Document Version**: 0.1.0
**Review Cycle**: Per milestone
