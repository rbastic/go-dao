package godao

import (
	"context"
	"database/sql"
	"errors"
)

var ErrZeroRowsAffected = errors.New("zero rows affected")

// Exec executes a single SQL statement, returning the last insert ID and number of rows affected.
func Exec(ctx context.Context, db *sql.DB, tx *sql.Tx, stmt *sql.Stmt, sqlText string, args ...interface{}) (lastID int64, numRows int64, err error) {

	if stmt == nil {
		if tx != nil {
			stmt, err = tx.PrepareContext(ctx, sqlText)
		} else {
			stmt, err = db.PrepareContext(ctx, sqlText)
		}
	}

	if err != nil {
		return -1, -1, err
	}

	result, err := stmt.ExecContext(ctx, args...)
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
