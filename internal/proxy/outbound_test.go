// outbound_test.go guards the credential swap in newOutboundRequest: the
// caller's ARVIS key never reaches the provider, on either path.
//
// Not covered: the full request lifecycle, which needs a database. Covered
// separately by the auth tests and manual runs.
//
// Run: go test ./internal/proxy/...
// Needs: nothing, no database or network.

package proxy

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kevinkiplangat432/arvis/internal/config"
)

const testArvisKey = "arvis_0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

func inboundRequest() *http.Request {
	in := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", strings.NewReader(`{"model":"x"}`))
	in.Header.Set("Authorization", "Bearer "+testArvisKey)
	in.Header.Set("Content-Type", "application/json")
	return in
}

func TestNewOutboundRequestSwapsCredentials(t *testing.T) {
	in := inboundRequest()
	provider := config.Provider{Name: "openai", BaseURL: "https://api.example.com", APIKey: "sk-provider"}

	out, err := newOutboundRequest(context.Background(), in, provider, []byte(`{"model":"masked"}`))
	if err != nil {
		t.Fatalf("newOutboundRequest: %v", err)
	}

	if got := out.Header.Get("Authorization"); got != "Bearer sk-provider" {
		t.Errorf("Authorization = %q, want the provider key", got)
	}
	for name, values := range out.Header {
		for _, v := range values {
			if strings.Contains(v, testArvisKey) {
				t.Errorf("ARVIS key leaked in header %s", name)
			}
		}
	}
	if out.Header.Get("Content-Type") != "application/json" {
		t.Error("other client headers should pass through")
	}
	if out.Method != http.MethodPost {
		t.Errorf("method = %s, want POST", out.Method)
	}
	if out.URL.String() != "https://api.example.com/v1/chat/completions" {
		t.Errorf("url = %s", out.URL)
	}
	body, _ := io.ReadAll(out.Body)
	if string(body) != `{"model":"masked"}` {
		t.Errorf("body = %s, want the supplied body", body)
	}
	if in.Header.Get("Authorization") != "Bearer "+testArvisKey {
		t.Error("inbound request was mutated, headers must be cloned")
	}
}

// A misconfigured provider with no key must still never forward the ARVIS key.
func TestNewOutboundRequestEmptyProviderKeyStillStripsArvisKey(t *testing.T) {
	provider := config.Provider{Name: "openai", BaseURL: "https://api.example.com", APIKey: ""}

	out, err := newOutboundRequest(context.Background(), inboundRequest(), provider, nil)
	if err != nil {
		t.Fatalf("newOutboundRequest: %v", err)
	}
	if strings.Contains(out.Header.Get("Authorization"), testArvisKey) {
		t.Error("ARVIS key forwarded when provider key was empty")
	}
}