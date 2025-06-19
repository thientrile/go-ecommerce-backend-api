# -------------------------
# Project Configuration
# -------------------------
APP_NAME            = server
IMAGE_NAME          = go-ecommerce-backend-api
TAG                 = latest
DOCKER_COMPOSE      = docker-compose.yml

# -------------------------
# Database Configuration
# -------------------------
DB_NAME             = shopDevgo
DB_USER             = root
DB_PASSWORD         = strongpassword123!
DB_CONTAINER        = go-ecommerce-mysql
BACKUP_DIR          = storage/backups

GOOSE_DRIVER        ?= mysql
GOOSE_DBSTRING      = $(DB_USER):$(DB_PASSWORD)@tcp(127.0.0.1:3306)/$(DB_NAME)
GOOSE_MIGRATION_DIR ?= sql/schema

# Detect OS
ifeq ($(OS),Windows_NT)
	DETECTED_OS := Windows
	MKDIR_P = if not exist $(subst /,\,$(BACKUP_DIR)) mkdir $(subst /,\,$(BACKUP_DIR))
	SET_ENV = set
	SHELL := powershell.exe
else
	DETECTED_OS := $(shell uname -s)
	MKDIR_P = mkdir -p $(BACKUP_DIR)
	SET_ENV = export
	SHELL := /bin/bash
endif

# -------------------------
# Initialization
# -------------------------
# Display help information
help:
	@echo "========================================="
	@echo "  Go E-commerce Backend API - Makefile"
	@echo "========================================="
	@echo ""
	@echo "DEVELOPMENT COMMANDS:"
	@echo "  init          - Initialize project dependencies"
	@echo "  build         - Build the application"
	@echo "  dev           - Start development server"
	@echo "  dev-services  - Start development services only"
	@echo "  dev-down      - Stop development services"
	@echo "  watch         - Start server with hot reload (Air)"
	@echo "  run-local     - Run with ENV=local"
	@echo "  test-local    - Test local environment"
	@echo ""	@echo "PRODUCTION COMMANDS:"
	@echo "  start         - Start production services"
	@echo "  stop          - Stop production services"
	@echo "  restart       - Restart production services"
	@echo "  status        - Check service status"
	@echo "  prod-up       - Start production environment"
	@echo "  prod-down     - Stop production environment"
	@echo "  prod-build    - Build and start production"
	@echo "  stop-all      - Stop all services (dev + prod)"
	@echo "  switch-to-dev - Switch from prod to dev mode"
	@echo "  switch-to-prod- Switch from dev to prod mode"
	@echo ""
	@echo "DOCKER COMMANDS:"
	@echo "  docker-up     - Start Docker containers"
	@echo "  docker-down   - Stop Docker containers"
	@echo "  docker-build  - Rebuild Docker images"
	@echo "  clean         - Remove Docker volumes"
	@echo "  logs          - Show Docker logs"
	@echo ""
	@echo "DATABASE COMMANDS:"
	@echo "  dump-schema   - Backup database schema"
	@echo "  dump-full     - Full database backup"
	@echo "  goose-create  - Create new migration (use: make goose-create name=migration_name)"
	@echo "  goose-up      - Run migrations"
	@echo "  goose-down    - Rollback last migration"
	@echo "  goose-reset   - Reset all migrations"
	@echo ""	@echo "DEVELOPMENT TOOLS:"
	@echo "  sqlgen        - Generate SQL code with sqlc"
	@echo "  swag          - Generate Swagger documentation"
	@echo ""
	@echo "DEBEZIUM COMMANDS:"
	@echo "  debezium-register     - Register MySQL connector with Debezium"
	@echo "  debezium-status       - Check Debezium connector status"
	@echo "  debezium-delete       - Delete Debezium connector"
	@echo "  debezium-restart      - Restart Debezium connector"
	@echo "  debezium-connectors   - List all Debezium connectors"
	@echo ""
	@echo "EXAMPLES:"
	@echo "  make start                    # Start production services"
	@echo "  make dev                      # Start development"
	@echo "  make goose-create name=users  # Create users migration"
	@echo ""

init:
	@echo "[INFO] Initializing project..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/air-verse/air@latest
	go mod download
	@echo "[SUCCESS] Initialization complete."

