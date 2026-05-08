package hardening

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type AIDEConfig struct {
	ConfigPath   string
	DatabasePath string
	BackupDir    string
}

func NewAIDEConfig() *AIDEConfig {
	return &AIDEConfig{
		ConfigPath:   "/etc/aide/aide.conf",
		DatabasePath: "/var/lib/aide/aide.db.gz",
		BackupDir:    "/var/backups/ironwall/aide",
	}
}

func (a *AIDEConfig) IsInstalled() bool {
	_, err := exec.LookPath("aide")
	return err == nil
}

func (a *AIDEConfig) EnsureInstalled() error {
	if a.IsInstalled() {
		return nil
	}

	fmt.Println("Installing AIDE...")
	cmd := exec.Command("apt-get", "install", "-y", "-qq", "aide", "aide-common")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install aide: %w\n%s", err, string(output))
	}
	fmt.Println("  ✓ AIDE installed")
	return nil
}

func (a *AIDEConfig) createDefaultConfig() error {
	if err := os.MkdirAll("/etc/aide", 0755); err != nil {
		return err
	}

	config := `# IronWall AIDE Configuration
database=file:@@{DBDIR}/aide.db.gz
database_out=file:@@{DBDIR}/aide.db.gz
gzip_dbout=yes
verbose=5
report_url=file:@@{LOGDIR}/aide.log

# Permissions, ownership, group, size, hash
Bin= p+i+n+u+g+s+b+m+c+sha256
Log= p+i+n+u+g+s+b+m+c+sha512
Dir= p+i+n+u+g
Perm= p+i+u+g
Static= p+i+n+u+g+s+b+m+c+sha256

# System binaries
/bin Bin
/sbin Bin
/usr/bin Bin
/usr/sbin Bin

# Libraries
/lib Dir
/usr/lib Dir

# Configuration files
/etc/passwd Perm
/etc/shadow Perm
/etc/group Perm
/etc/gshadow Perm
/etc/sudoers Perm
/etc/ssh/sshd_config Static
/etc/crontab Static
/etc/hosts.allow Static
/etc/hosts.deny Static

# System configuration
/etc/fstab Perm
/etc/issue.net Static
/etc/sysctl.conf Static
/etc/sysctl.d Dir

# Web application files
/var/www Dir

# Important system files
/boot Dir
/root Dir

# Exclude directories
!/var/log
!/var/spool
!/var/tmp
!/tmp
!/run
!/proc
!/sys
!/dev
`

	return os.WriteFile(a.ConfigPath, []byte(config), 0644)
}

func (a *AIDEConfig) Initialize() error {
	if err := a.EnsureInstalled(); err != nil {
		return err
	}

	if err := os.MkdirAll("/var/lib/aide", 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(a.BackupDir, 0755); err != nil {
		return err
	}

	if _, err := os.Stat(a.ConfigPath); os.IsNotExist(err) {
		if err := a.createDefaultConfig(); err != nil {
			return fmt.Errorf("failed to create config: %w", err)
		}
		fmt.Println("  ✓ Default config created")
	}

	fmt.Println("Initializing AIDE database...")
	cmd := exec.Command("aideinit", "--yes")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("aideinit failed: %w\n%s", err, string(output))
	}

	srcDb := "/var/lib/aide/aide.db.new.gz"
	if _, err := os.Stat(srcDb); err == nil {
		os.Rename(srcDb, a.DatabasePath)
	}

	backupDb := filepath.Join(a.BackupDir, "aide.db.gz")
	input, err := os.ReadFile(a.DatabasePath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(backupDb, input, 0644); err != nil {
		return err
	}

	fmt.Println("  ✓ AIDE database initialized")
	return nil
}

func (a *AIDEConfig) Check() error {
	if err := a.EnsureInstalled(); err != nil {
		return err
	}

	if _, err := os.Stat(a.DatabasePath); os.IsNotExist(err) {
		return fmt.Errorf("AIDE database not found, run init first")
	}

	fmt.Println("Running AIDE integrity check...")
	cmd := exec.Command("aide", "--check")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("aide check failed: %w\n%s", err, string(output))
	}

	fmt.Println(string(output))
	return nil
}

func (a *AIDEConfig) Update() error {
	if err := a.EnsureInstalled(); err != nil {
		return err
	}

	if _, err := os.Stat(a.DatabasePath); os.IsNotExist(err) {
		return a.Initialize()
	}

	fmt.Println("Updating AIDE database...")
	cmd := exec.Command("aide", "--update")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("aide update failed: %w\n%s", err, string(output))
	}

	srcDb := "/var/lib/aide/aide.db.new.gz"
	if _, err := os.Stat(srcDb); err == nil {
		os.Rename(srcDb, a.DatabasePath)
	}

	fmt.Println("  ✓ AIDE database updated")
	return nil
}

func (a *AIDEConfig) Apply() error {
	fmt.Println("Setting up AIDE file integrity monitoring...")

	if err := a.EnsureInstalled(); err != nil {
		return err
	}
	fmt.Println("  ✓ AIDE installed")

	if _, err := os.Stat(a.DatabasePath); os.IsNotExist(err) {
		if err := a.Initialize(); err != nil {
			return err
		}
	} else {
		if err := a.Update(); err != nil {
			return err
		}
	}

	serviceContent := `[Unit]
Description=AIDE Integrity Check
After=network.target

[Service]
Type=oneshot
ExecStart=/usr/bin/aide --check
`

	if err := os.WriteFile("/etc/systemd/system/aide-check.service", []byte(serviceContent), 0644); err != nil {
		return fmt.Errorf("failed to create systemd service: %w", err)
	}

	timerContent := `[Unit]
Description=Daily AIDE Integrity Check

[Timer]
OnCalendar=daily
RandomizedDelaySec=30min
Persistent=true

[Install]
WantedBy=timers.target
`

	timerPath := "/etc/systemd/system/aide-check.timer"
	if err := os.WriteFile(timerPath, []byte(timerContent), 0644); err != nil {
		return fmt.Errorf("failed to create timer: %w", err)
	}

	exec.Command("systemctl", "daemon-reload").Run()
	exec.Command("systemctl", "enable", "aide-check.timer").Run()
	exec.Command("systemctl", "start", "aide-check.timer").Run()

	fmt.Println("  ✓ AIDE daily check timer enabled")
	return nil
}

func (a *AIDEConfig) GetLastCheck() (string, error) {
	logFile := "/var/log/aide/aide.log"
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return "never", nil
	}

	content, err := os.ReadFile(logFile)
	if err != nil {
		return "unknown", err
	}

	lines := strings.Split(string(content), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], "AIDE") {
			return strings.TrimSpace(lines[i]), nil
		}
	}

	return "unknown", nil
}

func (a *AIDEConfig) DatabaseExists() bool {
	_, err := os.Stat(a.DatabasePath)
	return err == nil
}

func (a *AIDEConfig) Restore() error {
	backupDb := filepath.Join(a.BackupDir, "aide.db.gz")
	if _, err := os.Stat(backupDb); os.IsNotExist(err) {
		return fmt.Errorf("no backup found at %s", backupDb)
	}

	input, err := os.ReadFile(backupDb)
	if err != nil {
		return err
	}

	return os.WriteFile(a.DatabasePath, input, 0644)
}
