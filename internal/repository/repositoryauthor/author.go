package repositoryauthor

import (
	"context"
	"database/sql"
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

func (a *Author) Update(ctx context.Context, tx *sql.Tx, model *model.Author) error {
	_, err := tx.Exec(updateRecordSQL, model.Name, model.UpdatedBy, model.UpdatedAt, model.ID)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		a.Log.Error(ctx, errMsg, nil)

		return fmt.Errorf(errMsg)
	}

	return nil
}

func (a *Author) Delete(ctx context.Context, tx *sql.Tx, model *model.Author) error {
	_, err := tx.Exec(deleteRecordSQL, model.UpdatedBy, model.UpdatedAt, model.DeletedAt, model.ID)
	if err != nil {
		errMsg := fmt.Sprintf("there was an error saving record to db. %s", err)
		a.Log.Error(ctx, errMsg, nil)

		return fmt.Errorf(errMsg)
	}

	return nil
}

func (a *Author) GetByID(ctx context.Context, DB *sqlx.DB, ID int32, model *model.Author) error {
	err := DB.Select(model, selectByIDQuery, ID)
	if err != nil {
		a.Log.Info(ctx, "unable to fetch record", logga.KVPs{"ID": ID})

		return fmt.Errorf("unable to fetch record %d", ID)
	}

	return nil
}

// func (a *Author) GetAllXXX(ctx context.Context, DB *sqlx.DB, p *req.Request) ([]*model.Author, error) {
// 	var (
// 		authors []*model.Author
// 		rows    *sql.Rows
// 		err     error
// 	)
//
// 	countQ, countParams := repository.BuildQuery(selectCollectionCountQuery, p, true)
// 	count, err := repository.GetRecordCount(DB, countQ, countParams)
// 	if err != nil {
// 		a.Log.Error(ctx, "GetAll()", logga.KVPs{"err": err})
// 	}
// 	p.SetRecordCount(count)
//
// 	q, params := repository.BuildQuery(selectCollectionQuery, p, false)
//
// 	a.Log.Debug(ctx, "GetAll()", logga.KVPs{"query": q})
// 	if len(params) == 0 {
// 		rows, err = DB.Query(q)
// 	} else {
// 		rows, err = DB.Query(q, params...)
// 	}
// 	if err != nil {
// 		a.Log.Error(ctx, "GetAll() unable to fetch rows", logga.KVPs{"err": err})
//
// 		return authors, fmt.Errorf("unable to fetch rows")
// 	}
// 	defer rows.Close()
// 	if rows.Err() != nil {
// 		a.Log.Error(ctx, "error with rows.Err(). "+rows.Err().Error(), nil)
//
// 		return authors, err
// 	}
//
// 	for rows.Next() {
// 		var record model.Author
// 		if err = rows.Scan(&record.ID, &record.UserID, &record.Name, &record.CreatedBy, &record.UpdatedBy, &record.CreatedAt, &record.UpdatedAt); err != nil {
// 			return authors, fmt.Errorf("unable to fetch rows. %s", err)
// 		}
//
// 		authors = append(authors, &record)
// 	}
//
// 	return authors, nil
// }

func (a *Author) GetCollection(ctx context.Context, DB *sqlx.DB, p *req.Request) ([]*model.Author, error) {
	var (
		records []*model.Author
		// rows    *sql.Rows
		err error
	)

	countQ, countParams := repository.BuildQuery(selectCollectionCountQuery, p, true)
	count, err := repository.GetRecordCount(DB, countQ, countParams)
	if err != nil {
		a.Log.Error(ctx, "GetAll()", logga.KVPs{"err": err})
	}
	p.SetRecordCount(count)

	q, params := repository.BuildQuery(selectCollectionQuery, p, false)
	err = DB.Select(records, q, params)

	return records, nil

	// countQ, countParams := repository.BuildQuery(selectCollectionCountQuery, p, true)
	// var count int32 = 0
	// err = repository.GetRecordCount(DB, &count, countQ, countParams)
	// if err != nil {
	// 	a.Log.Error(ctx, "GetAll()", logga.KVPs{"err": err})
	// }
	// p.SetRecordCount(count)
	//
	// q, params := repository.BuildQuery(selectCollectionQuery, p, false)
	//
	// a.Log.Debug(ctx, "GetAll()", logga.KVPs{"query": q})
	// if len(params) == 0 {
	// 	rows, err = DB.Query(q)
	// } else {
	// 	rows, err = DB.Query(q, params...)
	// }
	// if err != nil {
	// 	a.Log.Error(ctx, "GetAll() unable to fetch rows", logga.KVPs{"err": err})
	//
	// 	return authors, fmt.Errorf("unable to fetch rows")
	// }
	// defer rows.Close()
	// if rows.Err() != nil {
	// 	a.Log.Error(ctx, "error with rows.Err(). "+rows.Err().Error(), nil)
	//
	// 	return authors, err
	// }
	//
	// for rows.Next() {
	// 	var record model.Author
	// 	if err = rows.Scan(&record.ID, &record.UserID, &record.Name, &record.CreatedBy, &record.UpdatedBy, &record.CreatedAt, &record.UpdatedAt); err != nil {
	// 		return authors, fmt.Errorf("unable to fetch rows. %s", err)
	// 	}
	//
	// 	authors = append(authors, &record)
	// }

	// return authors, nil
}
