package ylock

import (
	"runtime"
)

var cnt int

func init() {
	cnt = runtime.NumCPU() * 2
}

type Locker interface {
	TryLock(s string) bool
	Lock(s string)
	Unlock(s string)
}
