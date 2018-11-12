package sqli

// Delete delete row(s) from a table
func (d *Database) Delete(query DeleteQuery) error {
	sql := query.sql()

	_, err := d.execute(sql)
	if err != nil {
		d.log.Error("Unable to execute delete query: %s", err)
		return err
	}

	return nil
}

// DeleteQuery the delete query
type DeleteQuery struct {
	Table Table
	Where Where
}

func (q DeleteQuery) sql() string {
	sql := "DELETE FROM `" + stripName(q.Table.Name) + "` WHERE " + q.Where.sql() + ";"
	return sql
}
