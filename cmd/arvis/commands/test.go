package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		testDBURL = "postgres://arvis:arvis@localhost:5432/arvis_test?sslmode=disable"
	}

	fmt.Println("\nRunning ARVIS test suite...")
	fmt.Println("Test database:", testDBURL)
	fmt.Println(strings.Repeat("-", 50))

	cmd := exec.Command("go", "test", "./...", "-v")
	cmd.Env = append(os.Environ(), "TEST_DATABASE_URL="+testDBURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Println("All tests passed.")
	return nil
}