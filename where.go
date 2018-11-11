package sqli

import (
	"reflect"
	"strings"
)

// Where describes a complete WHERE clause.
// Only `WhereCondition` and `string` types are allowed.
// String values can be SQL operators like:
// AND, OR, (, ).
type Where []interface{}

// WhereCondition describes a where clause
type WhereCondition struct {
	Column   string
	Operator string
	Value    interface{}
}

// WhereEqual convenience method to form a equals where clause
func WhereEqual(column string, value interface{}) WhereCondition {
	return WhereCondition{
		Column:   column,
		Operator: "=",
		Value:    value,
	}
}

// WhereNotEqual convenience method to form a not equal where clause
func WhereNotEqual(column string, value interface{}) WhereCondition {
	return WhereCondition{
		Column:   column,
		Operator: "!=",
		Value:    value,
	}
}

// WhereGreaterThan convenience method to form a greater than where clause
func WhereGreaterThan(column string, value interface{}) WhereCondition {
	return WhereCondition{
		Column:   column,
		Operator: ">",
		Value:    value,
	}
}

// WhereLessThan convenience method to form a less than where clause
func WhereLessThan(column string, value interface{}) WhereCondition {
	return WhereCondition{
		Column:   column,
		Operator: "<",
		Value:    value,
	}
}

// WhereIn convenience method to form a in where clause
func WhereIn(column string, values []interface{}) WhereCondition {
	var valueStrings = make([]string, len(values))
	for i, value := range values {
		valueStrings[i] = sanitizeValue(value)
	}
	return WhereCondition{
		Column:   column,
		Operator: "IN",
		Value:    "(" + strings.Join(valueStrings, ",") + ")",
	}
}

func (w WhereCondition) sql() string {
	sql := "`" + stripName(w.Column) + "` " + w.Operator + " " + sanitizeValue(w.Value)
	return sql
}

func (q Where) sql() string {
	sql := ""
	for _, w := range q {
		typeOf := reflect.TypeOf(w).String()
		if typeOf == "string" {
			sql += " " + w.(string) + " "
		} else if typeOf == "sqli.WhereCondition" {
			sql += " " + w.(WhereCondition).sql() + " "
		}
	}
	return sql
}
