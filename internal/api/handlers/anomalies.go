package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevinkiplangat432/arvis/internal/store"
)

// Anomalies returns an isolated router for anomaly endpoints, mirroring
// the pattern in health.go and request.go.
func Anomalies(db *pgxpool.Pool) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		limit := 50
		if v := req.URL.Query().Get("limit"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		results, err := store.ListAnomalies(req.Context(), db, limit)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to list anomalies")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(results)
	})

	return r
}