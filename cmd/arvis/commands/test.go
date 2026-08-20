package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run the ARVIS test suite (vet, race-enabled tests, coverage)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTests()
	},
}

// runTests runs what any CI pipeline should run before anything ships.
// Stops at the first failing step rather than running everything and
// summarizing at the end — a failed vet buries real test output in
// noise, so there's no reason to let it get that far.
func runTests() error {
	steps := []struct {
		name string
		args []string
	}{
		{"go vet", []string{"vet", "./..."}},
		// -race catches data races the detector's goroutines (Phase 6/7)
		// are exactly the kind of code that can hide. -count=1 disables
		// Go's test result caching — a "rigorous" suite means every run
		// actually re-executes against current code, never reports a
		// stale cached pass.
		{"go test (race + coverage)", []string{"test", "./...", "-race", "-cover", "-count=1"}},
	}

	for _, step := range steps {
		fmt.Printf("\n==> %s\n", step.name)
		cmd := exec.Command("go", step.args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s failed: %w", step.name, err)
		}
	}

	fmt.Println("\nAll checks passed.")
	return nil
}