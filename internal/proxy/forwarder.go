package proxy

import (
	"context"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/kevinkiplangat432/arvis/internal/config"
	"github.com/kevinkiplangat432/arvis/internal/detector"
	"github.com/kevinkiplangat432/arvis/internal/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type forwarder struct {
	rp  *httputil.ReverseProxy
	cfg *config.Config
	db  *pgxpool.Pool
}

func newForwarder(rp *httputil.ReverseProxy, cfg *config.Config, db *pgxpool.Pool) *forwarder {
	return &forwarder{rp: rp, cfg: cfg, db: db}
}

func (f *forwarder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	rec := &statusRecorder{ResponseWriter: w, code: 200}
	f.rp.ServeHTTP(rec, r)

	req := store.Request{
		ID:         uuid.NewString(),
		Model:      r.Header.Get("X-Model"),
		LatencyMs:  int(time.Since(start).Milliseconds()),
		StatusCode: rec.code,
		CreatedAt:  time.Now(),
	}

	go func() {
		ctx := context.Background()
		_ = store.InsertRequest(ctx, f.db, req)

		for _, flag := range detector.Check(req) {
			_ = store.InsertAnomaly(ctx, f.db, store.Anomaly{
				ID:        uuid.NewString(),
				RequestID: req.ID,
				Rule:      flag.Rule,
				Detail:    flag.Detail,
				CreatedAt: time.Now(),
			})
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
