package godao

import (
	"fmt"
	"strings"
)

// INSERT will generate a INSERT statement.
func INSERT(tableName string, columns []string, rows int) string {
	return fmt.Sprintf("INSERT INTO %s ( %s ) VALUES %s", tableName, strings.Join(columns, ","), Ph(len(columns), rows))
}

// Ph generates the placeholders for SQL queries.
// For a bulk insert operation, rows is the number of rows you intend
// to insert, and columnsN is the number of fields per row.
func Ph(columnsN, rows int) string {

	inner := "( " + strings.TrimSuffix(strings.Repeat("?,", columnsN), ",") + " ),"
	return strings.TrimSuffix(strings.Repeat(inner, rows), ",")

}
