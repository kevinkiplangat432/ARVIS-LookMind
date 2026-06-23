package detector

import (
	"testing"

	"github.com/kevinkiplangat432/arvis/internal/store"
)

func TestCheck_NoFlags(t *testing.T) {
	r := store.Request{
		LatencyMs:    9999,
		PromptTokens: 1500,
		CompTokens:   1499,
		StatusCode:   200,
	}
	flags := Check(r)
	if len(flags) != 0 {
		t.Errorf("expected no flags, got %d: %+v", len(flags), flags)
	}
}

func TestCheck_HighLatency_ExactBoundary(t *testing.T) {
	tests := []struct {
		name      string
		latencyMs int
		wantFlag  bool
	}{
		{"below threshold", 9999, false},
		{"at threshold", 10000, false},
		{"above threshold", 10001, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := store.Request{LatencyMs: tt.latencyMs, StatusCode: 200}
			flags := Check(r)
			got := containsRule(flags, "high_latency")
			if got != tt.wantFlag {
				t.Errorf("latency=%d: expected high_latency=%v, got %v", tt.latencyMs, tt.wantFlag, got)
			}
		})
	}
}

func TestCheck_HighTokenUsage_ExactBoundary(t *testing.T) {
	tests := []struct {
		name         string
		promptTokens int
		compTokens   int
		wantFlag     bool
	}{
		{"below threshold", 1500, 1499, false},
		{"at threshold", 1500, 1500, false},
		{"above threshold by one", 1500, 1501, true},
		{"prompt alone exceeds", 3001, 0, true},
		{"comp alone exceeds", 0, 3001, true},
		{"both zero", 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := store.Request{
				PromptTokens: tt.promptTokens,
				CompTokens:   tt.compTokens,
				StatusCode:   200,
			}
			flags := Check(r)
			got := containsRule(flags, "high_token_usage")
			if got != tt.wantFlag {
				t.Errorf("prompt=%d comp=%d: expected high_token_usage=%v, got %v",
					tt.promptTokens, tt.compTokens, tt.wantFlag, got)
			}
		})
	}
}

func TestCheck_UpstreamError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantFlag   bool
	}{
		{"499 no flag", 499, false},
		{"500 flags", 500, true},
		{"503 flags", 503, true},
		{"599 flags", 599, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := store.Request{StatusCode: tt.statusCode}
			flags := Check(r)
			got := containsRule(flags, "upstream_error")
			if got != tt.wantFlag {
				t.Errorf("status=%d: expected upstream_error=%v, got %v", tt.statusCode, tt.wantFlag, got)
			}
		})
	}
}

func TestCheck_ClientError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantFlag   bool
	}{
		{"399 no flag", 399, false},
		{"400 flags", 400, true},
		{"403 flags", 403, true},
		{"404 flags", 404, true},
		{"499 flags", 499, true},
		{"500 no client_error flag", 500, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := store.Request{StatusCode: tt.statusCode}
			flags := Check(r)
			got := containsRule(flags, "client_error")
			if got != tt.wantFlag {
				t.Errorf("status=%d: expected client_error=%v, got %v", tt.statusCode, tt.wantFlag, got)
			}
		})
	}
}

func TestCheck_CorrectDetail(t *testing.T) {
	tests := []struct {
		name           string
		request        store.Request
		rule           string
		expectedDetail string
	}{
		{
			name:           "upstream_error detail includes status code",
			request:        store.Request{StatusCode: 502},
			rule:           "upstream_error",
			expectedDetail: "upstream returned 502",
		},
		{
			name:           "client_error detail includes status code",
			request:        store.Request{StatusCode: 403},
			rule:           "client_error",
			expectedDetail: "upstream returned 403",
		},
		{
			name:           "high_latency detail",
			request:        store.Request{LatencyMs: 10001, StatusCode: 200},
			rule:           "high_latency",
			expectedDetail: "latency exceeded 10s",
		},
		{
			name:           "high_token_usage detail",
			request:        store.Request{PromptTokens: 3001, StatusCode: 200},
			rule:           "high_token_usage",
			expectedDetail: "total tokens exceeded 3000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := Check(tt.request)
			for _, f := range flags {
				if f.Rule == tt.rule {
					if f.Detail != tt.expectedDetail {
						t.Errorf("rule=%s: expected detail %q, got %q", tt.rule, tt.expectedDetail, f.Detail)
					}
					return
				}
			}
			t.Errorf("rule %q not found in flags: %+v", tt.rule, flags)
		})
	}
}

func TestCheck_MultipleRulesFire(t *testing.T) {
	r := store.Request{
		LatencyMs:    15000,
		PromptTokens: 2000,
		CompTokens:   2000,
		StatusCode:   503,
	}
	flags := Check(r)
	expected := []string{"high_latency", "high_token_usage", "upstream_error"}
	for _, rule := range expected {
		if !containsRule(flags, rule) {
			t.Errorf("expected rule %q to fire, flags: %+v", rule, flags)
		}
	}
	if len(flags) != len(expected) {
		t.Errorf("expected exactly %d flags, got %d: %+v", len(expected), len(flags), flags)
	}
}

func TestCheck_ZeroValueRequest(t *testing.T) {
	r := store.Request{}
	flags := Check(r)
	if len(flags) != 0 {
		t.Errorf("zero value request should produce no flags, got %d: %+v", len(flags), flags)
	}
}

// containsRule is a test helper — not exported, not part of production code.
func containsRule(flags []Flag, rule string) bool {
	for _, f := range flags {
		if f.Rule == rule {
			return true
		}
	}
	return false
}