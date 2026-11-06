// auth.go is the front door of the proxy. It resolves the caller's ARVIS
// key to an Identity before anything else happens.
//
// Rules that must never be broken:
//  1. authenticate runs first on every request. No identity, no body read,
//     no policy check, no provider call.
//  2. It fails closed. Any error rejects the request.
//  3. The raw key is never logged, returned, or placed inside an error.
//  4. The Bearer scheme is matched case-insensitively (RFC 7235).
//     Strict clients and sloppy clients both get in, the key still decides.
//
// Dependencies: internal/auth (hashing) and internal/store (lookup).
// Must never import detector, policy, or tokenize.
// Usage: called once per request from Proxy.ServeHTTP.
// Thread safety: stateless, safe for concurrent use.
//
// Docs: docs/auth_and_identity.md

package proxy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kevinkiplangat432/arvis/internal/auth"
	"github.com/kevinkiplangat432/arvis/internal/store"
)

var (
	ErrMissingKey      = errors.New("missing or malformed authorization header")
	ErrUnknownKey      = errors.New("key does not match any known identity")
	ErrAuthUnavailable = errors.New("authentication temporarily unavailable")
)

// extractBearer pulls the token out of an Authorization header.
// Pure on purpose: no database, so it is easy to test and to fuzz.
func extractBearer(header string) (string, error) {
	scheme, token, found := strings.Cut(header, " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return "", ErrMissingKey
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", ErrMissingKey
	}
	return token, nil
}

// authenticate resolves the request's Bearer key to an Identity.
// Unknown keys are the caller's problem (401). A failing database is ours,
// and telling a customer their key is wrong while Postgres burns is rude.
func authenticate(ctx context.Context, db *pgxpool.Pool, r *http.Request) (*store.Identity, error) {
	rawKey, err := extractBearer(r.Header.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	hash := auth.HashKey(rawKey)
	identity, err := store.GetIdentityByKeyHash(ctx, db, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUnknownKey
		}
		return nil, fmt.Errorf("identity lookup failed: %w", err)
	}

	// Defense in depth: confirm the fetched hash in constant time.
	if !auth.EqualHashes(identity.KeyHash, hash) {
		return nil, ErrUnknownKey
	}
	return identity, nil
}

// authErrorStatus maps an authentication error to an HTTP status.
// It returns 0 when the client already hung up: nothing to write, and
// nothing to page anyone about. A canceled request is not an outage.
func authErrorStatus(ctx context.Context, err error) int {
	switch {
	case errors.Is(err, ErrMissingKey), errors.Is(err, ErrUnknownKey):
		return http.StatusUnauthorized
	case errors.Is(err, context.Canceled), errors.Is(ctx.Err(), context.Canceled):
		return 0
	default:
		return http.StatusServiceUnavailable
	}
}