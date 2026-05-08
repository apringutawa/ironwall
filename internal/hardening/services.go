package hardening

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type ServiceManager struct {
	EssentialServices []string
	DangerousServices []string
	OptionalServices  []string
}

type ServiceInfo struct {
	Name   string
	Status string
	Action string
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		EssentialServices: []string{
			"ssh", "sshd",
			"fail2ban",
			"systemd-journald",
			"systemd-logind",
			"cron", "crond",
			"rsyslog", "syslog-ng",
			"networking", "network",
			"systemd-timesyncd", "ntp", "chronyd",
			"dbus",
			"udev", "systemd-udevd",
		},
		DangerousServices: []string{
			"telnet", "telnetd",
			"rsh", "rsh-server",
			"rlogin", "rlogin-server",
			"rexec", "rexec-server",
			"vsftpd", "proftpd", "pure-ftpd",
			"tftpd", "tftp",
			"nfs-server", "nfs-kernel-server",
			"smbd", "nmbd", "samba",
			"snmpd",
			"bind9", "named",
			"slapd",
			"cups", "cups-browsed",
			"avahi-daemon",
			"xinetd",
			"rpcbind",
		},
		OptionalServices: []string{
			"apache2", "httpd", "nginx",
			"mysql", "mariadb", "postgresql",
			"redis-server", "redis",
			"docker", "containerd",
			"postfix", "sendmail", "exim4",
		},
	}
}

func (s *ServiceManager) IsSystemd() bool {
	_, err := os.Stat("/run/systemd/system")
	return err == nil
}

func (s *ServiceManager) detectServiceManager() string {
	if s.IsSystemd() {
		return "systemd"
	}
	if _, err := exec.LookPath("service"); err == nil {
		return "sysvinit"
	}
	return "unknown"
}

func (s *ServiceManager) ListEnabled() ([]ServiceInfo, error) {
	var services []ServiceInfo
	manager := s.detectServiceManager()

	switch manager {
	case "systemd":
		return s.listSystemdServices()
	case "sysvinit":
		return s.listSysvServices()
	}

	return services, nil
}

func (s *ServiceManager) listSystemdServices() ([]ServiceInfo, error) {
	var services []ServiceInfo

	cmd := exec.Command("systemctl", "list-unit-files", "--type=service", "--no-pager", "--no-legend")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^(\S+)\.service\s+(\S+)`)

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue
		}

		name := matches[1]
		status := matches[2]

		action := ""
		if s.isDangerous(name) {
			action = "disable"
		} else if s.isEssential(name) {
			action = "keep"
		} else if status == "enabled" {
			action = "review"
		}

		if action != "" {
			services = append(services, ServiceInfo{
				Name:   name,
				Status: status,
				Action: action,
			})
		}
	}

	return services, nil
}

func (s *ServiceManager) listSysvServices() ([]ServiceInfo, error) {
	var services []ServiceInfo

	cmd := exec.Command("service", "--status-all")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`\[(\+|\-)\]\s+(\S+)`)

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) < 3 {
			continue
		}

		status := matches[1]
		name := matches[2]

		statusStr := "inactive"
		if status == "+" {
			statusStr = "active"
		}

		action := ""
		if s.isDangerous(name) {
			action = "disable"
		} else if s.isEssential(name) {
			action = "keep"
		}

		if action != "" {
			services = append(services, ServiceInfo{
				Name:   name,
				Status: statusStr,
				Action: action,
			})
		}
	}

	return services, nil
}

func (s *ServiceManager) isEssential(name string) bool {
	for _, essential := range s.EssentialServices {
		if name == essential {
			return true
		}
	}
	return false
}

func (s *ServiceManager) isDangerous(name string) bool {
	for _, dangerous := range s.DangerousServices {
		if name == dangerous {
			return true
		}
	}
	return false
}

func (s *ServiceManager) Apply(dryRun bool) ([]string, error) {
	fmt.Println("Minimizing services...")

	if !s.IsSystemd() {
		return nil, fmt.Errorf("service minimization requires systemd")
	}

	var disabled []string

	dangerousServices := s.DangerousServices

	for _, svc := range dangerousServices {
		cmd := exec.Command("systemctl", "is-enabled", svc)
		if err := cmd.Run(); err != nil {
			continue
		}

		if dryRun {
			fmt.Printf("  [DRY RUN] Would disable: %s\n", svc)
			disabled = append(disabled, svc)
			continue
		}

		fmt.Printf("  Disabling dangerous service: %s\n", svc)
		exec.Command("systemctl", "stop", svc).Run()
		exec.Command("systemctl", "disable", svc).Run()
		exec.Command("systemctl", "mask", svc).Run()
		disabled = append(disabled, svc)
	}

	if !dryRun {
		exec.Command("systemctl", "daemon-reload").Run()
	}

	fmt.Printf("  ✓ %d dangerous services disabled\n", len(disabled))
	return disabled, nil
}

func (s *ServiceManager) Restore(services []string) error {
	for _, svc := range services {
		exec.Command("systemctl", "unmask", svc).Run()
		exec.Command("systemctl", "enable", svc).Run()
	}

	exec.Command("systemctl", "daemon-reload").Run()
	return nil
}

func (s *ServiceManager) GetStatus() (map[string]string, error) {
	status := make(map[string]string)

	if !s.IsSystemd() {
		return status, nil
	}

	for _, svc := range s.DangerousServices {
		cmd := exec.Command("systemctl", "is-active", svc)
		output, err := cmd.Output()
		if err == nil {
			status[svc] = strings.TrimSpace(string(output))
		}
	}

	return status, nil
}

func (s *ServiceManager) ScanForInsecureServices() ([]string, error) {
	cmd := exec.Command("ss", "-tlnp")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var insecure []string
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		for _, dangerous := range s.DangerousServices {
			if strings.Contains(line, dangerous) {
				insecure = append(insecure, fmt.Sprintf("%s: %s", dangerous, strings.TrimSpace(line)))
			}
		}
	}

	return insecure, nil
}
