// keys_fuzz_test.go fuzzes the properties of keys.go for arbitrary input.
//
// Properties: HashKey is deterministic and always yields 64 lowercase hex
// characters, the prefix never changes the hash, and EqualHashes agrees
// with == for any pair of strings and is symmetric.
//
// Run seeds only: go test ./internal/auth/...
// Run fuzzing:    go test ./internal/auth -run='^$' -fuzz='^FuzzHashKey$' -fuzztime=30s
//                 go test ./internal/auth -run='^$' -fuzz='^FuzzEqualHashes$' -fuzztime=30s
// Needs: nothing, no database or network.

package auth

import (
	"encoding/hex"
	"strings"
	"testing"
)

func FuzzHashKey(f *testing.F) {
	for _, seed := range []string{"", "abc", "arvis_abc", "arvis_", "arvis_arvis_x"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		h := HashKey(s)
		if len(h) != 64 {
			t.Fatalf("hash length = %d, want 64", len(h))
		}
		if _, err := hex.DecodeString(h); err != nil || h != strings.ToLower(h) {
			t.Fatalf("hash %q is not lowercase hex", h)
		}
		if h != HashKey(s) {
			t.Fatal("HashKey is not deterministic")
		}
		// Only one prefix is stripped, so inputs that already start with
		// the prefix are skipped. Generated keys are hex and cannot.
		if !strings.HasPrefix(s, keyPrefix) && HashKey(keyPrefix+s) != h {
			t.Fatalf("prefix changed the hash for %q", s)
		}
	})
}

func FuzzEqualHashes(f *testing.F) {
	f.Add("a", "a")
	f.Add("a", "b")
	f.Add("", "x")

	f.Fuzz(func(t *testing.T, a, b string) {
		if EqualHashes(a, b) != (a == b) {
			t.Fatalf("EqualHashes(%q, %q) disagrees with ==", a, b)
		}
		if EqualHashes(a, b) != EqualHashes(b, a) {
			t.Fatalf("EqualHashes is not symmetric for %q, %q", a, b)
		}
	})
}