package manager

import (
	"errors"
	"fmt"
	"sanguosha.com/baselib/util"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/ioservice"
)

// WorkerManager ...
type WorkerManager struct {
	worker iWorker
	io     ioservice.IOService
}

type iWorker interface {
	ID() uint32
	Post(f func())
	AfterPost(d time.Duration, f func()) (cancel func() bool)
	WorkerLen() int
}

// Init ...
func (p *WorkerManager) Init(worker iWorker) {
	p.worker = worker
	p.io = ioservice.NewIOService(fmt.Sprintf("worker-mgr-%d", worker.ID()), 10240)
}

// Run ...
func (p *WorkerManager) Run() {
	p.io.Run()
}

// Close ...
func (p *WorkerManager) Close() {
	p.io.Fini()
}

// Post ...
func (p *WorkerManager) Post(f func()) {
	p.io.Post(func() {
		p.worker.Post(f)
	})
}

// AfterPost ...
func (p *WorkerManager) AfterPost(d time.Duration, f func()) (cancel func() bool) {
	return p.io.AfterPost(d, func() {
		p.worker.Post(f)
	})
}

func (p *WorkerManager) panicHandler(pnc interface{}) {
	var err error
	switch pnc.(type) {
	case string:
		err = errors.New(pnc.(string))
	case error:
		err = pnc.(error)
	default:
		err = errors.New("unknown panic")
	}

	logrus.WithError(err).Error("workManager panic")
}

// Go ...
func (p *WorkerManager) Go(callback func()) {
	util.SafeGo(func() {
		defer func() {
			if err := recover(); err != nil {
				p.panicHandler(err)
			}
		}()
		callback()
	})
}

// Ticker ...
func (p *WorkerManager) Ticker(d time.Duration, callback func()) *time.Ticker {
	ticker := time.NewTicker(d)
	p.Go(func() {
		for range ticker.C {
			p.worker.Post(callback)
		}
	})
	return ticker
}

// WaitTimeout ...
func (p *WorkerManager) WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	p.Go(func() {
		defer close(c)
		wg.Wait()
	})
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
