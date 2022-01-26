package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
)

// 53	张飞	咆哮	被动技，无攻击次数限制且和武将对战伤害始终-2

type SkillPaoXiao struct {
	HeroSkill
}

func (ss *SkillPaoXiao) TriggerHandler() []TriggerHandler {
	th1 := TriggerHandler{
		name:         "SkillPaoXiao",
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
			ac.targetExtraDamage -= 2
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v受到攻击，被动技咆哮生效，受到的伤害-2", ac.targetCard.GetOwnInfo()))
		},
	}
	th2 := TriggerHandler{
		name:         "SkillPaoXiao",
		triggerTypes: []TriggerType{TriggerType_CheckAttackCntInTurn},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionCheck)
			if !ok {
				return
			}

			if !ac.card.HasSkill(ss.GetSkillId()) {
				return
			}

			setAttackCountAndNotify(ac.card, INF)
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v咆哮，攻击次数无限制", ac.card.GetOwnInfo()))
		},
	}
	th3 := TriggerHandler{
		name:         "SkillPaoXiao",
		triggerTypes: []TriggerType{TriggerType_AttackCard},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionAttackCard)
			if !ok {
				return
			}
			if ac.srcCard.isBack {
				return
			}
			if ac.srcCard.skillId != ss.GetSkillId() {
				return
			}
			ac.extraRetDamage -= 2
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v 受到的反击伤害-2", ac.srcCard.GetOwnInfo()))
		},
	}
	return []TriggerHandler{th1, th2, th3}
}
