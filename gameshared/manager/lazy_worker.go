package manager

import (
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/util"
	"time"
)

type LazyInfo struct {
	reqTime int64
	reqFunc func()
}

//将任务存起来慢慢做 做一点 然后等一会儿再做
type LazyWorker struct {
	worker    iWorker
	list      []*LazyInfo
	timeout   int64
	each      int64
	frameTime time.Duration
}

func NewLazyWorker(worker iWorker, each int64, frameTime time.Duration, timeout int64) *LazyWorker {
	if each < 1 {
		each = 1
	}
	dw := &LazyWorker{worker: worker, each: each, frameTime: frameTime, timeout: timeout}
	return dw
}

func (dw *LazyWorker) Post(f func()) {
	dw.list = append(dw.list, &LazyInfo{time.Now().Unix(), func() {
		defer util.Recover()
		f()
	}})

	if len(dw.list) == 1 {
		dw.worker.Post(dw.realRun)
		logrus.Debug("LazyWorker) Post")
	}
}

func (dw *LazyWorker) realRun() {
	total := int64(len(dw.list))
	if total == 0 {
		return
	}
	defer func() {
		if len(dw.list) != 0 {
			if dw.frameTime != 0 {
				dw.worker.AfterPost(dw.frameTime, dw.realRun)
			} else {
				dw.worker.Post(dw.realRun)
			}
		}
	}()

	n := dw.each
	if n > total {
		n = total
	}
	for ; n > 0; n-- {
		delayInfo := dw.list[0]
		if delayInfo == nil {
			dw.list = dw.list[1:]
			continue
		}
		if dw.timeout != 0 {
			nowTime := time.Now().Unix()
			if nowTime > delayInfo.reqTime {
				d := nowTime - delayInfo.reqTime
				if d >= dw.timeout {
					dw.list = dw.list[1:]
					continue
				}
			}
		}
		delayInfo.reqFunc()
		dw.list = dw.list[1:]
		logrus.Debug("LazyWorker) realRun")
	}
}
