# Deployment Guide

This guide covers deployment strategies for different environments.

## Deployment Overview

The application supports multiple deployment scenarios:

- **Local Development**: For development and testing
- **Docker Development**: Full containerized environment
- **Production**: Production-ready deployment

## Prerequisites

### System Requirements
- **Go**: 1.21 or higher
- **Docker**: 20.10 or higher
- **Docker Compose**: 2.0 or higher
- **MySQL**: 8.0 or higher
- **Redis**: 7.0 or higher

### Required Tools
```bash
# Install Go tools
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

## Local Development Deployment

### Quick Start
```bash
# 1. Clone repository
git clone <repository-url>
cd go-ecommerce-backend-api

# 2. Initialize project
make init

# 3. Setup domains
make setup-multi-domain

# 4. Start local environment
make dev-local

# 5. Run application
make run-local
```

### Manual Setup
```bash
# 1. Install dependencies
go mod download

# 2. Generate Swagger docs
make swag

# 3. Start database
docker-compose up -d mysql redis

# 4. Run migrations
make goose-up

# 5. Start application
ENV=local go run cmd/server/main.go
```

## Docker Development Deployment

### Quick Start
```bash
# 1. Setup domains
make setup-multi-domain

# 2. Start full Docker environment
make dev-docker
```

### Manual Docker Setup
```bash
# 1. Build images
docker-compose build

# 2. Start all services
docker-compose up -d

# 3. Check status
docker-compose ps

# 4. View logs
docker-compose logs -f backend
```

### Docker Services

| Service | Port | Description |
|---------|------|-------------|
| backend | 8002 | Go application |
| nginx | 80 | Reverse proxy |
| mysql | 3306 | Database |
| redis | 6379 | Cache/Sessions |
| kafka | 9092 | Message broker |
| kafka-ui | 8080 | Kafka management |

## Production Deployment

### Prerequisites
- **Server**: Linux Ubuntu 20.04+ or CentOS 8+
- **Memory**: 4GB RAM minimum
- **Storage**: 20GB+ SSD
- **Network**: Static IP address
- **SSL**: Valid SSL certificates

### Production Setup

#### 1. Server Preparation
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Create application user
sudo useradd -m -s /bin/bash appuser
sudo usermod -aG docker appuser
```

#### 2. Application Deployment
```bash
# Clone repository
git clone <repository-url> /opt/go-ecommerce-backend-api
cd /opt/go-ecommerce-backend-api

# Set ownership
sudo chown -R appuser:appuser /opt/go-ecommerce-backend-api

# Switch to app user
sudo su - appuser
cd /opt/go-ecommerce-backend-api
```

#### 3. Environment Configuration
```bash
# Create production environment file
cat > .env.production << EOF
ENV=production
DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-secure-password
DB_NAME=shopDevgo
REDIS_HOST=redis
REDIS_PORT=6379
JWT_SECRET=your-jwt-secret-key
EOF

# Create production docker-compose file
cp docker-compose.yml docker-compose.prod.yml
```

#### 4. Production Docker Compose
Update `docker-compose.prod.yml`:

```yaml
version: "3.9"
name: go-ecommerce-backend-api-prod

services:
  backend:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: go-ecommerce-backend-prod
    environment:
      - ENV=production
    command: ["/crm.shopdev.com", "configs/production.yaml"]
    volumes:
      - ./configs:/configs
      - ./storage/logs:/app/storage/logs
    restart: unless-stopped
    depends_on:
      - mysql
      - redis
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    container_name: nginx-prod
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    restart: unless-stopped
    depends_on:
      - backend
    networks:
      - app-network

  mysql:
    image: mysql:8.0
    container_name: mysql-prod
    environment:
      MYSQL_ROOT_PASSWORD: "${DB_PASSWORD}"
      MYSQL_DATABASE: "${DB_NAME}"
    volumes:
      - mysql-data:/var/lib/mysql
      - ./mysql-custom.cnf:/etc/mysql/conf.d/custom.cnf
    restart: unless-stopped
    networks:
      - app-network

  redis:
    image: redis:7.2-alpine
    container_name: redis-prod
    restart: unless-stopped
    networks:
      - app-network

volumes:
  mysql-data:

networks:
  app-network:
    driver: bridge
```

