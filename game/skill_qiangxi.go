package game

import (
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/game/core"
)

// 强袭：被动技，你无法攻击敌方军营。
type SkillQiangXi struct {
	HeroSkill
}

func (ss *SkillQiangXi) TriggerHandler() []TriggerHandler {
	th1 := TriggerHandler{
		name:         "SkillQiangXi",
		triggerTypes: []TriggerType{TriggerType_CheckAttackCamp},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionCheck)
			if !ok {
				return
			}
			if ac.card.isBack || !ac.card.HasSkill(ss.GetSkillId()) {
				return
			}
			ac.canAttackCamp = false
			logrus.Info("%v强袭生效，无法攻击敌方军营", ac.card.GetOwnInfo())
		},
	}

	return []TriggerHandler{th1}
}
