package sqli

import (
	"database/sql"

	"github.com/ecnepsnai/logtic"
)

// Database describes a sqli instance
type Database struct {
	log        *logtic.Source
	db         *sql.DB
	connection *Connection
	dbType     string
}

// Connection describes connection information to a SQL server
type Connection struct {
	Host     string
	Port     uint16
	Username string
	Password string
	Database string
}

const (
	// ServiceMySQL enum value for MySQL type SQL instances
	ServiceMySQL = "MySQL"
	// ServiceSQLite enum value for SQLite type SQL instances
	ServiceSQLite = "SQLite"
)

// Close close the database file
func (d *Database) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
