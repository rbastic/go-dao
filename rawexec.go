package godao

import (
	"context"
	"database/sql"
)

// RawExec executes a single SQL statement, returning the last insert ID and number of rows affected.
func RawExec(ctx context.Context, db *sql.DB, sqlText string, args ...interface{}) (lastID int64, numRows int64, err error) {

	result, err := db.ExecContext(ctx, sqlText, args...)
	if err != nil {
		return -1, -1, err
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return -1, -1, err
	}
	if rowsAff == 0 {
		return -1, -1, ErrZeroRowsAffected
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return -1, -1, err
	}

	return lastID, rowsAff, nil
}
