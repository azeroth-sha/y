package ylock

import (
	"sync"
)

type pMutex struct {
	once sync.Once
	mu   sync.Mutex
	size int
	dict map[string][2]chan *sync.Mutex
}

func (h *pMutex) init() {
	if h.size <= 0 {
		h.size = defaultSize
	}
	h.dict = make(map[string][2]chan *sync.Mutex)
}

func (h *pMutex) get(s string, l, b bool) *sync.Mutex {
	h.mu.Lock()
	h.once.Do(h.init)
	if _, y := h.dict[s]; !y {
		h.dict[s] = [2]chan *sync.Mutex{
			make(chan *sync.Mutex, h.size),
			make(chan *sync.Mutex, h.size),
		}
		for i := 0; i < h.size; i++ {
			h.dict[s][0] <- new(sync.Mutex)
		}
	}
	h.mu.Unlock()
	if l {
		if b {
			return <-h.dict[s][0]
		}
		select {
		case mu := <-h.dict[s][0]:
			return mu
		default:
			return nil
		}
	} else {
		return <-h.dict[s][1]
	}
}

func (h *pMutex) TryLock(s string) bool {
	mu := h.get(s, true, false)
	if mu == nil {
		return false
	}
	f := mu.TryLock()
	h.dict[s][1] <- mu
	return f
}

func (h *pMutex) Lock(s string) {
	mu := h.get(s, true, true)
	mu.Lock()
	h.dict[s][1] <- mu
}

func (h *pMutex) Unlock(s string) {
	mu := h.get(s, false, true)
	mu.Unlock()
	h.dict[s][0] <- mu
}

// NewPoolMutex new pool locker
func NewPoolMutex(n ...int) Locker {
	l := new(pMutex)
	if len(n) > 0 && n[0] > 0 {
		l.size = n[0]
	}
	return l
}
