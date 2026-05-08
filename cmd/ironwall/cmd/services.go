package cmd

import (
	"fmt"

	"github.com/apringutawa/ironwall/internal/hardening"
	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Manage system services",
	Long:  `Minimize attack surface by disabling dangerous and unnecessary services.`,
}

var servicesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List enabled services with security classification",
	Run: func(cmd *cobra.Command, args []string) {
		svc := hardening.NewServiceManager()

		services, err := svc.ListEnabled()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println("🔧 Service Status")
		fmt.Println("==================")

		var dangerous, essential, review []hardening.ServiceInfo
		for _, s := range services {
			switch s.Action {
			case "disable":
				dangerous = append(dangerous, s)
			case "keep":
				essential = append(essential, s)
			default:
				review = append(review, s)
			}
		}

		if len(dangerous) > 0 {
			fmt.Println("\n🛑 Dangerous Services (recommended to disable):")
			for _, s := range dangerous {
				fmt.Printf("  ⚠ %s (%s)\n", s.Name, s.Status)
			}
		}

		if len(essential) > 0 {
			fmt.Println("\n✅ Essential Services:")
			for _, s := range essential {
				fmt.Printf("  ✓ %s (%s)\n", s.Name, s.Status)
			}
		}

		if len(review) > 0 {
			fmt.Println("\n📋 Optional Services (review if needed):")
			for _, s := range review {
				fmt.Printf("  • %s (%s)\n", s.Name, s.Status)
			}
		}
	},
}

var servicesMinimizeCmd = &cobra.Command{
	Use:   "minimize",
	Short: "Disable dangerous services",
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		svc := hardening.NewServiceManager()

		if dryRun {
			fmt.Println("📋 Service Minimization (Dry Run)")
			fmt.Println("==================================")
			fmt.Println("The following dangerous services would be disabled:")
			for _, s := range svc.DangerousServices {
				fmt.Printf("  ⚠ %s\n", s)
			}
			fmt.Println("\nRun without --dry-run to apply.")
			return
		}

		if _, err := svc.Apply(false); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("\n✅ Service minimization completed")
	},
}

var servicesScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for insecure running services",
	Run: func(cmd *cobra.Command, args []string) {
		svc := hardening.NewServiceManager()

		fmt.Println("🔍 Scanning for insecure services...")
		fmt.Println()

		insecure, err := svc.ScanForInsecureServices()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(insecure) > 0 {
			fmt.Println("⚠️  Insecure services detected:")
			for _, s := range insecure {
				fmt.Printf("  %s\n", s)
			}
			fmt.Println("\nConsider running 'ironwall services minimize' to disable them.")
		} else {
			fmt.Println("✅ No insecure services detected")
		}
	},
}

func init() {
	rootCmd.AddCommand(servicesCmd)
	servicesCmd.AddCommand(servicesListCmd)
	servicesCmd.AddCommand(servicesMinimizeCmd)
	servicesCmd.AddCommand(servicesScanCmd)

	servicesMinimizeCmd.Flags().Bool("dry-run", false, "Preview changes without applying")
}
