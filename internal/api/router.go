package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevinkiplangat432/arvis/internal/api/handlers"
)

func NewRouter(db *pgxpool.Pool, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(withLogging(logger))

	r.Get("/health", handlers.Health)
	r.Get("/requests", handlers.ListRequests(db))
	r.Get("/anomalies", handlers.ListAnomalies(db))

	return r
}