package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
)

//21	贾诩	帷幕	被动技：被攻击时，受到伤害-2
type SkillWeiMu struct {
	HeroSkill
}

func (ss *SkillWeiMu) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillWeiMu",
		triggerTypes: []TriggerType{TriggerType_BeAttackCard},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionAttackCard)
			if !ok {
				return
			}
			//检查被攻击者 是否是 拥有帷幕 的贾诩
			if ac.targetCard.skillId == -1 {
				return
			}
			if ac.targetCard.skillCfg.SkillID != ss.GetSkillId() {
				return
			}
			ac.targetExtraDamage -= 2 //贾诩作为 被攻击者 target， 受到的伤害-2
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v 触发被动技：%v", ac.targetCard.GetOwnInfo(), ss.skillCfg.Name))
		},
	}
	return []TriggerHandler{th}
}
