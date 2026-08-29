<!--markdownlint-disable-->
# ARVIS — Comprehensive Developer & Operations Guide

This document consolidates the full CLI reference, test suite instructions, and Docker workflow for the ARVIS proxy server.  
It reflects the fixes applied during the setup: database connection, entrypoint script, port binding, and environment separation.

---

## 1. Repository Structure & Build

### 1.1 Always work from the repository root

```bash
cd ~/Documents/development/startup/ARVIS   # adjust to your path
```

The binary reads `migrations/` and `providers.yaml` relative to the current working directory.

### 1.2 Building the binary

```bash
go build -o arvis ./cmd/arvis/
```

**Never run** `go run cmd/arvis/main.go` – use the compiled binary or:

```bash
go run ./cmd/arvis/
```

---

## 2. Environment Variables

| Variable         | Description                               | Default                         |
| ---------------- | ----------------------------------------- | ------------------------------- |
| `DATABASE_URL`   | PostgreSQL connection string              | `postgres://arvis:arvis@localhost:5432/arvis?sslmode=disable` |
| `REDIS_ADDR`     | Redis address for the policy engine       | `localhost:6379`                |
| `PROVIDERS_FILE` | Path to `providers.yaml`                  | `providers.yaml`                |
| `PROXY_ADDR`     | Proxy listen address                      | `:8080`                         |
| `API_ADDR`       | API listen address                        | `:8081`                         |
| `MAX_TOKENS`     | Default token‑cap fallback                | `4096`                          |

**Important:** The server refuses to start without a resolvable `providers.yaml` and a reachable PostgreSQL database.  
Redis being unreachable does **not** stop the server – policy enforcement fails open and logs the error.

---

## 3. Docker Setup (Production / Development Stack)

### 3.1 Files

- `Dockerfile` – multi‑stage build (Go builder + Alpine runtime)
- `docker-compose.yml` – defines `postgres`, `redis`, and `arvis` services
- `docker-entrypoint.sh` – runs `migrate up` then `exec ./arvis server`

### 3.2 Docker Compose services

| Service  | Image              | Ports (host:container) | Purpose                         |
| -------- | ------------------ | ---------------------- | ------------------------------- |
| postgres | postgres:16-alpine | `127.0.0.1:5432:5432`  | Main database                   |
| redis    | redis:7-alpine     | `127.0.0.1:6379:6379`  | Policy engine cache             |
| arvis    | built from Dockerfile | `8080:8080`, `8081:8081` | The proxy + API server          |

**Security:** Ports are bound to `127.0.0.1` only, so they are not exposed to external networks.  
The `.env` file is mounted as environment variables for the `arvis` container.

### 3.3 Common Docker commands

| Command                                                         | Description                                                                                  |
| --------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| `docker compose build`                                          | Builds the image(s) without starting containers.                                             |
| `docker compose up`                                             | Starts all services, streams combined logs to terminal.                                      |
| `docker compose up -d`                                          | Starts in detached mode (background).                                                        |
| `docker compose up --build`                                     | Rebuilds and starts after code changes.                                                      |
| `docker compose down`                                           | Stops and removes containers, preserves volumes (data).                                      |
| `docker compose down -v`                                        | Stops, removes containers **and deletes volumes** (wipes DB).                                |
| `docker compose stop` / `start`                                 | Pauses / resumes containers without removing them.                                           |
| `docker compose restart arvis`                                  | Restarts only the ARVIS container (useful after config changes).                             |
| `docker compose ps`                                             | Shows status and port mappings of all containers.                                            |
| `docker compose logs -f arvis`                                  | Follow logs from the ARVIS container only.                                                   |
| `docker compose exec arvis sh`                                  | Opens an interactive shell inside the running ARVIS container.                               |
| `docker compose exec postgres psql -U arvis -d arvis`           | Opens `psql` prompt inside Postgres container.                                               |
| `docker compose run --rm arvis ./arvis <command>`               | Runs a one‑off command in a new container (does not touch the running one).                 |

---

## 4. Two Ways to Run CLI Commands

Because the database is available on `localhost:5432` (thanks to port mapping), you have **two equivalent ways** to execute any CLI command:

| Method                                                          | Command Example                                | When to use                                       |
| --------------------------------------------------------------- | ---------------------------------------------- | ------------------------------------------------- |
| **Host binary** (on your machine)                               | `./arvis identity create "Alice"`              | Fast, no `docker` overhead, daily development    |
| **Inside the container** (`docker compose exec`)                | `docker compose exec arvis ./arvis ...`        | To test exactly the container environment        |

### 4.1 Host binary prerequisites

- Set `DATABASE_URL` in your shell (or export it in `~/.bashrc`):
  ```bash
  export DATABASE_URL="postgres://arvis:arvis@localhost:5432/arvis?sslmode=disable"
  ```
- Ensure the Docker stack is running (`docker compose up -d`) so Postgres is reachable.

### 4.2 Inside‑container execution

The container already has the correct `DATABASE_URL` (pointing to `postgres` service name), so no extra setup is needed. Just run:

```bash
docker compose exec arvis ./arvis identity create "Alice"
```

---

## 5. CLI Commands (Full Reference)

### 5.1 Interactive Mode

Run without arguments:

```bash
./arvis
```

Shows a menu. Enter `1`, `2`, or `3` to select a mode.

### 5.2 Server

```bash
./arvis server
```

Starts the proxy on `:8080` and API on `:8081`.

### 5.3 Migrations

| Command                                         | Description                            |
| ----------------------------------------------- | -------------------------------------- |
| `./arvis migrate`                               | Runs migrations up (default).          |
| `./arvis migrate up`                            | Explicitly runs `up`.                  |
| `./arvis migrate down`                          | Rolls back **all** migrations.         |
| `./arvis migrate status`                        | Shows current version and dirty state. |
| `./arvis migrate force <ver>`                   | Forces a version (for dirty recovery). |
| `./arvis migrate --path <dir>`                  | Uses a custom migration directory.     |

