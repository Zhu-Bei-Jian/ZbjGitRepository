package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
)

type SkillGangLie struct {
	HeroSkill
}

func (ss *SkillGangLie) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillGangLie",
		triggerTypes: []TriggerType{TriggerType_BeAttackCard},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionAttackCard)
			if !ok {
				return
			}
			if ac.targetCard.skillId == -1 {
				return
			}
			if ac.targetCard.skillCfg.SkillID != ss.GetSkillId() {
				return
			}
			ac.extraRetDamage += ss.GetValue(1)

			g.GetCurrentPlayer().Log(fmt.Sprintf("%v触发被动技 %v", ac.targetCard.GetOwnInfo(), ac.targetCard.skillCfg.Name))
		},
	}
	return []TriggerHandler{th}
}
