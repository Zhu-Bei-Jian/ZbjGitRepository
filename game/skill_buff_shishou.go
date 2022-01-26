package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
)

//buff失守  军营每次受到伤害增加3点
type BuffShiShou struct {
	HeroBuff
}

func (ss *BuffShiShou) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "BuffShiShou",
		triggerTypes: []TriggerType{TriggerType_MakeDamageCamp},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionDamageCamp)
			if !ok {
				return
			}

			if !ac.player.Camp.HasBuff(ss.GetBuffId()) {
				return
			}
			for _, v := range ac.player.Camp.buffs { //每有一个 失守buff，军营受伤+3
				if v.buffCfg.BuffID == ss.GetBuffId() {
					ac.damage += 3
				}
			}

			g.GetCurrentPlayer().Log(fmt.Sprintf("触发buff:%v", ss.cfg.Name))
		},
	}

	return []TriggerHandler{th}
}
