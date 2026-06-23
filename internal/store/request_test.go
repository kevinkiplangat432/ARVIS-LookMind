package store_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kevinkiplangat432/arvis/internal/store"
)

func testDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		url = "postgres://arvis:arvis@localhost:5432/arvis_test?sslmode=disable"
	}
	db, err := pgxpool.New(context.Background(), url)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func cleanRequests(t *testing.T, db *pgxpool.Pool) {
	t.Helper()
	_, err := db.Exec(context.Background(), "DELETE FROM anomalies")
	if err != nil {
		t.Fatalf("failed to clean anomalies: %v", err)
	}
	_, err = db.Exec(context.Background(), "DELETE FROM requests")
	if err != nil {
		t.Fatalf("failed to clean requests: %v", err)
	}
}

func TestInsertRequest_AllFieldsPersisted(t *testing.T) {
	db := testDB(t)
	cleanRequests(t, db)

	now := time.Now().UTC().Truncate(time.Microsecond)
	req := store.Request{
		ID:           "test-id-001",
		Model:        "gpt-4",
		PromptTokens: 100,
		CompTokens:   200,
		LatencyMs:    350,
		StatusCode:   200,
		CreatedAt:    now,
	}

	if err := store.InsertRequest(context.Background(), db, req); err != nil {
		t.Fatalf("InsertRequest failed: %v", err)
	}

	rows, err := store.ListRequests(context.Background(), db, 1)
	if err != nil {
		t.Fatalf("ListRequests failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}

	got := rows[0]
	if got.ID != req.ID {
		t.Errorf("ID: expected %q, got %q", req.ID, got.ID)
	}
	if got.Model != req.Model {
		t.Errorf("Model: expected %q, got %q", req.Model, got.Model)
	}
	if got.PromptTokens != req.PromptTokens {
		t.Errorf("PromptTokens: expected %d, got %d", req.PromptTokens, got.PromptTokens)
	}
	if got.CompTokens != req.CompTokens {
		t.Errorf("CompTokens: expected %d, got %d", req.CompTokens, got.CompTokens)
	}
	if got.LatencyMs != req.LatencyMs {
		t.Errorf("LatencyMs: expected %d, got %d", req.LatencyMs, got.LatencyMs)
	}
	if got.StatusCode != req.StatusCode {
		t.Errorf("StatusCode: expected %d, got %d", req.StatusCode, got.StatusCode)
	}
}

func TestInsertRequest_DuplicateIDFails(t *testing.T) {
	db := testDB(t)
	cleanRequests(t, db)

	req := store.Request{
		ID:         "duplicate-id",
		Model:      "gpt-4",
		StatusCode: 200,
		CreatedAt:  time.Now().UTC(),
	}

	if err := store.InsertRequest(context.Background(), db, req); err != nil {
		t.Fatalf("first insert failed: %v", err)
	}
	if err := store.InsertRequest(context.Background(), db, req); err == nil {
		t.Error("expected duplicate insert to fail, got nil error")
	}
}

func TestListRequests_OrderedByCreatedAtDesc(t *testing.T) {
	db := testDB(t)
	cleanRequests(t, db)

	base := time.Now().UTC()
	requests := []store.Request{
		{ID: "req-oldest", Model: "gpt-4", StatusCode: 200, CreatedAt: base.Add(-2 * time.Second)},
		{ID: "req-middle", Model: "gpt-4", StatusCode: 200, CreatedAt: base.Add(-1 * time.Second)},
		{ID: "req-newest", Model: "gpt-4", StatusCode: 200, CreatedAt: base},
	}

	for _, r := range requests {
		if err := store.InsertRequest(context.Background(), db, r); err != nil {
			t.Fatalf("insert failed for %s: %v", r.ID, err)
		}
	}

	rows, err := store.ListRequests(context.Background(), db, 10)
	if err != nil {
		t.Fatalf("ListRequests failed: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
	if rows[0].ID != "req-newest" {
		t.Errorf("expected newest first, got %q", rows[0].ID)
	}
	if rows[2].ID != "req-oldest" {
		t.Errorf("expected oldest last, got %q", rows[2].ID)
	}
}

func TestListRequests_LimitIsRespected(t *testing.T) {
	db := testDB(t)
	cleanRequests(t, db)

	for i := 0; i < 5; i++ {
		r := store.Request{
			ID:         "req-limit-" + string(rune('0'+i)),
			Model:      "gpt-4",
			StatusCode: 200,
			CreatedAt:  time.Now().UTC(),
		}
		if err := store.InsertRequest(context.Background(), db, r); err != nil {
			t.Fatalf("insert failed: %v", err)
		}
	}

	rows, err := store.ListRequests(context.Background(), db, 3)
	if err != nil {
		t.Fatalf("ListRequests failed: %v", err)
	}
	if len(rows) != 3 {
		t.Errorf("expected 3 rows with limit=3, got %d", len(rows))
	}
}

func TestListRequests_EmptyTableReturnsEmptySlice(t *testing.T) {
	db := testDB(t)
	cleanRequests(t, db)

	rows, err := store.ListRequests(context.Background(), db, 10)
	if err != nil {
		t.Fatalf("ListRequests failed: %v", err)
	}
	if rows == nil {
		t.Error("expected empty slice, got nil — JSON will encode as null instead of []")
	}
}