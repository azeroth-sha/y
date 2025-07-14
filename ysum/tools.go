package ysum

import (
	"encoding/hex"
	"hash"
	"io"
	"os"
)

// Sum returns the checksum of data.
func Sum(h hash.Hash, data []byte) []byte {
	_, _ = h.Write(data)
	return h.Sum(nil)
}

// SumHex returns the checksum of data as a hex string.
func SumHex(h hash.Hash, data []byte) string {
	return hex.EncodeToString(Sum(h, data))
}

// SumFrom returns the checksum of the data read from r.
func SumFrom(h hash.Hash, r io.Reader) []byte {
	_, _ = io.Copy(h, r)
	return h.Sum(nil)
}

// SumFromHex returns the checksum of the data read from r as a hex string.
func SumFromHex(h hash.Hash, r io.Reader) string {
	return hex.EncodeToString(SumFrom(h, r))
}

// SumFile returns the checksum of the file n.
func SumFile(h hash.Hash, n string) ([]byte, error) {
	f, e := os.Open(n)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	return SumFrom(h, f), nil
}

// SumFileHex returns the checksum of the file n as a hex string.
func SumFileHex(h hash.Hash, n string) (string, error) {
	f, e := os.Open(n)
	if e != nil {
		return ``, e
	}
	defer f.Close()
	return SumFromHex(h, f), nil
}

// Sum16 returns the checksum of data.
func Sum16(h Hash16, data []byte) uint16 {
	_, _ = h.Write(data)
	return h.Sum16()
}

// SumFrom16 returns the checksum of the data read from r.
func SumFrom16(h Hash16, r io.Reader) uint16 {
	_, _ = io.Copy(h, r)
	return h.Sum16()
}

// SumFile16 returns the checksum of the file n.
func SumFile16(h Hash16, n string) (uint16, error) {
	f, e := os.Open(n)
	if e != nil {
		return 0, e
	}
	defer f.Close()
	return SumFrom16(h, f), nil
}

// Sum32 returns the checksum of data.
func Sum32(h hash.Hash32, data []byte) uint32 {
	_, _ = h.Write(data)
	return h.Sum32()
}

// SumFrom32 returns the checksum of the data read from r.
func SumFrom32(h hash.Hash32, r io.Reader) uint32 {
	_, _ = io.Copy(h, r)
	return h.Sum32()
}

// SumFile32 returns the checksum of the file n.
func SumFile32(h hash.Hash32, n string) (uint32, error) {
	f, e := os.Open(n)
	if e != nil {
		return 0, e
	}
	defer f.Close()
	return SumFrom32(h, f), nil
}

// Sum64 returns the checksum of data.
func Sum64(h hash.Hash64, data []byte) uint64 {
	_, _ = h.Write(data)
	return h.Sum64()
}

// SumFrom64 returns the checksum of the data read from r.
func SumFrom64(h hash.Hash64, r io.Reader) uint64 {
	_, _ = io.Copy(h, r)
	return h.Sum64()
}

// SumFile64 returns the checksum of the file n.
func SumFile64(h hash.Hash64, n string) (uint64, error) {
	f, e := os.Open(n)
	if e != nil {
		return 0, e
	}
	defer f.Close()
	return SumFrom64(h, f), nil
}
