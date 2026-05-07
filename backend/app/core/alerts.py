import httpx
from typing import Optional
from app.core.config import settings

class AlertService:
    def __init__(self):
        self.telegram_token = settings.TELEGRAM_BOT_TOKEN
        self.telegram_chat_id = settings.TELEGRAM_CHAT_ID

    async def send_telegram(self, message: str) -> bool:
        if not self.telegram_token or not self.telegram_chat_id:
            return False
        
        url = f"https://api.telegram.org/bot{self.telegram_token}/sendMessage"
        
        async with httpx.AsyncClient() as client:
            response = await client.post(url, json={
                "chat_id": self.telegram_chat_id,
                "text": message,
                "parse_mode": "HTML"
            })
        
        return response.status_code == 200

    async def send_discord(self, webhook_url: str, message: str) -> bool:
        if not webhook_url:
            return False
        
        async with httpx.AsyncClient() as client:
            response = await client.post(webhook_url, json={
                "content": message
            })
        
        return response.status_code == 204

    async def send_slack(self, webhook_url: str, message: str) -> bool:
        if not webhook_url:
            return False
        
        async with httpx.AsyncClient() as client:
            response = await client.post(webhook_url, json={
                "text": message
            })
        
        return response.status_code == 200

    async def send_alert(
        self,
        title: str,
        severity: str,
        description: str,
        server: str = "Unknown",
        action: str = "None",
        telegram: bool = True,
        discord_webhook: Optional[str] = None,
        slack_webhook: Optional[str] = None
    ) -> dict:
        severity_emoji = {
            "critical": "🔴",
            "warning": "⚠️",
            "info": "ℹ️"
        }
        
        emoji = severity_emoji.get(severity, "⚠️")
        
        message = f"""{emoji} <b>{title}</b>

<b>Server:</b> {server}

<b>Severity:</b> {severity.upper()}

<b>Details:</b>
{description}

<b>Action:</b> {action}
"""
        
        results = {
            "telegram": False,
            "discord": False,
            "slack": False
        }
        
        if telegram:
            results["telegram"] = await self.send_telegram(message)
        
        if discord_webhook:
            results["discord"] = await self.send_discord(discord_webhook, message)
        
        if slack_webhook:
            results["slack"] = await self.send_slack(slack_webhook, message)
        
        return results

    async def send_intrusion_alert(
        self,
        threat_type: str,
        description: str,
        server: str,
        action: str = "Process terminated"
    ):
        return await self.send_alert(
            title="Intrusion Detected",
            severity="critical",
            description=f"Threat: {threat_type}\n{description}",
            server=server,
            action=action
        )

    async def send_malware_alert(
        self,
        malware_type: str,
        file_path: str,
        server: str,
        action: str = "Quarantined"
    ):
        return await self.send_alert(
            title="Malware Detected",
            severity="critical",
            description=f"Type: {malware_type}\nLocation: {file_path}",
            server=server,
            action=action
        )

    async def send_brute_force_alert(
        self,
        ip: str,
        attempts: int,
        server: str,
        action: str = "IP blocked"
    ):
        return await self.send_alert(
            title="Brute Force Attempt",
            severity="warning",
            description=f"IP: {ip}\nAttempts: {attempts}",
            server=server,
            action=action
        )

    async def send_webshell_alert(
        self,
        file_path: str,
        server: str,
        action: str = "File quarantined"
    ):
        return await self.send_alert(
            title="Webshell Detected",
            severity="critical",
            description=f"Location: {file_path}",
            server=server,
            action=action
        )

    async def send_cron_alert(
        self,
        cron_path: str,
        suspicious_pattern: str,
        server: str,
        action: str = "Alert only"
    ):
        return await self.send_alert(
            title="Suspicious Cron Job",
            severity="warning",
            description=f"File: {cron_path}\nPattern: {suspicious_pattern}",
            server=server,
            action=action
        )

alert_service = AlertService()
