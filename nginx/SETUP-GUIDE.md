# Hosts File Configuration for Go E-commerce Backend API

## Overview
This guide helps you configure your local hosts file to access the Go E-commerce Backend API using custom domains both locally and in Docker environments.

## Domains Configuration

### Primary Domain (Docker & Local)
- **Docker**: `http://crm.shopdev.com/`
- **Local Development**: `http://crm.shopdev.test/`

## Windows Hosts File Configuration

### Location
The hosts file is located at: `C:\Windows\System32\drivers\etc\hosts`

### Required Entries
Add the following lines to your hosts file:

```
# Go E-commerce Backend API
127.0.0.1    crm.shopdev.com
127.0.0.1    crm.shopdev.test
::1          crm.shopdev.com
::1          crm.shopdev.test
```

### How to Edit Hosts File on Windows

1. **Open Command Prompt as Administrator**
   - Press `Win + R`, type `cmd`, then press `Ctrl + Shift + Enter`

2. **Edit the hosts file**
   ```cmd
   notepad C:\Windows\System32\drivers\etc\hosts
   ```

3. **Add the domain entries** (see above)

4. **Save and close the file**

5. **Flush DNS cache**
   ```cmd
   ipconfig /flushdns
   ```

## Testing Configuration

### 1. Test Domain Resolution
```cmd
ping crm.shopdev.com
ping crm.shopdev.test
```

### 2. Test with curl (if available)
```cmd
curl -I http://crm.shopdev.com/health
curl -I http://crm.shopdev.test/health
```

### 3. Test in Browser
- Docker: `http://crm.shopdev.com/`
- Local: `http://crm.shopdev.test/`

## Environment-Specific Access

### Docker Environment
When running with `docker-compose up`:
- **URL**: `http://crm.shopdev.com/`
- **Port**: 80 (nginx) → 8002 (Go app)
- **Internal**: nginx proxies to `crm.shopdev.com:8002`

### Local Development
When running Go app locally:
- **URL**: `http://crm.shopdev.test/`
- **Port**: 80 (nginx) → 8001 (Go app)
- **Direct**: `http://localhost:8001/` (without nginx)

## Nginx Configuration Files

### Available Configurations
1. **`nginx/conf.d/default.conf`** - Universal configuration
2. **`nginx/conf.d/docker.conf`** - Docker-optimized
3. **`nginx/conf.d/local.conf`** - Local development

### Switching Configurations

#### For Docker (default)
Uses `nginx/conf.d/docker.conf` automatically via docker-compose

#### For Local Development
1. Copy local configuration:
   ```cmd
   copy nginx\conf.d\local.conf nginx\conf.d\default.conf
   ```

2. Restart nginx service

## Troubleshooting

### Common Issues

1. **Domain not resolving**
   - Check hosts file entries
   - Flush DNS cache: `ipconfig /flushdns`
   - Restart browser

2. **Connection refused**
   - Verify Go application is running
   - Check port mappings in docker-compose.yml
   - Verify nginx configuration

3. **502 Bad Gateway**
   - Go application not responding
   - Check container logs: `docker-compose logs crm.shopdev.com`
   - Verify upstream configuration in nginx

4. **CORS Issues**
   - Check nginx CORS headers configuration
   - Verify API endpoint configuration

### Debugging Commands

```cmd
# Check running containers
docker-compose ps

# View logs
docker-compose logs nginx
docker-compose logs crm.shopdev.com

# Test nginx configuration
docker-compose exec nginx nginx -t

# Reload nginx configuration
docker-compose exec nginx nginx -s reload
```

## Security Notes

### Development vs Production
- Local configuration has relaxed CORS settings
- Docker configuration includes security headers
- Never use development settings in production

### HTTPS (Future Enhancement)
To enable HTTPS:
1. Generate SSL certificates
2. Place in `nginx/certs/` directory
3. Update nginx configuration to listen on port 443
4. Update docker-compose.yml to expose port 443

## API Endpoints

### Health Check
- Docker: `http://crm.shopdev.com/health`
- Local: `http://crm.shopdev.test/health`

### API Base
- Docker: `http://crm.shopdev.com/api/`
- Local: `http://crm.shopdev.test/api/`

### Swagger Documentation
- Docker: `http://crm.shopdev.com/swagger/`
- Local: `http://crm.shopdev.test/swagger/`

### Static Files
- Docker: `http://crm.shopdev.com/static/`
- Local: `http://crm.shopdev.test/static/`
