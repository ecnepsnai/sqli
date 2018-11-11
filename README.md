# sqli

SQLi is an SQL interface in Golang. It obscures the SQL away form your application and (eventually) allows you to use
multiple different types of SQL providers (MySQL, PostgreSQL, and SQLite).

Only SQLite is supported right now. MySQL and PostgreSQL are to come.

# Usage

## Connect

```golang
db, err := sqli.Open("database.db")
if err != nil {
	fmt.Printf("Unable to open sqlite db: %s\n", err)
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
data, err := db.SelectSingle(SelectQuery{
	Table: table,
	Where: Where{
		WhereEqual("id", 0),
	},
})
if err != nil {
	fmt.Printf("Error selecting single row: %s", err)
}
```

### Multiple Rows

```golang
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
		WhereGreaterThan("id", 100),
		"AND",
		WhereNotEqual("id", 101),
	},
	Limit: 50,
})
if err != nil {
	fmt.Printf("Error selecting multilpe row: %s", err)
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