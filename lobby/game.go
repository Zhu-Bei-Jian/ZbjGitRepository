package lobby

import (
	"fmt"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"
	"time"
)

type GameState int32

const (
	GameStateInitial GameState = 1
	GameStateRunning GameState = 2
	GameStateEnd     GameState = 3
)

type Game struct {
	gameId        string
	node          appframe.Server
	state         GameState
	gameMode      gameconf.GameModeTyp
	aiSvrID       uint32
	gameStartType gameconf.GameStartTyp
	roomId        uint32

	createTime int64

	playerCount int
	userIds     []uint64
}

func (g *Game) reqQuit(userId uint64, cb func(error)) {
	g.node.ReqSugar(&smsg.LoGaReqUserQuit{
		Userid: userId,
	}, func(resp *smsg.LoGaRespUserQuit, err error) {
		if err != nil {
			cb(err)
			return
		}

		if resp.ErrCode != 0 {
			cb(fmt.Errorf("errCode:%d", resp.ErrCode))
			return
		}
		cb(nil)
	}, time.Second*10)
}
