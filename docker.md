
<!--markdownlint-disable--># Docker Command Reference (for your ARVIS workflow)

## Building images

### `docker compose build`
Builds the image(s) defined in `docker-compose.yml` without starting any containers. Reads your `Dockerfile`, runs through each stage (in your case: the `golang:1.25-alpine` builder stage, then the `alpine:3.19` runtime stage), and produces a tagged image (e.g. `arvis-arvis:latest`) sitting ready on disk.
**Run when:** you changed Go code, the `Dockerfile`, or `go.mod`/`go.sum` and want to build without immediately running the stack (e.g. to check the build succeeds before committing).

### `docker compose up --build`
Combines build + start. Rebuilds any image whose build context changed, then creates and starts all containers.
**Run when:** you've changed your own source code (like your `router.go` fix earlier). Docker caches layers, so if only your Go files changed, it skips re-downloading `go mod`, and just reruns `COPY . .` and `go build` — usually fast.
**Common mistake:** forgetting `--build` after a code change (you hit this earlier) — Compose happily reuses the stale image and you keep chasing a bug you already fixed.

### `docker compose build --no-cache`
Forces a full rebuild ignoring all cached layers, including `go mod download`.
**Run when:** you suspect a corrupted cache, a dependency was updated but the lockfile didn't reflect it properly, or "it works on a clean machine but not here" debugging.

---

## Starting and stopping

### `docker compose up`
Starts all services (`postgres`, `redis`, `arvis` in your case) using existing images — no rebuild. Attaches your terminal to the combined logs of every container (that's the interleaved `postgres-1 | ...`, `redis-1 | ...`, `arvis-1 | ...` output you've been reading).
**Run when:** nothing changed since the last build, you're just resuming work for the day.

### `docker compose up -d`
Same as above, but detached — starts everything in the background and gives you your terminal prompt back immediately instead of streaming logs.
**Run when:** you don't need to watch logs live, e.g. you're about to run `curl` commands against the API and don't want log noise interrupting your terminal.

### `docker compose down`
Stops and **removes** all containers and the default network created for this project. Named volumes (like your `arvis_pgdata`) are preserved by default — your Postgres data survives.
**Run when:** you want a clean slate for containers/networking, e.g. before switching branches with a different compose setup, or to clear out a broken container state (like your port-conflict situation earlier).

### `docker compose down -v`
Same as `down`, but also deletes named volumes — this **wipes your Postgres data**.
**Run when:** you want to test migrations from a truly empty database, or your local dev data is garbage and you want to start fresh. Be deliberate with this one; it's destructive.

### `docker compose stop`
Stops running containers without removing them (keeps containers, networks, and volumes intact — just paused).
**Run when:** you want to pause everything temporarily (e.g. laptop going to sleep, switching to a different project) but plan to resume the exact same containers shortly.

### `docker compose start`
Restarts containers previously stopped with `stop` (as opposed to `up`, which also handles creation).
**Run when:** resuming after `docker compose stop`.

### `docker compose restart`
Stops then starts containers again, in place — doesn't rebuild, doesn't recreate.
**Run when:** a service is misbehaving (e.g. Redis connection got weird) and you just want to bounce it without a full rebuild.

### `docker compose restart arvis`
Restarts only the named service (`arvis` here), leaving `postgres` and `redis` untouched.
**Run when:** you want to reset just your app process without disturbing the database connection state or losing Redis data.

---

## Watching what's happening

### `docker compose ps`
Lists containers for this project, their status (`Up`, `Exited`, `Restarting`), and port mappings.
**Run when:** you want a quick "is everything actually running?" check without wading through logs.

### `docker compose logs`
Shows logs from all services since they started (not live-streaming by default — it's a snapshot).
**Run when:** you missed something that scrolled by, or you're checking after running in detached mode (`-d`).

### `docker compose logs -f`
Same as above but **follows** the logs live, like `tail -f`. Ctrl+C to stop watching (doesn't stop the containers).
**Run when:** you're actively debugging and want to watch requests come in as you hit endpoints.

### `docker compose logs -f arvis`
Follows logs from just the `arvis` service, filtering out Postgres/Redis noise.
**Run when:** you specifically care about your app's output, e.g. watching for a panic after a fix.

### `docker compose logs --tail 50 arvis`
Shows just the last 50 lines from `arvis`, without streaming.
**Run when:** you want a quick recent snapshot without scrolling through the full crash-loop history (relevant after what you just went through).

---

## Getting inside a running container

### `docker compose exec arvis sh`
Opens an interactive shell **inside the already-running** `arvis` container.
**Run when:** you want to poke around the live filesystem, check env vars actually loaded (`env | grep DB`), or manually run a binary command like `./arvis migrate status` against the live DB.

### `docker compose exec postgres psql -U arvis -d arvis`
Opens a `psql` prompt inside the Postgres container, logged in as your `arvis` user against the `arvis` database.
**Run when:** you want to manually inspect tables — e.g. `SELECT * FROM anomalies LIMIT 10;` — to sanity-check what your app is actually writing.

### `docker compose run --rm arvis sh`
Starts a **new, one-off** container from the `arvis` image (not the already-running one) and drops you into a shell; `--rm` auto-deletes the container when you exit.
**Run when:** you want an isolated environment to test something without touching your actual running app container — e.g. testing a migration command in isolation.

---

## Cleaning up disk space

### `docker image prune`
Deletes dangling (untagged, unused) images — typically old layers left behind by repeated builds.
**Run when:** `docker compose build` output starts looking cluttered with old intermediate images and you want to reclaim disk space.

### `docker system prune`
More aggressive — removes all stopped containers, unused networks, and dangling images across your whole machine (not just this project).
**Run when:** doing general Docker housekeeping. Use with a bit of care if you have other projects' containers stopped that you plan to resume.

### `docker volume ls`
Lists all named volumes on your machine, including `arvis_pgdata`.
**Run when:** checking whether your Postgres volume still exists, or spotting leftover volumes from old experiments.

---

## Quick reference — your typical day-to-day loop

```bash
# Changed Go code → rebuild and start, watch logs
docker compose up --build

# Resuming later, nothing changed
docker compose up

# Done for now, keep data
docker compose down

# Want a totally clean DB to retest migrations from scratch
docker compose down -v
docker compose up --build

# Something's acting weird, just bounce the app
docker compose restart arvis

# Peek inside the DB manually
docker compose exec postgres psql -U arvis -d arvis
```