package ystudio

import (
	"time"

	"github.com/azeroth-sha/y/ylog"
)

// AppendOption represents the configuration options for the Append operation.
type AppendOption func(*appendConfig)

type appendConfig struct {
	waitDur time.Duration
	logger  ylog.Logger
}

// WithWaitDuration sets the wait duration for the Append operation.
func WithWaitDuration(d time.Duration) AppendOption {
	return func(c *appendConfig) {
		if d < 0 {
			d = -1
		}
		c.waitDur = d
	}
}

// WithLogger sets the logger for the Append operation.
func WithLogger(logger ylog.Logger) AppendOption {
	return func(c *appendConfig) {
		if logger == nil {
			logger = ylog.DefaultLog()
		}
		c.logger = logger
	}
}

// DefaultOptions returns the default options for the Append operation.
func DefaultOptions() []AppendOption {
	return []AppendOption{
		WithWaitDuration(-1),
		WithLogger(ylog.DefaultLog()),
	}
}
