# Feature Tracker

> What is built, what is in progress, what is planned.
>
> This file tracks every feature against the current Go + React architecture.
> For the story of how we got here — including the original Python SDK phase — see [`docs/HISTORY.md`](../docs/HISTORY.md).
> For the technical deep-dive on any item, see [`DOCUMENTATION.md`](../DOCUMENTATION.md).

---

## Status Key

| Symbol | Meaning |
|--------|---------|
| ✅ | Built and working |
| 🔨 | In progress |
| 📋 | Scoped and ready to build |
| 💡 | Proposed — not yet scoped |
| ❌ | Dropped |

---

## Phase 1 — Core Foundation
> Goal: prove the full pipeline works locally. Every item here must be ✅ before Phase 2 begins.

| # | Feature | Status | Location |
|---|---------|--------|----------|
| 1.1 | Go reverse proxy — intercepts every request to the LLM API | ✅ | `internal/proxy/proxy.go` |
| 1.2 | Request forwarder — proxies to upstream, captures latency and status | ✅ | `internal/proxy/forwarder.go` |
| 1.3 | Async audit log writer — goroutine writes to Postgres, never blocks response | ✅ | `internal/proxy/forwarder.go` |
| 1.4 | Anomaly rule: high latency (>10s) | ✅ | `internal/detector/rules.go` |
| 1.5 | Anomaly rule: high token usage (>3000 total tokens) | ✅ | `internal/detector/rules.go` |
| 1.6 | Anomaly rule: upstream 5xx error | ✅ | `internal/detector/rules.go` |
| 1.7 | Postgres schema — `requests` and `anomalies` tables | ✅ | `migrations/` |
| 1.8 | Store layer — `InsertRequest`, `ListRequests`, `InsertAnomaly`, `ListAnomalies` | ✅ | `internal/store/` |
| 1.9 | `GET /health` endpoint | ✅ | `internal/api/handlers/health.go` |
| 1.10 | `GET /requests?limit=N` endpoint | ✅ | `internal/api/handlers/requests.go` |
| 1.11 | `GET /anomalies?limit=N` endpoint | ✅ | `internal/api/handlers/anomalies.go` |
| 1.12 | chi router — wires all API routes | ✅ | `internal/api/router.go` |
| 1.13 | Environment-based config with sensible defaults | ✅ | `internal/config/config.go` |
| 1.14 | Single entry point — proxy goroutine + API server in one binary | ✅ | `cmd/arvis/main.go` |
| 1.15 | React dashboard — live request feed, polls every 5s | ✅ | `dashboard/src/pages/Dashboard.tsx` |
| 1.16 | React dashboard — anomaly feed | ✅ | `dashboard/src/components/AnomalyFeed.tsx` |
| 1.17 | React dashboard — aggregate stat cards (requests, tokens, avg latency) | ✅ | `dashboard/src/components/StatCards.tsx` |
| 1.18 | React dashboard — full request log page | ✅ | `dashboard/src/pages/Requests.tsx` |
| 1.19 | Vite proxy — `/api/*` forwarded to Go API in development | ✅ | `dashboard/vite.config.ts` |
| 1.20 | `docker-compose.yml` — single command starts Postgres + Go binary | ✅ | `docker-compose.yml` |
| 1.21 | Migrations auto-run on first Postgres container start | ✅ | `docker-compose.yml` → `migrations/` |
| 1.22 | Seed script — inserts demo data | ✅ | `scripts/seed.go` |

---

## Phase 2 — Make It Sellable
> Goal: features that turn the proof of concept into something a paying client would trust with live traffic.

| # | Feature | Status | Location |
|---|---------|--------|----------|
| 2.1 | PII detection — Kenyan National ID regex | 📋 | `internal/masking/pii.go` (to create) |
| 2.2 | PII detection — M-PESA phone number regex | 📋 | `internal/masking/pii.go` |
| 2.3 | PII detection — KRA PIN regex | 📋 | `internal/masking/pii.go` |
| 2.4 | PII tokenisation — replace value with `[KENYAN_ID_TOKEN_1]` before forwarding | 📋 | `internal/masking/pii.go` |
| 2.5 | PII token storage — encrypted original value in `pii_tokens` table | 📋 | `internal/store/pii.go` (to create) |
| 2.6 | Redis-backed policy engine — read rules per client key on every request | 📋 | `internal/policy/rules.go` (to create) |
| 2.7 | Daily token budget enforcement via Redis | 📋 | `internal/policy/rules.go` |
| 2.8 | Kill switch — terminate response stream mid-flight on policy violation | 📋 | `internal/proxy/forwarder.go` |
| 2.9 | Request body parsing — extract model + token counts from OpenAI JSON | 📋 | `internal/proxy/forwarder.go` |
| 2.10 | Per-client API key authentication (`X-Client-Key` header) | 📋 | `internal/proxy/middleware.go` |
| 2.11 | `POST /clients` — register a new client | 📋 | `internal/api/handlers/clients.go` (to create) |
| 2.12 | `PUT /clients/{id}/rules` — update rules, write to Redis | 📋 | `internal/api/handlers/clients.go` |
| 2.13 | `GET /report` — aggregate PII blocks, kill switch events, token spend | 📋 | `internal/api/handlers/reports.go` (to create) |
| 2.14 | VPS deployment — `docker-compose.prod.yml` with `restart: always` | 📋 | `docker-compose.prod.yml` (to create) |
| 2.15 | Nginx HTTPS termination — Let's Encrypt, subdomain routing | 📋 | `infra/nginx/nginx.conf` (to create) |
| 2.16 | Postgres and Redis bound to `127.0.0.1` in production | 📋 | Production security requirement |
| 2.17 | ODPC data processor registration submitted | 📋 | Kenya Data Protection Act 2019 |

