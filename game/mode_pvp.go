package game

import "sanguosha.com/sgs_herox/game/kit/fsm"

type GamePVP struct {
	base *GameBase
}

func NewGamePVP() *GameBase {
	return nil
}

func (g *GamePVP) InitFSM(base *GameBase) *fsm.FSM {
	g.base = base
	states := fsm.States{
		"*": fsm.State{
			Transitions: fsm.Transitions{
				"stop": "end",
			},
		},
		"prepare":     NewGamePrepareState(base),
		"init_card":   NewInitCardsState(base),
		"phase_begin": NewPhaseBeginState(base),
		"phase_draw":  NewPhaseDrawState(base),
		"phase_main":  NewPhaseMainState(base),
		"phase_end":   NewPhaseEndState(base),
		"end":         NewGameEndState(base),
	}

	callbacks := fsm.Callbacks{}
	result := fsm.NewFSM("prepare", states, callbacks)
	return result
}
