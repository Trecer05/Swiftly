.PHONY: build build-auth build-chat build-task \
        run run-auth run-chat run-task \
        docker-build docker-up docker-down docker-restart \
        logs logs-auth logs-chat logs-task clean

# Имя бинарников
BIN_DIR := ./bin

# Пути к main.go
AUTH_CMD := ./cmd/servers/auth
CHAT_CMD := ./cmd/servers/chat
TASK_CMD := ./cmd/servers/task_tracker

# ===========================
#        TEST & LINT
# ===========================

go-test:
	go test ./...

lint:
	golangci-lint run ./...

# ===========================
#        BUILD (LOCAL)
# ===========================

build: build-auth build-chat build-task

build-auth:
	@echo "Building auth service..."
	@go build -o $(BIN_DIR)/auth $(AUTH_CMD)

build-chat:
	@echo "Building chat service..."
	@go build -o $(BIN_DIR)/chat $(CHAT_CMD)

build-task:
	@echo "Building task_tracker service..."
	@go build -o $(BIN_DIR)/task_tracker $(TASK_CMD)

# ===========================
#        RUN LOCAL
# ===========================

run: run-auth run-chat run-task

run-auth:
	@echo "Running auth service..."
	@go run $(AUTH_CMD)

run-chat:
	@echo "Running chat service..."
	@go run $(CHAT_CMD)

run-task:
	@echo "Running task_tracker service..."
	@go run $(TASK_CMD)

# ===========================
#        CLEAN
# ===========================

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
	@go clean

# ===========================
#        DOCKER
# ===========================

docker-build:
	@echo "Building all Docker images..."
	@docker compose build

docker-up:
	@echo "Starting all services..."
	@docker compose up -d

docker-down:
	@echo "Stopping all services..."
	@docker compose down

docker-restart:
	@echo "Restarting all services..."
	@docker compose down
	@docker compose up -d

# ===========================
#        LOGS
# ===========================

logs:
	@docker compose logs -f

logs-auth:
	@docker compose logs -f auth_service

logs-chat:
	@docker compose logs -f chat_service

logs-task:
	@docker compose logs -f task_service
