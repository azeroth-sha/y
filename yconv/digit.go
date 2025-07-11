package yconv

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/azeroth-sha/y/internal"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Digit convert any to digit.
func Digit[T Number](v any) (T, error) {
	val, err := toDigit[T](v)
	return val, err
}

// MustDigit convert any to digit and lose error.
func MustDigit[T Number](v any) T {
	val, _ := toDigit[T](v)
	return val
}

/*
  Package private
*/

func toDigit[T Number](v any) (T, error) {
	var val T
	switch vv := v.(type) {
	case string:
		return parseDigit[T](vv)
	case []byte:
		return parseDigit[T](string(vv))
	case bool:
		val = internal.Or[T](vv, 1, 0)
	case int:
		val = T(vv)
	case int8:
		val = T(vv)
	case int16:
		val = T(vv)
	case int32:
		val = T(vv)
	case int64:
		val = T(vv)
	case uint:
		val = T(vv)
	case uint8:
		val = T(vv)
	case uint16:
		val = T(vv)
	case uint32:
		val = T(vv)
	case uint64:
		val = T(vv)
	default:
		return refToDigit[T](v)
	}
	return val, nil
}

func refToDigit[T Number](v any) (T, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		return toDigit[T](rv.Elem().Interface())
	}
	switch rv.Kind() {
	case reflect.Bool:
		return internal.Or[T](rv.Bool(), 1, 0), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return toDigit[T](rv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return toDigit[T](rv.Uint())
	case reflect.String:
		return parseDigit[T](rv.String())
	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return toDigit[T](rv.Bytes())
		}
		fallthrough
	default:
		return 0, fmt.Errorf(`unsupported type: %T`, v)
	}
}

func parseDigit[T Number](str string) (T, error) {
	num, err := strconv.ParseInt(str, 0, 0)
	return T(num), err
}
