package core

import (
	"github.com/golang/protobuf/proto"
	"reflect"
	"time"
)

type waitManager struct {
	waitPlayers map[Player]*LastOp

	waitCallback func(bool)
	waitTimer    *time.Timer
	waitEndTime  int64
	waiting      bool

	waitMsgTypes map[reflect.Type]string

	core *LogicCore
}

type LastOp struct {
	SendTime int64
	Msg      proto.Message
}

func newWaitManager(l *LogicCore) *waitManager {
	return &waitManager{
		core:         l,
		waitMsgTypes: make(map[reflect.Type]string),
		waitPlayers:  make(map[Player]*LastOp),
	}
}

func (wm *waitManager) clear() {
	wm.waitPlayers = make(map[Player]*LastOp)
	wm.waiting = false
	wm.waitCallback = nil
	wm.waitTimer = nil
	wm.waitEndTime = 0
	wm.waitMsgTypes = make(map[reflect.Type]string)
}

func (wm *waitManager) waitInfo() (waitPlayers []Player, endTime int64) {
	for p, _ := range wm.waitPlayers {
		waitPlayers = append(waitPlayers, p)
	}
	endTime = wm.waitEndTime
	return
}

func (wm *waitManager) addWaitPlayer(p Player, msg proto.Message) {
	wm.waitPlayers[p] = &LastOp{
		Msg:      msg,
		SendTime: time.Now().Unix(),
	}
}

func (wm *waitManager) wait(duration time.Duration, waitCallback func(bool)) {

	f := func(timeout bool) {
		wm.clear()
		if waitCallback != nil {
			waitCallback(timeout)
		}
	}

	wm.waitCallback = f
	wm.waiting = true
	wm.waitEndTime = time.Now().Unix() + int64(duration/time.Second)
	wm.waitTimer = time.AfterFunc(duration, func() {
		wm.core.DoNow(func() {
			f(true)
		})
	})
}

func (wm *waitManager) resetIfWaitTimeLong(duration time.Duration) (int64, bool) {
	leftSec := wm.waitEndTime - time.Now().Unix()

	//所剩时间小于重置时间，不用重置
	if leftSec <= int64(duration/time.Second) {
		return 0, false
	}

	if !wm.waitTimer.Stop() {
		return 0, false
	}

	wm.waitEndTime = time.Now().Unix() + int64(duration/time.Second)
	wm.waitTimer.Reset(duration)
	return wm.waitEndTime, true
}

func (wm *waitManager) addWaitMsg(list ...proto.Message) {
	for _, v := range list {
		wm.waitMsgTypes[reflect.TypeOf(v)] = proto.MessageName(v)
	}
}

func (wm *waitManager) isWaitMsg(m proto.Message) bool {
	for k, _ := range wm.waitMsgTypes {
		if reflect.TypeOf(m) == k {
			return true
		}
	}

	return false
}

func (wm *waitManager) isPlayerWait(player Player) bool {
	lastOp, exist := wm.waitPlayers[player]
	if !exist {
		return false
	}

	if lastOp == nil {
		return false
	}

	return true
}

func (wm *waitManager) isWaitPlayerMsg(p Player, m proto.Message) bool {
	_, exist := wm.waitPlayers[p]
	if !exist {
		return false
	}

	if !wm.isWaitMsg(m) {
		return false
	}
	return true
}

func (wm *waitManager) getPlayerLastOp(p Player) (*LastOp, bool) {
	lastOp, exist := wm.waitPlayers[p]
	if exist {
		return lastOp, true
	}
	if lastOp == nil {
		return nil, false
	}
	return nil, false
}

func (wm *waitManager) isWait(player Player, m proto.Message) bool {
	if !wm.isWaitMsg(m) {
		return false
	}

	if !wm.isPlayerWait(player) {
		return false
	}

	return true
}

func (wm *waitManager) answer(player Player) {
	wm.setAnswer(player)
	if wm.isAllAnswer() {
		wm.forceWaitCallback()
	}
}

func (wm *waitManager) forceWaitCallback() {
	wm.waiting = false
	if ok := wm.waitTimer.Stop(); ok {
		wm.waitCallback(false)
	}
}

func (wm *waitManager) forceTimeout() {
	wm.waiting = false
	if ok := wm.waitTimer.Stop(); ok {
		wm.waitCallback(true)
	}
}

func (wm *waitManager) forceStop() {
	wm.waiting = false
	wm.waitTimer.Stop()
}

func (wm *waitManager) setAnswer(player Player) {
	wm.waitPlayers[player] = nil
}

func (wm *waitManager) isAllAnswer() bool {
	for _, v := range wm.waitPlayers {
		if v != nil {
			return false
		}
	}
	return true
}

func (wm *waitManager) isWaiting() bool {
	return wm.waiting
}