#### 5. SSL Configuration
```bash
# Create SSL directory
mkdir -p nginx/ssl

# Copy SSL certificates
cp your-domain.crt nginx/ssl/
cp your-domain.key nginx/ssl/

# Update nginx configuration for SSL
```

#### 6. Database Setup
```bash
# Start database
docker-compose -f docker-compose.prod.yml up -d mysql

# Wait for MySQL to be ready
sleep 30

# Run migrations
GOOSE_DBSTRING="root:${DB_PASSWORD}@tcp(localhost:3306)/${DB_NAME}" make goose-up
```

#### 7. Application Startup
```bash
# Start all services
docker-compose -f docker-compose.prod.yml up -d

# Check status
docker-compose -f docker-compose.prod.yml ps

# View logs
docker-compose -f docker-compose.prod.yml logs -f backend
```

### Production Monitoring

#### Health Checks
```bash
# Application health
curl -f http://localhost/health || exit 1

# Database connectivity
docker-compose -f docker-compose.prod.yml exec mysql mysqladmin ping

# Redis connectivity
docker-compose -f docker-compose.prod.yml exec redis redis-cli ping
```

#### Log Monitoring
```bash
# Application logs
docker-compose -f docker-compose.prod.yml logs -f backend

# Nginx logs
docker-compose -f docker-compose.prod.yml logs -f nginx

# Database logs
docker-compose -f docker-compose.prod.yml logs -f mysql
```

### Backup and Recovery

#### Database Backup
```bash
# Create backup script
cat > scripts/backup-db.sh << EOF
#!/bin/bash
BACKUP_DIR="./storage/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="backup_${DATE}.sql"

mkdir -p $BACKUP_DIR
docker-compose -f docker-compose.prod.yml exec mysql mysqldump -u root -p${DB_PASSWORD} ${DB_NAME} > ${BACKUP_DIR}/${BACKUP_FILE}
gzip ${BACKUP_DIR}/${BACKUP_FILE}
echo "Backup created: ${BACKUP_DIR}/${BACKUP_FILE}.gz"
EOF

chmod +x scripts/backup-db.sh
```

#### Application Backup
```bash
# Backup application data
tar -czf backup-$(date +%Y%m%d).tar.gz \
  configs/ \
  storage/logs/ \
  storage/backups/ \
  docker-compose.prod.yml \
  .env.production
```

### Updates and Maintenance

#### Application Updates
```bash
# 1. Backup current version
./scripts/backup-db.sh

# 2. Pull latest changes
git pull origin main

# 3. Rebuild and restart
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml build
docker-compose -f docker-compose.prod.yml up -d

# 4. Run migrations if needed
make goose-up

# 5. Verify deployment
curl -f http://localhost/health
```

#### System Maintenance
```bash
# Clean unused Docker resources
docker system prune -f

# Update Docker images
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d

# Rotate logs
find storage/logs/ -name "*.log" -mtime +30 -delete
```

## Troubleshooting

### Common Issues

#### Port Conflicts
```bash
# Check port usage
netstat -tulpn | grep :80
netstat -tulpn | grep :3306

# Change ports in docker-compose file if needed
```

#### Database Connection Issues
```bash
# Check MySQL status
docker-compose -f docker-compose.prod.yml logs mysql

# Test connection
docker-compose -f docker-compose.prod.yml exec mysql mysql -u root -p

# Reset password if needed
docker-compose -f docker-compose.prod.yml exec mysql mysql -u root -p -e "ALTER USER 'root'@'%' IDENTIFIED BY 'new-password';"
```

#### Application Not Starting
```bash
# Check application logs
docker-compose -f docker-compose.prod.yml logs backend

# Check configuration
docker-compose -f docker-compose.prod.yml exec backend cat configs/production.yaml

# Restart services
docker-compose -f docker-compose.prod.yml restart backend
```

### Performance Tuning

#### MySQL Optimization
```bash
# Update MySQL configuration
cp mysql-prod.cnf mysql-custom.cnf

# Restart MySQL
docker-compose -f docker-compose.prod.yml restart mysql
```

#### Application Optimization
- Enable Go build optimizations
- Configure proper logging levels
- Set up connection pooling
- Implement caching strategies

#### Nginx Optimization
- Enable gzip compression
- Set up proper caching headers
- Configure rate limiting
- Optimize SSL settings
