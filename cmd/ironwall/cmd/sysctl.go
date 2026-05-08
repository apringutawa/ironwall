package cmd

import (
	"fmt"

	"github.com/apringutawa/ironwall/internal/hardening"
	"github.com/spf13/cobra"
)

var sysctlCmd = &cobra.Command{
	Use:   "sysctl",
	Short: "Manage kernel hardening parameters",
	Long:  `Apply and verify kernel hardening sysctl parameters.`,
}

var sysctlApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply kernel hardening parameters",
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		sysctl := hardening.NewSysctlConfig()

		if dryRun {
			fmt.Println("📋 Kernel Hardening Parameters (Dry Run)")
			fmt.Println("=========================================")
			fmt.Println("The following parameters will be applied:")
			fmt.Println()
			for _, p := range sysctl.GetParams() {
				fmt.Printf("  %-50s = %-6s  (%s)\n", p.Key, p.Value, p.Desc)
			}
			fmt.Println("\nConfig file:", sysctl.ConfigPath)
			return
		}

		if err := sysctl.Apply(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("\n✅ Kernel hardening applied successfully")
	},
}

var sysctlVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify kernel parameters are applied",
	Run: func(cmd *cobra.Command, args []string) {
		sysctl := hardening.NewSysctlConfig()

		fmt.Println("Verifying kernel parameters...")
		fmt.Println()

		for _, p := range sysctl.GetParams() {
			current, err := sysctl.GetCurrent(p.Key)
			if err != nil {
				fmt.Printf("  ⚠ %-50s = ERROR: %v\n", p.Key, err)
				continue
			}
			if current == p.Value {
				fmt.Printf("  ✓ %-50s = %s\n", p.Key, current)
			} else {
				fmt.Printf("  ✗ %-50s = %s (expected: %s)\n", p.Key, current, p.Value)
			}
		}
	},
}

var sysctlRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore previous sysctl configuration",
	Run: func(cmd *cobra.Command, args []string) {
		sysctl := hardening.NewSysctlConfig()
		if err := sysctl.Restore(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("✅ Previous sysctl configuration restored")
	},
}

func init() {
	rootCmd.AddCommand(sysctlCmd)
	sysctlCmd.AddCommand(sysctlApplyCmd)
	sysctlCmd.AddCommand(sysctlVerifyCmd)
	sysctlCmd.AddCommand(sysctlRestoreCmd)

	sysctlApplyCmd.Flags().Bool("dry-run", false, "Preview changes without applying")
}
