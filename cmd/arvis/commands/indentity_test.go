// identity_test.go guards the rules in identity.go: name validation, the
// key is delivered only after the identity is stored, ambiguous inserts are
// reconciled, and an undeliverable key rolls the identity back.
//
// Not covered: the real Postgres adapter and cobra wiring, which are thin
// pass throughs exercised by store tests and manual runs.
//
// Run: go test ./cmd/arvis/commands/...
// Needs: nothing, storage is faked.

package commands

import (
	"bytes"
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"

	"github.com/kevinkiplangat432/arvis/internal/auth"
	"github.com/kevinkiplangat432/arvis/internal/store"
)

type fakeStore struct {
	rows               map[string]store.Identity
	insertErr          error
	commitDespiteError bool
	getErr             error
	deleteErr          error
}

func newFakeStore() *fakeStore { return &fakeStore{rows: map[string]store.Identity{}} }

func (f *fakeStore) Insert(_ context.Context, i store.Identity) error {
	if f.insertErr == nil || f.commitDespiteError {
		f.rows[i.ID] = i
	}
	return f.insertErr
}

func (f *fakeStore) Get(_ context.Context, id string) (*store.Identity, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	i, ok := f.rows[id]
	if !ok {
		return nil, pgx.ErrNoRows
	}
	return &i, nil
}

func (f *fakeStore) Delete(_ context.Context, id string) error {
	if f.deleteErr != nil {
		return f.deleteErr
	}
	delete(f.rows, id)
	return nil
}

type brokenWriter struct{}

func (brokenWriter) Write([]byte) (int, error) { return 0, errors.New("broken pipe") }

var keyPattern = regexp.MustCompile(`arvis_[0-9a-f]{64}`)

func onlyRow(t *testing.T, f *fakeStore) store.Identity {
	t.Helper()
	if len(f.rows) != 1 {
		t.Fatalf("expected exactly one stored identity, got %d", len(f.rows))
	}
	var row store.Identity
	for _, r := range f.rows {
		row = r
	}
	return row
}

func TestValidateIdentityName(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{"plain name", "safaricom-billing", "safaricom-billing", false},
		{"trims spaces", "  jane doe  ", "jane doe", false},
		{"empty", "", "", true},
		{"only spaces", "   ", "", true},
		{"only tabs and newlines", "\t\n", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateIdentityName(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIssueIdentityHappyPath(t *testing.T) {
	f := newFakeStore()
	var out, errOut bytes.Buffer

	if err := issueIdentity(context.Background(), f, &out, &errOut, "svc"); err != nil {
		t.Fatalf("issueIdentity: %v", err)
	}

	key := keyPattern.FindString(out.String())
	if key == "" {
		t.Fatalf("no key in output: %q", out.String())
	}
	row := onlyRow(t, f)
	if row.KeyHash != auth.HashKey(key) {
		t.Error("stored hash does not match the printed key")
	}
	if strings.Contains(row.KeyHash, key) {
		t.Error("raw key found in stored value")
	}
	if !strings.Contains(out.String(), row.ID) {
		t.Error("identity ID missing from output")
	}
}

func TestIssueIdentityInsertFailsNothingPrinted(t *testing.T) {
	f := newFakeStore()
	f.insertErr = errors.New("constraint violation")
	var out, errOut bytes.Buffer

	if err := issueIdentity(context.Background(), f, &out, &errOut, "svc"); err == nil {
		t.Fatal("expected an error")
	}
	if out.Len() != 0 {
		t.Errorf("key was printed despite failed insert: %q", out.String())
	}
	if len(f.rows) != 0 {
		t.Error("a row exists after a definite failure")
	}
}

func TestIssueIdentityAmbiguousInsertThatCommitted(t *testing.T) {
	f := newFakeStore()
	f.insertErr = errors.New("connection reset")
	f.commitDespiteError = true
	var out, errOut bytes.Buffer

	if err := issueIdentity(context.Background(), f, &out, &errOut, "svc"); err != nil {
		t.Fatalf("committed insert should succeed after reconcile: %v", err)
	}
	if keyPattern.FindString(out.String()) == "" {
		t.Error("key was not delivered after a reconciled insert")
	}
	onlyRow(t, f)
}

func TestIssueIdentityAmbiguousInsertUnresolved(t *testing.T) {
	f := newFakeStore()
	f.insertErr = errors.New("connection reset")
	f.commitDespiteError = true
	f.getErr = errors.New("still down")
	var out, errOut bytes.Buffer

	err := issueIdentity(context.Background(), f, &out, &errOut, "svc")
	if err == nil {
		t.Fatal("expected an error")
	}
	if out.Len() != 0 {
		t.Error("key printed although the outcome was unconfirmed")
	}
	if !strings.Contains(err.Error(), onlyRow(t, f).ID) {
		t.Error("error should name the identity that may exist")
	}
}

func TestIssueIdentityDeliveryFailureRollsBack(t *testing.T) {
	f := newFakeStore()
	var errOut bytes.Buffer

	err := issueIdentity(context.Background(), f, brokenWriter{}, &errOut, "svc")
	if err == nil {
		t.Fatal("expected an error")
	}
	if len(f.rows) != 0 {
		t.Error("identity was not rolled back")
	}
}

func TestIssueIdentityDeliveryFailureRollbackFails(t *testing.T) {
	f := newFakeStore()
	f.deleteErr = errors.New("db gone")
	var errOut bytes.Buffer

	err := issueIdentity(context.Background(), f, brokenWriter{}, &errOut, "svc")
	if err == nil {
		t.Fatal("expected an error")
	}
	id := onlyRow(t, f).ID
	if !strings.Contains(errOut.String(), id) {
		t.Errorf("stderr should name identity %s for manual cleanup, got %q", id, errOut.String())
	}
	if !strings.Contains(err.Error(), id) {
		t.Error("returned error should name the identity too")
	}
}