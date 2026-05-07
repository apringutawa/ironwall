package monitor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type CronMonitor struct {
	CronPaths  []string
	BackupDir  string
	AlertChan  chan string
}

func NewCronMonitor() *CronMonitor {
	return &CronMonitor{
		CronPaths: []string{
			"/etc/crontab",
			"/etc/cron.d",
			"/var/spool/cron",
		},
		BackupDir:  "/var/backups/ironwall/cron",
		AlertChan:  make(chan string, 100),
	}
}

func (c *CronMonitor) Backup() error {
	if err := os.MkdirAll(c.BackupDir, 0700); err != nil {
		return err
	}

	for _, path := range c.CronPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		backupPath := filepath.Join(c.BackupDir, filepath.Base(path)+"."+time.Now().Format("20060102150405"))

		cmd := exec.Command("cp", "-r", path, backupPath)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (c *CronMonitor) DetectSuspicious(content string) []string {
	suspiciousPatterns := []string{
		`curl.*\|.*bash`,
		`wget.*\|.*sh`,
		`curl.*\|.*sh`,
		`wget.*\|.*bash`,
		`python.*-c.*import`,
		`/dev/tcp`,
		`xmrig`,
		`kinsing`,
		`base64.*-d`,
		`eval\s*\(`,
		`chmod.*\+x.*http`,
		`/tmp/.*`,
		`crontab.*-r`,
	}

	var detected []string
	for _, pattern := range suspiciousPatterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(content) {
			detected = append(detected, pattern)
		}
	}

	return detected
}

func (c *CronMonitor) ReadCronFiles() (map[string]string, error) {
	cronContents := make(map[string]string)

	for _, path := range c.CronPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		fileInfo, err := os.Stat(path)
		if err != nil {
			continue
		}

		if fileInfo.IsDir() {
			files, err := filepath.Glob(filepath.Join(path, "*"))
			if err != nil {
				continue
			}

			for _, file := range files {
				content, err := os.ReadFile(file)
				if err != nil {
					continue
				}
				cronContents[file] = string(content)
			}
		} else {
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			cronContents[path] = string(content)
		}
	}

	return cronContents, nil
}

func (c *CronMonitor) Scan() ([]string, error) {
	fmt.Println("Scanning cron jobs...")

	cronContents, err := c.ReadCronFiles()
	if err != nil {
		return nil, err
	}

	var suspicious []string
	for path, content := range cronContents {
		detected := c.DetectSuspicious(content)
		if len(detected) > 0 {
			suspicious = append(suspicious, fmt.Sprintf("%s: %s", path, strings.Join(detected, ", ")))
		}
	}

	return suspicious, nil
}

func (c *CronMonitor) Start() error {
	fmt.Println("Starting cron monitor...")

	if err := c.Backup(); err != nil {
		return err
	}

	go c.monitor()

	return nil
}

func (c *CronMonitor) monitor() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		cronContents, err := c.ReadCronFiles()
		if err != nil {
			continue
		}

		for path, content := range cronContents {
			detected := c.DetectSuspicious(content)
			if len(detected) > 0 {
				alert := fmt.Sprintf("Suspicious cron detected in %s: %s", path, strings.Join(detected, ", "))
				c.AlertChan <- alert
			}
		}
	}
}

func (c *CronMonitor) Stop() {
	close(c.AlertChan)
}
