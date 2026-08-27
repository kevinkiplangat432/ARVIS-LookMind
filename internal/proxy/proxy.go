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
	"github.com/kevinkiplangat432/arvis/internal/tokenize"
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
	requestID := uuid.NewString() // generated up front now — tokenize needs it before forwarding, not just for logging afterward

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

	violation, err := policy.Check(ctx, p.rdb, identity.ID, provider.Name, requestBodyBytes)
	if err != nil {
		p.logger.Error("policy check failed, proceeding without enforcement", "error", err.Error())
	}
	if violation != nil {
		p.logBlockedRequest(requestID, identity, provider, model, *violation, int(time.Since(start).Milliseconds()))
		p.writeError(w, http.StatusForbidden, fmt.Errorf("request blocked by policy: %s", violation.Detail))
		return
	}

	syncFlags := p.detector.CheckSync(ctx, identity.ID, model)

	// Tokenization fails CLOSED, deliberately the opposite of policy's
	// fail-open Redis behavior above. A Redis outage during a policy
	// check means one request goes unchecked — bad, but recoverable.
	// A Redis outage during tokenization, if it failed open, would
	// mean real customer PII goes to an external provider unmasked —
	// exactly the thing this feature exists to prevent. Silence is
	// not an acceptable failure mode here.
	tokenizedBody, err := tokenize.Tokenize(ctx, p.rdb, requestID, requestBodyBytes)
	if err != nil {
		p.writeError(w, http.StatusServiceUnavailable, fmt.Errorf("tokenization unavailable, request rejected rather than sent unmasked: %w", err))
		return
	}

	outboundURL := provider.BaseURL + r.URL.Path
	outbound, err := http.NewRequestWithContext(ctx, r.Method, outboundURL, io.NopCloser(bytes.NewReader(tokenizedBody)))
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

	// Detokenize fails OPEN, on purpose, opposite of Tokenize above —
	// the asymmetry is intentional. A failure here means the caller
	// gets a response still containing tokens like [KENYAN_NATIONAL_ID_1]
	// instead of the real value: degraded, but not a leak in either
	// direction, so there's no safety reason to reject the request.
	reconstructed, err := tokenize.Detokenize(ctx, p.rdb, requestID, respBody)
	if err != nil {
		p.logger.Error("detokenization failed, returning response with tokens intact", "error", err.Error())
		reconstructed = respBody
	}

	for k, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(reconstructed)

	latencyMs := int(time.Since(start).Milliseconds())
	go p.logAndDetect(requestID, identity, provider, model, resp.StatusCode, latencyMs, requestBodyBytes, reconstructed, syncFlags)
}

func (p *Proxy) logAndDetect(requestID string, identity *store.Identity, provider config.Provider, model string, statusCode, latencyMs int, requestBody, responseBody []byte, syncFlags []detector.Flag) {
	ctx := context.Background()
	promptTokens, completionTokens := extractTokenUsage(responseBody)
	totalTokens := promptTokens + completionTokens

	req := store.Request{
		ID:           requestID,
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

	// requestBody/responseBody here are the real, reconstructed values
	// (pre-tokenize / post-detokenize) — the detector's job is to flag
	// what an employee actually typed and actually received, not the
	// masked version that briefly existed only for the provider call.
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

func (p *Proxy) logBlockedRequest(requestID string, identity *store.Identity, provider config.Provider, model string, violation policy.Violation, latencyMs int) {
	ctx := context.Background()

	req := store.Request{
		ID:         requestID,
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