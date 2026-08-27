package commands

import (
	"context"
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/kevinkiplangat432/arvis/internal/policy"
)

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Manage governance policies — budgets and blocked topics",
}

var (
	budgetScope   string
	budgetID      string
	budgetDaily   int
	budgetMonthly int
	budgetEnabled bool
)

var policySetBudgetCmd = &cobra.Command{
	Use:   "set-budget",
	Short: "Set a token budget for an identity, a provider, or the whole org",
	RunE: func(cmd *cobra.Command, args []string) error {
		scope := policy.BudgetScope(budgetScope)
		if scope != policy.ScopeIdentity && scope != policy.ScopeProvider && scope != policy.ScopeGlobal {
			return fmt.Errorf("--scope must be identity, provider, or global")
		}
		if scope != policy.ScopeGlobal && budgetID == "" {
			return fmt.Errorf("--id is required for scope %q", scope)
		}

		rdb, err := policy.Connect(cfg.RedisAddr)
		if err != nil {
			return err
		}
		defer rdb.Close()

		b := policy.Budget{DailyLimit: budgetDaily, MonthlyLimit: budgetMonthly, Enabled: budgetEnabled}
		if err := policy.SetBudget(context.Background(), rdb, scope, budgetID, b); err != nil {
			return fmt.Errorf("failed to set budget: %w", err)
		}

		fmt.Printf("Budget set — scope: %s, id: %s, daily: %d, monthly: %d, enabled: %v\n", scope, budgetID, budgetDaily, budgetMonthly, budgetEnabled)
		return nil
	},
}

var policyBlockTopicCmd = &cobra.Command{
	Use:   "block-topic [key]",
	Short: "Block a topic category org-wide (see 'policy list-topics' for available keys)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, ok := policy.GetTopic(args[0]); !ok {
			return fmt.Errorf("unknown topic key %q — run 'arvis policy list-topics' to see valid keys", args[0])
		}
		rdb, err := policy.Connect(cfg.RedisAddr)
		if err != nil {
			return err
		}
		defer rdb.Close()

		if err := policy.BlockTopic(context.Background(), rdb, args[0]); err != nil {
			return fmt.Errorf("failed to block topic: %w", err)
		}
		fmt.Printf("Blocked: %s\n", args[0])
		return nil
	},
}

var policyUnblockTopicCmd = &cobra.Command{
	Use:   "unblock-topic [key]",
	Short: "Unblock a previously blocked topic category",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rdb, err := policy.Connect(cfg.RedisAddr)
		if err != nil {
			return err
		}
		defer rdb.Close()

		if err := policy.UnblockTopic(context.Background(), rdb, args[0]); err != nil {
			return fmt.Errorf("failed to unblock topic: %w", err)
		}
		fmt.Printf("Unblocked: %s\n", args[0])
		return nil
	},
}

var policyListTopicsCmd = &cobra.Command{
	Use:   "list-topics",
	Short: "List all topic categories available to block",
	RunE: func(cmd *cobra.Command, args []string) error {
		topics := policy.ListTopics()
		sort.Slice(topics, func(i, j int) bool { return topics[i].Key < topics[j].Key })
		for _, t := range topics {
			fmt.Printf("%-32s %-45s [%s]\n", t.Key, t.Name, t.Source)
		}
		return nil
	},
}

var policyShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show currently blocked topics",
	RunE: func(cmd *cobra.Command, args []string) error {
		rdb, err := policy.Connect(cfg.RedisAddr)
		if err != nil {
			return err
		}
		defer rdb.Close()

		blocked, err := policy.ListBlockedTopics(context.Background(), rdb)
		if err != nil {
			return fmt.Errorf("failed to list blocked topics: %w", err)
		}
		if len(blocked) == 0 {
			fmt.Println("No topics currently blocked.")
			return nil
		}
		fmt.Println("Currently blocked:")
		for _, key := range blocked {
			fmt.Println(" -", key)
		}
		return nil
	},
}
var policyBudgetsCmd = &cobra.Command{
	Use:   "budgets",
	Short: "Show all configured budgets and current usage against them",
	RunE: func(cmd *cobra.Command, args []string) error {
		rdb, err := policy.Connect(cfg.RedisAddr)
		if err != nil {
			return err
		}
		defer rdb.Close()

		statuses, err := policy.ListBudgetStatuses(context.Background(), rdb)
		if err != nil {
			return fmt.Errorf("failed to list budget statuses: %w", err)
		}
		if len(statuses) == 0 {
			fmt.Println("No budgets configured.")
			return nil
		}

		for _, s := range statuses {
			label := string(s.Scope)
			if s.ID != "" {
				label += ":" + s.ID
			}
			enabled := "disabled"
			if s.Budget.Enabled {
				enabled = "enabled"
			}
			fmt.Printf("%-30s [%s]  daily: %d/%d   monthly: %d/%d\n",
				label, enabled, s.DailyUsed, s.Budget.DailyLimit, s.MonthlyUsed, s.Budget.MonthlyLimit)
		}
		return nil
	},
}
func init() {
	policySetBudgetCmd.Flags().StringVar(&budgetScope, "scope", "", "identity, provider, or global (required)")
	policySetBudgetCmd.Flags().StringVar(&budgetID, "id", "", "identity or provider name (required unless scope=global)")
	policySetBudgetCmd.Flags().IntVar(&budgetDaily, "daily", 0, "daily token limit")
	policySetBudgetCmd.Flags().IntVar(&budgetMonthly, "monthly", 0, "monthly token limit")
	policySetBudgetCmd.Flags().BoolVar(&budgetEnabled, "enabled", true, "whether this budget is enforced")

	
	policyCmd.AddCommand(policySetBudgetCmd, policyBlockTopicCmd, policyUnblockTopicCmd, policyListTopicsCmd, policyShowCmd, policyBudgetsCmd)
	rootCmd.AddCommand(policyCmd)

}