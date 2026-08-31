# `migrate` and `atlas` install into GOPATH/bin, which is not on the default
# macOS PATH. Resolve them here so the targets work from a bare shell.
export PATH := $(PATH):$(shell go env GOPATH)/bin

DB_URL ?= postgresql://postgres:password@localhost:5432/ecommerce_shop?sslmode=disable
MIGRATIONS_DIR := db/migrations

.PHONY: help build run dev lint format db-diff db-hash db-inspect migrate-up migrate-down migrate-reset migrate-status docker-up docker-down

help:
	@echo "Available Commands:"
	@echo "make build              - Build the application"
	@echo "make run                - Run the application"
	@echo "make dev                - Run the application in development mode"
	@echo "make lint               - Run linter on the codebase"
	@echo "make format             - Format the code and re-arrange imports"
	@echo ""
	@echo "make db-diff name=xxx   - Generate a migration from the GORM models"
	@echo "make db-inspect         - Print the DDL the models currently describe"
	@echo "make db-hash            - Re-hash atlas.sum after hand-editing a migration"
	@echo "make migrate-up         - Apply pending migrations"
	@echo "make migrate-down       - Roll back the last migration"
	@echo "make migrate-reset      - Roll back every migration"
	@echo "make migrate-status     - Show the currently applied version"

build:
	@go build -o bin/app ./cmd/api

run:
	@go run ./cmd/api

dev:
	@go run ./cmd/api

lint:format
	@golangci-lint run ./...

format:
	@gofmt -s -w .

# Diff internal/models against db/migrations and write the SQL for whatever
# changed. Requires Docker: Atlas spins up a throwaway Postgres to normalise
# the schema, then throws it away. Your real database is never touched.
db-diff:
ifndef name
	$(error usage: make db-diff name=add_wishlist)
endif
	@atlas migrate diff $(name) --env gorm

db-inspect:
	@go run -C db/loader .

db-hash:
	@atlas migrate hash --env gorm

migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-reset:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

migrate-status:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

docker-up:
	@docker compose -f docker/docker-compose.yml up -d

docker-down:
	@docker compose -f docker/docker-compose.yml down
