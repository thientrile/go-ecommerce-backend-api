# Enhanced Logging System

## Tổng quan
Hệ thống logging được cải thiện để dễ đọc hơn trên console với màu sắc và format đẹp mắt, đồng thời vẫn ghi file logs dạng JSON cho việc phân tích sau này.

## Tính năng chính

### 1. Console Output với màu sắc
- **DEBUG**: Cyan `[DBG]`
- **INFO**: Green `[INF]`  
- **WARN**: Yellow `[WRN]`
- **ERROR**: Red `[ERR]`
- **FATAL**: Red Bold `[FTL]`

### 2. Dual Output
- **Console**: Format dễ đọc với màu sắc và icons
- **File**: JSON format để phân tích và monitoring

### 3. Helper Methods
- `PrintStartupBanner()`: Banner đẹp khi khởi động
- `LogInitStep()`: Log từng bước initialization 
- `LogDBConnection()`: Log kết nối database
- `LogHTTPRequest()`: Log HTTP requests với icons
- `LogShutdown()`: Log khi shutdown

## Ví dụ sử dụng

### Trong Initialization
```go
// Khởi động
global.Logger.PrintStartupBanner("my-service", "v1.0.0", "8080", "dev")

// Bắt đầu init component
global.Logger.LogInitStart("Database")

// Kết quả init
global.Logger.LogInitStep("Database", true, nil)  // success
global.Logger.LogInitStep("Redis", false, err)    // failed
```

### Trong Business Logic
```go
// Info với context
global.Logger.Info("🔐 User login attempt", 
    zap.String("email", email),
    zap.String("ip", clientIP))

// Error với details
global.Logger.Error("❌ Payment processing failed", 
    zap.String("orderID", orderID),
    zap.Float64("amount", amount),
    zap.Error(err))

// Database operations
global.Logger.LogDBConnection("MySQL", "localhost:3306", true, nil)
```

### HTTP Request Logging
```go
// Middleware tự động log với format:
// ✅ GET /api/users from 192.168.1.100 (23ms, 1234 bytes)
// ❌ POST /api/login from 192.168.1.100 (156ms, 45 bytes)
```

## Log Levels và khi nào sử dụng

### TRACE
- Chi tiết cực kỳ nhỏ (thường không dùng trong production)
- VD: "Entering function X with params Y"

### DEBUG  
- Thông tin debug cho development
- VD: "Database query: SELECT * FROM users WHERE id = 123"

### INFO
- Thông tin general về flow của ứng dụng
- VD: "User logged in", "Order created", "Email sent"

### WARN
- Cảnh báo về tình huống bất thường nhưng không phải lỗi
- VD: "Retry attempt 3/5", "Cache miss", "Deprecated API used"

### ERROR
- Lỗi trong business logic nhưng ứng dụng vẫn chạy được
- VD: "Failed to send email", "Database timeout", "Invalid user input"

### FATAL
- Lỗi nghiêm trọng khiến ứng dụng không thể tiếp tục
- VD: "Cannot connect to database", "Missing critical config"
- Sẽ tự động shutdown ứng dụng

## Icons được sử dụng

- 🚀 Server startup
- 🔧 Initialization 
- ✅ Success operations
- ❌ Failed operations  
- ⚠️ Warnings
- 💾 Database operations
- 🌐 HTTP requests
- 🔐 Authentication
- 📧 Email operations
- 💰 Payment operations
- 🛑 Shutdown

## Best Practices

### 1. Sử dụng structured logging
```go
// Good ✅
global.Logger.Info("User created", 
    zap.String("userID", userID),
    zap.String("email", email))

// Bad ❌  
global.Logger.Info(fmt.Sprintf("User %s created with email %s", userID, email))
```

### 2. Thêm context hữu ích
```go
global.Logger.Error("Database connection failed",
    zap.String("host", dbHost),
    zap.Int("port", dbPort),
    zap.String("database", dbName),
    zap.Duration("timeout", timeout),
    zap.Error(err))
```

### 3. Sử dụng log level phù hợp
```go
// Init errors -> FATAL (vì không thể chạy được)
global.Logger.Fatal("Cannot initialize database", zap.Error(err))

// Business errors -> ERROR (ứng dụng vẫn chạy được)  
global.Logger.Error("Failed to process payment", zap.Error(err))

// Warnings -> WARN
global.Logger.Warn("High memory usage detected", zap.Float64("usage", 85.5))
```

### 4. Performance considerations
- Sử dụng DEBUG level cho detailed logs
- Tránh log quá nhiều trong tight loops
- Sử dụng sampling cho high-frequency events

## Configuration

Trong config file, có thể điều chỉnh:
```yaml
logger:
  log_level: "info"        # trace, debug, info, warn, error, fatal
  file_log_path: "./storage/logs"
  file_log_name: "app"
  max_size: 10             # MB
  max_backups: 5
  max_age: 30              # days  
  compress: true
```

## Monitoring và Alerts

File logs JSON có thể được sử dụng với:
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Grafana + Loki
- Prometheus alerts
- Custom log analysis tools

Example JSON log entry:
```json
{
  "level": "error",
  "time": "2025-06-19T15:04:05.123+07:00",
  "caller": "service/user.go:45",
  "msg": "Failed to create user",
  "userID": "12345",
  "email": "user@example.com",
  "error": "duplicate email address",
  "stacktrace": "..."
}
```
