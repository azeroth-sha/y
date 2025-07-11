package ygrace

import (
	"cmp"
	"errors"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/azeroth-sha/y/logger"
)

var (
	ErrServing = errors.New("server is running")
)

type Grace interface {
	Serv() error
	Down() error
	Register(Service) error
}

type grace struct {
	mu      *sync.Mutex
	opt     *option
	srv     []Service
	running int32
	closed  chan struct{}
}

func (g *grace) Serv() error {
	if atomic.SwapInt32(&g.running, 1) != 0 {
		return ErrServing
	}
	defer func() {
		_ = g.Down()
	}()
	g.closed = make(chan struct{})
	all := make([]Service, 0, len(g.srv))
	for _, srv := range g.srv {
		all = append(all, srv)
	}
	slices.SortFunc(all, func(a, b Service) int {
		return cmp.Compare(a.Priority(), b.Priority())
	})
	log := g.opt.log
	for _, srv := range all {
		if h, y := srv.(ServWait); y {
			log.Info("waiting:", srv.Name())
			time.Sleep(h.ServWait())
		}
		log.Info("serving:", srv.Name())
		go g.serv(srv)
		log.Info("served:", srv.Name())
	}
	return nil
}

func (g *grace) Down() error {
	if atomic.SwapInt32(&g.running, 0) != 1 {
		return nil
	}
	close(g.closed)
	all := make([]Service, 0, len(g.srv))
	for _, srv := range g.srv {
		all = append(all, srv)
	}
	slices.SortFunc(all, func(a, b Service) int {
		return cmp.Compare(b.Priority(), a.Priority())
	})
	log := g.opt.log
	for _, srv := range all {
		if h, y := srv.(DownWait); y {
			log.Info("waiting:", srv.Name())
			time.Sleep(h.DownWait())
		}
		log.Info("downing:", srv.Name())
		g.down(srv)
		log.Info("downed:", srv.Name())
	}
	return nil
}

func (g *grace) Register(svr Service) error {
	if g.isRunning() {
		return ErrServing
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.srv = append(g.srv, svr)
	return nil
}

// New returns a new Grace instance.
func New(opts ...Option) Grace {
	g := &grace{
		closed: nil,
		mu:     new(sync.Mutex),
		opt: &option{
			dur: time.Second,
			log: logger.DefaultLog(),
		},
		running: 0,
		srv:     make([]Service, 0),
	}
	for _, opt := range opts {
		opt(g.opt)
	}
	return g
}

/*
  Package private
*/

func (g *grace) down(srv Service) {
	log := g.opt.log
	defer func() {
		if rec := recover(); rec != nil {
			log.Errorf("shutdown panic: [%s] -> %v", srv.Name(), rec)
		}
	}()
	if err := srv.Down(log); err != nil {
		log.Errorf("shutdown error: [%s] -> %v", srv.Name(), err)
	}
}

func (g *grace) serv(srv Service) {
	log := g.opt.log
	dur := g.opt.dur
EXIT:
	for g.isRunning() {
		select {
		case <-g.closed:
			break EXIT
		default:
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						log.Errorf("running panic: [%s] -> %v", srv.Name(), rec)
					}
				}()
				if err := srv.Serv(log); err != nil {
					log.Errorf("running error: [%s] -> %v", srv.Name(), err)
				}
			}()
			if dur > 0 && g.isRunning() {
				time.Sleep(dur)
			}
		}
	}
}

func (g *grace) isRunning() bool {
	return atomic.LoadInt32(&g.running) == 1
}
