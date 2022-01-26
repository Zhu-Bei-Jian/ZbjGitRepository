package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
)

//buff节命 使武将免疫下一次收到的伤害
type BuffJieMing struct {
	HeroBuff
}

func (ss *BuffJieMing) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "BuffJieMing",
		triggerTypes: []TriggerType{TriggerType_OnDamage},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionDamageCard)
			if !ok {
				return
			}

			card := ac.card
			if !card.HasBuff(ss.GetBuffId()) {
				return
			}

			ac.damage = 0
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v 触发buff：%v", ac.card.GetOwnInfo(), ss.cfg.Name))
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v 触发buff：%v", ac.card.GetOwnInfo(), ss.cfg.Name))
			StartBuffOnEffect(card, ss.GetBuffId())
		},
	}

	th2 := TriggerHandler{
		name:         "BuffJieMing",
		triggerTypes: []TriggerType{TriggerType_LoseSkill},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionDamageCard)
			if !ok {
				return
			}

			card := ac.card
			if !card.HasBuff(ss.GetBuffId()) {
				return
			}

			ac.damage = 0
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v 触发buff：%v", ac.card.GetOwnInfo(), ss.cfg.Name))

			StartBuffOnEffect(card, ss.GetBuffId())
		},
	}

	return []TriggerHandler{th, th2}
}

//TODO 目前每个触发的buff技能都要加这个，这个逻辑应该放在框架层
func StartBuffOnEffect(card *Card, buffId int32) {
	card.BuffManager.AddBuffUseCount(buffId, 1)
	buffIds := card.BuffManager.ExpireBuffInTimes()

	if len(buffIds) > 0 {
		StartLoseBuffs(card, buffIds)
	}

}