# -------------------------
# Build & Dev
# -------------------------
build:
	@echo "[INFO] Building application..."
	go build -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/

dev: swag 
	@echo "[INFO] Starting development server with ENV=local..."
	go run ./cmd/$(APP_NAME)/

dev-services:
	@echo "[INFO] Stopping production services first..."
ifeq ($(OS),Windows_NT)
	@powershell -Command "try { docker-compose --profile prod down } catch { }"
else
	@docker-compose --profile prod down 2>/dev/null || true
endif
	@echo "[INFO] Starting development services (MySQL, Redis, Kafka, Nginx-Dev)..."
	docker-compose --profile dev up -d
	@echo "[SUCCESS] Development services started."

dev-down:
	@echo "[INFO] Stopping development services..."
	docker-compose --profile dev down
	@echo "[SUCCESS] Development services stopped."

# Development with hot reload using Air
watch: dev-services 
	@echo "[INFO] Starting development server with hot reload using Air..."
	air

# Local development with ENV=local
run-local: swag
	@echo "[INFO] Starting local development server with ENV=local..."
	@$(SET_ENV) ENV=local && go run ./cmd/$(APP_NAME)/

# Test the local environment
test-local:
	@echo "[INFO] Testing local environment..."
	@$(SET_ENV) ENV=local && go run ./cmd/$(APP_NAME)/ &
	@powershell -Command "Start-Sleep 3; try { Invoke-WebRequest -Uri 'http://crm.shopdev.test/health' -UseBasicParsing -TimeoutSec 5; Write-Host '✅ Local environment working' } catch { Write-Host '❌ Local environment failed' }"

run: docker-up build
	@echo "[INFO] Running application..."
	./bin/$(APP_NAME)

# -------------------------
# Docker Controls
# -------------------------
docker-up:
	@echo "[INFO] Starting Docker containers..."
	docker-compose up -d
	@echo "[SUCCESS] Docker containers started."

docker-down:
	@echo "[INFO] Stopping Docker containers..."
	docker-compose down
	@echo "[SUCCESS] Docker containers stopped."

docker-build:
	@echo "[INFO] Rebuilding Docker images..."
	docker-compose up -d --build
	@echo "[SUCCESS] Docker images rebuilt."

# Production mode - with Go backend container
prod-up:
	@echo "[INFO] Stopping development services first..."
ifeq ($(OS),Windows_NT)
	@powershell -Command "try { docker-compose --profile dev down } catch { }"
else
	@docker-compose --profile dev down 2>/dev/null || true
endif
	@echo "[INFO] Starting production environment..."
	docker-compose --profile prod up -d
	@echo "[SUCCESS] Production environment started."

prod-down:
	@echo "[INFO] Stopping production environment..."
	docker-compose --profile prod down
	@echo "[SUCCESS] Production environment stopped."

prod-build:
	@echo "[INFO] Stopping development services first..."
ifeq ($(OS),Windows_NT)
	@powershell -Command "try { docker-compose --profile dev down } catch { }"
else
	@docker-compose --profile dev down 2>/dev/null || true
endif
	@echo "[INFO] Building and starting production environment..."
	docker-compose --profile prod up -d --build
	@echo "[SUCCESS] Production environment built and started."

# Start production services
start: prod-up
	@echo "[SUCCESS] Production services started successfully!"	@echo "[INFO] Application endpoints:"
	@echo "  - API: http://crm.shopdev.com"
	@echo "  - Admin: http://admin.shopdev.com"
	@echo "  - Swagger: http://crm.shopdev.com/swagger/index.html"
	@echo "[INFO] Checking service status in 5 seconds..."
ifeq ($(OS),Windows_NT)
	@powershell -Command "Start-Sleep 5"
	@powershell -Command "docker ps --filter 'label=com.docker.compose.project=go-ecommerce-backend-api'"
else
	@sleep 5
	@$(MAKE) status
endif

# Stop production services
stop: prod-down
	@echo "[SUCCESS] Production services stopped successfully!"

# Stop all services (dev + prod)
stop-all:
	@echo "[INFO] Stopping all services (development and production)..."
ifeq ($(OS),Windows_NT)
	@powershell -Command "try { docker-compose --profile dev down } catch { }"
	@powershell -Command "try { docker-compose --profile prod down } catch { }"
