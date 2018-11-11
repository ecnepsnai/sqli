package sqli

import (
	"database/sql"

	"github.com/ecnepsnai/logtic"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

// Database describes a sqli instance
type Database struct {
	log *logtic.Source
	db  *sql.DB
}

// Open open a new database file at the specified path
func Open(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	log := logtic.Connect("sql")
	if err != nil {
		log.Error("Unable to open sqlite database file '%s': %s", path, err)
		return nil, err
	}

	database := Database{
		log: log,
		db:  db,
	}

	log.Info("Connected to database file '%s'", path)

	return &database, nil
}

// Close close the database file
func (d *Database) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
