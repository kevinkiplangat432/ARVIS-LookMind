# Sovereign AI Proxy

> **Version**: 0.1.0 (Alpha) · **Status**: Active Development  
> **Compliance**: Kenya Data Protection Act 2019 · Data residency enforced at the infrastructure level

## Overview

Sovereign AI Proxy is a transparent, real-time governance layer that sits between an organisation's systems and any LLM API. Every prompt and response passes through the proxy — PII is detected and masked before data leaves the organisation, policy rules are enforced before the request is forwarded, and a complete audit trail is written to Postgres for compliance and reporting.

The system is built as a monorepo containing three independent services: a high-performance Go proxy engine, a FastAPI management API, and a React control dashboard. They share a Postgres database and a Redis cache but never call each other's code directly — making each service independently deployable and scalable.

---

## Why This Exists

Enterprises deploying LLM-powered applications face a consistent set of risks:

- Sensitive data — customer identifiers, financial references, personal information — leaking into third-party AI APIs
- No visibility into what employees or internal systems are sending to LLMs
- No mechanism to enforce token budgets, block topics, or terminate unsafe responses in real time
- No audit trail for regulatory compliance or internal governance

Sovereign AI Proxy addresses all of these at the infrastructure level, requiring zero changes to existing applications. A company redirects their LLM base URL to the proxy endpoint — everything else works as before, but every request is now governed, logged, and auditable.

---

## Architecture

```
Client Application
       │
       ▼
┌──────────────────────────────────────────────┐
│            Go Proxy Engine (:8080)           │
│                                              │
│  1. Authenticate company API key             │
│  2. Read Redis — budget + policy rules       │
│  3. Scan request for PII → tokenise          │
│  4. Forward clean request to LLM API         │
│  5. Monitor response stream → kill switch    │
│  6. Write audit log to Postgres (async)      │
└──────────────────────────────────────────────┘
       │                         │
       ▼                         ▼
  LLM API                   PostgreSQL
  (OpenAI /               (audit_logs,
   Anthropic /             pii_tokens,
   any compatible)         companies)
                                │
                                ▼
                     ┌──────────────────────┐
                     │  FastAPI Dashboard   │
                     │  (:8000)             │
                     │  Auth · Rules        │
                     │  Logs · Reports      │
                     └──────────────────────┘
                                │
                                ▼
                     ┌──────────────────────┐
                     │   React Frontend     │
                     │   (:3000)            │
                     │  Live feed · Reports │
                     │  Kill switch · Rules │
                     └──────────────────────┘
```

**Service communication rules:**
- Go reads policy rules from Redis on every request (< 1ms, hot path)
- FastAPI writes rule updates to Redis — Go enforces them on the next request, no restart needed
- Go writes audit logs to Postgres asynchronously — never on the critical path
- FastAPI reads from Postgres to power the dashboard — it never touches the proxy hot path
- React calls FastAPI REST endpoints — standard JSON over HTTP
- Go and FastAPI share a database. They never call each other directly

---

## Repository Structure

