package game

import (
	"fmt"
	"github.com/sirupsen/logrus"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/smsg"
	"time"
)

func (g *Game) onPrepare() {

}

func (g *Game) onStart() {
	AhoCorasick.SetMap(map[string]string{
		"zhubeijian": "15558352672",
		"panshi":     "17602222208",
		"zhangjie":   "17764549625",
	})
}

func (g *Game) onInitCards() {
	for _, v := range g.base.players {
		//paiku
		p := v
		var heroIds []int32
		cardPoolCount := g.base.config.CardPoolCount
		if p.user.userBrief.Nickname == "hi_1" || p.user.userBrief.Nickname == "hi_2" || p.user.userBrief.Nickname == "hi_beijian" { //测试账号 走固定卡牌，便于debug
			heroIds = g.base.GetHeroIdsByNames("曹仁", "马谡", "黄月英", "刘禅", "司马徽")
			//shuffleFisherYates(heroIds)
			for _, heroId := range heroIds {
				card := g.genPoolCard(heroId)
				p.HandCardPool.AddCard(card)
				if p.HandCardPool.Count() > int(cardPoolCount) {
					break
				}
			}
		} else { //正常情况 牌库从entity 加载
			g.base.ReqEntity(&smsg.ReqUserData{
				Userid: p.user.userId,
			}, func(respData *smsg.RespUserData, err error) {
				if err != nil {
					logrus.Warnf("返回类型错误")
					return
				}

				if respData.CardGroup == nil || respData.CardGroup.GroupId < 0 { //如果玩家从未编辑过牌库 或者指定随机卡组，那么随机一套卡组
					heroIds = p.game.config.Hero.HeroIds()
				} else { //将玩家最近指定使用的套牌作为游戏内套牌
					heroIds = respData.CardGroup.HeroIds
				}
				shuffleFisherYates(heroIds)
				for _, heroId := range heroIds {
					card := g.genPoolCard(heroId)
					p.HandCardPool.AddCard(card)
					if p.HandCardPool.Count() > int(cardPoolCount) {
						break
					}
				}
			}, time.Second*30)
		}

	}
}

func (g *Game) genPoolCard(heroId int32) *gamedef.PoolCard {
	return &gamedef.PoolCard{
		CardId: g.newCardId(),
		HeroId: heroId,
	}
}

func (g *Game) newCardId() int32 {
	g.cardIdIndex++
	return g.cardIdIndex
}

func (g *Game) onPhaseBegin() {
	game := g.base
	ad := &ActionDataBase{}
	game.PostActData(ad)

	ad.PostActStream(func() {
		nextPlayer := game.GetNextAlivePlayer(game.curPlayer)
		game.curPlayer = nextPlayer.GetSeatID() //回合开始，更新当前玩家
		game.roundCount++
		game.setPhaseAndSync(gamedef.GamePhase_PhaseBegin)
	})

	ad.PostActStream(func() {
		game.GetCurrentPlayer().onEnterPhase(gamedef.GamePhase_PhaseBegin)
		game.board.onEnterPhase(gamedef.GamePhase_PhaseBegin)
	})

	ad.PostActStream(func() {
		//触发 在回合开始生效的技能
		game.GetTriggerMgr().Trigger(TriggerType_PhaseBegin, &ActionDataBase{})
	})
}

func (g *Game) onPhaseDraw() {

	if g.base.roundCount == 1 {
		for _, p := range g.base.players {
			g.base.DrawHandCards(p)
		}
	} else {
		p := g.base.GetCurrentPlayer()
		g.base.DrawHandCards(p)
	}

}

func (g *Game) onPhaseMain() {

}

func (g *Game) onPhaseEnd() {

	ad := &ActionDataBase{}
	g.base.PostActData(ad)

	ad.PostActStream(func() {
		g.base.GetTriggerMgr().Trigger(TriggerType_PhaseEnd, &ActionDataBase{})
	})
	ad.PostActStream(func() {
		for _, rows := range g.base.board.cells {
			for _, cell := range rows {
				if !cell.HasCard() {
					continue
				}
				cell.canTurnUp = true
			}
		}
	})
	ad.PostActStream(func() {
		if g.base.roundCount%6 == 0 {
			if uidToLog[g.base.roomId] == "" {
				return
			}
			info := fmt.Sprintf("%v和%v的对局 roomId: %v ,gameUuid:%v ,roomNo:%v：第%v至第%v回合操作记录:", g.base.players[0].user.userBrief.Nickname, g.base.players[1].user.userBrief.Nickname, g.base.roomId, g.base.gameUUID, g.base.roomNO, (g.base.roundCount+1)/2-2, (g.base.roundCount+1)/2) + uidToLog[g.base.roomId]

			g.base.DingSendMsg(info)
			//logrus.Info(info)
			uidToLog[g.base.roomId] = ""
		}
	})

}

func (g *Game) onEnd() {
	g.onGameOver()
}
