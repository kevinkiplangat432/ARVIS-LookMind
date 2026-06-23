package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Anomaly struct {
	ID        string    `json:"id"`
	RequestID string    `json:"request_id"`
	Rule      string    `json:"rule"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

func InsertAnomaly(ctx context.Context, db *pgxpool.Pool, a Anomaly) error {
	_, err := db.Exec(ctx,
		`INSERT INTO anomalies (id, request_id, rule, detail, created_at) VALUES ($1,$2,$3,$4,$5)`,
		a.ID, a.RequestID, a.Rule, a.Detail, a.CreatedAt,
	)
	return err
}

func ListAnomalies(ctx context.Context, db *pgxpool.Pool, limit int) ([]Anomaly, error) {
	rows, err := db.Query(ctx,
		`SELECT id, request_id, rule, detail, created_at FROM anomalies ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Anomaly, 0)
	for rows.Next() {
		var a Anomaly
		if err := rows.Scan(&a.ID, &a.RequestID, &a.Rule, &a.Detail, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}