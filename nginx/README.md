# Nginx Configuration for Go E-commerce Backend API

## Quick Start

### 1. Configure Hosts File (Required)
Run as Administrator:
```powershell
cd scripts
.\configure-hosts.ps1
```

### 2. Start with Docker (Recommended)
```bash
docker-compose up -d
```
Access: `http://crm.shopdev.com/`

### 3. Start Local Development
```powershell
# Install nginx (if not installed)
choco install nginx

# Start Go application
go run cmd/server/main.go

# Start nginx with local config
cd scripts
.\start-local.ps1 -Start
```
Access: `http://crm.shopdev.test/`

## Configuration Files

| File | Purpose |
|------|---------|
| `nginx/conf.d/default.conf` | Universal configuration (both environments) |
| `nginx/conf.d/docker.conf` | Docker-optimized configuration |
| `nginx/conf.d/local.conf` | Local development configuration |
| `nginx/nginx.conf` | Main nginx configuration |

## URLs

### Docker Environment
- Application: `http://crm.shopdev.com/`
- API: `http://crm.shopdev.com/api/`
- Health: `http://crm.shopdev.com/health`
- Swagger: `http://crm.shopdev.com/swagger/`

### Local Development
- Application: `http://crm.shopdev.test/`
- API: `http://crm.shopdev.test/api/`
- Health: `http://crm.shopdev.test/health`
- Swagger: `http://crm.shopdev.test/swagger/`

## Port Mapping

| Environment | External | Nginx | Go App |
|-------------|----------|-------|---------|
| Docker | 80 | 80 | 8002 |
| Local | 80 | 80 | 8001 |

## Features

- ✅ CORS support for API endpoints
- ✅ Security headers
- ✅ Gzip compression
- ✅ Health check endpoints
- ✅ Static file serving
- ✅ WebSocket support
- ✅ Rate limiting zones
- ✅ Error page handling
- ✅ SSL-ready (certificates needed)

## Troubleshooting

### Common Commands
```bash
# Check nginx status
.\scripts\start-local.ps1 -Status

# View Docker logs
docker-compose logs nginx
docker-compose logs crm.shopdev.com

# Test nginx config
nginx -t

# Restart nginx
.\scripts\start-local.ps1 -Restart
```

### DNS Issues
```cmd
# Flush DNS cache
ipconfig /flushdns

# Test domain resolution
ping crm.shopdev.com
ping crm.shopdev.test
```

See `nginx/SETUP-GUIDE.md` for detailed setup instructions.
