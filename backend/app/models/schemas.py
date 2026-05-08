from pydantic import BaseModel
from datetime import datetime
from typing import Optional, List, Dict, Any

class ServerCreate(BaseModel):
    hostname: str
    ip_address: str
    os_type: str
    os_version: str

class ServerResponse(BaseModel):
    id: int
    hostname: str
    ip_address: str
    os_type: str
    os_version: str
    installed_at: datetime
    last_scan: Optional[datetime] = None
    
    class Config:
        from_attributes = True

class EventCreate(BaseModel):
    server_id: int
    event_type: str
    severity: str
    description: str
    details: Optional[Dict[str, Any]] = None

class EventResponse(BaseModel):
    id: int
    server_id: int
    event_type: str
    severity: str
    description: str
    details: Optional[Dict[str, Any]] = None
    created_at: datetime
    
    class Config:
        from_attributes = True

class ModuleUpdate(BaseModel):
    enabled: Optional[bool] = None
    status: Optional[str] = None

class ModuleResponse(BaseModel):
    id: int
    name: str
    enabled: bool
    status: str
    last_run: Optional[datetime] = None
    
    class Config:
        from_attributes = True

class ConfigUpdate(BaseModel):
    key: str
    value: str

class AlertConfig(BaseModel):
    telegram_token: Optional[str] = None
    telegram_chat_id: Optional[str] = None
    discord_webhook: Optional[str] = None
    slack_webhook: Optional[str] = None
    email: Optional[str] = None

class FirewallBlockRequest(BaseModel):
    ip: str
    reason: Optional[str] = None

class FirewallAllowRequest(BaseModel):
    port: int
    protocol: str = "tcp"

class SystemStatus(BaseModel):
    status: str
    uptime: str
    cpu_usage: float
    memory_usage: float
    modules: List[ModuleResponse]
    blocked_ips_count: int
    events_24h: int

class ScanRequest(BaseModel):
    paths: Optional[List[str]] = None
    quick: bool = False
    full: bool = False

class ScanResult(BaseModel):
    files_scanned: int
    malware_found: int
    rootkits_found: int
    webshells_found: int
    suspicious_processes: int
    timestamp: datetime

class HardeningStatus(BaseModel):
    kernel_hardening: dict
    account_hardening: dict
    file_integrity: dict
    service_minimization: dict
    security_audit: dict

class AuditResult(BaseModel):
    hardening_score: float
    warnings: int
    suggestions: int
    tests_performed: int
    status: str
    last_run: Optional[datetime] = None

class ServiceInfo(BaseModel):
    name: str
    status: str
    classification: str

class KernelParam(BaseModel):
    key: str
    value: str
    current: Optional[str] = None
    description: str
    applied: bool
