// auth_fuzz_test.go fuzzes extractBearer, the parser that reads attacker
// controlled Authorization headers on the front door.
//
// Properties: never panics, an error always means no token, a nil error
// always means a non empty trimmed token that parses back to itself.
//
// Run seeds only: go test ./internal/proxy/...
// Run fuzzing:    go test ./internal/proxy -run='^$' -fuzz='^FuzzExtractBearer$' -fuzztime=60s
// Needs: nothing, no database or network.

package proxy

import (
	"errors"
	"strings"
	"testing"
)

func FuzzExtractBearer(f *testing.F) {
	for _, seed := range []string{
		"", "Bearer arvis_abc", "bearer x", "BEARER   x  ", "Bearer", "Bearer ",
		"Basic abc", "Bearer\tx", "Bearer a b", "Bearer \x00",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, header string) {
		token, err := extractBearer(header)
		if err != nil {
			if !errors.Is(err, ErrMissingKey) {
				t.Fatalf("unexpected error type: %v", err)
			}
			if token != "" {
				t.Fatalf("error returned together with token %q", token)
			}
			return
		}
		if token == "" {
			t.Fatal("nil error with empty token")
		}
		if token != strings.TrimSpace(token) {
			t.Fatalf("token %q has surrounding whitespace", token)
		}
		again, err := extractBearer("Bearer " + token)
		if err != nil || again != token {
			t.Fatalf("round trip failed: %q -> %q (%v)", token, again, err)
		}
	})
}