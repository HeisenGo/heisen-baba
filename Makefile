DOCKER_COMPOSE := docker-compose.yaml
BIN_DIR := ./bin

# Run Docker Compose
.PHONY: dockerize
dockerize:
	@echo "Running Docker Compose..."
	rm -rf $(BIN_DIR)
	docker compose down --remove-orphans
	docker compose -f $(DOCKER_COMPOSE) up --build -d