-include .env
export

MIGRATIONS_DIR := database/migrations
DATABASE_URL   ?= postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
MIGRATE        := go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.19.0

.PHONY: test test-integration migrate-up migrate-down run

## test: testes unitários (domínio e services)
test:
	go test ./...

## test-integration: testes de repositório contra Postgres 15 (requer Docker)
test-integration:
	go test -tags integration -count=1 ./...

## migrate-up: aplica todas as migrations pendentes no banco do .env
migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

## migrate-down: desfaz a última migration aplicada
migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

## run: sobe a API com as configurações do .env
run:
	go run ./cmd
