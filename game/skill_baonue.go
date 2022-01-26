package game

import (
	"sanguosha.com/sgs_herox/game/core"
)

// 董卓-暴虐  重装；（重装：上场、攻击、移动需要花费2点行动点）（重装的武将只能正面放置）
type SkillBaoNue struct {
	HeroSkill
}

func (ss *SkillBaoNue) TriggerHandler() []TriggerHandler {
	th1 := TriggerHandler{
		name:         "SkillBaoNue",
		triggerTypes: []TriggerType{TriggerType_CheckAP_PlaceCard},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionCheckPlaceCard)
			if !ok {
				return
			}

			if ac.skillCfg.SkillID != ss.GetSkillId() {
				return
			}
			ac.ap = g.config.HeavyCost
		},
	}

	th2 := TriggerHandler{
		name:         "SkillBaoNue",
		triggerTypes: []TriggerType{TriggerType_CheckHeavyAP},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionCheck)
			if !ok {
				return
			}

			if ac.card.skillId != ss.GetSkillId() {
				return
			}
			ac.ap[MOVE_CARD] = g.config.HeavyCost
			ac.ap[ATTACK_CAMP] = g.config.HeavyCost
			ac.ap[ATTACK_CARD] = g.config.HeavyCost

		},
	}

	return []TriggerHandler{th1, th2}
}
