// auth_test.go guards the rules in auth.go: header parsing, case
// insensitive scheme, and the mapping from errors to HTTP statuses.
//
// Not covered: the database paths of authenticate (unknown key, outage,
// hash mismatch). They need a lookup seam, deferred on purpose, see
// docs/auth_and_identity.md.
//
// Run: go test ./internal/proxy/...
// Needs: nothing, no database or network.

package proxy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtractBearer(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		want    string
		wantErr bool
	}{
		{"standard", "Bearer arvis_abc", "arvis_abc", false},
		{"lowercase scheme", "bearer arvis_abc", "arvis_abc", false},
		{"uppercase scheme", "BEARER arvis_abc", "arvis_abc", false},
		{"extra spaces before token", "Bearer    arvis_abc", "arvis_abc", false},
		{"trailing whitespace", "Bearer arvis_abc  ", "arvis_abc", false},
		{"no header", "", "", true},
		{"wrong scheme", "Basic dXNlcjpwYXNz", "", true},
		{"scheme only", "Bearer", "", true},
		{"scheme and one space", "Bearer ", "", true},
		{"scheme and many spaces", "Bearer     ", "", true},
		{"no space between scheme and token", "Bearerarvis_abc", "", true},
		{"tab instead of space", "Bearer\tarvis_abc", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractBearer(tt.header)
			if tt.wantErr {
				if !errors.Is(err, ErrMissingKey) {
					t.Errorf("err = %v, want ErrMissingKey", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// Rejections return before the database is touched, so a nil pool is safe.
func TestAuthenticateRejectsBeforeDatabase(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)

	identity, err := authenticate(context.Background(), nil, req)
	if !errors.Is(err, ErrMissingKey) {
		t.Errorf("err = %v, want ErrMissingKey", err)
	}
	if identity != nil {
		t.Error("identity should be nil on rejection")
	}
}

func TestAuthErrorStatus(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name string
		ctx  context.Context
		err  error
		want int
	}{
		{"missing key", context.Background(), ErrMissingKey, http.StatusUnauthorized},
		{"unknown key", context.Background(), ErrUnknownKey, http.StatusUnauthorized},
		{"wrapped unknown key", context.Background(), fmt.Errorf("wrap: %w", ErrUnknownKey), http.StatusUnauthorized},
		{"database down", context.Background(), errors.New("connection refused"), http.StatusServiceUnavailable},
		{"deadline exceeded is our problem", context.Background(), context.DeadlineExceeded, http.StatusServiceUnavailable},
		{"canceled error", context.Background(), fmt.Errorf("lookup: %w", context.Canceled), 0},
		{"canceled context, opaque error", canceled, errors.New("driver said something odd"), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := authErrorStatus(tt.ctx, tt.err); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}