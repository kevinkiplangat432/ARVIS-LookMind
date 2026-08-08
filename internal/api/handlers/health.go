package api

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// toutes return a completely isolated router for health endpoints
func Health() chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "UP"}`))
	})

	return r
}