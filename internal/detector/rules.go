package detector

import (
    "fmt"

    "github.com/kevinkiplangat432/arvis/internal/store"
)

type Flag struct {
    Rule   string
    Detail string
}

func Check(r store.Request) []Flag {
    var flags []Flag

    if r.LatencyMs > 10000 {
        flags = append(flags, Flag{Rule: "high_latency", Detail: "latency exceeded 10s"})
    }
    if r.PromptTokens+r.CompTokens > 3000 {
        flags = append(flags, Flag{Rule: "high_token_usage", Detail: "total tokens exceeded 3000"})
    }
    if r.StatusCode >= 500 {
        flags = append(flags, Flag{Rule: "upstream_error", Detail: fmt.Sprintf("upstream returned %d", r.StatusCode)})
    }
    if r.StatusCode >= 400 && r.StatusCode < 500 {
        flags = append(flags, Flag{Rule: "client_error", Detail: fmt.Sprintf("upstream returned %d", r.StatusCode)})
    }

    return flags
}