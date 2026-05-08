package hardening

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SysctlConfig struct {
	ConfigPath string
	BackupPath string
}

type SysctlParam struct {
	Key   string
	Value string
	Desc  string
}

func NewSysctlConfig() *SysctlConfig {
	return &SysctlConfig{
		ConfigPath: "/etc/sysctl.d/99-ironwall.conf",
		BackupPath: "/var/backups/ironwall/sysctl.conf",
	}
}

func (s *SysctlConfig) GetParams() []SysctlParam {
	return []SysctlParam{
		// IP Spoofing protection
		{"net.ipv4.conf.all.rp_filter", "1", "Enable IP spoofing protection"},
		{"net.ipv4.conf.default.rp_filter", "1", "Enable IP spoofing protection (default)"},

		// Ignore ICMP redirects
		{"net.ipv4.conf.all.accept_redirects", "0", "Ignore ICMP redirects"},
		{"net.ipv4.conf.default.accept_redirects", "0", "Ignore ICMP redirects (default)"},
		{"net.ipv6.conf.all.accept_redirects", "0", "Ignore IPv6 ICMP redirects"},
		{"net.ipv6.conf.default.accept_redirects", "0", "Ignore IPv6 ICMP redirects (default)"},

		// Ignore source-routed packets
		{"net.ipv4.conf.all.accept_source_route", "0", "Ignore source-routed packets"},
		{"net.ipv4.conf.default.accept_source_route", "0", "Ignore source-routed packets (default)"},
		{"net.ipv6.conf.all.accept_source_route", "0", "Ignore IPv6 source-routed packets"},
		{"net.ipv6.conf.default.accept_source_route", "0", "Ignore IPv6 source-routed packets (default)"},

		// SYN flood protection
		{"net.ipv4.tcp_syncookies", "1", "Enable SYN cookies"},
		{"net.ipv4.tcp_syn_retries", "2", "Reduce SYN retries"},
		{"net.ipv4.tcp_synack_retries", "2", "Reduce SYN-ACK retries"},
		{"net.ipv4.tcp_max_syn_backlog", "2048", "Increase SYN backlog"},

		// ASLR
		{"kernel.randomize_va_space", "2", "Enable full ASLR"},

		// IP forwarding
		{"net.ipv4.ip_forward", "0", "Disable IP forwarding"},
		{"net.ipv6.conf.all.forwarding", "0", "Disable IPv6 forwarding"},

		// Kernel panic
		{"kernel.panic", "10", "Kernel panic timeout"},
		{"kernel.panic_on_oops", "1", "Panic on kernel oops"},

		// Magic SysRq
		{"kernel.sysrq", "0", "Disable Magic SysRq"},

		// Core dumps
		{"fs.suid_dumpable", "0", "Disable core dumps for setuid"},

		// Network hardening
		{"net.ipv4.tcp_rfc1337", "1", "Enable TCP RFC 1337"},
		{"net.ipv4.icmp_echo_ignore_broadcasts", "1", "Ignore ICMP echo broadcasts"},
		{"net.ipv4.icmp_ignore_bogus_error_responses", "1", "Ignore bogus ICMP errors"},
		{"net.ipv4.conf.all.log_martians", "1", "Log martian packets"},
		{"net.ipv4.conf.default.log_martians", "1", "Log martian packets (default)"},

		// TCP hardening
		{"net.ipv4.tcp_timestamps", "0", "Disable TCP timestamps"},
		{"net.ipv4.tcp_fin_timeout", "15", "Reduce FIN timeout"},
		{"net.ipv4.tcp_tw_reuse", "1", "Enable TIME_WAIT reuse"},
	}
}

func (s *SysctlConfig) Backup() error {
	if _, err := os.Stat(s.ConfigPath); os.IsNotExist(err) {
		return nil
	}
	input, err := os.ReadFile(s.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read sysctl config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(s.BackupPath), 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}
	return os.WriteFile(s.BackupPath, input, 0644)
}

func (s *SysctlConfig) Restore() error {
	if _, err := os.Stat(s.BackupPath); os.IsNotExist(err) {
		return fmt.Errorf("no backup found at %s", s.BackupPath)
	}
	input, err := os.ReadFile(s.BackupPath)
	if err != nil {
		return err
	}
	return os.WriteFile(s.ConfigPath, input, 0644)
}

func (s *SysctlConfig) Apply() error {
	fmt.Println("Applying kernel hardening (sysctl)...")

	if err := s.Backup(); err != nil {
		fmt.Printf("  ⚠ Backup warning: %v\n", err)
	}

	var lines []string
	lines = append(lines,
		"# IronWall Kernel Hardening",
		"# Applied: "+fmt.Sprintf("%v", os.Getpid()),
		"",
	)

	for _, param := range s.GetParams() {
		lines = append(lines, fmt.Sprintf("# %s", param.Desc))
		lines = append(lines, fmt.Sprintf("%s = %s", param.Key, param.Value))
		lines = append(lines, "")
	}

	content := strings.Join(lines, "\n")
	if err := os.MkdirAll("/etc/sysctl.d", 0755); err != nil {
		return fmt.Errorf("failed to create sysctl.d: %w", err)
	}
	if err := os.WriteFile(s.ConfigPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write sysctl config: %w", err)
	}
	fmt.Println("  ✓ Config written to", s.ConfigPath)

	cmd := exec.Command("sysctl", "-p", s.ConfigPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to apply sysctl: %w\n%s", err, string(output))
	}
	fmt.Println("  ✓ Kernel parameters applied")

	return nil
}

func (s *SysctlConfig) ApplyParam(key, value string) error {
	cmd := exec.Command("sysctl", "-w", fmt.Sprintf("%s=%s", key, value))
	return cmd.Run()
}

func (s *SysctlConfig) GetCurrent(key string) (string, error) {
	cmd := exec.Command("sysctl", "-n", key)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (s *SysctlConfig) Verify() ([]string, error) {
	var failed []string

	for _, param := range s.GetParams() {
		current, err := s.GetCurrent(param.Key)
		if err != nil || current != param.Value {
			failed = append(failed, param.Key)
		}
	}

	return failed, nil
}
