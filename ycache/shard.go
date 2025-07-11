package ycache

import (
	"sync"
	"time"
)

type shard struct {
	mu   *sync.RWMutex
	dict map[string]*item
}

func (s *shard) Has(cnf *options, key string) (has bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if i, y := s.dict[key]; y {
		if i.IsExpired() {
			s.expired(cnf, key)
			return false
		}
		return true
	}
	return false
}

func (s *shard) Set(cnf *options, key string, val any, opts ...ItemOption) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i, y := s.dict[key]; y && i.IsExpired() {
		s.expired(cnf, key)
	}
	s.dict[key] = newItem(cnf, key, val, opts...)
}

func (s *shard) Get(cnf *options, key string) (val any, has bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if i, y := s.dict[key]; y {
		if i.IsExpired() {
			s.expired(cnf, key)
			return nil, false
		}
		return i.val, true
	}
	return nil, false
}

func (s *shard) Del(cnf *options, key string) (has bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i, y := s.dict[key]; y {
		if i.IsExpired() {
			s.expired(cnf, key)
			return false
		} else {
			delete(s.dict, key)
			return true
		}
	}
	return false
}

func (s *shard) GetDel(cnf *options, key string) (val any, has bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i, y := s.dict[key]; y {
		if i.IsExpired() {
			s.expired(cnf, key)
			return nil, false
		} else {
			delete(s.dict, key)
			return i.val, true
		}
	}
	return nil, false
}

func (s *shard) GetSet(cnf *options, key string, newVal any, opts ...ItemOption) (oldVal any, has bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i, y := s.dict[key]; y {
		if i.IsExpired() {
			s.expired(cnf, key)
			return nil, false
		} else {
			oldVal = i.val
			s.dict[key] = newItem(cnf, key, newVal, opts...)
			return oldVal, true
		}
	} else {
		s.dict[key] = newItem(cnf, key, newVal, opts...)
		return nil, false
	}
}

func (s *shard) SetX(cnf *options, key string, val any, opts ...ItemOption) (ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i, y := s.dict[key]; y && i.IsExpired() {
		s.expired(cnf, key)
	}
	if _, y := s.dict[key]; y {
		return false
	} else {
		s.dict[key] = newItem(cnf, key, val, opts...)
		return true
	}
}

func (s *shard) DelExpired(cnf *options, key string) (has, exp bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i, y := s.dict[key]; y {
		has = true
		if exp = i.IsExpired(); exp {
			s.expired(cnf, key)
		}
	}
	return has, exp
}

func (s *shard) All(cnf *options) (kvs map[string]any) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	kvs = make(map[string]any)
	for k, v := range s.dict {
		if v.IsExpired() {
			s.expired(cnf, k)
		} else {
			kvs[k] = v.val
		}
	}
	return kvs
}

func (s *shard) Len(_ *options) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.dict)
}

func (s *shard) Clear(cnf *options) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	count := 0
	for k, v := range s.dict {
		if v.IsExpired() {
			s.expired(cnf, k)
		} else {
			delete(s.dict, k)
			count++
		}
	}
	return count
}

func (s *shard) TTL(cnf *options, key string) (ttl time.Duration, has bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if i, y := s.dict[key]; y {
		if i.IsExpired() {
			s.expired(cnf, key)
			return 0, false
		}
		return i.exp.Sub(time.Now()), true
	}
	return 0, false
}

func (s *shard) Expire(cnf *options, key string, ttl time.Duration) (ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if i, y := s.dict[key]; y {
		if i.IsExpired() {
			s.expired(cnf, key)
			return false
		}
		i.exp = time.Now().Add(ttl)
		return true
	}
	return false
}

func (s *shard) Check(cnf *options) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range s.dict {
		if v.IsExpired() {
			s.expired(cnf, k)
		}
	}
}

/*
  Package private
*/

func (s *shard) expired(cnf *options, key string) {
	i := s.dict[key]
	delete(s.dict, key)
	defer func() {
		if rec := recover(); rec != nil {
			cnf.log.Errorf("ycache: expired panic: %v", rec)
		}
	}()
	if i.hand != nil {
		i.hand(key, i.val)
	}
}

func newShard() *shard {
	return &shard{
		mu:   new(sync.RWMutex),
		dict: make(map[string]*item),
	}
}
