package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type countRes struct {
	Count int32 `db:"count"`
}

// GetRecordCount will retrieve the number of records for a given query for pagination responses.
// The function expects only one column in the query. An example would be `SELECT COUNT(*) count FROM {table}1`.
func GetRecordCount(DB *sqlx.DB, q string, params []any) (int32, error) {
	var (
		err error
		// row *sql.Rows
		c countRes
	)
	if params == nil {
		err = DB.Get(&c, q)
		if err != nil {
			return c.Count, fmt.Errorf("unable to run query. %s", err)
		}
	} else {
		err = DB.Get(&c, q, params...)
		if err != nil {
			return c.Count, fmt.Errorf("unable to fetch row. %s", err)
		}
	}

	return c.Count, err
}
