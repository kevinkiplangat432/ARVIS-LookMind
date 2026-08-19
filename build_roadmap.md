<!--Markdownlint-disable-->
# ARVIS Finishing Roadmap — v0.6.0 → v0.7.0

## Where we actually are (honest framing)
Of the four pillars — Observability, Auditing, Governance, Compliance —
ARVIS today is: Observability (mostly built), early Auditing (immutable
requests/anomalies tables exist). Governance (policy engine, kill switch)
and Compliance (mapping to a named standard) are ahead of us. This
roadmap gets Observability + Auditing to a real, end-to-end, on-prem-
deployable state. Governance/Compliance/Dashboard/PDF/ML are later
chapters, not tonight.

## Phase order and why
Database schema comes *before* proxy/detector this time — you're right
that it doesn't capture enough yet (no identity linkage, no provider
column, no anomaly categorization). Building the proxy against next
week's schema instead of tonight's would mean redoing it. Identity
comes right after schema, since the proxy needs something to attribute
a request to before it can log one meaningfully.

---

### Phase 0 — Bugfix (on `main`, no branch needed)
- Fix `db.go`: pool gets closed before it's pinged.
- **~1 commit.**

### Phase 1 — `feature/handlers-finish`
- Finish `request.go` (wire to `store.ListRequests`, `?limit` param)
- Add `anomalies.go` handler (same isolated-router pattern as `health.go`)
- Mount both in `router.go`
- **~4 commits** (one per file + router mount)

### Phase 2 — `feature/schema-v2`
- Migration 000003: `identities` table (id, name, key_hash, created_at)
- Migration 000004: alter `requests` — add `identity_id` FK, `provider`
  column, `model_requested` vs `model_used` if they can differ
- Migration 000005: alter `anomalies` — add `category` (volume/content/
  latency), `severity`, `status` (open/reviewed/dismissed)
- Update `store.Request` / `store.Anomaly` structs + insert/list queries
  to match
- **~6 commits**

### Phase 3 — `feature/identity`
- `identity create [name]` Cobra command
- `internal/store/identities.go` — create/lookup by hashed key
- Raw key printed once at creation, never stored in plaintext
- **~3 commits**

### Phase 4 — `feature/provider-config`
- `providers.yaml` — up to 20 providers, `${ENV_VAR}` interpolation for
  keys, never raw keys in the file
- `config` package refactor: replace single `TargetURL`/`APIKey` with a
  loaded provider list
- **~3 commits**

### Phase 5 — `feature/proxy-core`
- `internal/proxy/proxy.go` — the actual reverse proxy on `cfg.ProxyAddr`
- Auth: caller sends their ARVIS-issued key, proxy resolves `identity_id`
- Routing: read `model` from the request body, map to a provider, proxy
  injects the *real* provider key (caller never sees or needs it)
- **~4 commits**

### Phase 6 — `feature/proxy-logging`
- Capture latency, status, token counts per call
- Write via `store.InsertRequest`, now tagged with `identity_id` +
  `provider`
- **~2 commits**

### Phase 7 — `feature/detector`
- `internal/detector/rules.go` — `Rule` interface, `Check() []Flag`
- `volume.go` — sync, in-memory rate/window counters, runs before
  forwarding
- `latency.go` — async, rolling-stats outlier detection
- `content.go` — async, regex-based PII rules (Kenyan ID, KRA PIN,
  M-PESA to start — this is where the local differentiator lives)
- `detector.go` — orchestrator: sync rules block nothing but must
  finish before response; async rules run in a goroutine after the
  response is already sent back to the caller
- Wire into the proxy pipeline, writes via `store.InsertAnomaly`
- **~7 commits**

### Phase 8 — `feature/server-wiring`
- `server.go`: run proxy (8080) and API (8081) concurrently
- **~2 commits**

### Phase 9 — `feature/cli-hardening`
- `test.go`: real test suite — unit tests across packages + integration
  tests against a live test DB + a smoke pass hitting the running
  proxy/API
- `migrate.go`: absolute migrations path (not relative to cwd), add
  `status`/`force` subcommands, handle dirty state
- **~4 commits**

### → Tag v0.7.0
Once Phase 9 lands and the full pipeline works end to end (request in,
routed, logged, detected, attributable to an identity), that's the
natural v0.7.0 cut.

**Total: ~29 commits** — close enough to your 30 that the real number
will land wherever it lands once we're actually writing.

---

## Later chapters (not tonight, no branch yet)
- PDF audit export (possible digital signing later)
- Governance: Redis-backed policy engine + kill switch
- Enterprise dashboard
- GNN/Bayesian/RL modeling layer (feeds from [[arvis-ml-research]] docs)
- Federated cross-institution learning (long-horizon research bet)