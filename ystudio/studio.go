package ystudio

import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/azeroth-sha/y/ylog"
)

var (
	ErrJobNotExist  = errors.New("job not exist")
	ErrJobEventFull = errors.New("event full")
	ErrStudioClosed = errors.New("studio closed")
)

type (
	WorkHandler  func(event Event)
	AbortHandler func(event Event, rec any)
)

// Studio is a task manager
type Studio interface {
	// Append a task
	Append(name string, param any, opts ...AppendOption) Reply
	// SetJob set a job
	SetJob(name string, hand WorkHandler, capSize int) error
	// Release the studio
	Release()
	// Count the studio task count
	Count() int
	// Len the studio task length
	Len(name string) int
}

type studio struct {
	mu      *sync.Mutex
	running int32
	closed  chan struct{}
	abort   AbortHandler
	event   map[string]chan Event
	jobs    map[string]WorkHandler
}

func (s *studio) Append(name string, param any, opts ...AppendOption) Reply {
	cfg := &appendConfig{
		waitDur: -1,
		logger:  ylog.DefaultLog(),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	t := newTask(name, param, cfg.logger)
	var ch chan Event

	s.mu.Lock()
	if val, ok := s.event[name]; ok {
		s.mu.Unlock()
		ch = val
	} else {
		s.mu.Unlock()
		_ = t.Release(nil, ErrJobNotExist)
		return t
	}

	switch {
	case cfg.waitDur == -1:
		select {
		case <-s.closed:
			_ = t.Release(nil, ErrStudioClosed)
		case ch <- t:
		}
	case cfg.waitDur == 0:
		select {
		case <-s.closed:
			_ = t.Release(nil, ErrStudioClosed)
		case ch <- t:
		default:
			_ = t.Release(nil, ErrJobEventFull)
		}
	default:
		select {
		case <-s.closed:
			_ = t.Release(nil, ErrStudioClosed)
		case ch <- t:
		case <-time.After(cfg.waitDur):
			_ = t.Release(nil, ErrJobEventFull)
		}
	}
	return t
}

//func (s *studio) Append(name string, param any, args ...any) Reply {
//	logEntry := ylog.DefaultLog()
//	waitDur := time.Duration(-1)
//	for i := 0; i < len(args); i++ {
//		val := args[i]
//		switch vv := val.(type) {
//		case time.Duration:
//			waitDur = vv
//		case ylog.Logger:
//			logEntry = vv
//		}
//	}
//	t := newTask(name, param, logEntry)
//	var ch chan Event
//	s.mu.Lock()
//	if atomic.LoadInt32(&s.running) != 1 {
//		_ = t.Release(nil, ErrStudioClosed)
//		return t
//	} else if val, ok := s.event[name]; ok {
//		s.mu.Unlock()
//		ch = val
//	} else {
//		s.mu.Unlock()
//		_ = t.Release(nil, ErrJobNotExist)
//		return t
//	}
//	if waitDur <= -1 { // 阻塞等待
//		select {
//		case <-s.closed:
//			_ = t.Release(nil, ErrStudioClosed)
//		case ch <- t:
//		}
//	} else if waitDur == 0 { // 不阻塞等待
//		select {
//		case <-s.closed:
//			_ = t.Release(nil, ErrStudioClosed)
//		case ch <- t:
//		default:
//			_ = t.Release(nil, ErrJobEventFull)
//		}
//	} else { // 等待指定时长
//		select {
//		case <-s.closed:
//			_ = t.Release(nil, ErrStudioClosed)
//		case ch <- t:
//		case <-time.After(waitDur):
//			_ = t.Release(nil, ErrJobEventFull)
//		}
//	}
//	return t
//}

func (s *studio) SetJob(name string, hand WorkHandler, capSize int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if atomic.LoadInt32(&s.running) != 1 {
		return ErrStudioClosed
	}
	if v, ok := s.event[name]; ok {
		close(v)
		delete(s.event, name)
		delete(s.jobs, name)
	}
	ch := make(chan Event, capSize)
	s.event[name] = ch
	s.jobs[name] = hand
	go s.working(ch, hand)
	return nil
}

func (s *studio) Release() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if atomic.SwapInt32(&s.running, 0) != 1 {
		return
	}
	close(s.closed)
	for n, v := range s.event {
		close(v)
		delete(s.event, n)
	}
	for n, _ := range s.jobs {
		delete(s.jobs, n)
	}
}

func (s *studio) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	cnt := 0
	for _, v := range s.event {
		cnt += len(v)
	}
	return cnt
}

func (s *studio) Len(name string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	cnt := 0
	if ch, ok := s.event[name]; ok {
		cnt = len(ch)
	}
	return cnt
}

// New return a new studio.
func New(opts ...any) Studio {
	s := &studio{
		mu:      new(sync.Mutex),
		running: 1,
		closed:  make(chan struct{}),
		abort:   nil,
		event:   make(map[string]chan Event),
		jobs:    make(map[string]WorkHandler),
	}
	for i := 0; i < len(opts); i++ {
		val := opts[i]
		switch vv := val.(type) {
		case AbortHandler:
			s.abort = vv
		}
	}
	runtime.SetFinalizer(s, func(s *studio) {
		s.Release()
	})
	return s
}

/*
  Package private
*/

func (s *studio) working(ch chan Event, fun WorkHandler) {
	//ch := s.event[name]
	//fun := s.jobs[name]
EXIT:
	for {
		select {
		case <-s.closed:
			break EXIT
		case event, ok := <-ch:
			if !ok {
				break EXIT
			}
			s.run(event, fun)
		}
	}
}

func (s *studio) run(e Event, f WorkHandler) {
	abort := s.abort
	entry := e.Logger()
	defer func() {
		rec := recover()
		if rec != nil && abort != nil {
			abort(e, rec)
		}
		if rec != nil {
			err := fmt.Errorf(`panic: %v`, rec)
			_ = e.Release(nil, err)
			entry.Errorf(`task error: %s -> %v`, e.Name(), err)
		} else {
			_ = e.Release(nil)
			entry.Infof(`task done: %s`, e.Name())
		}
	}()
	f(e)
}
