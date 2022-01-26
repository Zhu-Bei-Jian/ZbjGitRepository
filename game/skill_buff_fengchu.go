package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
)

//凤雏 已方军营免疫两回合伤害
type BuffFengChu struct {
	HeroBuff
}

func (ss *BuffFengChu) TriggerHandler() []TriggerHandler {
	th1 := TriggerHandler{
		name:         "BuffFengChu",
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
			g.GetCurrentPlayer().Log(fmt.Sprintf("玩家%v的庞统被动技凤雏生效。 %v 无法攻击敌方军营", g.players[(ac.card.owner.seatId+1)%2].user.userBrief.Nickname, ac.card.GetOwnInfo()))

		},
	}

	th2 := TriggerHandler{
		name:         "BuffFengChu",
		triggerTypes: []TriggerType{TriggerType_MakeDamageCamp},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionDamageCamp)
			if !ok {
				return
			}
			if ac.player.HasBuff(ss.GetBuffId()) {
				ac.extraDamageToCamp -= INF
			}

		},
	}

	return []TriggerHandler{th1, th2}
}
