.PHONY: dev build test down seed

dev:
	docker compose up -d
run:
	go run ./cmd/arvis/ server

build:
	go build -o arvis ./cmd/arvis/

test:
	TEST_DATABASE_URL="postgres://arvis:arvis@localhost:5432/arvis_test?sslmode=disable" go test ./... -v

test-unit:
	go test ./internal/detector/... -v

test-integration:
	TEST_DATABASE_URL="postgres://arvis:arvis@localhost:5432/arvis_test?sslmode=disable" go test ./internal/store/... -v

down:
	docker compose down

seed:
	go run scripts/seed.go

migrate:
	psql $$DATABASE_URL -f migrations/001_create_requests.sql
	psql $$DATABASE_URL -f migrations/002_create_anomalies.sql
