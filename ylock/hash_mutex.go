package ylock

import (
	"hash/fnv"
	"sync"

	"github.com/azeroth-sha/y/internal"
	"github.com/azeroth-sha/y/ysum"
)

type hMutex struct {
	once sync.Once
	mu   sync.Mutex
	size uint32
	dict map[uint32]*sync.Mutex
}

func (h *hMutex) init() {
	if h.size <= 0 {
		h.size = uint32(cnt)
	}
	h.dict = make(map[uint32]*sync.Mutex)
}

func (h *hMutex) get(s string) *sync.Mutex {
	h.mu.Lock()
	h.once.Do(h.init)
	sum := ysum.Sum32(fnv.New32a(), internal.ToBytes(s)) % h.size
	if _, y := h.dict[sum]; !y {
		h.dict[sum] = new(sync.Mutex)
	}
	h.mu.Unlock()
	return h.dict[sum]
}

func (h *hMutex) TryLock(s string) bool {
	mu := h.get(s)
	return mu.TryLock()
}

func (h *hMutex) Lock(s string) {
	mu := h.get(s)
	mu.Lock()
}

func (h *hMutex) Unlock(s string) {
	mu := h.get(s)
	mu.Unlock()
}

// NewHashMutex new hash locker
func NewHashMutex(n ...int) Locker {
	l := new(hMutex)
	if len(n) > 0 {
		l.size = uint32(n[0])
	}
	return l
}
