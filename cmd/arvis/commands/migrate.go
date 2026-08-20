package commands

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var migrationsPath string

var migrateCmd = &cobra.Command{
	Use:   "migrate [up|down|status|force] [args]",
	Short: "Run database migrations",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMigrate(args[0], args[1:])
	},
}

func init() {
	migrateCmd.Flags().StringVar(&migrationsPath, "path", "migrations", "path to the migrations directory")
}

// newMigrator resolves migrationsPath to an absolute path before
// handing it to golang-migrate. The old version hardcoded a relative
// "file://migrations", which only worked if arvis happened to be run
// from the project root — a real bug for an on-prem binary that gets
// invoked from wherever an ops team's deploy scripts live.
func newMigrator() (*migrate.Migrate, error) {
	abs, err := filepath.Abs(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve migrations path %q: %w", migrationsPath, err)
	}
	m, err := migrate.New("file://"+abs, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to init migrations from %q: %w", abs, err)
	}
	return m, nil
}

func runMigrate(direction string, args []string) error {
	m, err := newMigrator()
	if err != nil {
		return err
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

	case "status":
		version, dirty, err := m.Version()
		if err != nil {
			return fmt.Errorf("failed to read migration status: %w", err)
		}
		fmt.Printf("current version: %d, dirty: %v\n", version, dirty)

	case "force":
		if len(args) != 1 {
			return fmt.Errorf("force requires exactly one argument: the version to force to")
		}
		version, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid version %q: %w", args[0], err)
		}
		if err := m.Force(version); err != nil {
			return fmt.Errorf("force failed: %w", err)
		}
		fmt.Printf("forced to version %d — dirty flag cleared, but verify your schema actually matches this version before running up/down again\n", version)

	default:
		return fmt.Errorf("unknown direction %q — use up, down, status, or force", direction)
	}

	return nil
}