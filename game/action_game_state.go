package game

import "sanguosha.com/sgs_herox/game/core"

type ActionDataBase struct {
	core.ActionDataCore
	game *GameBase
}

type ActDataGameState struct {
	ActionDataBase
	game   *GameBase
	state  string
	status bool
}

func NewActDataGameState(g *GameBase, state string, statue bool) *ActDataGameState {
	ad := &ActDataGameState{
		ActionDataBase: ActionDataBase{},
		game:           g,
		state:          state,
		status:         statue,
	}
	return ad
}
