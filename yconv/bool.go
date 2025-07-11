package yconv

import (
	"fmt"
	"reflect"
	"strconv"
)

// Bool converts v to bool.
func Bool(v any) (bool, error) {
	return toBool(v)
}

// MustBool converts v to bool and lose error.
func MustBool(v any) bool {
	vv, _ := toBool(v)
	return vv
}

/*
  Package private
*/

func toBool(v any) (bool, error) {
	var val bool
	switch vv := v.(type) {
	case string:
		return parseBool(vv)
	case bool:
		val = vv
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		val = vv == 1
	default:
		return refToBool(v)
	}
	return val, nil
}

func refToBool(v any) (bool, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		return toBool(rv.Elem().Interface())
	}
	switch rv.Kind() {
	case reflect.Bool:
		return rv.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 1, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() == 1, nil
	case reflect.String:
		return parseBool(rv.String())
	default:
		return false, fmt.Errorf(`unsupported type: %T`, v)
	}
}

func parseBool(str string) (bool, error) {
	return strconv.ParseBool(str)
}
