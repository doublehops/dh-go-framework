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


## Container stuff ##
build_container:
	docker rm dhapi-container; docker run -it --name dhapi-container dhapi /bin/bash

conn_container:
	docker run -it --name dhapi-container dhapi /bin/bash

## Build test container
build_test_container:
	docker build -t api-test-image -f Dockerfile_test .

run_test_container:
	docker run --rm --name app-container --network mynetwork -p 8088:8088 api-test-image


