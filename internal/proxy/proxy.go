package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kevinkiplangat432/arvis/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg *config.Config, db *pgxpool.Pool) http.Handler {
	target, _ := url.Parse(cfg.TargetURL)
	rp := httputil.NewSingleHostReverseProxy(target)

	return withMiddleware(newForwarder(rp, cfg, db))
}
