# 📊 Monitoring Guide - Go E-commerce Backend API

## 🎯 Tổng quan

Hệ thống monitoring của Go E-commerce Backend API sử dụng stack Prometheus + Grafana để theo dõi hiệu suất và trạng thái của ứng dụng. Hệ thống bao gồm:

- **Prometheus**: Thu thập và lưu trữ metrics
- **Grafana**: Trực quan hóa dữ liệu và dashboard
- **Exporters**: Thu thập metrics từ các dịch vụ khác nhau

## 🏗️ Kiến trúc Monitoring

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Go Backend    │    │   MySQL DB      │    │   Redis Cache   │
│   (Port 8002)   │    │   (Port 3306)   │    │   (Port 6379)   │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │                      │                      │
┌─────────▼───────┐    ┌─────────▼───────┐    ┌─────────▼───────┐
│ Node Exporter   │    │ MySQL Exporter  │    │ Redis Exporter  │
│   (Port 9100)   │    │   (Port 9104)   │    │   (Port 9121)   │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │              ┌───────▼────────┐             │
          │              │ Kafka Exporter │             │
          │              │   (Port 9308)  │             │
          │              └───────┬────────┘             │
          │                      │                      │
          └──────────────────────▼──────────────────────┘
                                 │
                        ┌─────────▼───────┐
                        │   Prometheus    │
                        │   (Port 9090)   │
                        └─────────┬───────┘
                                  │
                        ┌─────────▼───────┐
                        │    Grafana      │
                        │   (Port 3000)   │
                        └─────────────────┘
```

## 🚀 Khởi chạy hệ thống Monitoring

### 1. Development Mode

```bash
# Khởi chạy tất cả dịch vụ cơ bản + monitoring
docker-compose --profile dev up -d

# Hoặc chỉ khởi chạy monitoring services
docker-compose --profile dev up -d prometheus grafana mysql-exporter redis-exporter node-exporter kafka-exporter
```

### 2. Production Mode

```bash
# Khởi chạy full stack production
docker-compose --profile prod up -d

# Hoặc chỉ monitoring
docker-compose --profile prod up -d prometheus grafana mysql-exporter redis-exporter node-exporter kafka-exporter
```

### 3. Kiểm tra trạng thái

```bash
# Xem trạng thái tất cả containers
docker-compose ps

# Xem logs của monitoring services
docker-compose logs prometheus
docker-compose logs grafana
docker-compose logs mysql-exporter
```

## 🔗 Truy cập các dịch vụ

| Dịch vụ | URL | Thông tin đăng nhập |
|---------|-----|-------------------|
| **Grafana** | http://localhost:3000 | admin / admin123 |
| **Prometheus** | http://localhost:9090 | Không cần đăng nhập |
| **Kafka UI** | http://localhost:8080 | Không cần đăng nhập |

## 📈 Metrics và Exporters

### 1. Node Exporter (System Metrics)
- **Port**: 9100
- **Metrics**: CPU, Memory, Disk, Network
- **Endpoint**: http://localhost:9100/metrics

### 2. MySQL Exporter (Database Metrics)
- **Port**: 9104
- **Metrics**: Connections, Queries, InnoDB status
- **Endpoint**: http://localhost:9104/metrics

### 3. Redis Exporter (Cache Metrics)
- **Port**: 9121
- **Metrics**: Memory usage, Connections, Commands
- **Endpoint**: http://localhost:9121/metrics

### 4. Kafka Exporter (Message Queue Metrics)
- **Port**: 9308
- **Metrics**: Topics, Partitions, Consumer lag
- **Endpoint**: http://localhost:9308/metrics

## 📊 Cấu hình Prometheus

File cấu hình: `monitoring/prometheus/prometheus.yml`

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'node-exporter'
    static_configs:
      - targets: ['node-exporter:9100']

  - job_name: 'mysql-exporter'
    static_configs:
      - targets: ['mysql-exporter:9104']

  - job_name: 'redis-exporter'
    static_configs:
      - targets: ['redis-exporter:9121']

  - job_name: 'kafka-exporter'
    static_configs:
      - targets: ['kafka-exporter:9308']
```

## 🎨 Grafana Dashboard

### 1. Data Source Setup

1. Truy cập Grafana: http://localhost:3000
2. Đăng nhập: admin/admin123
3. Thêm Data Source:
   - Type: Prometheus
   - URL: http://prometheus:9090
   - Access: Server (default)

