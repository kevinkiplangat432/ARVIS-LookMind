// Package auth provides API key generation and hashing for ARVIS machine
// callers (employee tools, services, AI agents). Human SSO is out of scope.
//
// Purpose:
//  1. Issue a random key once and store only its SHA-256 digest.
//  2. Turn an incoming key into the same digest so the proxy can resolve it
//     to an identity.
//
// Rules that must never be broken:
//  1. Raw keys are never logged, persisted, or kept after issuance.
//  2. Only the output of HashKey is stored or queried.
//  3. The "arvis_" prefix is presentation only and is never part of the hash.
//  4. Hashes are compared with EqualHashes, never with ==.
//  5. This package never touches the database, HTTP, or config.
//
// Dependencies: standard library only. Must never import store, proxy, or config.
// Usage: cmd/identity calls GenerateKey at issuance. proxy/auth.go calls
// HashKey on every request.
// Thread safety: stateless, no globals, safe for concurrent use.
//
// Docs: docs/auth_and_identity.md
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"
)

const keyPrefix = "arvis_"

// GenerateKey returns a new key for the caller and its hash for storage.
// The key is 32 random bytes from OS entropy. The prefix is added after
// hashing input is chosen, so the hash covers only the secret.
// Keys are shown once. Lose one and you get a new key, not a refund.
func GenerateKey() (raw string, hash string, err error) {
	buf := make([]byte, 32)

	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("failed to generate key: %w", err)
	}
	secretHex := hex.EncodeToString(buf)
	raw = keyPrefix + secretHex

	return raw, HashKey(secretHex), nil
}

// HashKey returns the hex SHA-256 of a key's secret part. It accepts the
// full key or the bare secret and strips one leading "arvis_" first, so
// prefix changes never invalidate stored hashes.
// Plain SHA-256 is deliberate: the input is 256 random bits, so a slow
// password hash would add latency to every request and no security.
func HashKey(rawOrSecret string) string {
	secret := strings.TrimPrefix(rawOrSecret, keyPrefix)
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}

// EqualHashes compares two hex hashes in constant time.
func EqualHashes(hashA, hashB string) bool {
	return subtle.ConstantTimeCompare([]byte(hashA), []byte(hashB)) == 1
}