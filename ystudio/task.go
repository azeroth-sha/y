package ystudio

import (
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/azeroth-sha/y/ylog"
)

var (
	ErrDone = errors.New("event done")
)

type Event interface {
	Name() (name string)
	Occur() (occur time.Time)
	Request() (v any)
	Logger() (log ylog.Logger)
	Release(v any, err ...error) error
}

type Reply interface {
	Wait() (v any, err error)
}

type task struct {
	mu       *sync.Mutex
	name     string
	occur    time.Time
	log      ylog.Logger
	request  any
	response chan any
	err      error
	done     bool
}

func (t *task) Wait() (any, error) {
	if resp, ok := <-t.response; ok {
		return resp, t.err
	}
	return nil, ErrDone
}

func (t *task) Name() (name string) {
	return t.name
}

func (t *task) Occur() (occur time.Time) {
	return t.occur
}

func (t *task) Request() (v any) {
	return t.request
}

func (t *task) Logger() (log ylog.Logger) {
	return t.log
}

func (t *task) Release(v any, err ...error) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.done {
		return ErrDone
	}
	t.done = true
	t.err = errors.Join(err...)
	t.response <- v
	close(t.response)
	return nil
}

// newTask returns a new task
func newTask(name string, req any, log ylog.Logger) *task {
	t := &task{
		mu:       new(sync.Mutex),
		name:     name,
		occur:    time.Now(),
		log:      log,
		request:  req,
		response: make(chan any, 1),
		err:      nil,
		done:     false,
	}
	runtime.SetFinalizer(t, func(e *task) {
		_ = e.Release(nil)
	})
	return t
}
