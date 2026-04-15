# Feature Tracker

> Track what is built, what is in progress, and what is planned.  
> Update status as work progresses. Phases are not strict — a Phase 3 feature can be pulled forward if there is demand.

---

## Status Key

| Symbol | Meaning |
|--------|---------|
| ✅ | Built and tested |
| 🔨 | In progress |
| 📋 | Planned — scoped and ready to build |
| 💡 | Proposed — not yet scoped |
| ❌ | Dropped |

---

## Phase 1 — Core Foundation
> Goal: prove the full pipeline works locally. Every item here must be ✅ before Phase 2.

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 1.1 | Go reverse proxy — intercepts requests to LLM API | 📋 | `proxy/interceptor/proxy.go` |
| 1.2 | Company API key authentication (`X-Company-Key` header) | 📋 | `proxy/main.go` |
| 1.3 | PII detection — Kenyan National ID regex | 📋 | `proxy/masking/pii.go` |
| 1.4 | PII detection — M-PESA phone number regex | 📋 | `proxy/masking/pii.go` |
| 1.5 | PII detection — KRA PIN regex | 📋 | `proxy/masking/pii.go` |
| 1.6 | PII tokenisation — replace value with `[KENYAN_ID_TOKEN_1]` | 📋 | `proxy/masking/pii.go` |
| 1.7 | PII token storage — encrypted original value in `pii_tokens` table | 📋 | `proxy/audit/logger.go` |
| 1.8 | Redis-backed policy engine — read rules per company on every request | 📋 | `proxy/policy/rules.go` |
| 1.9 | Daily token budget enforcement via Redis | 📋 | `proxy/policy/rules.go` |
| 1.10 | Token stream monitor — read LLM response as stream | 📋 | `proxy/stream/monitor.go` |
| 1.11 | Kill switch — terminate connection mid-stream on trigger | 📋 | `proxy/stream/monitor.go` |
| 1.12 | Async audit log writer to Postgres | 📋 | `proxy/audit/logger.go` |
| 1.13 | Postgres schema — `companies`, `audit_logs`, `pii_tokens`, `company_rules` | 📋 | `infra/postgres/schema.sql` |
| 1.14 | FastAPI `/auth/login` — issues JWT tokens | 📋 | `dashboard/routes/auth.py` |
| 1.15 | FastAPI `POST /company/new` — register a company | 📋 | `dashboard/routes/companies.py` |
| 1.16 | FastAPI `PUT /company/rules` — update rules, write to Redis | 📋 | `dashboard/routes/companies.py` |
| 1.17 | FastAPI `GET /logs` — read audit logs from Postgres | 📋 | `dashboard/routes/logs.py` |
| 1.18 | FastAPI `GET /report` — aggregate PII blocks and cost data | 📋 | `dashboard/routes/reports.py` |
| 1.19 | React login page | 📋 | `frontend/src/pages/Login.jsx` |
| 1.20 | React live feed — real-time traffic table | 📋 | `frontend/src/pages/LiveFeed.jsx` |
| 1.21 | React compliance report view | 📋 | `frontend/src/pages/Report.jsx` |
| 1.22 | React settings page — budget and rule configuration | 📋 | `frontend/src/pages/Settings.jsx` |
| 1.23 | KillSwitchBadge component — red alert when kill switch fires | 📋 | `frontend/src/components/KillSwitchBadge.jsx` |
| 1.24 | PIICounter component — running count of blocks today | 📋 | `frontend/src/components/PIICounter.jsx` |
| 1.25 | `docker-compose.yml` — single command starts all 5 containers | 📋 | `infra/docker-compose.yml` |
| 1.26 | `/health` endpoint on Go proxy | 📋 | `proxy/main.go` |

---

