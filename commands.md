<!--markdownlint-disable-->
# ARVIS Test Suite & CLI Reference

## Building the Binary

Always build from the repo root:

```bash
go build -o arvis ./cmd/arvis/
````

Never run:

```bash
go run cmd/arvis/main.go
```

Use the compiled binary or:

```bash
go run ./cmd/arvis/
```

---

## Environment Variables

| Variable         | Description                         | Default          |
| ---------------- | ----------------------------------- | ---------------- |
| `DATABASE_URL`   | PostgreSQL connection string        | Local ARVIS DB   |
| `REDIS_ADDR`     | Redis address for the policy engine | `localhost:6379` |
| `PROVIDERS_FILE` | Path to `providers.yaml`            | `providers.yaml` |
| `PROXY_ADDR`     | Proxy listen address                | `:8080`          |
| `API_ADDR`       | API listen address                  | `:8081`          |
| `MAX_TOKENS`     | Default token-cap fallback          | `4096`           |

The server will refuse to start without:

* A resolvable `providers.yaml`
* A reachable PostgreSQL database

It now fails loudly on both rather than starting half-configured.

Redis being unreachable **does not** stop the server. Policy enforcement fails open and logs the error instead.

See `internal/proxy/proxy.go`'s Redis error handling if this behavior needs to be revisited.

---

## CLI Commands

### Interactive Mode

Run without arguments:

```bash
./arvis
```

This shows a menu. Enter `1`, `2`, or `3` to select a mode.

### Core Subcommands

```bash
./arvis server
```

Start the proxy on `:8080` and API on `:8081`.

```bash
./arvis test
```

Run:

* `go vet`
* `go test -race`
* Coverage
* Tests with `-count=1`
* All packages

```bash
./arvis migrate
```

Run migrations up by default.

```bash
./arvis migrate up
```

Run migrations up explicitly.

```bash
./arvis migrate down
```

Roll back all migrations.

```bash
./arvis migrate status
```

Show the current migration version and dirty state.

```bash
./arvis migrate force [ver]
```

Force the migration version for dirty-state recovery.

```bash
./arvis migrate --path <dir>
```

Use a migrations directory other than `./migrations`.

```bash
./arvis --help
```

Show all commands.

```bash
./arvis server --help
```

Show help for the server command.

---

## Identity

Create a new identity and issue a one-time key:

```bash
./arvis identity create "[name]"
```

---

## Governance — Policy Engine

### Set Budgets

Identity-level budget:

```bash
./arvis policy set-budget \
  --scope identity \
  --id <id> \
  --daily <n> \
  --monthly <n> \
  --enabled
```

Provider-level budget:

```bash
./arvis policy set-budget \
  --scope provider \
  --id openai \
  --daily <n>
```

Global budget:

```bash
./arvis policy set-budget \
  --scope global \
  --daily <n>
```

### Topic Blocking

Block a topic:

```bash
./arvis policy block-topic [key]
```

See `list-topics` for valid topic keys.

Unblock a topic:

```bash
./arvis policy unblock-topic [key]
```

List all blockable topic categories and their sources:

```bash
./arvis policy list-topics
```

Show currently blocked topics:

```bash
./arvis policy show
```

Show every configured budget and current usage:

```bash
./arvis policy budgets
```

---

## Observability

### List Requests

List all requests:

```bash
./arvis requests list
```

Filter by identity:

```bash
./arvis requests list --identity <id>
```

Limit the number of results:

```bash
./arvis requests list --limit N
```

Both filters can be combined:

```bash
./arvis requests list --identity <id> --limit N
```

### List Anomalies

List anomalies:

```bash
./arvis anomalies list
```

Filter by category:

```bash
./arvis anomalies list --category X
```

Filter by severity:

```bash
./arvis anomalies list --severity X
```

Filter by status:

```bash
./arvis anomalies list --status X
```

Limit results:

```bash
./arvis anomalies list --limit N
```

Filters can be combined:

```bash
./arvis anomalies list \
  --category X \
  --severity X \
  --status X \
  --limit N
```

### Resolve an Anomaly

Mark an anomaly as reviewed:

```bash
./arvis anomalies resolve [id] --status reviewed
```

Dismiss an anomaly:

```bash
./arvis anomalies resolve [id] --status dismissed
```

---

## Audit Export

Export an audit report:

```bash
./arvis audit export \
  [--identity <id>] \
  [--from YYYY-MM-DD] \
  [--to YYYY-MM-DD] \
  [--out path.pdf]
```

### Default Behavior

| Option       | Default                                                |
| ------------ | ------------------------------------------------------ |
| `--identity` | Organization-wide report                               |
| `--from`     | 30 days before `--to`                                  |
| `--to`       | Today                                                  |
| `--out`      | `arvis-audit-<timestamp>.pdf` in the current directory |

For example:

```bash
./arvis audit export
```

Organization-wide report using the default date range and output filename.

```bash
./arvis audit export --identity <id>
```

Export an identity-specific report.

```bash
./arvis audit export \
  --from 2026-08-01 \
  --to 2026-08-22 \
  --out report.pdf
