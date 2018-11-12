package sqli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/ecnepsnai/logtic"

	_ "github.com/mattn/go-sqlite3"
)

var db *Database
var tempDir string
var table = Table{
	Name: "Test",
	Columns: []Column{
		Column{
			Name:          "id",
			Type:          TypeInteger,
			Length:        16,
			NotNull:       true,
			PrimaryKey:    true,
			AutoIncrement: true,
			Default:       0,
		},
		Column{
			Name:    "value",
			Type:    TypeString,
			Length:  128,
			NotNull: true,
			Unique:  true,
		},
	},
}
var file *logtic.File
var log *logtic.Source

type testObject struct {
	id    int
	value string
}

func setupSQLite() {
	d, err := SQLite(path.Join(tempDir, "sqlite.db"))
	if err != nil {
		fmt.Printf("Unable to open sqlite db: %s\n", err)
		os.Exit(1)
	}
	db = d
}

func setupMySQL() {
	d, err := MySQL(Connection{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "root",
		Database: "crt",
	})
	if err != nil {
		fmt.Printf("Unable to connect to mysql db: %s\n", err)
		os.Exit(1)
	}
	db = d

	db.execute("DROP TABLE IF EXISTS `" + stripName(table.Name) + "`;")
}

func setupLog(verbose bool) {
	tmpDir, err := ioutil.TempDir("", "sqli")
	if err != nil {
		fmt.Printf("Unable to make temporary directory: %s\n", err)
		os.Exit(1)
	}
	tempDir = tmpDir

	level := logtic.LevelWarn
	if verbose {
		level = logtic.LevelDebug
	}

	f, s, err := logtic.New(path.Join(tmpDir, "sqli.log"), level, "test")
	if err != nil {
		fmt.Printf("Unable to open logtic instance: %s\n", err)
		os.Exit(1)
	}
	file = f
	log = s
}

func testdownTest() {
	file.Close()
	os.RemoveAll(tempDir)
}

func TestMain(m *testing.M) {
	verbose := false
	for _, arg := range os.Args {
		if arg == "-test.v=true" {
			verbose = true
		}
	}

	setupLog(verbose)
	log.Info("Running test suite with SQLite type DB")
	setupSQLite()
	retCode := m.Run()
	if retCode > 0 {
		testdownTest()
		os.Exit(retCode)
	}
	log.Info("Running test suite with MySQL type DB")
	setupMySQL()
	retCode = m.Run()
	testdownTest()
	os.Exit(retCode)
}

func TestCreateTable(t *testing.T) {
	err := db.CreateTable(table)
	if err != nil {
		t.Errorf("Error creating table: %s", err)
		t.FailNow()
	}
}

func TestInsert(t *testing.T) {
	err := db.Insert(InsertQuery{
		Table: table,
		Values: map[string]interface{}{
			"id":    0,
			"value": "insert test",
		},
	})
	if err != nil {
		t.Errorf("Error inserting row in table: %s", err)
		t.FailNow()
	}
}

func TestUpsert(t *testing.T) {
	err := db.Insert(InsertQuery{
		Table: table,
		Values: map[string]interface{}{
			"id":    100,
			"value": "2nd insert test",
		},
	})
	if err != nil {
		t.Errorf("Error inserting row in table: %s", err)
		t.FailNow()
	}

	err = db.Upsert(InsertQuery{
		Table: table,
		Values: map[string]interface{}{
			"id":    100,
			"value": "upsert test",
		},
	})
	if err != nil {
		t.Errorf("Error inserting row in table: %s", err)
		t.FailNow()
	}
}

func TestUpdate(t *testing.T) {
	rowID := 800

	err := db.Insert(InsertQuery{
		Table: table,
		Values: map[string]interface{}{
			"id":    rowID,
			"value": "update test",
		},
	})
	if err != nil {
		t.Errorf("Error inserting row in table: %s", err)
		t.FailNow()
	}

	expectedValue := "updated test"

	err = db.Update(UpdateQuery{
		Table: table,
		Values: map[string]interface{}{
			"value": expectedValue,
		},
		Where: Where{
			WhereEqual("id", rowID),
		},
	})
	if err != nil {
		t.Errorf("Error updating row in table: %s", err)
		t.FailNow()
	}

	row := db.SelectSingle(SelectQuery{
		Table: table,
		Where: Where{
			WhereEqual("id", rowID),
		},
	})
	data := struct {
		id    int
		value string
	}{}
	if err := row.Scan(&data.id, &data.value); err != nil {
		t.Errorf("Error selecting single row: %s", err)
		t.FailNow()
	}
	returnValue := data.value
	if returnValue != expectedValue {
		t.Errorf("Incorrect data returned from SELECT. Expected '%s' got '%s'", expectedValue, returnValue)
		t.FailNow()
	}
}

func TestDelete(t *testing.T) {
	err := db.Insert(InsertQuery{
		Table: table,
		Values: map[string]interface{}{
			"id":    500,
			"value": "delete me",
		},
	})
	if err != nil {
		t.Errorf("Error inserting row in table: %s", err)
		t.FailNow()
	}

	err = db.Delete(DeleteQuery{
		Table: table,
		Where: Where{
			WhereEqual("id", 500),
		},
	})
	if err != nil {
		t.Errorf("Error inserting row in table: %s", err)
		t.FailNow()
	}
}

func TestSelect(t *testing.T) {
	var expectedID int64 = 500

	err := db.InsertMany([]InsertQuery{
		InsertQuery{
			Table: table,
			Values: map[string]interface{}{
				"id":    expectedID,
				"value": "find me 1",
			},
		},
		InsertQuery{
			Table: table,
			Values: map[string]interface{}{
				"id":    expectedID + 10,
				"value": "find me 2",
			},
		},
		InsertQuery{
			Table: table,
			Values: map[string]interface{}{
				"id":    expectedID + 20,
				"value": "find me 3",
			},
		},
		InsertQuery{
			Table: table,
			Values: map[string]interface{}{
				"id":    expectedID + 30,
				"value": "find me 4",
			},
		},
	})
	if err != nil {
		t.Errorf("Error inserting row in table: %s", err)
		t.FailNow()
	}

	row := db.SelectSingle(SelectQuery{
		Table: table,
		Where: Where{
			WhereEqual("id", expectedID),
		},
	})

	data := testObject{}
	if err := row.Scan(&data.id, &data.value); err != nil {
		t.Errorf("Error selecting single row: %s", err)
		t.FailNow()
	}
	if data.value != "find me 1" {
		t.Errorf("Returned value was not correct. Expected '%s' got '%s'", "find me 1", data.value)
		t.FailNow()
	}
	var results []testObject
	err = db.Select(SelectQuery{
		Table: table,
		Columns: []string{
			"value",
		},
		Order: Order{
			Column:     "id",
			Descending: true,
		},
		Where: Where{
			WhereGreaterThan("id", expectedID),
			"AND",
			WhereNotEqual("id", expectedID+10),
		},
		Limit: 50,
	}, func(row Row) error {
		to := testObject{}
		if err := row.Scan(&to.value); err != nil {
			t.Errorf("Error selecting multilpe row: %s", err)
			t.FailNow()
			return err
		}
		results = append(results, to)
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting multilpe row: %s", err)
		t.FailNow()
	}

	expectedLength := 3
	gotLength := len(results)
	if gotLength != expectedLength {
		t.Errorf("Returned number of results was not correct. Expected %d got %d", expectedLength, gotLength)
		t.FailNow()
	}
}