### 2. Import Dashboard

Có thể import các dashboard từ Grafana.com:

- **Node Exporter Full**: ID 1860
- **MySQL Overview**: ID 7362
- **Redis Dashboard**: ID 763
- **Kafka Overview**: ID 7589

### 3. Custom Dashboard

Tạo dashboard tùy chỉnh để monitor:
- API Response Time
- Request Rate
- Error Rate
- Database Connection Pool
- Cache Hit Rate

## 🔧 Queries phổ biến

### System Metrics
```promql
# CPU Usage
100 - (avg by(instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)

# Memory Usage
(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100

# Disk Usage
100 - ((node_filesystem_avail_bytes * 100) / node_filesystem_size_bytes)
```

### MySQL Metrics
```promql
# Active Connections
mysql_global_status_threads_connected

# Queries per second
rate(mysql_global_status_questions[5m])

# InnoDB Buffer Pool Usage
mysql_global_status_innodb_buffer_pool_pages_data / mysql_global_status_innodb_buffer_pool_pages_total * 100
```

### Redis Metrics
```promql
# Memory Usage
redis_memory_used_bytes

# Connected Clients
redis_connected_clients

# Commands per second
rate(redis_commands_processed_total[5m])
```

## 🚨 Alerting Rules

### 1. Cấu hình Alert Rules

File: `monitoring/prometheus/alert-rules.yml`

```yaml
groups:
  - name: system-alerts
    rules:
      - alert: HighCPUUsage
        expr: 100 - (avg by(instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100) > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High CPU usage detected"

      - alert: HighMemoryUsage
        expr: (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100 > 85
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage detected"

  - name: database-alerts
    rules:
      - alert: MySQLDown
        expr: mysql_up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "MySQL is down"

      - alert: HighMySQLConnections
        expr: mysql_global_status_threads_connected > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High MySQL connections"
```

### 2. Notification Channels

Cấu hình thông báo qua:
- Email
- Slack
- Webhook
- PagerDuty

## 🔍 Troubleshooting

### 1. Prometheus không thu thập được metrics

```bash
# Kiểm tra targets trong Prometheus
curl http://localhost:9090/api/v1/targets

# Kiểm tra logs
docker-compose logs prometheus
```

### 2. Grafana không kết nối được Prometheus

```bash
# Kiểm tra network connectivity
docker exec -it go-ecommerce-grafana nslookup prometheus

# Test connection
curl -v http://prometheus:9090/api/v1/query?query=up
```

### 3. Exporter không hoạt động

```bash
# MySQL Exporter
docker-compose logs mysql-exporter
curl http://localhost:9104/metrics

# Redis Exporter
docker-compose logs redis-exporter
curl http://localhost:9121/metrics
```

## 📝 Best Practices

### 1. Retention Policy

- **Prometheus**: Giữ dữ liệu 200 giờ (8+ ngày)
- **Grafana**: Backup dashboard definitions
- **Long-term storage**: Sử dụng remote storage như Thanos

### 2. Resource Management

```yaml
# Giới hạn resource cho containers
resources:
  limits:
    memory: 512M
    cpus: '0.5'
  reservations:
    memory: 256M
    cpus: '0.25'
```

### 3. Security

- Thay đổi password mặc định của Grafana
- Sử dụng TLS cho production
- Restrict network access
- Enable authentication

## 🔄 Backup và Restore

### 1. Backup Grafana

```bash
# Backup Grafana data
docker run --rm -v go-ecommerce-backend-api_grafana-data:/from alpine tar czf - -C /from . > grafana-backup.tar.gz

# Backup dashboards
curl -H "Authorization: Bearer <api-key>" http://localhost:3000/api/search?type=dash-db | jq
```

### 2. Backup Prometheus

```bash
# Backup Prometheus data
docker run --rm -v go-ecommerce-backend-api_prometheus-data:/from alpine tar czf - -C /from . > prometheus-backup.tar.gz
```

## 📚 Tài liệu tham khảo

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Node Exporter](https://github.com/prometheus/node_exporter)
- [MySQL Exporter](https://github.com/prometheus/mysqld_exporter)
- [Redis Exporter](https://github.com/oliver006/redis_exporter)

---

**Lưu ý**: Đảm bảo tất cả các dịch vụ dependencies (MySQL, Redis, Kafka) đang chạy trước khi khởi động monitoring stack.
