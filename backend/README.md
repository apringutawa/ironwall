# IronWall Backend

FastAPI backend for IronWall security monitoring system.

## Setup

```bash
cd backend
pip install -r requirements.txt
cp .env.example .env
uvicorn app.main:app --reload --port 8001
```

## API Endpoints

### Status
- `GET /api/v1/status` - System health and status
- `GET /api/v1/modules` - List all security modules
- `POST /api/v1/modules/{name}/enable` - Enable a module
- `POST /api/v1/modules/{name}/disable` - Disable a module

### Events
- `GET /api/v1/events` - Get security events
- `POST /api/v1/events` - Create a new event
- `GET /api/v1/events/stats` - Event statistics

### Scan
- `POST /api/v1/scan` - Run security scan
- `POST /api/v1/protect` - Enable all protections
- `POST /api/v1/unprotect` - Disable protections
- `POST /api/v1/rollback` - Rollback to previous state

### Firewall
- `GET /api/v1/firewall/status` - Firewall status
- `GET /api/v1/firewall/rules` - List firewall rules
- `POST /api/v1/firewall/block` - Block an IP
- `POST /api/v1/firewall/unblock` - Unblock an IP
- `GET /api/v1/firewall/blocked` - List blocked IPs
- `POST /api/v1/firewall/allow` - Allow a port

### Alerts
- `POST /api/v1/alerts/configure` - Configure alert settings
- `GET /api/v1/alerts/configure` - Get alert configuration
- `POST /api/v1/alerts/test` - Send test alert

## Database

Uses SQLite by default. The database file is created at `backend/ironwall.db`.

## Environment Variables

```env
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id
DASHBOARD_PORT=8080
```

## Docker

```bash
docker build -t ironwall-backend .
docker run -p 8001:8001 ironwall-backend
```
