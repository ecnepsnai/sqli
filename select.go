package sqli

import (
	"fmt"
	"strings"
)

// SelectSingle select a single row from the database. Will return (nil, nil) if nothing found.
func (d *Database) SelectSingle(query SelectQuery) (map[string]interface{}, error) {
	rows, err := d.Select(query)
	if err != nil {
		d.log.Error("Unable to execute select query: %s", err)
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return rows[0], nil
}

// Select select multiple rows from the database. Will return (nil, nil) if nothing found.
func (d *Database) Select(query SelectQuery) ([]map[string]interface{}, error) {
	rows, err := d.query(query.sql())
	if err != nil {
		d.log.Error("Unable to execute select query: %s", err)
		return nil, err
	}

	var results []map[string]interface{}

	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			d.log.Error("Error scanning rows into map: %s", err)
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		results = append(results, m)
	}

	return results, nil
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
