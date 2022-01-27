package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
)

//缔盟：翻牌技，指定一己方武将获得技能：该武将受到反击时，伤害-3
type BuffDiMeng struct {
	HeroBuff
}

func (ss *BuffDiMeng) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "BuffDiMeng",
		triggerTypes: []TriggerType{TriggerType_BeRetAttack},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionAttackCard)
			if !ok {
				return
			}

			card := ac.srcCard
			if !card.HasBuff(ss.GetBuffId()) {
				return
			}
			for _, v := range card.buffs {
				if v.buffCfg.BuffID == ss.GetBuffId() {
					ac.extraRetDamage -= ss.GetCfg().BuffDamage
				}
			}

			g.GetCurrentPlayer().Log(fmt.Sprintf("触发buff：%v", ss.cfg.Name))
		},
	}
	return []TriggerHandler{th}
}
