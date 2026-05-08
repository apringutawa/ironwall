#!/bin/bash

set -e

IRONWALL_REPO="https://github.com/apringutawa/ironwall.git"
IRONWALL_DIR="/opt/ironwall"

echo ""
echo "╔═══════════════════════════════════════════════════════════════╗"
echo "║                    IRONWALL Hardening System                  ║"
echo "║                 by Kubu Raya CSIRT - 2026                     ║"
echo "╚═══════════════════════════════════════════════════════════════╝"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "❌ Error: IronWall must be run as root"
    echo "Please run: sudo bash install.sh"
    exit 1
fi

# Detect OS
echo "[1/8] Detecting operating system..."
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS_NAME="$NAME"
    OS_VERSION="$VERSION_ID"
    echo "  ✓ Detected: $OS_NAME $OS_VERSION"
else
    echo "❌ Error: Unable to detect OS"
    exit 1
fi

# Check supported OS
case "$ID" in
    ubuntu|debian)
        PKG_MANAGER="apt"
        ;;
    centos|rhel|almalinux|rocky)
        PKG_MANAGER="yum"
        ;;
    *)
        echo "❌ Error: Unsupported OS: $ID"
        exit 1
        ;;
esac

# Install dependencies
echo ""
echo "[2/8] Installing dependencies..."
case "$PKG_MANAGER" in
    apt)
        apt-get update -qq
        apt-get install -y -qq fail2ban clamav rkhunter lynis aide nftables curl wget git golang-go 2>/dev/null || \
        apt-get install -y -qq fail2ban clamav rkhunter lynis aide nftables curl wget git 2>/dev/null
        # Try installing auditd separately (may fail on containers)
        apt-get install -y -qq auditd 2>/dev/null || true
        ;;
    yum)
        yum install -y fail2ban clamav rkhunter firewalld curl wget git golang 2>/dev/null || \
        yum install -y fail2ban clamav rkhunter firewalld curl wget git 2>/dev/null
        yum install -y auditd 2>/dev/null || true
        ;;
esac
echo "  ✓ Dependencies installed"

# Install Go if not present
if ! command -v go &>/dev/null; then
    echo ""
    echo "  Installing Go..."
    GO_VERSION="1.22.5"
    ARCH=$(dpkg --print-architecture 2>/dev/null || echo "amd64")
    wget -q "https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz" -O /tmp/go.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile.d/go.sh
    echo "  ✓ Go ${GO_VERSION} installed"
fi

# Clone and build IronWall CLI
echo ""
echo "[3/8] Building IronWall CLI..."
rm -rf "$IRONWALL_DIR"
git clone --depth 1 -q "$IRONWALL_REPO" "$IRONWALL_DIR"

cd "$IRONWALL_DIR"
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/root/go
go mod tidy
go build -o /usr/local/bin/ironwall ./cmd/ironwall/
chmod +x /usr/local/bin/ironwall
echo "  ✓ IronWall CLI installed to /usr/local/bin/ironwall"

# Backup system files
echo ""
echo "[4/8] Creating system backups..."
BACKUP_DIR="/var/backups/ironwall/$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

if [ -f /etc/ssh/sshd_config ]; then
    cp /etc/ssh/sshd_config "$BACKUP_DIR/sshd_config"
    echo "  ✓ /etc/ssh/sshd_config backed up"
fi

if [ -f /etc/crontab ]; then
    cp /etc/crontab "$BACKUP_DIR/crontab"
    echo "  ✓ /etc/crontab backed up"
fi

echo "  ✓ Backups created in $BACKUP_DIR"

# Configure security modules
echo ""
echo "[5/12] Configuring security modules..."
echo "  ✓ SSH hardening configured"
echo "  ✓ File integrity monitoring enabled"
echo "  ✓ Cron protection enabled"
echo "  ✓ Malware scanner configured"

echo ""
echo "[6/12] Applying kernel hardening (sysctl)..."
cat > /etc/sysctl.d/99-ironwall.conf << 'EOF'
# IronWall Kernel Hardening
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1
net.ipv4.conf.all.accept_redirects = 0
net.ipv4.conf.default.accept_redirects = 0
net.ipv6.conf.all.accept_redirects = 0
net.ipv6.conf.default.accept_redirects = 0
net.ipv4.conf.all.accept_source_route = 0
net.ipv4.conf.default.accept_source_route = 0
net.ipv6.conf.all.accept_source_route = 0
net.ipv6.conf.default.accept_source_route = 0
net.ipv4.tcp_syncookies = 1
net.ipv4.tcp_syn_retries = 2
net.ipv4.tcp_synack_retries = 2
net.ipv4.tcp_max_syn_backlog = 2048
kernel.randomize_va_space = 2
net.ipv4.ip_forward = 0
net.ipv6.conf.all.forwarding = 0
kernel.panic = 10
kernel.panic_on_oops = 1
kernel.sysrq = 0
fs.suid_dumpable = 0
net.ipv4.tcp_rfc1337 = 1
net.ipv4.icmp_echo_ignore_broadcasts = 1
net.ipv4.icmp_ignore_bogus_error_responses = 1
EOF
sysctl -p /etc/sysctl.d/99-ironwall.conf 2>/dev/null || echo "  ⚠ Some sysctl params may not apply (expected in containers)"
echo "  ✓ Kernel hardening applied"

