// Package repositoryauthor provides database operations for author records.
package repositoryauthor

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model/author"
	"github.com/doublehops/dh-go-framework/internal/repository"
	req "github.com/doublehops/dh-go-framework/internal/request"
)

// Author handles database operations for author records.
type Author struct {
	Log *logga.Logga
}

// New creates a new author Author repository with the provided logger.
func New(logger *logga.Logga) *Author {
	return &Author{
		Log: logger,
	}
}

// Create inserts a new author record within the provided transaction.
func (a *Author) Create(ctx context.Context, tx *sqlx.Tx, record *author.Author) error {
	result, err := tx.NamedExec(insertRecordSQL, record)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		a.Log.Error(ctx, errMsg, nil)

		return errors.New(errMsg)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	record.ID = int32(lastInsertID) //nolint:gosec

	return nil
}

// Update saves changes to an existing author record within the provided transaction.
func (a *Author) Update(ctx context.Context, tx *sqlx.Tx, model *author.Author) error {
	_, err := tx.NamedExec(updateRecordSQL, model)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		a.Log.Error(ctx, errMsg, nil)

		return errors.New(errMsg)
	}

	return nil
}

// Delete soft-deletes an author record within the provided transaction.
func (a *Author) Delete(ctx context.Context, tx *sqlx.Tx, model *author.Author) error {
	_, err := tx.NamedExec(deleteRecordSQL, model)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		a.Log.Error(ctx, errMsg, nil)

		return errors.New(errMsg)
	}

	return nil
}

// GetByID retrieves an author record by its primary key.
func (a *Author) GetByID(ctx context.Context, db *sqlx.DB, id int32, record *author.Author) error {
	err := db.Get(record, selectByIDQuery, id)
	if err != nil {
		a.Log.Error(ctx, "unable to fetch record", logga.KVPs{"ID": id})

		return fmt.Errorf("unable to fetch record %d", id)
	}

	return nil
}

// GetCollection retrieves a paginated collection of author records.
func (a *Author) GetCollection(ctx context.Context, db *sqlx.DB, p *req.Request) ([]*author.Author, error) {
	var (
		records []*author.Author
		err     error
	)

	countQ, countParams := repository.BuildQuery(selectCollectionCountQuery, p, true)
	count, err := repository.GetRecordCount(db, countQ, countParams)
	if err != nil {
		a.Log.Error(ctx, "GetCollection()", logga.KVPs{"err": err})
	}
	p.SetRecordCount(count)

	q, params := repository.BuildQuery(selectCollectionQuery, p, false)
	err = db.Select(&records, q, params...)
	if err != nil {
		return records, fmt.Errorf("unable to retrieve records: %s", err.Error())
	}

	return records, nil
}
