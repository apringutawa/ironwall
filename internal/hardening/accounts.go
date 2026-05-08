package hardening

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type AccountConfig struct {
	PamConfigPath string
	LoginDefsPath string
	BackupDir     string
}

func NewAccountConfig() *AccountConfig {
	return &AccountConfig{
		PamConfigPath: "/etc/pam.d/common-password",
		LoginDefsPath: "/etc/login.defs",
		BackupDir:     "/var/backups/ironwall/accounts",
	}
}

func (a *AccountConfig) Backup() error {
	if err := os.MkdirAll(a.BackupDir, 0755); err != nil {
		return err
	}

	files := []string{a.PamConfigPath, a.LoginDefsPath}
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}
		input, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		backupPath := filepath.Join(a.BackupDir, filepath.Base(file))
		if err := os.WriteFile(backupPath, input, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (a *AccountConfig) Apply() error {
	fmt.Println("Applying user account hardening...")

	if err := a.Backup(); err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}
	fmt.Println("  ✓ Backup created")

	if err := a.HardenPasswordPolicy(); err != nil {
		fmt.Printf("  ⚠ Password policy: %v\n", err)
	} else {
		fmt.Println("  ✓ Password policy enforced (min length 12, complexity enabled)")
	}

	if err := a.ConfigureFaillock(); err != nil {
		fmt.Printf("  ⚠ Faillock: %v\n", err)
	} else {
		fmt.Println("  ✓ Faillock configured (5 attempts, 15min lockout)")
	}

	if err := a.HardenSudo(); err != nil {
		fmt.Printf("  ⚠ Sudo audit: %v\n", err)
	} else {
		fmt.Println("  ✓ Sudo audit logging enabled")
	}

	if err := a.CheckUIDZero(); err != nil {
		fmt.Printf("  ⚠ UID 0 check: %v\n", err)
	}

	if err := a.CheckEmptyPasswords(); err != nil {
		fmt.Printf("  ⚠ Empty password check: %v\n", err)
	}

	if err := a.SetPasswordAging(); err != nil {
		fmt.Printf("  ⚠ Password aging: %v\n", err)
	} else {
		fmt.Println("  ✓ Password aging configured (90 days max)")
	}

	return nil
}

func (a *AccountConfig) HardenPasswordPolicy() error {
	if _, err := os.Stat(a.PamConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("pam config not found at %s", a.PamConfigPath)
	}

	content, err := os.ReadFile(a.PamConfigPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	newLines := make([]string, 0, len(lines))
	pamUpdated := false

	for _, line := range lines {
		if strings.Contains(line, "pam_unix.so") && !strings.HasPrefix(line, "#") {
			if !strings.Contains(line, "minlen=") {
				newLines = append(newLines, line)
				newLines = append(newLines, "password requisite pam_pwquality.so retry=3 minlen=12 difok=3 ucredit=-1 lcredit=-1 dcredit=-1 ocredit=-1")
				pamUpdated = true
				continue
			}
		}
		newLines = append(newLines, line)
	}

	if pamUpdated {
		return os.WriteFile(a.PamConfigPath, []byte(strings.Join(newLines, "\n")), 0644)
	}
	return nil
}

func (a *AccountConfig) ConfigureFaillock() error {
	faillockPath := "/etc/pam.d/common-auth"
	if _, err := os.Stat(faillockPath); os.IsNotExist(err) {
		faillockPath = "/etc/pam.d/system-auth"
		if _, err := os.Stat(faillockPath); os.IsNotExist(err) {
			return fmt.Errorf("pam auth config not found")
		}
	}

	content, err := os.ReadFile(faillockPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	hasFaillock := false
	for _, line := range lines {
		if strings.Contains(line, "pam_faillock.so") {
			hasFaillock = true
			break
		}
	}

	if hasFaillock {
		return nil
	}

	entries := []string{
		"auth    required    pam_faillock.so preauth silent deny=5 unlock_time=900",
		"auth    [default=die]  pam_faillock.so authfail deny=5 unlock_time=900",
		"account required    pam_faillock.so",
	}

	newContent := string(content) + "\n" + strings.Join(entries, "\n") + "\n"
	return os.WriteFile(faillockPath, []byte(newContent), 0644)
}

func (a *AccountConfig) HardenSudo() error {
	sudoersFile := "/etc/sudoers.d/ironwall"
	content := []byte("# IronWall sudo audit configuration\nDefaults    log_input, log_output\nDefaults    logfile=/var/log/sudo.log\n")

	if err := os.WriteFile(sudoersFile, content, 0440); err != nil {
		return fmt.Errorf("failed to write sudoers config: %w", err)
	}
	return nil
}

func (a *AccountConfig) CheckUIDZero() error {
	content, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) >= 3 && parts[2] == "0" && parts[0] != "root" {
			fmt.Printf("  ⚠ WARNING: Non-root user %s has UID 0!\n", parts[0])
		}
	}
	return nil
}