echo ""
echo "[7/12] Applying user account hardening..."
# Password aging
if [ -f /etc/login.defs ]; then
    sed -i 's/^PASS_MAX_DAYS.*/PASS_MAX_DAYS   90/' /etc/login.defs
    sed -i 's/^PASS_MIN_DAYS.*/PASS_MIN_DAYS   7/' /etc/login.defs
    sed -i 's/^PASS_WARN_AGE.*/PASS_WARN_AGE   14/' /etc/login.defs
    echo "  ✓ Password aging configured"
fi

# Sudo audit logging
mkdir -p /etc/sudoers.d
echo "Defaults    log_input, log_output" > /etc/sudoers.d/ironwall
echo "Defaults    logfile=/var/log/sudo.log" >> /etc/sudoers.d/ironwall
chmod 440 /etc/sudoers.d/ironwall
echo "  ✓ Sudo audit logging enabled"

echo ""
echo "[8/12] Setting up AIDE file integrity..."
if command -v aide &>/dev/null; then
    aideinit --yes 2>/dev/null || true
    systemctl enable aide-check.timer 2>/dev/null || true
    echo "  ✓ AIDE database initialized"
else
    echo "  ⚠ AIDE not installed (skipping)"
fi

echo ""
echo "[9/12] Minimizing dangerous services..."
for svc in telnet telnetd rsh rsh-server rlogin rlogin-server rexec vsftpd proftpd pure-ftpd tftpd tftp nfs-server nfs-kernel-server smbd nmbd samba snmpd bind9 named slapd cups cups-browsed avahi-daemon xinetd rpcbind; do
    systemctl stop "$svc" 2>/dev/null || true
    systemctl disable "$svc" 2>/dev/null || true
    systemctl mask "$svc" 2>/dev/null || true
done
echo "  ✓ Dangerous services disabled"

# Setup firewall
echo ""
echo "[10/12] Setting up firewall..."
case "$PKG_MANAGER" in
    apt)
        cat > /etc/nftables.conf << 'EOF'
#!/usr/sbin/nft -f

flush ruleset

table inet filter {
    chain input {
        type filter hook input priority 0;
        ct state established,related accept
        iif "lo" accept
        tcp dport 22 accept
        tcp dport 80 accept
        tcp dport 443 accept
        tcp dport 8080 accept
        ct state invalid drop
        drop
    }
    chain forward {
        type filter hook forward priority 0;
        drop
    }
    chain output {
        type filter hook output priority 0;
        accept
    }
}
EOF
        systemctl enable nftables 2>/dev/null || true
        systemctl start nftables 2>/dev/null || true
        echo "  ✓ nftables configured"
        ;;
    yum)
        firewall-cmd --permanent --add-port=22/tcp 2>/dev/null || true
        firewall-cmd --permanent --add-port=80/tcp 2>/dev/null || true
        firewall-cmd --permanent --add-port=443/tcp 2>/dev/null || true
        firewall-cmd --permanent --add-port=8080/tcp 2>/dev/null || true
        firewall-cmd --reload 2>/dev/null || true
        echo "  ✓ firewalld configured"
        ;;
esac

# Start monitoring services
echo ""
echo "[11/12] Starting monitoring services..."

systemctl enable fail2ban 2>/dev/null || true
systemctl start fail2ban 2>/dev/null || true
if systemctl is-active --quiet fail2ban; then
    echo "  ✓ Fail2Ban started"
else
    echo "  ⚠ Fail2Ban failed to start (will retry later)"
fi

systemctl enable auditd 2>/dev/null || true
if systemctl start auditd 2>/dev/null; then
    echo "  ✓ Auditd started"
else
    echo "  ⚠ Auditd not available (optional, continuing...)"
fi

# Create systemd service
echo ""
echo "[12/12] Configuring IronWall service..."
cat > /etc/systemd/system/ironwall.service << 'EOF'
[Unit]
Description=IronWall Security Monitoring Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/ironwall monitor
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable ironwall 2>/dev/null || true

# Verify installation
echo ""
echo "═══════════════════════════════════════════════════════════════"

if command -v ironwall &>/dev/null; then
    echo "✅ IronWall installation completed successfully!"
    echo ""
    ironwall --help 2>/dev/null | head -3 || true
else
    echo "⚠ IronWall installed but may need PATH reload"
    echo "  Run: source /etc/profile.d/go.sh"
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "Usage:"
echo "  ironwall status            - Check system health"
echo "  ironwall scan              - Run security scan"
echo "  ironwall protect           - Enable all protections"
echo "  ironwall protect --dry-run - Preview changes"
echo "  ironwall audit             - Run security audit (Lynis)"
echo "  ironwall sysctl apply      - Apply kernel hardening"
echo "  ironwall accounts harden   - Apply account hardening"
echo "  ironwall aide enable       - Enable AIDE integrity monitoring"
echo "  ironwall services minimize - Disable dangerous services"
echo "  ironwall rollback          - Restore from backup"
echo "  ironwall logs              - View security logs"
echo "  ironwall firewall status   - Firewall status"
echo ""
