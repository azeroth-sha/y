package ylock

import (
	"sync"
)

type nMutex struct {
	once sync.Once
	mu   sync.Mutex
	dict map[string]*sync.Mutex
}

func (n *nMutex) init() {
	n.dict = make(map[string]*sync.Mutex)
}

func (n *nMutex) get(s string) *sync.Mutex {
	n.mu.Lock()
	n.once.Do(n.init)
	if _, y := n.dict[s]; !y {
		n.dict[s] = new(sync.Mutex)
	}
	n.mu.Unlock()
	return n.dict[s]
}

func (n *nMutex) TryLock(s string) bool {
	mu := n.get(s)
	return mu.TryLock()
}

func (n *nMutex) Lock(s string) {
	mu := n.get(s)
	mu.Lock()
}

func (n *nMutex) Unlock(s string) {
	mu := n.get(s)
	mu.Unlock()
}

// NewNameLocker new name locker
func NewNameLocker() Locker {
	l := new(nMutex)
	return l
}
