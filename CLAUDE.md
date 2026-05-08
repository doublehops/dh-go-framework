# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

**Run the server:**
```bash
go run cmd/server/run.go -config config.json
```

**Run migrations:**
```bash
go run cmd/migrate/migrate.go -action up
go run cmd/migrate/migrate.go -action down -number 1
go run cmd/migrate/migrate.go -action create -name <migration_name>
```

**Run scaffolding** (generates handler/service/repository/model files from a DB table):
```bash
make scaffold table=<table_name>
```

**Run tests:**
```bash
go test ./...
go test ./internal/some/package/... -v -run TestFunctionName
```

**Lint:**
```bash
golangci-lint run --config ci/.golangci-lint.yml
```

## Configuration

The app loads config from a JSON file (default `config.json`, test uses `config_test.json`) and resolves env vars from `.env`. JSON values prefixed with `ENV:` (e.g. `"ENV:MYSQL_HOST"`) are replaced at runtime with the matching env var. See `internal/config/config.go` for the full `Config` struct.

The database for development is MariaDB 10.5 via Docker:
```bash
docker compose up -d
```

## Architecture

This is a layered REST API framework. Each resource follows a strict 4-layer pattern:

```
Handler → Service → Repository → DB
```

- **Handler** (`internal/handlers/<resource>/`) — decodes HTTP request, validates input, calls service, writes JSON response using helpers from `internal/request/responses.go`.
- **Service** (`internal/service/<resource>service/`) — business logic, transaction management. Embeds `*service.App` which holds `DB *sqlx.DB` and `Log *logga.Logga`. Soft-delete and ownership permission checks live here via `BaseModel` and `service.App.HasPermission`.
- **Repository** (`internal/repository/<resource>repository/`) — raw SQL via `sqlx`. Write operations receive a `*sqlx.Tx`; reads receive `*sqlx.DB`. SQL strings are in a companion `sql.go` file.
- **Model** (`internal/model/<resource>/`) — struct embedding `model.BaseModel`, plus a `Validate()` method returning `req.ErrMsgs`.

**Routing:** Routes are registered in `internal/routes/routes.go` using `httprouter-group`. All routes are prefixed `/v1`. Protected routes use `middleware.AuthMiddleware` (currently a stub that injects a hardcoded userID). To add a new resource, register its group in `GetV1Routes()`.

**BaseModel** (`internal/model/basemodel.go`) — all models embed this. It provides `id`, `user_id`, `created_by`, `updated_by`, `created_at`, `updated_at`, `deleted_at`. Soft-deletes are handled by setting `deleted_at`. `SetCreated`/`SetUpdated`/`SetDeleted` pull the authenticated user ID from the request context via `app.UserIDKey`.

**Validation** (`internal/validator/`) — field-level validation rules are declared as `[]validator.Rule` on the model. Each rule references one or more `ValidationFuncs`. `validator.RunValidation` returns a map of field → error messages.

**Request/Response** (`internal/request/`) — `GetRequestParams` parses pagination, filtering, and sorting from query params. Filters are declared per-handler as `[]req.FilterRule`; sortable fields are whitelisted. Response helpers (`GetSingleItemResp`, `GetListResp`, `GetValidateErrResp`, etc.) standardise the JSON envelope.

**Migrations** (`migrations/`) — raw SQL files named `<timestamp>_<name>.up.sql` / `.down.sql` where timestamp is `YYYYMMDDHHmmss` (14 digits, no separators). Multiple statements per file are separated by `;`. Tracking is via a `schema_migrations` table managed by [golang-migrate/migrate](https://github.com/golang-migrate/migrate). If a migration fails mid-run the DB is marked dirty; resolve with `migrate force <version>` before retrying.

**Scaffolding** (`internal/scaffold/`, `cmd/scaffold/`) — reads a live DB table schema and generates the four layers above. Output paths are configured in `cmd/scaffold/config.json`.
