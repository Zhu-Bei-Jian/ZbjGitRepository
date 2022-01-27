package appframe

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

// Worker 定义如何执行 ListenMsg 时设置的回调函数
type Worker interface {
	// 执行工作任务
	Post(func())

	// 等待当前所有任务执行完毕
	WaitDone()

	GroutineCount() int32
}

type parallelWorker struct {
	wg            sync.WaitGroup
	ch            chan struct{}
	recover       bool
	groutineCount int32
}

func (w *parallelWorker) GroutineCount() int32 {
	return atomic.LoadInt32(&w.groutineCount)
}

func (w *parallelWorker) Post(f func()) {
	w.wg.Add(1)
	atomic.AddInt32(&w.groutineCount, 1)
	if w.ch == nil {
		go func() {
			defer func() {
				atomic.AddInt32(&w.groutineCount, -1)
				w.wg.Done()
				if w.recover {
					if err := recover(); err != nil {
						stack := debug.Stack()
						logrus.WithFields(logrus.Fields{
							"err":   err,
							"stack": string(stack),
						}).Error("Recover")

						os.Stderr.Write([]byte(fmt.Sprintf("%v\n", err)))
						os.Stderr.Write(stack)
					}
				}
			}()
			f()
		}()
	} else {
		<-w.ch
		go func() {
			defer func() {
				atomic.AddInt32(&w.groutineCount, -1)
				w.wg.Done()
				w.ch <- struct{}{}
				if w.recover {
					if err := recover(); err != nil {
						stack := debug.Stack()
						logrus.WithFields(logrus.Fields{
							"err":   err,
							"stack": string(stack),
						}).Error("Recover")

						os.Stderr.Write([]byte(fmt.Sprintf("%v\n", err)))
						os.Stderr.Write(stack)
					}
				}
			}()
			f()
		}()
	}
}
func (w *parallelWorker) WaitDone() {
	w.wg.Wait()
}

// NewParallelWorker 创建并发的执行器.
// maxGoroutines 为最大 goroutine 并发数量, maxGoroutines <= 0 时, 表示不限.
// recover 表示是否需要启用 recover 来保护 panic.
func NewParallelWorker(maxGoroutines int, recover bool) Worker {
	w := &parallelWorker{recover: recover}
	if maxGoroutines > 0 {
		w.ch = make(chan struct{}, maxGoroutines)
		for i := 0; i < maxGoroutines; i++ {
			w.ch <- struct{}{}
		}
	}
	return w
}
