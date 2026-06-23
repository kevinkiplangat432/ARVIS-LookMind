package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/kevinkiplangat432/arvis/internal/store"
)

func TestInsertAnomaly_AllFieldsPersisted(t *testing.T) {
	db := testDB(t)
	cleanRequests(t, db)

	req := store.Request{
		ID:         "req-for-anomaly",
		Model:      "gpt-4",
		StatusCode: 200,
		CreatedAt:  time.Now().UTC(),
	}
	if err := store.InsertRequest(context.Background(), db, req); err != nil {
		t.Fatalf("setup InsertRequest failed: %v", err)
	}

	anomaly := store.Anomaly{
		ID:        "anomaly-001",
		RequestID: req.ID,
		Rule:      "client_error",
		Detail:    "upstream returned 403",
		CreatedAt: time.Now().UTC(),
	}
	if err := store.InsertAnomaly(context.Background(), db, anomaly); err != nil {
		t.Fatalf("InsertAnomaly failed: %v", err)
	}

	rows, err := store.ListAnomalies(context.Background(), db, 1)
	if err != nil {
		t.Fatalf("ListAnomalies failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(rows))
	}

	got := rows[0]
	if got.ID != anomaly.ID {
		t.Errorf("ID: expected %q, got %q", anomaly.ID, got.ID)
	}
	if got.RequestID != anomaly.RequestID {
		t.Errorf("RequestID: expected %q, got %q", anomaly.RequestID, got.RequestID)
	}
	if got.Rule != anomaly.Rule {
		t.Errorf("Rule: expected %q, got %q", anomaly.Rule, got.Rule)
	}
	if got.Detail != anomaly.Detail {
		t.Errorf("Detail: expected %q, got %q", anomaly.Detail, got.Detail)
	}
}

func TestInsertAnomaly_OrphanedRequestIDFails(t *testing.T) {
	db := testDB(t)
	cleanRequests(t, db)

	anomaly := store.Anomaly{
		ID:        "orphan-anomaly",
		RequestID: "nonexistent-request-id",
		Rule:      "client_error",
		Detail:    "upstream returned 404",
		CreatedAt: time.Now().UTC(),
	}

	err := store.InsertAnomaly(context.Background(), db, anomaly)
	if err == nil {
		t.Error("expected foreign key violation, got nil error — audit trail integrity is broken")
	}
}

func TestInsertAnomaly_DuplicateIDFails(t *testing.T) {
	db := testDB(t)
	cleanRequests(t, db)

	req := store.Request{
		ID:         "req-dup-anomaly",
		Model:      "gpt-4",
		StatusCode: 200,
		CreatedAt:  time.Now().UTC(),
	}
	if err := store.InsertRequest(context.Background(), db, req); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	anomaly := store.Anomaly{
		ID:        "dup-anomaly-id",
		RequestID: req.ID,
		Rule:      "client_error",
		Detail:    "upstream returned 403",
		CreatedAt: time.Now().UTC(),
	}

	if err := store.InsertAnomaly(context.Background(), db, anomaly); err != nil {
		t.Fatalf("first insert failed: %v", err)
	}
	if err := store.InsertAnomaly(context.Background(), db, anomaly); err == nil {
		t.Error("expected duplicate anomaly ID to fail, got nil error")
	}
}

func TestListAnomalies_OrderedByCreatedAtDesc(t *testing.T) {
	db := testDB(t)
	cleanRequests(t, db)

	req := store.Request{
		ID:         "req-order-test",
		Model:      "gpt-4",
		StatusCode: 200,
		CreatedAt:  time.Now().UTC(),
	}
	if err := store.InsertRequest(context.Background(), db, req); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	base := time.Now().UTC()
	anomalies := []store.Anomaly{
		{ID: "a-oldest", RequestID: req.ID, Rule: "high_latency", Detail: "latency exceeded 10s", CreatedAt: base.Add(-2 * time.Second)},
		{ID: "a-middle", RequestID: req.ID, Rule: "client_error", Detail: "upstream returned 404", CreatedAt: base.Add(-1 * time.Second)},
		{ID: "a-newest", RequestID: req.ID, Rule: "upstream_error", Detail: "upstream returned 500", CreatedAt: base},
	}

	for _, a := range anomalies {
		if err := store.InsertAnomaly(context.Background(), db, a); err != nil {
			t.Fatalf("insert failed for %s: %v", a.ID, err)
		}
	}

	rows, err := store.ListAnomalies(context.Background(), db, 10)
	if err != nil {
		t.Fatalf("ListAnomalies failed: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}
	if rows[0].ID != "a-newest" {
		t.Errorf("expected newest first, got %q", rows[0].ID)
	}
	if rows[2].ID != "a-oldest" {
		t.Errorf("expected oldest last, got %q", rows[2].ID)
	}
}

func TestListAnomalies_EmptyTableReturnsEmptySlice(t *testing.T) {
	db := testDB(t)
	cleanRequests(t, db)

	rows, err := store.ListAnomalies(context.Background(), db, 10)
	if err != nil {
		t.Fatalf("ListAnomalies failed: %v", err)
	}
	if rows == nil {
		t.Error("expected empty slice, got nil — JSON will encode as null instead of []")
	}
}