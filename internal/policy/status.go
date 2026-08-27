package policy

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/redis/go-redis/v9"
)

type BudgetStatus struct {
	Scope       BudgetScope
	ID          string
	Budget      Budget
	DailyUsed   int
	MonthlyUsed int
}

// ListBudgetStatuses scans every configured budget key and pairs each
// one with its current usage — this is entirely a read path built for
// display, "arvis policy budgets" is its only caller. Uses Scan, not
// Keys — Keys blocks the whole Redis instance on a large keyspace,
// Scan doesn't, and there's no reason a CLI status command should
// carry that risk into a production Redis.
func ListBudgetStatuses(ctx context.Context, rdb *redis.Client) ([]BudgetStatus, error) {
	var out []BudgetStatus
	iter := rdb.Scan(ctx, 0, "policy:budget:*", 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		data, err := rdb.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}
		var b Budget
		if err := json.Unmarshal(data, &b); err != nil {
			continue
		}

		scope, id := parseBudgetKey(key)
		daily, monthly, err := GetUsage(ctx, rdb, scope, id)
		if err != nil {
			continue
		}

		out = append(out, BudgetStatus{Scope: scope, ID: id, Budget: b, DailyUsed: daily, MonthlyUsed: monthly})
	}

	return out, iter.Err()
}

func parseBudgetKey(key string) (scope BudgetScope, id string) {
	parts := strings.SplitN(key, ":", 4)
	if len(parts) < 3 {
		return "", ""
	}
	if parts[2] == "global" {
		return ScopeGlobal, ""
	}
	scope = BudgetScope(parts[2])
	if len(parts) == 4 {
		id = parts[3]
	}
	return scope, id
}