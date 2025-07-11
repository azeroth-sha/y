package internal

import (
	"unsafe"
)

// ToString converts byte slice to string without a memory allocation.
//
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func ToString(bs []byte) string {
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}

// ToBytes converts string to byte slice without a memory allocation.
//
// For more details, see https://github.com/golang/go/issues/53003#issuecomment-1140276077.
func ToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
