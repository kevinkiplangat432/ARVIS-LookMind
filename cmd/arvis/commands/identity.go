package commands

import (
	"context"
	"fmt"
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

func createIdentity(name string) error {
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

	fmt.Println("Identity created.")
	fmt.Println("  ID:  ", identity.ID)
	fmt.Println("  Name:", identity.Name)
	fmt.Println()
	fmt.Println("Key (shown once — store it securely, ARVIS cannot recover it):")
	fmt.Println("  ", rawKey)

	return nil
}

func init() {
	identityCmd.AddCommand(identityCreateCmd)
	rootCmd.AddCommand(identityCmd)
}