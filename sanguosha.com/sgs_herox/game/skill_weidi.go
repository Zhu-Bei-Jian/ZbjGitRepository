package game

import (
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//伪帝：被动技，当你处于正面时，你的回合结束获得+1/+1
type SkillWeiDi struct {
	HeroSkill
}

func (ss *SkillWeiDi) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillWeiDi",
		triggerTypes: []TriggerType{TriggerType_PhaseEnd},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			p := g.GetCurrentPlayer()
			cards := g.skillCards(p, ss.GetSkillId())
			for _, card := range cards {
				if card.isBack {
					continue
				}
				if card.owner.seatId != p.seatId {
					continue
				}

				g.PostActData(ss)
				ss.PostActStream(func() {
					StartGetBuff(card, ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
				})

			}

		},
	}
	return []TriggerHandler{th}
}
