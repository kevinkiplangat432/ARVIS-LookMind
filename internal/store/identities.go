package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Identity struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	KeyHash   string    `json:"-"` // never serialized out, even internally
	CreatedAt time.Time `json:"created_at"`
}

func InsertIdentity(ctx context.Context, db *pgxpool.Pool, i Identity) error {
	_, err := db.Exec(ctx,
		`INSERT INTO identities (id, name, key_hash, created_at) VALUES ($1,$2,$3,$4)`,
		i.ID, i.Name, i.KeyHash, i.CreatedAt,
	)
	return err
}

// GetIdentityByKeyHash is what the proxy will call on every incoming
// request in Phase 5 — hot path, keeps it to a single indexed lookup.
func GetIdentityByKeyHash(ctx context.Context, db *pgxpool.Pool, keyHash string) (*Identity, error) {
	var i Identity
	err := db.QueryRow(ctx,
		`SELECT id, name, key_hash, created_at FROM identities WHERE key_hash = $1`, keyHash,
	).Scan(&i.ID, &i.Name, &i.KeyHash, &i.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

// GetIdentityByID is the audit export's lookup — reports are usually
// requested by identity ID (from the dashboard or CLI), not by key
// hash the way the proxy's auth path needs.
func GetIdentityByID(ctx context.Context, db *pgxpool.Pool, id string) (*Identity, error) {
	var i Identity
	err := db.QueryRow(ctx,
		`SELECT id, name, key_hash, created_at FROM identities WHERE id = $1`, id,
	).Scan(&i.ID, &i.Name, &i.KeyHash, &i.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &i, nil
}