package policy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type BudgetScope string

const (
	ScopeIdentity BudgetScope = "identity"
	ScopeProvider BudgetScope = "provider"
	ScopeGlobal   BudgetScope = "global"
)

type Budget struct {
	DailyLimit   int  `json:"daily_limit"`
	MonthlyLimit int  `json:"monthly_limit"`
	Enabled      bool `json:"enabled"`
}

func budgetKey(scope BudgetScope, id string) string {
	if scope == ScopeGlobal {
		return "policy:budget:global"
	}
	return fmt.Sprintf("policy:budget:%s:%s", scope, id)
}

func SetBudget(ctx context.Context, rdb *redis.Client, scope BudgetScope, id string, b Budget) error {
	data, err := json.Marshal(b)
	if err != nil {
		return fmt.Errorf("failed to marshal budget: %w", err)
	}
	return rdb.Set(ctx, budgetKey(scope, id), data, 0).Err()
}

// GetBudget returns (nil, nil) if no budget is configured for this
// scope/id — "no budget set" and "budget set but disabled" are
// deliberately different states, callers need to tell them apart.
func GetBudget(ctx context.Context, rdb *redis.Client, scope BudgetScope, id string) (*Budget, error) {
	data, err := rdb.Get(ctx, budgetKey(scope, id)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get budget: %w", err)
	}
	var b Budget
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("failed to unmarshal budget: %w", err)
	}
	return &b, nil
}

const blockedTopicsKey = "policy:blocked_topics"

func BlockTopic(ctx context.Context, rdb *redis.Client, topicKey string) error {
	return rdb.SAdd(ctx, blockedTopicsKey, topicKey).Err()
}

func UnblockTopic(ctx context.Context, rdb *redis.Client, topicKey string) error {
	return rdb.SRem(ctx, blockedTopicsKey, topicKey).Err()
}

func ListBlockedTopics(ctx context.Context, rdb *redis.Client) ([]string, error) {
	return rdb.SMembers(ctx, blockedTopicsKey).Result()
}