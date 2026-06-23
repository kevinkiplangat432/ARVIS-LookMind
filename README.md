<!-- markdownlint-disable MD031 MD040 MD033 MD041 MD036 MD060 -->
<p align="center">
  <img src="LookMind.png" alt="LookMind" width="300" />
</p>

# ARVIS — AI Request Visibility & Intelligence System

> **Version**: 0.1.0 (Alpha) · **Status**: Active Development
> **Compliance**: Kenya Data Protection Act 2019 · Data residency enforced at the infrastructure level

---

## The Short Version

Organisations are sending sensitive data to AI APIs with no visibility into what is being sent, no way to detect unusual behaviour, and no audit trail for compliance. ARVIS fixes this at the infrastructure level — no code changes required. Redirect your LLM base URL to the proxy and every request is intercepted, logged, analysed, and surfaced in a live dashboard.

This project started in February 2026 as a Python SDK. It pivoted to a Go reverse proxy because the SDK required adoption — every team had to integrate it manually before governance could occur. The proxy requires nothing. See the full story in [`docs/HISTORY.md`](docs/HISTORY.md).

---

## What It Does

A Go binary runs two HTTP servers from a single process:

- **Proxy** on `:8080` — intercepts every LLM request, forwards it upstream, records metadata and anomalies asynchronously
- **API** on `:8081` — serves the request log and anomaly feed as JSON to the dashboard

A React dashboard on `:3000` polls the API every 5 seconds and renders a live feed of requests, flagged anomalies, and aggregate stats.
```

Your Application
      │
      ▼
┌─────────────────────────────────────────────┐
│           Go Process (one binary)           │
│                                             │
│  Proxy :8080          API :8081             │
│  ─────────────        ────────────          │
│  intercept            GET /health           │
│  forward              GET /requests         │
│  detect anomalies     GET /anomalies        │
│  log async                                  │
└─────────────────────────────────────────────┘
      │                        │
      ▼                        ▼
  LLM API               PostgreSQL
  (any provider)        (requests,
                         anomalies)
                               │
                               ▼
                     React Dashboard :3000
                     live feed · stats · anomalies
```

---

## Architecture Decisions

The entire backend is Go. One language, one binary, one process. No Python. No FastAPI. No separate management service. The proxy and the API share the same database pool and run in the same `main()`.

All database writes happen in goroutines after the response is sent — the proxy never blocks on Postgres. If the database is slow or unavailable, requests continue to be proxied. The audit trail catches up when connectivity resumes.

Audit records are append-only. Nothing in `requests` or `anomalies` is ever updated or deleted.

For the full engineering rationale, see [`DOCUMENTATION.md`](DOCUMENTATION.md).

---

## Getting Started

### Prerequisites

- Docker and Docker Compose
- An API key for your target LLM provider

### Run locally

```bash
git clone https://github.com/your-org/arvis.git
cd arvis

cp .env.example .env
# Set TARGET_URL and API_KEY in .env

docker compose up -d
```

| Service         | URL                    |
|-----------------|------------------------|
| Go Proxy        | <http://localhost:8080>  |
| Go API          | <http://localhost:8081>  |
| React Dashboard | <http://localhost:3000>  |

```bash
curl http://localhost:8081/health
# {"status":"ok"}
```

### Environment variables

```bash
PROXY_ADDR=:8080
API_ADDR=:8081
DATABASE_URL=postgres://arvis:arvis@localhost:5432/arvis?sslmode=disable
TARGET_URL=https://api.openai.com
API_KEY=your_llm_api_key_here
MAX_TOKENS=4096
```

### Makefile

```bash
make dev      # docker compose up -d
make build    # go build ./cmd/arvis
make test     # go test ./...
make seed     # insert demo data
make down     # docker compose down
```

---

## Repository Structure

