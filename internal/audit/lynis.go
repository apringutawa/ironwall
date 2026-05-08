package audit

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type LynisScanner struct {
	ReportPath string
	LogPath    string
}

type LynisResult struct {
	HardeningIndex   float64
	Warnings         int
	Suggestions      int
	TestsPerformed   int
	Status           string
	Recommendations  []string
	InstallAvailable bool
}

func NewLynisScanner() *LynisScanner {
	return &LynisScanner{
		ReportPath: "/var/log/lynis-report.dat",
		LogPath:    "/var/log/lynis.log",
	}
}

func (l *LynisScanner) IsInstalled() bool {
	_, err := exec.LookPath("lynis")
	return err == nil
}

func (l *LynisScanner) Install() error {
	fmt.Println("Installing Lynis...")
	cmd := exec.Command("apt-get", "install", "-y", "-qq", "lynis")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install lynis: %w\n%s", err, string(output))
	}
	fmt.Println("  ✓ Lynis installed")
	return nil
}

func (l *LynisScanner) EnsureInstalled() error {
	if l.IsInstalled() {
		return nil
	}
	return l.Install()
}

func (l *LynisScanner) RunAudit() (*LynisResult, error) {
	if err := l.EnsureInstalled(); err != nil {
		return nil, err
	}

	fmt.Println("Running Lynis security audit...")
	cmd := exec.Command("lynis", "audit", "system", "--quick", "--no-colors")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("lynis audit failed: %w", err)
	}

	logFile := "/var/log/lynis.log"
	if _, err := os.Stat(logFile); err == nil {
		return l.ParseReport(logFile)
	}

	return &LynisResult{
		Status:          "completed",
		TestsPerformed:  len(strings.Split(string(output), "\n")),
		InstallAvailable: true,
	}, nil
}

func (l *LynisScanner) ParseReport(logPath string) (*LynisResult, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := &LynisResult{
		InstallAvailable: true,
		Status:          "unknown",
	}

	reIndex := regexp.MustCompile(`hardening_index\s*[=:]\s*([0-9.]+)`)
	reWarnings := regexp.MustCompile(`warnings\s*[=:]\s*(\d+)`)
	reSuggestions := regexp.MustCompile(`suggestions\s*[=:]\s*(\d+)`)
	reTests := regexp.MustCompile(`tests_performed\s*[=:]\s*(\d+)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if matches := reIndex.FindStringSubmatch(line); len(matches) > 1 {
			result.HardeningIndex, _ = strconv.ParseFloat(matches[1], 64)
		}
		if matches := reWarnings.FindStringSubmatch(line); len(matches) > 1 {
			result.Warnings, _ = strconv.Atoi(matches[1])
		}
		if matches := reSuggestions.FindStringSubmatch(line); len(matches) > 1 {
			result.Suggestions, _ = strconv.Atoi(matches[1])
		}
		if matches := reTests.FindStringSubmatch(line); len(matches) > 1 {
			result.TestsPerformed, _ = strconv.Atoi(matches[1])
		}
	}

	if result.HardeningIndex >= 60 {
		result.Status = "good"
	} else if result.HardeningIndex >= 40 {
		result.Status = "fair"
	} else {
		result.Status = "poor"
	}

	return result, scanner.Err()
}

func (l *LynisScanner) GetHardeningScore() float64 {
	result, err := l.ParseReport(l.LogPath)
	if err != nil || result == nil {
		return 0
	}
	return result.HardeningIndex
}

func (l *LynisScanner) GetRecommendations() ([]string, error) {
	file, err := os.Open(l.LogPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var recommendations []string
	re := regexp.MustCompile(`(?i)suggestion|warning|recommendation`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) && !strings.HasPrefix(line, "#") {
			recommendations = append(recommendations, strings.TrimSpace(line))
		}
	}

	return recommendations, scanner.Err()
}
