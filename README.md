# ARVIS — AI Request Visibility & Intelligence System

> **Version**: 0.1.0 (Alpha) · **Status**: Active Development
> **Compliance**: Kenya Data Protection Act 2019 · Data residency enforced at the infrastructure level

## Overview

ARVIS is a transparent, real-time governance layer that sits between an organisation's systems and any LLM API. Every prompt and response passes through the proxy — requests are intercepted and logged, anomalies are detected, policy rules are enforced, and a complete audit trail is written to Postgres.

The system is two services: a Go backend that handles both the reverse proxy and the management API, and a React dashboard. They share a Postgres database. The Go binary runs two HTTP servers from a single process — the proxy on `:8080` and the API on `:8081`.

---

## Why This Exists

Enterprises deploying LLM-powered applications face a consistent set of risks:

- Sensitive data leaking into third-party AI APIs
- No visibility into what internal systems are sending to LLMs
- No mechanism to enforce token budgets or block anomalous requests
- No audit trail for regulatory compliance or internal governance

ARVIS addresses all of these at the infrastructure level, requiring zero changes to existing applications. Redirect your LLM base URL to the proxy — everything else works as before.

---

## Architecture

```
Client Application
       │
       ▼
┌──────────────────────────────────────────────┐
│          Go Proxy Engine (:8080)             │
│                                              │
│  1. Intercept request                        │
│  2. Log request metadata                     │
│  3. Forward to LLM API                       │
│  4. Record latency + status                  │
│  5. Run anomaly detection rules              │
│  6. Write to Postgres (async)                │
└──────────────────────────────────────────────┘
       │                         │
       ▼                         ▼
  LLM API                   PostgreSQL
  (OpenAI /                 (requests,
   Anthropic /               anomalies)
   any compatible)
                                │
                                ▼
                     ┌──────────────────────┐
                     │  Go Management API   │
                     │  (:8081)             │
                     │  /requests           │
                     │  /anomalies          │
                     │  /health             │
                     └──────────────────────┘
                                │
                                ▼
                     ┌──────────────────────┐
                     │   React Dashboard    │
                     │   (:3000)            │
                     │  Live feed · Stats   │
                     │  Anomaly feed        │
                     └──────────────────────┘
```

**Service communication rules:**
- Go proxy intercepts every request and writes logs to Postgres asynchronously — never on the critical path
- Go API reads from Postgres to power the dashboard
- React calls Go API REST endpoints — standard JSON over HTTP
- The proxy and API are one Go binary, two HTTP listeners

---

## Repository Structure

```
arvis/
├── cmd/
│   └── arvis/
│       └── main.go              ← entry point, starts proxy + API
│
├── internal/
│   ├── proxy/
│   │   ├── proxy.go             ← wires ReverseProxy, middleware, forwarder
│   │   ├── forwarder.go         ← forwards to LLM, records latency, triggers detector
│   │   └── middleware.go        ← request logging middleware
│   │
│   ├── api/
│   │   ├── router.go            ← chi router, registers all routes
│   │   ├── handlers/
│   │   │   ├── requests.go      ← GET /requests
│   │   │   ├── anomalies.go     ← GET /anomalies
│   │   │   └── health.go        ← GET /health
│   │   └── middleware.go        ← API-level logging
│   │
│   ├── store/
│   │   ├── db.go                ← pgxpool connection
│   │   ├── requests.go          ← InsertRequest, ListRequests
│   │   └── anomalies.go         ← InsertAnomaly, ListAnomalies
│   │
│   ├── detector/
│   │   └── rules.go             ← rule-based anomaly detection
│   │
│   └── config/
│       └── config.go            ← loads config from environment
│
├── dashboard/                   ← React + TypeScript frontend
│   ├── src/
│   │   ├── components/
│   │   │   ├── RequestTable.tsx
│   │   │   ├── AnomalyFeed.tsx
│   │   │   └── StatCards.tsx
│   │   ├── pages/
│   │   │   ├── Dashboard.tsx
│   │   │   └── Requests.tsx
│   │   ├── api/
│   │   │   └── client.ts        ← calls Go API on :8081
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
│
├── migrations/
│   ├── 001_create_requests.sql
│   └── 002_create_anomalies.sql
│
├── scripts/
│   └── seed.go                  ← demo data
│
├── docker-compose.yml
├── .env.example
├── Makefile
├── go.mod
└── README.md
```

---

## Core Features

### Request Interception and Logging
Every request the proxy handles is recorded: timestamp, model, prompt tokens, completion tokens, latency, and status code. Logs are written asynchronously — never on the critical proxy path.

### Rule-Based Anomaly Detection
After each request is logged, the detector runs a set of rules against it. Currently detects: high latency (>10s), high token usage (>3000 total tokens), and upstream 5xx errors. Any triggered rule produces an anomaly record linked to the request.

### Management API
A chi-based REST API on `:8081` exposes the request log and anomaly feed to the dashboard. Pagination via `?limit=N`. Runs in the same process as the proxy — no inter-service calls.

