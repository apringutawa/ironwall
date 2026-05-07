package scanner

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type MalwareScanner struct {
	ClamAVPath string
	ScanPaths  []string
}

func NewMalwareScanner() *MalwareScanner {
	return &MalwareScanner{
		ClamAVPath: "/var/lib/clamav",
		ScanPaths:  []string{"/tmp", "/var/www", "/home", "/dev/shm"},
	}
}

func (s *MalwareScanner) UpdateSignatures() error {
	fmt.Println("Updating ClamAV signatures...")
	cmd := exec.Command("freshclam")
	return cmd.Run()
}

func (s *MalwareScanner) Scan() (int, error) {
	fmt.Println("Scanning for malware...")

	totalMalware := 0

	for _, path := range s.ScanPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		fmt.Printf("  Scanning %s...\n", path)

		cmd := exec.Command("clamdscan", "--no-summary", path)
		output, err := cmd.CombinedOutput()

		if err != nil {
			malwareCount := s.parseScanOutput(string(output))
			totalMalware += malwareCount
		}
	}

	return totalMalware, nil
}

func (s *MalwareScanner) parseScanOutput(output string) int {
	count := 0
	lines := strings.Split(output, "\n")
	re := regexp.MustCompile(`(?i)(FOUND|Infected)`)
	for _, line := range lines {
		if re.MatchString(line) {
			count++
		}
	}
	return count
}

func (s *MalwareScanner) ScanProcesses() ([]string, error) {
	fmt.Println("Scanning running processes...")

	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	suspiciousPatterns := []string{
		"xmrig",
		"kinsing",
		"cryptominer",
		"minerd",
		"cpuminer",
		"stratum",
		"reverse.*shell",
		"nc.*-e",
		"bash.*-i",
		"python.*socket",
	}

	var suspicious []string
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		for _, pattern := range suspiciousPatterns {
			re := regexp.MustCompile(pattern)
			if re.MatchString(line) {
				suspicious = append(suspicious, line)
				break
			}
		}
	}

	return suspicious, nil
}

func (s *MalwareScanner) KillProcess(pid int) error {
	cmd := exec.Command("kill", "-9", fmt.Sprintf("%d", pid))
	return cmd.Run()
}

func (s *MalwareScanner) Quarantine(path string) error {
	quarantineDir := "/var/quarantine/ironwall"
	if err := os.MkdirAll(quarantineDir, 0700); err != nil {
		return err
	}

	cmd := exec.Command("mv", path, quarantineDir)
	return cmd.Run()
}

func (s *MalwareScanner) ScanWebshell() ([]string, error) {
	fmt.Println("Scanning for webshells...")

	phpPaths := []string{"/var/www", "/home", "/tmp"}
	var webshells []string

	patterns := []string{
		`eval\s*\(\s*base64_decode`,
		`system\s*\(\s*\$_GET`,
		`shell_exec`,
		`exec\s*\(\s*\$_`,
		`preg_replace.*\xe2\x80\x8b`,
		`assert\s*\(`,
	}

	for _, path := range phpPaths {
		cmd := exec.Command("find", path, "-name", "*.php", "-type", "f")
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		files := strings.Split(string(output), "\n")
		for _, file := range files {
			if file == "" {
				continue
			}

			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			for _, pattern := range patterns {
				re := regexp.MustCompile(pattern)
				if re.Match(content) {
					webshells = append(webshells, file)
					break
				}
			}
		}
	}

	return webshells, nil
}
