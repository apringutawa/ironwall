package cmd

import (
	"fmt"

	"github.com/apringutawa/ironwall/internal/audit"
	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Run security audit with Lynis",
	Long:  `Run a comprehensive security audit using Lynis and get hardening score.`,
	Run: func(cmd *cobra.Command, args []string) {
		scanner := audit.NewLynisScanner()

		if !scanner.IsInstalled() {
			fmt.Println("Lynis not installed. Use 'ironwall audit install' to install it first.")
			return
		}

		result, err := scanner.RunAudit()
		if err != nil {
			fmt.Printf("Error running audit: %v\n", err)
			return
		}

		fmt.Println("\n🛡️  Lynis Audit Results")
		fmt.Println("=======================")
		fmt.Printf("  Hardening Index: %.1f/100\n", result.HardeningIndex)
		fmt.Printf("  Status: %s\n", result.Status)
		fmt.Printf("  Warnings: %d\n", result.Warnings)
		fmt.Printf("  Suggestions: %d\n", result.Suggestions)
		fmt.Printf("  Tests Performed: %d\n", result.TestsPerformed)

		if result.HardeningIndex < 60 {
			fmt.Println("\n⚠️  Hardening score is below 60. Consider running 'ironwall protect' to improve.")
		} else {
			fmt.Println("\n✅ Server hardening score is acceptable.")
		}
	},
}

var auditInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Lynis security scanner",
	Run: func(cmd *cobra.Command, args []string) {
		scanner := audit.NewLynisScanner()
		if err := scanner.EnsureInstalled(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("✅ Lynis installed successfully")
	},
}

var auditScoreCmd = &cobra.Command{
	Use:   "score",
	Short: "Show hardening score from last audit",
	Run: func(cmd *cobra.Command, args []string) {
		scanner := audit.NewLynisScanner()
		score := scanner.GetHardeningScore()
		fmt.Printf("Last hardening score: %.1f/100\n", score)

		if score >= 80 {
			fmt.Println("Status: Excellent")
		} else if score >= 60 {
			fmt.Println("Status: Good")
		} else if score >= 40 {
			fmt.Println("Status: Fair")
		} else {
			fmt.Println("Status: Poor")
		}
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
	auditCmd.AddCommand(auditInstallCmd)
	auditCmd.AddCommand(auditScoreCmd)
}
