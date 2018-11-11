package sqli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

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
			NotNull:       true,
			PrimaryKey:    true,
			AutoIncrement: true,
			Default:       0,
		},
		Column{
			Name:    "value",
			Type:    TypeText,
			NotNull: true,
		},
	},
}

func setupTest() {
	tmpDir, err := ioutil.TempDir("", "sqlite")
	if err != nil {
		fmt.Printf("Unable to make temporary directory: %s\n", err)
		os.Exit(1)
	}
	tempDir = tmpDir

	d, err := Open(path.Join(tempDir, "sqlite.db"))
	if err != nil {
		fmt.Printf("Unable to open sqlite db: %s\n", err)
		os.Exit(1)
	}
	db = d
}

func testdownTest() {
	db.Close()
	os.RemoveAll(tempDir)
}

func TestMain(m *testing.M) {
	setupTest()
	retCode := m.Run()
	testdownTest()
	os.Exit(retCode)
}

func TestCreateTable(t *testing.T) {
	err := db.CreateTable(table)
	if err != nil {
		t.Errorf("Error creating table: %s", err)
		t.Fail()
	}
}

func TestInsert(t *testing.T) {
	err := db.Insert(InsertQuery{
		Table: table,
		Values: map[string]interface{}{
			"id":    0,
			"value": "hello world",
		},
	})
	if err != nil {
		t.Errorf("Error inserting row in table: %s", err)
		t.Fail()
	}
}

func TestUpsert(t *testing.T) {
	err := db.Insert(InsertQuery{
		Table: table,
		Values: map[string]interface{}{
			"id":    100,
			"value": "hello world",
		},
	})
	if err != nil {
		t.Errorf("Error inserting row in table: %s", err)
		t.Fail()
	}

	err = db.Upsert(InsertQuery{
		Table: table,
		Values: map[string]interface{}{
			"id":    100,
			"value": "hello again",
		},
	})
	if err != nil {
		t.Errorf("Error inserting row in table: %s", err)
		t.Fail()
	}
}

func TestUpdate(t *testing.T) {
	rowID := 800

	err := db.Insert(InsertQuery{
		Table: table,
		Values: map[string]interface{}{
			"id":    rowID,
			"value": "hello world",
		},
	})
	if err != nil {
		t.Errorf("Error inserting row in table: %s", err)
		t.Fail()
	}

	expectedValue := "hello again!"

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
		t.Fail()
	}

	data, err := db.SelectSingle(SelectQuery{
		Table: table,
		Where: Where{
			WhereEqual("id", rowID),
		},
	})
	if err != nil {
		t.Errorf("Error selecting single row: %s", err)
		t.Fail()
	}
	if data == nil {
		t.Error("No rows returned from SELECT")
		t.Fail()
	}
	returnValue := string(data["value"].([]byte))
	if returnValue != expectedValue {
		t.Errorf("Incorrect data returned from SELECT. Expected '%s' got '%s'", expectedValue, returnValue)
		t.Fail()
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
		t.Fail()
	}

	err = db.Delete(DeleteQuery{
		Table: table,
		Where: Where{
			WhereEqual("id", 500),
		},
	})
	if err != nil {
		t.Errorf("Error inserting row in table: %s", err)
		t.Fail()
	}
}

func TestSelect(t *testing.T) {
	var expectedID int64 = 500

	err := db.InsertMany([]InsertQuery{
		InsertQuery{
			Table: table,
			Values: map[string]interface{}{
				"id":    expectedID,
				"value": "find me",
			},
		},
		InsertQuery{
			Table: table,
			Values: map[string]interface{}{
				"id":    expectedID + 10,
				"value": "find me",
			},
		},
		InsertQuery{
			Table: table,
			Values: map[string]interface{}{
				"id":    expectedID + 20,
				"value": "find me",
			},
		},
		InsertQuery{
			Table: table,
			Values: map[string]interface{}{
				"id":    expectedID + 30,
				"value": "find me",
			},
		},
	})
	if err != nil {
		t.Errorf("Error inserting row in table: %s", err)
		t.Fail()
	}

	data, err := db.SelectSingle(SelectQuery{
		Table: table,
		Where: Where{
			WhereEqual("id", expectedID),
		},
	})
	if err != nil {
		t.Errorf("Error selecting single row: %s", err)
		t.Fail()
	}
	if data == nil {
		t.Error("No rows returned from SELECT")
		t.Fail()
	}
	if data["id"] != expectedID {
		t.Errorf("Incorrect data returned from SELECT. Expected %d got %d", expectedID, data["id"])
		t.Fail()
	}

	results, err := db.Select(SelectQuery{
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
	})
	if err != nil {
		t.Errorf("Error selecting multilpe row: %s", err)
		t.Fail()
	}

	if len(results) < 1 {
		t.Error("No rows returned from SELECT")
		t.Fail()
	}
}
