package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ghoststack/ghoststack/internal/audit"
	"github.com/spf13/cobra"
)

func newAuditCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "Run privacy and security audits",
	}

	cmd.AddCommand(newAuditPrivacyCommand())
	cmd.AddCommand(newAuditSecurityCommand())
	cmd.AddCommand(newAuditAllCommand())

	return cmd
}

func newAuditPrivacyCommand() *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "privacy",
		Short: "Run privacy audit checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := audit.RunPrivacyAudit(cmd.Context())
			if err != nil {
				return fmt.Errorf("privacy audit: %w", err)
			}

			if jsonOut {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(report)
			}

			printAuditReport("Privacy Audit", report)
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")
	return cmd
}

func newAuditSecurityCommand() *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "security",
		Short: "Run security audit checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := audit.RunSecurityAudit(cmd.Context())
			if err != nil {
				return fmt.Errorf("security audit: %w", err)
			}

			if jsonOut {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(report)
			}

			printAuditReport("Security Audit", report)
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")
	return cmd
}

func newAuditAllCommand() *cobra.Command {
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "all",
		Short: "Run all audits (privacy + security)",
		RunE: func(cmd *cobra.Command, args []string) error {
			privReport, err := audit.RunPrivacyAudit(cmd.Context())
			if err != nil {
				return fmt.Errorf("privacy audit: %w", err)
			}

			secReport, err := audit.RunSecurityAudit(cmd.Context())
			if err != nil {
				return fmt.Errorf("security audit: %w", err)
			}

			if jsonOut {
				out := map[string]any{
					"privacy":  privReport,
					"security": secReport,
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(out)
			}

			printAuditReport("Privacy Audit", privReport)
			fmt.Fprintln(os.Stdout)
			printAuditReport("Security Audit", secReport)

			totalScore := (privReport.Score + secReport.Score) / 2
			fmt.Fprintln(os.Stdout)
			fmt.Fprintf(os.Stdout, "Overall Score: %d/100\n", totalScore)
			if totalScore >= 80 {
				fmt.Fprintln(os.Stdout, "Status: GOOD")
			} else if totalScore >= 50 {
				fmt.Fprintln(os.Stdout, "Status: FAIR")
			} else {
				fmt.Fprintln(os.Stdout, "Status: POOR")
			}

			if privReport.Failed > 0 || secReport.Failed > 0 {
				return fmt.Errorf("audit completed with %d failures", privReport.Failed+secReport.Failed)
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")
	return cmd
}

func printAuditReport(title string, report *audit.AuditReport) {
	fmt.Fprintf(os.Stdout, "=== %s ===\n", title)
	fmt.Fprintf(os.Stdout, "Score: %d/100 (%d pass, %d fail)\n", report.Score, report.Passed, report.Failed)
	fmt.Fprintln(os.Stdout)

	for _, r := range report.Results {
		icon := "✓"
		switch r.Severity {
		case audit.SeverityPass:
			icon = "✓"
		case audit.SeverityInfo:
			icon = "ℹ"
		case audit.SeverityWarn:
			icon = "⚠"
		case audit.SeverityFail:
			icon = "✗"
		}
		fmt.Fprintf(os.Stdout, "  %s [%s] %s\n", icon, r.Severity, r.Name)
		fmt.Fprintf(os.Stdout, "         %s\n", r.Message)
		if r.Fix != "" {
			fmt.Fprintf(os.Stdout, "         Fix: %s\n", r.Fix)
		}
		fmt.Fprintln(os.Stdout)
	}
}
