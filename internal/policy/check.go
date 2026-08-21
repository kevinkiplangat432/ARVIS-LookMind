package policy

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"strings"

	"github.com/redis/go-redis/v9"
)

type Violation struct {
	Rule     string
	Detail   string
	Severity string
}

// CheckBudgets checks identity, then provider, then global — in that
// order, since a more specific budget being exceeded is usually the
// more useful thing to report first. Returns the first violation
// found; enforcement only needs one reason to say no.
func CheckBudgets(ctx context.Context, rdb *redis.Client, identityID, providerName string) (*Violation, error) {
	scopes := []struct {
		scope BudgetScope
		id    string
	}{
		{ScopeIdentity, identityID},
		{ScopeProvider, providerName},
		{ScopeGlobal, ""},
	}

	for _, s := range scopes {
		budget, err := GetBudget(ctx, rdb, s.scope, s.id)
		if err != nil {
			return nil, err
		}
		if budget == nil || !budget.Enabled {
			continue
		}

		if budget.DailyLimit > 0 {
			daily, err := getUsage(ctx, rdb, s.scope, s.id, "daily")
			if err != nil {
				return nil, err
			}
			if daily >= budget.DailyLimit {
				return &Violation{
					Rule: "budget_exceeded_daily", Severity: "high",
					Detail: fmt.Sprintf("%s budget exceeded: %d/%d tokens used today", s.scope, daily, budget.DailyLimit),
				}, nil
			}
		}

		if budget.MonthlyLimit > 0 {
			monthly, err := getUsage(ctx, rdb, s.scope, s.id, "monthly")
			if err != nil {
				return nil, err
			}
			if monthly >= budget.MonthlyLimit {
				return &Violation{
					Rule: "budget_exceeded_monthly", Severity: "high",
					Detail: fmt.Sprintf("%s budget exceeded: %d/%d tokens used this month", s.scope, monthly, budget.MonthlyLimit),
				}, nil
			}
		}
	}

	return nil, nil
}

func usageKey(scope BudgetScope, id, period string) string {
	now := time.Now().UTC()
	suffix := now.Format("2006-01-02")
	if period == "monthly" {
		suffix = now.Format("2006-01")
	}
	base := "policy:usage:" + string(scope)
	if id != "" {
		base += ":" + id
	}
	return base + ":" + period + ":" + suffix
}

func getUsage(ctx context.Context, rdb *redis.Client, scope BudgetScope, id, period string) (int, error) {
	val, err := rdb.Get(ctx, usageKey(scope, id, period)).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	n, _ := strconv.Atoi(val)
	return n, nil
}

// RecordUsage always writes, regardless of whether a budget is
// currently enabled for this scope — so enabling a budget later
// starts from real history instead of zero. TTLs (48h, 32d) exist
// only to keep Redis from accumulating stale keys forever; they're
// deliberately longer than the period itself so a counter never
// expires mid-period.
func RecordUsage(ctx context.Context, rdb *redis.Client, scope BudgetScope, id string, tokens int) error {
	dailyKey := usageKey(scope, id, "daily")
	monthlyKey := usageKey(scope, id, "monthly")

	pipe := rdb.TxPipeline()
	pipe.IncrBy(ctx, dailyKey, int64(tokens))
	pipe.Expire(ctx, dailyKey, 48*time.Hour)
	pipe.IncrBy(ctx, monthlyKey, int64(tokens))
	pipe.Expire(ctx, monthlyKey, 32*24*time.Hour)
	_, err := pipe.Exec(ctx)
	return err
}

// RecordUsageAll writes to identity, provider, and global counters in
// one call — every completed request counts against all three scopes
// at once, since which ones are actually enforced can change anytime
// via the dashboard toggle without needing historical data backfilled.
func RecordUsageAll(ctx context.Context, rdb *redis.Client, identityID, providerName string, tokens int) error {
	if err := RecordUsage(ctx, rdb, ScopeIdentity, identityID, tokens); err != nil {
		return err
	}
	if err := RecordUsage(ctx, rdb, ScopeProvider, providerName, tokens); err != nil {
		return err
	}
	return RecordUsage(ctx, rdb, ScopeGlobal, "", tokens)
}

// CheckTopics scans the raw request body against every currently
// blocked topic's keywords. Whole-body, case-insensitive substring
// matching — the same pragmatic approach as detector.ContentRule, not
// a semantic classifier. Good enough to prove the mechanism; a real
// classifier is a drop-in upgrade later since it'd satisfy this same
// function's contract.
func CheckTopics(ctx context.Context, rdb *redis.Client, requestBody []byte) (*Violation, error) {
	blocked, err := ListBlockedTopics(ctx, rdb)
	if err != nil {
		return nil, err
	}
	if len(blocked) == 0 {
		return nil, nil
	}

	body := strings.ToLower(string(requestBody))
	for _, key := range blocked {
		topic, ok := GetTopic(key)
		if !ok {
			continue
		}
		for _, kw := range topic.Keywords {
			if strings.Contains(body, strings.ToLower(kw)) {
				return &Violation{
					Rule: "blocked_topic_" + topic.Key, Severity: "high",
					Detail: fmt.Sprintf("request matches blocked topic %q (%s)", topic.Name, topic.Source),
				}, nil
			}
		}
	}
	return nil, nil
}

// Check is the single entry point the proxy calls — budgets first
// (cheaper, no body scan needed), topics second.
func Check(ctx context.Context, rdb *redis.Client, identityID, providerName string, requestBody []byte) (*Violation, error) {
	if v, err := CheckBudgets(ctx, rdb, identityID, providerName); err != nil || v != nil {
		return v, err
	}
	return CheckTopics(ctx, rdb, requestBody)
}