package backup

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type BackupManager struct {
	BackupDir string
}

func NewBackupManager() *BackupManager {
	return &BackupManager{
		BackupDir: "/var/backups/ironwall",
	}
}

func (b *BackupManager) CreateBackup() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(b.BackupDir, timestamp)

	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return "", err
	}

	files := []string{
		"/etc/ssh/sshd_config",
		"/etc/crontab",
		"/etc/passwd",
		"/etc/shadow",
		"/etc/group",
	}

	dirs := []string{
		"/etc/ssh",
		"/etc/cron.d",
		"/var/spool/cron",
	}

	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}
		if err := b.copyFile(file, backupPath); err != nil {
			return "", err
		}
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		if err := b.copyDir(dir, backupPath); err != nil {
			return "", err
		}
	}

	return backupPath, nil
}

func (b *BackupManager) copyFile(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	dstPath := filepath.Join(dst, filepath.Base(src))
	output, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}

func (b *BackupManager) copyDir(src, dst string) error {
	dirName := filepath.Base(src)
	dstPath := filepath.Join(dst, dirName)

	cmd := exec.Command("cp", "-r", src, dstPath)
	return cmd.Run()
}

func (b *BackupManager) ListBackups() ([]string, error) {
	var backups []string

	entries, err := os.ReadDir(b.BackupDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			backups = append(backups, entry.Name())
		}
	}

	return backups, nil
}

func (b *BackupManager) Restore(backupID string) error {
	backupPath := filepath.Join(b.BackupDir, backupID)

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup not found: %s", backupID)
	}

	files := map[string]string{
		"sshd_config": "/etc/ssh/sshd_config",
		"crontab":     "/etc/crontab",
		"passwd":      "/etc/passwd",
		"shadow":      "/etc/shadow",
		"group":       "/etc/group",
	}

	for src, dst := range files {
		srcPath := filepath.Join(backupPath, src)
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			continue
		}

		cmd := exec.Command("cp", srcPath, dst)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// Restore directories
	dirs := []string{"ssh", "cron.d"}
	for _, dir := range dirs {
		srcPath := filepath.Join(backupPath, dir)
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			continue
		}

		switch dir {
		case "ssh":
			cmd := exec.Command("cp", "-r", srcPath+"/.", "/etc/ssh/")
			cmd.Run()
		case "cron.d":
			cmd := exec.Command("cp", "-r", srcPath+"/.", "/etc/cron.d/")
			cmd.Run()
		}
	}

	return nil
}

func (b *BackupManager) UnlockFiles() error {
	files := []string{
		"/etc/passwd",
		"/etc/shadow",
		"/etc/group",
		"/etc/ssh/sshd_config",
	}

	for _, file := range files {
		cmd := exec.Command("chattr", "-i", file)
		cmd.Run()
	}

	return nil
}

func (b *BackupManager) LockFiles() error {
	files := []string{
		"/etc/passwd",
		"/etc/shadow",
		"/etc/group",
	}

	for _, file := range files {
		cmd := exec.Command("chattr", "+i", file)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}
