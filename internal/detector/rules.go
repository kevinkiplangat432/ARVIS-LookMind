package detector

import "github.com/kevinkiplangat432/arvis/internal/store"

type Flag struct {
	Rule   string
	Detail string
}

// Check runs all rules against a request and returns any flags.
func Check(r store.Request) []Flag {
	var flags []Flag

	if r.LatencyMs > 10000 {
		flags = append(flags, Flag{Rule: "high_latency", Detail: "latency exceeded 10s"})
	}
	if r.PromptTokens+r.CompTokens > 3000 {
		flags = append(flags, Flag{Rule: "high_token_usage", Detail: "total tokens exceeded 3000"})
	}
	if r.StatusCode >= 500 {
		flags = append(flags, Flag{Rule: "upstream_error", Detail: "upstream returned 5xx"})
	}

	return flags
}
