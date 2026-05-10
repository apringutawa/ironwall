<div align="center">
  <br/>
  <pre style="font-family: monospace; font-size: 1.2em; line-height: 1.4; color: #00ff88;">
╔═══════════════════════════════════════════════════════════════╗
║                    IRONWALL Hardening System                  ║
║                 by Kubu Raya CSIRT - 2026                     ║
╚═══════════════════════════════════════════════════════════════╝
  </pre>
  <br/>
  <p><strong>Automated Linux System Hardening &amp; Threat Prevention Toolkit</strong></p>
  <p>
    <a href="https://github.com/apringutawa/ironwall/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue?style=flat-square" alt="License"/></a>
    <a href="https://go.dev/"><img src="https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go" alt="Go"/></a>
    <a href="https://www.python.org/"><img src="https://img.shields.io/badge/python-3.11+-3776AB?style=flat-square&logo=python" alt="Python"/></a>
    <a href="https://nextjs.org/"><img src="https://img.shields.io/badge/next.js-15-000000?style=flat-square&logo=nextdotjs" alt="Next.js"/></a>
    <img src="https://img.shields.io/badge/platform-linux-important?style=flat-square&logo=linux" alt="Linux"/>
    <img src="https://img.shields.io/badge/stability-stable-success?style=flat-square" alt="Stable"/>
  </p>
  <br/>
</div>

---

## Overview

IronWall is a comprehensive, one-shot security solution for Linux servers. It automates the entire hardening process from SSH protection to kernel tuning, file integrity monitoring, malware detection, and real-time alerting.

> **One command to harden your server:**
> ```bash
> curl -sSL https://raw.githubusercontent.com/apringutawa/ironwall/main/scripts/install.sh | sudo bash
> ```

---

## Features

### 🔐 SSH Hardening
- Disable root login & password authentication
- Enforce ed25519/RSA public key authentication
- Optional port change with automatic SELinux/AppArmor adjustment
- Automatic backup before any modification

### 🛡️ Firewall Automation
- Auto-detects `nftables` / `iptables` / `firewalld`
- Default deny ingress policy
- Allow ports: `22`, `80`, `443`, `8080`
- IP block/unblock, port management

### 🚫 Fail2Ban Integration
- Automatic brute-force protection for SSH, HTTP, HTTPS
- Custom jail configurations
- Real-time ban/unban monitoring

### 🦠 Malware Scanner
- ClamAV integration for file scanning
- Cryptominer detection (xmrig, kinsing, cpuminer, etc.)
- Webshell analysis (PHP eval, base64, encoded payloads)
- Process monitoring & auto-quarantine

### ⏰ Cron Protection
- Real-time monitoring of `/etc/crontab`, `/etc/cron.d/`, user crontabs
- Suspicious pattern detection: `curl|bash`, `wget|sh`, reverse shells
- Alert generation with full context

### 📁 AIDE File Integrity
- File integrity database with SHA-256/SHA-512 hashing
- Daily automated integrity checks via systemd timer
- Fallback when `auditd` is unavailable (containers, VPS)
- Baseline snapshots & change detection alerts

### ⚙️ Kernel Sysctl Hardening (27 parameters)
- ASLR (`kernel.randomize_va_space = 2`)
- IP spoofing protection (strict rp_filter)
- SYN flood protection (syncookies, reduced retries)
- ICMP redirect ignore, source-routed packet rejection
- TCP stack hardening (timestamps off, FIN timeout reduced)
- Magic SysRq disabled, core dumps restricted

### 👤 User Account Hardening
- PAM password policy: min length 12, complexity (upper, lower, digit, other)
- Faillock: 5 failed attempts → 15-minute lockout
- Sudo I/O audit logging (input/output logged to `/var/log/sudo.log`)
- Password aging: max 90 days, min 7 days, warn 14 days
- UID 0 check (non-root users with UID 0 flagged)
- Empty password check

### 🔧 Service Minimization
- **Dangerous services auto-disabled:** telnet, rsh, rlogin, rexec, ftp, tftp, nfs, samba, snmpd, bind, slapd, cups, avahi-daemon, xinetd, rpcbind
- Systemd masking prevents re-activation
- Dry-run mode for preview

### 📊 Security Audit (Lynis)
- Comprehensive system audit with 300+ tests
- Hardening index score (0–100)
- Actionable recommendations & warnings
- Auto-installs Lynis if missing

---

## Quick Start

### One-Click Installation

```bash
curl -sSL https://raw.githubusercontent.com/apringutawa/ironwall/main/scripts/install.sh | sudo bash
```

Or via `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/apringutawa/ironwall/main/scripts/install.sh | sudo bash
```

The installer will:
1. Detect your OS (Ubuntu/Debian/CentOS/RHEL)
2. Install all dependencies (ClamAV, Lynis, AIDE, nftables, Fail2Ban, etc.)
3. Auto-install Go & build the `ironwall` CLI
4. Back up existing system configuration
5. Apply **all** hardening modules (SSH, kernel, accounts, AIDE, firewall, services)
6. Start monitoring services (Fail2Ban, AIDE, systemd timer)
7. Install the IronWall systemd service

