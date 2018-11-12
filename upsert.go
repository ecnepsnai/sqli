package sqli

import (
	"fmt"
	"strings"
)

// Upsert a combonation of INSERT OR IGNORE and UPDATE. Inserts a row if it doesn't already exist otherwise updating it.
// A unique or primary key column must be included in the query
func (d *Database) Upsert(query InsertQuery) error {
	var sql string
	var err error
	if d.dbType == ServiceSQLite {
		sql, err = d.upsertSQLite(query)
	} else if d.dbType == ServiceMySQL {
		sql, err = d.upsertMySQL(query)
	}
	if err != nil {
		d.log.Error("Invalid insert query provided for upsert: %s", err.Error())
		return err
	}

	_, err = d.execute(sql)
	if err != nil {
		d.log.Error("Unable to perform INSERT and UPDATE transaction: %s", err)
		return err
	}

	return nil
}

func (d *Database) upsertSQLite(query InsertQuery) (string, error) {
	sql := "BEGIN TRANSACTION;"
	query.IgnoreConflict = true
	sql += query.sql(d)

	var whereColumn string
	var whereValue interface{}

	columns := mapKeys(query.Values)
	for _, columnName := range columns {
		column := query.Table.ColumnByName(columnName)
		if column.PrimaryKey || column.Unique {
			whereColumn = columnName
			whereValue = query.Values[columnName]
			break
		}
	}

	if whereColumn == "" {
		return "", fmt.Errorf("Cannot upsert without including a primary key or unique column value")
	}

	update := UpdateQuery{
		Table:  query.Table,
		Values: query.Values,
		Where: Where{
			WhereEqual(whereColumn, whereValue),
		},
	}

	sql += update.sql()
	sql += "END TRANSACTION;"

	return sql, nil
}

func (d *Database) upsertMySQL(query InsertQuery) (string, error) {
	query.IgnoreConflict = false
	sql := query.sql(d)

	// Append to the insert query
	sql = strings.TrimSuffix(sql, ";")

	sql += " ON DUPLICATE KEY UPDATE "

	columns := mapKeys(query.Values)
	columnNames := make([]string, len(columns))
	columnValues := make([]string, len(columns))
	for i, column := range columns {
		columnNames[i] = stripName(column)
		columnValues[i] = sanitizeValue(query.Values[column])
	}

	updateStrings := make([]string, len(columns))
	for i, columnName := range columns {
		updateStrings[i] = "`" + stripName(columnName) + "` = " + sanitizeValue(query.Values[columnName])
	}

	sql += strings.Join(updateStrings, ",")

	return sql, nil
}
