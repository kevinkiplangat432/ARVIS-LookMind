package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Request struct {
	ID           string    `json:"id"`
	IdentityID   *string   `json:"identity_id,omitempty"`
	Provider     string    `json:"provider"`
	Model        string    `json:"model"`
	PromptTokens int       `json:"prompt_tokens"`
	CompTokens   int       `json:"completion_tokens"`
	LatencyMs    int       `json:"latency_ms"`
	StatusCode   int       `json:"status_code"`
	CreatedAt    time.Time `json:"created_at"`
}

func InsertRequest(ctx context.Context, db *pgxpool.Pool, r Request) error {
	_, err := db.Exec(ctx,
		`INSERT INTO requests (id, identity_id, provider, model, prompt_tokens, completion_tokens, latency_ms, status_code, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		r.ID, r.IdentityID, r.Provider, r.Model, r.PromptTokens, r.CompTokens, r.LatencyMs, r.StatusCode, r.CreatedAt,
	)
	return err
}

func ListRequests(ctx context.Context, db *pgxpool.Pool, limit int) ([]Request, error) {
	rows, err := db.Query(ctx,
		`SELECT id, identity_id, provider, model, prompt_tokens, completion_tokens, latency_ms, status_code, created_at
		 FROM requests ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Request, 0)
	for rows.Next() {
		var r Request
		if err := rows.Scan(&r.ID, &r.IdentityID, &r.Provider, &r.Model, &r.PromptTokens, &r.CompTokens, &r.LatencyMs, &r.StatusCode, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}