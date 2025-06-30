-include .env
export

.PHONY: startup stop down logs init restart migrate rollback migrate-docker rollback-docker help

startup:
	docker-compose -f docker-compose.yml up --build -d

stop:
	docker-compose -f docker-compose.yml down

down:
	docker-compose -f docker-compose.yml down -v

logs:
	@if [ -z "$(SERVICE)" ]; then \
	  docker-compose -f docker-compose.yml logs -f; \
	else \
	  docker-compose -f docker-compose.yml logs -f $(SERVICE); \
	fi

migrate:
	migrate -path ./migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable" up

rollback:
	migrate -path ./migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable" down 1

migrate-docker:
	docker-compose run --env-file .env migrate

rollback-docker:
	docker-compose run --env-file .env migrate down 1

help:
	@echo "Available commands:"
	@echo "  startup         Build and start containers in detached mode"
	@echo "  stop            Stop containers"
	@echo "  down            Stop and remove containers, networks, volumes"
	@echo "  logs            Show logs (set SERVICE=<name> to filter)"
	@echo "  init            Run the Go application"
	@echo "  restart         Restart containers"
	@echo "  migrate         Run database migrations locally"
	@echo "  rollback        Rollback last database migration locally"
	@echo "  migrate-docker  Run database migrations using Docker container"
	@echo "  rollback-docker Rollback last migration using Docker container"


test:
	go test ./...

init:
	go run cmd/main.go

restart: down startup
