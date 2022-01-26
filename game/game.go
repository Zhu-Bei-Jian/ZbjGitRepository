package game

import (
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/game/kit/fsm"
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

type Game struct {
	base        *GameBase
	winSeat     int32
	cardIdIndex int32
}

func (g *Game) NewFSM(base *GameBase) *fsm.FSM {

	states := fsm.States{
		"*": fsm.State{
			Transitions: fsm.Transitions{
				"stop": "end",
			},
		},
		"init":        NewGameInitState(base),
		"prepare":     NewGamePrepareState(base),
		"start":       NewGameStartState(base),
		"init_card":   NewInitCardsState(base),
		"phase_begin": NewPhaseBeginState(base),
		"phase_draw":  NewPhaseDrawState(base),
		"phase_main":  NewPhaseMainState(base),
		"phase_end":   NewPhaseEndState(base),
		"end":         NewGameEndState(base),
	}

	callbacks := fsm.Callbacks{}
	result := fsm.NewFSM("init", states, callbacks)
	return result
}

func newGame(worker core.Worker, setting *gamedef.RoomSetting, gameUUID string, roomId, roomNO uint32, config *conf.GameConfig) *Game {
	base := newGameBase(worker, setting, gameUUID, roomId, roomNO, config)
	g := &Game{
		base: base,
	}
	fsm := g.NewFSM(base)
	base.InitFSM(fsm)
	base.modelImp = g
	return g
}

func (g *Game) OnEnterState(state string) {
	switch state {
	case "prepare":
		g.onPrepare()
	case "start":
		g.onStart()
	case "init_card":
		g.onInitCards()
	case "phase_begin":
		g.onPhaseBegin()
	case "phase_draw":
		g.onPhaseDraw()
	case "phase_main":
		g.onPhaseMain()
	case "phase_end":
		g.onPhaseEnd()
	case "end":
		g.onEnd()
	default:

	}
}
