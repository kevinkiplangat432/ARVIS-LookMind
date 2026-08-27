package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/kevinkiplangat432/arvis/internal/store"
)

var (
	reqIdentity string
	reqLimit    int
)

var requestsCmd = &cobra.Command{
	Use:   "requests",
	Short: "Inspect logged requests",
}

var requestsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := store.Connect(cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer db.Close()

		results, err := store.ListRequests(context.Background(), db, reqLimit)
		if err != nil {
			return fmt.Errorf("failed to list requests: %w", err)
		}

		for _, r := range results {
			identity := "unattributed"
			if r.IdentityID != nil {
				identity = *r.IdentityID
			}
			if reqIdentity != "" && identity != reqIdentity {
				continue
			}
			fmt.Printf("%-20s %-10s %-20s %-8s %5dms  status=%d  tokens=%d/%d\n",
				r.CreatedAt.Format("2006-01-02 15:04:05"), r.Provider, r.Model, identity, r.LatencyMs, r.StatusCode, r.PromptTokens, r.CompTokens)
		}
		return nil
	},
}

func init() {
	requestsListCmd.Flags().StringVar(&reqIdentity, "identity", "", "filter to one identity ID")
	requestsListCmd.Flags().IntVar(&reqLimit, "limit", 50, "max results")

	requestsCmd.AddCommand(requestsListCmd)
	rootCmd.AddCommand(requestsCmd)
}