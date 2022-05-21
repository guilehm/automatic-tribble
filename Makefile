DOCKER_COMPOSE=docker-compose
APP_NAME?=tribble
DATABASE_TESTS_URL=postgres://postgres:postgres@db:5432/tests?sslmode=disable

build:
	@echo "Building the app"
	-$(DOCKER_COMPOSE) build tribble

run:
	-$(DOCKER_COMPOSE) up

up:
	-$(DOCKER_COMPOSE) up -d

up-db:
	-$(DOCKER_COMPOSE) up -d postgres 

stop:
	@echo "Stopping containers"
	-$(DOCKER_COMPOSE) stop

down:
	@echo "Removing containers"
	-$(DOCKER_COMPOSE) down

remove:
	@echo "Removing containers and volumes"
	-$(DOCKER_COMPOSE) down -v

setup: build up-db

test:
	@echo "Running tests"
	-$(DOCKER_COMPOSE) exec -e DATABASE_URL=$(DATABASE_TESTS_URL) $(APP_NAME) go run ./tests/tests.go
	-$(DOCKER_COMPOSE) exec -e DATABASE_URL=$(DATABASE_TESTS_URL) $(APP_NAME) go test -v ./...
