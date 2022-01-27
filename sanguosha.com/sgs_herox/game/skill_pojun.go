package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
)

//68	徐盛	破军	被动技：攻击军营，造成双倍伤害
type SkillPoJun struct {
	HeroSkill
}

func (ss *SkillPoJun) PreUseSkill() {
	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})
	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})

}
func (ss *SkillPoJun) TriggerHandler() []TriggerHandler {
	th1 := TriggerHandler{
		name:         "SkillPoJun",
		triggerTypes: []TriggerType{TriggerType_AttackCamp},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionAttackCamp)
			if !ok {
				return
			}
			if ac.srcCard.GetSkillId() != ss.GetSkillId() {
				return
			}
			ac.multi = 2
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v 破军 被动技：攻击军营，造成双倍伤害", ac.srcCard.GetOwnInfo()))
		},
	}

	return []TriggerHandler{th1}
}
