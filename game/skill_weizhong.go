package game

import (
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/game/core"
)

type SkillWeiZhong struct {
	HeroSkill
}

//威重：被动技：回合结束时，生命值-3
func (ss *SkillWeiZhong) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillWeiZhong",
		triggerTypes: []TriggerType{TriggerType_PhaseEnd},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {

			var cards []*Card
			for _, rows := range g.board.cells {
				for _, v := range rows {
					if !v.HasCard() {
						continue
					}
					if v.owner == g.GetCurrentPlayer() && v.skillId == ss.GetSkillId() && !v.isBack {
						cards = append(cards, v.Card)
					}
				}
			}
			if len(cards) == 1 {
				logrus.Infof("%v 触发被动技：回合结束时，生命值-%v", cards[0].GetOwnInfo(), ss.GetValue(1))
				NewActionDamageCard(g, cards[0], cards[0], nil, ss.GetValue(1), ss.GetSkillId()).DoDamage()
			}

		},
	}
	return []TriggerHandler{th}
}
