package proxy

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kevinkiplangat432/arvis/internal/config"
)

// Proxy is the reverse proxy that sits in front of every configured AI
// provider. Exported as a type, not a bare handler func, so Phase 6's
// latency/token capture and Phase 7's detector hooks have somewhere to
// attach state (an HTTP client, in-memory rate counters) without
// changing this type's shape from the outside.
type Proxy struct {
	cfg    *config.Config
	db     *pgxpool.Pool
	logger *slog.Logger
	client *http.Client
}

func New(cfg *config.Config, db *pgxpool.Pool, logger *slog.Logger) *Proxy {
	return &Proxy{
		cfg:    cfg,
		db:     db,
		logger: logger,
		client: &http.Client{},
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	identity, err := authenticate(ctx, p.db, r)
	if err != nil {
		p.writeError(w, http.StatusUnauthorized, err)
		return
	}

	provider, body, err := resolveProvider(p.cfg, r)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, ErrUnknownModel) {
			status = http.StatusBadRequest
		}
		p.writeError(w, status, err)
		return
	}

	outboundURL := provider.BaseURL + r.URL.Path
	outbound, err := http.NewRequestWithContext(ctx, r.Method, outboundURL, body)
	if err != nil {
		p.writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to build outbound request: %w", err))
		return
	}
	outbound.Header = r.Header.Clone()
	// The caller's ARVIS key was only for authenticating to ARVIS itself
	// — it must never be forwarded upstream. The real provider key is
	// injected here instead, which is what lets a client hold exactly
	// one credential regardless of how many providers sit behind ARVIS.
	outbound.Header.Set("Authorization", "Bearer "+provider.APIKey)

	resp, err := p.client.Do(outbound)
	if err != nil {
		p.writeError(w, http.StatusBadGateway, fmt.Errorf("upstream provider %q unreachable: %w", provider.Name, err))
		return
	}
	defer resp.Body.Close()

	for k, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	// identity and provider are already resolved by this point in the
	// request — Phase 6 attaches latency/token capture and
	// store.InsertRequest right here, once the response has actually
	// reached the caller, so logging never adds to response latency.
	_ = identity
}

func (p *Proxy) writeError(w http.ResponseWriter, status int, err error) {
	p.logger.Warn("proxy request rejected", "status", status, "error", err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error": %q}`, err.Error())
}