package audit

import (
	"fmt"
	"sort"

	"github.com/go-pdf/fpdf"
)

// RenderPDF is deliberately plain — a compliance document earns trust
// by being clear and complete, not by being visually elaborate. Every
// section here maps directly to something an auditor would actually
// ask for: who, when, what happened, and the specific incidents.
func RenderPDF(r *Report, outPath string) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()

	pdf.SetFont("Helvetica", "B", 18)
	pdf.CellFormat(0, 10, "ARVIS Audit Report", "", 1, "L", false, 0, "")

	pdf.SetFont("Helvetica", "", 10)
	scope := "Organisation-wide"
	if r.Identity != nil {
		scope = fmt.Sprintf("Identity: %s (%s)", r.Identity.Name, r.Identity.ID)
	}
	pdf.CellFormat(0, 6, scope, "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Period: %s to %s", r.From.Format("2006-01-02"), r.To.Format("2006-01-02")), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 6, "Generated: "+r.GeneratedAt.Format("2006-01-02 15:04 MST"), "", 1, "L", false, 0, "")
	pdf.Ln(6)

	pdf.SetFont("Helvetica", "B", 13)
	pdf.CellFormat(0, 8, "Summary", "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	summaryLines := []string{
		fmt.Sprintf("Total requests: %d", r.Summary.TotalRequests),
		fmt.Sprintf("Total anomalies: %d", r.Summary.TotalAnomalies),
		fmt.Sprintf("Blocked by policy (pre-flight): %d", r.Summary.BlockedCount),
		fmt.Sprintf("Terminated by kill switch (mid-stream): %d", r.Summary.KilledCount),
		fmt.Sprintf("Total prompt tokens: %d", r.Summary.TotalPromptTok),
		fmt.Sprintf("Total completion tokens: %d", r.Summary.TotalCompleteTok),
	}
	for _, line := range summaryLines {
		pdf.CellFormat(0, 6, line, "", 1, "L", false, 0, "")
	}
	pdf.Ln(4)

	if len(r.Summary.ByCategory) > 0 {
		pdf.SetFont("Helvetica", "B", 11)
		pdf.CellFormat(0, 7, "Anomalies by category", "", 1, "L", false, 0, "")
		pdf.SetFont("Helvetica", "", 10)
		categories := make([]string, 0, len(r.Summary.ByCategory))
		for cat := range r.Summary.ByCategory {
			categories = append(categories, cat)
		}
		sort.Strings(categories)
		for _, cat := range categories {
			pdf.CellFormat(0, 6, fmt.Sprintf("  %s: %d", cat, r.Summary.ByCategory[cat]), "", 1, "L", false, 0, "")
		}
		pdf.Ln(4)
	}

	pdf.SetFont("Helvetica", "B", 13)
	pdf.CellFormat(0, 8, "Anomaly Log", "", 1, "L", false, 0, "")

	if len(r.Anomalies) == 0 {
		pdf.SetFont("Helvetica", "I", 10)
		pdf.CellFormat(0, 6, "No anomalies recorded in this period.", "", 1, "L", false, 0, "")
	}

	pdf.SetFont("Helvetica", "B", 8)
	widths := []float64{28, 30, 22, 20, 90}
	headers := []string{"Timestamp", "Rule", "Category", "Severity", "Detail"}
	for i, h := range headers {
		pdf.CellFormat(widths[i], 6, h, "1", 0, "L", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Helvetica", "", 7)
	for _, a := range r.Anomalies {
		pdf.CellFormat(widths[0], 6, a.CreatedAt.Format("01-02 15:04"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[1], 6, a.Rule, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 6, a.Category, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 6, a.Severity, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[4], 6, truncate(a.Detail, 90), "1", 0, "L", false, 0, "")
		pdf.Ln(-1)
	}

	return pdf.OutputFileAndClose(outPath)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}