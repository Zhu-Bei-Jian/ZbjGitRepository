package core

import (
	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/game/kit/fsm"
	"sync/atomic"
	"time"
)

type LogicCore struct {
	fsm *fsm.FSM

	worker Worker

	actionData *arraystack.Stack

	delayActionData []interface{}

	priorityFuncCount int32

	waitMgr *waitManager

	waitParallelMgr *WaitParallelManager
}

func (l *LogicCore) Init(worker Worker) {
	l.worker = worker
	l.actionData = arraystack.New()
	l.waitMgr = newWaitManager(l)
	l.waitParallelMgr = newWaitParallelManager()
}

func (l *LogicCore) InitFSM(fsm *fsm.FSM) {
	l.fsm = fsm
}

func (g *LogicCore) PostActData(v IActionData) {
	g.delayActionData = append(g.delayActionData, v)
}

func (g *LogicCore) DoNow(f func()) {
	atomic.AddInt32(&g.priorityFuncCount, 1)
	g.worker.Post(func() {
		atomic.AddInt32(&g.priorityFuncCount, -1)
		if f != nil {
			f()
		}
		g.ActivateAction()
	})
}

func (g *LogicCore) hasPriorityFunc() bool {
	return atomic.LoadInt32(&g.priorityFuncCount) != 0
}

func (g *LogicCore) ClearAction() {
	g.actionData.Clear()
	g.delayActionData = nil
}

func (g *LogicCore) ActivateAction() {
	for {
		if g.hasPriorityFunc() {
			return
		}

		if g.waitParallelMgr.WaitingCount() > 0 {
			return
		}

		if g.waitMgr.isWaiting() {
			return
		}

		length := len(g.delayActionData)

		if length > 0 {
			for i := length; i != 0; i-- {
				d := g.delayActionData[i-1]
				g.actionData.Push(d)
			}
			g.delayActionData = nil
		}

		peek, ok := g.actionData.Peek()
		if !ok {
			return
		}
		data := peek.(IActionData)
		if data == nil {
			g.actionData.Pop()
			logrus.Error("action data peek error")
			continue
		}

		if data.IsStop() {
			//data.Clear()
			g.actionData.Pop()
			continue
		}

		front, exist := data.Front()
		if !exist {
			g.actionData.Pop()
			continue
		}

		front()
	}
}

func (g *LogicCore) FSMEvent(event string, args ...interface{}) {
	g.DoNow(func() {
		g.fsm.Event(event, args)
	})
}

func (g *LogicCore) ActDataPeek() (value IActionData) {
	v, ok := g.actionData.Peek()
	if !ok {
		return nil
	}
	return v.(IActionData)
}

func (g *LogicCore) StartWaitingNone(sec int64, callback func()) {
	cb := func(timout bool) {
		if callback != nil {
			callback()
		}
	}
	g.waitMgr.clear()
	g.waitMgr.wait(time.Duration(sec)*time.Second, cb)
}

func (g *LogicCore) StartWaitingNoneFloat(sec float32, callback func()) {
	cb := func(timout bool) {
		if callback != nil {
			callback()
		}
	}
	g.waitMgr.clear()
	var count = int64(sec * 1000)
	g.waitMgr.wait(time.Duration(count)*time.Millisecond, cb)
}

func (g *LogicCore) StartWaiting(player []Player, sec int64, opMsg proto.Message, msgAllow []proto.Message, callback func(timeout bool)) {
	g.waitMgr.clear()
	for _, v := range player {
		g.waitMgr.addWaitPlayer(v, opMsg)
	}
	g.waitMgr.addWaitMsg(msgAllow...)
	g.waitMgr.wait(time.Duration(sec)*time.Second, callback)
}

func (g *LogicCore) CheckAllow(p Player, msg proto.Message, f func()) {
	if !g.waitMgr.isWaiting() {
		return
	}
	if !g.waitMgr.isWaitPlayerMsg(p, msg) {
		return
	}
	f()
}

func (g *LogicCore) SetPlayerAnswered(p Player) {
	g.waitMgr.answer(p)
}

func (g *LogicCore) GetPlayerLastOp(p Player) (*LastOp, bool) {
	return g.waitMgr.getPlayerLastOp(p)
}

func (g *LogicCore) StopWaiting() {
	g.waitMgr.forceStop()
}

func (g *LogicCore) ForceTimeout() {
	g.waitMgr.forceTimeout()
}

func (g *LogicCore) WaitingParallel(duration time.Duration, onTimeout func()) *WaitParallelTimer {
	return g.waitParallelMgr.StartWaiting(duration, func() {
		g.DoNow(func() {
			onTimeout()
		})
	})
}