```
sovereign-proxy/              ← monorepo root
│
├── proxy/                    ← Go engine (the core product)
│   ├── main.go               ← entry point, HTTP server
│   ├── go.mod / go.sum       ← dependencies
│   ├── interceptor/
│   │   └── proxy.go          ← reverse proxy logic (httputil.ReverseProxy)
│   ├── masking/
│   │   └── pii.go            ← Kenyan ID, M-PESA, KRA PIN detection + tokenisation
│   ├── policy/
│   │   └── rules.go          ← Redis-backed budget and rule enforcement
│   ├── audit/
│   │   └── logger.go         ← async audit log writer to Postgres
│   ├── stream/
│   │   └── monitor.go        ← token stream monitor, kill switch trigger
│   └── Dockerfile
│
├── dashboard/                ← FastAPI management API
│   ├── main.py               ← entry point, route registration
│   ├── requirements.txt
│   ├── routes/
│   │   ├── auth.py           ← POST /auth/login — issues JWT tokens
│   │   ├── companies.py      ← POST /company/new, PUT /company/rules
│   │   ├── logs.py           ← GET /logs — audit log reader
│   │   └── reports.py        ← GET /report — PII blocks, cost aggregates
│   ├── models/
│   │   ├── company.py        ← Pydantic: company data shape
│   │   └── log_entry.py      ← Pydantic: audit log row shape
│   ├── db/
│   │   └── connection.py     ← asyncpg connection pool
│   └── Dockerfile
│
├── frontend/                 ← React control dashboard
│   ├── src/
│   │   ├── App.jsx           ← root component, routing
│   │   ├── pages/
│   │   │   ├── Login.jsx     ← company manager login
│   │   │   ├── LiveFeed.jsx  ← real-time traffic table
│   │   │   ├── Report.jsx    ← compliance report view
│   │   │   └── Settings.jsx  ← budget and rule configuration
│   │   └── components/
│   │       ├── KillSwitchBadge.jsx  ← alert badge when kill switch fires
│   │       └── PIICounter.jsx       ← running count of PII blocks today
│   ├── package.json
│   └── Dockerfile
│
├── infra/                            ← shared infrastructure
│   ├── docker-compose.yml            ← single command local development
│   ├── docker-compose.prod.yml       ← production overrides
│   ├── nginx/
│   │   └── nginx.conf                ← routes traffic to Go and FastAPI by subdomain
│   ├── redis/
│   │   └── redis.conf                ← persistence settings, memory limits
│   └── postgres/
│       └── schema.sql                ← single source of truth for all tables
│
├── .env.example              ← all required environment variables with placeholders
├── .gitignore
├── Makefile                  ← make dev · make test · make deploy
└── README.md
```

---

## Core Features

### PII Masking and Tokenisation
Detects and replaces sensitive identifiers before any data is forwarded to an external LLM API. Supported patterns include Kenyan National ID numbers, M-PESA phone references, and KRA PIN formats. The original value is stored encrypted in Postgres; the LLM receives a reversible token such as `[KENYAN_ID_TOKEN_1]`. The real value never leaves the organisation's infrastructure.

### Real-Time Kill Switch
The Go engine reads the LLM response as a token stream. If a trigger condition is met mid-response — a policy violation, a blocked topic, a budget threshold breach — the connection is terminated immediately. The partial response is discarded and the event is logged with a violation reason.

### Redis-Backed Policy Engine
Each company's rules — daily token budget, allowed models, blocked topics — are stored in Redis and read on every request in under 1ms. When a company manager updates rules via the dashboard, FastAPI writes to Redis and the proxy enforces the new rules on the very next request with no restart required.

### Immutable Audit Logs
Every request the proxy handles produces an audit log row in Postgres: timestamp, company ID, PII detected flag, PII type, kill switch flag, tokens used, model, and violation reason. Logs are written asynchronously and are never modified after creation.

### Multi-Tenant Company Isolation
Each company is issued a unique API key. All logs, rules, PII tokens, and reports are scoped to that company's ID. One deployment serves multiple organisations with complete data isolation between tenants.

### Compliance Reporting
The FastAPI dashboard aggregates audit log data into compliance reports: total PII blocks, PII types detected, kill switch events, token spend by day, and estimated cost saved by blocking unsafe requests.

---

## Database Schema

Defined in `infra/postgres/schema.sql` — the single source of truth for all tables. Runs automatically on first Postgres start via Docker.