## Phase 2 — Make It Sellable
> Goal: features that turn the proof of concept into a product a paying company would trust.

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 2.1 | Multi-LLM support — configurable upstream target via `OPENAI_BASE_URL` | 📋 | Anthropic, Mistral, any OpenAI-compatible endpoint |
| 2.2 | Per-agent identity tagging — attribute each request to an employee, service, or agent ID | 💡 | New header `X-Agent-ID`, stored in `audit_logs` |
| 2.3 | Prompt injection detection — detect and block jailbreak attempts before forwarding | 💡 | Pattern matching + optional LLM-based classifier |
| 2.4 | Cost attribution — token spend broken down by department, team, or agent | 💡 | Requires agent identity (2.2) to be built first |
| 2.5 | Webhook alerts — fire to Slack / Teams / PagerDuty on kill switch or budget breach | 💡 | Configurable per company in `company_rules` |
| 2.6 | Rate limiting per agent — requests-per-minute throttling per identity | 💡 | Redis counter, separate from daily budget |
| 2.7 | VPS deployment — `docker-compose.prod.yml` with `restart: always` | 📋 | `infra/docker-compose.prod.yml` |
| 2.8 | Nginx HTTPS termination — Let's Encrypt certificate, subdomain routing | 📋 | `infra/nginx/nginx.conf` |
| 2.9 | Postgres and Redis bound to `127.0.0.1` — not internet-accessible | 📋 | Production security requirement |
| 2.10 | ODPC data processor registration submitted | 📋 | Compliance — Kenya Data Protection Act 2019 |

---

## Phase 3 — Make It Sticky
> Goal: features that make the proxy hard to remove once embedded in a company's infrastructure.

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 3.1 | Policy-as-code — companies define rules in YAML committed to their own repos | 💡 | GitOps for AI governance |
| 3.2 | Semantic policy rules — meaning-based blocking beyond keyword matching | 💡 | e.g. block prompts requesting competitor analysis |
| 3.3 | Shadow mode — proxy logs everything, blocks nothing, for zero-risk onboarding | 💡 | Per-company flag in `company_rules` |
| 3.4 | Data residency routing — route sensitive prompts to on-prem or regional models | 💡 | Based on data classification tag on the request |
| 3.5 | Anomaly detection — baseline normal agent behaviour, alert on deviations | 💡 | Insider threat detection for AI traffic |
| 3.6 | Forensic replay — replay any historical session exactly as it occurred | 💡 | Requires full request/response body storage |
| 3.7 | Load testing baseline — proxy handles 1,000+ concurrent requests < 200ms added latency | 📋 | Phase 3 readiness gate |
| 3.8 | PII false positive rate below 1% — validated against representative production text | 📋 | Phase 3 readiness gate |

---

## Phase 4 — Platform Scale
> Goal: features that expand the proxy from a tool into a platform. Build these after the first enterprise contract.

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 4.1 | Compliance packs — pre-built rule sets for CBK, GDPR, HIPAA, SOC 2 | 💡 | One-click compliance for regulated industries |
| 4.2 | Agent marketplace governance — govern third-party AI agents regardless of vendor | 💡 | Proxy becomes the universal trust layer |
| 4.3 | Federated proxy network — multiple nodes across regions, centrally managed | 💡 | Sell to governments and telcos |
| 4.4 | AI spend optimisation — analyse prompt patterns, suggest caching, auto-compress prompts | 💡 | Save companies money while governing them |
| 4.5 | Sovereign LLM routing — intelligent routing between sovereign and cloud LLMs | 💡 | Based on data classification, not just config |
| 4.6 | Prometheus + Grafana monitoring stack | 💡 | `infra/` — request throughput, latency, Redis memory, Postgres query time |
| 4.7 | Managed database migration — move Postgres to DBaaS with automated backups | 💡 | Required before any enterprise SLA commitment |
| 4.8 | Container orchestration — Docker Swarm or Kubernetes for zero-downtime deploys | 💡 | Extract Go proxy into its own repo at this point |

---

## Dropped / Deferred

| # | Feature | Reason |
|---|---------|--------|
| — | — | — |

---

## Notes

- Features marked 💡 are not yet scoped. Before moving one to 📋, define: what file it lives in, what the input/output is, and what the acceptance test looks like.
- The Phase 1 checklist maps directly to the `docker-compose up` proof of concept. Nothing in Phase 2+ should be started until every Phase 1 item is ✅.
- Agent identity (2.2) is a dependency for cost attribution (2.4) and anomaly detection (3.5). Build it early in Phase 2.
