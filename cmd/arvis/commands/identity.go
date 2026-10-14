// identity.go is the CLI front for managing ARVIS identities.
//
// Rules that must never be broken:
//  1. The raw key is printed once, to stdout, after the insert succeeds.
//     It is never logged, written to disk, or retrievable later.
//  2. If the insert fails, no key is shown. We never hand out a key that
//     does not exist.
//  3. Only the hash goes to the database.
//
// Dependencies: internal/auth, internal/store, cobra. Must never import proxy.
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
	"strings"
	"time"

	"github.com/google/uuid"
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
		return createIdentity(args[0])
	},
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

func createIdentity(rawName string) error {
	name, err := validateIdentityName(rawName)
	if err != nil {
		return err
	}

	db, err := store.Connect(cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

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

	if err := store.InsertIdentity(context.Background(), db, identity); err != nil {
		return fmt.Errorf("failed to create identity: %w", err)
	}

	// Print only after the insert succeeds (rule 2).
	fmt.Println("Identity created.")
	fmt.Println("  ID:  ", identity.ID)
	fmt.Println("  Name:", identity.Name)
	fmt.Println()
	fmt.Println("Key (shown once, store it securely, ARVIS cannot recover it):")
	fmt.Println("  ", rawKey)

	return nil
}

func init() {
	identityCmd.AddCommand(identityCreateCmd)
	rootCmd.AddCommand(identityCmd)
}