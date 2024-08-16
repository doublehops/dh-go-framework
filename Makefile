.SILENT:

run:
	go run cmd/server/run.go -config ./config.json

run_for_test:
	go run cmd/server/run.go -config ./config_test.json

gofmt:
	gofumpt -l -w .

# `test` is a reserved keyword in makefiles.
tst:
	go test ./... -cover

# To fix permissions issue, run: sudo chown -R $(whoami) .db-data/
lint:
	golangci-lint --config ./ci/.golangci-lint.yml run

test:
	go test ./... -cover

SHELL := /bin/bash
docker_up:
	source .env && docker compose -f docker-compose.yml up -d

docker_down:
	docker compose -f docker-compose.yml down

# make scaffold model=<table_name>
scaffold:
	go run ./cmd/scaffold/run.go -config ./config.json -table $(table)

migrate:
	go run ./cmd/migrate/migrate.go -action up

# Run migrations for test database.
migrate_test:
	go run ./cmd/migrate/migrate.go -action up -config config_test.json
