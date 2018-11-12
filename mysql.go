package sqli

import (
	"database/sql"
	"fmt"

	"github.com/ecnepsnai/logtic"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

// MySQL connects to a the specified MySQL server
func MySQL(connection Connection) (*Database, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", connection.Username, connection.Password, connection.Host, connection.Port, connection.Database))
	log := logtic.Connect("sql")
	if err != nil {
		log.Error("Unable to connect to MySQL server at '%s@%s:%d': %s", connection.Username, connection.Host, connection.Port, err.Error())
		return nil, err
	}

	database := Database{
		log:    log,
		db:     db,
		dbType: ServiceMySQL,
	}

	log.Info("Connected to MySQL server at '%s@%s:%d'", connection.Username, connection.Host, connection.Port)

	return &database, nil
}
