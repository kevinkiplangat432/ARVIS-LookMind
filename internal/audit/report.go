package audit

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kevinkiplangat432/arvis/internal/store"
)

// Report is everything a PDF export needs, gathered once so pdf.go
// only has to render, never query.
type Report struct {
	Identity    *store.Identity // nil for an org-wide report
	From        time.Time
	To          time.Time
	GeneratedAt time.Time
	Requests    []store.Request
	Anomalies   []store.Anomaly
	Summary     Summary
}

type Summary struct {
	TotalRequests    int
	TotalAnomalies   int
	BlockedCount     int // pre-flight policy blocks, status 403
	KilledCount      int // kill-switch terminations, status 499
	ByCategory       map[string]int
	TotalPromptTok   int
	TotalCompleteTok int
}

func BuildReport(ctx context.Context, db *pgxpool.Pool, identityID *string, from, to time.Time) (*Report, error) {
	var identity *store.Identity
	if identityID != nil {
		i, err := store.GetIdentityByID(ctx, db, *identityID)
		if err != nil {
			return nil, fmt.Errorf("failed to look up identity: %w", err)
		}
		identity = i
	}

	requests, err := store.ListRequestsInRange(ctx, db, identityID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to list requests: %w", err)
	}

	anomalies, err := store.ListAnomaliesInRange(ctx, db, identityID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to list anomalies: %w", err)
	}

	summary := Summary{
		TotalRequests:  len(requests),
		TotalAnomalies: len(anomalies),
		ByCategory:     make(map[string]int),
	}
	for _, r := range requests {
		summary.TotalPromptTok += r.PromptTokens
		summary.TotalCompleteTok += r.CompTokens
		switch r.StatusCode {
		case 403:
			summary.BlockedCount++
		case 499:
			summary.KilledCount++
		}
	}
	for _, a := range anomalies {
		summary.ByCategory[a.Category]++
	}

	return &Report{
		Identity:    identity,
		From:        from,
		To:          to,
		GeneratedAt: time.Now(),
		Requests:    requests,
		Anomalies:   anomalies,
		Summary:     summary,
	}, nil
}