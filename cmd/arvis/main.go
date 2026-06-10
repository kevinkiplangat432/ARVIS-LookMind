package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	fmt.Println(`

	 █████╗  ██████╗  ██╗   ██╗ ██╗ ███████╗
	██╔══██╗ ██╔══██╗ ██║   ██║ ██║ ██╔════╝
	███████║ ██████╔╝ ██║   ██║ ██║ ███████╗
	██╔══██║ ██╔══██╗ ╚██╗ ██╔╝ ██║ ╚════██║
	██║  ██║ ██║  ██║  ╚████╔╝  ██║ ███████║
	╚═╝  ╚═╝ ╚═╝  ╚═╝   ╚═══╝   ╚═╝ ╚══════╝
    
  AI Runtime Visibility & Intelligence System
  Version 0.1.0
	`)


	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"alive","service":"arvis","version":"0.1.0"}`))
	})

	fmt.Println("ARVIS API listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}