else
	@docker-compose --profile dev down 2>/dev/null || true
	@docker-compose --profile prod down 2>/dev/null || true
endif
	@echo "[SUCCESS] All services stopped!"

# Switch to development mode
switch-to-dev: 
	@echo "[INFO] Switching to development mode..."
ifeq ($(OS),Windows_NT)
	@powershell -Command "try { docker-compose --profile prod down } catch { }"
else
	@docker-compose --profile prod down 2>/dev/null || true
endif
	@$(MAKE) dev-services
	@echo "[SUCCESS] Switched to development mode!"

# Switch to production mode  
switch-to-prod:
	@echo "[INFO] Switching to production mode..."
ifeq ($(OS),Windows_NT)
	@powershell -Command "try { docker-compose --profile dev down } catch { }"
else
	@docker-compose --profile dev down 2>/dev/null || true
endif
	@$(MAKE) prod-up
	@echo "[SUCCESS] Switched to production mode!"

# Restart production services
restart: stop start
	@echo "[SUCCESS] Production services restarted successfully!"

# Check service status
status:
	@echo "[INFO] Checking service status..."
	@echo "========================================="
	@echo "DEVELOPMENT SERVICES:"
	@docker-compose --profile dev ps 2>/dev/null || echo "No development services running"
	@echo ""
	@echo "PRODUCTION SERVICES:"
	@docker-compose --profile prod ps 2>/dev/null || echo "No production services running"
	@echo "========================================="
	@echo "[INFO] Service health checks:"
	@powershell -Command "try { Invoke-WebRequest -Uri 'http://crm.shopdev.test/health' -UseBasicParsing -TimeoutSec 5; Write-Host '✅ API service is healthy' } catch { Write-Host '❌ API service is not responding' }"

clean:
	@echo "[WARNING] Removing Docker volumes..."
	docker-compose -f $(DOCKER_COMPOSE) down -v
	@echo "[SUCCESS] Volumes removed."

logs:
	docker-compose -f $(DOCKER_COMPOSE) logs -f

# -------------------------
# Nginx Controls
# -------------------------


# -------------------------
# Database Backups
# -------------------------
dump-schema:
	@echo "[INFO] Dumping database schema..."
	@$(MKDIR_P)
	docker exec -i $(DB_CONTAINER) mysqldump \
		-u$(DB_USER) -p$(DB_PASSWORD) \
		--databases $(DB_NAME) \
		--add-drop-database --add-drop-table --add-drop-trigger \
		--add-locks --no-data > $(BACKUP_DIR)/$(DB_NAME)_schema.sql
	@echo "[SUCCESS] Schema saved to $(BACKUP_DIR)/$(DB_NAME)_schema.sql"

dump-full:
	@echo "[INFO] Dumping full database..."
	@$(MKDIR_P)
	docker exec -i $(DB_CONTAINER) mysqldump \
		-u$(DB_USER) -p$(DB_PASSWORD) \
		--databases $(DB_NAME) \
		--add-drop-database --add-drop-table --add-drop-trigger \
		--add-locks > $(BACKUP_DIR)/$(DB_NAME)_full.sql
	@echo "[SUCCESS] Full dump saved to $(BACKUP_DIR)/$(DB_NAME)_full.sql"

# -------------------------
# Goose Migrations
# -------------------------
goose-create:
	@echo "[INFO] Creating new migration: $(name)"
	$(SET_ENV) GOOSE_DRIVER=$(GOOSE_DRIVER) && \
	$(SET_ENV) GOOSE_DBSTRING=$(GOOSE_DBSTRING) && \
	goose -dir=$(GOOSE_MIGRATION_DIR) create $(name) sql
	@echo "[SUCCESS] Migration created."

goose-up:
	$(SET_ENV) GOOSE_DRIVER=$(GOOSE_DRIVER) && \
	$(SET_ENV) GOOSE_DBSTRING=$(GOOSE_DBSTRING) && \
	goose -dir=$(GOOSE_MIGRATION_DIR) up

goose-down:
	$(SET_ENV) GOOSE_DRIVER=$(GOOSE_DRIVER) && \
	$(SET_ENV) GOOSE_DBSTRING=$(GOOSE_DBSTRING) && \
	goose -dir=$(GOOSE_MIGRATION_DIR) down

