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

// InsertRequest writes one Request into the database.
//
// Signature note: ctx context.Context is always the first argument on
// any function that talks to the database. It's how a caller can
// cancel or time-bound this specific call without InsertRequest
// itself needing any opinion on why. You'll pass ctx straight through
// to db.Exec below.
func InsertRequest(ctx context.Context, db *pgxpool.Pool, r Request) error {

	_, err := db.Exec(ctx, )
	// TODO: db.Exec(ctx, sql, args...) running an INSERT.
	//
	// Decide deliberately: name every column explicitly in the INSERT
	// (INSERT INTO requests (id, model, ...) VALUES ($1, $2, ...)) or
	// rely on column order matching the table. Naming them explicitly
	// means a future migration that reorders columns can't silently
	// scramble your data, that's worth the extra typing.
	//
	// Use $1, $2, $3... placeholders for r's fields, never string
	// concatenation, that's not style, that's what stops SQL
	// injection.
	return nil
}

// ListRequests returns up to limit requests, most recent first.
func ListRequests(ctx context.Context, db *pgxpool.Pool, limit int) ([]Request, error) {
	// TODO: db.Query(ctx, sql, limit) running a SELECT with
	// ORDER BY created_at DESC LIMIT $1, matching the index you
	// already built in the migration.
	//
	// Check the error from Query, then immediately:
	//     defer rows.Close()
	// Skipping this leaks a connection back to the pool every single
	// call. Under real traffic that's how pgxpool quietly runs out of
	// connections with no obvious cause days later.
	//
	// Loop with rows.Next(), and inside the loop, rows.Scan(&r.ID,
	// &r.Model, ...) in the exact same column order as your SELECT.
	//
	// Initialize your result as:
	//     out := make([]Request, 0)
	// not `var out []Request`. Both range identically in Go, but they
	// encode differently to JSON: nil becomes null, an empty slice
	// becomes []. "No requests yet" should always be [], not null, or
	// your dashboard's frontend has to special-case it.
	return nil, nil
}

var _ = time.Now // remove once CreatedAt is wired into real logic above