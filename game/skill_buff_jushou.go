package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
)

type BuffJuShou struct {
	HeroBuff
}

func (ss *BuffJuShou) TriggerHandler() []TriggerHandler {
	th1 := TriggerHandler{
		name:         "BuffJuShou",
		triggerTypes: []TriggerType{TriggerType_CheckAttackCamp},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionCheck)
			if !ok {
				return
			}
			if !g.players[(ac.card.owner.seatId+1)%2].Camp.HasBuff(ss.GetBuffId()) {
				return
			}
			ac.canAttackCamp = false
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v 曹仁据守生效。 %v 无法攻击敌方军营", g.players[(ac.card.owner.seatId+1)%2].user.userBrief.Nickname, ac.card.GetOwnInfo()))

		},
	}

	return []TriggerHandler{th1}
}