### Manual Installation (Development)

```bash
git clone https://github.com/apringutawa/ironwall.git
cd ironwall

# Backend
cd backend && pip install -r requirements.txt && cd ..

# Frontend
cd frontend && npm install && cd ..

# CLI
go build -o ironwall cmd/ironwall/main.go
sudo mv ironwall /usr/local/bin/

# Run installer
sudo ironwall install
```

---

## CLI Reference

### System

| Command | Description |
|---------|-------------|
| `ironwall status` | Show system health & module status |
| `ironwall status --json` | JSON-formatted status |
| `ironwall scan` | Run malware scan (ClamAV + processes) |
| `ironwall logs` | View security logs |
| `ironwall logs --follow` | Tail logs in real-time |
| `ironwall logs --lines 100` | Show last N lines |

### Protection

| Command | Description |
|---------|-------------|
| `ironwall protect` | Apply all security protections |
| `ironwall protect --dry-run` | Preview changes without applying |
| `ironwall unprotect --force` | Disable all protections |
| `ironwall rollback` | Restore system from backup |

### Firewall

| Command | Description |
|---------|-------------|
| `ironwall firewall status` | Show firewall status |
| `ironwall firewall block <IP>` | Block an IP address |
| `ironwall firewall unblock <IP>` | Unblock an IP address |
| `ironwall firewall allow <port> <tcp\|udp>` | Allow a port |

### Kernel Hardening

| Command | Description |
|---------|-------------|
| `ironwall sysctl apply` | Apply all kernel hardening parameters |
| `ironwall sysctl apply --dry-run` | Preview kernel changes |
| `ironwall sysctl verify` | Verify all parameters are applied |
| `ironwall sysctl restore` | Restore previous sysctl configuration |

### Account Hardening

| Command | Description |
|---------|-------------|
| `ironwall accounts harden` | PAM policy, faillock, sudo audit, password aging |
| `ironwall accounts list` | List all human users |
| `ironwall accounts check` | Security audit (UID 0, empty passwords, SSH config) |
| `ironwall accounts restore` | Restore previous account configuration |

### AIDE Integrity

| Command | Description |
|---------|-------------|
| `ironwall aide init` | Initialize AIDE database |
| `ironwall aide check` | Run integrity check |
| `ironwall aide update` | Update AIDE database baseline |
| `ironwall aide enable` | Enable daily automated integrity checks |

### Service Minimization

| Command | Description |
|---------|-------------|
| `ironwall services list` | Classify all services (dangerous/essential/optional) |
| `ironwall services minimize` | Disable & mask all dangerous services |
| `ironwall services minimize --dry-run` | Preview service changes |
| `ironwall services scan` | Scan for insecure running services |

### Security Audit

| Command | Description |
|---------|-------------|
| `ironwall audit` | Run full Lynis security audit |
| `ironwall audit install` | Install Lynis |
| `ironwall audit score` | Show last hardening index score |

---

## Web Dashboard

Access the dashboard at `http://your-server-ip:8080`

Features:
- Real-time system status (CPU, memory, uptime)
- Security module management (enable/disable)
- Security events viewer with severity filters
- Firewall rule management (block/unblock IPs)
- Hardening module controls (apply kernel, accounts, AIDE, services)
- Service classification & dangerous service detection
- Alert configuration (Telegram, Discord, Slack, Email)

### API Documentation

Once the backend is running:
- **Swagger UI**: `http://localhost:8001/docs`
- **ReDoc**: `http://localhost:8001/redoc`

### Default Ports

| Port | Service | Purpose |
|------|---------|---------|
| `22` | SSH | Secure shell access |
| `80` | HTTP | Web server (optional) |
| `443` | HTTPS | Web server (optional) |
| `8080` | Dashboard | IronWall web UI |
| `8001` | API | IronWall REST API |

---

## Architecture

```
ironwall/
├── cmd/ironwall/              # Go CLI entry point
│   ├── main.go
│   └── cmd/                   # 13 Cobra commands
├── internal/                  # Go security modules
│   ├── audit/                 # Lynis integration
│   ├── hardening/             # SSH, sysctl, accounts, AIDE, services
│   ├── firewall/              # nftables/iptables/firewalld
│   ├── scanner/               # ClamAV, cryptominer, webshell
│   ├── monitor/               # Cron job monitoring
│   └── backup/                # Backup & restore
├── backend/                   # FastAPI + SQLAlchemy + SQLite
│   └── app/
│       ├── api/               # 6 API route modules
│       ├── models/            # Database ORM + Pydantic schemas
│       └── core/              # Config, database, alerts
├── frontend/                  # Next.js 15 + React 19 + TypeScript
│   └── src/
│       ├── app/               # 7 pages (App Router)
│       ├── components/        # React components
│       └── lib/               # API client (axios)
├── configs/                   # Default security configs
├── scripts/                   # Installer script
└── docs/                      # Documentation
```

