# IronWall - Automated Linux System Hardening & Threat Prevention Toolkit

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/)
[![Python](https://img.shields.io/badge/python-3.11+-blue.svg)](https://www.python.org/)
[![Next.js](https://img.shields.io/badge/next.js-15+-black.svg)](https://nextjs.org/)

IronWall is a comprehensive security solution that automates server hardening, malware detection, intrusion prevention, and real-time monitoring for Linux servers.

## Features

- 🔐 **SSH Hardening** - Disable root login, enforce key-based authentication
- 🛡️ **Firewall Automation** - Auto-configure nftables/iptables with sensible defaults
- 🚫 **Fail2Ban Integration** - Automatic brute force protection
- 🦠 **Malware Scanner** - Detect and quarantine malware, cryptominers, webshells
- ⏰ **Cron Protection** - Monitor and alert on suspicious cron jobs
- 📁 **File Integrity** - Protect critical system files with immutable flags
- 🔔 **Real-time Alerts** - Telegram, Discord, Slack, Email notifications
- 📊 **Web Dashboard** - Monitor security status and manage modules

## Quick Start

### One-Line Installation

```bash
curl -sSL https://ironwall.sh/install.sh | sudo bash
```

Or:

```bash
wget -O - https://ironwall.sh/install.sh | sudo bash
```

## Usage

### CLI Commands

```bash
# Check system status
ironwall status

# Run security scan
ironwall scan

# Enable all protections
ironwall protect

# Disable protections
ironwall unprotect --force

# Rollback to previous configuration
ironwall rollback --force

# View logs
ironwall logs
ironwall logs --follow

# Manage firewall
ironwall firewall status
ironwall firewall block 192.168.1.100
ironwall firewall unblock 192.168.1.100
ironwall firewall list
```

### Web Dashboard

Access the dashboard at `http://localhost:8080` or `http://your-server-ip:8080`

## Supported Operating Systems

### Phase 1 (Current)
- Ubuntu 20.04+
- Ubuntu 22.04+
- Debian 11+
- Debian 12+

### Phase 2 (Planned)
- CentOS
- AlmaLinux
- Rocky Linux

## Architecture

```
ironwall/
├── cmd/ironwall/          # Go CLI binary
├── internal/              # Go internal packages
│   ├── hardening/         # SSH hardening
│   ├── firewall/          # Firewall management
│   ├── scanner/           # Malware scanner
│   ├── monitor/           # Cron monitoring
│   └── backup/            # Backup/restore
├── backend/               # FastAPI + SQLite
├── frontend/              # Next.js dashboard
├── configs/               # Default configs
└── scripts/               # Installer scripts
```

## Development

### Backend

```bash
cd backend
pip install -r requirements.txt
uvicorn app.main:app --reload --port 8001
```

### Frontend

```bash
cd frontend
npm install
npm run dev
```

### CLI

```bash
go build -o ironwall cmd/ironwall/main.go
./ironwall status
```

## API Documentation

Once the backend is running, visit:
- Swagger UI: `http://localhost:8001/docs`
- ReDoc: `http://localhost:8001/redoc`

## Configuration

### Backend (.env)
```env
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id
DASHBOARD_PORT=8080
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8001
```

## Security Modules

### SSH Hardening
- Disable root login
- Disable password authentication
- Enforce public key authentication
- Change default SSH port (optional)
- Rate limiting

### Firewall
- Auto-configure nftables/iptables
- Default deny policy
- Allow only required ports (22, 80, 443, 8080)
- Rate limiting
- Anti-port scanning

### Malware Scanner
- ClamAV integration
- Process monitoring for cryptominers (xmrig, kinsing)
- Webshell detection
- Reverse shell detection
- Auto-kill and quarantine

### Cron Protection
- Monitor /etc/crontab and user crontabs
- Detect suspicious patterns (curl|bash, wget|sh)
- Real-time alerts
- Whitelist mechanism

### File Integrity
- Protect critical files with chattr +i
- Monitor /etc/passwd, /etc/shadow, /etc/ssh/
- Baseline snapshots
- Real-time change detection

## Performance

- CPU overhead: <5%
- RAM usage: <300MB (including dashboard)
- Alert latency: <3 seconds
- Boot impact: minimal

## Contributing

Contributions are welcome! Please read our contributing guidelines.

## License

MIT License - see LICENSE file for details

## Support

- Documentation: https://ironwall.sh/docs
- Issues: https://github.com/yourusername/ironwall/issues
- Discord: https://discord.gg/ironwall

## Roadmap

### Phase 1 (MVP) ✅
- CLI tool
- SSH hardening
- Firewall automation
- Malware scanner
- Cron protection
- File integrity
- Web dashboard

### Phase 2 (Planned)
- WAF (ModSecurity)
- Real-time intrusion detection
- eBPF monitoring
- Docker protection

### Phase 3 (Future)
- AI anomaly detection
- Cluster protection
- Kubernetes hardening
- Zero trust mode

---

**⚠️ Warning**: IronWall makes significant changes to your system configuration. Always test in a non-production environment first and ensure you have backups.
