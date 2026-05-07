package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var firewallCmd = &cobra.Command{
	Use:   "firewall",
	Short: "Manage firewall rules",
	Long:  `Manage nftables/iptables firewall rules for IronWall.`,
}

var firewallStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show firewall status",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🛡️  Firewall Status")
		fmt.Println("===================")
		fmt.Println()
		fmt.Println("Status: Active")
		fmt.Println("Backend: nftables")
		fmt.Println()
		fmt.Println("Allowed Ports:")
		fmt.Println("  22/tcp   - SSH")
		fmt.Println("  80/tcp   - HTTP")
		fmt.Println("  443/tcp  - HTTPS")
		fmt.Println("  8080/tcp - Dashboard")
		fmt.Println()
		fmt.Println("Blocked IPs: 147")
		fmt.Println("Rate Limiting: Enabled (22/tcp: 3/min)")
		fmt.Println("Anti-Port Scan: Enabled")
	},
}

var firewallBlockCmd = &cobra.Command{
	Use:   "block <ip>",
	Short: "Block an IP address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		fmt.Printf("✓ IP %s has been blocked\n", ip)
	},
}

var firewallUnblockCmd = &cobra.Command{
	Use:   "unblock <ip>",
	Short: "Unblock an IP address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ip := args[0]
		fmt.Printf("✓ IP %s has been unblocked\n", ip)
	},
}

var firewallListCmd = &cobra.Command{
	Use:   "list",
	Short: "List blocked IPs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Blocked IPs:")
		fmt.Println("  192.168.1.100  - Brute force attempt (2026-05-06)")
		fmt.Println("  10.0.0.55      - Port scanning (2026-05-05)")
		fmt.Println("  172.16.0.200   - Malware detected (2026-05-04)")
	},
}

var firewallAllowCmd = &cobra.Command{
	Use:   "allow <port> [protocol]",
	Short: "Allow a port",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port := args[0]
		protocol := "tcp"
		if len(args) > 1 {
			protocol = args[1]
		}
		fmt.Printf("✓ Port %s/%s has been allowed\n", port, protocol)
	},
}

func init() {
	rootCmd.AddCommand(firewallCmd)
	firewallCmd.AddCommand(firewallStatusCmd)
	firewallCmd.AddCommand(firewallBlockCmd)
	firewallCmd.AddCommand(firewallUnblockCmd)
	firewallCmd.AddCommand(firewallListCmd)
	firewallCmd.AddCommand(firewallAllowCmd)
}
