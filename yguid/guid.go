package yguid

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/big"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/azeroth-sha/y/internal"
	"github.com/azeroth-sha/y/yconv"
	"github.com/azeroth-sha/y/yrand"
	"github.com/azeroth-sha/y/ysum"
)

const (
	BLen = 12 // GUID字节长度
	SLen = 20 // GUID字符长度
	Base = 36 // GUID转换进制
)

var (
	ErrInvalid = errors.New("invalid GUID")
)

type GUID [BLen]byte

func (g GUID) Time() time.Time {
	unix := uint32(0)
	_ = binary.Read(bytes.NewReader(g[:4]), endian, &unix)
	return time.Unix(yconv.MustDigit[int64](unix), 0)
}

func (g GUID) Mark() uint32 {
	mrk := uint32(0)
	_ = binary.Read(bytes.NewReader(g[4:8]), endian, &mrk)
	return mrk
}

func (g GUID) Index() uint16 {
	idx := uint16(0)
	_ = binary.Read(bytes.NewReader(g[8:10]), endian, &idx)
	return idx
}

func (g GUID) Rand() uint16 {
	num := uint16(0)
	_ = binary.Read(bytes.NewReader(g[10:12]), endian, &num)
	return num
}

func (g GUID) String() string {
	bInt := getInt()
	defer putInt(bInt)
	id := bInt.SetBytes(g[:]).Text(Base)
	if len(id) < SLen {
		id = strings.Repeat("0", SLen-len(id)) + id
	}
	return id
}

func (g GUID) Bytes() []byte {
	return g[:]
}

// String GUID to string
func String() string {
	return New().String()
}

// New GUID
func New() GUID {
	return NewWith(uint32(time.Now().Unix()), hostMark)
}

// NewWith GUID with unix time and mark
func NewWith(unixSec, mark uint32) GUID {
	var id GUID
	endian.PutUint32(id[:4], unixSec)
	endian.PutUint32(id[4:8], mark)
	endian.PutUint16(id[8:10], getSerial())
	endian.PutUint16(id[10:12], yrand.Uint16())
	return id
}

// Parse GUID from string
func Parse(s string) (id GUID, _ error) {
	if len(s) != SLen {
		return id, ErrInvalid
	}
	bInt := getInt()
	defer putInt(bInt)
	if _, ok := bInt.SetString(s, Base); !ok {
		return id, ErrInvalid
	} else {
		_ = copy(id[:], bInt.Bytes())
		return id, nil
	}
}

// MustParse parse GUID from string and lose error
func MustParse(s string) GUID {
	id, _ := Parse(s)
	return id
}

/*
  Package private
*/

var (
	endian   = binary.BigEndian
	bigPool  = &sync.Pool{New: func() any { return new(big.Int) }}
	hostMark = uint32(0)
	serial   = uint32(0)
)

func init() {
	hostMark = uint32(getHostID())<<16 | uint32(os.Getpid())
	serial = yrand.Uint32()
}

func getInt() *big.Int {
	return bigPool.Get().(*big.Int)
}

func putInt(bInt *big.Int) {
	bInt.SetInt64(0)
	bigPool.Put(bInt)
}

func getHostID() uint16 {
	hid, err := internal.HostID()
	if hid == "" || err != nil {
		hid = yrand.StringBy(16, yrand.AlphaNum)
	}
	h := ysum.NewCrc16()
	_, _ = h.Write([]byte(hid))
	return h.Sum16()
}

func getSerial() uint16 {
	n := atomic.AddUint32(&serial, 1)
	return uint16(n)
}
