# Tên app
APP_NAME = server

# File docker-compose
DOCKER_COMPOSE = docker-compose.yml

# Mục run: khởi động docker + chạy app
run: docker-up
	go run ./cmd/$(APP_NAME)/

# Khởi động docker compose (nếu chưa chạy)
docker-up:
	@echo "🚀 Starting Docker containers..."
	@docker-compose -f $(DOCKER_COMPOSE) up -d
	@echo "✅ Docker containers started."

# Tắt docker
docker-down:
	@echo "🛑 Stopping Docker containers..."
	@docker-compose -f $(DOCKER_COMPOSE) down
	@echo "✅ Docker containers stopped."

# Xem log
logs:
	@docker-compose logs -f

# Rebuild Docker (nếu có service cần build)
docker-build:
	@docker-compose -f $(DOCKER_COMPOSE) build

# Clean volumes (cẩn thận nhé!)
clean:
	@docker-compose down -v
