package ycache

import (
	"time"
)

type item struct {
	val  any
	exp  time.Time
	hand ExpireFun
}

func (i *item) IsExpired() bool {
	return !i.exp.IsZero() && i.exp.Before(time.Now())
}

func newItem(cnf *options, _ string, val any, opts ...ItemOption) *item {
	i := &item{val: val}
	if cnf.expire != nil {
		i.hand = cnf.expire
	}
	for _, opt := range opts {
		opt(i)
	}
	return i
}
