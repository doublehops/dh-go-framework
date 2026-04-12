package usersessionrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model/usersession"
	"github.com/doublehops/dh-go-framework/internal/repository"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
)

// Repo handles database operations for user session records.
type Repo struct {
	l *logga.Logga
}

// New creates a new user session Repo with the provided logger.
func New(logger *logga.Logga) *Repo {
	return &Repo{
		l: logger,
	}
}

// Create inserts a new user session record within the provided transaction.
func (r *Repo) Create(ctx context.Context, tx *sqlx.Tx, record *usersession.UserSession) error {
	result, err := tx.NamedExec(insertRecordSQL, record)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		r.l.Error(ctx, errMsg, nil)

		return errors.New(errMsg)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	record.ID = int32(lastInsertID) //nolint:gosec

	return nil
}

// UpdateLastRequest updates the last_request timestamp on a session record within the provided transaction.
func (r *Repo) UpdateLastRequest(ctx context.Context, tx *sqlx.Tx, record *usersession.UserSession) error {
	_, err := tx.NamedExec(updateLastRequestSQL, record)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error updating record in db. %s", err)
		r.l.Error(ctx, errMsg, nil)

		return errors.New(errMsg)
	}

	return nil
}

// Delete soft-deletes a user session record within the provided transaction.
func (r *Repo) Delete(ctx context.Context, tx *sqlx.Tx, record *usersession.UserSession) error {
	_, err := tx.NamedExec(deleteRecordSQL, record)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error deleting record in db. %s", err)
		r.l.Error(ctx, errMsg, nil)

		return errors.New(errMsg)
	}

	return nil
}

// GetByID retrieves a user session record by its primary key.
func (r *Repo) GetByID(ctx context.Context, db *sqlx.DB, id int32, record *usersession.UserSession) error {
	err := db.Get(record, selectByIDQuery, id)
	if err != nil {
		r.l.Error(ctx, service.UnableToRetrieveRecord, logga.KVPs{"ID": id})

		return fmt.Errorf("%s %d", service.UnableToRetrieveRecord, id)
	}

	return nil
}

// GetByToken retrieves a user session by its authentication token.
func (r *Repo) GetByToken(ctx context.Context, db *sqlx.DB, token string, record *usersession.UserSession) error {
	err := db.Get(record, selectByTokenQuery, token)
	if err != nil {
		r.l.Error(ctx, service.UnableToRetrieveRecord, logga.KVPs{"token": token})

		return fmt.Errorf("%s token %s", service.UnableToRetrieveRecord, token)
	}

	return nil
}

// GetByUserID retrieves a user session by the associated user ID.
func (r *Repo) GetByUserID(ctx context.Context, db *sqlx.DB, userID int32, record *usersession.UserSession) error {
	err := db.Get(record, selectByUserIDQuery, userID)
	if err != nil {
		r.l.Error(ctx, service.UnableToRetrieveRecord, logga.KVPs{"userID": userID})

		return fmt.Errorf("%s userID %d", service.UnableToRetrieveRecord, userID)
	}

	return nil
}

// GetCollection retrieves a paginated collection of user session records.
func (r *Repo) GetCollection(ctx context.Context, db *sqlx.DB, p *req.Request) ([]*usersession.UserSession, error) {
	var (
		records []*usersession.UserSession
		err     error
	)

	countQ, countParams := repository.BuildQuery(selectCollectionCountQuery, p, true)
	count, err := repository.GetRecordCount(db, countQ, countParams)
	if err != nil {
		r.l.Error(ctx, "GetCollection()", logga.KVPs{"err": err})
	}
	p.SetRecordCount(count)

	q, params := repository.BuildQuery(selectCollectionQuery, p, false)
	err = db.Select(&records, q, params...)
	if err != nil {
		return records, fmt.Errorf("%s: %s", service.UnableToRetrieveRecord, err.Error())
	}

	return records, nil
}

// SetLastRequestNow updates the session's last_request field to the current time via a direct DB call.
func (r *Repo) SetLastRequestNow(ctx context.Context, db *sqlx.DB, record *usersession.UserSession) error {
	var err error

	_, err = db.NamedExec(updateLastRequestSQL, record)
	if err != nil {
		r.l.Error(ctx, "SetLastRequestNow()", logga.KVPs{"err": err})

		return err
	}

	return nil
}
