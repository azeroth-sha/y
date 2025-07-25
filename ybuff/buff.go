package ybuff

import (
	"bytes"
	"sync"
)

// Get a buffer from pool
func Get() *bytes.Buffer {
	return pool.Get().(*bytes.Buffer)
}

// Put a buffer to pool
func Put(buf ...*bytes.Buffer) {
	for _, b := range buf {
		if b == nil {
			continue
		}
		b.Reset()
		pool.Put(b)
	}
}

/*
  Package private
*/

var pool = &sync.Pool{New: newBuf}

func newBuf() any {
	return bytes.NewBuffer(make([]byte, 0, 512))
}
