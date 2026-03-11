# Makefile for SGH Assessment Project
# Provides a clean, cross-platform interface for common tasks

# Detect the operating system
ifeq ($(OS),Windows_NT)
	DETECTED_OS := Windows
	SCRIPT_EXT := .bat
	SCRIPT_PREFIX := scripts\\
else
	DETECTED_OS := $(shell uname -s)
	SCRIPT_EXT := .sh
	SCRIPT_PREFIX := ./scripts/
endif

# Detect docker compose command (docker-compose vs docker compose)
ifeq ($(DETECTED_OS),Windows)
	# On Windows, use docker compose (modern syntax)
	DOCKER_COMPOSE := docker compose
else
	# On Unix-like systems, detect which command is available
	DOCKER_COMPOSE := $(shell if command -v docker-compose >/dev/null 2>&1; then echo "docker-compose"; elif command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then echo "docker compose"; else echo "docker compose"; fi)
endif

# Color output for better readability (works on Unix-like systems)
ifneq ($(DETECTED_OS),Windows)
	GREEN := \033[0;32m
	YELLOW := \033[0;33m
	NC := \033[0m # No Color
else
	GREEN :=
	YELLOW :=
	NC :=
endif

.PHONY: help run start stop restart test clean status logs

# Default target - show help
help:
	@echo "$(GREEN)SGH Assessment - Available Commands:$(NC)"
	@echo ""
	@echo "  $(YELLOW)make run$(NC)      - Build and start the Docker containers"
	@echo "  $(YELLOW)make start$(NC)    - Alias for 'make run'"
	@echo "  $(YELLOW)make stop$(NC)     - Stop and remove the containers"
	@echo "  $(YELLOW)make restart$(NC)  - Restart the entire application"
	@echo "  $(YELLOW)make test$(NC)     - Run backend unit tests"
	@echo "  $(YELLOW)make logs$(NC)     - View container logs"
	@echo "  $(YELLOW)make status$(NC)   - Show running containers"
	@echo "  $(YELLOW)make clean$(NC)    - Stop containers and remove volumes"
	@echo ""
	@echo "Detected OS: $(DETECTED_OS)"

# Start the application
run:
	@echo "$(GREEN)Starting the application...$(NC)"
	@$(SCRIPT_PREFIX)run$(SCRIPT_EXT)

# Alias for run
start: run

# Stop the application
stop:
	@echo "$(YELLOW)Stopping the application...$(NC)"
	@$(SCRIPT_PREFIX)stop$(SCRIPT_EXT)

# Restart the application
restart:
	@echo "$(YELLOW)Restarting the application...$(NC)"
	@$(SCRIPT_PREFIX)restart$(SCRIPT_EXT)

# Run backend tests
test:
	@echo "$(GREEN)Running backend tests...$(NC)"
	@$(SCRIPT_PREFIX)test$(SCRIPT_EXT)

# View logs from all containers
logs:
	@echo "$(GREEN)Viewing container logs...$(NC)"
	@$(DOCKER_COMPOSE) logs -f

# Show status of running containers
status:
	@echo "$(GREEN)Container Status:$(NC)"
	@$(DOCKER_COMPOSE) ps

# Clean up everything including volumes
clean:
	@echo "$(YELLOW)Cleaning up containers and volumes...$(NC)"
	@$(DOCKER_COMPOSE) down -v
	@echo "$(GREEN)Cleanup complete!$(NC)"
