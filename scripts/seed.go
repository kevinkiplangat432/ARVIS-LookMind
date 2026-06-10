//go:build ignore

package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/kevinkiplangat432/arvis/internal/store"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	db, err := store.Connect("postgres://arvis:arvis@localhost:5432/arvis?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	models := []string{"gpt-4o", "gpt-4o-mini", "claude-3-5-sonnet"}
	for i := 0; i < 20; i++ {
		req := store.Request{
			ID:           uuid.NewString(),
			Model:        models[rand.Intn(len(models))],
			PromptTokens: rand.Intn(2000),
			CompTokens:   rand.Intn(1000),
			LatencyMs:    rand.Intn(5000),
			StatusCode:   200,
			CreatedAt:    time.Now().Add(-time.Duration(i) * time.Minute),
		}
		if err := store.InsertRequest(context.Background(), db, req); err != nil {
			log.Printf("seed request: %v", err)
		}
	}
	log.Println("seeded 20 requests")
}
