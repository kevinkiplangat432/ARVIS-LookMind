package store

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestIdentityNeverSerializesKeyHash(t *testing.T) {
	i := Identity{ID: "1", Name: "svc", KeyHash: "deadbeef", CreatedAt: time.Now()}

	out, err := json.Marshal(i)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	s := string(out)
	if strings.Contains(s, "deadbeef") || strings.Contains(s, "key_hash") {
		t.Errorf("KeyHash leaked into JSON: %s", s)
	}
}

func testDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	db, err := Connect(url)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	t.Cleanup(db.Close)
	return db
}

func newTestIdentity(t *testing.T, db *pgxpool.Pool) Identity {
	t.Helper()
	i := Identity{
		ID:        uuid.NewString(),
		Name:      "test-identity",
		KeyHash:   uuid.NewString(),
		CreatedAt: time.Now(),
	}
	if err := InsertIdentity(context.Background(), db, i); err != nil {
		t.Fatalf("InsertIdentity: %v", err)
	}
	t.Cleanup(func() {
		db.Exec(context.Background(), `DELETE FROM identities WHERE id = $1`, i.ID)
	})
	return i
}

func TestIdentityRoundTrip(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()
	want := newTestIdentity(t, db)

	byHash, err := GetIdentityByKeyHash(ctx, db, want.KeyHash)
	if err != nil {
		t.Fatalf("GetIdentityByKeyHash: %v", err)
	}
	byID, err := GetIdentityByID(ctx, db, want.ID)
	if err != nil {
		t.Fatalf("GetIdentityByID: %v", err)
	}

	for name, got := range map[string]*Identity{"by hash": byHash, "by id": byID} {
		if got.ID != want.ID || got.Name != want.Name || got.KeyHash != want.KeyHash {
			t.Errorf("%s: got %+v, want %+v", name, got, want)
		}
	}
}

func TestIdentityLookupMissReturnsErrNoRows(t *testing.T) {
	db := testDB(t)
	ctx := context.Background()

	if _, err := GetIdentityByKeyHash(ctx, db, "no-such-hash"); !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("by hash: err = %v, want pgx.ErrNoRows", err)
	}
	if _, err := GetIdentityByID(ctx, db, uuid.NewString()); !errors.Is(err, pgx.ErrNoRows) {
		t.Errorf("by id: err = %v, want pgx.ErrNoRows", err)
	}
}

// If this fails, key_hash is missing its UNIQUE constraint.
func TestInsertDuplicateKeyHashRejected(t *testing.T) {
	db := testDB(t)
	first := newTestIdentity(t, db)

	dup := Identity{ID: uuid.NewString(), Name: "dup", KeyHash: first.KeyHash, CreatedAt: time.Now()}
	if err := InsertIdentity(context.Background(), db, dup); err == nil {
		db.Exec(context.Background(), `DELETE FROM identities WHERE id = $1`, dup.ID)
		t.Error("duplicate key_hash was accepted")
	}
}