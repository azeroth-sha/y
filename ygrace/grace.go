package ygrace

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/azeroth-sha/y/ylog"
)

type Service interface {
	Serv(ylog.Logger) error
	Down(ylog.Logger) error
}

type Grace interface {
	Serv()
	Down()
	Register(name string, svr Service)
}

type grace struct {
	running int32
	closed  chan struct{}
	wait    sync.WaitGroup
	once    sync.Once
	opts    option
	svrMu   sync.Mutex
	svrList map[string]Service
}

func (g *grace) Serv() {
	if atomic.SwapInt32(&g.running, 1) != 0 {
		return
	}
	g.closed = make(chan struct{})
	for name, svr := range g.svrList {
		g.wait.Add(1)
		go g.serv(name, svr)
	}
}

func (g *grace) Down() {
	if atomic.SwapInt32(&g.running, 0) != 1 {
		return
	}
	close(g.closed)
	for name, svr := range g.svrList {
		g.down(name, svr)
	}
	g.wait.Wait()
}

func (g *grace) Register(name string, svr Service) {
	g.svrMu.Lock()
	defer g.svrMu.Unlock()
	if g.svrList == nil {
		g.svrList = make(map[string]Service)
	}
	g.svrList[name] = svr
}

// New returns a new Grace.
func New(opts ...Option) Grace {
	g := new(grace)
	g.opts = option{
		dur: time.Second,
		log: ylog.DefaultLog(),
	}
	for _, opt := range opts {
		opt(&g.opts)
	}
	return g
}

/*
  Package private
*/

func (g *grace) down(name string, svr Service) {
	log := g.opts.log
	log.Infof("stopping: %s", name)
	defer log.Infof("stopped: %s", name)
	if err := svr.Down(log); err != nil {
		log.Errorf("stoping error: %s -> %v", name, err)
	}
}

func (g *grace) serv(name string, svr Service) {
	defer g.wait.Done()
	dur := g.opts.dur
	log := g.opts.log
EXIT:
	for atomic.LoadInt32(&g.running) == 1 {
		log.Infof("starting: %s", name)
		g.run(log, name, svr)()
		if dur > 0 {
			log.Infof("waiting: %s -> %s", name, dur)
			select {
			case <-time.After(dur):
			case <-g.closed:
				break EXIT
			}
		}
	}
}

func (g *grace) run(l ylog.Logger, name string, svr Service) func() {
	return func() {
		defer func() {
			if rec := recover(); rec != nil {
				l.Errorf("running panic: %s -> %v", name, rec)
			}
		}()
		l.Infof("running: %s", name)
		if err := svr.Serv(l); err != nil {
			l.Errorf("running error: %s -> %v", name, err)
		}
	}
}
