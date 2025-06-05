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
DB_PASSWORD = 1
MIGRATIONS_DIR = migrations

# Get MySQL container name from docker-compose
DB_CONTAINER = go-ecommerce-mysql

.PHONY: build up down clean logs rebuild status



# Default target: run app
run: 
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
