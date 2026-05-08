# AGENTS.md - IronWall

## Project Overview

IronWall is an automated Linux system hardening and threat prevention toolkit. It provides comprehensive security including SSH hardening, firewall automation, malware detection, intrusion prevention, and real-time monitoring.

## Architecture

### Components

1. **CLI Tool (Go)** - Command-line interface for system administration
2. **Backend API (FastAPI)** - REST API for dashboard and monitoring
3. **Frontend Dashboard (Next.js)** - Web-based management interface
4. **Security Modules (Go)** - Core security functionality

### Tech Stack

- **CLI**: Go 1.21+ with Cobra
- **Backend**: FastAPI + SQLAlchemy + SQLite
- **Frontend**: Next.js 15 + React 19 + TypeScript + Tailwind CSS 4
- **Security Tools**: ClamAV, Fail2Ban, nftables/iptables, auditd, Lynis, AIDE

## Project Structure

```
ironwall/
├── cmd/ironwall/          # Go CLI entry point
│   ├── main.go
│   └── cmd/               # Cobra commands
├── internal/              # Go internal packages
│   ├── audit/             # Lynis security audit
│   ├── hardening/         # SSH, sysctl, accounts, AIDE, services
│   ├── firewall/          # Firewall management
│   ├── scanner/           # Malware scanner
│   ├── monitor/           # Cron monitoring
│   └── backup/            # Backup/restore
├── backend/               # FastAPI backend
│   ├── app/
│   │   ├── api/           # API routes
│   │   ├── models/        # Database models & schemas
│   │   └── core/          # Config, database, alerts
│   └── requirements.txt
├── frontend/              # Next.js frontend
│   ├── src/
│   │   ├── app/           # Next.js pages (App Router)
│   │   ├── components/    # React components
│   │   └── lib/           # API client
│   └── package.json
├── configs/               # Default security configs
├── scripts/               # Installation scripts
└── docs/                  # Documentation
```

## Development Workflow

### Backend Development

```bash
cd backend
pip install -r requirements.txt
cp .env.example .env
uvicorn app.main:app --reload --port 8001
```

API docs: http://localhost:8001/docs

### Frontend Development

```bash
cd frontend
npm install
cp .env.example .env.local
npm run dev
```

Dashboard: http://localhost:3000

### CLI Development

```bash
# Build CLI
go build -o ironwall cmd/ironwall/main.go

# Test commands
./ironwall status
./ironwall scan
./ironwall protect --dry-run
```

## Security Modules

### 1. SSH Hardening (`internal/hardening/ssh.go`)
- Disables root login
- Disables password authentication
- Enforces public key authentication
- Optional port change
- Automatic backup before changes

### 2. Firewall (`internal/firewall/firewall.go`)
- Auto-detects nftables/iptables/firewalld
- Configures default deny policy
- Allows ports: 22, 80, 443, 8080
- IP blocking/unblocking
- Port management

### 3. Malware Scanner (`internal/scanner/scanner.go`)
- ClamAV integration
- Process scanning for cryptominers
- Webshell detection
- Quarantine functionality

### 4. Cron Monitor (`internal/monitor/cron.go`)
- Monitors /etc/crontab and user crontabs
- Detects suspicious patterns
- Real-time monitoring
- Alert generation

### 5. Backup Manager (`internal/backup/backup.go`)
- Creates backups before changes
- Restores from backup
- File locking/unlocking (chattr)

### 6. Lynis Audit (`internal/audit/lynis.go`)
- Runs comprehensive security audit
- Parses hardening score and recommendations
- Auto-installs Lynis if missing
- Tests performed tracking

### 7. Kernel Sysctl Hardening (`internal/hardening/sysctl.go`)
- 27 kernel parameters for network hardening
- ASLR, anti-spoofing, SYN flood protection
- ICMP redirect/ignore configuration
- TCP stack hardening

### 8. User Account Hardening (`internal/hardening/accounts.go`)
- PAM password policy (min length 12, complexity)
- Faillock lockout (5 attempts, 15min)
- Sudo I/O audit logging
- Password aging (90 days max)
- UID 0 check, empty password check

### 9. AIDE Integrity (`internal/hardening/aide.go`)
- File integrity database initialization
- Daily automated integrity checks
- Fallback when auditd unavailable
- Systemd timer integration

### 10. Service Minimization (`internal/hardening/services.go`)
- Detects dangerous services (telnet, ftp, nfs, etc.)
- Systemd/SysVinit auto-detection
- Disable + mask dangerous services
- Dry-run support

## API Endpoints

### Status & Modules
- `GET /api/v1/status` - System health
- `GET /api/v1/modules` - List modules
- `POST /api/v1/modules/{name}/enable` - Enable module
- `POST /api/v1/modules/{name}/disable` - Disable module

### Events
- `GET /api/v1/events` - Get security events
- `POST /api/v1/events` - Create event
- `GET /api/v1/events/stats` - Event statistics

