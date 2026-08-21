package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/kevinkiplangat432/arvis/internal/config"
	"github.com/kevinkiplangat432/arvis/internal/detector"
	"github.com/kevinkiplangat432/arvis/internal/policy"
	"github.com/kevinkiplangat432/arvis/internal/store"
)

type Proxy struct {
	cfg      *config.Config
	db       *pgxpool.Pool
	rdb      *redis.Client
	logger   *slog.Logger
	client   *http.Client
	detector *detector.Detector
}

func New(cfg *config.Config, db *pgxpool.Pool, rdb *redis.Client, logger *slog.Logger) *Proxy {
	return &Proxy{
		cfg:      cfg,
		db:       db,
		rdb:      rdb,
		logger:   logger,
		client:   &http.Client{},
		detector: detector.New(db, logger),
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

	provider, model, bodyReader, err := resolveProvider(p.cfg, r)
	if err != nil {
		p.writeError(w, http.StatusBadRequest, err)
		return
	}

	requestBodyBytes, err := io.ReadAll(bodyReader)
	if err != nil {
		p.writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to buffer request body: %w", err))
		return
	}

	// Governance runs before anything touches the network. This is the
	// actual line between Auditing (always records, never acts) and
	// Governance (can say no) — a blocked request never reaches a
	// provider, but it IS still logged, just as a 403 instead of a
	// real response.
	violation, err := policy.Check(ctx, p.rdb, identity.ID, provider.Name, requestBodyBytes)
	if err != nil {
		// Fail open on a Redis outage: taking down the whole proxy
		// because its policy store briefly hiccuped is a worse outcome
		// than one request going unchecked. Worth revisiting once
		// Redis has real production HA behind it.
		p.logger.Error("policy check failed, proceeding without enforcement", "error", err.Error())
	}
	if violation != nil {
		p.logBlockedRequest(identity, provider, model, *violation, int(time.Since(start).Milliseconds()))
		p.writeError(w, http.StatusForbidden, fmt.Errorf("request blocked by policy: %s", violation.Detail))
		return
	}

	syncFlags := p.detector.CheckSync(ctx, identity.ID, model)


	if isStreaming(requestBodyBytes) {
		blocked := policy_ListBlockedTopicsSafe(ctx, p.rdb, p.logger)
		p.serveStreaming(w, r, identity, provider, model, requestBodyBytes, blocked, start)
		_ = syncFlags // sync detector flags aren't wired into streamed responses yet — the pre-flight policy check already ran on the initiating prompt either way
		return
	}

	outboundURL := provider.BaseURL + r.URL.Path
	outbound, err := http.NewRequestWithContext(ctx, r.Method, outboundURL, io.NopCloser(bytes.NewReader(requestBodyBytes)))
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
	go p.logAndDetect(identity, provider, model, resp.StatusCode, latencyMs, requestBodyBytes, respBody, syncFlags)
}

func policy_ListBlockedTopicsSafe(ctx context.Context, rdb *redis.Client, logger *slog.Logger) []string {
	blocked, err := policy.ListBlockedTopics(ctx, rdb)
	if err != nil {
		logger.Error("failed to fetch blocked topics for streaming check", "error", err.Error())
		return nil
	}
	return blocked
}

func (p *Proxy) logAndDetect(identity *store.Identity, provider config.Provider, model string, statusCode, latencyMs int, requestBody, responseBody []byte, syncFlags []detector.Flag) {
	ctx := context.Background()
	promptTokens, completionTokens := extractTokenUsage(responseBody)
	totalTokens := promptTokens + completionTokens

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

	if err := store.InsertRequest(ctx, p.db, req); err != nil {
		p.logger.Error("failed to log request", "error", err.Error())
		return
	}

	if err := policy.RecordUsageAll(ctx, p.rdb, identity.ID, provider.Name, totalTokens); err != nil {
		p.logger.Error("failed to record policy usage", "error", err.Error())
	}

	p.detector.RunAsync(req.ID, syncFlags, detector.AsyncInput{
		IdentityID:   identity.ID,
		Provider:     provider.Name,
		Model:        model,
		RequestBody:  requestBody,
		ResponseBody: responseBody,
		LatencyMs:    latencyMs,
		StatusCode:   statusCode,
	})
}

// logBlockedRequest records a policy-blocked call exactly the way a
// forwarded one is recorded — a 403 in the requests table, plus a
// governance-category anomaly attached to it. A blocked request
// belongs in the audit trail every bit as much as an allowed one.
func (p *Proxy) logBlockedRequest(identity *store.Identity, provider config.Provider, model string, violation policy.Violation, latencyMs int) {
	ctx := context.Background()

	req := store.Request{
		ID:         uuid.NewString(),
		IdentityID: &identity.ID,
		Provider:   provider.Name,
		Model:      model,
		LatencyMs:  latencyMs,
		StatusCode: http.StatusForbidden,
		CreatedAt:  time.Now(),
	}
	if err := store.InsertRequest(ctx, p.db, req); err != nil {
		p.logger.Error("failed to log blocked request", "error", err.Error())
		return
	}

	anomaly := store.Anomaly{
		ID:        uuid.NewString(),
		RequestID: req.ID,
		Rule:      violation.Rule,
		Detail:    violation.Detail,
		Category:  "governance",
		Severity:  violation.Severity,
		Status:    "open",
		CreatedAt: time.Now(),
	}
	if err := store.InsertAnomaly(ctx, p.db, anomaly); err != nil {
		p.logger.Error("failed to log policy violation", "error", err.Error())
	}
}

func (p *Proxy) writeError(w http.ResponseWriter, status int, err error) {
	p.logger.Warn("proxy request rejected", "status", status, "error", err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error": %q}`, err.Error())
}