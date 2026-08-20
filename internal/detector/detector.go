package detector

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kevinkiplangat432/arvis/internal/store"
)

// Detector owns every rule and is what the proxy actually talks to.
// Two separate entry points on purpose — CheckSync must finish before
// a request is forwarded, RunAsync must never block a response
// that's already gone back to the caller. Conflating them would be
// the easiest way to silently reintroduce the latency the sync/async
// split exists to avoid.
type Detector struct {
	db         *pgxpool.Pool
	logger     *slog.Logger
	syncRules  []SyncRule
	asyncRules []AsyncRule
}

func New(db *pgxpool.Pool, logger *slog.Logger) *Detector {
	return &Detector{
		db:     db,
		logger: logger,
		syncRules: []SyncRule{
			NewVolumeRule(100, time.Minute),
		},
		asyncRules: []AsyncRule{
			NewLatencyRule(50),
			NewContentRule(),
		},
	}
}

// CheckSync runs before forwarding. Returning flags doesn't block the
// request yet — deciding to actually block on a flag is a policy-
// engine decision, explicitly out of scope until governance work
// later. For now every sync flag is logged exactly like an async one,
// just computed earlier.
func (d *Detector) CheckSync(ctx context.Context, identityID, model string) []Flag {
	var flags []Flag
	for _, rule := range d.syncRules {
		flags = append(flags, rule.CheckSync(ctx, identityID, model)...)
	}
	return flags
}

// RunAsync should always be called with `go`. Runs every async rule
// and persists every flag from both sync and async rules against the
// given request ID.
func (d *Detector) RunAsync(requestID string, syncFlags []Flag, in AsyncInput) {
	ctx := context.Background()

	flags := append([]Flag{}, syncFlags...)
	for _, rule := range d.asyncRules {
		flags = append(flags, rule.CheckAsync(ctx, in)...)
	}

	for _, f := range flags {
		anomaly := store.Anomaly{
			ID:        uuid.NewString(),
			RequestID: requestID,
			Rule:      f.Rule,
			Detail:    f.Detail,
			Category:  f.Category,
			Severity:  f.Severity,
			Status:    "open",
			CreatedAt: time.Now(),
		}
		if err := store.InsertAnomaly(ctx, d.db, anomaly); err != nil {
			d.logger.Error("failed to insert anomaly", "error", err.Error())
		}
	}
}