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
	_, _ = b.WriteString("\r\n")
	_, _ = b.WriteTo(e.output)
}

func (e *Entry) outBuf(l Level, args ...any) *bytes.Buffer {
	if e.level < l {
		return nil
	}
	buf := ybuff.Get()
	e.wrtPre(buf, l)
	buf.WriteString(`msg=`)
	buf.WriteString(fmt.Sprint(args...))
	return buf
}

func (e *Entry) outFmtBuf(l Level, format string, args ...any) *bytes.Buffer {
	if e.level < l {
		return nil
	}
	buf := ybuff.Get()
	e.wrtPre(buf, l)
	buf.WriteString(`msg=`)
	buf.WriteString(fmt.Sprintf(format, args...))
	return buf
}

func (e *Entry) wrtPre(b *bytes.Buffer, l Level) {
	b.WriteString(`level=`)
	b.WriteString(l.String())
	b.WriteString(` `)
	if e.outTime {
		b.WriteString(`time=`)
		b.WriteString(time.Now().Format(e.timeFmt))
		b.WriteString(` `)
	}
	if e.outFile {
		name, line := e.callInfo()
		b.WriteString(`file=`)
		b.WriteString(name)
		if e.fileLin {
			b.WriteString(`:`)
			b.WriteString(fmt.Sprint(line))
		}
		b.WriteString(` `)
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
