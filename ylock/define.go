package ylock

import (
	"runtime"
)

var defaultSize int

func init() {
	defaultSize = runtime.NumCPU() * 2
}

type Locker interface {
	TryLock(s string) bool
	Lock(s string)
	Unlock(s string)
}
