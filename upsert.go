package sqli

import (
	"fmt"
)

// Upsert a combonation of INSERT OR IGNORE and UPDATE. Inserts a row if it doesn't already exist otherwise updating it.
// A unique or primary key column must be included in the query
func (d *Database) Upsert(query InsertQuery) error {
	sql := "BEGIN TRANSACTION;"
	query.IgnoreConflict = true
	sql += query.sql()

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
		return fmt.Errorf("Cannot upsert without including a primary key or unique column value")
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

	_, err := d.executeNonQuery(sql)
	if err != nil {
		d.log.Error("Unable to perform INSERT and UPDATE transaction: %s", err)
		return err
	}

	return nil
}
