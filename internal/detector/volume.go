package detector

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// VolumeRule flags an identity making more than Threshold requests
// within Window. Purely in-memory — a restart resets counts, an
// acceptable trade-off since this rule's whole point is speed, not
// perfect historical accuracy (the requests table already has that).
type VolumeRule struct {
	Threshold int
	Window    time.Duration

	mu     sync.Mutex
	counts map[string][]time.Time
}

func NewVolumeRule(threshold int, window time.Duration) *VolumeRule {
	return &VolumeRule{
		Threshold: threshold,
		Window:    window,
		counts:    make(map[string][]time.Time),
	}
}

func (v *VolumeRule) CheckSync(ctx context.Context, identityID, model string) []Flag {
	v.mu.Lock()
	defer v.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-v.Window)

	recent := v.counts[identityID][:0]
	for _, t := range v.counts[identityID] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}
	recent = append(recent, now)
	v.counts[identityID] = recent

	if len(recent) > v.Threshold {
		return []Flag{{
			Rule:     "volume_spike",
			Category: "volume",
			Severity: "medium",
			Detail:   fmt.Sprintf("identity made %d requests in the last %s (threshold %d)", len(recent), v.Window, v.Threshold),
		}}
	}
	return nil
}