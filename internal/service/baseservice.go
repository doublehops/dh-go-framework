// Package service provides shared service-level utilities and base types for business logic layers.
package service

import (
	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model"
)

// UnableToSaveRecord and related constants are common error message strings used across services.
const (
	UnableToSaveRecord        = "unable to save record"
	UnableToUpdateRecord      = "unable to update record"
	UnableToRetrieveRecord    = "unable to retrieve record"
	UnableToCommitTransaction = "unable to commit transaction"
)

// App holds shared application dependencies (DB and logger) injected into all services.
type App struct {
	DB  *sqlx.DB
	Log *logga.Logga
}

// HasPermission will check whether the authenticated user has authorisation for the requested record. This function
// can be overwritten in each service.
func (a *App) HasPermission(userID int32, record model.Model) bool {
	return userID == record.GetUserID()
}
