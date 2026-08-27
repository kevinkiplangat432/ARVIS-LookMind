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
	Category  string    `json:"category"`
	Severity  string    `json:"severity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func InsertAnomaly(ctx context.Context, db *pgxpool.Pool, a Anomaly) error {
	_, err := db.Exec(ctx,
		`INSERT INTO anomalies (id, request_id, rule, detail, category, severity, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		a.ID, a.RequestID, a.Rule, a.Detail, a.Category, a.Severity, a.Status, a.CreatedAt,
	)
	return err
}

func ListAnomalies(ctx context.Context, db *pgxpool.Pool, limit int) ([]Anomaly, error) {
	rows, err := db.Query(ctx,
		`SELECT id, request_id, rule, detail, category, severity, status, created_at
		 FROM anomalies ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Anomaly, 0)
	for rows.Next() {
		var a Anomaly
		if err := rows.Scan(&a.ID, &a.RequestID, &a.Rule, &a.Detail, &a.Category, &a.Severity, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

// UpdateAnomalyStatus is the only mutation this package allows on the
// anomalies table — and it only ever touches status. rule, detail,
// category, severity, and created_at stay untouched forever, since
// those are the actual audit record. "Reviewed" or "dismissed" is a
// human's judgment about the record, not a correction to it.
func UpdateAnomalyStatus(ctx context.Context, db *pgxpool.Pool, id, status string) error {
	_, err := db.Exec(ctx, `UPDATE anomalies SET status = $1 WHERE id = $2`, status, id)
	return err
}

type AnomalyFilter struct {
	Category string
	Severity string
	Status   string
}

func ListAnomaliesFiltered(ctx context.Context, db *pgxpool.Pool, f AnomalyFilter, limit int) ([]Anomaly, error) {
	rows, err := db.Query(ctx,
		`SELECT id, request_id, rule, detail, category, severity, status, created_at
		 FROM anomalies
		 WHERE ($1 = '' OR category = $1)
		   AND ($2 = '' OR severity = $2)
		   AND ($3 = '' OR status = $3)
		 ORDER BY created_at DESC LIMIT $4`,
		f.Category, f.Severity, f.Status, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Anomaly, 0)
	for rows.Next() {
		var a Anomaly
		if err := rows.Scan(&a.ID, &a.RequestID, &a.Rule, &a.Detail, &a.Category, &a.Severity, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}