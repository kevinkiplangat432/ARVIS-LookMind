package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/kevinkiplangat432/arvis/internal/store"
)

var (
	anomCategory string
	anomSeverity string
	anomStatus   string
	anomLimit    int
)

var anomaliesCmd = &cobra.Command{
	Use:   "anomalies",
	Short: "Inspect and manage flagged anomalies",
}

var anomaliesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List anomalies, optionally filtered",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := store.Connect(cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer db.Close()

		filter := store.AnomalyFilter{Category: anomCategory, Severity: anomSeverity, Status: anomStatus}
		results, err := store.ListAnomaliesFiltered(context.Background(), db, filter, anomLimit)
		if err != nil {
			return fmt.Errorf("failed to list anomalies: %w", err)
		}

		for _, a := range results {
			fmt.Printf("%-36s %-20s %-10s %-8s %-8s %s\n", a.ID, a.CreatedAt.Format("2006-01-02 15:04:05"), a.Category, a.Severity, a.Status, a.Detail)
		}
		return nil
	},
}

var anomaliesResolveCmd = &cobra.Command{
	Use:   "resolve [anomaly-id]",
	Short: "Mark an anomaly as reviewed or dismissed",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if anomStatus != "reviewed" && anomStatus != "dismissed" {
			return fmt.Errorf("--status must be 'reviewed' or 'dismissed'")
		}

		db, err := store.Connect(cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer db.Close()

		if err := store.UpdateAnomalyStatus(context.Background(), db, args[0], anomStatus); err != nil {
			return fmt.Errorf("failed to update anomaly: %w", err)
		}
		fmt.Printf("Anomaly %s marked as %s\n", args[0], anomStatus)
		return nil
	},
}

func init() {
	anomaliesListCmd.Flags().StringVar(&anomCategory, "category", "", "filter by category (volume, content, latency, governance)")
	anomaliesListCmd.Flags().StringVar(&anomSeverity, "severity", "", "filter by severity")
	anomaliesListCmd.Flags().StringVar(&anomStatus, "status", "", "filter by status (open, reviewed, dismissed)")
	anomaliesListCmd.Flags().IntVar(&anomLimit, "limit", 50, "max results")

	anomaliesResolveCmd.Flags().StringVar(&anomStatus, "status", "reviewed", "reviewed or dismissed")

	anomaliesCmd.AddCommand(anomaliesListCmd, anomaliesResolveCmd)
	rootCmd.AddCommand(anomaliesCmd)
}