### 5.4 Identity

```bash
./arvis identity create "[name]"
```

Creates a new identity and outputs a one‑time key.

### 5.5 Governance – Policy Engine

**Budgets**  
Identity‑level:

```bash
./arvis policy set-budget \
  --scope identity \
  --id <id> \
  --daily <n> \
  --monthly <n> \
  --enabled
```

Provider‑level:

```bash
./arvis policy set-budget \
  --scope provider \
  --id openai \
  --daily <n>
```

Global:

```bash
./arvis policy set-budget \
  --scope global \
  --daily <n>
```

**Topic blocking**

```bash
./arvis policy block-topic [key]
./arvis policy unblock-topic [key]
./arvis policy list-topics       # shows all valid topic keys
./arvis policy show              # shows currently blocked topics
./arvis policy budgets           # shows every budget & current usage
```

### 5.6 Observability

**Requests**

```bash
./arvis requests list
./arvis requests list --identity <id>
./arvis requests list --limit N
# combine: --identity + --limit
```

**Anomalies**

```bash
./arvis anomalies list
./arvis anomalies list --category X --severity X --status X --limit N
./arvis anomalies resolve [id] --status reviewed
./arvis anomalies resolve [id] --status dismissed
```

### 5.7 Audit Export

```bash
./arvis audit export \
  [--identity <id>] \
  [--from YYYY-MM-DD] \
  [--to YYYY-MM-DD] \
  [--out path.pdf]
```

Defaults: `--to` = today, `--from` = 30 days before, `--out` = `arvis-audit-<timestamp>.pdf`.

### 5.8 Testing

```bash
# Full test suite (vet, race, coverage, all packages)
./arvis test
```

**Make targets** (if available):

```bash
make test            # all tests
make test-unit       # no DB required
make test-integration # requires arvis_test
```

**Direct `go test`** (example for store tests):

```bash
TEST_DATABASE_URL="postgres://arvis:arvis@localhost:5432/arvis_test?sslmode=disable" \
go test ./internal/store/... -v
```

---

## 6. Database Setup

### 6.1 Development database

- Name: `arvis` – used by the live server and CLI.
- Created automatically by the Postgres container when `POSTGRES_DB=arvis` is set.
- **Tests never touch it.**

### 6.2 Test database

- Name: `arvis_test`
- Must be created manually **once**:

```bash
# If using Docker:
docker compose exec postgres psql -U arvis -d arvis -c "CREATE DATABASE arvis_test OWNER arvis;"

# If using native Postgres:
sudo -u postgres psql -c "CREATE DATABASE arvis_test OWNER arvis;"
```

- Apply migrations to the test DB:

```bash
DATABASE_URL="postgres://arvis:arvis@localhost:5432/arvis_test?sslmode=disable" \
./arvis migrate up
```

- The test suite will wipe and rewrite this database, so it’s safe to reuse.

---

## 7. Redis

- Development Redis is on `localhost:6379` (mapped from container).
- **No separate test Redis instance** exists yet. Policy tests (budgets, topic blocking) will share the same Redis.  
  **TODO:** Use a dedicated DB index or flush test keys between runs.

---

## 8. Recommended Day‑to‑Day Workflows

### 8.1 Starting the full stack (development)

```bash
docker compose up --build   # rebuilds if code changed, streams logs
```

or detached:

```bash
docker compose up -d
```

The container’s entrypoint automatically runs `migrate up` before starting the server.

### 8.2 Running admin commands quickly (host binary)

```bash
export DATABASE_URL="postgres://arvis:arvis@localhost:5432/arvis?sslmode=disable"
./arvis identity create "New User"
./arvis policies list
./arvis audit export --out report.pdf
```

### 8.3 Running the server locally for fast iteration (without Docker for the app)

```bash
# Start only dependencies in Docker
docker compose up -d postgres redis

# Build and run the server on host
go build -o arvis ./cmd/arvis/
./arvis server
```

This gives instant rebuilds without Docker overhead.

### 8.4 Resetting everything to a clean state

```bash
docker compose down -v   # deletes volumes, including Postgres data
docker compose up --build
# then re‑create test DB if needed
```

---

## 9. Troubleshooting Quick Reference

| Symptom                                    | Likely cause                           | Fix                                                |
| ------------------------------------------ | -------------------------------------- | -------------------------------------------------- |
| `database ping failed: context deadline`   | `DATABASE_URL` points to wrong host   | Set to `localhost` for host binary                 |
| Port `8080` already in use                 | Another ARVIS instance running         | Stop other process or remap ports                  |
| Container exits immediately                | Entrypoint `sh -c` conflict            | Use the corrected `docker-entrypoint.sh` (provided)|
| Migrations not applied on startup          | Entrypoint not set correctly           | Ensure `ENTRYPOINT ["/app/docker-entrypoint.sh"]`  |
| `arvis_test` not found                     | Test DB not created                    | Run `CREATE DATABASE arvis_test OWNER arvis;`      |

---

## 10. Final Notes

- Always run commands from the repository root to avoid path‑related errors.
- The binary and the containerised server share the same database – they are **two interfaces to one system**, not two separate systems.
- Use the host binary for CLI speed; use `docker compose exec` to verify container‑specific behaviour.
- Keep `providers.yaml` up‑to‑date (it is mounted as a volume in Compose, so changes on the host reflect inside the container without rebuild).

---

*This document supersedes the earlier separate CLI and Docker references.*  
*Last updated: 2026-08-29 after successful integration of entrypoint script and port binding.*