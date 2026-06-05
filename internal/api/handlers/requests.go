package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/kevinkiplangat432/arvis/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ListRequests(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 50
		}
		rows, err := store.ListRequests(r.Context(), db, limit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	}
}
