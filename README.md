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

To run:
```
make scaffold table=<table_name>
```