### React Dashboard
Polls the Go API every 5 seconds. Shows live request feed, anomaly feed, and aggregate stats (total requests, total tokens, average latency). No build-time API config needed — Vite proxies `/api/*` to `localhost:8081` in development.

---

## Database Schema

Defined in `migrations/` — runs automatically on first Postgres start via Docker.

```sql
CREATE TABLE requests (
    id                  TEXT PRIMARY KEY,
    model               TEXT,
    prompt_tokens       INTEGER DEFAULT 0,
    completion_tokens   INTEGER DEFAULT 0,
    latency_ms          INTEGER DEFAULT 0,
    status_code         INTEGER DEFAULT 200,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE anomalies (
    id          TEXT PRIMARY KEY,
    request_id  TEXT REFERENCES requests(id),
    rule        TEXT NOT NULL,
    detail      TEXT,
    created_at  TIMESTAMPTZ DEFAULT NOW()
);
```

---

## Getting Started

### Prerequisites
- Docker and Docker Compose
- An API key for your target LLM provider

### Local Development

```bash
git clone https://github.com/your-org/arvis.git
cd arvis

cp .env.example .env
# Edit .env — add your LLM API key and target URL

docker compose up -d
```

| Service        | URL                   |
|----------------|-----------------------|
| Go Proxy       | http://localhost:8080 |
| Go API         | http://localhost:8081 |
| React Dashboard| http://localhost:3000 |

### Verify

```bash
curl http://localhost:8081/health
# {"status":"ok"}
```

### Makefile shortcuts

```bash
make dev      # docker compose up -d
make build    # go build ./cmd/arvis
make test     # go test ./...
make down     # docker compose down
make seed     # insert demo data
make migrate  # run SQL migrations manually
```

---

## Environment Variables

```bash
PROXY_ADDR=:8080
API_ADDR=:8081
DATABASE_URL=postgres://arvis:arvis@localhost:5432/arvis?sslmode=disable
TARGET_URL=https://api.openai.com
API_KEY=your_llm_api_key_here
MAX_TOKENS=4096
```

---

## Deployment Phases

### Phase 1 — Local Proof of Concept
All containers on a single machine via `docker compose up`. Proves the full pipeline: request enters proxy → logged to Postgres → anomaly detected → visible in React dashboard.

**Done when:**
- Proxied requests appear in the dashboard live feed
- Anomalies are detected and shown in the anomaly feed
- `docker compose down && docker compose up` recovers cleanly with all data intact

### Phase 2 — Pilot VPS Deployment
Docker Compose on a cloud VPS. Nginx terminates HTTPS and routes subdomains. Postgres is bound to `127.0.0.1` only. Pilot clients point at the proxy with a single base URL change.

**Done when:**
- Proxy reachable over HTTPS with a valid TLS certificate
- At least one external client routing live traffic through the proxy
- Postgres not internet-accessible

### Phase 3 — Production Scale
Multiple stateless Go instances behind a load balancer. Postgres moves to a managed database with automated backups. Prometheus + Grafana for observability.

**Done when:**
- Load testing confirms 1,000+ concurrent requests with < 200ms added latency
- Managed database with automated backups and failover in place
- Monitoring and alerting stack live

---

## Tech Stack

| Layer | Technology |
|---|---|
| Proxy engine | Go — `httputil.ReverseProxy` |
| Management API | Go — chi router |
| Anomaly detection | Go — rule engine |
| Audit storage | PostgreSQL 16 |
| Frontend | React + Vite + TypeScript |
| Container runtime | Docker + Docker Compose |
| Production routing | Nginx + Let's Encrypt |
| Observability (Phase 3) | Prometheus + Grafana |

---

## Git Workflow

```
main      — production-ready only
dev       — active development
feature/* — one branch per feature
fix/*     — bug fix branches
```

Commit format: `service: what you did (present tense, concise)`

```bash
git commit -m 'proxy: record latency and status per request'
git commit -m 'detector: add high token usage rule'
git commit -m 'api: add pagination to GET /requests'
git commit -m 'dashboard: auto-poll anomaly feed every 5s'
```

---

## Security and Compliance

- Postgres is never exposed to the public internet
- Credentials managed via environment variables, never committed
- Audit logs are append-only — no update or delete on `requests` or `anomalies`
- Data residency enforced at the infrastructure level

---

## Roadmap

- **PII detection and masking** — detect Kenyan National ID, M-PESA, KRA PIN before forwarding
- **Redis-backed policy engine** — token budgets and blocked topics enforced per client key
- **Kill switch** — terminate response stream mid-flight on policy violation
- **Per-key identity** — attribute requests to specific employees or agents
- **Webhook alerts** — Slack/PagerDuty when anomaly rules fire
- **Semantic rules** — meaning-based blocking beyond keyword matching
- **Shadow mode** — log everything, block nothing, for zero-risk onboarding
- **Cost attribution** — token spend breakdown by team or agent

---

## Contact

- Email: kiplangatkevin335@gmail.com
