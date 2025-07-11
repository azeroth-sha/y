package yconv

import (
	"fmt"
	"reflect"
	"strconv"
)

// String convert any to string
func String(v any) (string, error) {
	return toString(v)
}

// MustString convert any to string and lose error
func MustString(v any) string {
	vv, _ := toString(v)
	return vv
}

/*
  Package private
*/

func toString(v any) (string, error) {
	var val string
	switch vv := v.(type) {
	case bool:
		val = strconv.FormatBool(vv)
	case string:
		val = vv
	case []byte:
		val = string(vv)
	case []rune:
		val = string(vv)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		val = fmt.Sprintf("%d", vv)
	case float32, float64:
		val = fmt.Sprintf("%g", vv)
	default:
		return refToString(v)
	}
	return val, nil
}

func refToString(v any) (string, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		return toString(rv.Elem().Interface())
	}
	switch rv.Kind() {
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return toString(rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return toString(rv.Uint())
	case reflect.Float32, reflect.Float64:
		return toString(rv.Float())
	case reflect.String:
		return rv.String(), nil
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return toString(rv.Bytes())
		}
		fallthrough
	default:
		return ``, fmt.Errorf(`unsupported type: %T`, v)
	}
}
