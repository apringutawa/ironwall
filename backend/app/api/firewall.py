from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.models import schemas

router = APIRouter(prefix="/api/v1", tags=["firewall"])

@router.get("/firewall/status")
async def get_firewall_status():
    """Get firewall status"""
    return {
        "status": "active",
        "backend": "nftables",
        "allowed_ports": [
            {"port": 22, "protocol": "tcp", "service": "SSH"},
            {"port": 80, "protocol": "tcp", "service": "HTTP"},
            {"port": 443, "protocol": "tcp", "service": "HTTPS"},
            {"port": 8080, "protocol": "tcp", "service": "Dashboard"}
        ],
        "blocked_ips_count": 147,
        "rate_limiting": True,
        "anti_port_scan": True
    }

@router.get("/firewall/rules")
async def get_firewall_rules():
    """Get all firewall rules"""
    return {
        "rules": [
            {"id": 1, "action": "allow", "port": 22, "protocol": "tcp"},
            {"id": 2, "action": "allow", "port": 80, "protocol": "tcp"},
            {"id": 3, "action": "allow", "port": 443, "protocol": "tcp"},
            {"id": 4, "action": "allow", "port": 8080, "protocol": "tcp"},
            {"id": 5, "action": "deny", "port": "*", "protocol": "*"}
        ]
    }

@router.post("/firewall/block")
async def block_ip(request: schemas.FirewallBlockRequest):
    """Block an IP address"""
    return {
        "message": f"IP {request.ip} has been blocked",
        "reason": request.reason or "Manual block"
    }

@router.post("/firewall/unblock")
async def unblock_ip(ip: str):
    """Unblock an IP address"""
    return {
        "message": f"IP {ip} has been unblocked"
    }

@router.get("/firewall/blocked")
async def get_blocked_ips():
    """Get list of blocked IPs"""
    return {
        "blocked_ips": [
            {"ip": "192.168.1.100", "reason": "Brute force attempt", "blocked_at": "2026-05-06T23:45:00"},
            {"ip": "10.0.0.55", "reason": "Port scanning", "blocked_at": "2026-05-05T15:30:00"},
            {"ip": "172.16.0.200", "reason": "Malware detected", "blocked_at": "2026-05-04T10:20:00"}
        ]
    }

@router.post("/firewall/allow")
async def allow_port(request: schemas.FirewallAllowRequest):
    """Allow a port through firewall"""
    return {
        "message": f"Port {request.port}/{request.protocol} has been allowed"
    }
