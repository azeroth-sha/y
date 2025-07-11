package y

import (
	"reflect"

	"github.com/azeroth-sha/y/internal"
)

// Or returns y if cond is true, n otherwise.
func Or[T any](cond bool, y, n T) T {
	return internal.Or(cond, y, n)
}

// IsNil returns true if v is nil.
func IsNil(v any) bool {
	return v == nil || reflect.ValueOf(v).IsNil()
}

// IsEmpty returns true if v is empty.
func IsEmpty(v any) bool {
	return IsNil(v) || reflect.ValueOf(v).IsZero()
}
