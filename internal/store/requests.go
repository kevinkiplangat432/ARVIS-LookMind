package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Request struct {
	ID           string    `json:"id"`
	Model        string    `json:"model"`
	PromptTokens int       `json:"prompt_tokens"`
	CompTokens   int       `json:"completion_tokens"`
	LatencyMs    int       `json:"latency_ms"`
	StatusCode   int       `json:"status_code"`
	CreatedAt    time.Time `json:"created_at"`
}

func InsertRequest(ctx context.Context, db *pgxpool.Pool, r Request) error {
	_, err := db.Exec(ctx,
		`INSERT INTO requests (id, model, prompt_tokens, completion_tokens, latency_ms, status_code, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		r.ID, r.Model, r.PromptTokens, r.CompTokens, r.LatencyMs, r.StatusCode, r.CreatedAt,
	)
	return err
}

func ListRequests(ctx context.Context, db *pgxpool.Pool, limit int) ([]Request, error) {
	rows, err := db.Query(ctx,
		`SELECT id, model, prompt_tokens, completion_tokens, latency_ms, status_code, created_at
		 FROM requests ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Request
	for rows.Next() {
		var r Request
		if err := rows.Scan(&r.ID, &r.Model, &r.PromptTokens, &r.CompTokens, &r.LatencyMs, &r.StatusCode, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}
