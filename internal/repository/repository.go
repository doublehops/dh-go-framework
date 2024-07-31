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
	// defer row.Close()
	// if row.Err() != nil {
	// 	return fmt.Errorf("error in row.Err(). " + row.Err().Error())
	// }
	//
	// if err != nil {
	// 	return fmt.Errorf("unable to run count query. %s", err)
	// }
	//
	// for row.Next() {
	// 	err = row.Scan(&count)
	// 	if err != nil {
	// 		return fmt.Errorf("unable to scan query result. %s", err)
	// 	}
	// }

	return c.Count, err
}
