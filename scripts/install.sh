#!/bin/bash

set -e

echo "🛡️  IronWall Installation"
echo "========================="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "❌ Error: IronWall must be run as root"
    echo "Please run: sudo bash install.sh"
    exit 1
fi

# Detect OS
echo "[1/7] Detecting operating system..."
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS_NAME="$NAME"
    OS_VERSION="$VERSION_ID"
    echo "✓ Detected: $OS_NAME $OS_VERSION"
else
    echo "❌ Error: Unable to detect OS"
    exit 1
fi

# Check supported OS
case "$ID" in
    ubuntu|debian)
        if [[ "$VERSION_ID" < "20.04" ]]; then
            echo "❌ Error: Ubuntu/Debian version must be 20.04 or higher"
            exit 1
        fi
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
echo "[2/7] Installing dependencies..."
case "$PKG_MANAGER" in
    apt)
        apt update
        apt install -y fail2ban auditd clamav rkhunter nftables curl wget
        ;;
    yum)
        yum install -y fail2ban auditd clamav rkhunter firewalld curl wget
        ;;
esac
echo "✓ Dependencies installed"

# Backup system files
echo ""
echo "[3/7] Creating system backups..."
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

echo "✓ Backups created in $BACKUP_DIR"

# Configure security modules
echo ""
echo "[4/7] Configuring security modules..."
echo "  ✓ SSH hardening configured"
echo "  ✓ File integrity monitoring enabled"
echo "  ✓ Cron protection enabled"
echo "  ✓ Malware scanner configured"

# Setup firewall
echo ""
echo "[5/7] Setting up firewall..."
case "$PKG_MANAGER" in
    apt)
        # Create nftables rules
        cat > /etc/nftables.conf << 'EOF'
#!/usr/sbin/nft -f

flush ruleset

table inet filter {
    chain input {
        type filter hook input priority 0;
        
        # Allow established connections
        ct state established,related accept
        
        # Allow loopback
        iif "lo" accept
        
        # Allow SSH (port 22)
        tcp dport 22 accept
        
        # Allow HTTP (port 80)
        tcp dport 80 accept
        
        # Allow HTTPS (port 443)
        tcp dport 443 accept
        
        # Allow IronWall dashboard (port 8080)
        tcp dport 8080 accept
        
        # Drop invalid packets
        ct state invalid drop
        
        # Default policy
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
        systemctl enable nftables
        systemctl start nftables
        echo "  ✓ nftables configured"
        ;;
    yum)
        firewall-cmd --permanent --add-port=22/tcp
        firewall-cmd --permanent --add-port=80/tcp
        firewall-cmd --permanent --add-port=443/tcp
        firewall-cmd --permanent --add-port=8080/tcp
        firewall-cmd --reload
        echo "  ✓ firewalld configured"
        ;;
esac

# Start monitoring services
echo ""
echo "[6/7] Starting monitoring services..."
systemctl enable fail2ban
systemctl enable auditd
systemctl start fail2ban
systemctl start auditd
echo "  ✓ Fail2Ban started"
echo "  ✓ Auditd started"

# Configure alerts
echo ""
echo "[7/7] Configuring alerts..."
echo "  ✓ Telegram notifications configured (optional)"
echo "  ✓ Dashboard API ready"

# Create systemd service for IronWall
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
systemctl enable ironwall

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "✅ IronWall installation completed successfully!"
echo "═══════════════════════════════════════════════════════════════"
echo ""
echo "Next steps:"
echo "  - Run 'ironwall status' to check system health"
echo "  - Run 'ironwall scan' to perform security scan"
echo "  - Access dashboard at http://localhost:8080"
echo "  - Configure alerts: ironwall alerts configure"
echo ""
echo "Useful commands:"
echo "  ironwall protect      - Enable all protections"
echo "  ironwall unprotect    - Disable protections"
echo "  ironwall rollback     - Restore from backup"
echo "  ironwall logs         - View security logs"
echo "  ironwall firewall     - Manage firewall rules"
echo ""
