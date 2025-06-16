# Multi-Domain Development Guide

This guide explains how to set up and use the multi-domain development environment for seamless switching between local and Docker environments.

## Overview

The application supports two development modes:

- **Local Development**: `crm.shopdev.test` - Run Go backend locally with nginx in Docker
- **Docker Environment**: `crm.shopdev.com` - Full Docker stack deployment

## Environment Detection

The application automatically detects the environment using the `ENV` environment variable:

```bash
ENV=local   # Local development (default)
ENV=docker  # Docker environment
```

## Domain Configuration

| Environment | Domain | Config File | Backend Location |
|-------------|--------|-------------|------------------|
| Local | `crm.shopdev.test` | `configs/local.yaml` | Host machine |
| Docker | `crm.shopdev.com` | `configs/docker.yaml` | Docker container |

## Quick Start

### 1. Setup Domains

Run as Administrator in PowerShell:

```bash
# Setup both domains
make setup-multi-domain

# Or setup individually
make setup-local-domain   # Only crm.shopdev.test
make setup-docker-domain  # Only crm.shopdev.com
```

### 2. Start Local Development

```bash
# Start supporting services (nginx, mysql, redis)
make dev-local

# In another terminal, run Go backend
make run-local
# OR manually with environment variable
$env:ENV="local"; go run cmd/server/main.go
```

### 3. Start Docker Environment

```bash
# Start full Docker stack
make dev-docker
```

## Access Points

### Local Development (`crm.shopdev.test`)
- 🌐 **Main API**: http://crm.shopdev.test
- 📊 **Swagger UI**: http://crm.shopdev.test/swagger/index.html
- 🏥 **Health Check**: http://crm.shopdev.test/health

### Docker Environment (`crm.shopdev.com`)
- 🌐 **Main API**: http://crm.shopdev.com
- 📊 **Swagger UI**: http://crm.shopdev.com/swagger/index.html
- 🏥 **Health Check**: http://crm.shopdev.com/health
- 🎛️ **Kafka UI**: http://crm.shopdev.com/kafka-ui/

## How It Works

### Environment Detection
The application reads the `ENV` environment variable in `internal/initialize/loadconfig.init.go`:

```go
env := os.Getenv("ENV")
if env == "" {
    env = "local"
}
```

### Swagger Host Detection
The Swagger configuration in `internal/initialize/swag.init.go` automatically sets the host:

```go
env := os.Getenv("ENV")
if env == "docker" {
    host = "crm.shopdev.com"
} else {
    host = "crm.shopdev.test"
}
```

### Nginx Routing
- **Local**: nginx proxies `crm.shopdev.test` → `host.docker.internal:8002`
- **Docker**: nginx proxies `crm.shopdev.com` → `go-ecommerce-backend:8001`

## Troubleshooting

### Domain Not Accessible
1. Ensure you ran the domain setup as Administrator
2. Check Windows hosts file: `C:\Windows\System32\drivers\etc\hosts`
3. Flush DNS cache: `ipconfig /flushdns`

### Services Not Starting
```bash
# Check Docker services
docker-compose ps

# Check backend health
make test-domains

# View logs
make logs
```

### Wrong Swagger Host
Ensure the `ENV` variable is set correctly:
```bash
# Check current environment
echo $env:ENV

# Set environment
$env:ENV="local"   # or "docker"
```

## Available Make Commands

| Command | Description |
|---------|-------------|
| `make setup-multi-domain` | Setup both domains |
| `make setup-local-domain` | Setup local domain only |
| `make setup-docker-domain` | Setup Docker domain only |
| `make remove-domains` | Remove all domains |
| `make dev-local` | Start local development environment |
| `make dev-docker` | Start Docker environment |
| `make run-local` | Run Go server with ENV=local |
| `make test-domains` | Test domain accessibility |
| `make check-backends` | Check backend health |

## Configuration Files

- `configs/local.yaml` - Local development configuration
- `configs/docker.yaml` - Docker environment configuration
- `nginx/conf.d/local-server.conf` - Local nginx configuration
- `nginx/conf.d/docker-server.conf` - Docker nginx configuration
- `scripts/setup-multi-domain.ps1` - Domain setup script
