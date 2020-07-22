SHELL := /bin/bash

# Running from within docker compose

run: up

up:
	docker-compose -f compose.yaml up --detach --remove-orphans

down:
	docker-compose -f compose.yaml down --remove-orphans

logs:
	docker-compose -f compose.yaml logs -f

# Administration

schema:
	go run app/admin/main.go schema

seed:
	go run app/admin/main.go seed

# Running tests within the local computer

test:
	go test ./... -count=1
	staticcheck ./...