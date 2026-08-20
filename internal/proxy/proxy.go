package proxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kevinkiplangat432/arvis/internal/config"
	"github.com/kevinkiplangat432/arvis/internal/store"
)

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
	start := time.Now()
	ctx := r.Context()

	identity, err := authenticate(ctx, p.db, r)
	if err != nil {
		p.writeError(w, http.StatusUnauthorized, err)
		return
	}

	provider, model, body, err := resolveProvider(p.cfg, r)
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
	outbound.Header.Set("Authorization", "Bearer "+provider.APIKey)

	resp, err := p.client.Do(outbound)
	if err != nil {
		p.writeError(w, http.StatusBadGateway, fmt.Errorf("upstream provider %q unreachable: %w", provider.Name, err))
		return
	}
	defer resp.Body.Close()

	// Buffered rather than streamed, deliberately: token usage lives at
	// the end of a typical JSON response body, so it can't be read
	// without the whole thing. This is a real trade-off — it adds
	// latency proportional to response size and won't work for
	// streaming (SSE) responses at all. Both are known limitations to
	// revisit once streaming support is actually needed.
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		p.writeError(w, http.StatusBadGateway, fmt.Errorf("failed to read upstream response: %w", err))
		return
	}

	for k, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)

	latencyMs := int(time.Since(start).Milliseconds())
	go p.logRequest(identity, provider, model, resp.StatusCode, latencyMs, respBody)
}

func (p *Proxy) logRequest(identity *store.Identity, provider config.Provider, model string, statusCode, latencyMs int, respBody []byte) {
	promptTokens, completionTokens := extractTokenUsage(respBody)

	req := store.Request{
		ID:           uuid.NewString(),
		IdentityID:   &identity.ID,
		Provider:     provider.Name,
		Model:        model,
		PromptTokens: promptTokens,
		CompTokens:   completionTokens,
		LatencyMs:    latencyMs,
		StatusCode:   statusCode,
		CreatedAt:    time.Now(),
	}

	if err := store.InsertRequest(context.Background(), p.db, req); err != nil {
		p.logger.Error("failed to log request", "error", err.Error())
	}
}

func (p *Proxy) writeError(w http.ResponseWriter, status int, err error) {
	p.logger.Warn("proxy request rejected", "status", status, "error", err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error": %q}`, err.Error())
}