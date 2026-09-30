package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/kevinkiplangat432/arvis/internal/audit"
	"github.com/kevinkiplangat432/arvis/internal/store"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Generate audit reports",
}

var (
	auditIdentityID string
	auditFrom       string
	auditTo         string
	auditOut        string
)

var auditExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export an audit report as a PDF",
	RunE: func(cmd *cobra.Command, args []string) error {
		to := time.Now()
		if auditTo != "" {
			t, err := time.Parse("2006-01-02", auditTo)
			if err != nil {
				return fmt.Errorf("invalid --to date, expected YYYY-MM-DD: %w", err)
			}
			to = t
		}

		from := to.AddDate(0, 0, -30) // default: trailing 30 days
		if auditFrom != "" {
			f, err := time.Parse("2006-01-02", auditFrom)
			if err != nil {
				return fmt.Errorf("invalid --from date, expected YYYY-MM-DD: %w", err)
			}
			from = f
		}

		var identityID *string
		if auditIdentityID != "" {
			identityID = &auditIdentityID
		}

		outPath := auditOut
		if outPath == "" {
			outPath = fmt.Sprintf("arvis-audit-%s.pdf", time.Now().Format("20060102-150405"))
		}

		db, err := store.Connect(cfg.DatabaseURL)
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		defer db.Close()

		report, err := audit.BuildReport(context.Background(), db, identityID, from, to)
		if err != nil {
			return fmt.Errorf("failed to build report: %w", err)
		}

		if err := audit.RenderPDF(report, outPath); err != nil {
			return fmt.Errorf("failed to render PDF: %w", err)
		}

		fmt.Printf("Audit report written to %s (%d requests, %d anomalies)\n", outPath, report.Summary.TotalRequests, report.Summary.TotalAnomalies)
		return nil
	},
}

func init() {
	auditExportCmd.Flags().StringVar(&auditIdentityID, "identity", "", "filter to one identity (omit for org-wide)")
	auditExportCmd.Flags().StringVar(&auditFrom, "from", "", "start date, YYYY-MM-DD (default: 30 days ago)")
	auditExportCmd.Flags().StringVar(&auditTo, "to", "", "end date, YYYY-MM-DD (default: today)")
	auditExportCmd.Flags().StringVar(&auditOut, "out", "", "output file path (default: arvis-audit-<timestamp>.pdf)")

	auditCmd.AddCommand(auditExportCmd)
	rootCmd.AddCommand(auditCmd)
}