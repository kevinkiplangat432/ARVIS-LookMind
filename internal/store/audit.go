package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ListRequestsInRange and ListAnomaliesInRange exist specifically for
// audit exports — everything else in the store package serves the
// live dashboard/proxy path (recent-N, via limit). An audit report
// needs a bounded date window instead, optionally scoped to one
// identity. identityID is nil for an org-wide export.

func ListRequestsInRange(ctx context.Context, db *pgxpool.Pool, identityID *string, from, to time.Time) ([]Request, error) {
	rows, err := db.Query(ctx,
		`SELECT id, identity_id, provider, model, prompt_tokens, completion_tokens, latency_ms, status_code, created_at
		 FROM requests
		 WHERE created_at BETWEEN $1 AND $2
		   AND ($3::text IS NULL OR identity_id = $3)
		 ORDER BY created_at DESC`,
		from, to, identityID,
	)
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

func ListAnomaliesInRange(ctx context.Context, db *pgxpool.Pool, identityID *string, from, to time.Time) ([]Anomaly, error) {
	rows, err := db.Query(ctx,
		`SELECT a.id, a.request_id, a.rule, a.detail, a.category, a.severity, a.status, a.created_at
		 FROM anomalies a
		 JOIN requests r ON r.id = a.request_id
		 WHERE r.created_at BETWEEN $1 AND $2
		   AND ($3::text IS NULL OR r.identity_id = $3)
		 ORDER BY a.created_at DESC`,
		from, to, identityID,
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