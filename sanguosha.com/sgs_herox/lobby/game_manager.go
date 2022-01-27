package lobby

import (
	"errors"
	"fmt"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameshared/manager"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"time"
)

var (
	ErrNoPlayer       = errors.New("game no player")
	ErrServerLoadFull = errors.New("server load full")
	ErrNoAiServer     = errors.New("no AvailableNode Ai")
	ErrNoGameServer   = errors.New("no AvailableNode Game")
	ErrSeasonUnStart  = errors.New("season not open")
	ErrAIHeroConfig   = errors.New("AI hero config error")
)

type GameManager struct {
	app   *appframe.Application
	games map[string]*Game

	idMgr *manager.IDManager
}

func newGameManager(app *appframe.Application) *GameManager {
	gm := new(GameManager)
	gm.app = app
	gm.games = make(map[string]*Game)
	gm.idMgr = &manager.IDManager{}
	gm.idMgr.Init(app.ID())
	return gm
}

func (p *GameManager) newGameId() string {
	return fmt.Sprintf("%d", p.idMgr.GeneratePKID())
}

func (p *GameManager) createGame(gameMode gameconf.GameModeTyp, roomId uint32, aiSvrID uint32, playerCount int, gsType gameconf.GameStartTyp, userIds []uint64) (*Game, error) {
	gameServer, err := p.getAvailableNodeForGame()
	if err != nil {
		return nil, err
	}

	gameId := p.newGameId()

	g := &Game{
		gameId:        gameId,
		node:          gameServer,
		state:         GameStateInitial,
		gameMode:      gameMode,
		aiSvrID:       aiSvrID,
		createTime:    time.Now().Unix(),
		roomId:        roomId,
		playerCount:   playerCount,
		userIds:       userIds,
		gameStartType: gsType,
	}
	p.games[g.gameId] = g
	return g, nil
}

func (p *GameManager) findGame(gameId string) (*Game, bool) {
	g, ok := p.games[gameId]
	return g, ok
}

// 游戏服通知游戏结束时调用
func (p *GameManager) deleteGame(gameId string) {
	delete(p.games, gameId)
}

func (p *GameManager) getAvailableNodeForGame() (appframe.Server, error) {
	/*serverId, ok := p.getGameSvrID()
	if !ok {
		return nil, errors.New("avaliablegame")
	}
	return p.app.GetServer(serverId), nil*/
	return serverStateService.GetLoadableServer(sgs_herox.SvrTypeGame)
}

//TODO 负载控制及版本更新控制
func (p *GameManager) getGameSvrID() (uint32, bool) {
	serverIds := p.app.GetAvailableServerIDs(sgs_herox.SvrTypeGame)
	if len(serverIds) == 0 {
		return 0, false
	}

	return serverIds[gameutil.Rand()%len(serverIds)], true
}
