// auth.go is the front door of the proxy. It resolves the caller's ARVIS key to an identity before anything else happens/
// 
// Rules that must never be broken:
// 	1. authenticate runs first on every request. No identity, no body read, no policy check, no providers call.
// 	2. Its fails closed. Any error rejects the request.
// 	3. The raw key is never logged, returned, or placed inside an error.
//
//  Dependencies: internal/auth (hashing) and internal/store (lookup).
// Must never import detertor, policy, or tokenize. 
// Usage: called once per request from proxy.ServeHttp.
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
	ErrMissingKey = errors.New("missing or malformed authorization header")
	ErrUnknownKey = errors.New("key does not match any known identity")
	ErrAuthUnavailable = errors.New("authentication temporarily unavailable")
)

// authenticate resolves the request's Bearer key to an identity.
// Unknown keys are the caller's problem (401). A failing database is ours.
// and telling a customer their key is wrong while Postgres burns, is rude.
func authenticate(ctx context.Context, db *pgxpool.Pool, r *http.Request) (*store.Identity, error) {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return nil, ErrMissingKey
	}
	rawKey := strings.TrimPrefix(header, "Bearer ")
	if rawKey == "" {
		return nil, ErrMissingKey
	}

	hash := auth.HashKey(rawKey)
	identity, err := store.GetIdentityByKeyHash(ctx, db, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUnknownKey
		}
		return nil, fmt.Errorf("identity lookup failed: %w", err)
	}
	// defense in depth: confirm the fetched hash in constant time.
	if !auth.EqualHashes(identity.KeyHash, hash) {
		return nil, ErrUnknownKey
	}
	return identity, nil
}