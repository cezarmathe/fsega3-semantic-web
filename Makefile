#
# Makefile
#

DOCKER ?= podman
DOCKER_COMPOSE ?= $(DOCKER)-compose

run-backend:
	go run backend/main.go
.PHONY: run-backend

run-frontend:
	cd web && yarn build && yarn start --port 3000
.PHONY: run-frontend

local-up:
	cp json-server/db.original.json json-server/db.json
	$(DOCKER_COMPOSE) up --force-recreate --remove-orphans -d
.PHONY: local-up
