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
	case "*string":
		s, ok := value.(*string)
		if !ok {
			return ""
		}
		if s == nil {
			return ""
		}
		return sanitizeValue(*s)
	case "int":
		return fmt.Sprintf("%d", value.(int))
	case "*int":
		s, ok := value.(*int)
		if !ok {
			return ""
		}
		if s == nil {
			return ""
		}
		return sanitizeValue(*s)
	case "int16":
		return fmt.Sprintf("%d", value.(int16))
	case "*int16":
		s, ok := value.(*int16)
		if !ok {
			return ""
		}
		if s == nil {
			return ""
		}
		return sanitizeValue(*s)
	case "int32":
		return fmt.Sprintf("%d", value.(int32))
	case "*int32":
		s, ok := value.(*int32)
		if !ok {
			return ""
		}
		if s == nil {
			return ""
		}
		return sanitizeValue(*s)
	case "int64":
		return fmt.Sprintf("%d", value.(int64))
	case "*int64":
		s, ok := value.(*int64)
		if !ok {
			return ""
		}
		if s == nil {
			return ""
		}
		return sanitizeValue(*s)
	case "uint":
		return fmt.Sprintf("%d", value.(uint))
	case "*uint":
		s, ok := value.(*uint)
		if !ok {
			return ""
		}
		if s == nil {
			return ""
		}
		return sanitizeValue(*s)
	case "uint16":
		return fmt.Sprintf("%d", value.(uint16))
	case "*uint16":
		s, ok := value.(*uint16)
		if !ok {
			return ""
		}
		if s == nil {
			return ""
		}
		return sanitizeValue(*s)
	case "uint32":
		return fmt.Sprintf("%d", value.(uint32))
	case "*uint32":
		s, ok := value.(*uint32)
		if !ok {
			return ""
		}
		if s == nil {
			return ""
		}
		return sanitizeValue(*s)
	case "uint64":
		return fmt.Sprintf("%d", value.(uint64))
	case "*uint64":
		s, ok := value.(*uint64)
		if !ok {
			return ""
		}
		if s == nil {
			return ""
		}
		return sanitizeValue(*s)
	case "float32":
	case "float64":
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case "[]uint8":
		return fmt.Sprintf("'%x'", value.([]byte))
	case "*[]uint8":
		s, ok := value.(*[]uint8)
		if !ok {
			return ""
		}
		if s == nil {
			return ""
		}
		return sanitizeValue(*s)
	case "bool":
		if value.(bool) {
			return "1"
		}
		return "0"
	}

	return ""
}
