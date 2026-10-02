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

var ErrMissingKey = errors.New("missing or malformed authorization header")
var ErrUnknownKey = errors.New("key does not match any known identity")

// authenticate extracts the caller's ARVIS key from the Authorization
// header, hashes it, and resolves it to an Identity. This runs before
// anything else on every request — an unattributed request never
// reaches a provider, full stop.
func authenticate(ctx context.Context, db *pgxpool.Pool, r *http.Request) (*store.Identity, error) {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return nil, ErrMissingKey
	}
	rawKey := strings.TrimPrefix(header, "Bearer ")
	if rawKey == "" {
		return nil, ErrMissingKey
	}

	hashkey := auth.HashKey(rawKey)
	identity, err := store.GetIdentityByKeyHash(ctx, db, hashkey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUnknownKey
		}
		return nil, fmt.Errorf("identity lookup failed: %w", err)
	}

	if !auth.EqualHashes(identity.KeyHash, hash) {
		return nil, ErrUnknownKey
	}
	return identity, nil
}