package proxy

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/kevinkiplangat432/arvis/internal/config"
	"github.com/kevinkiplangat432/arvis/internal/policy"
	"github.com/kevinkiplangat432/arvis/internal/store"
)

// killSwitchWindow is how much recently-streamed text is kept for
// matching. Small on purpose — this check runs on every chunk, and a
// large rolling buffer would work against the low-latency goal
// everything else in the proxy has held to.
const killSwitchWindow = 2000

// isStreaming reports whether the request explicitly asked for a
// streamed response, the standard OpenAI/Anthropic "stream": true
// flag. Only these requests take the streaming path — everything
// else keeps using the proven buffered path, unchanged.
func isStreaming(body []byte) bool {
	var parsed struct {
		Stream bool `json:"stream"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return false
	}
	return parsed.Stream
}

// serveStreaming relays an SSE response chunk by chunk, forwarding
// each one to the client immediately, then checking a rolling buffer
// of recently-streamed text against the same blocked-topic catalog
// the pre-flight policy check uses. A match stops the stream
// immediately — everything already sent stays sent, nothing further
// does.
func (p *Proxy) serveStreaming(w http.ResponseWriter, r *http.Request, identity *store.Identity, provider config.Provider, model string, requestBody []byte, blockedTopics []string, start time.Time) {
	ctx := r.Context()

	outboundURL := provider.BaseURL + r.URL.Path
	outbound, err := http.NewRequestWithContext(ctx, r.Method, outboundURL, io.NopCloser(bytes.NewReader(requestBody)))
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

	flusher, ok := w.(http.Flusher)
	if !ok {
		p.writeError(w, http.StatusInternalServerError, fmt.Errorf("streaming not supported by this response writer"))
		return
	}

	for k, values := range resp.Header {
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	reader := bufio.NewReader(resp.Body)
	var rolling strings.Builder
	var killed *policy.Violation

	for {
		line, readErr := reader.ReadString('\n')
		if len(line) > 0 {
			// Forward first — the client never waits on the scan.
			w.Write([]byte(line))
			flusher.Flush()

			rolling.WriteString(line)
			if rolling.Len() > killSwitchWindow {
				trimmed := rolling.String()
				rolling.Reset()
				rolling.WriteString(trimmed[len(trimmed)-killSwitchWindow:])
			}

			if v := policy.MatchText(rolling.String(), blockedTopics); v != nil {
				killed = v
				fmt.Fprintf(w, "\ndata: {\"arvis_kill_switch\": %q}\n\n", v.Detail)
				flusher.Flush()
				break
			}
		}
		if readErr != nil {
			break
		}
	}

	latencyMs := int(time.Since(start).Milliseconds())
	if killed != nil {
		p.logKilledStream(identity, provider, model, *killed, latencyMs)
		return
	}

	go p.logStreamedRequest(identity, provider, model, resp.StatusCode, latencyMs)
}

func (p *Proxy) logKilledStream(identity *store.Identity, provider config.Provider, model string, violation policy.Violation, latencyMs int) {
	ctx := context.Background()

	req := store.Request{
		ID:         uuid.NewString(),
		IdentityID: &identity.ID,
		Provider:   provider.Name,
		Model:      model,
		LatencyMs:  latencyMs,
		StatusCode: 499, // non-standard, but the common convention for "connection closed before completion" — here, ARVIS closed it
		CreatedAt:  time.Now(),
	}
	if err := store.InsertRequest(ctx, p.db, req); err != nil {
		p.logger.Error("failed to log killed stream", "error", err.Error())
		return
	}

	anomaly := store.Anomaly{
		ID:        uuid.NewString(),
		RequestID: req.ID,
		Rule:      "kill_switch_" + violation.Rule,
		Detail:    "stream terminated mid-flight: " + violation.Detail,
		Category:  "governance",
		Severity:  violation.Severity,
		Status:    "open",
		CreatedAt: time.Now(),
	}
	if err := store.InsertAnomaly(ctx, p.db, anomaly); err != nil {
		p.logger.Error("failed to log kill switch anomaly", "error", err.Error())
	}
}

// logStreamedRequest handles the common case: a stream that completed
// normally. Token usage is deliberately left at zero here — most
// providers only send real usage figures in the final SSE event, and
// parsing that reliably is more than this pass was scoped for. Worth
// a real follow-up once streaming volume justifies it.
func (p *Proxy) logStreamedRequest(identity *store.Identity, provider config.Provider, model string, statusCode, latencyMs int) {
	ctx := context.Background()
	req := store.Request{
		ID:         uuid.NewString(),
		IdentityID: &identity.ID,
		Provider:   provider.Name,
		Model:      model,
		LatencyMs:  latencyMs,
		StatusCode: statusCode,
		CreatedAt:  time.Now(),
	}
	if err := store.InsertRequest(ctx, p.db, req); err != nil {
		p.logger.Error("failed to log streamed request", "error", err.Error())
	}
}