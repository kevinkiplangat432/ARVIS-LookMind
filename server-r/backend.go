package main

import (
	"net/http"
	"log/slog"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Hello, from the backend!")
}

func main(){
	
	slog.Info("starting server........")
	http.HandleFunc("/", HomeHandler)
	http.ListenAndServe(":5555", nil)
}