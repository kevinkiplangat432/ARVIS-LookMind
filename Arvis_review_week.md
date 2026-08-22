<!--markdownlint-disable-->
# ARVIS Review Plan

## How to use this
Phase 1 is bottom-up: one file, one Go concept, in the order things
were actually built. Phase 2 is top-down: pick a feature, trace it
start to finish, and see if you can explain the whole path without
opening this doc. Do them together as you build the CLI in v0.10.0 —
each CLI command you build is the excuse to actually read the package
underneath it, not just skim it.

---

## Phase 1 — File by File (Go patterns)

### Foundations
| File | What to notice |
|---|---|
| `store/db.go` | Why `Ping` has to come *after* the pool is created, not before it's closed — the actual bug we fixed. What does `pgxpool.New` guarantee vs. what does `Ping` guarantee? |
| `handlers/health.go` | The isolated-sub-router pattern — every handler file returns its own `chi.Router`, mounted later. Why return a router instead of a bare `http.HandlerFunc`? |
| `config/config.go` | `Load()` vs `ResolveProviders()` being two separate calls — why does `migrate` never need providers, but `server` always does? |

### Store layer
| File | What to notice |
|---|---|
| `store/requests.go`, `store/anomalies.go`, `store/identities.go` | Every query is hand-written SQL, no ORM. `Scan` order must match `SELECT` column order exactly — where would this break silently if they drifted? |
| `store/audit.go` | The `$3::text IS NULL OR identity_id = $3` pattern — one query serving both "filtered" and "org-wide" callers. Could you write this two different ways and explain the tradeoff? |

### Concurrency patterns
| File | What to notice |
|---|---|
| `proxy/proxy.go` — `logAndDetect` | Why is this called with `go` and not called directly? What happens to the HTTP response if this function panics? |
| `detector/volume.go`, `detector/latency.go` | Both hold a `sync.Mutex` around a `map`. Why do they need one and `ContentRule` doesn't? |
| `commands/server.go` | `errgroup.Group` running two `http.ListenAndServe` calls — why does one server dying bring down the other? |

### Interfaces and abstraction
| File | What to notice |
|---|---|
| `detector/rule.go` | `SyncRule` vs `AsyncRule` — two interfaces instead of one, on purpose. What would go wrong if `ContentRule` implemented `SyncRule` instead? |
| `policy/check.go` — `MatchText` | Pulled out of `CheckTopics` specifically so `streaming.go` could reuse it. Find both call sites — what's actually shared, and what stayed different? |

### Auth and security-sensitive code
| File | What to notice |
|---|---|
| `auth/keys.go` | Why `crypto/rand`, never `math/rand`. Why is the raw key never stored, only its hash? |
| `proxy/auth.go` | Fails closed (missing key = rejected) — contrast with `policy.Check`'s Redis error path, which fails *open*. Can you explain, in your own words, why those two failure modes are opposite on purpose? |

### CLI plumbing
| File | What to notice |
|---|---|
| `commands/identity.go`, `commands/policy.go`, `commands/audit.go` | Every command self-registers via its own `init()`. Why does this mean `root.go` never has to be touched when a new command is added? |

---

## Phase 2 — Feature by Feature (end to end)

For each feature: list of files touched, in call order, then one
question. If you can answer the question out loud without opening
any file, you understand the feature.

### Feature: Identity & Attribution
`auth/keys.go` → `store/identities.go` → `commands/identity.go` (create)
`proxy/auth.go` → `store/identities.go` (lookup on every request)

**Question:** A bad prompt comes in. Trace, from the raw HTTP request,
exactly how ARVIS knows *which employee* sent it — every function it
passes through, in order.

### Feature: Multi-Provider Routing
`config/providers.go` → `config/config.go` (`ResolveProviders`)
`proxy/routing.go` (`resolveProvider`) → `proxy/proxy.go`

**Question:** `providers.yaml` has 5 providers, 40 total models. A
request comes in for `"model": "claude-3-opus"`. What's the exact
lookup path from raw JSON bytes to the provider's real API key being
attached to the outbound request?

### Feature: Observability (Logging)
`proxy/proxy.go` (`logAndDetect`) → `proxy/usage.go` → `store/requests.go`

**Question:** Why does token-usage extraction never fail the request,
even if the provider's response is malformed? What's the actual cost
of that leniency?

### Feature: Detection (Volume / Latency / Content)
`detector/rule.go` → `detector/volume.go` + `detector/latency.go` + `detector/content.go` → `detector/detector.go` → `store/anomalies.go`

**Question:** Which rules run *before* the response reaches the
caller, and which run *after*? What would break if that split were
reversed?

### Feature: Governance (Budgets & Blocked Topics)
`commands/policy.go` (configure) → `policy/policy.go` (Redis storage) → `policy/check.go` (`Check`) → `proxy/proxy.go` (enforcement point)

**Question:** A budget is exceeded mid-day. Walk through exactly what
happens to the *next* request that identity sends — every check it
hits, in order, before it would normally reach a provider.

### Feature: Kill Switch
`proxy/streaming.go` (`isStreaming`, `serveStreaming`) → `policy/check.go` (`MatchText`, reused) → `store/requests.go` (status 499)

**Question:** A blocked keyword lands exactly on a chunk boundary —
half in one SSE chunk, half in the next. What actually happens? Is
this a bug, or a known, accepted limitation? Where is that decision
written down?

### Feature: Audit Export
`commands/audit.go` → `audit/report.go` (`BuildReport`) → `store/audit.go` → `audit/pdf.go` (`RenderPDF`)

**Question:** The PDF shows "Terminated by kill switch: 3." Trace
backward — which table, which column, which status code produces
that number?

---

## What "done" looks like
Phase 1 done when you could explain any file in this list to someone
else, cold, without re-reading it first. Phase 2 done when you could
answer all seven questions above out loud, in order, without opening
the codebase. That's the actual bar for the dashboard being safe to
build on top of this.