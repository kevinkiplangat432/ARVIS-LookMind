package commands

import (
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
	return nil
}