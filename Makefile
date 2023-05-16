#
# Makefile
#

run-backend:
	go run backend/main.go
.PHONY: run-backend

run-frontend:
	cd web && yarn dev --port 4000
.PHONY: run-frontend

local-up:
	podman-compose up --force-recreate --remove-orphans -d
.PHONY: local-up
