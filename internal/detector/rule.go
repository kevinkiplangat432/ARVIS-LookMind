package detector

import "context"

// Flag is what a Rule emits when it thinks something is off. Category
// and Severity map directly onto the anomalies table's columns.
type Flag struct {
	Rule     string
	Category string
	Severity string
	Detail   string
}

// SyncRule runs before a request is forwarded and must be fast — it
// blocks the caller's request. Only cheap, in-memory checks belong
// here.
type SyncRule interface {
	CheckSync(ctx context.Context, identityID, model string) []Flag
}

// AsyncRule runs after the response has already gone back to the
// caller, so it's safe to be slower. Never blocks a client.
type AsyncRule interface {
	CheckAsync(ctx context.Context, in AsyncInput) []Flag
}

// AsyncInput is everything an async rule might need. Kept as one
// struct rather than a long parameter list, since rules will likely
// need more fields over time without every implementation's
// signature having to change.
type AsyncInput struct {
	IdentityID   string
	Provider     string
	Model        string
	RequestBody  []byte
	ResponseBody []byte
	LatencyMs    int
	StatusCode   int
}