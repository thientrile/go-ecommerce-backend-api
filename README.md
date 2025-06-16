# Go E-commerce Backend API

A robust, scalable e-commerce backend API built with Go, featuring multi-domain development environments, comprehensive authentication, and modern microservices architecture.

## 🚀 Features

- **Multi-Domain Development**: Seamless switching between local (`crm.shopdev.test`) and Docker (`crm.shopdev.com`) environments
- **JWT Authentication**: Secure user authentication with JWT tokens
- **RESTful API**: Clean, well-documented REST API endpoints
- **Database Management**: MySQL with migration support using Goose
- **Caching**: Redis integration for session management and caching
- **Message Queue**: Kafka integration for asynchronous processing
- **API Documentation**: Auto-generated Swagger documentation
- **Docker Support**: Full containerization with Docker Compose
- **Nginx Proxy**: Reverse proxy configuration for production
- **Comprehensive Logging**: Structured logging with different levels
- **Email Integration**: Email verification and notifications
- **Two-Factor Authentication**: Enhanced security with 2FA support

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client/Web    │    │   Mobile App    │    │   Admin Panel   │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │       Nginx Proxy         │
                    └─────────────┬─────────────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │     Go Backend API        │
                    │   (Gin Framework)         │
                    └─────────────┬─────────────┘
                                  │
          ┌───────────────────────┼───────────────────────┐
          │                       │                       │
   ┌──────┴──────┐       ┌────────┴────────┐     ┌────────┴────────┐
   │   MySQL     │       │     Redis       │     │     Kafka       │
   │ (Database)  │       │   (Cache)       │     │ (Message Queue) │
   └─────────────┘       └─────────────────┘     └─────────────────┘
```

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: MySQL 8.0
- **Cache**: Redis 7.2
- **Message Queue**: Apache Kafka
- **Documentation**: Swagger/OpenAPI
- **Authentication**: JWT
- **ORM**: SQLC + Goose Migrations
- **Containerization**: Docker & Docker Compose
- **Reverse Proxy**: Nginx
- **Logging**: Structured logging with rotation

## 📋 Prerequisites

- **Go**: 1.21 or higher
- **Docker**: 20.10 or higher
- **Docker Compose**: 2.0 or higher
- **Git**: Latest version
- **Make**: Build automation tool

## 🚀 Quick Start

### 1. Clone the Repository
```bash
git clone https://github.com/your-username/go-ecommerce-backend-api.git
cd go-ecommerce-backend-api
```

### 2. Initialize Project
```bash
# Install dependencies and tools
make init
```

### 3. Setup Development Environment
```bash
# Setup both local and Docker domains (Run as Administrator)
make setup-multi-domain

# Start local development environment
make dev-local

# In another terminal, run the Go backend
make run-local
```

### 4. Access the Application
- **Local Development**: http://crm.shopdev.test
- **Swagger Documentation**: http://crm.shopdev.test/swagger/index.html
- **Health Check**: http://crm.shopdev.test/health

### Alternative: Full Docker Environment
```bash
# Start complete Docker stack
make dev-docker

# Access at: http://crm.shopdev.com
```

## 📚 Documentation

Comprehensive documentation is available in the `docs/` directory:

- **[Multi-Domain Setup Guide](docs/multi-domain-setup.md)**: Environment configuration and domain management
- **[Development Guide](docs/development-guide.md)**: Development workflow and project structure
- **[Deployment Guide](docs/deployment-guide.md)**: Production deployment strategies
- **[API Documentation](docs/api-documentation.md)**: Complete API reference and examples

## 🔧 Available Commands

### Development Commands
```bash
make init              # Initialize project dependencies
make dev               # Start development server
make run-local         # Run server with local environment
make build             # Build application binary
make swag              # Generate Swagger documentation
```

### Database Commands
```bash
make goose-up          # Run database migrations
make goose-down        # Rollback last migration
make goose-create      # Create new migration
make sqlgen            # Generate SQL queries
```

### Docker Commands
```bash
make docker-up         # Start Docker services
make docker-down       # Stop Docker services
make docker-build      # Rebuild Docker images
make dev-docker        # Start full Docker environment
```

### Domain Management
```bash
make setup-multi-domain    # Setup both domains
make setup-local-domain    # Setup local domain only
make setup-docker-domain   # Setup Docker domain only
make remove-domains        # Remove all domains
make test-domains          # Test domain accessibility
```

### Testing & Monitoring
```bash
make test              # Run tests
make test-domains      # Test domain configuration
make check-backends    # Check backend health
make logs              # View Docker logs
```

## 🌍 Environment Configuration

The application supports multiple environments with automatic detection:

| Environment | Domain | Config File | Description |
|-------------|--------|-------------|-------------|
| **Local** | `crm.shopdev.test` | `configs/local.yaml` | Development on host machine |
| **Docker** | `crm.shopdev.com` | `configs/docker.yaml` | Full containerized environment |
| **Production** | Custom domain | `configs/production.yaml` | Production deployment |

Environment is automatically detected using the `ENV` environment variable:
```bash
ENV=local    # Local development (default)
ENV=docker   # Docker environment
ENV=production  # Production environment
```

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication:

1. **Register**: Create a new user account
2. **Verify**: Verify email with OTP code
3. **Login**: Authenticate and receive JWT token
4. **Protected Routes**: Include JWT token in Authorization header

Example:
```bash
# Login
curl -X POST http://crm.shopdev.test/v1/2025/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Use token for protected routes
curl -X GET http://crm.shopdev.test/v1/2025/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📊 API Endpoints

