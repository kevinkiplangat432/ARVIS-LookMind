package api

import (
	"net/http"

	"github.com/kevinkiplangat432/arvis/internal/api/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handlers.Health)
	r.Get("/requests", handlers.ListRequests(db))
	r.Get("/anomalies", handlers.ListAnomalies(db))

	return r
}