---

## Phase 3 — Make It Sticky
> Goal: features that make ARVIS hard to remove once embedded in a client's infrastructure.

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 3.1 | Per-agent identity tagging — attribute each request to an employee, service, or agent | 💡 | New header `X-Agent-ID`, stored in `requests` |
| 3.2 | Webhook alerts — POST to Slack/Teams/PagerDuty on kill switch or budget breach | 💡 | Configurable per client in policy rules |
| 3.3 | Rate limiting per agent — requests-per-minute throttling per identity | 💡 | Redis counter, separate from daily budget |
| 3.4 | Prompt injection detection — detect and block jailbreak attempts before forwarding | 💡 | Pattern matching + optional LLM classifier |
| 3.5 | Cost attribution — token spend broken down by department, team, or agent | 💡 | Requires agent identity (3.1) first |
| 3.6 | Shadow mode — log everything, block nothing, for zero-risk client onboarding | 💡 | Per-client flag in policy rules |
| 3.7 | Semantic policy rules — meaning-based blocking beyond keyword matching | 💡 | Embedding similarity against blocked topics |
| 3.8 | Multi-LLM upstream routing — configurable target per request | 💡 | Anthropic, Mistral, any OpenAI-compatible endpoint |
| 3.9 | Load test gate — 1,000+ concurrent requests with <200ms added latency | 📋 | Phase 3 readiness gate |
| 3.10 | PII false positive rate below 1% against representative production text | 📋 | Phase 3 readiness gate |

---

## Phase 4 — Platform Scale
> Build these after the first enterprise contract is signed.

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 4.1 | Forensic replay — replay any historical session exactly as it occurred | 💡 | Requires full request/response body storage |
| 4.2 | Anomaly detection ML — baseline normal behaviour, alert on deviations | 💡 | Insider threat detection for AI traffic |
| 4.3 | Compliance packs — pre-built rule sets for CBK, GDPR, HIPAA, SOC 2 | 💡 | One-click compliance for regulated industries |
| 4.4 | Data residency routing — route sensitive prompts to on-prem or regional models | 💡 | Based on data classification, not just config |
| 4.5 | Sovereign LLM routing — intelligent routing between sovereign and cloud models | 💡 | For governments and telcos |
| 4.6 | Prometheus + Grafana monitoring stack | 💡 | Request throughput, latency percentiles, Redis memory, Postgres query time |
| 4.7 | Managed database migration — move Postgres to DBaaS with automated backups | 💡 | Required before any enterprise SLA commitment |
| 4.8 | Container orchestration — Kubernetes or ECS for zero-downtime deploys | 💡 | Extract proxy into its own deployable unit at this point |
| 4.9 | AI spend optimisation — analyse prompt patterns, suggest caching, auto-compress | 💡 | Save clients money while governing them |

---

## Dropped

| # | Feature | Reason |
|---|---------|--------|
| — | Python SDK / FastAPI management API | Replaced by Go-only backend. The SDK required manual adoption by every team. The proxy requires nothing. See [`docs/HISTORY.md`](../docs/HISTORY.md) for full context. |

---

## Notes

- Items marked 💡 are not yet scoped. Before moving one to 📋, define: what file it lives in, what the input/output contract is, and what the acceptance test looks like.
- Phase 1 is the `docker compose up` proof of concept. Nothing in Phase 2+ should be started until every Phase 1 item is ✅.
- Agent identity (3.1) is a dependency for cost attribution (3.5) and ML anomaly detection (4.2). Pull it forward to early Phase 3.
- PII masking (2.1–2.5) and the Redis policy engine (2.6–2.7) are the two features most likely to determine whether a pilot client converts to a paying client. Prioritise them first in Phase 2.
