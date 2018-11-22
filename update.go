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
	var columnNames []string
	var columnValues []string
	for _, column := range columns {
		value := sanitizeValue(q.Values[column])
		if len(value) > 0 {
			columnNames = append(columnNames, stripName(column))
			columnValues = append(columnValues, value)
		}
	}

	updateStrings := make([]string, len(columnNames))
	for i, columnName := range columnNames {
		updateStrings[i] = "`" + columnName + "` = " + columnValues[i]
	}

	sql += strings.Join(updateStrings, ",")

	if q.Where != nil {
		sql += " WHERE " + q.Where.sql()
	}

	sql += ";"

	return sql
}
