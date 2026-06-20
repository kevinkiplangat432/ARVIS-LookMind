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