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
		if o.dur = d; o.dur < 0 {
			o.dur = 0
		}
	}
}

func WithLogger(l ylog.Logger) Option {
	return func(o *option) {
		if o.log = l; o.log == nil {
			o.log = ylog.NewLogger()
		}
	}
}
