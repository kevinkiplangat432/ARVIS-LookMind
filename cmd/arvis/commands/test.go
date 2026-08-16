package commands

import (
	"github.com/spf13/cobra"
)
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run the ARVIS test suite",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTests()
	},
}

func runTests() error {
	return nil
}