```sql
-- Companies registered on the platform
CREATE TABLE companies (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT NOT NULL,
    api_key       TEXT UNIQUE NOT NULL,       -- hashed value of X-Company-Key header
    budget_daily  INTEGER DEFAULT 10000,      -- max tokens per day
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

-- Every request the proxy handles
CREATE TABLE audit_logs (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id       UUID REFERENCES companies(id),
    timestamp        TIMESTAMPTZ DEFAULT NOW(),
    pii_detected     BOOLEAN DEFAULT FALSE,
    pii_type         TEXT,                    -- 'KENYAN_ID' | 'MPESA_PHONE' | 'KRA_PIN'
    kill_switch      BOOLEAN DEFAULT FALSE,
    tokens_used      INTEGER,
    model            TEXT,
    violation_reason TEXT
);

-- PII token map — original values never leave the organisation
CREATE TABLE pii_tokens (
    token          TEXT PRIMARY KEY,          -- '[KENYAN_ID_TOKEN_1]'
    original_value TEXT NOT NULL,             -- encrypted at rest
    company_id     UUID REFERENCES companies(id),
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

-- Per-company policy rules (also cached in Redis)
CREATE TABLE company_rules (
    company_id      UUID REFERENCES companies(id) PRIMARY KEY,
    blocked_topics  TEXT[],
    allowed_models  TEXT[],
    budget_daily    INTEGER DEFAULT 10000,
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for fast dashboard queries
CREATE INDEX idx_logs_company   ON audit_logs(company_id);
CREATE INDEX idx_logs_timestamp ON audit_logs(timestamp DESC);
CREATE INDEX idx_logs_pii       ON audit_logs(pii_detected) WHERE pii_detected = TRUE;
```

---

## Getting Started

### Prerequisites
- Docker and Docker Compose
- An API key for your target LLM provider

### Local Development

```bash
git clone https://github.com/your-org/sovereign-proxy.git
cd sovereign-proxy

cp .env.example .env
# Edit .env — add your LLM API key and set a strong JWT secret

docker-compose -f infra/docker-compose.yml up -d
```

| Service          | URL                     |
|------------------|-------------------------|
| Go Proxy         | http://localhost:8080   |
| FastAPI Dashboard| http://localhost:8000   |
| React Frontend   | http://localhost:3000   |

### Verify the proxy is running

```bash
curl http://localhost:8080/health
# {"status": "ok"}
```

### Test PII masking end-to-end

```bash
curl -X POST http://localhost:8080/v1/chat/completions \
  -H 'Content-Type: application/json' \
  -H 'X-Company-Key: test_key_123' \
  -d '{"messages":[{"role":"user","content":"My ID is 23456789"}]}'
```

Check Postgres — `audit_logs` should show `pii_detected = true`. The `pii_tokens` table should contain the token mapping. The value `23456789` should not appear anywhere in the forwarded request.

### Makefile shortcuts

```bash
make dev      # start all containers
make logs     # tail all container logs
make test     # run Go and Python test suites
make down     # stop all containers
make deploy   # deploy to production VPS
```

---

## Environment Variables

All variables are defined in `.env.example`. Copy to `.env` and fill in real values. The `.env` file is never committed to git.

```bash
# Go Proxy
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_API_KEY=your_llm_api_key_here
PROXY_PORT=8080
REDIS_URL=redis://redis:6379
POSTGRES_URL=postgres://user:pass@postgres:5432/sovereign

# FastAPI Dashboard
JWT_SECRET=change_this_to_a_long_random_string
JWT_EXPIRE_HOURS=24

# Postgres
POSTGRES_USER=user
POSTGRES_PASSWORD=change_in_production
POSTGRES_DB=sovereign

# Admin
ADMIN_API_KEY=your_internal_admin_key

# Phase 2+
DOMAIN=proxy.yourcompany.co.ke
SSL_EMAIL=you@yourcompany.co.ke
```

---

## Deployment Phases

### Phase 1 — Local Proof of Concept
All five containers run on a single machine via `docker-compose up`. The goal is to prove the full pipeline works end-to-end: a request containing a sensitive identifier enters the proxy, the identifier is masked, the clean request is forwarded to the LLM, and the event is logged and visible in the dashboard.

**Done when:**
- A request containing a Kenyan ID number reaches the proxy and the ID never appears in the upstream request
- The original value is stored encrypted in `pii_tokens`
- `audit_logs` shows `pii_detected = true`
- Redis correctly blocks a request when a company exceeds their daily token budget
- The React dashboard shows the live feed including the PII block event
- `docker-compose down && docker-compose up` recovers cleanly with all data intact

### Phase 2 — Pilot VPS Deployment
The same Docker Compose setup runs on a cloud VPS hosted in-region. Nginx terminates HTTPS and routes subdomains to the correct containers. Postgres and Redis are bound to `127.0.0.1` only — not accessible from the internet. Pilot clients point their existing applications at the proxy endpoint with a single base URL change.

