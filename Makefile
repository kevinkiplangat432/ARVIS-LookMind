.PHONY: dev build test down seed

dev:
	docker compose up -d

build:
	go build ./cmd/arvis

test:
	go test ./...

down:
	docker compose down

seed:
	go run scripts/seed.go

migrate:
	psql $$DATABASE_URL -f migrations/001_create_requests.sql
	psql $$DATABASE_URL -f migrations/002_create_anomalies.sql
