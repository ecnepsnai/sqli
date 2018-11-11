package sqli

import (
	"strings"
)

// Table describes a database table
type Table struct {
	Name    string
	Columns []Column
}

var (
	// TypeInteger an interger column type
	TypeInteger = "INTEGER"
	// TypeText a text column type
	TypeText = "TEXT"
	// TypeBlob a blob column type
	TypeBlob = "BLOB"
	// TypeReal a real column type
	TypeReal = "REAL"
	// TypeNumeric a numberic column type
	TypeNumeric = "NUMERIC"
)

// Column describes a table column
type Column struct {
	Name          string
	Type          string
	NotNull       bool
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	Default       interface{}
}

// CreateTable create a new table in the database. Will not fail if the table already exists
func (d *Database) CreateTable(table Table) error {
	sql := "CREATE TABLE IF NOT EXISTS `" + stripName(table.Name) + "` ("
	columnStrings := make([]string, len(table.Columns))
	for i, column := range table.Columns {
		col := "`" + stripName(column.Name) + "` " + column.Type
		if column.NotNull {
			col += " NOT NULL "
		}
		if column.Default != nil {
			col += " DEFAULT " + sanitizeValue(column.Default) + " "
		}
		if column.PrimaryKey {
			col += " PRIMARY KEY "
		}
		if column.AutoIncrement {
			col += " AUTOINCREMENT "
		}
		columnStrings[i] = col
	}

	sql += strings.Join(columnStrings, ", ")
	sql += ");"

	_, err := d.executeNonQuery(sql)
	if err != nil {
		d.log.Error("Unable to create new table: %s", err)
		return err
	}

	return nil
}

// ColumnByName get the column with the specified name
func (t Table) ColumnByName(column string) *Column {
	for _, c := range t.Columns {
		if c.Name == column {
			return &c
		}
	}

	return nil
}
