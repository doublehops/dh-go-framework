package repositoryauthor

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/doublehops/dh-go-framework/internal/logga"
	"github.com/doublehops/dh-go-framework/internal/model"
	"github.com/doublehops/dh-go-framework/internal/repository"
	req "github.com/doublehops/dh-go-framework/internal/request"
)

type Author struct {
	Log *logga.Logga
}

func New(logger *logga.Logga) *Author {
	return &Author{
		Log: logger,
	}
}

func (a *Author) Create(ctx context.Context, tx *sqlx.Tx, record *model.Author) error {
	result, err := tx.NamedExec(insertRecordSQL, record)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		a.Log.Error(ctx, errMsg, nil)

		return fmt.Errorf(errMsg)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	record.ID = int32(lastInsertID)

	return nil
}

func (a *Author) Update(ctx context.Context, tx *sqlx.Tx, model *model.Author) error {
	_, err := tx.NamedExec(updateRecordSQL, model)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		a.Log.Error(ctx, errMsg, nil)

		return fmt.Errorf(errMsg)
	}

	return nil
}

func (a *Author) Delete(ctx context.Context, tx *sqlx.Tx, model *model.Author) error {
	_, err := tx.NamedExec(deleteRecordSQL, model)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		a.Log.Error(ctx, errMsg, nil)

		return fmt.Errorf(errMsg)
	}

	return nil
}

func (a *Author) GetByID(ctx context.Context, DB *sqlx.DB, ID int32, record *model.Author) error {
	err := DB.Get(record, selectByIDQuery, ID)
	if err != nil {
		a.Log.Error(ctx, "unable to fetch record", logga.KVPs{"ID": ID})

		return fmt.Errorf("unable to fetch record %d", ID)
	}

	return nil
}

func (a *Author) GetCollection(ctx context.Context, DB *sqlx.DB, p *req.Request) ([]*model.Author, error) {
	var (
		records []*model.Author
		err     error
	)

	countQ, countParams := repository.BuildQuery(selectCollectionCountQuery, p, true)
	count, err := repository.GetRecordCount(DB, countQ, countParams)
	if err != nil {
		a.Log.Error(ctx, "GetCollection()", logga.KVPs{"err": err})
	}
	p.SetRecordCount(count)

	q, params := repository.BuildQuery(selectCollectionQuery, p, false)
	err = DB.Select(&records, q, params...)
	if err != nil {
		return records, fmt.Errorf("unable to retrieve records: %s", err.Error())
	}

	return records, nil
}
