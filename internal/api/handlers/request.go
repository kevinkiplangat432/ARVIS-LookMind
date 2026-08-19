package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevinkiplangat432/arvis/internal/store"
)

// Requests returns an isolated router for request-log endpoints, same
// pattern as Health() in health.go, an isolated sub-router that gets
// mounted onto the main one.
func Requests(db *pgxpool.Pool) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		// TODO: call store.ListRequests(context, db, limit).
		//
		// Where should limit come from? Hardcoding it works for now,
		// but think about whether a query parameter (?limit=50) makes
		// more sense for something a dashboard will call repeatedly.
		//
		// TODO: handle the error case first (what status code should
		// a failed database query return to the caller?), then encode
		// the successful result as JSON with json.NewEncoder(w).Encode(...).
		//
		// Don't forget w.Header().Set("Content-Type", "application/json")
		// before writing, same as health.go already does.
		_ = context.Background()
		_ = json.NewEncoder
		_ = store.ListRequests
	})

	return r
}