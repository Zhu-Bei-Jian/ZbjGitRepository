package game

import (
	"errors"
)

type GameManager struct {
	games map[string]*GameBase
}

func newGameManager() *GameManager {
	return &GameManager{
		games: map[string]*GameBase{},
	}
}

func (gm *GameManager) allGames() map[string]*GameBase {
	return gm.games
}

func (gm *GameManager) add(gameId string, g *GameBase) error {
	if _, ok := gm.games[gameId]; ok {
		return errors.New("gameId exist")
	}
	gm.games[gameId] = g
	return nil
}

func (gm *GameManager) findGame(gameId string) (*GameBase, bool) {
	if g, ok := gm.games[gameId]; ok {
		return g, true
	}
	return nil, false
}

func (gm *GameManager) deleteGame(gameId string) {
	delete(gm.games, gameId)
}
