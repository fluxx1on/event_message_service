golangci-lint:
	golangci-lint run -v ./internal/...
.PHONY: lint

create-db:
	psql -h localhost -d postgres -c "CREATE DATABASE notify;"
.PHONY: create-db

migrate-db:
	go build -o ./bin/migrate cmd/migrate/main.go
	DB_PATH=./config/db/postgres.yaml ./bin/migrate -m
.PHONY: migrate-db

delete-migration:
	go build -o ./bin/migrate cmd/migrate/main.go
	DB_PATH=./config/db/postgres.yaml ./bin/migrate

build:
	go mod tidy && go mod download
	go build -o ./bin/server ./cmd/server
.PHONY: build

swag-v1:
	swag init -g internal/transport/http/v1/router.go
.PHONY: swag-v1

run: swag-v1 build
	CONFIG_PATH=./config/develop.yaml DB_PATH=./config/db/postgres.yaml ./bin/server

setup: create-db migrate-db

test:
	go test -v -race -covermode atomic -coverprofile=coverage.txt ./internal/...
.PHONY: test

docker-up:
	docker compose up --build -d postgres nats app nats-streaming && docker compose logs -f

docker-integration:
	docker compose up --build --abort-on-container-exit --exit-code-from integration
.PHONY: integration

docker-down:
	docker compose down --remove-orphans

test-all: lint test integration