```

Export a report for a specific date range and output path.

---

# Setup

Run this once, or after adding new migrations.

## Create the Test Database

If `arvis_test` does not exist yet:

```bash
sudo -u postgres psql -c "CREATE DATABASE arvis_test OWNER arvis;"
```

## Apply Migrations

There are currently five migrations covering:

1. Requests
2. Anomalies
3. Identities
4. Requests alterations
5. Anomalies alterations

Apply them to the test database:

```bash
DATABASE_URL="postgres://arvis:arvis@localhost:5432/arvis_test?sslmode=disable" \
./arvis migrate up
```

After adding new migration files, run the same command again:

```bash
DATABASE_URL="postgres://arvis:arvis@localhost:5432/arvis_test?sslmode=disable" \
./arvis migrate up
```

---

## Redis Setup

Redis needs to be running locally for anything touching the policy package, including:

* Budgets
* Blocked topics
* Kill switch topic lookup

There is currently **no separate test Redis instance** like there is for PostgreSQL.

### Current Limitation

Real Redis isolation should be added before policy tests are written.

Possible approaches:

* Use a dedicated Redis DB index for tests.
* Flush test keys between test runs.
* Use another Redis instance/container for testing.

This prevents tests from leaving state behind that could affect subsequent tests.

---

# Running Tests

## Using the ARVIS Binary

Recommended:

```bash
TEST_DATABASE_URL="postgres://arvis:arvis@localhost:5432/arvis_test?sslmode=disable" \
./arvis test
```

This runs the complete test workflow.

---

## Using Make

Run all tests:

```bash
make test
```

Run unit tests only:

```bash
make test-unit
```

Unit tests do not require a database.

Run integration tests only:

```bash
make test-integration
```

Integration tests require `arvis_test`.

---

## Using `go test` Directly

Run detector tests:

```bash
go test ./internal/detector/... -v
```

Run store tests:

```bash
TEST_DATABASE_URL="postgres://arvis:arvis@localhost:5432/arvis_test?sslmode=disable" \
go test ./internal/store/... -v
```

> **Always run `make test` before every commit.**

---

# Test Coverage

## Existing Tests

### `internal/detector/rules_test.go`

Covers:

* Anomaly detection rules
* Boundary values
* Multi-rule firing

### `internal/store/requests_test.go`

Covers:

* Request insertion
* Duplicate ID rejection
* Ordering
* Result limits
* Empty slices

### `internal/store/anomalies_test.go`

Covers:

* Anomaly insertion
* Foreign-key enforcement
* Ordering
* Empty slices

---

## Not Yet Covered

These are **real testing gaps**, not items to ignore.

### `internal/policy/*`

Needs tests for:

* Budget checks
* Usage counters
* Topic matching
* `MatchText`

`MatchText` especially deserves boundary tests around the chunk-split limitation identified in `streaming.go`.

### `internal/proxy/*`

Needs tests for:

* Authentication
* Routing
* Streaming relay
* Kill-switch termination

### `internal/audit/*`

Needs tests for:

* `BuildReport` correctness
* PDF generation

### `internal/store/identities.go`

Needs tests for:

* Key-hash lookup
* ID lookup

### `internal/config/providers.go`

Needs tests for:

* Environment-variable interpolation
* Duplicate detection
* 20-provider cap

---

## Phase 2 Review Alignment

This list exactly follows the Phase 2 review guide's feature list in the same order.

Each untested feature in the review guide corresponds to an untested package or component here.

Tests should be written **as each feature is reviewed**, rather than postponing all testing until the end of the review.

---

# Databases

| Database     | Purpose              | Test Access                           |
| ------------ | -------------------- | ------------------------------------- |
| `arvis`      | Development database | **Never touched by tests**            |
| `arvis_test` | Test database        | Wiped and rewritten by the test suite |

Tests must use `arvis_test` and must never operate against the development database.

---

# Redis

## Development Redis

```text
localhost:6379
```

Used by:

* Live server
* CLI
* Policy engine

## Test Redis

A separate test Redis instance does **not yet exist**.

> **TODO:** Add real Redis test isolation before policy tests are implemented.

Recommended options include:

* Dedicated Redis DB index
* Dedicated Redis instance/container
* Explicit test-key cleanup

---

# Working Directory

## Always Run From the Repository Root

All commands assume the repository is located at:

```text
~/Documents/development/ARVIS
```

Run commands from the repository root:

```bash
cd ~/Documents/development/ARVIS
```

The binary reads the following paths relative to the current working directory:

```text
migrations/
providers.yaml
```

Running the binary from another directory can therefore cause:

* Migration-loading errors
* Provider-loading errors

### Running From Another Directory

If you genuinely need to run migrations from elsewhere, specify the migration directory explicitly:

```bash
./arvis migrate --path <dir>
```

For normal development and testing, **always run ARVIS from the repository root**.

```
```
