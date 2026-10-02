// Package auth provides API key generation and hashing for ARVIS machine callers (employees, tools, services, AI agents). Human SSO is out of scope
//
// Purpose:
// 1. Issue a random key once, and store only its SHA-256 digest.
// 2. Turn an incomming key into the same digest so the proxy can resolve to an identity.
//
// Rules that must never be broken:
// 1. Raw keys are never logged, persisted, or kept in memory after issuance.
// 2. Only the output of HashKey is stored or queried.
// 3. The "arvis_" prefix is presentation only and is never part of the hash.
// 4. This package never touches the database, HTTP, or config.
// 5. Hashes are compared with EqualHashes, Never with ==.
//
// Dependencies: stardard library only (crypto/rand, crypto/sha256, crypto/subtle, encoding/hex).
// 
// Usage:
// 1. cmd/identity: GenerateKey at issuance, store hash, show raw ket once.
// 2. proxy/auth.go: HashKey on every request, the look up by digest.
// 
// Thread safety: stateless, no globals, safe for concurrent use.
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
// 32 random bytes (256 bits) come from OS entropy via crypto/rand.
// The prefix is added after hashing, so the hash covers only the secret.
func GenerateKey() (raw string, hash string, err error) {
	buf := make([]byte, 32)

	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("failed to generate key: %w", err)
	}
	secretHex := hex.EncodeToString(buf)
	raw = keyPrefix + secretHex
	
	return raw, HashKey(secretHex), nil
}

// HashKey returns the hex SHA-256 of a key's secrete part. 
// It accepts the full key or the bare secret and strips one leading "arvis_" first, so prefix changes never invalidate stored hashes.
//Plain SHA-256 is intentional: keys are 256-bit random, so slow password hashing adds costs on the hot path with no security gain.
func HashKey(rawOrSecret string) string {
	secret := strings.TrimPrefix(rawOrSecret, keyPrefix)
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}

// EqualHashes compares two hex hashes in constant time.
// Use it for every has comparison done in Go code.
func EqualHashes(hashA, hashB string) bool {
	return subtle.ConstantTimeCompare([]byte(hashA), []byte(hashB)) == 1
}
