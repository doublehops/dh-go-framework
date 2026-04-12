package userrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/repository"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
)

// Repo handles database operations for user records.
type Repo struct {
	l *logga.Logga
}

// New creates a new user Repo with the provided logger.
func New(logger *logga.Logga) *Repo {
	return &Repo{
		l: logger,
	}
}

// Create inserts a new user record within the provided transaction.
func (r *Repo) Create(ctx context.Context, tx *sqlx.Tx, record *user.User) error {
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

// Update saves changes to an existing user record within the provided transaction.
func (r *Repo) Update(ctx context.Context, tx *sqlx.Tx, record *user.User) error {
	_, err := tx.NamedExec(updateRecordSQL, record)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		r.l.Error(ctx, errMsg, nil)

		return errors.New(errMsg)
	}

	return nil
}

// Delete soft-deletes a user record within the provided transaction.
func (r *Repo) Delete(ctx context.Context, tx *sqlx.Tx, record *user.User) error {
	_, err := tx.NamedExec(deleteRecordSQL, record)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		r.l.Error(ctx, errMsg, nil)

		return errors.New(errMsg)
	}

	return nil
}

// GetByID retrieves a user record by its primary key.
func (r *Repo) GetByID(ctx context.Context, db *sqlx.DB, id int32, record *user.User) error {
	err := db.Get(record, selectByIDQuery, id)
	if err != nil {
		r.l.Error(ctx, service.UnableToRetrieveRecord, logga.KVPs{"ID": id})

		return fmt.Errorf("%s %d", service.UnableToRetrieveRecord, id)
	}

	return nil
}

// GetByEmailAddress retrieves a user record by their email address.
func (r *Repo) GetByEmailAddress(ctx context.Context, db *sqlx.DB, emailAddress string, record *user.User) error {
	err := db.Get(record, selectByEmailAddressQuery, emailAddress)
	if err != nil {
		r.l.Error(ctx, service.UnableToRetrieveRecord, logga.KVPs{"emailAddress": emailAddress})

		return fmt.Errorf("%s %s", service.UnableToRetrieveRecord, emailAddress)
	}

	return nil
}

// GetCollection retrieves a paginated collection of user records.
func (r *Repo) GetCollection(ctx context.Context, db *sqlx.DB, p *req.Request) ([]*user.User, error) {
	var (
		records []*user.User
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
