package auth

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	raw, hash, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}

	if !strings.HasPrefix(raw, keyPrefix) {
		t.Errorf("raw key %q is missing the prefix", raw)
	}
	secret := strings.TrimPrefix(raw, keyPrefix)
	if len(secret) != 64 {
		t.Errorf("secret length = %d, want 64", len(secret))
	}
	if _, err := hex.DecodeString(secret); err != nil {
		t.Errorf("secret is not valid hex: %v", err)
	}
	if hash != HashKey(raw) {
		t.Error("returned hash does not match HashKey(raw)")
	}
	if hash != HashKey(secret) {
		t.Error("returned hash does not match HashKey(bare secret)")
	}
}

func TestGenerateKeyUnique(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		raw, _, err := GenerateKey()
		if err != nil {
			t.Fatalf("GenerateKey: %v", err)
		}
		if seen[raw] {
			t.Fatalf("duplicate key after %d generations", i)
		}
		seen[raw] = true
	}
}

func TestHashKey(t *testing.T) {
	const abc = "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	const empty = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"bare secret, known vector", "abc", abc},
		{"prefixed hashes same as bare", "arvis_abc", abc},
		{"empty input", "", empty},
		{"prefix only equals empty", "arvis_", empty},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HashKey(tt.in); got != tt.want {
				t.Errorf("HashKey(%q) = %s, want %s", tt.in, got, tt.want)
			}
		})
	}
}

func TestHashKeyStripsOnlyOnePrefix(t *testing.T) {
	if HashKey("arvis_arvis_abc") == HashKey("abc") {
		t.Error("HashKey stripped more than one prefix")
	}
}

func TestHashKeyFormat(t *testing.T) {
	h := HashKey("anything")
	if len(h) != 64 {
		t.Errorf("hash length = %d, want 64", len(h))
	}
	if h != strings.ToLower(h) {
		t.Error("hash is not lowercase hex")
	}
}

func TestEqualHashes(t *testing.T) {
	a := HashKey("one")
	b := HashKey("two")

	tests := []struct {
		name string
		x, y string
		want bool
	}{
		{"identical", a, a, true},
		{"different", a, b, false},
		{"different length", a, a[:10], false},
		{"both empty", "", "", true},
		{"empty versus hash", "", a, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EqualHashes(tt.x, tt.y); got != tt.want {
				t.Errorf("EqualHashes(%q, %q) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}