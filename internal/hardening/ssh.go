package hardening

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type SSHConfig struct {
	ConfigPath string
	BackupPath string
}

func NewSSHConfig() *SSHConfig {
	return &SSHConfig{
		ConfigPath: "/etc/ssh/sshd_config",
		BackupPath: "/var/backups/ironwall/sshd_config",
	}
}

func (s *SSHConfig) Backup() error {
	input, err := os.ReadFile(s.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read sshd_config: %w", err)
	}

	if err := os.MkdirAll("/var/backups/ironwall", 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	if err := os.WriteFile(s.BackupPath, input, 0644); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}

	return nil
}

func (s *SSHConfig) parseConfig() (map[string]string, error) {
	config := make(map[string]string)

	file, err := os.Open(s.ConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return config, scanner.Err()
}

func (s *SSHConfig) updateConfig(key, value string) error {
	input, err := os.ReadFile(s.ConfigPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")
	found := false
	re := regexp.MustCompile(`^#?\s*` + regexp.QuoteMeta(key) + `\s+`)

	for i, line := range lines {
		if re.MatchString(line) {
			lines[i] = key + " " + value
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, key+" "+value)
	}

	output := strings.Join(lines, "\n")
	return os.WriteFile(s.ConfigPath, []byte(output), 0644)
}

func (s *SSHConfig) DisableRootLogin() error {
	return s.updateConfig("PermitRootLogin", "no")
}

func (s *SSHConfig) DisablePasswordAuth() error {
	return s.updateConfig("PasswordAuthentication", "no")
}

func (s *SSHConfig) EnablePublicKeyAuth() error {
	return s.updateConfig("PubkeyAuthentication", "yes")
}

func (s *SSHConfig) ChangePort(port int) error {
	return s.updateConfig("Port", fmt.Sprintf("%d", port))
}

func (s *SSHConfig) RestartSSH() error {
	cmd := exec.Command("systemctl", "restart", "sshd")
	if err := cmd.Run(); err != nil {
		cmd = exec.Command("systemctl", "restart", "ssh")
		return cmd.Run()
	}
	return nil
}

func (s *SSHConfig) Apply(sshPort int) error {
	fmt.Println("Applying SSH hardening...")

	if err := s.Backup(); err != nil {
		return err
	}
	fmt.Println("  ✓ Backup created")

	if err := s.DisableRootLogin(); err != nil {
		return err
	}
	fmt.Println("  ✓ Root login disabled")

	if err := s.DisablePasswordAuth(); err != nil {
		return err
	}
	fmt.Println("  ✓ Password authentication disabled")

	if err := s.EnablePublicKeyAuth(); err != nil {
		return err
	}
	fmt.Println("  ✓ Public key authentication enabled")

	if sshPort != 22 {
		if err := s.ChangePort(sshPort); err != nil {
			return err
		}
		fmt.Printf("  ✓ SSH port changed to %d\n", sshPort)
	}

	if err := s.RestartSSH(); err != nil {
		return err
	}
	fmt.Println("  ✓ SSH service restarted")

	return nil
}

func (s *SSHConfig) Restore() error {
	input, err := os.ReadFile(s.BackupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup: %w", err)
	}

	if err := os.WriteFile(s.ConfigPath, input, 0644); err != nil {
		return fmt.Errorf("failed to restore config: %w", err)
	}

	return s.RestartSSH()
}
