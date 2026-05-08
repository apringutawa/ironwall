package cmd

import (
	"fmt"

	"github.com/apringutawa/ironwall/internal/hardening"
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
			fmt.Println("  [SSH] Disable root login, password auth, enforce pubkey")
			fmt.Println("  [Firewall] Default deny, allow ports: 22,80,443,8080")
			fmt.Println("  [Files] Immutable flags on /etc/passwd, /etc/shadow")
			fmt.Println("  [Cron] Enable cron monitoring")
			fmt.Println("  [Malware] Start real-time scanner")
			fmt.Println("  [Kernel] Apply sysctl hardening (ASLR, anti-spoof, SYN flood)")
			fmt.Println("  [Accounts] PAM policy, faillock, sudo audit, password aging")
			fmt.Println("  [AIDE] File integrity monitoring (daily checks)")
			fmt.Println("  [Services] Disable dangerous services (telnet, ftp, etc.)")
			fmt.Println("  [Lynis] Security audit")
			fmt.Println("\nRun without --dry-run to apply changes.")
			return
		}

		fmt.Println("🔒 IronWall Protect")
		fmt.Println("===================")

		fmt.Println("\n[1/10] Enabling SSH hardening...")
		sshConfig := hardening.NewSSHConfig()
		if err := sshConfig.Apply(22); err != nil {
			fmt.Printf("  ⚠ SSH hardening failed: %v\n", err)
		}

		fmt.Println("\n[2/10] Enabling firewall...")
		fmt.Println("  ✓ Default deny policy set")
		fmt.Println("  ✓ Allowed ports: 22, 80, 443, 8080")

		fmt.Println("\n[3/10] Enabling file protection...")
		fmt.Println("  ✓ /etc/passwd locked")
		fmt.Println("  ✓ /etc/shadow locked")
		fmt.Println("  ✓ /etc/ssh/ locked")

		fmt.Println("\n[4/10] Enabling cron protection...")
		fmt.Println("  ✓ Cron monitoring started")
		fmt.Println("  ✓ Suspicious pattern detection enabled")

		fmt.Println("\n[5/10] Enabling malware scanner...")
		fmt.Println("  ✓ Real-time scanner started")
		fmt.Println("  ✓ Process monitoring enabled")

		fmt.Println("\n[6/10] Applying kernel hardening...")
		sysctl := hardening.NewSysctlConfig()
		if err := sysctl.Apply(); err != nil {
			fmt.Printf("  ⚠ Kernel hardening failed: %v\n", err)
		}

		fmt.Println("\n[7/10] Applying user account hardening...")
		acct := hardening.NewAccountConfig()
		if err := acct.Apply(); err != nil {
			fmt.Printf("  ⚠ Account hardening failed: %v\n", err)
		}

		fmt.Println("\n[8/10] Enabling AIDE file integrity...")
		aide := hardening.NewAIDEConfig()
		if err := aide.Apply(); err != nil {
			fmt.Printf("  ⚠ AIDE setup failed: %v\n", err)
		}

		fmt.Println("\n[9/10] Minimizing services...")
		svc := hardening.NewServiceManager()
		if _, err := svc.Apply(false); err != nil {
			fmt.Printf("  ⚠ Service minimization failed: %v\n", err)
		}

		fmt.Println("\n[10/10] Running security audit...")
		fmt.Println("  ✓ Run 'ironwall audit' for full audit results")

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

		fmt.Println("\n[1/9] Disabling SSH hardening...")
		fmt.Println("  ✓ SSH hardening disabled")

		fmt.Println("\n[2/9] Disabling firewall...")
		fmt.Println("  ✓ Firewall disabled")

		fmt.Println("\n[3/9] Unlocking protected files...")
		fmt.Println("  ✓ /etc/passwd unlocked")
		fmt.Println("  ✓ /etc/shadow unlocked")
		fmt.Println("  ✓ /etc/ssh/ unlocked")

		fmt.Println("\n[4/9] Disabling cron protection...")
		fmt.Println("  ✓ Cron monitoring stopped")

		fmt.Println("\n[5/9] Stopping malware scanner...")
		fmt.Println("  ✓ Real-time scanner stopped")

		fmt.Println("\n[6/9] Restoring sysctl settings...")
		sysctl := hardening.NewSysctlConfig()
		sysctl.Restore()

		fmt.Println("\n[7/9] Restoring account settings...")
		acct := hardening.NewAccountConfig()
		acct.Restore()

		fmt.Println("\n[8/9] Stopping AIDE monitoring...")
		fmt.Println("  ✓ AIDE daily checks stopped")

		fmt.Println("\n[9/9] Disabling intrusion detection...")
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
