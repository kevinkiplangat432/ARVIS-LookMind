// identity_fuzz_test.go fuzzes validateIdentityName.
//
// Properties: an error happens exactly when the input is blank, accepted
// names have no surrounding whitespace, and validation is idempotent.
//
// Run seeds only: go test ./cmd/arvis/commands/...
// Run fuzzing:    go test ./cmd/arvis/commands -run='^$' -fuzz='^FuzzValidateIdentityName$' -fuzztime=30s
// Needs: nothing, no database or network.

package commands

import (
	"strings"
	"testing"
)

func FuzzValidateIdentityName(f *testing.F) {
	for _, seed := range []string{"", " ", "jane", "  jane  ", "\t\n", "a b"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, in string) {
		got, err := validateIdentityName(in)
		blank := strings.TrimSpace(in) == ""
		if (err != nil) != blank {
			t.Fatalf("err = %v for blank = %v, input %q", err, blank, in)
		}
		if err != nil {
			return
		}
		if got != strings.TrimSpace(got) {
			t.Fatalf("result %q has surrounding whitespace", got)
		}
		again, err := validateIdentityName(got)
		if err != nil || again != got {
			t.Fatalf("not idempotent: %q -> %q (%v)", got, again, err)
		}
	})
}