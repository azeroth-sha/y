package yrand

import (
	"bytes"
	"crypto/rand"
	"io"
	"sync"

	"github.com/azeroth-sha/y/ybuff"
	"github.com/azeroth-sha/y/yconst"
)

const (
	Numeral  = "0123456789"
	Lower    = "abcdefghijklmnopqrstuvwxyz"
	Upper    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Alpha    = Lower + Upper
	AlphaNum = Numeral + Alpha
	Symbol   = "!@#$%^&*()_+-=[]{}|;':\",./<>?"
	AllChars = AlphaNum + Symbol
)

func Uint8() uint8 {
	buf := getBuf(1)
	defer ybuff.Put(buf)
	n, _ := buf.ReadByte()
	return n
}

func Uint16() uint16 {
	buf := getBuf(2)
	defer ybuff.Put(buf)
	b1, _ := buf.ReadByte()
	b2, _ := buf.ReadByte()
	return uint16(b1) | uint16(b2)<<8
}

func Uint32() uint32 {
	buf := getBuf(4)
	defer ybuff.Put(buf)
	n1, _ := buf.ReadByte()
	n2, _ := buf.ReadByte()
	n3, _ := buf.ReadByte()
	n4, _ := buf.ReadByte()
	return uint32(n1) | uint32(n2)<<8 | uint32(n3)<<16 | uint32(n4)<<24
}

func Uint64() uint64 {
	buf := getBuf(8)
	defer ybuff.Put(buf)
	b1, _ := buf.ReadByte()
	b2, _ := buf.ReadByte()
	b3, _ := buf.ReadByte()
	b4, _ := buf.ReadByte()
	b5, _ := buf.ReadByte()
	b6, _ := buf.ReadByte()
	b7, _ := buf.ReadByte()
	b8, _ := buf.ReadByte()
	return uint64(b1) | uint64(b2)<<8 | uint64(b3)<<16 | uint64(b4)<<24 |
		uint64(b5)<<32 | uint64(b6)<<40 | uint64(b7)<<48 | uint64(b8)<<56
}
func Uint() uint {
	buf := getBuf(yconst.IntCap)
	defer ybuff.Put(buf)
	var (
		b byte
		n uint
	)
	for i := 0; i < yconst.IntCap; i++ {
		b, _ = buf.ReadByte()
		n |= uint(b) << uint(yconst.IntCap-i-1) * 8
	}
	return n
}

func Int8() int8 {
	if num := int8(Uint8()); num < 0 {
		return -num
	} else {
		return num
	}
}

func Int16() int16 {
	if num := int16(Uint16()); num < 0 {
		return -num
	} else {
		return num
	}
}

func Int32() int32 {
	if num := int32(Uint32()); num < 0 {
		return -num
	} else {
		return num
	}
}

func Int64() int64 {
	if num := int64(Uint64()); num < 0 {
		return -num
	} else {
		return num
	}
}

func Int() int {
	if num := int(Uint()); num < 0 {
		return -num
	} else {
		return num
	}
}

func Bytes(n int) []byte {
	buf := getBuf(n)
	defer ybuff.Put(buf)
	bts := make([]byte, n)
	_, _ = buf.Read(bts)
	return bts
}

func CharsBy(n int, dict []byte) []byte {
	dictLen := len(dict)
	bts := make([]byte, n)
	for i := 0; i < n; i++ {
		bts[i] = dict[Int()%dictLen]
	}
	return bts
}

func StringBy(n int, dict string) string {
	return string(CharsBy(n, []byte(dict)))
}

/*
  Package private
*/

var (
	reader   = rand.Reader
	readerMu = new(sync.Mutex)
	capLen   = 4 << 10
	pool     = sync.Pool{New: newBuf}
)

func getBuf(n int) *bytes.Buffer {
	buf := ybuff.Get()
	for cnt := buf.Len(); buf.Len() < n; cnt = buf.Len() {
		readBuf(buf, n-cnt)
	}
	return buf
}

func readBuf(w *bytes.Buffer, n int) {
	r := pool.Get().(*bytes.Buffer)
	if _, _ = io.CopyN(w, r, int64(n)); r.Len() > 0 {
		pool.Put(r)
	} else {
		ybuff.Put(r)
	}
}

func newBuf() any {
	buf := ybuff.Get()
	read(buf)
	return buf
}

func read(w io.Writer) {
	readerMu.Lock()
	defer readerMu.Unlock()
	_, _ = io.CopyN(w, reader, int64(capLen))
}