func (a *AccountConfig) CheckEmptyPasswords() error {
	content, err := os.ReadFile("/etc/shadow")
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) >= 2 && parts[1] == "" {
			fmt.Printf("  ⚠ WARNING: User %s has empty password!\n", parts[0])
		}
	}
	return nil
}

func (a *AccountConfig) SetPasswordAging() error {
	content, err := os.ReadFile(a.LoginDefsPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	updated := false

	for i, line := range lines {
		if strings.HasPrefix(line, "PASS_MAX_DAYS") {
			lines[i] = "PASS_MAX_DAYS   90"
			updated = true
		}
		if strings.HasPrefix(line, "PASS_MIN_DAYS") {
			lines[i] = "PASS_MIN_DAYS   7"
			updated = true
		}
		if strings.HasPrefix(line, "PASS_WARN_AGE") {
			lines[i] = "PASS_WARN_AGE   14"
			updated = true
		}
	}

	if updated {
		return os.WriteFile(a.LoginDefsPath, []byte(strings.Join(lines, "\n")), 0644)
	}
	return nil
}

func (a *AccountConfig) Restore() error {
	backupFiles := []string{"common-password", "login.defs"}
	for _, file := range backupFiles {
		backupPath := filepath.Join(a.BackupDir, file)
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			continue
		}
		input, err := os.ReadFile(backupPath)
		if err != nil {
			return err
		}

		var dstPath string
		switch file {
		case "common-password":
			dstPath = a.PamConfigPath
		case "login.defs":
			dstPath = a.LoginDefsPath
		}

		if err := os.WriteFile(dstPath, input, 0644); err != nil {
			return err
		}
	}

	ironwallSudo := "/etc/sudoers.d/ironwall"
	if _, err := os.Stat(ironwallSudo); err == nil {
		os.Remove(ironwallSudo)
	}

	return nil
}

func (a *AccountConfig) ListUsers() ([]map[string]string, error) {
	content, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return nil, err
	}

	shadowContent, err := os.ReadFile("/etc/shadow")
	if err != nil {
		return nil, err
	}

	shadowMap := make(map[string]string)
	for _, line := range strings.Split(string(shadowContent), "\n") {
		parts := strings.Split(line, ":")
		if len(parts) >= 2 {
			shadowMap[parts[0]] = parts[1]
		}
	}

	lines := strings.Split(string(content), "\n")
	var users []map[string]string

	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 7 {
			continue
		}

		uid, _ := strconv.Atoi(parts[2])
		if uid < 1000 && uid != 0 {
			continue
		}

		user := map[string]string{
			"username": parts[0],
			"uid":      parts[2],
			"gid":      parts[3],
			"shell":    parts[6],
			"home":     parts[5],
		}

		if hash, ok := shadowMap[parts[0]]; ok {
			if hash == "" || strings.HasPrefix(hash, "!") || strings.HasPrefix(hash, "*") {
				user["locked"] = "true"
			}
		}

		users = append(users, user)
	}

	return users, nil
}

func (a *AccountConfig) LockUser(username string) error {
	cmd := exec.Command("passwd", "-l", username)
	return cmd.Run()
}

func (a *AccountConfig) UnlockUser(username string) error {
	cmd := exec.Command("passwd", "-u", username)
	return cmd.Run()
}

func (a *AccountConfig) CheckSshdConfig() error {
	sshdPath := "/etc/ssh/sshd_config"
	if _, err := os.Stat(sshdPath); os.IsNotExist(err) {
		return nil
	}

	content, err := os.ReadFile(sshdPath)
	if err != nil {
		return err
	}

	checks := map[string]string{
		`^PermitEmptyPasswords\s+yes`: "PermitEmptyPasswords is enabled in SSH",
		`^PermitRootLogin\s+yes`:      "Root login is allowed in SSH",
		`^PasswordAuthentication\s+yes`: "Password authentication is enabled in SSH",
	}

	for pattern, msg := range checks {
		re := regexp.MustCompile(pattern)
		if re.MatchString(string(content)) {
			fmt.Printf("  ⚠ %s\n", msg)
		}
	}

	return nil
}
