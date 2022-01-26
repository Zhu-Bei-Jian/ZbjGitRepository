package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/gameutil"
)

//43	小乔	天香	被动技：受到攻击时，敌方随机一名武将受到同样的伤害

type SkillTianXiang struct {
	HeroSkill
}

func (ss *SkillTianXiang) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillTianXiang",
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
			var enemies []*Card
			for _, rows := range g.board.cells {
				for _, cell := range rows {
					if !cell.HasCard() {
						continue
					}
					cd := cell.Card
					if cd.owner != ac.targetCard.owner { //和 小乔敌对的卡牌
						enemies = append(enemies, cd)
					}
				}
			}
			if len(enemies) < 1 {
				return
			}
			damage := ac.targetExtraDamage + ac.srcCard.attack
			if damage <= 0 {
				return
			}
			if len(enemies) == 0 {
				return
			}
			index := gameutil.Intn(len(enemies))
			NewActionDamageCard(g, enemies[index], ac.targetCard, nil, damage, ss.GetSkillId()).DoDamage()
			g.GetCurrentPlayer().Log(fmt.Sprintf("触发被动技：%v", ss.skillCfg.Name))
		},
	}
	return []TriggerHandler{th}
}
