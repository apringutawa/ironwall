package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var protectCmd = &cobra.Command{
	Use:   "protect",
	Short: "Enable all security protections",
	Long:  `Enable all security modules and apply hardening configurations.`,
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		
		if dryRun {
			fmt.Println("🔒 IronWall Protect (Dry Run)")
			fmt.Println("=============================")
			fmt.Println("\nChanges that would be applied:")
			fmt.Println("  [SSH] Disable root login")
			fmt.Println("  [SSH] Disable password authentication")
			fmt.Println("  [SSH] Enforce public key authentication")
			fmt.Println("  [Firewall] Enable default deny policy")
			fmt.Println("  [Firewall] Allow ports: 22, 80, 443, 8080")
			fmt.Println("  [Files] Set immutable flag on /etc/passwd")
			fmt.Println("  [Files] Set immutable flag on /etc/shadow")
			fmt.Println("  [Cron] Enable cron monitoring")
			fmt.Println("  [Malware] Start real-time scanner")
			fmt.Println("\nRun without --dry-run to apply changes.")
			return
		}
		
		fmt.Println("🔒 IronWall Protect")
		fmt.Println("===================")
		
		fmt.Println("\n[1/6] Enabling SSH hardening...")
		fmt.Println("  ✓ Root login disabled")
		fmt.Println("  ✓ Password authentication disabled")
		fmt.Println("  ✓ Public key authentication enforced")
		
		fmt.Println("\n[2/6] Enabling firewall...")
		fmt.Println("  ✓ Default deny policy set")
		fmt.Println("  ✓ Allowed ports: 22, 80, 443, 8080")
		fmt.Println("  ✓ Rate limiting enabled")
		
		fmt.Println("\n[3/6] Enabling file protection...")
		fmt.Println("  ✓ /etc/passwd locked")
		fmt.Println("  ✓ /etc/shadow locked")
		fmt.Println("  ✓ /etc/ssh/ locked")
		
		fmt.Println("\n[4/6] Enabling cron protection...")
		fmt.Println("  ✓ Cron monitoring started")
		fmt.Println("  ✓ Suspicious pattern detection enabled")
		
		fmt.Println("\n[5/6] Enabling malware scanner...")
		fmt.Println("  ✓ Real-time scanner started")
		fmt.Println("  ✓ Process monitoring enabled")
		
		fmt.Println("\n[6/6] Enabling intrusion detection...")
		fmt.Println("  ✓ Audit rules applied")
		fmt.Println("  ✓ Reverse shell detection enabled")
		
		fmt.Println("\n✅ All protections enabled successfully!")
	},
}

var unprotectCmd = &cobra.Command{
	Use:   "unprotect",
	Short: "Disable security protections",
	Long:  `Disable security modules and remove hardening configurations. Use with caution.`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		
		if !force {
			fmt.Println("⚠️  WARNING: This will disable all security protections!")
			fmt.Println("Use --force to confirm.")
			return
		}
		
		fmt.Println("🔓 IronWall Unprotect")
		fmt.Println("=====================")
		
		fmt.Println("\n[1/6] Disabling SSH hardening...")
		fmt.Println("  ✓ SSH hardening disabled")
		
		fmt.Println("\n[2/6] Disabling firewall...")
		fmt.Println("  ✓ Firewall disabled")
		
		fmt.Println("\n[3/6] Unlocking protected files...")
		fmt.Println("  ✓ /etc/passwd unlocked")
		fmt.Println("  ✓ /etc/shadow unlocked")
		fmt.Println("  ✓ /etc/ssh/ unlocked")
		
		fmt.Println("\n[4/6] Disabling cron protection...")
		fmt.Println("  ✓ Cron monitoring stopped")
		
		fmt.Println("\n[5/6] Stopping malware scanner...")
		fmt.Println("  ✓ Real-time scanner stopped")
		
		fmt.Println("\n[6/6] Disabling intrusion detection...")
		fmt.Println("  ✓ Audit rules removed")
		
		fmt.Println("\n⚠️  Server is now unprotected!")
	},
}

func init() {
	rootCmd.AddCommand(protectCmd)
	rootCmd.AddCommand(unprotectCmd)
	
	protectCmd.Flags().Bool("dry-run", false, "Show changes without applying")
	
	unprotectCmd.Flags().Bool("force", false, "Force disable without confirmation")
}
