package sqli

import (
	"strings"
)

// Insert insert a new row into the table
func (d *Database) Insert(query InsertQuery) error {
	sql := query.sql(d)

	_, err := d.execute(sql)
	if err != nil {
		d.log.Error("Unable to execute insert query: %s", err)
		return err
	}

	return nil
}

// InsertMany insert multiple new rows into the table
func (d *Database) InsertMany(queries []InsertQuery) error {
	rows := make([]string, len(queries))
	for i, query := range queries {
		rows[i] = query.sql(d)
	}

	_, err := d.executeMany(rows)
	if err != nil {
		d.log.Error("Unable to perform multiple INSERT transaction: %s", err)
		return err
	}

	return nil
}

// InsertQuery describes a insert query
type InsertQuery struct {
	Table          Table
	Values         map[string]interface{}
	IgnoreConflict bool
}

func (q InsertQuery) sql(d *Database) string {
	sql := "INSERT "
	if q.IgnoreConflict {
		if d.dbType == ServiceMySQL {
			sql += " IGNORE "
		} else if d.dbType == ServiceSQLite {
			sql += " OR IGNORE "
		}
	}

	sql += " INTO `" + stripName(q.Table.Name) + "`("

	columns := mapKeys(q.Values)
	columnNames := make([]string, len(columns))
	columnValues := make([]string, len(columns))
	for i, column := range columns {
		value := sanitizeValue(q.Values[column])
		if len(value) > 0 {
			columnNames[i] = stripName(column)
			columnValues[i] = value
		}
	}

	sql += strings.Join(columnNames, ",")
	sql += ") VALUES ("
	sql += strings.Join(columnValues, ",")
	sql += ") "

	sql += ";"

	return sql
}
