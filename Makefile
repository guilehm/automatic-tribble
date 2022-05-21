DOCKER_COMPOSE=docker-compose
APP_NAME?=tribble
DATABASE_TESTS_URL=postgres://postgres:postgres@db:5432/tests?sslmode=disable

build:
	@echo "BUILDING THE APP"
	-$(DOCKER_COMPOSE) build $(APP_NAME)

run:
	-$(DOCKER_COMPOSE) up

up-db:
	-$(DOCKER_COMPOSE) up -d postgres

stop:
	@echo "STOPPING CONTAINERS"
	-$(DOCKER_COMPOSE) stop

down:
	@echo "REMOVING CONTAINERS"
	-$(DOCKER_COMPOSE) down

remove:
	@echo "REMOVING CONTAINERS AND VOLUMES"
	-$(DOCKER_COMPOSE) down -v

test:
	@echo "PREPARING DATABASE FOR TESTS\n"
	-$(DOCKER_COMPOSE) exec -e DATABASE_URL=$(DATABASE_TESTS_URL) $(APP_NAME) go run ./tests/prepare-database.go
	@echo "RUNNING TESTS\n"
	-$(DOCKER_COMPOSE) exec -e DATABASE_URL=$(DATABASE_TESTS_URL) $(APP_NAME) go test -v ./...

setup: build up-db
