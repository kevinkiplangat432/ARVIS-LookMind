package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevinkiplangat432/arvis/internal/store"
)

// Requests returns an isolated router for request-log endpoints, same
// pattern as Health() in health.go — an isolated sub-router that gets
// mounted onto the main one.
func Requests(db *pgxpool.Pool) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		limit := 50
		if v := req.URL.Query().Get("limit"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		results, err := store.ListRequests(req.Context(), db, limit)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to list requests")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results)
	})

	return r
}

// writeJSONError is a small shared helper so every handler returns
// errors in the same shape. Lives here since request.go is the first
// handler that needs it; anomalies.go will reuse it.
func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}