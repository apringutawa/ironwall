package firewall

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Firewall struct {
	Backend string
}

func NewFirewall() *Firewall {
	backend := detectBackend()
	return &Firewall{Backend: backend}
}

func detectBackend() string {
	if _, err := exec.LookPath("nft"); err == nil {
		return "nftables"
	}
	if _, err := exec.LookPath("iptables"); err == nil {
		return "iptables"
	}
	if _, err := exec.LookPath("firewall-cmd"); err == nil {
		return "firewalld"
	}
	return "unknown"
}

func (f *Firewall) Setup() error {
	fmt.Printf("Setting up firewall (%s)...\n", f.Backend)

	switch f.Backend {
	case "nftables":
		return f.setupNftables()
	case "iptables":
		return f.setupIptables()
	case "firewalld":
		return f.setupFirewalld()
	default:
		return fmt.Errorf("no supported firewall backend found")
	}
}

func (f *Firewall) setupNftables() error {
	rules := `#!/usr/sbin/nft -f

flush ruleset

table inet filter {
    chain input {
        type filter hook input priority 0;
        
        # Allow established connections
        ct state established,related accept
        
        # Allow loopback
        iif "lo" accept
        
        # Allow SSH (port 22)
        tcp dport 22 accept
        
        # Allow HTTP (port 80)
        tcp dport 80 accept
        
        # Allow HTTPS (port 443)
        tcp dport 443 accept
        
        # Allow IronWall dashboard (port 8080)
        tcp dport 8080 accept
        
        # Drop invalid packets
        ct state invalid drop
        
        # Default policy
        drop
    }
    
    chain forward {
        type filter hook forward priority 0;
        drop
    }
    
    chain output {
        type filter hook output priority 0;
        accept
    }
}
`

	if err := os.WriteFile("/etc/nftables.conf", []byte(rules), 0644); err != nil {
		return err
	}

	cmd := exec.Command("systemctl", "enable", "nftables")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("systemctl", "restart", "nftables")
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("  ✓ nftables configured")
	return nil
}

func (f *Firewall) setupIptables() error {
	commands := [][]string{
		{"iptables", "-F"},
		{"iptables", "-P", "INPUT", "DROP"},
		{"iptables", "-P", "FORWARD", "DROP"},
		{"iptables", "-P", "OUTPUT", "ACCEPT"},
		{"iptables", "-A", "INPUT", "-m", "state", "--state", "ESTABLISHED,RELATED", "-j", "ACCEPT"},
		{"iptables", "-A", "INPUT", "-i", "lo", "-j", "ACCEPT"},
		{"iptables", "-A", "INPUT", "-p", "tcp", "--dport", "22", "-j", "ACCEPT"},
		{"iptables", "-A", "INPUT", "-p", "tcp", "--dport", "80", "-j", "ACCEPT"},
		{"iptables", "-A", "INPUT", "-p", "tcp", "--dport", "443", "-j", "ACCEPT"},
		{"iptables", "-A", "INPUT", "-p", "tcp", "--dport", "8080", "-j", "ACCEPT"},
	}

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			return fmt.Errorf("failed to run %s: %w", strings.Join(cmd, " "), err)
		}
	}

	cmd := exec.Command("iptables-save")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	if err := os.WriteFile("/etc/iptables/rules.v4", output, 0644); err != nil {
		return err
	}

	fmt.Println("  ✓ iptables configured")
	return nil
}

func (f *Firewall) setupFirewalld() error {
	commands := [][]string{
		{"firewall-cmd", "--permanent", "--add-port=22/tcp"},
		{"firewall-cmd", "--permanent", "--add-port=80/tcp"},
		{"firewall-cmd", "--permanent", "--add-port=443/tcp"},
		{"firewall-cmd", "--permanent", "--add-port=8080/tcp"},
		{"firewall-cmd", "--reload"},
	}

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			return fmt.Errorf("failed to run %s: %w", strings.Join(cmd, " "), err)
		}
	}

	fmt.Println("  ✓ firewalld configured")
	return nil
}

func (f *Firewall) BlockIP(ip string) error {
	switch f.Backend {
	case "nftables":
		cmd := exec.Command("nft", "add", "element", "inet", "filter", "blacklist", "{", ip, "}")
		return cmd.Run()
	case "iptables":
		cmd := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-j", "DROP")
		return cmd.Run()
	case "firewalld":
		cmd := exec.Command("firewall-cmd", "--permanent", "--add-rich-rule=rule family='ipv4' source address='"+ip+"' reject")
		if err := cmd.Run(); err != nil {
			return err
		}
		cmd = exec.Command("firewall-cmd", "--reload")
		return cmd.Run()
	}
	return fmt.Errorf("unsupported backend")
}

func (f *Firewall) UnblockIP(ip string) error {
	switch f.Backend {
	case "nftables":
		cmd := exec.Command("nft", "delete", "element", "inet", "filter", "blacklist", "{", ip, "}")
		return cmd.Run()
	case "iptables":
		cmd := exec.Command("iptables", "-D", "INPUT", "-s", ip, "-j", "DROP")
		return cmd.Run()
	case "firewalld":
		cmd := exec.Command("firewall-cmd", "--permanent", "--remove-rich-rule=rule family='ipv4' source address='"+ip+"' reject")
		if err := cmd.Run(); err != nil {
			return err
		}
		cmd = exec.Command("firewall-cmd", "--reload")
		return cmd.Run()
	}
	return fmt.Errorf("unsupported backend")
}

func (f *Firewall) AllowPort(port int, protocol string) error {
	portStr := fmt.Sprintf("%d", port)

	switch f.Backend {
	case "nftables":
		cmd := exec.Command("nft", "add", "rule", "inet", "filter", "input", protocol, "dport", portStr, "accept")
		return cmd.Run()
	case "iptables":
		cmd := exec.Command("iptables", "-A", "INPUT", "-p", protocol, "--dport", portStr, "-j", "ACCEPT")
		return cmd.Run()
	case "firewalld":
		cmd := exec.Command("firewall-cmd", "--permanent", fmt.Sprintf("--add-port=%s/%s", portStr, protocol))
		if err := cmd.Run(); err != nil {
			return err
		}
		cmd = exec.Command("firewall-cmd", "--reload")
		return cmd.Run()
	}
	return fmt.Errorf("unsupported backend")
}
