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
GOOSE_DBSTRING = $(DB_USER):$(DB_PASSWORD)@tcp(127.0.0.1:3306)/$(DB_NAME)
BACKUP_DIR = storate/backups
GOOSE_MIGRATION_DIR ?= sql/schema

# Get MySQL container name from docker-compose
DB_CONTAINER = go-ecommerce-mysql



# Default target: build app
build:
	go build -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/

# Default target: run app
dev: swag
	@echo "[INFO] Starting development server..."
	@go run ./cmd/$(APP_NAME)/


# Run app (build + docker-up + run)
run: docker-up build
	@echo "[INFO] Running the application..."
	@./bin/$(APP_NAME)

	

# Start Docker containers
docker-up:
	@echo "[INFO] Starting Docker containers..."
	@docker-compose -f $(DOCKER_COMPOSE) up -d
	@echo "[SUCCESS] Docker containers started."

# Stop Docker containers
docker-down:
	@echo "[STOP] Stopping Docker containers..."
	@docker-compose -f $(DOCKER_COMPOSE) down
	@echo "[SUCCESS] Docker containers stopped."

# Rebuild Docker images
docker-build:
	@echo "[INFO] Rebuilding Docker images..."
	@docker-compose -f $(DOCKER_COMPOSE) up -d --build
	@echo "[SUCCESS] Docker images rebuilt."

# View logs
logs:
	@docker-compose -f $(DOCKER_COMPOSE) logs -f

# Clean volumes (DANGER)
clean:
	@echo "[WARNING] Removing Docker volumes..."
	@docker-compose -f $(DOCKER_COMPOSE) down -v
	@echo "[SUCCESS] Docker volumes removed."

# Dump schema only (no data)
dump-schema:
	@echo "[INFO] Dumping schema only..."
	@if not exist "$(BACKUP_DIR)" mkdir $(BACKUP_DIR)
	@docker exec -i $(DB_CONTAINER) mysqldump \
	-u$(DB_USER) -p$(DB_PASSWORD) \
	--databases $(DB_NAME) \
	--add-drop-database \
	--add-drop-table \
	--add-drop-trigger \
	--add-locks \
	--no-data > $(BACKUP_DIR)/$(DB_NAME)_schema.sql
	@echo "[SUCCESS] Schema dumped to $(BACKUP_DIR)/$(DB_NAME)_schema.sql"

# Dump full database (schema + data)
dump-full:
	@echo "[INFO] Dumping full database..."
	@if not exist "$(BACKUP_DIR)" mkdir $(BACKUP_DIR)
	@docker exec -i $(DB_CONTAINER) mysqldump \
	-u$(DB_USER) -p$(DB_PASSWORD) \
	--databases $(DB_NAME) \
	--add-drop-database \
	--add-drop-table \
	--add-drop-trigger \
	--add-locks > $(BACKUP_DIR)/$(DB_NAME)_full.sql
	@echo "[SUCCESS] Full dump saved to $(BACKUP_DIR)/$(DB_NAME)_full.sql"
# ================== Goose DB Migration ==================
goose-create:
	@echo [INFO] Creating new migration file...
	@set "GOOSE_DRIVER=$(GOOSE_DRIVER)" && \
	 set "GOOSE_DBSTRING=$(GOOSE_DBSTRING)" && \
	 goose -dir=$(GOOSE_MIGRATION_DIR) create $(name) sql
	@echo [SUCCESS] Migration file created in $(GOOSE_MIGRATION_DIR)


goose-up:
	@set GOOSE_DRIVER=$(GOOSE_DRIVER)&& set GOOSE_DBSTRING=$(GOOSE_DBSTRING)&& goose -dir=$(GOOSE_MIGRATION_DIR) up

goose-down:
	@set GOOSE_DRIVER=$(GOOSE_DRIVER)&& set GOOSE_DBSTRING=$(GOOSE_DBSTRING)&& goose -dir=$(GOOSE_MIGRATION_DIR) down

goose-reset:
	@set GOOSE_DRIVER=$(GOOSE_DRIVER)&& set GOOSE_DBSTRING=$(GOOSE_DBSTRING)&& goose -dir=$(GOOSE_MIGRATION_DIR) reset
goose-up-by-one:
	@set GOOSE_DRIVER=$(GOOSE_DRIVER)&& set GOOSE_DBSTRING=$(GOOSE_DBSTRING)&& goose -dir=$(GOOSE_MIGRATION_DIR) up-by-one
sqlgen:
	@echo "[INFO] Generating SQL files..."
	@sqlc generate
	@echo "[SUCCESS] SQL files generated."
swag:
	swag init -g ./internal/initialize/swag.init.go -o ./cmd/swag/docs

.PHONY: run docker-up docker-down docker-build clean logs \
        dump-schema dump-full goose-up goose-down goose-reset \
		sqlgen dev goose-create goose-up-by-one swag