### Authentication
- `POST /v1/2025/auth/register` - User registration
- `POST /v1/2025/auth/login` - User login
- `POST /v1/2025/auth/verify` - Email verification
- `POST /v1/2025/auth/logout` - User logout

### User Management
- `GET /v1/2025/users/profile` - Get user profile
- `PUT /v1/2025/users/profile` - Update user profile
- `POST /v1/2025/users/change-password` - Change password

### Health & Monitoring
- `GET /health` - Application health check
- `GET /swagger/*` - API documentation

For complete API documentation, visit the Swagger UI or check [API Documentation](docs/api-documentation.md).

## 🐳 Docker Services

The application includes the following Docker services:

| Service | Port | Description |
|---------|------|-------------|
| **backend** | 8001 | Go application server |
| **nginx** | 80/443 | Reverse proxy and load balancer |
| **mysql** | 3306 | Primary database |
| **redis** | 6379 | Caching and session storage |
| **kafka** | 9092 | Message broker |
| **kafka-ui** | 8080 | Kafka management interface |

## 📁 Project Structure

```
go-ecommerce-backend-api/
├── cmd/                    # Application entry points
│   ├── server/            # Main server application
│   ├── cli/               # CLI tools and utilities
│   └── swag/docs/         # Generated Swagger documentation
├── configs/               # Configuration files
│   ├── local.yaml         # Local development config
│   ├── docker.yaml        # Docker environment config
│   └── production.yaml    # Production config
├── docs/                  # Project documentation
├── internal/              # Private application code
│   ├── controller/        # HTTP request handlers
│   ├── service/           # Business logic layer
│   ├── repo/              # Data access layer
│   ├── model/             # Data models and structures
│   ├── middleware/        # HTTP middleware
│   └── initialize/        # Application initialization
├── nginx/                 # Nginx configuration
│   ├── conf.d/           # Server configurations
│   └── nginx.conf        # Main nginx config
├── pkg/                   # Public packages
│   ├── logger/           # Logging utilities
│   ├── response/         # HTTP response helpers
│   └── utils/            # Common utilities
├── scripts/               # Automation scripts
├── sql/                   # Database files
│   ├── queries/          # SQL queries for SQLC
│   └── schema/           # Database migrations
├── storage/               # Storage directories
│   ├── logs/             # Application logs
│   └── backups/          # Database backups
├── tests/                 # Test files
├── docker-compose.yml     # Docker composition
├── Dockerfile            # Container definition
├── Makefile              # Build automation
└── README.md             # This file
```

## 🔧 Development

### Adding New Features

1. **Create Models**: Define data structures in `internal/model/`
2. **Database**: Add migrations in `sql/schema/` and queries in `sql/queries/`
3. **Repository**: Implement data access in `internal/repo/`
4. **Service**: Add business logic in `internal/service/`
5. **Controller**: Create HTTP handlers in `internal/controller/`
6. **Routes**: Register routes in `internal/routers/`
7. **Documentation**: Update Swagger annotations

### Testing

```bash
# Run all tests
go test ./...

# Run specific tests
go test ./tests/api/
go test ./internal/service/

# Run with coverage
go test -cover ./...
```

### Database Management

```bash
# Create new migration
make goose-create name=create_products_table

# Run migrations
make goose-up

# Rollback migration
make goose-down

# Generate Go code from SQL
make sqlgen
```

## 🚀 Deployment

### Development Deployment
Follow the Quick Start guide above.

### Production Deployment
See the comprehensive [Deployment Guide](docs/deployment-guide.md) for production setup, including:
- Server preparation
- SSL configuration
- Database setup
- Monitoring and backup strategies

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Commit changes**: `git commit -m 'Add amazing feature'`
4. **Push to branch**: `git push origin feature/amazing-feature`
5. **Open a Pull Request**

### Development Guidelines
- Follow Go conventions and best practices
- Write tests for new features
- Update documentation
- Ensure all tests pass
- Add Swagger annotations for API endpoints

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: Check the `docs/` directory
- **Issues**: Open an issue on GitHub
- **Discussions**: Use GitHub Discussions for questions
- **Email**: contact@yourcompany.com

## 🎯 Roadmap

- [ ] **Product Management**: Product catalog and inventory
- [ ] **Order Processing**: Shopping cart and order management
- [ ] **Payment Integration**: Multiple payment gateways
- [ ] **User Roles**: Admin, vendor, customer roles
- [ ] **Analytics**: Business intelligence and reporting
- [ ] **Mobile API**: Enhanced mobile app support
- [ ] **Microservices**: Split into domain-specific services
- [ ] **GraphQL**: GraphQL API support
- [ ] **Real-time**: WebSocket integration
- [ ] **Kubernetes**: K8s deployment configs

## 📈 Status

- **Build Status**: [![Build Status](https://github.com/your-username/go-ecommerce-backend-api/workflows/CI/badge.svg)](https://github.com/your-username/go-ecommerce-backend-api/actions)
- **Coverage**: [![Coverage](https://codecov.io/gh/your-username/go-ecommerce-backend-api/branch/main/graph/badge.svg)](https://codecov.io/gh/your-username/go-ecommerce-backend-api)
- **Go Report**: [![Go Report Card](https://goreportcard.com/badge/github.com/your-username/go-ecommerce-backend-api)](https://goreportcard.com/report/github.com/your-username/go-ecommerce-backend-api)
- **Version**: ![Version](https://img.shields.io/github/v/release/your-username/go-ecommerce-backend-api)

---

**Built with ❤️ using Go and modern technologies**
