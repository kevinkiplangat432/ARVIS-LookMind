package commands

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kevinkiplangat432/arvis/internal/config"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate [up|down]",
	Short: "Run database migrations",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		direction := "up"
		if len(args) > 0 {
			direction = args[0]
		}
		return runMigrate(direction)
	},
}

func runMigrate(direction string) error {
	cfg := config.Load()

	m, err := migrate.New("file://migrations", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to init migrations: %w", err)
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migrate up failed: %w", err)
		}
		fmt.Println("migrations applied")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migrate down failed: %w", err)
		}
		fmt.Println("migrations rolled back")
	default:
		return fmt.Errorf("unknown direction %q — use up or down", direction)
	}

	return nil
}