### Tech Stack

| Layer | Technology |
|-------|-----------|
| CLI | Go 1.21+ / Cobra |
| Backend | FastAPI / SQLAlchemy / SQLite |
| Frontend | Next.js 15 / React 19 / TypeScript / Tailwind CSS 4 |
| Security Tools | ClamAV, Fail2Ban, Lynis, AIDE, nftables, auditd |

---

## Configuration

### Backend (`backend/.env`)

```env
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id
DASHBOARD_PORT=8080
```

### Frontend (`frontend/.env.local`)

```env
NEXT_PUBLIC_API_URL=http://your-server-ip:8001
```

### File Locations

| Resource | Path |
|----------|------|
| Backups | `/var/backups/ironwall/` |
| Database | `backend/ironwall.db` |
| Firewall Rules | `/etc/nftables.conf` |
| Kernel Parameters | `/etc/sysctl.d/99-ironwall.conf` |
| Sudo Audit | `/etc/sudoers.d/ironwall` |
| AIDE Config | `/etc/aide/aide.conf` |
| AIDE Database | `/var/lib/aide/aide.db.gz` |
| Logs | `/var/log/ironwall/` |
| Sudo Log | `/var/log/sudo.log` |

---

## Supported Operating Systems

| OS | Status |
|----|--------|
| Ubuntu 20.04 LTS | ✅ Full support |
| Ubuntu 22.04 LTS | ✅ Full support |
| Ubuntu 24.04 LTS | ✅ Full support |
| Debian 11 | ✅ Full support |
| Debian 12 | ✅ Full support |
| CentOS 7 | ⏳ Planned |
| AlmaLinux 9 | ⏳ Planned |
| Rocky Linux 9 | ⏳ Planned |

---

## Performance

| Metric | Typical |
|--------|---------|
| CPU Overhead | < 5% |
| RAM Usage | < 300 MB (incl. dashboard) |
| Alert Latency | < 3 seconds |
| Installation Time | 2–5 minutes |
| Boot Impact | Minimal (lazy-loaded services) |

---

## Development

### Prerequisites

- Go 1.21+
- Python 3.11+
- Node.js 20+
- Linux (Ubuntu/Debian recommended)

### Setup

```bash
# Clone
git clone https://github.com/apringutawa/ironwall.git
cd ironwall

# Backend
cd backend
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8001

# Frontend (new terminal)
cd frontend
npm install
npm run dev

# CLI
go build -o ironwall cmd/ironwall/main.go
./ironwall status
```

### Testing

```bash
# Backend
cd backend && pytest

# Frontend
cd frontend && npm test

# CLI
go test ./...
```

---

## Roadmap

### Phase 1 — MVP ✅

| Feature | Status |
|---------|--------|
| CLI with Cobra commands | ✅ |
| SSH hardening | ✅ |
| Firewall automation | ✅ |
| Malware scanner | ✅ |
| Cron monitoring | ✅ |
| File integrity (chattr) | ✅ |
| Web dashboard | ✅ |
| Alert system | ✅ |
| **Lynis security audit** | ✅ |
| **Kernel sysctl hardening** | ✅ |
| **User account hardening** | ✅ |
| **AIDE file integrity** | ✅ |
| **Service minimization** | ✅ |

### Phase 2 — In Progress 🚧

| Feature | Status |
|---------|--------|
| WAF integration (ModSecurity) | ⏳ |
| Real-time IDS | ⏳ |
| eBPF monitoring | ⏳ |
| Docker container protection | ⏳ |
| Auto security updates | ⏳ |

### Phase 3 — Future 🔮

| Feature | Status |
|---------|--------|
| AI anomaly detection | 📋 |
| Multi-server management | 📋 |
| Compliance reporting (CIS, NIST) | 📋 |
| Kubernetes admission control | 📋 |
| Zero trust mode | 📋 |

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Commit changes (`git commit -m 'Add my feature'`)
4. Push to branch (`git push origin feature/my-feature`)
5. Open a Pull Request

### Guidelines

- Follow existing code structure and patterns
- Add tests for new functionality
- Update `AGENTS.md` when adding modules
- Test on supported OS versions
- Use `--dry-run` flag for destructive operations

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

## Contact

- **Repository**: [github.com/apringutawa/ironwall](https://github.com/apringutawa/ironwall)
- **Issues**: [GitHub Issues](https://github.com/apringutawa/ironwall/issues)
- **Organization**: [Kubu Raya CSIRT](https://csirt.kuburayakab.go.id)

---

<div align="center">
  <br/>
  <p>
    <strong>IRONWALL Hardening System</strong><br/>
    by <strong>Kubu Raya CSIRT</strong> — 2026
  </p>
  <p>
    <sub>Built with ❤️ for the security community</sub>
  </p>
  <br/>
  <pre style="font-family: monospace; color: #00ff88;">
╔═══════════════════════════════════════════════════════════════╗
║              One Command to Harden Them All                   ║
╚═══════════════════════════════════════════════════════════════╝
  </pre>
  <br/>
</div>
