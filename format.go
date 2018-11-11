package sqli

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func stripName(name string) string {
	return strings.Replace(name, "`", "``", -1)
}

func sanitizeValue(value interface{}) string {
	typeOf := reflect.TypeOf(value).String()

	switch typeOf {
	case "string":
		return "'" + strings.Replace(value.(string), "'", "''", -1) + "'"
	case "int":
		return fmt.Sprintf("%d", value.(int))
	case "int16":
		return fmt.Sprintf("%d", value.(int16))
	case "int32":
		return fmt.Sprintf("%d", value.(int32))
	case "int64":
		return fmt.Sprintf("%d", value.(int64))
	case "uint":
		return fmt.Sprintf("%d", value.(uint))
	case "uint16":
		return fmt.Sprintf("%d", value.(uint16))
	case "uint32":
		return fmt.Sprintf("%d", value.(uint32))
	case "uint64":
		return fmt.Sprintf("%d", value.(uint64))
	case "float32":
	case "float64":
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case "[]uint8":
		return fmt.Sprintf("'%x'", value.([]byte))
	case "bool":
		if value.(bool) {
			return "1"
		}
		return "0"
	}

	return fmt.Sprintf("'%s'", value)
}
