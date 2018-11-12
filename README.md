# sqli

SQLi is an SQL interface in Golang. It masks the SQL away form your application and allows you to use
multiple different types of SQL providers (MySQL, PostgreSQL*, and SQLite).

*PostgreSQL support coming soon.

# Usage

## Connect

### SQLite

```golang
db, err := sqli.SQLite("database.db")
if err != nil {
	fmt.Printf("Unable to open sqlite db: %s\n", err)
	os.Exit(1)
}
defer db.Close()
```

### MySQL

```golang
db, err := sqli.MySQL(sqli.Connection{
	Host: "127.0.0.1",
	Port: 3306,
	Username: "sqli_user",
	Password: "not this",
	Database: "sqli",
})
if err != nil {
	fmt.Printf("Unable to connect to mysql server: %s\n", err)
	os.Exit(1)
}
defer db.Close()
```

## Make a Table

```golang
table := sqli.Table{
	Name: "Test",
	Columns: []sqli.Column{
		Column{
			Name:          "id",
			Type:          sqli.TypeInteger,
			NotNull:       true,
			PrimaryKey:    true,
			AutoIncrement: true,
			Default:       0,
		},
		sqli.Column{
			Name:    "value",
			Type:    sqli.TypeText,
			NotNull: true,
		},
	},
}
err := db.CreateTable(table)
if err != nil {
	fmt.Printf("Error creating table: %s", err.Error())
}
```

## Insert

```golang
err := db.Insert(sqli.InsertQuery{
	Table: table,
	Values: map[string]interface{}{
		"id":    0,
		"value": "hello world",
	},
})
if err != nil {
	fmt.Printf("Error inserting row in table: %s", err.Error())
}
```

### Upsert

Use `sqli.Database.Upsert` to "UPSERT" a row. UPSERT is the term used to describe an INSERT and UPDATE query. Upsert
takes the same parameters as Insert, but requires that you include at least one Unique or Primary Key value.

## Update

```golang
err = db.Update(UpdateQuery{
	Table: table,
	Values: map[string]interface{}{
		"value": "hello again!",
	},
	Where: Where{
		WhereEqual("id", 0),
	},
})
if err != nil {
	fmt.Printf("Error updating row in table: %s", err)
}
```

## Select

### Single Row

```golang
row := db.SelectSingle(sqli.SelectQuery{
	Table: table,
	Where: sqli.Where{
		sqli.WhereEqual("id", 1),
	},
})
type exampleData struct{
	id    int
	value string
}
data := exampleData{}
if err := row.Scan(&data.id, &data.value); err != nil {
	fmt.Printf("Error selecting single row: %s", err.Error())
}
```

### Multiple Rows

```golang
type exampleData struct{
	id    int
	value string
}
var results []exampleData
err = db.Select(sqli.SelectQuery{
	Table: table,
	Columns: []string{
		"value",
	},
	Order: sqli.Order{
		Column:     "id",
		Descending: true,
	},
	Where: sqli.Where{
		sqli.WhereGreaterThan("id", expectedID),
		"AND",
		sqli.WhereNotEqual("id", expectedID+10),
	},
	Limit: 50,
}, func(row Row) error {
	to := exampleData{}
	if err := row.Scan(&to.value); err != nil {
		fmt.Printf("Error selecting multilpe row: %s", err.Error())
		return err
	}
	results = append(results, to)
	return nil
})
if err != nil {
	fmt.Printf("Error selecting multilpe row: %s", err.Error())
}
```

## Delete

```golang
err = db.Delete(DeleteQuery{
	Table: table,
	Where: Where{
		WhereEqual("id", 500),
	},
})
if err != nil {
	fmt.Printf("Error inserting row in table: %s", err)
}
```