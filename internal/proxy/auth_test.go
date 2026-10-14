package proxy

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"	
)

func TestAuthenticateRejectsBadHeaders(t *testing.T) {
	tests := []struct {
		name string
		header string
	}{
		{"no header", ""},
		{"wrong scheme", "Basic dXNlcjpwYXNz"},
		{"bearer without key","Bearer"},
		{"bearer without space", "Bearer"},
		{"bearer without space ", "Bearer"},
		{"lowercase schema", "bearer arvis_abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)

			}

			identity, err := authenticate(context.Background(), nil, req)
			if  !errors.Is(err, ErrMissingKey) {
				t.Errorf("err = %v, want ErrMissingKey", err)
			}
			if identity != nil {
				t.Error ("identity should be nil on rejection")
			}
		} )
	}
}

