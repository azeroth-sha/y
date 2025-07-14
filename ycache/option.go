package ycache

import (
	"runtime"
	"time"

	"github.com/azeroth-sha/y/ylog"
)

const DefaultInterval = time.Second

// ExpireFun item expire callback
type ExpireFun func(key string, val any)

type Option func(opts *options)

type options struct {
	interval  time.Duration
	expire    ExpireFun
	shardSize int
	log       ylog.Logger
}

// WithInterval set expire check interval
func WithInterval(d time.Duration) Option {
	return func(opts *options) {
		if d <= 0 {
			d = DefaultInterval
		}
		opts.interval = d
	}
}

// WithExpireFun set default expire callback
func WithExpireFun(h ExpireFun) Option {
	return func(opts *options) {
		opts.expire = h
	}
}

// WithShardSize set shard size
func WithShardSize(count int) Option {
	return func(opts *options) {
		if count <= 0 {
			count = runtime.NumCPU() * 4
		}
		opts.shardSize = count
	}
}

// WithLogger set logger
func WithLogger(l ylog.Logger) Option {
	return func(opts *options) {
		if l == nil {
			l = ylog.DefaultLog()
		}
		opts.log = l
	}
}

// ItemOption item option
type ItemOption func(i *item)

// WithItemFun set item expire callback
func WithItemFun(h ExpireFun) ItemOption {
	return func(i *item) {
		i.hand = h
	}
}

// WithItemTime set item expire time
func WithItemTime(t time.Time) ItemOption {
	return func(i *item) {
		i.exp = t
	}
}

// WithItemDur set item expire duration
func WithItemDur(d time.Duration) ItemOption {
	return func(i *item) {
		i.exp = time.Now().Add(d)
	}
}
