<p align="center">
  <img src="LookMind.png" alt="LookMind" width="300" />
</p>

# ARVIS — Technical Documentation

> **Version**: 0.1.0 (Alpha) · **Last Updated**: 2025
> **Maintained by**: Kevin — kiplangatkevin335@gmail.com

This document is the technical reference for ARVIS. It covers architecture, data flow, every internal package, the API contract, the dashboard, testing, and deployment in depth.

For context on why this system exists and how it evolved from a Python SDK, start with [`docs/HISTORY.md`](docs/HISTORY.md).
For a high-level introduction and quick-start, see [`README.md`](README.md).
For the feature build tracker, see [`FEATURES.md`](FEATURES.md).

---

## Table of Contents

1. [System Architecture](#system-architecture)
2. [Go Backend](#go-backend)
3. [Anomaly Detection](#anomaly-detection)
4. [Database Schema](#database-schema)
5. [Management API](#management-api)
6. [React Dashboard](#react-dashboard)
7. [Configuration](#configuration)
8. [Testing Strategy](#testing-strategy)
9. [Deployment](#deployment)
10. [Engineering Principles](#engineering-principles)
11. [Tech Stack Decisions](#tech-stack-decisions)

---

## System Architecture

```
Client Application
       │
       ▼
┌──────────────────────────────────────────────────┐
│              Go Process (one binary)             │
│                                                  │
│  ┌──────────────────────────────────────────┐   │
│  │  Proxy Server (:8080)                    │   │
│  │  httputil.ReverseProxy                   │   │
│  │  → intercept → log metadata → forward   │   │
│  │  → capture latency + status              │   │
│  │  → goroutine: write DB + run detector    │   │
│  └──────────────────────────────────────────┘   │
│                                                  │
│  ┌──────────────────────────────────────────┐   │
│  │  API Server (:8081)                      │   │
│  │  chi router                              │   │
│  │  GET /health                             │   │
│  │  GET /requests?limit=N                   │   │
│  │  GET /anomalies?limit=N                  │   │
│  └──────────────────────────────────────────┘   │
└──────────────────────────────────────────────────┘
       │                        │
       ▼                        ▼
  LLM API               PostgreSQL
  (any provider)        requests, anomalies
                               │
                               ▼
                     React Dashboard (:3000)
                     polls Go API every 5s
```

### Data Flow

1. Client sends a request to `:8080`
2. Proxy middleware logs method, path, and latency to stdout
3. Forwarder wraps the `ResponseWriter` in a `statusRecorder` to capture the status code
4. `httputil.ReverseProxy` forwards the request to the upstream LLM API
5. After the response is written, a goroutine is spawned:
   - Builds a `store.Request` struct from captured metadata
   - Calls `store.InsertRequest` — writes one row to `requests`
   - Calls `detector.Check` — evaluates all rules against the request
   - For each triggered rule, calls `store.InsertAnomaly` — writes one row to `anomalies`
6. The goroutine runs entirely off the critical path — the response has already been returned to the client
7. React dashboard polls `GET /requests` and `GET /anomalies` every 5 seconds and renders the result

### Why one binary, two listeners

Running the proxy and the API in a single process eliminates inter-service latency, simplifies deployment to a single Docker image, and means both servers share one database connection pool. The tradeoff is that they cannot be scaled independently — this is acceptable until Phase 3, where the proxy will be extracted into its own deployable unit.

---

## Go Backend

### Entry Point — `cmd/arvis/main.go`

```go
// Proxy runs in a goroutine so main can block on the API server
go func() {
    log.Fatal(http.ListenAndServe(cfg.ProxyAddr, proxy.New(cfg, db)))
}()
log.Fatal(http.ListenAndServe(cfg.APIAddr, api.NewRouter(db)))
```

Loads config from environment, opens the pgxpool connection, starts the proxy as a goroutine, then blocks on the API server. If either server exits, the binary exits and Docker restarts it.

### Config — `internal/config/config.go`

All configuration is read from environment variables. Every variable has a default that allows the binary to start and serve without a `.env` file present.

| Variable | Default | Purpose |
|---|---|---|
| `PROXY_ADDR` | `:8080` | Proxy listen address |
| `API_ADDR` | `:8081` | API listen address |
| `DATABASE_URL` | `postgres://arvis:arvis@localhost:5432/arvis?sslmode=disable` | Postgres connection string |
| `TARGET_URL` | `https://api.openai.com` | Upstream LLM API base URL |
| `API_KEY` | `""` | LLM API key — forwarded on every proxied request |
| `MAX_TOKENS` | `4096` | Maximum tokens per request |

### Proxy — `internal/proxy/`

**`proxy.go`** — `New(cfg, db)` is the public constructor. It parses the target URL, creates an `httputil.ReverseProxy`, wires the forwarder, and wraps it with the logging middleware.

**`forwarder.go`** — implements `http.Handler`. The core logic:

```
ServeHTTP(w, r):
  1. record start time
  2. wrap w in statusRecorder
  3. call ReverseProxy.ServeHTTP (blocks until upstream responds)
  4. go func() { InsertRequest + detector.Check + InsertAnomalies }
```

The goroutine in step 4 uses `context.Background()` — it is detached from the request context so it continues even after the HTTP connection closes.

`statusRecorder` is a minimal `http.ResponseWriter` wrapper that intercepts `WriteHeader` to capture the status code:

```go
type statusRecorder struct {
    http.ResponseWriter
    code int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.code = code
    r.ResponseWriter.WriteHeader(code)
}
```

**`middleware.go`** — `withMiddleware` wraps any `http.Handler` with a logger that prints method, path, and elapsed time after the handler returns.

### Store — `internal/store/`

**`db.go`** — `Connect(url)` calls `pgxpool.New` and returns the pool. The pool is passed to both `proxy.New` and `api.NewRouter` — they share the same set of connections.

**`requests.go`**

```go
type Request struct {
    ID           string
    Model        string
    PromptTokens int
    CompTokens   int
    LatencyMs    int
    StatusCode   int
    CreatedAt    time.Time
}

func InsertRequest(ctx, db, r) error
func ListRequests(ctx, db, limit) ([]Request, error)
```

`ListRequests` runs `ORDER BY created_at DESC LIMIT $1`. It returns an empty slice (not nil) when there are no rows.

**`anomalies.go`** — same pattern:

```go
type Anomaly struct {
    ID        string
    RequestID string
    Rule      string
    Detail    string
    CreatedAt time.Time
}

func InsertAnomaly(ctx, db, a) error
func ListAnomalies(ctx, db, limit) ([]Anomaly, error)
```

All store functions take a `context.Context` as the first argument and return errors to the caller. They never log internally.

---

## Anomaly Detection

`internal/detector/rules.go` exports a single function:

```go
func Check(r store.Request) []Flag
```

A `Flag` is a `{Rule string, Detail string}` pair. `Check` runs all rules and returns every flag that triggers — zero, one, or many per request.

### Current Rules

| Rule constant | Condition | Detail message |
|---|---|---|
| `high_latency` | `r.LatencyMs > 10000` | `"latency exceeded 10s"` |
| `high_token_usage` | `r.PromptTokens + r.CompTokens > 3000` | `"total tokens exceeded 3000"` |
| `upstream_error` | `r.StatusCode >= 500` | `"upstream returned 5xx"` |

### Adding a rule

Add a condition block to `Check`. No other file needs to change:

```go
if r.PromptTokens > 2000 {
    flags = append(flags, Flag{
        Rule:   "large_prompt",
        Detail: "prompt tokens exceeded 2000",
    })
}
```

The forwarder automatically converts every returned flag into an `Anomaly` row linked to the request by `request_id`.

---

## Database Schema

Migration files live in `migrations/` and are mounted into Postgres via `docker-entrypoint-initdb.d`. They run automatically on first container start.

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

Both tables are append-only. No `UPDATE` or `DELETE` operations are ever performed. The `request_id` foreign key ensures every anomaly is traceable to the request that produced it.

### Planned additions (Phase 2)

```sql
-- Phase 2: per-client isolation
CREATE TABLE clients (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    api_key     TEXT UNIQUE NOT NULL,  -- hashed
    created_at  TIMESTAMPTZ DEFAULT NOW()
);

-- Phase 2: PII tokenisation
CREATE TABLE pii_tokens (
    token          TEXT PRIMARY KEY,   -- e.g. [KENYAN_ID_TOKEN_1]
    original_value TEXT NOT NULL,      -- encrypted at rest
    client_id      TEXT REFERENCES clients(id),
    created_at     TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Management API

`internal/api/router.go` — `NewRouter(db)` returns a `chi.Mux`. chi's built-in `middleware.Logger` and `middleware.Recoverer` are applied globally.

### `GET /health`

No authentication. Returns 200 immediately.

```json
{"status": "ok"}
```

### `GET /requests?limit=N`

Returns up to N requests ordered by `created_at DESC`. If `limit` is absent or zero, defaults to 50. Encoded as a JSON array — empty array `[]` when there are no rows.

```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
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
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "rule": "high_latency",
    "detail": "latency exceeded 10s",
    "created_at": "2025-01-01T12:00:05Z"
  }
]
```

### Error responses

All handlers return `500 Internal Server Error` with a plain text body if the database query fails. There are no 4xx responses from the current endpoints — invalid `limit` values are silently coerced to the default.

---

## React Dashboard

### Stack

| Layer | Technology | Why |
|---|---|---|
| Framework | React 18 | Component model, large ecosystem |
| Build tool | Vite | Fast HMR, simple config |
| Language | TypeScript | Type safety on API responses |
| Routing | React Router v6 | Standard, well-documented |
| API calls | native `fetch` | No additional dependency needed |

No component library, no state management library. The dependency surface is intentionally small.

### Development proxy — `vite.config.ts`

```typescript
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:8081',
      rewrite: path => path.replace(/^\/api/, '')
    }
  }
}
```

The dashboard calls `/api/requests`. Vite rewrites this to `http://localhost:8081/requests`. In production, Nginx performs the same rewrite — no code change needed between environments.

### API client — `src/api/client.ts`

```typescript
export async function getRequests(limit = 50): Promise<Request[]>
export async function getAnomalies(limit = 50): Promise<Anomaly[]>
```

Both functions throw on non-2xx responses. The caller catches and ignores errors — polling continues regardless.

### Pages

**`Dashboard.tsx`** — the primary view. On mount and every 5 seconds:
1. Calls `getRequests()` and `getAnomalies()` in parallel
2. Computes stats from the requests array: total count, sum of all tokens, mean latency
3. Renders `StatCards`, `RequestTable`, and `AnomalyFeed`

**`Requests.tsx`** — full log view. Loads 200 requests on mount, no polling. Intended for investigating a specific time window.

### Components

**`StatCards.tsx`** — three number cards: total requests, total tokens, average latency (ms).

**`RequestTable.tsx`** — columns: time, model, tokens (prompt + completion), latency, status code.

**`AnomalyFeed.tsx`** — columns: time, rule, detail, request ID (first 8 characters).

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

Copy to `.env` before running. The `.gitignore` excludes `.env`. Never commit it.

### `docker-compose.yml`

- `postgres` service mounts `./migrations` into `/docker-entrypoint-initdb.d` — schema is created automatically on first start
- `arvis` service depends on `postgres`, reads from `.env`, exposes `:8080` and `:8081`
- A named volume `pgdata` persists data across `docker compose down` and `up`

---

## Testing Strategy

### Unit — detector rules

No database, no network. Fast to run, easy to reason about:

```go
func TestHighLatency(t *testing.T) {
    flags := detector.Check(store.Request{LatencyMs: 15000, StatusCode: 200})
    require.Len(t, flags, 1)
    require.Equal(t, "high_latency", flags[0].Rule)
}

func TestNoFlagsOnCleanRequest(t *testing.T) {
    flags := detector.Check(store.Request{LatencyMs: 200, StatusCode: 200, PromptTokens: 100, CompTokens: 50})
    require.Empty(t, flags)
}
```

### Integration — store layer

Requires a running Postgres. Use Docker or `testcontainers-go`:

```go
// Insert a request, list it back, assert round-trip fidelity
func TestInsertAndListRequest(t *testing.T) {
    db := connectTestDB(t)
    r := store.Request{ID: uuid.NewString(), Model: "gpt-4o", LatencyMs: 100, StatusCode: 200, CreatedAt: time.Now()}
    require.NoError(t, store.InsertRequest(ctx, db, r))

    rows, err := store.ListRequests(ctx, db, 10)
    require.NoError(t, err)
    require.Len(t, rows, 1)
    require.Equal(t, r.ID, rows[0].ID)
}
```

### End-to-end

`docker compose up`, send a request through the proxy, assert the row appears via `GET /requests`:

```bash
docker compose up -d
curl -s -o /dev/null http://localhost:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4o","messages":[{"role":"user","content":"hello"}]}'

sleep 1
curl http://localhost:8081/requests | jq '.[0].latency_ms'
# should print a number
```

Run all Go tests:

```bash
make test
# go test ./...
```

---

## Deployment

### Local (Phase 1)

```bash
make dev        # docker compose up -d
make logs       # docker compose logs -f
make seed       # insert demo data
make down       # docker compose down
```

### VPS — Phase 2

1. Copy `docker-compose.yml`, `.env`, `migrations/` to the server
2. Bind Postgres to `127.0.0.1` in `docker-compose.yml`
3. Install Nginx and configure subdomain routing:

```nginx
server {
    listen 443 ssl;
    server_name proxy.yourcompany.co.ke;
    location / { proxy_pass http://127.0.0.1:8080; }
}

server {
    listen 443 ssl;
    server_name api.yourcompany.co.ke;
    location / { proxy_pass http://127.0.0.1:8081; }
}
```

4. Obtain TLS certificate with Certbot: `certbot --nginx -d proxy.yourcompany.co.ke`
5. `docker compose up -d`

### Production — Phase 3

Multiple stateless Go instances behind a load balancer. Postgres moves to a managed database (RDS, Cloud SQL, or equivalent). Container orchestration (ECS Fargate or Kubernetes) handles zero-downtime deploys and automatic restarts. Prometheus and Grafana provide observability. See [`FEATURES.md`](FEATURES.md) Phase 4 for the full list of infrastructure work.

---

## Engineering Principles

### Go for everything backend
One language, one binary, one process, one container. The proxy and API are two `http.ListenAndServe` calls in `main()`. No inter-service calls, no serialisation overhead between components, no separate deployment pipeline.

### Async writes on the critical path
`receive → forward → respond` is the only work on the critical path. Database writes, anomaly detection, and any future rule evaluation all happen in goroutines after the response is returned. The proxy adds microseconds of overhead, not milliseconds.

### Append-only audit trail
No row is ever updated or deleted. This makes the audit trail tamper-evident by construction, not by policy.

### Fail open
If Postgres is unavailable, the proxy continues proxying. Log write errors go to stderr. The proxy is not a single point of failure for the application it governs.

### Minimal dependencies
Go standard library + chi (routing) + pgx (Postgres). No ORM. No framework. Every external dependency earns its place.

---

## Tech Stack Decisions

| Layer | Technology | Reason |
|---|---|---|
| Proxy engine | Go `httputil.ReverseProxy` | Standard library, zero-copy streaming, microsecond overhead |
| Management API | Go chi | Lightweight router, idiomatic middleware, no magic |
| Database driver | pgx v5 | Native Postgres wire protocol, built-in connection pooling |
| Database | PostgreSQL 16 | ACID compliance, JSONB for future flexibility, append-only possible |
| Frontend | React + Vite + TypeScript | Component model, fast builds, type-safe API responses |
| Containers | Docker + Docker Compose | Reproducible local dev, identical config on VPS |
| Production routing | Nginx | HTTPS termination, subdomain routing, well-understood |
| Observability | Prometheus + Grafana (Phase 3) | Industry standard, open source, integrates with Go metrics |

---

**Document Version**: 0.1.0 · **Review**: per milestone
