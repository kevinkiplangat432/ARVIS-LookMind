// Package auth provides secure cryptographic key generation and deterministic
// hashing mechanisms for the ARVIS system.
//
// This lives in its own package because both the identity CLI command and
// the proxy's incoming-request auth need the exact same behavior.
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"crypto/subtle"
)

const keyPrefix = "arvis_"

// GenerateKey returns a new random ARVIS key ready to hand to a caller,
// plus its SHA-256 hash ready to store. 
//
// The raw key is never persisted anywhere; only the output of [HashKey] 
// is ever written to the database.
func GenerateKey() (raw string, hash string, err error) {
	buf := make([]byte, 32)

	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("failed to generate key: %w", err)
	}
	secretHex := hex.EncodeToString(buf)
	raw = keyPrefix + secretHex
	
	return raw, HashKey(secretHex), nil
}

// HashKey hashes a raw key or secret payload deterministically, so an incoming request's
// key can be looked up by comparing hashes. 
//
// The raw key itself is never stored or compared directly.
func HashKey(rawOrSecret string) string {
	secret := strings.TrimPrefix(rawOrSecret, keyPrefix)
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}

// EqualHashes compares two hex hashes in constant time to prevent timing leaks.
func EqualHashes(hashA, hashB string) bool {
	return subtle.ConstantTimeCompare([]byte(hashA), []byte(hashB)) == 1
}
