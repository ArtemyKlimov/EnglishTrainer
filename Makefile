include .env
.PHONY: build
build:
		go build -v ./cmd/main.go

.PHONY: test
test:
		go test -v -race -timeout 30s ./...

.PHONY: migrate_up
migrate_up:
		migrate -path migrations -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOSTNAME):$(PG_PORT)/$(PG_DBNAME)?sslmode=disable" up

.PHONY: migrate_down
migrate_down:
		migrate -path migrations -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOSTNAME):$(PG_PORT)/$(PG_DBNAME)?sslmode=disable" down

.DEFAULT_GOAL := build

