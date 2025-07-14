package ygrace

import (
	"time"

	"github.com/azeroth-sha/y/ylog"
)

type option struct {
	dur time.Duration
	log ylog.Logger
}

type Option func(*option)

func WithDuration(d time.Duration) Option {
	return func(o *option) {
		o.dur = d
	}
}

func WithLogger(l ylog.Logger) Option {
	return func(o *option) {
		if l == nil {
			l = ylog.DefaultLog()
		}
		o.log = l
	}
}
