package ylog

import (
	"io"
)

type Option func(*Entry)

func WithLevel(level Level) Option {
	return func(e *Entry) {
		e.level = level
	}
}

func WithFile(enable, line bool, depth ...int) Option {
	return func(e *Entry) {
		e.outFile = enable
		e.fileLin = line
		if len(depth) > 1 {
			e.fileDep = depth[1]
		}
	}
}

func WithTime(enable bool, format ...string) Option {
	return func(e *Entry) {
		e.outTime = enable
		if len(format) > 0 {
			e.timeFmt = format[0]
		}
	}
}

func WithWriter(writer io.WriteCloser) Option {
	return func(e *Entry) {
		e.output = writer
	}
}
