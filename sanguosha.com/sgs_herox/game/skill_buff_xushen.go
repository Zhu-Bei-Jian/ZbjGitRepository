package game

import (
	"sanguosha.com/sgs_herox/game/core"
)

//许身 翻牌技：使一名其他武将获得一个技能：你的回合结束，回复2点生命值
type BuffXuShen struct {
	HeroBuff
}

func (ss *BuffXuShen) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "BuffXuShen",
		triggerTypes: []TriggerType{TriggerType_PhaseEnd},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			cfg := g.GetBuffCfg(ss.GetBuffId())
			for _, rows := range g.board.cells {
				for _, v := range rows {
					if !v.HasCard() || !v.HasBuff(ss.GetBuffId()) || v.owner.seatId != g.GetCurrentPlayer().seatId {
						continue
					}
					oldHP := v.GetHP()
					var value int32 = 0
					for _, bf := range v.Card.buffs {
						if bf.buffCfg.BuffID == ss.GetBuffId() {
							value += cfg.GetBuffHP() * bf.buffCount
						}
					}
					v.AddHP(value)
					SyncChangeHP(v.Card, oldHP, v.GetHP(), nil, ss.GetBuffId())

				}
			}

		},
	}

	return []TriggerHandler{th}
}
