package godao

import (
	"context"
	"database/sql"
)

// Query returns a []map[string]*interface{} result set or an error, for a
// given query and a set of arguments. Transactions are optionally supported -
// nontransacted queries should just pass nil for the third parameter.
func Query(ctx context.Context, db *sql.DB, tx *sql.Tx, sqlText string, args ...interface{}) (results []map[string]interface{}, err error) {
	var stmt *sql.Stmt

	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, sqlText)
	} else {
		stmt, err = db.PrepareContext(ctx, sqlText)
	}

	if err != nil {
		return nil, err
	}

	// trying to be clever here, like ... if len(args) == 0 then pass nothing,
	// doesn't work...
	var rows * sql.Rows
	if len(args) > 0 {
		rows, err = stmt.QueryContext(ctx, args...)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = stmt.QueryContext(ctx)
		if err != nil {
			return nil, err
		}
	}

	defer rows.Close()

	// NOTE(rbastic): I've tried many different versions of the code below
	// over the past 2-3 years, so far, this looks like the cleanest.
	//
	// https://kylewbanks.com/blog/query-result-to-map-in-golang

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		results = append(results, m)
	}

	err = rows.Err()

	return results, nil
}
