package sqli

import (
	"fmt"
	"strings"
)

// SelectSingle select a single row from the database.
func (d *Database) SelectSingle(query SelectQuery) Row {
	sql := query.sql()
	d.log.Debug("%s", sql)
	return Row{
		row: d.db.QueryRow(sql),
	}
}

// Select select multiple rows from the database. Calls row() for each row returned
func (d *Database) Select(query SelectQuery, row func(row Row) error) error {
	sql := query.sql()
	d.log.Debug("%s", sql)
	rows, err := d.db.Query(sql)
	if err != nil {
		d.log.Error("Error performing SELECT query: %s", err.Error())
		return err
	}

	for rows.Next() {
		rowErr := row(Row{
			rows: rows,
		})
		if rowErr != nil {
			d.log.Error("Error during row handling: %s", rowErr.Error())
			return rowErr
		}
	}
	rows.Close()

	return nil
}

// SelectQuery describes a select query
type SelectQuery struct {
	Table   Table
	Columns []string
	Where   Where
	Order   Order
	Limit   uint
}

// Order describes an order clause for a select query
type Order struct {
	Column     string
	Descending bool
}

func (q SelectQuery) sql() string {
	sql := "SELECT "
	if q.Columns == nil {
		sql += " * "
	} else {
		columnNames := make([]string, len(q.Columns))
		for i, columnName := range q.Columns {
			columnNames[i] = "`" + stripName(columnName) + "`"
		}

		sql += " (" + strings.Join(columnNames, ",") + ") "
	}

	sql += " FROM `" + stripName(q.Table.Name) + "` "

	if q.Where != nil {
		sql += " WHERE " + q.Where.sql()
	}

	if q.Order.Column != "" {
		sql += q.Order.sql()
	}

	if q.Limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d ", q.Limit)
	}

	sql += ";"

	return sql
}

func (o Order) sql() string {
	sql := " ORDER BY `" + stripName(o.Column) + "` "
	if o.Descending {
		sql += " DESC "
	} else {
		sql += " ASC "
	}

	return sql
}
