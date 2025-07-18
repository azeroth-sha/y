package ytime

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/azeroth-sha/y/yconst"
)

// Format convert any to string.
func Format(v any, f ...string) string {
	t := MustToTime(v)
	if len(f) > 0 {
		return t.Format(f[0])
	}
	return t.Format(RFC3339)
}

// ToTime convert any to time.Time.
func ToTime(v any) (t time.Time, err error) {
	switch vv := v.(type) {
	case sql.NullTime:
		if vv.Valid {
			t = vv.Time
		}
	case *sql.NullTime:
		if vv != nil && vv.Valid {
			t = vv.Time
		}
	case *time.Time:
		if vv != nil {
			t = *vv
		}
	case time.Time:
		t = vv
	case int32:
		t = time.Unix(int64(vv), 0)
	case uint32:
		t = time.Unix(int64(vv), 0)
	case int64:
		sec := vv
		ms := int64(0)
		if sec > int64(yconst.MaxUint32) {
			ms = sec % 1000
			sec /= 1000
		}
		t = time.Unix(sec, ms*int64(time.Millisecond))
	case uint64:
		sec := int64(vv)
		ms := int64(0)
		if sec > int64(yconst.MaxUint32) {
			ms = sec % 1000
			sec /= 1000
		}
		t = time.Unix(sec, ms*int64(time.Millisecond))
	case string:
		return ParseTime(vv)
	}
	return t.In(time.Local), nil
}

// MustToTime convert any to time.Time and lose error.
func MustToTime(v any) time.Time {
	t, _ := ToTime(v)
	return t
}

// Parse convert string to time.Time
func Parse(layout, str string) (time.Time, error) {
	return time.Parse(layout, str)
}

// MustParse convert string to time.Time and lose error.
func MustParse(layout, str string) time.Time {
	t, _ := Parse(layout, str)
	return t
}

// ParseTime convert string to time.Time
func ParseTime(str string, loc ...*time.Location) (time.Time, error) {
	l := time.Local
	if len(loc) > 0 {
		l = loc[0]
	}
	for _, f := range layouts {
		if f.z {
			if t, err := time.Parse(f.f, str); err == nil {
				return t.In(l), nil
			}
		} else if t, err := time.ParseInLocation(f.f, str, l); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time format: %s", str)
}

// MustParseTime convert string to time.Time and lose error.
func MustParseTime(str string, loc ...*time.Location) time.Time {
	t, _ := ParseTime(str, loc...)
	return t
}

// Unix convert any to unix timestamp .
func Unix(v any) (int64, error) {
	if t, e := ToTime(v); e == nil {
		return t.Unix(), nil
	} else {
		return 0, e
	}
}

// MustUnix convert any to unix timestamp and lose error.
func MustUnix(v any) int64 {
	sec, _ := Unix(v)
	return sec
}

// UnixMilli convert any to unix milli timestamp .
func UnixMilli(v any) (int64, error) {
	if t, e := ToTime(v); e == nil {
		return t.UnixMilli(), nil
	} else {
		return 0, e
	}
}

// MustUnixMilli convert any to unix milli timestamp and lose error.
func MustUnixMilli(v any) int64 {
	ms, _ := UnixMilli(v)
	return ms
}

// UnixMicro convert any to unix micro timestamp .
func UnixMicro(v any) (int64, error) {
	if t, e := ToTime(v); e == nil {
		return t.UnixMicro(), nil
	} else {
		return 0, e
	}
}

// MustUnixMicro convert any to unix micro timestamp and lose error.
func MustUnixMicro(v any) int64 {
	ns, _ := UnixMicro(v)
	return ns
}

// UnixNano convert any to unix nano timestamp .
func UnixNano(v any) (int64, error) {
	if t, e := ToTime(v); e == nil {
		return t.UnixNano(), nil
	} else {
		return 0, e
	}
}

// MustUnixNano convert any to unix nano timestamp and lose error.
func MustUnixNano(v any) int64 {
	ns, _ := UnixNano(v)
	return ns
}

/*
Package private
*/
type format struct {
	f string // format string
	z bool   // has zone info
}

var layouts = []*format{
	{Stamp, false},
	{StampMilli, false},
	{StampMicro, false},
	{StampNano, false},
	{DateTime, false},
	{DateOnly, false},
	{TimeOnly, false},
	{Standard, true},
	{Layout, true},
	{ANSIC, false},
	{UnixDate, false},
	{RubyDate, true},
	{RFC822, false},
	{RFC822Z, true},
	{RFC850, false},
	{RFC1123, false},
	{RFC1123Z, true},
	{RFC3339, true},
	{RFC3339Nano, true},
	{Kitchen, false},
}
