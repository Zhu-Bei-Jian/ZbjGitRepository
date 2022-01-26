package ioservice

import (
	"expvar"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/util"
)

// NewIOService 创建ioservice
func NewIOService(name string, rcpChanLen int) IOService {
	var io = new(impl)

	io.rpcChan = make(chan *rpcEvent, rcpChanLen)

	tempname := fmt.Sprintf("IO服务-%s", name)
	io.rpcLen = expvar.NewInt(tempname)
	io.name = tempname

	return io
}

// impl ioservice的具体实现
type impl struct {
	IOService

	rpcChan chan *rpcEvent
	rpcLen  *expvar.Int

	closed int32
	finiWg sync.WaitGroup

	name string

	lastPrintTime int64
}

// Post 传递f到goroutine上执行
func (io *impl) Post(f func()) {
	io.RPCGo(f)
}

// TryPost 传递f到goroutine上执行
func (io *impl) TryPost(f func()) bool {
	event := newRPCEvent(f)
	select {
	case io.rpcChan <- newRPCEvent(f):
		if atomic.LoadInt32(&io.closed) == 0 {
			io.rpcLen.Add(1)
		} else {
			logrus.WithFields(logrus.Fields{
				"io":       io,
				"event":    event,
				"function": event.f,
			}).Error("ioservice has fini(TryPost)")

			if event.retChan != nil {
				event.retChan <- nil
			}
		}
		return true
	default:
		return false
	}
}

func (io *impl) AfterPost(d time.Duration, f func()) func() bool {
	return io.AfterRPCGo(d, f)
}

// goroutine safe
func (io *impl) RPCGo(f interface{}, args ...interface{}) {
	io.AfterRPCGo(0, f, args...)
}

func (io *impl) AfterRPCGo(d time.Duration, f interface{}, args ...interface{}) func() bool {
	event := newRPCEvent(f, args...)

	if d <= 0 {
		io.pushEvent(event)
		return func() bool {
			return false
		}
	}

	var stop int32
	t := time.AfterFunc(d, func() {
		if atomic.CompareAndSwapInt32(&stop, 0, 1) {
			io.pushEvent(event)
		}
	})

	return func() bool {
		if atomic.CompareAndSwapInt32(&stop, 0, 1) {
			t.Stop()
			return true
		}
		return false
	}
}

// goroutine safe
func (io *impl) RPCCall(f interface{}, args ...interface{}) interface{} {
	return io.AfterRPCCall(0, f, args...)
}

//
func (io *impl) AfterRPCCall(d time.Duration, f interface{}, args ...interface{}) interface{} {
	event := newRPCEvent2(f, args...)

	if d <= 0 {
		io.pushEvent(event)

	} else {
		time.AfterFunc(d, func() {
			io.pushEvent(event)
		})
	}

	if event.retChan != nil {
		return <-event.retChan
	}

	return nil
}

//
func (io *impl) Init() {

}

//
func (io *impl) Run() {
	io.finiWg.Add(1)

	go func() {
		defer util.Recover()
		defer func() {
			// 由于defer的调用比较耗内存，不能对外部的每个函数进行defer，所以采用了以下方式
			// 挂了重启，关了退出
			if atomic.LoadInt32(&io.closed) == 0 {
				logrus.Error("IO catch err")
				io.Run()
			}
			io.finiWg.Done()
		}()

		for curEvent := range io.rpcChan {
			io.rpcLen.Add(-1)
			curEvent.doRPC()
		}
	}()
}

func(io *impl)WorkerLen()int{
	return len(io.rpcChan)
}

//
func (io *impl) Fini() {
	if atomic.CompareAndSwapInt32(&io.closed, 0, 1) {
		close(io.rpcChan)
		io.finiWg.Wait()
	}
}

func (io *impl) pushEvent(event *rpcEvent) {
	if atomic.LoadInt32(&io.closed) == 0 {
		io.rpcChan <- event
		io.rpcLen.Add(1)

		//队列有点大了 打印日志(1秒打印一次)
		nowTime := time.Now().Unix()
		if (io.lastPrintTime != nowTime) && (len(io.rpcChan) >= (cap(io.rpcChan) / 2)) {
			io.lastPrintTime = nowTime
			logrus.WithFields(logrus.Fields{
				"eventCount": io.rpcLen.Value(),
				"name ":      io.name,
			}).Warning("ioservice too long")
		}
	} else {
		logrus.WithFields(logrus.Fields{
			"io":       io,
			"event":    event,
			"function": event.f,
		}).Error("ioservice has fini")

		if event.retChan != nil {
			event.retChan <- nil
		}
	}
}
