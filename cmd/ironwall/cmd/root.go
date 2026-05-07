package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ironwall",
	Short: "IronWall - Automated Linux System Hardening Toolkit",
	Long: `IronWall is an automated Linux server hardening and threat prevention toolkit.
It provides comprehensive security including SSH hardening, firewall automation,
malware detection, intrusion prevention, and real-time monitoring.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
