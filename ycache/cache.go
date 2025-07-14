package ycache

import (
	"hash/fnv"
	"runtime"
	"time"

	"github.com/azeroth-sha/y/internal"
	"github.com/azeroth-sha/y/logger"
)

// Cache interface
type Cache interface {
	// Has returns whether the key exists.
	Has(key string) (has bool)
	// Set key-value, if the key exists, it will be overwritten.
	Set(key string, val any, opts ...ItemOption)
	// Get returns the value of the key, if the key does not exist, it returns false.
	Get(key string) (val any, has bool)
	// Del deletes the key.
	Del(key string) (ok bool)
	// GetDel returns the value of the key and deletes the key.
	GetDel(key string) (val any, has bool)
	// GetSet returns the value of the key and sets the new value.
	GetSet(key string, newVal any, opts ...ItemOption) (oldVal any, has bool)
	// SetX sets the key-value, if the key exists, it will not be overwritten.
	SetX(key string, val any, opts ...ItemOption) (ok bool)
	// DelExpired deletes the expired key.
	DelExpired(key string) (has, exp bool)
	// All returns all key-value.
	All() (kvs map[string]any)
	// Len returns the number of key-value.
	Len() int
	// Clear deletes all key-value.
	Clear() int
	// TTL returns the remaining time of the key.
	TTL(key string) (ttl time.Duration, has bool)
	// Expire sets the remaining time of the key.
	Expire(key string, ttl time.Duration) (ok bool)
}

type cache struct {
	opts   *options
	bucket []*shard
	closed chan struct{}
}

func (c *cache) Has(key string) (has bool) {
	s := c.getShard(key)
	return s.Has(c.opts, key)
}

func (c *cache) Set(key string, val any, opts ...ItemOption) {
	s := c.getShard(key)
	s.Set(c.opts, key, val, opts...)
}

func (c *cache) Get(key string) (val any, has bool) {
	s := c.getShard(key)
	return s.Get(c.opts, key)
}

func (c *cache) Del(key string) (ok bool) {
	s := c.getShard(key)
	return s.Del(c.opts, key)
}

func (c *cache) GetDel(key string) (val any, has bool) {
	s := c.getShard(key)
	return s.GetDel(c.opts, key)
}

func (c *cache) GetSet(key string, newVal any, opts ...ItemOption) (oldVal any, has bool) {
	s := c.getShard(key)
	return s.GetSet(c.opts, key, newVal, opts...)
}

func (c *cache) SetX(key string, val any, opts ...ItemOption) (ok bool) {
	s := c.getShard(key)
	return s.SetX(c.opts, key, val, opts...)
}

func (c *cache) DelExpired(key string) (has, exp bool) {
	s := c.getShard(key)
	return s.DelExpired(c.opts, key)
}

func (c *cache) All() (kvs map[string]any) {
	kvs = make(map[string]any)
	for _, s := range c.bucket {
		shardKvs := s.All(c.opts)
		for k, v := range shardKvs {
			kvs[k] = v
		}
	}
	return kvs
}

func (c *cache) Len() int {
	count := 0
	for _, s := range c.bucket {
		count += s.Len(c.opts)
	}
	return count
}

func (c *cache) Clear() int {
	count := 0
	for _, s := range c.bucket {
		count += s.Clear(c.opts)
	}
	return count
}

func (c *cache) TTL(key string) (ttl time.Duration, has bool) {
	s := c.getShard(key)
	return s.TTL(c.opts, key)
}

func (c *cache) Expire(key string, ttl time.Duration) (ok bool) {
	s := c.getShard(key)
	return s.Expire(c.opts, key, ttl)
}

func New(opts ...Option) Cache {
	cnf := &options{
		interval:  DefaultInterval,
		expire:    nil,
		shardSize: runtime.NumCPU() * 4,
		log:       logger.DefaultLog(),
	}
	for _, opt := range opts {
		opt(cnf)
	}
	bucket := &cache{
		opts:   cnf,
		bucket: make([]*shard, 0, cnf.shardSize),
		closed: make(chan struct{}),
	}
	for i := 0; i < cnf.shardSize; i++ {
		bucket.bucket = append(bucket.bucket, newShard())
	}
	go bucket.ticker()
	runtime.SetFinalizer(bucket, func(c *cache) {
		close(c.closed)
	})
	return bucket
}

/*
  Package private
*/

func (c *cache) getShard(k string) *shard {
	h := fnv.New32a()
	_, _ = h.Write(internal.ToBytes(k))
	return c.bucket[int(h.Sum32())%c.opts.shardSize]
}

func (c *cache) ticker() {
	dur := c.opts.interval
	if dur <= 0 {
		return
	}
	tk := time.NewTicker(dur)
	defer tk.Stop()
EXIT:
	for {
		select {
		case <-c.closed:
			break EXIT
		case <-tk.C:
			c.checkAll()
		}
	}
}

func (c *cache) checkAll() {
	for _, s := range c.bucket {
		s.Check(c.opts)
	}
}
