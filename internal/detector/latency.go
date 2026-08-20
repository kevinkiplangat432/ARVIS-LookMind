package detector

import (
	"context"
	"fmt"
	"sync"
)

// LatencyRule flags requests whose latency is a significant outlier
// against that provider's own recent history, plus error status codes
// outright. Runs async — a request's own latency is only knowable
// after it's finished.
type LatencyRule struct {
	mu      sync.Mutex
	history map[string][]int // keyed by provider name
	window  int
}

func NewLatencyRule(window int) *LatencyRule {
	return &LatencyRule{
		history: make(map[string][]int),
		window:  window,
	}
}

func (l *LatencyRule) CheckAsync(ctx context.Context, in AsyncInput) []Flag {
	var flags []Flag

	switch {
	case in.StatusCode >= 500:
		flags = append(flags, Flag{
			Rule: "upstream_error", Category: "latency", Severity: "high",
			Detail: fmt.Sprintf("provider %q returned status %d", in.Provider, in.StatusCode),
		})
	case in.StatusCode >= 400:
		flags = append(flags, Flag{
			Rule: "client_error", Category: "latency", Severity: "low",
			Detail: fmt.Sprintf("provider %q returned status %d", in.Provider, in.StatusCode),
		})
	}

	l.mu.Lock()
	hist := l.history[in.Provider]
	avg := average(hist)
	hist = append(hist, in.LatencyMs)
	if len(hist) > l.window {
		hist = hist[len(hist)-l.window:]
	}
	l.history[in.Provider] = hist
	l.mu.Unlock()

	// Only judge once there's a real baseline — flagging against zero
	// prior samples is just noise.
	if avg > 0 && float64(in.LatencyMs) > avg*3 {
		flags = append(flags, Flag{
			Rule: "latency_outlier", Category: "latency", Severity: "medium",
			Detail: fmt.Sprintf("request took %dms, provider %q's recent average is %.0fms", in.LatencyMs, in.Provider, avg),
		})
	}

	return flags
}

func average(vals []int) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0
	for _, v := range vals {
		sum += v
	}
	return float64(sum) / float64(len(vals))
}