```
arvis/
├── cmd/arvis/main.go            ← entry point — starts proxy + API
├── internal/
│   ├── config/config.go         ← env config with defaults
│   ├── proxy/                   ← reverse proxy, middleware, forwarder
│   ├── api/                     ← chi router, handlers, middleware
│   ├── store/                   ← Postgres queries (requests, anomalies)
│   └── detector/rules.go        ← rule-based anomaly detection
├── dashboard/                   ← React + TypeScript frontend
│   └── src/
│       ├── api/client.ts        ← calls Go API on :8081
│       ├── components/          ← RequestTable, AnomalyFeed, StatCards
│       └── pages/               ← Dashboard, Requests
├── migrations/                  ← SQL run automatically on first Postgres start
├── scripts/seed.go              ← demo data
├── docs/HISTORY.md              ← origin story and the pivot
├── FEATURES.md                  ← build tracker — what's done, what's next
├── DOCUMENTATION.md             ← full technical reference
├── docker-compose.yml
├── .env.example
├── Makefile
└── go.mod
```

---

## Current Anomaly Rules

Three rules run against every logged request. See `internal/detector/rules.go`.

| Rule | Trigger | Notes |
|---|---|---|
| `high_latency` | latency > 10s | Possible upstream issue or prompt too large |
| `high_token_usage` | total tokens > 3000 | Cost risk, potential data dump |
| `upstream_error` | status code ≥ 500 | LLM provider failure |

---

## Database Schema

```sql
CREATE TABLE requests (
    id                TEXT PRIMARY KEY,
    model             TEXT,
    prompt_tokens     INTEGER DEFAULT 0,
    completion_tokens INTEGER DEFAULT 0,
    latency_ms        INTEGER DEFAULT 0,
    status_code       INTEGER DEFAULT 200,
    created_at        TIMESTAMPTZ DEFAULT NOW()
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

## Deployment Phases

### Phase 1 — Local Proof of Concept ← we are here

`docker compose up` on a single machine. Pipeline proven end-to-end: request enters proxy → logged to Postgres → anomaly detected → visible in the dashboard.

Done when:

- Proxied requests appear in the live feed
- Anomalies are detected and displayed
- `docker compose down && up` recovers cleanly with all data intact

### Phase 2 — Pilot VPS

Docker Compose on a cloud VPS. Nginx terminates HTTPS. Postgres bound to `127.0.0.1`. Pilot clients redirect a single base URL — no other changes to their systems.

### Phase 3 — Production Scale

Multiple stateless Go instances behind a load balancer. Managed Postgres with automated backups. Prometheus + Grafana observability stack.

---

## What's Coming

The next capabilities in order of priority:

1. **PII detection and masking** — Kenyan National ID, M-PESA, KRA PIN detected and tokenised before any data leaves the organisation
2. **Redis-backed policy engine** — per-client token budgets and blocked topics, enforced on every request with no restart required
3. **Kill switch** — terminate the response stream mid-flight when a policy is violated
4. **Per-key identity** — attribute every request to a specific employee, service, or agent
5. **Webhook alerts** — Slack/PagerDuty when anomaly rules fire

---

## Tech Stack

| Layer | Technology |
|---|---|
| Proxy + API | Go — `httputil.ReverseProxy` + chi |
| Database | PostgreSQL 16 |
| Frontend | React + Vite + TypeScript |
| Containers | Docker + Docker Compose |
| Production routing | Nginx + Let's Encrypt |
| Observability (Phase 3) | Prometheus + Grafana |

---

## Security and Compliance

- Postgres is never exposed to the public internet
- Credentials managed via environment variables — never committed
- Audit logs are append-only
- Data residency enforced at the infrastructure level
- Architecture targets Kenya Data Protection Act 2019 compliance

---

## Git Workflow

```
main      — production-ready only
dev       — active development
feature/* — one branch per feature
fix/*     — bug fixes
```

Commit format: `service: what you did (present tense)`

```bash
git commit -m 'proxy: record latency and status per request'
git commit -m 'detector: add high token usage rule'
git commit -m 'api: add pagination to GET /requests'
git commit -m 'dashboard: auto-poll anomaly feed every 5s'
```

---

## Further Reading

| Document | Purpose |
|---|---|
| [`docs/HISTORY.md`](docs/HISTORY.md) | Origin story — why this exists and how it evolved |
| [`docs/vision.md`](docs/vision.md) | Vision — the full picture of what ARVIS is being built toward |
| [`DOCUMENTATION.md`](DOCUMENTATION.md) | Technical reference — architecture, data flow, API spec, deployment |

---

## Contact

Kevin — <kiplangatkevin335@gmail.com>
