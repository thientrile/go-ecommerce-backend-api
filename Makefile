# Makefile for go-ecommerce-backend-api

# Project Config
APP_NAME = server
IMAGE_NAME = go-ecommerce-backend-api
TAG = latest

# Docker Compose File
DOCKER_COMPOSE = docker-compose.yml

# Database Config
DB_NAME = shopDevgo
DB_USER = root
DB_PASSWORD = strongpassword123!
GOOSE_DRIVER ?= mysql
GOOSE_DBSTRING = "$(DB_USER):$(DB_PASSWORD)@tcp(127.0.0.1:33306)/$(DB_NAME)"
MIGRATIONS_DIR = migrations
GOOSE_MIGRATION_DIR = $(MIGRATIONS_DIR)/sql

# Get MySQL container name from docker-compose
DB_CONTAINER = go-ecommerce-mysql

.PHONY: run docker-up docker-down docker-build clean logs \
        dump-schema dump-full goose-up goose-down goose-reset

# Default target: build app
build:
	go build -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/

# Default target: run app
dev:
	@echo "🚀 Starting development server..."
	@go run ./cmd/$(APP_NAME)/


run:
	@echo "🚀 Starting the application..."
	docker-up
	go run ./cmd/$(APP_NAME)/

# Start Docker containers
docker-up:
	@echo "🚀 Starting Docker containers..."
	@docker-compose -f $(DOCKER_COMPOSE) up -d
	@echo "✅ Docker containers started."

# Stop Docker containers
docker-down:
	@echo "🛑 Stopping Docker containers..."
	@docker-compose -f $(DOCKER_COMPOSE) down
	@echo "✅ Docker containers stopped."

# Rebuild Docker images
docker-build:
	@echo "🔧 Rebuilding Docker images..."
	@docker-compose -f $(DOCKER_COMPOSE) up --build
	@echo "✅ Docker images rebuilt."

# View logs
logs:
	@docker-compose -f $(DOCKER_COMPOSE) logs -f

# Clean volumes (DANGER)
clean:
	@echo "⚠️ Removing Docker volumes..."
	@docker-compose -f $(DOCKER_COMPOSE) down -v
	@echo "✅ Docker volumes removed."

# Dump schema only (no data)
dump-schema:
	@echo "📦 Dumping schema only..."
	@if not exist "$(MIGRATIONS_DIR)" mkdir $(MIGRATIONS_DIR)
	@docker exec -i $(DB_CONTAINER) mysqldump \
	-u$(DB_USER) -p$(DB_PASSWORD) \
	--databases $(DB_NAME) \
	--add-drop-database \
	--add-drop-table \
	--add-drop-trigger \
	--add-locks \
	--no-data > $(MIGRATIONS_DIR)/$(DB_NAME)_schema.sql
	@echo "✅ Schema dumped to $(MIGRATIONS_DIR)/$(DB_NAME)_schema.sql"

# Dump full database (schema + data)
dump-full:
	@echo "📦 Dumping full database..."
	@if not exist "$(MIGRATIONS_DIR)" mkdir $(MIGRATIONS_DIR)
	@docker exec -i $(DB_CONTAINER) mysqldump \
	-u$(DB_USER) -p$(DB_PASSWORD) \
	--databases $(DB_NAME) \
	--add-drop-database \
	--add-drop-table \
	--add-drop-trigger \
	--add-locks > $(MIGRATIONS_DIR)/$(DB_NAME)_full.sql
	@echo "✅ Full dump saved to $(MIGRATIONS_DIR)/$(DB_NAME)_full.sql"
# ================== Goose DB Migration ==================
goose-up:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) up

goose-down:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) down

goose-reset:
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) reset