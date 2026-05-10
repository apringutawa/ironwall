<div align="center">
  <br/>
  <pre style="font-family: monospace; font-size: 1.2em; line-height: 1.4; color: #00ff88;">
╔═══════════════════════════════════════════════════════════════╗
║                    IRONWALL Hardening System                  ║
║                 by KUBU RAYA CSIRT - 2026                     ║
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

## Quick Install

```bash
curl -sSL https://raw.githubusercontent.com/apringutawa/ironwall/main/scripts/install.sh | sudo bash
```

---

## Screenshots

> 📸 **Screenshots**: Add your installation and dashboard screenshots to `docs/screenshots/` and update this section.  
> Run the tool first, then capture images of: installation process, web dashboard, Lynis audit output, and hardening dashboard.

---

## Features

| # | Module | Description |
|---|--------|-------------|
| 1 | 🔐 **SSH Hardening** | Disable root login, enforce key auth, port change |
| 2 | 🛡️ **Firewall Automation** | nftables/iptables/firewalld with default deny |
| 3 | 🚫 **Fail2Ban** | Brute-force protection for SSH, HTTP, HTTPS |
| 4 | 🦠 **Malware Scanner** | ClamAV, cryptominer, webshell, reverse shell detection |
| 5 | ⏰ **Cron Protection** | Real-time cron monitoring & suspicious pattern detection |
| 6 | 📁 **AIDE Integrity** | File integrity DB with daily automated checks |
| 7 | ⚙️ **Kernel Hardening** | 27 sysctl params: ASLR, anti-spoof, SYN flood, TCP stack |
| 8 | 👤 **Account Hardening** | PAM policy, faillock, sudo audit, password aging |
| 9 | 🔧 **Service Minimization** | Disable & mask telnet, ftp, nfs, samba, cups, 15+ services |
| 10 | 📊 **Security Audit** | Lynis: 300+ tests, hardening index score, recommendations |

---

## CLI Commands

### System
```
ironwall status                  Show system health
ironwall scan                    Run malware scan
ironwall logs                    View security logs
ironwall logs --follow           Tail logs in real-time
```

### Protection
```
ironwall protect                 Apply all protections
ironwall protect --dry-run       Preview changes
ironwall unprotect --force       Disable all protections
ironwall rollback                Restore from backup
```

### Firewall
```
ironwall firewall status         Show firewall status
ironwall firewall block <IP>     Block an IP
ironwall firewall unblock <IP>   Unblock an IP
```

### Hardening Modules
```
ironwall sysctl apply            Apply kernel hardening
ironwall sysctl verify           Verify kernel params
ironwall accounts harden         PAM, faillock, sudo audit
ironwall accounts list           List system users
ironwall accounts check          Security audit
ironwall aide init               Init AIDE database
ironwall aide check              Run integrity check
ironwall aide enable             Enable daily checks
ironwall services list           List services by class
ironwall services minimize       Disable dangerous services
ironwall audit                   Run Lynis security audit
ironwall audit score             Show hardening score
```

---

## Web Dashboard

Access at `http://your-server-ip:8080`

- System health monitoring
- Module management (enable/disable)
- Security events viewer
- Firewall rule management
- Hardening controls (kernel, accounts, AIDE, services)
- Service classification & dangerous service detection
- Alert configuration (Telegram, Discord, Slack, Email)

### API Docs
- **Swagger UI**: `http://localhost:8001/docs`
- **ReDoc**: `http://localhost:8001/redoc`

---

## Architecture

```
ironwall/
├── cmd/ironwall/              # Go CLI (Cobra)
│   └── cmd/                   # 13 subcommands
├── internal/                  # Security modules
│   ├── audit/                 # Lynis integration
│   ├── hardening/             # SSH, sysctl, accounts, AIDE, services
│   ├── firewall/              # nftables/iptables/firewalld
│   ├── scanner/               # Malware detection
│   ├── monitor/               # Cron monitoring
│   └── backup/                # Backup & restore
├── backend/                   # FastAPI + SQLite
│   └── app/
│       ├── api/               # 6 route modules
│       ├── models/            # ORM + schemas
│       └── core/              # Config, alerts
├── frontend/                  # Next.js 15 + React 19
├── configs/                   # Security configs
├── scripts/                   # Installer
└── docs/screenshots/          # Screenshots
```

### Tech Stack

| Layer | Stack |
|-------|-------|
| CLI | Go 1.21+ / Cobra |
| Backend | FastAPI / SQLAlchemy / SQLite |
| Frontend | Next.js 15 / React 19 / TypeScript / Tailwind CSS 4 |
| Security | ClamAV, Fail2Ban, Lynis, AIDE, nftables, auditd |

---

## Supported OS

| OS | Status |
|----|--------|
| Ubuntu 20.04 / 22.04 / 24.04 | ✅ |
| Debian 11 / 12 | ✅ |
| CentOS 7 / AlmaLinux / Rocky Linux | ⏳ |

---

## Configuration

### Backend (`backend/.env`)
```env
TELEGRAM_BOT_TOKEN=your_token
TELEGRAM_CHAT_ID=your_chat_id
DASHBOARD_PORT=8080
```

### Frontend (`frontend/.env.local`)
```env
NEXT_PUBLIC_API_URL=http://localhost:8001
```

---

## Collaboration

**IRONWALL** dikembangkan oleh **KUBU RAYA CSIRT** untuk kepentingan keamanan siber di lingkungan Pemerintah Kabupaten Kubu Raya dan masyarakat luas.

Kami mengundang kontribusi dari:
- **Security researchers** — laporkan celah, usulkan fitur keamanan baru
- **System administrators** — test di lingkungan produksi, beri feedback
- **Developers** — bantu pengembangan modul, fix bug, perbaiki dokumentasi

### Cara Berkontribusi
1. Fork repository ini
2. Buat branch: `git checkout -b feature/nama-fitur`
3. Commit: `git commit -m 'Add fitur baru'`
4. Push: `git push origin feature/nama-fitur`
5. Open Pull Request

### Kontak
- **GitHub Issues**: [github.com/apringutawa/ironwall/issues](https://github.com/apringutawa/ironwall/issues)
- **Organisasi**: [Kubu Raya CSIRT](https://csirt.kuburayakab.go.id)

---

## Roadmap

| Phase | Status | Fitur |
|-------|--------|-------|
| **Phase 1** | ✅ Selesai | CLI, SSH, firewall, scanner, cron, AIDE, kernel, accounts, services, audit, dashboard, alerts |
| **Phase 2** | 🚧 Proses | WAF (ModSecurity), IDS real-time, eBPF, Docker protection, auto security updates |
| **Phase 3** | 🔮 Rencana | AI anomaly detection, multi-server, compliance (CIS/NIST), K8s hardening, zero trust |

---

## License

MIT License — see [LICENSE](LICENSE)

---

<div align="center">
  <br/>
  <pre style="font-family: monospace; font-size: 1.2em; line-height: 1.4; color: #00ff88;">
╔═══════════════════════════════════════════════════════════════╗
║              One Command to Harden Them All                   ║
║                 by KUBU RAYA CSIRT - 2026                     ║
╚═══════════════════════════════════════════════════════════════╝
  </pre>
  <br/>
  <p><sub>Built for the security community</sub></p>
  <br/>
</div>
