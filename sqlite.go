package sqli

import (
	"database/sql"

	"github.com/ecnepsnai/logtic"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

// SQLite open a new SQLite database at the provided path
func SQLite(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	log := logtic.Connect("sql")
	if err != nil {
		log.Error("Unable to open sqlite database file '%s': %s", path, err.Error())
		return nil, err
	}

	database := Database{
		log:    log,
		db:     db,
		dbType: ServiceSQLite,
	}

	log.Info("Connected to database file '%s'", path)

	return &database, nil
}
