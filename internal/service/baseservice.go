package service

import (
	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model"
)

const (
	UnableToSaveRecord        = "unable to save record"
	UnableToUpdateRecord      = "unable to update record"
	UnableToRetrieveRecord    = "unable to retrieve record"
	UnableToCommitTransaction = "unable to commit transaction"
)

type App struct {
	DB  *sqlx.DB
	Log *logga.Logga
}

// HasPermission will check whether the authenticated user has authorisation for the requested record. This function
// can be overwritten in each service.
func (a *App) HasPermission(userID int32, record model.Model) bool {
	return userID == record.GetUserID()
}
