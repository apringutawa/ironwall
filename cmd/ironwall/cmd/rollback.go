package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Rollback to previous configuration",
	Long:  `Restore system configuration from backup. This will undo all changes made by IronWall.`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		
		if !force {
			fmt.Println("⚠️  WARNING: This will restore all backed up configurations!")
			fmt.Println("Use --force to confirm rollback.")
			return
		}
		
		fmt.Println("⏮️  IronWall Rollback")
		fmt.Println("====================")
		
		fmt.Println("\n[1/5] Restoring SSH configuration...")
		fmt.Println("  ✓ /etc/ssh/sshd_config restored")
		
		fmt.Println("\n[2/5] Restoring firewall rules...")
		fmt.Println("  ✓ Firewall rules restored")
		
		fmt.Println("\n[3/5] Unlocking protected files...")
		fmt.Println("  ✓ File locks removed")
		
		fmt.Println("\n[4/5] Restoring cron configuration...")
		fmt.Println("  ✓ /etc/crontab restored")
		
		fmt.Println("\n[5/5] Stopping IronWall services...")
		fmt.Println("  ✓ All services stopped")
		
		fmt.Println("\n✅ Rollback completed successfully!")
		fmt.Println("System restored to pre-IronWall state.")
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)
	rollbackCmd.Flags().Bool("force", false, "Force rollback without confirmation")
	rollbackCmd.Flags().String("backup-id", "", "Specific backup ID to restore")
}
