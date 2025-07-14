package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
}

var entry Logger

func init() {
	entry = logrus.New()
}

// DefaultLog default log
func DefaultLog() Logger {
	return entry
}

// SetDefault set default log
func SetDefault(v Logger) {
	entry = v
}
