package sqli

import (
	"database/sql"
	"fmt"
)

func (d *Database) executeNonQuery(sql string) (sql.Result, error) {
	fmt.Printf("%s\n", sql)
	d.log.Debug("%s", sql)
	return d.db.Exec(sql, nil)
}

func (d *Database) query(sql string) (*sql.Rows, error) {
	fmt.Printf("%s\n", sql)
	d.log.Debug("%s", sql)
	return d.db.Query(sql)
}
