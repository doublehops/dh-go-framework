package userrepository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model/user"
	"github.com/doublehops/dh-go-framework/internal/repository"
	req "github.com/doublehops/dh-go-framework/internal/request"
	"github.com/doublehops/dh-go-framework/internal/service"
)

type Repo struct {
	l *logga.Logga
}

func New(logger *logga.Logga) *Repo {
	return &Repo{
		l: logger,
	}
}

func (r *Repo) Create(ctx context.Context, tx *sqlx.Tx, record *user.User) error {
	result, err := tx.NamedExec(insertRecordSQL, record)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		r.l.Error(ctx, errMsg, nil)

		return fmt.Errorf(errMsg)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	record.ID = int32(lastInsertID)

	return nil
}

func (r *Repo) Update(ctx context.Context, tx *sqlx.Tx, record *user.User) error {
	_, err := tx.NamedExec(updateRecordSQL, record)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		r.l.Error(ctx, errMsg, nil)

		return fmt.Errorf(errMsg)
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, tx *sqlx.Tx, record *user.User) error {
	_, err := tx.NamedExec(deleteRecordSQL, record)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		r.l.Error(ctx, errMsg, nil)

		return fmt.Errorf(errMsg)
	}

	return nil
}

func (r *Repo) GetByID(ctx context.Context, DB *sqlx.DB, ID int32, record *user.User) error {
	err := DB.Get(record, selectByIDQuery, ID)
	if err != nil {
		r.l.Error(ctx, service.UnableToRetrieveRecord, logga.KVPs{"ID": ID})

		return fmt.Errorf("%s %d", service.UnableToRetrieveRecord, ID)
	}

	return nil
}

func (r *Repo) GetByEmailAddress(ctx context.Context, DB *sqlx.DB, emailAddress string, record *user.User) error {
	err := DB.Get(record, selectByEmailAddressQuery, emailAddress)
	if err != nil {
		r.l.Error(ctx, service.UnableToRetrieveRecord, logga.KVPs{"emailAddress": emailAddress})

		return fmt.Errorf("%s %s", service.UnableToRetrieveRecord, emailAddress)
	}

	return nil
}

func (r *Repo) GetCollection(ctx context.Context, DB *sqlx.DB, p *req.Request) ([]*user.User, error) {
	var (
		records []*user.User
		err     error
	)

	countQ, countParams := repository.BuildQuery(selectCollectionCountQuery, p, true)
	count, err := repository.GetRecordCount(DB, countQ, countParams)
	if err != nil {
		r.l.Error(ctx, "GetCollection()", logga.KVPs{"err": err})
	}
	p.SetRecordCount(count)

	q, params := repository.BuildQuery(selectCollectionQuery, p, false)
	err = DB.Select(&records, q, params...)
	if err != nil {
		return records, fmt.Errorf("%s: %s", service.UnableToRetrieveRecord, err.Error())
	}

	return records, nil
}
