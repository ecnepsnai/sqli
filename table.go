package sqli

import (
	"fmt"
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
	// TypeText a text column type that cannot be a primary key
	TypeText = "TEXT"
	// TypeString a string column type
	TypeString = "STRING"
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
	Length        uint
	NotNull       bool
	PrimaryKey    bool
	AutoIncrement bool
	Unique        bool
	Default       interface{}
}

// CreateTable create a new table in the database. Will not fail if the table already exists
func (d *Database) CreateTable(table Table) error {
	var sql string
	var err error
	if d.dbType == ServiceSQLite {
		sql, err = table.sqlite()
	} else if d.dbType == ServiceMySQL {
		sql, err = table.mysql()
	}

	if err != nil {
		d.log.Error("Invalid create table parameters specified: %s", err.Error())
		return err
	}

	_, err = d.execute(sql)
	if err != nil {
		d.log.Error("Unable to create new table: %s", err)
		return err
	}

	return nil
}

func (table Table) sqlite() (string, error) {
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

	return sql, nil
}

func (table Table) mysql() (string, error) {
	sql := "CREATE TABLE IF NOT EXISTS `" + stripName(table.Name) + "` ("
	var columnStrings []string
	for _, column := range table.Columns {
		if column.PrimaryKey && column.Type == TypeText {
			return "", fmt.Errorf("TEXT type columns cannot be primary keys")
		}

		col := "`" + stripName(column.Name) + "` "
		if column.Type == TypeString {
			col += "VARCHAR"
		} else if column.Type == TypeInteger {
			col += "BIGINT"
		} else {
			col += column.Type
		}
		if column.Type == TypeInteger || column.Type == TypeString {
			col += fmt.Sprintf(" (%d) ", column.Length)
		}

		if column.NotNull {
			col += " NOT NULL "
		}
		if column.Default != nil && !column.PrimaryKey {
			col += " DEFAULT " + sanitizeValue(column.Default) + " "
		}
		if column.AutoIncrement {
			col += " AUTO_INCREMENT "
		}

		columnStrings = append(columnStrings, col)
	}

	for _, column := range table.Columns {
		if !column.PrimaryKey && !column.Unique {
			continue
		}
		if column.PrimaryKey && column.Unique {
			return "", fmt.Errorf("Column '%s' cannot be both a parimary key and unique", column.Name)
		}

		var index string
		if column.PrimaryKey {
			index = "PRIMARY KEY (`" + stripName(column.Name) + "`)"
		}
		if column.Unique {
			index = "UNIQUE KEY `" + stripName(column.Name) + "` (`" + stripName(column.Name) + "`)"
		}

		columnStrings = append(columnStrings, index)
	}

	sql += strings.Join(columnStrings, ", ")
	sql += ");"

	return sql, nil
}

// ColumnByName get the column with the specified name
func (table Table) ColumnByName(column string) *Column {
	for _, c := range table.Columns {
		if c.Name == column {
			return &c
		}
	}

	return nil
}
