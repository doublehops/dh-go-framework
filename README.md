This is a RESTful API framework built in Golang that is hopefully well designed for
quick development and easy maintainability. It will work with dependencies made 
available in the repository `github.com/doublehops/dhapi`.

## Basic Requirements
The basic requirements from a RESTful API framework are:
- ~~parameters in URL - such as `api/user/123`~~
- ~~easy validation~~
- ~~database integration~~
- ~~database migrations~~
- ~~pagination~~
- `includes` parameter to include related models in response. 
- ~~ability to include prebuilt and custom middleware~~
- ~~easy filtering in collection requests~~
- ~~sorting and ordering in collection requests~~
- ~~CRUD scaffolding~~
- ~~user model, login and authentication, etc...~~
- ~~ci/cd~~
- documentation

Good to have:
- command line tool to list endpoints

## Directory Structure
- `./cmd/api/main.go` # Start API service
- `./cmd/migrate/migrate.go` # Run database migrations
- `./internal/routes/` # Contains definitions of API routes. One file per model
- `./internal/handlers/` # Contains handlers for incoming API requests
- `./internal/models/` # Contain data models
- `./internal/service/` # Service layer that contains business logic of each model/endpoints
- `./internal/repository/` # Contains data retrieval functions
- `./migrations/` # Contains database migration files
- `./internal/middleware/` # Contains API middleware
- `./config.json` # Application configuration

## Architecture

Each resource is split into handler → service → repository → DB. Services wire their own
repository in `New(app)`, so handlers only ever talk to the service. Single-record operations
live on a record wrapper, in an active-record style:

```go
a, err := authorSrv.GetByID(ctx, id)          // *AuthorRecord, error when not found
if !a.CanBeAccessedBy(requestingUser) { ... } // record-level permission
err = a.Update(ctx, payload)                  // copies editable fields, saves, reloads
err = a.Delete(ctx)                           // soft delete

created, err := authorSrv.NewRecord(&author.Author{Name: "Jo"}).Create(ctx)
```

All IDs are UUIDv7 strings (`CHAR(36)`), assigned by the repository on create via `repository.NewID()`.

Global roles (`standard_user`, `organisation_admin`, `master_admin`) are in the `user_role` table.
Restrict a route group with `middleware.RequireRoles(...)` inside the auth middleware. The
`/v1/user` endpoints are limited to admins; `GET /v1/me` returns the logged-in user.

## Security defaults

- `POST /v1/auth/login` is rate limited per client IP (10 attempts a minute) and takes the same
  time whether or not the email exists. Inactive users can't log in.
- Sessions expire after 24 hours idle and 7 days in total.
- Sorting is restricted to each handler's whitelisted fields, and `order` to `ASC`/`DESC`.
- The database password is never logged.

## Migrations

Migration files live in `./migrations/` and are named `<timestamp>_<name>.up.sql` / `<timestamp>_<name>.down.sql` where the timestamp is 14 digits (`YYYYMMDDHHmmss`). Multiple SQL statements in one file are separated by `;`. Run state is tracked in a `schema_migrations` table managed by [golang-migrate/migrate](https://github.com/golang-migrate/migrate).

**Create a new migration:**
```bash
go run ./cmd/migrate/migrate.go -action create -name <migration_name>
```

**Run all pending migrations:**
```bash
go run ./cmd/migrate/migrate.go -action up
```

**Run a specific number of migrations:**
```bash
go run ./cmd/migrate/migrate.go -action up -number 1
```

**Rollback one migration:**
```bash
go run ./cmd/migrate/migrate.go -action down
```

**Rollback a specific number of migrations:**
```bash
go run ./cmd/migrate/migrate.go -action down -number 2
```

If a migration fails mid-run the database is marked dirty. Resolve with:
```bash
migrate force <version>
```

## Scaffolding
Scaffolding is a tool that will read the table definition in the database and create the CRUD routes, handlers, service and
repository layers so that you can get started creating the new endpoints for a new model much faster than creating each file manually.
Make sure the database is running and app database settings configured.

The table must use a `CHAR(36)` `id` and the standard audit columns (`created_by`, `updated_by`,
`created_at`, `updated_at`, `deleted_at`). The generated model, repository, service (with record
and repository interface) and handler follow the `author` example. After generating, fill in the
record's `Update` method, which copies no fields until you do, and review the generated
validation rules.

To run:
```
make scaffold table=<table_name>
```