**Done when:**
- The proxy is reachable over HTTPS on a production domain with a valid TLS certificate
- At least one external client is routing live traffic through the proxy
- Postgres and Redis ports are not internet-accessible
- Kill switch is tested and terminates a response mid-stream correctly
- A compliance report has been generated and delivered to a pilot client
- Data processor registration with the relevant data protection authority is submitted

### Phase 3 — Production Scale
Multiple stateless Go proxy instances run behind a load balancer. Postgres moves to a managed database service with automated backups and replication. Container orchestration handles zero-downtime deployments and automatic restarts. A Prometheus and Grafana monitoring stack provides real-time visibility into request throughput, latency percentiles, and infrastructure health.

**Done when:**
- Load testing confirms the proxy handles 1,000+ concurrent requests with < 200ms added latency
- PII false positive rate is below 1% against representative production text
- Managed database with automated backups and failover is in place
- Monitoring and alerting stack is live
- A documented case study exists: volume of PII blocks, cost saved, compliance events over a defined period

---

## Git Workflow

```
main        — production-ready code only, never commit directly
dev         — active development, merge to main when a feature is complete
feature/*   — one branch per feature, e.g. feature/pii-masking
fix/*       — bug fix branches, e.g. fix/kill-switch-latency
```

**Commit message format:**

```bash
git commit -m 'proxy: add Kenyan national ID regex to PII masker'
git commit -m 'proxy: implement Redis budget check before forwarding'
git commit -m 'dashboard: add GET /report endpoint with weekly summary'
git commit -m 'infra: add docker-compose.prod.yml with Nginx config'
git commit -m 'frontend: add KillSwitchBadge to live feed table'
git commit -m 'schema: add pii_tokens table with company_id foreign key'
```

Pattern: `service: what you did (present tense, concise)`

---

## Security and Compliance

- PII original values are encrypted at rest in Postgres
- Audit logs are immutable — no update or delete operations are permitted on `audit_logs`
- Postgres and Redis are never exposed to the public internet in any deployment phase
- All dashboard access is authenticated via short-lived JWT tokens
- Credentials are managed via environment variables and never committed to version control
- Data residency is enforced at the infrastructure level — all data remains within the deployment region
- Architecture supports registration with national data protection authorities

---

## Roadmap

- **Per-agent identity tagging** — attribute each request to a specific employee, service, or agent
- **Prompt injection detection** — detect and block jailbreak and prompt injection attempts before they reach the LLM
- **Cost attribution** — break down token spend by department, team, or agent for internal chargeback
- **Webhook alerts** — fire to Slack, Teams, or PagerDuty when kill switch triggers or budget threshold is hit
- **Semantic policy rules** — meaning-based blocking beyond keyword matching
- **Shadow mode** — proxy logs everything and blocks nothing, for zero-risk enterprise onboarding
- **Data residency routing** — route sensitive prompt categories to on-premises or regional models based on data classification
- **Compliance packs** — pre-built rule sets for specific regulatory frameworks, deployable in one step
- **Forensic replay** — replay any historical session exactly as it occurred for incident investigation
- **Anomaly detection** — baseline normal agent behaviour and alert on statistical deviations

---

## Tech Stack

| Layer | Technology |
|---|---|
| Proxy engine | Go — `httputil.ReverseProxy` |
| Policy cache | Redis 7 |
| Audit storage | PostgreSQL 16 |
| Management API | FastAPI + asyncpg |
| Frontend | React + Vite |
| Container runtime | Docker + Docker Compose |
| Production routing | Nginx + Let's Encrypt |
| Observability (Phase 3) | Prometheus + Grafana |

---

## Contributing

1. Fork the repository
2. Create a feature branch — `git checkout -b feature/your-feature`
3. Commit using the message format described above
4. Push and open a pull request against `dev`

---

## License

MIT License — see [LICENSE](./LICENSE) for details.

---

## Contact

- GitHub Issues: [open an issue](https://github.com/your-org/sovereign-proxy/issues)
- Email: kiplangatkevin335@gmail.com
