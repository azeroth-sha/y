package yconst

const (
	IntSize = 32 << (^uint(0) >> 63) // 32 or 64
	IntCap  = IntSize / 8

	MaxInt   int   = 1<<(IntSize-1) - 1
	MaxInt8  int8  = 1<<7 - 1
	MaxInt16 int16 = 1<<15 - 1
	MaxInt32 int32 = 1<<31 - 1
	MaxInt64 int64 = 1<<63 - 1

	MinInt   int   = -MaxInt - 1
	MinInt8  int8  = -MaxInt8 - 1
	MinInt16 int16 = -MaxInt16 - 1
	MinInt32 int32 = -MaxInt32 - 1
	MinInt64 int64 = -MaxInt64 - 1

	MaxUint   uint   = 1<<IntSize - 1
	MaxUint8  uint8  = 1<<8 - 1
	MaxUint16 uint16 = 1<<16 - 1
	MaxUint32 uint32 = 1<<32 - 1
	MaxUint64 uint64 = 1<<64 - 1

	MaxFloat32 = 0x1p127 * (1 + (1 - 0x1p-23))  // 3.40282346638528859811704183484516925440e+38
	MaxFloat64 = 0x1p1023 * (1 + (1 - 0x1p-52)) // 1.79769313486231570814527423731704356798070e+308
)
