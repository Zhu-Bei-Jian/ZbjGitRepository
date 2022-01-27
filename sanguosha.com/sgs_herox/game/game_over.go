package game

import (
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/proto/cmsg"
)

func (g *Game) CheckIsOver() bool {
	players := g.base.GetAlivePlayers()
	if len(players) == 1 {
		g.winSeat = players[0].GetSeatID()
		return true
	}

	players = g.base.CanActionPlayers()
	if len(players) == 1 {
		g.winSeat = players[0].GetSeatID()
		return true
	}
	if len(players) == 0 {
		if len(g.base.players) != 2 {
			logrus.Fatal("对局玩家数量不为2")
			return true
		}
		if g.base.players[0].hp == g.base.players[1].hp {
			g.winSeat = -1
		} else if g.base.players[0].hp > g.base.players[1].hp {
			g.winSeat = g.base.players[0].GetSeatID()
		} else {
			g.winSeat = g.base.players[1].GetSeatID()
		}

		return true
	}

	return false
}

func (g *Game) onGameOver() {
	g.base.Send2All(&cmsg.SNoticeGameOver{WinSeat: g.winSeat})

}
