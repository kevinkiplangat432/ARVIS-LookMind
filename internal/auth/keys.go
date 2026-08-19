package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateKey returns a new random ARVIS key ready to hand to a caller,
// plus its SHA-256 hash ready to store. The raw key is never persisted
// anywhere — only HashKey's output is ever written to the database. This
// lives in its own package because both the identity CLI command and
// the proxy's incoming-request auth (Phase 5) need the exact same
// hashing behavior.
func GenerateKey() (raw string, hash string, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("failed to generate key: %w", err)
	}
	raw = "arvis_" + hex.EncodeToString(buf)
	return raw, HashKey(raw), nil
}

// HashKey hashes a raw key deterministically, so an incoming request's
// key can be looked up by comparing hashes — the raw key itself is
// never stored or compared directly.
func HashKey(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}