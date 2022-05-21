DOCKER_COMPOSE=docker-compose


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

