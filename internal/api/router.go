package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kevinkiplangat432/arvis/internal/api/handlers"

)

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger) //log all incoming request
	r.Use(middleware.Recoverer) // prevent the api from crashing on panic

	// mount the sub routers
	r.Mount("/api/v1/health", handlers.Health())

	return r 

}