package yconv

import (
	"github.com/azeroth-sha/y/internal"
)

// ToString converts []byte to string
func ToString(bs []byte) string {
	return internal.ToString(bs)
}

// ToBytes converts string to []byte
func ToBytes(s string) []byte {
	return internal.ToBytes(s)
}
