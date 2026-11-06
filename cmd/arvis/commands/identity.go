// identity.go is the CLI front for managing ARVIS identities.
//
// Rules that must never be broken:
//  1. The raw key is written once, to the command's output, and only after
//     the identity is confirmed stored. It is never logged, written to
//     disk, or retrievable later.
//  2. If the key cannot be delivered, the identity is rolled back. We never
//     leave an identity whose key nobody has seen.
//  3. Only the hash goes to the database.
//  4. An insert that errors is reconciled by lookup before it counts as a failure.
//
// Dependencies: internal/auth, internal/store, cobra, pgx (error values only).
// Must never import proxy.
// Usage: arvis identity create "<name>"
// Thread safety: one goroutine per CLI invocation.
//
// Treat the printed key like a house key: once handed over, we keep no spare.
//
// Docs: docs/auth_and_identity.md

package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"

	"github.com/kevinkiplangat432/arvis/internal/auth"
	"github.com/kevinkiplangat432/arvis/internal/store"
)

var identityCmd = &cobra.Command{
	Use:   "identity",
	Short: "Manage ARVIS identities (employees and services)",
}

var identityCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Register a new employee or service and issue them a key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return createIdentity(cmd, args[0])
	},
}

// identityStore is the slice of storage issueIdentity needs, a seam so the
// failure paths can be tested without Postgres.
type identityStore interface {
	Insert(ctx context.Context, i store.Identity) error
	Get(ctx context.Context, id string) (*store.Identity, error)
	Delete(ctx context.Context, id string) error
}

type pgIdentityStore struct{ db *pgxpool.Pool }

func (p pgIdentityStore) Insert(ctx context.Context, i store.Identity) error {
	return store.InsertIdentity(ctx, p.db, i)
}

func (p pgIdentityStore) Get(ctx context.Context, id string) (*store.Identity, error) {
	return store.GetIdentityByID(ctx, p.db, id)
}

func (p pgIdentityStore) Delete(ctx context.Context, id string) error {
	return store.DeleteIdentity(ctx, p.db, id)
}

// validateIdentityName trims the name and rejects blanks, so the audit
// trail never shows an identity called "".
func validateIdentityName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("identity name cannot be empty")
	}
	return name, nil
}

func createIdentity(cmd *cobra.Command, rawName string) error {
	name, err := validateIdentityName(rawName)
	if err != nil {
		return err
	}

	db, err := store.Connect(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	return issueIdentity(context.Background(), pgIdentityStore{db}, cmd.OutOrStdout(), cmd.ErrOrStderr(), name)
}

// issueIdentity registers a validated name and hands the key to out.
func issueIdentity(ctx context.Context, ids identityStore, out, errOut io.Writer, name string) error {
	rawKey, keyHash, err := auth.GenerateKey()
	if err != nil {
		return err
	}

	identity := store.Identity{
		ID:        uuid.NewString(),
		Name:      name,
		KeyHash:   keyHash,
		CreatedAt: time.Now(),
	}

	if err := ids.Insert(ctx, identity); err != nil {
		// The insert may have committed before the reply was lost, so ask
		// the database before declaring failure. The ID is fresh, so a hit
		// can only be ours.
		if _, getErr := ids.Get(ctx, identity.ID); getErr != nil {
			if errors.Is(getErr, pgx.ErrNoRows) {
				return fmt.Errorf("failed to create identity: %w", err)
			}
			return fmt.Errorf("failed to create identity and could not confirm the outcome, identity %s may exist without a key, check and remove it: %w", identity.ID, err)
		}
	}

	if _, err := io.WriteString(out, deliveryMessage(identity, rawKey)); err != nil {
		if delErr := ids.Delete(ctx, identity.ID); delErr != nil {
			fmt.Fprintf(errOut, "identity %s was created but its key was not delivered, delete it by hand\n", identity.ID)
			return fmt.Errorf("key delivery failed (%v) and rollback failed, remove identity %s by hand: %w", err, identity.ID, delErr)
		}
		return fmt.Errorf("key delivery failed, identity rolled back, safe to retry: %w", err)
	}
	return nil
}

func deliveryMessage(i store.Identity, rawKey string) string {
	return fmt.Sprintf(
		"Identity created.\n  ID:   %s\n  Name: %s\n\nKey (shown once, store it securely, ARVIS cannot recover it):\n  %s\n",
		i.ID, i.Name, rawKey,
	)
}

func init() {
	identityCmd.AddCommand(identityCreateCmd)
	rootCmd.AddCommand(identityCmd)
}