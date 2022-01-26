package core

import (
	"time"
)

type WaitParallelManager struct {
	index int64

	timers map[int64]*WaitParallelTimer
}

func newWaitParallelManager() *WaitParallelManager {
	return &WaitParallelManager{
		index:  0,
		timers: make(map[int64]*WaitParallelTimer),
	}
}

func (wm *WaitParallelManager) newIndex() int64 {
	wm.index++
	return wm.index
}

type WaitParallelTimer struct {
	timer *time.Timer
	index int64
	mgr   *WaitParallelManager
}

func (p *WaitParallelTimer) Stop() bool {
	ok := p.timer.Stop()
	p.mgr.remove(p.index)
	return ok
}

func (wm *WaitParallelManager) StartWaiting(t time.Duration, waitCallback func()) *WaitParallelTimer {
	index := wm.newIndex()

	timer := time.AfterFunc(t, func() {
		wm.remove(index)
		if waitCallback != nil {
			waitCallback()
		}
	})

	waitTimer := &WaitParallelTimer{
		timer: timer,
		index: index,
		mgr:   wm,
	}
	wm.timers[index] = waitTimer
	return waitTimer
}

func (wm *WaitParallelManager) WaitingCount() int {
	return len(wm.timers)
}

func (wm *WaitParallelManager) remove(index int64) {
	delete(wm.timers, index)
}
