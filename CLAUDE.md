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

Optional host settings: `corsOrigin` (exact frontend origin; empty disables CORS) and `trustForwardedFor` (key the login rate limiter on the last `X-Forwarded-For` entry — only when every request arrives through one reverse proxy that appends to it, such as an AWS ALB). Unset `ENV:` values resolve to an empty string.

Logging: `logging.writer` (`ENV:LOG_WRITER`) is `stdout` (default when empty) or `file`; `file` writes to `logging.filePath` (`ENV:LOG_FILE_PATH`) and to stdout via `io.MultiWriter`.

The database for development is MariaDB 10.5 via Docker:
```bash
docker compose up -d
```

## Architecture

This is a layered REST API framework. Each resource follows a strict 4-layer pattern:

```
Handler → Service → Repository → DB
```

- **Handler** (`internal/handlers/<resource>/`) — decodes HTTP request, calls `Normalise()` then `Validate()` on the model, calls the service/record, writes JSON response using helpers from `internal/request/responses.go`. Handlers never touch repositories.
- **Service** (`internal/service/<resource>service/`) — business logic, transaction management. Embeds `*service.App` which holds `DB *sqlx.DB` and `Log *logga.Logga`. Split into three files:
  - `<resource>service.go` — `New(app)` wires the concrete repository internally; collection-level methods (`GetAll`), `GetByID(ctx, id)` (returns `(*<Resource>Record, error)`, error when not found) and `NewRecord(model)`.
  - `<resource>record.go` — `<Resource>Record{Record *model; svc}` with the single-record operations `Create(ctx)`, `Update(ctx, payload)` (copies only editable fields), `Delete(ctx)`, and record-level permission checks such as `CanBeAccessedBy(user)`.
  - `repository.go` — the repository interface the service depends on.
- **Repository** (`internal/repository/<resource>repository/`) — raw SQL via `sqlx`. Write operations receive a `*sqlx.Tx`; reads receive `*sqlx.DB`. `Create` assigns a new UUID to `record.ID`. SQL strings are in a companion `sql.go` file.
- **Model** (`internal/model/<resource>/`) — struct embedding `model.BaseModel`, plus `Validate()` returning `req.ErrMsgs` and optionally `Normalise()` (trim/lower-case input before validation).

See `author` (handler, service, record, repository, model, tests) for the reference implementation.

**IDs:** all primary and foreign keys are UUIDv7 strings stored as `CHAR(36)`, generated in the repository's `Create` with `repository.NewID()` (`internal/repository/id.go`). Use it for every new id rather than calling `uuid` directly, so the version stays in one place.

**Routing:** Routes are registered in `internal/routes/routes.go` using `httprouter-group`. All routes are prefixed `/v1`. Protected routes use `authMW` (`middleware.AuthMiddleware`), which validates the Bearer session token and puts the user ID (`app.UserIDKey`) and the `*userservice.UserRecord` (`app.UserRecordKey`) in the request context. To add a new resource, register its group in `GetV1Routes()`.

**Authorisation:** global roles live in `user_role` (IDs in `internal/model/userrole`); each user has one `role_id`. Restrict a route group with `middleware.RequireRoles(roleIDs...)`, always composed inside `authMW`. Record-level checks are methods on the record (e.g. `AuthorRecord.CanBeAccessedBy`), called from the handler with `h.base.GetRequestingUserRecord(ctx)`.

**BaseModel** (`internal/model/basemodel.go`) — all models embed this. It provides `id`, `created_by`, `updated_by`, `created_at`, `updated_at`, `deleted_at`. Soft-deletes are handled by setting `deleted_at`. `SetCreated`/`SetUpdated`/`SetDeleted` pull the authenticated user ID from the request context via `app.UserIDKey`.

**Validation** (`internal/validator/`) — field-level validation rules are declared as `[]validator.Rule` on the model. Each rule references one or more `ValidationFuncs`. `validator.RunValidation` returns a map of field → error messages. `validator.RecordExistsByID(othermodel.Other{}, "")` checks a referenced (non-deleted) row exists, using the snake_case of the model's type name as the table; it needs `validator.Init(db)` at startup (done in `cmd/server/run.go`). When a request payload differs from the stored model (e.g. a public submission form), give it its own request struct with its own rules rather than loosening the model's.

**Request/Response** (`internal/request/`) — `GetRequestParams` parses pagination, filtering, and sorting from query params. Filters are declared per-handler as `[]req.FilterRule`; sortable fields are whitelisted. Response helpers (`GetSingleItemResp`, `GetListResp`, `GetValidateErrResp`, etc.) standardise the JSON envelope.

**Migrations** (`migrations/`) — raw SQL files named `<timestamp>_<name>.up.sql` / `.down.sql` where timestamp is `YYYYMMDDHHmmss` (14 digits, no separators). Multiple statements per file are separated by `;`. Tracking is via a `schema_migrations` table managed by [golang-migrate/migrate](https://github.com/golang-migrate/migrate). If a migration fails mid-run the DB is marked dirty; resolve with `migrate force <version>` before retrying.

**Scaffolding** (`internal/scaffold/`, `cmd/scaffold/`) — reads a live DB table schema and generates the four layers above (service, record and repository interface included). The table must use `CHAR(36)` UUID ids and the standard audit columns. After generating, fill in the record's `Update` (it copies no fields until you do) and refine the generated validation rules. Output paths are configured in `cmd/scaffold/config.json`.

**Testing:** handler tests start the server as a subprocess on port `0` (OS-assigned, so packages can run in parallel) via `testtools.StartTestServer(cfg)`, run migrations, and create users with `testtools.CreateTestUser` / `CreateTestUserWithRole`. They need the test database from `config_test.json`.
