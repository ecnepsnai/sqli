package sqli

import (
	"database/sql"
	"fmt"
)

// Row describes a single row in a query result
type Row struct {
	row  *sql.Row
	rows *sql.Rows
}

// Scan scan values for each column in the row in to the provided pointers
func (r Row) Scan(dest ...interface{}) error {
	if r.row != nil {
		return r.row.Scan(dest...)
	}
	if r.rows != nil {
		return r.rows.Scan(dest...)
	}

	return fmt.Errorf("nothing to scan")
}
