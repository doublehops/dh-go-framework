// Package service provides shared service-level utilities and base types for business logic layers.
package service

import (
	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/logga"
)

// UnableToSaveRecord and related constants are common error message strings used across services.
const (
	UnableToSaveRecord        = "unable to save record"
	UnableToUpdateRecord      = "unable to update record"
	UnableToRetrieveRecord    = "unable to retrieve record"
	UnableToCommitTransaction = "unable to commit transaction"
)

// App holds shared application dependencies (DB and logger) injected into all services.
//
// Authorisation is not handled here. AuthMiddleware puts the requesting *userservice.UserRecord in
// the request context; route-level role checks use middleware.RequireRoles, and record-level checks
// are methods on the resource's record (e.g. AuthorRecord.CanBeAccessedBy).
type App struct {
	DB  *sqlx.DB
	Log *logga.Logga
}
