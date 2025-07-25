package ylog

// Level log level
type Level uint8

func (l Level) String() string {
	switch l {
	case LevelError:
		return `Error`
	case LevelWarn:
		return `Warn`
	case LevelInfo:
		return `Info`
	case LevelDebug:
		return `Debug`
	default:
		return `Unknown`
	}
}

const (
	LevelError Level = iota + 1
	LevelWarn
	LevelInfo
	LevelDebug
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
	entry = NewLogger()
}

// DefaultLog default log
func DefaultLog() Logger {
	return entry
}

// SetDefault set default log
func SetDefault(v Logger) {
	entry = v
}
