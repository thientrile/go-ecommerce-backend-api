# Development Guide

This guide covers the development workflow, project structure, and common tasks.

## Project Structure

```
go-ecommerce-backend-api/
├── cmd/                    # Application entry points
│   ├── server/            # Main server application
│   ├── cli/               # CLI tools
│   └── swag/docs/         # Generated Swagger documentation
├── configs/               # Configuration files
│   ├── local.yaml         # Local development config
│   ├── docker.yaml        # Docker environment config
│   └── production.yaml    # Production config
├── docs/                  # Documentation
├── internal/              # Private application code
│   ├── controller/        # HTTP handlers
│   ├── service/           # Business logic
│   ├── repo/              # Data access layer
│   ├── model/             # Data models
│   ├── middleware/        # HTTP middleware
│   └── initialize/        # Application initialization
├── nginx/                 # Nginx configuration
├── pkg/                   # Public packages
├── scripts/               # Utility scripts
├── sql/                   # Database files
├── storage/               # Storage directories
└── tests/                 # Test files
```

## Development Workflow

### 1. Setup Development Environment

```bash
# Initialize project dependencies
make init

# Setup domains for development
make setup-multi-domain

# Start local development
make dev-local
```

### 2. Database Management

```bash
# Create new migration
make goose-create name=create_users_table

# Run migrations
make goose-up

# Rollback migration
make goose-down

# Generate SQL queries
make sqlgen
```

### 3. API Documentation

```bash
# Generate Swagger documentation
make swag

# View documentation
# Local: http://crm.shopdev.test/swagger/index.html
# Docker: http://crm.shopdev.com/swagger/index.html
```

### 4. Testing

```bash
# Run tests
go test ./...

# Test API endpoints
make test-domains

# Check backend health
make check-backends
```

## Configuration Management

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENV` | Environment mode (local/docker) | `local` |

### Configuration Files

- **`configs/local.yaml`**: Local development settings
- **`configs/docker.yaml`**: Docker environment settings
- **`configs/production.yaml`**: Production settings

Example configuration:
```yaml
server:
  port: 8002
  mode: dev

mysql:
  host: 127.0.0.1
  port: 3306
  username: root
  password: strongpassword123!
  dbname: shopDevgo

logger:
  log_level: debug
  file_log_path: "./storage/logs"
```

## API Development

### Creating New Endpoints

1. **Define Model** (`internal/model/`)
```go
type User struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

2. **Create Repository** (`internal/repo/`)
```go
func (r *UserRepo) CreateUser(user *model.User) error {
    // Database operations
}
```

3. **Implement Service** (`internal/service/`)
```go
func (s *UserService) CreateUser(user *model.User) error {
    // Business logic
    return s.userRepo.CreateUser(user)
}
```

4. **Add Controller** (`internal/controller/`)
```go
// @Summary Create User
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.User true "User data"
// @Success 200 {object} response.ResponseData
// @Router /users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
    // Handle HTTP request
}
```

5. **Register Routes** (`internal/routers/`)
```go
func (r *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
    userRouter := Router.Group("/users")
    {
        userRouter.POST("", c.userController.CreateUser)
    }
}
```

### Database Operations

#### Using SQLC
1. Write SQL queries in `sql/queries/`
2. Generate Go code: `make sqlgen`
3. Use generated functions in repositories

#### Using Goose Migrations
1. Create migration: `make goose-create name=migration_name`
2. Edit SQL files in `sql/schema/`
3. Apply migrations: `make goose-up`

## Debugging

### Logging
Logs are stored in `storage/logs/` with different levels:
- `debug`: Detailed information
- `info`: General information
- `warn`: Warning messages
- `error`: Error messages

### Health Checks
- **Local**: http://crm.shopdev.test/health
- **Docker**: http://crm.shopdev.com/health

### Database Connection
Check database connectivity in logs:
```
Server Port:: 8002
Server Host:: [jwt-key]
```

## Docker Development

### Services
- **MySQL**: Database server
- **Redis**: Caching and sessions
- **Kafka**: Message broker
- **Nginx**: Reverse proxy
- **Backend**: Go application

### Commands
```bash
# Start all services
docker-compose up -d

# Start specific service
docker-compose up -d mysql

# View logs
docker-compose logs backend

# Execute commands in container
docker-compose exec backend sh
```

## Production Deployment

### Build Application
```bash
# Build binary
make build

# Build Docker image
make docker-build
```

### Environment Setup
1. Set `ENV=production`
2. Use `configs/production.yaml`
3. Configure production database
4. Set up SSL certificates
5. Configure production nginx

## Common Issues

### Port Conflicts
- Change ports in `docker-compose.yml`
- Update configuration files accordingly

### Database Connection
- Ensure MySQL is running
- Check connection parameters
- Verify database exists

### Nginx Routing
- Check nginx configuration
- Verify upstream definitions
- Test proxy settings

### Domain Resolution
- Ensure hosts file is updated
- Flush DNS cache
- Check firewall settings