### Firewall
- `GET /api/v1/firewall/status` - Firewall status
- `POST /api/v1/firewall/block` - Block IP
- `POST /api/v1/firewall/unblock` - Unblock IP
- `GET /api/v1/firewall/blocked` - List blocked IPs

### Alerts
- `POST /api/v1/alerts/configure` - Configure alerts
- `GET /api/v1/alerts/configure` - Get alert config
- `POST /api/v1/alerts/test` - Test alert

### Hardening
- `GET /api/v1/hardening/status` - Hardening module status
- `POST /api/v1/hardening/kernel/apply` - Apply kernel sysctl params
- `POST /api/v1/hardening/accounts/apply` - Apply account hardening
- `POST /api/v1/hardening/aide/init` - Init AIDE database
- `POST /api/v1/hardening/aide/check` - Run AIDE check
- `POST /api/v1/hardening/services/minimize` - Disable dangerous services
- `GET /api/v1/hardening/services/list` - List services by class
- `GET /api/v1/hardening/audit/status` - Lynis audit status
- `POST /api/v1/hardening/audit/run` - Run Lynis audit
- `POST /api/v1/hardening/audit/install` - Install Lynis

## CLI Commands

```bash
ironwall install              # Install and configure IronWall
ironwall status               # Show system status
ironwall scan                 # Run security scan
ironwall protect              # Enable all protections
ironwall unprotect            # Disable protections
ironwall rollback             # Restore from backup
ironwall logs                 # View logs
ironwall firewall status      # Firewall status
ironwall firewall block       # Block IP
ironwall firewall unblock     # Unblock IP
ironwall audit                # Run security audit (Lynis)
ironwall audit install        # Install Lynis
ironwall audit score          # Show hardening score
ironwall sysctl apply         # Apply kernel hardening
ironwall sysctl verify        # Verify kernel params
ironwall sysctl restore       # Restore sysctl backup
ironwall accounts harden      # Apply account hardening
ironwall accounts list        # List system users
ironwall accounts check       # Check account security
ironwall accounts restore     # Restore account config
ironwall aide init            # Init AIDE database
ironwall aide check           # Run integrity check
ironwall aide update          # Update AIDE database
ironwall aide enable          # Enable daily monitoring
ironwall services list        # List services by class
ironwall services minimize    # Disable dangerous services
ironwall services scan        # Scan insecure services
```

## Configuration

### Backend (.env)
```env
TELEGRAM_BOT_TOKEN=your_token
TELEGRAM_CHAT_ID=your_chat_id
DASHBOARD_PORT=8080
```

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8001
```

## Installation

### Development
```bash
# Clone and setup
git clone <repo>
cd ironwall

# Backend
cd backend && pip install -r requirements.txt

# Frontend
cd frontend && npm install

# CLI
go build -o ironwall cmd/ironwall/main.go
```

### Production
```bash
curl -sSL https://ironwall.sh/install.sh | sudo bash
```

## Testing

### Backend Tests
```bash
cd backend
pytest
```

### Frontend Tests
```bash
cd frontend
npm test
```

### CLI Tests
```bash
go test ./...
```

## Deployment

### Docker Compose
```bash
docker-compose up --build
```

### Manual
```bash
# Build CLI
go build -o ironwall cmd/ironwall/main.go
sudo mv ironwall /usr/local/bin/

# Run installer
sudo ironwall install
```

## Important Notes

### Security Considerations
- Always backup before applying changes
- Test in non-production environment first
- Use `--dry-run` flag to preview changes
- Keep recovery access method available

### File Locations
- Backups: `/var/backups/ironwall/`
- Database: `backend/ironwall.db`
- Configs: `/etc/nftables.conf`, `/etc/ssh/sshd_config`
- Logs: `/var/log/ironwall/`

### Supported OS
- Ubuntu 20.04+, 22.04+
- Debian 11+, 12+
- CentOS/AlmaLinux/Rocky (planned)

## Troubleshooting

### Backend won't start
- Check Python version (3.11+)
- Verify dependencies: `pip install -r requirements.txt`
- Check port 8001 availability

### Frontend won't start
- Check Node version (20+)
- Run `npm install` again
- Verify API URL in `.env.local`

### CLI errors
- Ensure running as root: `sudo ironwall`
- Check Go version (1.21+)
- Verify system compatibility

## Contributing

When adding new features:
1. Follow existing code structure
2. Add tests for new functionality
3. Update documentation
4. Test on supported OS versions
5. Create PR with clear description

## Roadmap

### Phase 1 (MVP) ✅
- CLI tool with core commands
- SSH hardening
- Firewall automation
- Malware scanner
- Cron protection
- File integrity
- Web dashboard
- Alert system
- Lynis security audit ✅
- Kernel sysctl hardening ✅
- User account hardening ✅
- AIDE file integrity ✅
- Service minimization ✅

### Phase 2 (Planned)
- WAF integration (ModSecurity)
- Real-time intrusion detection
- eBPF monitoring
- Docker container protection
- Kubernetes support

### Phase 3 (Future)
- AI anomaly detection
- Multi-server management
- Compliance reporting
- Zero trust mode
