package sqli

import (
	"strings"
)

// Update update a row in the database
func (d *Database) Update(query UpdateQuery) error {
	sql := query.sql()

	_, err := d.execute(sql)
	if err != nil {
		d.log.Error("Unable to perform update query: %s", err)
		return err
	}

	return nil
}

// UpdateQuery describes an update query
type UpdateQuery struct {
	Table  Table
	Values map[string]interface{}
	Where  Where
}

func (q UpdateQuery) sql() string {
	sql := "UPDATE `" + stripName(q.Table.Name) + "` SET "

	columns := mapKeys(q.Values)
	columnNames := make([]string, len(columns))
	columnValues := make([]string, len(columns))
	for i, column := range columns {
		columnNames[i] = stripName(column)
		columnValues[i] = sanitizeValue(q.Values[column])
	}

	updateStrings := make([]string, len(columns))
	for i, columnName := range columns {
		updateStrings[i] = "`" + stripName(columnName) + "` = " + sanitizeValue(q.Values[columnName])
	}

	sql += strings.Join(updateStrings, ",")

	if q.Where != nil {
		sql += " WHERE " + q.Where.sql()
	}

	sql += ";"

	return sql
}
