package ylog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/azeroth-sha/y/ybuff"
)

type Entry struct {
	level   Level          // 日志等级
	outFile bool           // 输出文件名
	fileDep int            // 文件深度
	fileLin bool           // 文件行号
	outTime bool           // 输出时间
	timeFmt string         // 时间格式
	output  io.WriteCloser // 日志输出
}

func (e *Entry) Debug(args ...any) {
	buf := e.outBuf(LevelDebug, args...)
	defer ybuff.Put(buf)
	e.wrtBuf(buf)
}

func (e *Entry) Debugf(format string, args ...any) {
	buf := e.outFmtBuf(LevelDebug, format, args...)
	defer ybuff.Put(buf)
	e.wrtBuf(buf)
}

func (e *Entry) Info(args ...any) {
	buf := e.outBuf(LevelInfo, args...)
	defer ybuff.Put(buf)
	e.wrtBuf(buf)
}

func (e *Entry) Infof(format string, args ...any) {
	buf := e.outFmtBuf(LevelInfo, format, args...)
	defer ybuff.Put(buf)
	e.wrtBuf(buf)
}

func (e *Entry) Warn(args ...any) {
	buf := e.outBuf(LevelWarn, args...)
	defer ybuff.Put(buf)
	e.wrtBuf(buf)
}

func (e *Entry) Warnf(format string, args ...any) {
	buf := e.outFmtBuf(LevelWarn, format, args...)
	defer ybuff.Put(buf)
	e.wrtBuf(buf)
}

func (e *Entry) Error(args ...any) {
	buf := e.outBuf(LevelError, args...)
	defer ybuff.Put(buf)
	e.wrtBuf(buf)
}

func (e *Entry) Errorf(format string, args ...any) {
	buf := e.outFmtBuf(LevelError, format, args...)
	defer ybuff.Put(buf)
	e.wrtBuf(buf)
}

// NewLogger returns a new logger
func NewLogger(opts ...Option) Logger {
	l := &Entry{
		level:   LevelInfo,
		outFile: false,
		fileDep: 4,
		fileLin: true,
		outTime: true,
		timeFmt: time.RFC3339,
		output:  &nopCloseWriter{os.Stdout},
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

/*
  Package private
*/

func (e *Entry) wrtBuf(b *bytes.Buffer) {
	if b == nil {
		return
	}
	_, _ = fmt.Fprint(e.output, b.String())
}

func (e *Entry) outBuf(l Level, args ...any) *bytes.Buffer {
	if e.level < l {
		return nil
	}
	buf := ybuff.Get()
	e.wrtPre(buf, l)
	_, _ = fmt.Fprint(buf, fmt.Sprintf(`message=%s`, fmt.Sprint(args...)))
	return buf
}

func (e *Entry) outFmtBuf(l Level, format string, args ...any) *bytes.Buffer {
	if e.level < l {
		return nil
	}
	buf := ybuff.Get()
	e.wrtPre(buf, l)
	_, _ = fmt.Fprintf(buf, fmt.Sprintf(`message=%s`, fmt.Sprintf(format, args...)))
	return buf
}

func (e *Entry) wrtPre(b *bytes.Buffer, l Level) {
	_, _ = fmt.Fprintf(b, "level=%s ", l.String())
	if e.outTime {
		_, _ = fmt.Fprintf(b, "time=%s ", time.Now().Format(e.timeFmt))
	}
	if e.outFile {
		name, line := e.callInfo()
		if e.fileLin {
			_, _ = fmt.Fprintf(b, "file=%s:%d ", name, line)
		} else {
			_, _ = fmt.Fprintf(b, "file=%s ", name)
		}
	}
}

func (e *Entry) callInfo() (string, int) {
	_, file, line, ok := runtime.Caller(e.fileDep)
	if !ok {
		return "", 0
	} else {
		d, n := path.Split(file)
		d = path.Base(d)
		return path.Join(d, n), line
	}
}

type nopCloseWriter struct {
	io.Writer
}

func (nopCloseWriter) Close() error { return nil }
