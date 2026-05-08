package cmd

import (
	"fmt"

	"github.com/apringutawa/ironwall/internal/hardening"
	"github.com/spf13/cobra"
)

var aideCmd = &cobra.Command{
	Use:   "aide",
	Short: "Manage AIDE file integrity monitoring",
	Long:  `Initialize, check, and manage AIDE file integrity database.`,
}

var aideInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize AIDE database",
	Run: func(cmd *cobra.Command, args []string) {
		aide := hardening.NewAIDEConfig()
		if err := aide.Initialize(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("✅ AIDE database initialized")
	},
}

var aideCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Run AIDE integrity check",
	Run: func(cmd *cobra.Command, args []string) {
		aide := hardening.NewAIDEConfig()
		if err := aide.Check(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("✅ AIDE integrity check completed")
	},
}

var aideUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update AIDE database",
	Run: func(cmd *cobra.Command, args []string) {
		aide := hardening.NewAIDEConfig()
		if err := aide.Update(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("✅ AIDE database updated")
	},
}

var aideEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable AIDE monitoring with daily checks",
	Run: func(cmd *cobra.Command, args []string) {
		aide := hardening.NewAIDEConfig()
		if err := aide.Apply(); err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("✅ AIDE monitoring enabled with daily checks")
	},
}

func init() {
	rootCmd.AddCommand(aideCmd)
	aideCmd.AddCommand(aideInitCmd)
	aideCmd.AddCommand(aideCheckCmd)
	aideCmd.AddCommand(aideUpdateCmd)
	aideCmd.AddCommand(aideEnableCmd)
}
