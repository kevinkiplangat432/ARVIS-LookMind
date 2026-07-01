package proxy

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevinkiplangat432/arvis/internal/config"
	"github.com/kevinkiplangat432/arvis/internal/detector"
	"github.com/kevinkiplangat432/arvis/internal/store"
)

type Proxy struct {
	rp     *httputil.ReverseProxy
	cfg    *config.Config
	db     *pgxpool.Pool
	logger *slog.Logger
}

func New(rp *httputil.ReverseProxy, cfg *config.Config, db *pgxpool.Pool, logger *slog.Logger) *Proxy {
	return &Proxy{rp: rp, cfg: cfg, db: db, logger: logger}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	rec := &statusRecorder{ResponseWriter: w, code: http.StatusOK}
	p.rp.ServeHTTP(rec, r)

	latency := time.Since(start).Milliseconds()

	p.logger.Info("proxy forwarded",
		"method", r.Method,
		"path", r.URL.Path,
		"status", rec.code,
		"latency_ms", latency,
	)

	req := store.Request{
		ID:         uuid.NewString(),
		Model:      r.Header.Get("X-Model"),
		LatencyMs:  int(latency),
		StatusCode: rec.code,
		CreatedAt:  time.Now(),
	}

	go func() {
		ctx := context.Background()
		if err := store.InsertRequest(ctx, p.db, req); err != nil {
			p.logger.Error("failed to insert request", "error", err, "request_id", req.ID)
		}

		for _, flag := range detector.Check(req) {
			if err := store.InsertAnomaly(ctx, p.db, store.Anomaly{
				ID:        uuid.NewString(),
				RequestID: req.ID,
				Rule:      flag.Rule,
				Detail:    flag.Detail,
				CreatedAt: time.Now(),
			}); err != nil {
				p.logger.Error("failed to insert anomaly", "error", err, "request_id", req.ID, "rule", flag.Rule)
			}
		}
	}()
}

type statusRecorder struct {
	http.ResponseWriter
	code int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}


// code review 1 
//Lines 16-21 (type Proxy struct): Dependency injection of the reverse proxy, config, and database pool. This is excellent for unit testing and modularity.
//Lines 27-31 (ServeHTTP Start): Captures the request arrival time and wraps the response writer to intercept the HTTP status code.
//Lines 42-48 (store.Request building): Prepares the telemetry data for persistence. Using uuid.NewString() for every request is standard but can become a CPU bottleneck at extremely high RPS.
//Lines 50-68 (go func() Background Processing): Offloads database writes to a background goroutine so the user doesn't wait for the database before getting their response. This is essential for low-latency proxies.



// SOC 2 Type II Analysis: Fail
//Context Mismanagement (Fail): You are using context.Background() inside the background goroutine [Lines ]. SOC 2 requires Availability and Processing Integrity. During a server shutdown, these background tasks will be "orphaned"—they don't respect the server's graceful shutdown deadline and could be killed mid-write, leading to missing audit logs.
//Header Leakage (Risk): You are pulling r.Header.Get("X-Model"). If an attacker or a misconfigured client sends sensitive data in headers that you later log or store without sanitization, you violate the Privacy and Confidentiality criteria.



// scaling to 1000 rps 
// The Issue: Goroutine and Connection Exhaustion
//Unbounded Goroutines (Critical): Spawning a new goroutine for every request (go func()) at 1,000 RPS will lead to thousands of concurrent goroutines. If the database slows down even slightly, these goroutines will pile up, consuming all available RAM and potentially crashing the process
// Database Contention: Direct inserts for every single request at high throughput create massive "lock contention" in PostgreSQL
// Standard Transport Limits: Without a custom http.Transport (as discussed in your first file review), the proxy will hit the default limit of 2 idle connections per host, causing requests to queue up and latency to spike
