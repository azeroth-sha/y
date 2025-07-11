package ygrace

import (
	"time"

	"github.com/azeroth-sha/y/logger"
)

type option struct {
	dur time.Duration
	log logger.Logger
}

type Option func(*option)

func WithDuration(d time.Duration) Option {
	return func(o *option) {
		o.dur = d
	}
}

func WithLogger(l logger.Logger) Option {
	return func(o *option) {
		if l == nil {
			l = logger.DefaultLog()
		}
		o.log = l
	}
}
