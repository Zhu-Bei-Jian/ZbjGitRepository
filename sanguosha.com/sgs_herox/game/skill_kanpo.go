package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
)

//62	诸葛亮	看破	被动技，己方武将攻击军营不会受到伤害。

type SkillKanPo struct {
	HeroSkill
}

func (ss *SkillKanPo) TriggerHandler() []TriggerHandler {
	th1 := TriggerHandler{
		name:         "SkillKanPo",
		triggerTypes: []TriggerType{TriggerType_BeRetAttackByCamp},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionAttackCamp)
			if !ok {
				return
			}
			cards := FindAllMyOwnCards(ac.srcCard)
			hasKanPo := false
			for _, v := range cards {
				if v.HasSkill(ss.GetSkillId()) {
					hasKanPo = true
					break
				}
			}
			if !hasKanPo {
				return
			}
			ac.retDamage -= INF
			g.GetCurrentPlayer().Log(fmt.Sprintf("诸葛亮 被动技-看破 ：己方武将攻击军营不会受到伤害。此次攻击为%v 发动,不会受到敌方军营的反击伤害", ac.srcCard.GetOwnInfo()))
		},
	}

	return []TriggerHandler{th1}
}

// 翻面
