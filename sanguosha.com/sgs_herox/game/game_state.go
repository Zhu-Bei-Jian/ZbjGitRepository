package game

import (
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/game/kit/fsm"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

func NewGameInitState(g *GameBase) fsm.State {
	state := fsm.State{}

	state.Transitions = fsm.Transitions{
		"prepare": "prepare",
	}

	state.OnEnter = OnEnterGameInit(g)
	state.OnLeave = OnCommonLeaveState(g)
	return state
}

func OnEnterGameInit(g *GameBase) func(*fsm.Event) {
	return func(e *fsm.Event) {
		logrus.Debug("进入初始阶段")
	}
}

func NewGamePrepareState(g *GameBase) fsm.State {
	state := fsm.State{}

	state.Transitions = fsm.Transitions{
		"end":   "end",
		"start": "start",
	}

	state.OnEnter = OnEnterGamePrepare(g)
	state.OnLeave = OnCommonLeaveState(g)
	return state
}

func OnEnterGamePrepare(g *GameBase) func(*fsm.Event) {
	return func(e *fsm.Event) {
		logrus.Debug("进入准备阶段")
		//g.ClearAction()
		data := NewActDataGameState(g, "Prepare", true)
		g.PostActData(data)

		data.PostActStream(func() {
			g.setPhaseAndSync(gamedef.GamePhase_Ready)

			g.SetAllClientReady(true)
			g.Send2All(&cmsg.SNoticeGameReady{
				GameMode: g.setting.GameMode,
				Seats:    nil,
				GameUUID: g.gameUUID,
				RoomId:   g.roomId,
				VoiceId:  "",
			})
			g.SetAllClientReady(false)
			g.checkStateSync()
		})

		data.PostActStream(func() {
			g.modelImp.OnEnterState(e.Dst)
		})

		data.PostActStream(func() {
			g.FSMEvent("start")
		})
	}
}

func OnCommonLeaveState(g *GameBase) func(*fsm.Event) {
	return func(e *fsm.Event) {
		g.ClearAction()
	}
}

func NewGameStartState(g *GameBase) fsm.State {
	state := fsm.State{}

	state.Transitions = fsm.Transitions{
		"init_card": "init_card",
	}

	state.OnEnter = OnEnterGameStart(g)
	state.OnLeave = OnCommonLeaveState(g)
	return state
}

func OnEnterGameStart(g *GameBase) func(*fsm.Event) {
	return func(e *fsm.Event) {
		logrus.Debug("进入开始阶段")

		data := NewActDataGameState(g, "Start", true)
		g.PostActData(data)

		data.PostActStream(func() {
			g.setPhaseAndSync(gamedef.GamePhase_Start)
			g.noticeGameStart()
		})

		data.PostActStream(func() {
			g.modelImp.OnEnterState(e.Dst)
		})

		data.PostActStream(func() {
			g.StartWaitingNone(3, nil)
		})

		data.PostActStream(func() {
			g.FSMEvent("init_card")
		})
	}
}

func NewInitCardsState(g *GameBase) fsm.State {
	state := fsm.State{}

	state.Transitions = fsm.Transitions{
		"phase_begin": "phase_begin",
	}
	state.InternalEvents = fsm.Callbacks{}
	state.OnEnter = OnEnterInitCards(g)
	state.OnLeave = OnCommonLeaveState(g)

	return state
}

func OnEnterInitCards(g *GameBase) func(e *fsm.Event) {
	return func(e *fsm.Event) {
		logrus.Debug("进入初始牌准备阶段")
		//g.ClearAction()
		data := NewActDataGameState(g, "InitCards", true)
		g.PostActData(data)

		data.PostActStream(func() {
			g.setPhaseAndSync(gamedef.GamePhase_InitCard)
			g.modelImp.OnEnterState(e.Dst)
		})

		data.PostActStream(func() {
			g.FSMEvent("phase_begin")
		})
	}
}

func NewPhaseBeginState(g *GameBase) fsm.State {
	state := fsm.State{}

	state.Transitions = fsm.Transitions{
		"phase_begin": "phase_begin",
		"phase_main":  "phase_main",
		"phase_draw":  "phase_draw",
	}

	state.InternalEvents = fsm.Callbacks{}
	state.OnEnter = OnEnterPhaseBeginState(g)
	state.OnLeave = OnCommonLeaveState(g)

	return state
}

func OnEnterPhaseBeginState(g *GameBase) func(e *fsm.Event) {
	return func(e *fsm.Event) {

		//g.ClearAction()
		data := NewActDataGameState(g, "PhaseBegin", true)
		g.PostActData(data)

		data.PostActStream(func() {
			g.modelImp.OnEnterState(e.Dst)
		})

		data.PostActStream(func() {
			g.FSMEvent("phase_draw")
		})
	}
}

func NewPhaseDrawState(g *GameBase) fsm.State {
	state := fsm.State{}
	state.Transitions = fsm.Transitions{
		"phase_main":    "phase_main",
		"phase_cleanup": "phase_cleanup",
	}
	state.InternalEvents = fsm.Callbacks{}
	state.OnEnter = OnEnterPhaseDrawState(g)
	state.OnLeave = OnCommonLeaveState(g)
	return state
}

func OnEnterPhaseDrawState(g *GameBase) func(e *fsm.Event) {
	return func(e *fsm.Event) {
		//logrus.Debug("进入补牌阶段")
		//g.ClearAction()

		data := NewActDataGameState(g, "PhaseDraw", true)
		g.PostActData(data)
		data.PostActStream(func() {
			g.setPhaseAndSync(gamedef.GamePhase_PhaseDraw)
			g.modelImp.OnEnterState(e.Dst)
		})

		data.PostActStream(func() {
			g.GetTriggerMgr().Trigger(TriggerType_DrawStart, &ActionDataBase{})
		})

		data.PostActStream(func() {
			g.FSMEvent("phase_main")
		})
	}
}

func NewPhaseMainState(g *GameBase) fsm.State {
	state := fsm.State{}

	state.Transitions = fsm.Transitions{
		"phase_main":  "phase_main",
		"phase_begin": "phase_begin",
		"phase_end":   "phase_end",
	}

	state.InternalEvents = fsm.Callbacks{}
	state.OnEnter = OnEnterPhaseMainState(g)
	state.OnLeave = OnCommonLeaveState(g)

	return state
}

func OnEnterPhaseMainState(g *GameBase) func(e *fsm.Event) {
	return func(e *fsm.Event) {
		curPlayer := g.GetCurrentPlayer()

		data := NewActDataGameState(g, "PhaseMain", true)
		g.PostActData(data)
		data.PostActStream(func() {
			g.setPhaseAndSync(gamedef.GamePhase_PhaseMain)
			g.modelImp.OnEnterState(e.Dst)
		})

		data.PostActStream(func() {
			playStateWaitPlay(g, curPlayer, data)
		})

	}
}

func playStateWaitPlay(g *GameBase, curPlayer *Player, state *ActDataGameState) {
	NoticeOpActionStart(curPlayer, func(timeout bool) {
		if timeout {
			state.status = false
		}
	})

	state.PostActStream(func() {
		playStateCheckEnd(g, curPlayer, state)
	})
}

func playStateCheckEnd(g *GameBase, curPlayer *Player, state *ActDataGameState) {
	if curPlayer.IsDead() {
		g.FSMEvent("phase_begin")
		return
	}

	if g.phaseLeftSec() <= 0 {
		state.status = false
	}

	if state.status {
		state.PostActStream(func() {
			playStateWaitPlay(g, curPlayer, state)
		})
	} else {
		state.PostActStream(func() {
			g.FSMEvent("phase_end")
		})
	}
}

func NewPhaseEndState(g *GameBase) fsm.State {
	state := fsm.State{}
	state.Transitions = fsm.Transitions{
		"phase_begin": "phase_begin",
	}
	state.InternalEvents = fsm.Callbacks{}
	state.OnEnter = OnEnterPhaseEndState(g)
	state.OnLeave = OnCommonLeaveState(g)
	return state
}

func OnEnterPhaseEndState(g *GameBase) func(e *fsm.Event) {
	return func(e *fsm.Event) {
		//g.ClearAction()
		data := NewActDataGameState(g, "PhaseEnd", true)
		g.PostActData(data)
		data.PostActStream(func() {
			g.setPhaseAndSync(gamedef.GamePhase_PhaseEnd)
			g.modelImp.OnEnterState(e.Dst)
		})

		data.PostActStream(func() {
			g.FSMEvent("phase_begin")
		})
	}
}

func NewGameEndState(g *GameBase) fsm.State {
	state := fsm.State{}
	state.Transitions = fsm.Transitions{
		"end": "end",
	}
	state.OnEnter = OnEnterGameEndState(g)
	state.OnLeave = OnCommonLeaveState(g)
	return state
}

func OnEnterGameEndState(g *GameBase) func(e *fsm.Event) {
	return func(e *fsm.Event) {
		//g.ClearAction()
		logrus.Debug("进入游戏结束阶段")
		data := NewActDataGameState(g, "GameEnd", true)
		g.PostActData(data)

		//标记游戏已结束
		g.SetOver()

		data.PostActStream(func() {
			g.setPhaseAndSync(gamedef.GamePhase_End)
			g.modelImp.OnEnterState(e.Dst)
		})

		data.PostActStream(func() {
			g.Clear()
		})
	}
}

func OnLeaveGameEndState(g *GameBase) func(e *fsm.Event) {
	return func(e *fsm.Event) {
	}
}
