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
- CRUD scaffolding
- user model, login and authentication, etc...
- ci/cd
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
- `./internal/migrations/` # Contains database migration files
- `./internal/middleware/` # Contains API middleware
- `./config.json` # Application configuration