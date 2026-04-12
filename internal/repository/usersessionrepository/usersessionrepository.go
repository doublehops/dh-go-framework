package usersessionrepository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model/usersession"
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

func (r *Repo) Create(ctx context.Context, tx *sqlx.Tx, record *usersession.UserSession) error {
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

func (r *Repo) UpdateLastRequest(ctx context.Context, tx *sqlx.Tx, record *usersession.UserSession) error {
	_, err := tx.NamedExec(updateLastRequestSQL, record)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error updating record in db. %s", err)
		r.l.Error(ctx, errMsg, nil)

		return fmt.Errorf(errMsg)
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, tx *sqlx.Tx, record *usersession.UserSession) error {
	_, err := tx.NamedExec(deleteRecordSQL, record)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error deleting record in db. %s", err)
		r.l.Error(ctx, errMsg, nil)

		return fmt.Errorf(errMsg)
	}

	return nil
}

func (r *Repo) GetByID(ctx context.Context, DB *sqlx.DB, ID int32, record *usersession.UserSession) error {
	err := DB.Get(record, selectByIDQuery, ID)
	if err != nil {
		r.l.Error(ctx, service.UnableToRetrieveRecord, logga.KVPs{"ID": ID})

		return fmt.Errorf("%s %d", service.UnableToRetrieveRecord, ID)
	}

	return nil
}

func (r *Repo) GetByToken(ctx context.Context, DB *sqlx.DB, token string, record *usersession.UserSession) error {
	err := DB.Get(record, selectByTokenQuery, token)
	if err != nil {
		r.l.Error(ctx, service.UnableToRetrieveRecord, logga.KVPs{"token": token})

		return fmt.Errorf("%s token %s", service.UnableToRetrieveRecord, token)
	}

	return nil
}

func (r *Repo) GetByUserID(ctx context.Context, DB *sqlx.DB, userID int32, record *usersession.UserSession) error {
	err := DB.Get(record, selectByUserIDQuery, userID)
	if err != nil {
		r.l.Error(ctx, service.UnableToRetrieveRecord, logga.KVPs{"userID": userID})

		return fmt.Errorf("%s userID %d", service.UnableToRetrieveRecord, userID)
	}

	return nil
}

func (r *Repo) GetCollection(ctx context.Context, DB *sqlx.DB, p *req.Request) ([]*usersession.UserSession, error) {
	var (
		records []*usersession.UserSession
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

func (r *Repo) SetLastRequestNow(ctx context.Context, DB *sqlx.DB, record *usersession.UserSession) error {
	var (
		err error
	)

	_, err = DB.NamedExec(updateLastRequestSQL, record)
	if err != nil {
		r.l.Error(ctx, "SetLastRequestNow()", logga.KVPs{"err": err})

		return err
	}

	return nil
}
