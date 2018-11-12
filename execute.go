package sqli

import (
	"database/sql"
)

func (d *Database) execute(sql string) (sql.Result, error) {
	d.log.Debug("%s", sql)
	return d.db.Exec(sql)
}

func (d *Database) executeMany(queries []string) ([]sql.Result, error) {
	tx, err := d.db.Begin()
	if err != nil {
		d.log.Error("Unable to start transaction: %s", err.Error())
		return nil, err
	}

	results := make([]sql.Result, len(queries))
	for i, query := range queries {
		d.log.Debug("%s", query)
		result, err := tx.Exec(query)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		results[i] = result
	}

	if err := tx.Commit(); err != nil {
		d.log.Error("Unable to commit transaction: %s", err.Error())
		return nil, err
	}

	return results, nil
}
