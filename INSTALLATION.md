# IronWall Installation Guide

## Prerequisites

### For Linux Server (Production)
- Ubuntu 20.04+ / Debian 11+ / CentOS 8+
- Root or sudo access
- Minimum 2GB RAM
- 10GB free disk space

### For Development (Windows/Mac/Linux)
- Go 1.21+
- Python 3.11+
- Node.js 20+
- npm or yarn

## Quick Installation (Linux)

### One-Line Install

```bash
curl -sSL https://raw.githubusercontent.com/apringutawa/ironwall/main/scripts/install.sh | sudo bash
```

Or using wget:

```bash
wget -qO- https://raw.githubusercontent.com/apringutawa/ironwall/main/scripts/install.sh | sudo bash
```

### Manual Installation

```bash
# Clone repository
git clone https://github.com/apringutawa/ironwall.git
cd ironwall

# Build CLI (requires Go)
go build -o ironwall cmd/ironwall/main.go
sudo mv ironwall /usr/local/bin/

# Run installer
sudo ironwall install
```

## Development Setup

### 1. Backend Setup

```bash
cd ironwall/backend

# Create virtual environment (recommended)
python -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install fastapi uvicorn sqlalchemy pydantic pydantic-settings python-multipart aiosqlite httpx

# Create .env file
cp .env.example .env

# Run backend
uvicorn app.main:app --host 0.0.0.0 --port 8001
```

Backend will be available at: http://localhost:8001
API docs: http://localhost:8001/docs

### 2. Frontend Setup

```bash
cd ironwall/frontend

# Install dependencies
npm install

# Create .env.local
cp .env.example .env.local

# Run development server
npm run dev
```

Frontend will be available at: http://localhost:3000

### 3. CLI Development

```bash
cd ironwall

# Download Go dependencies
go mod download

# Build CLI
go build -o ironwall cmd/ironwall/main.go

# Test CLI
./ironwall --help
```

## Docker Installation

### Using Docker Compose

```bash
cd ironwall
docker-compose up --build
```

Services:
- Backend: http://localhost:8001
- Frontend: http://localhost:3000

### Manual Docker Build

```bash
# Build backend
cd backend
docker build -t ironwall-backend .
docker run -p 8001:8001 ironwall-backend

# Build frontend
cd frontend
docker build -t ironwall-frontend .
docker run -p 3000:3000 ironwall-frontend
```

## Configuration

### Backend Configuration (.env)

```env
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_CHAT_ID=your_chat_id_here
DASHBOARD_PORT=8080
```

### Frontend Configuration (.env.local)

```env
NEXT_PUBLIC_API_URL=http://localhost:8001
```

## Verification

### Test Backend

```bash
curl http://localhost:8001/health
# Expected: {"status":"healthy"}

curl http://localhost:8001/api/v1/status
# Expected: JSON with system status
```

### Test Frontend

Open browser: http://localhost:3000

You should see the IronWall dashboard.

### Test CLI

```bash
ironwall status
ironwall scan --dry-run
ironwall protect --dry-run
```

## Troubleshooting

### Backend Issues

**Error: Module not found**
```bash
pip install --upgrade pip
pip install -r requirements.txt
```

**Error: Port 8001 already in use**
```bash
# Change port in command
uvicorn app.main:app --port 8002
```

### Frontend Issues

**Error: Cannot find module**
```bash
rm -rf node_modules package-lock.json
npm install
```

**Error: Build failed**
```bash
npm run build
# Check error messages and fix TypeScript errors
```

### CLI Issues

**Error: go: command not found**
- Install Go from https://golang.org/dl/

**Error: package not found**
```bash
go mod tidy
go mod download
```

## Production Deployment

### 1. Backend (Systemd Service)

Create `/etc/systemd/system/ironwall-backend.service`:

```ini
[Unit]
Description=IronWall Backend API
After=network.target

[Service]
Type=simple
User=ironwall
WorkingDirectory=/opt/ironwall/backend
Environment="PATH=/opt/ironwall/venv/bin"
ExecStart=/opt/ironwall/venv/bin/uvicorn app.main:app --host 0.0.0.0 --port 8001
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable ironwall-backend
sudo systemctl start ironwall-backend
```

### 2. Frontend (PM2 or Systemd)

Build production:
```bash
cd frontend
npm run build
npm start
```

Or use PM2:
```bash
npm install -g pm2
pm2 start npm --name "ironwall-frontend" -- start
pm2 save
pm2 startup
```

### 3. Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location /api/ {
        proxy_pass http://localhost:8001/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Security Notes

⚠️ **Important Security Considerations:**

1. **Change default ports** in production
2. **Use HTTPS** with SSL certificates (Let's Encrypt)
3. **Configure firewall** to allow only necessary ports
4. **Set strong passwords** for any authentication
5. **Keep dependencies updated** regularly
6. **Backup configuration** before making changes
7. **Test in staging** before production deployment

## Support

- Documentation: https://github.com/apringutawa/ironwall/wiki
- Issues: https://github.com/apringutawa/ironwall/issues
- Repository: https://github.com/apringutawa/ironwall

## Next Steps

After installation:

1. Run `ironwall status` to check system health
2. Run `ironwall scan` to perform initial security scan
3. Configure alerts in dashboard settings
4. Review and customize security modules
5. Set up automated backups

For detailed usage, see [README.md](README.md)