goose-reset:
	$(SET_ENV) GOOSE_DRIVER=$(GOOSE_DRIVER) && \
	$(SET_ENV) GOOSE_DBSTRING=$(GOOSE_DBSTRING) && \
	goose -dir=$(GOOSE_MIGRATION_DIR) reset

goose-up-by-one:
	$(SET_ENV) GOOSE_DRIVER=$(GOOSE_DRIVER) && \
	$(SET_ENV) GOOSE_DBSTRING=$(GOOSE_DBSTRING) && \
	goose -dir=$(GOOSE_MIGRATION_DIR) up-by-one

# -------------------------
# SQL & Swagger
# -------------------------
sqlgen:
	@echo "[INFO] Generating SQL code..."
	sqlc generate
	@echo "[SUCCESS] SQL code generated."

swag:
	@echo "[INFO] Generating Swagger docs..."
	swag init -g ./internal/initialize/swag.init.go -o ./cmd/swag/docs
	@echo "[SUCCESS] Swagger docs generated."



# -------------------------
# Debezium Management
# -------------------------

# Register MySQL connector with Debezium
debezium-register:
	@echo "[INFO] Registering MySQL connector with Debezium..."
	@echo "[INFO] Waiting for Debezium to be ready..."
ifeq ($(OS),Windows_NT)
	@powershell -Command "Start-Sleep 10"
	@curl.exe -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" http://localhost:8083/connectors/ -d @register-mysql.json
else
	@sleep 10
	@curl -i -X POST -H "Accept:application/json" -H "Content-Type:application/json" http://localhost:8083/connectors/ -d @register-mysql.json
endif
	@echo ""
	@echo "[SUCCESS] MySQL connector registered successfully!"

# Check Debezium connector status
debezium-status:
	@echo "[INFO] Checking Debezium connector status..."
ifeq ($(OS),Windows_NT)
	@curl.exe -s http://localhost:8083/connectors/mysql-connector/status
else
	@curl -s http://localhost:8083/connectors/mysql-connector/status | jq '.'
endif

# List all Debezium connectors
debezium-connectors:
	@echo "[INFO] Listing all Debezium connectors..."
ifeq ($(OS),Windows_NT)
	@curl.exe -s http://localhost:8083/connectors
else
	@curl -s http://localhost:8083/connectors | jq '.'
endif

# Delete Debezium connector
debezium-delete:
	@echo "[INFO] Deleting MySQL connector..."
ifeq ($(OS),Windows_NT)
	@curl.exe -i -X DELETE http://localhost:8083/connectors/mysql-connector
else
	@curl -i -X DELETE http://localhost:8083/connectors/mysql-connector
endif
	@echo ""
	@echo "[SUCCESS] MySQL connector deleted!"

# Restart Debezium connector
debezium-restart:
	@echo "[INFO] Restarting MySQL connector..."
ifeq ($(OS),Windows_NT)
	@curl.exe -i -X POST http://localhost:8083/connectors/mysql-connector/restart
else
	@curl -i -X POST http://localhost:8083/connectors/mysql-connector/restart
endif
	@echo ""
	@echo "[SUCCESS] MySQL connector restarted!"

# Check Debezium service health
debezium-health:
	@echo "[INFO] Checking Debezium service health..."
ifeq ($(OS),Windows_NT)
	@curl.exe -s http://localhost:8083/
else
	@curl -s http://localhost:8083/
endif

# Auto setup: Start services and register connector
debezium-setup: dev-services
	@echo "[INFO] Setting up Debezium with MySQL connector..."
	@echo "[INFO] Waiting for services to be ready..."
ifeq ($(OS),Windows_NT)
	@powershell -Command "Start-Sleep 30"
else
	@sleep 30
endif
	@$(MAKE) debezium-register
	@echo "[SUCCESS] Debezium setup completed!"

# -------------------------
# Phony Targets
# -------------------------
.PHONY: help init build dev watch run \
	docker-up docker-down docker-build clean logs \
	dump-schema dump-full \
	goose-create goose-up goose-down goose-reset goose-up-by-one \
	sqlgen swag start stop restart status stop-all switch-to-dev switch-to-prod \
	debezium-register debezium-status debezium-delete debezium-restart debezium-connectors debezium-health debezium-setup\
	debezium-register debezium-status debezium-delete debezium-restart debezium-connectors debezium-health debezium-setup
