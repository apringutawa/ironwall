package cmd

import (
	"fmt"

	"github.com/apringutawa/ironwall/internal/hardening"
	"github.com/spf13/cobra"
)

var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage user account hardening",
	Long:  `Harden user accounts: password policy, faillock, sudo audit, and more.`,
}

var accountsHardenCmd = &cobra.Command{
	Use:   "harden",
	Short: "Apply user account hardening",
	Run: func(cmd *cobra.Command, args []string) {
		acct := hardening.NewAccountConfig()
		if err := acct.Apply(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("\n✅ Account hardening applied successfully")
	},
}

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all human users on the system",
	Run: func(cmd *cobra.Command, args []string) {
		acct := hardening.NewAccountConfig()
		users, err := acct.ListUsers()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println("👤 System Users")
		fmt.Println("================")
		for _, u := range users {
			locked := u["locked"]
			lockIcon := "✓"
			if locked == "true" {
				lockIcon = "🔒"
			}
			fmt.Printf("  %s %-15s UID:%-5s Shell: %s\n", lockIcon, u["username"], u["uid"], u["shell"])
		}
	},
}

var accountsCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for security issues in user accounts",
	Run: func(cmd *cobra.Command, args []string) {
		acct := hardening.NewAccountConfig()

		fmt.Println("🔍 Security Check: User Accounts")
		fmt.Println("==================================")

		acct.CheckUIDZero()
		acct.CheckEmptyPasswords()
		acct.CheckSshdConfig()

		fmt.Println("\n✅ Check completed")
	},
}

var accountsRestoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore previous account configuration",
	Run: func(cmd *cobra.Command, args []string) {
		acct := hardening.NewAccountConfig()
		if err := acct.Restore(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("✅ Account configuration restored")
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.AddCommand(accountsHardenCmd)
	accountsCmd.AddCommand(accountsListCmd)
	accountsCmd.AddCommand(accountsCheckCmd)
	accountsCmd.AddCommand(accountsRestoreCmd)
}
