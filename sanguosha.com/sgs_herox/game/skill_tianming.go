package game

import (
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

type SkillTianMing struct {
	HeroSkill
}

func (ss *SkillTianMing) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillTianMing",
		triggerTypes: []TriggerType{TriggerType_AfterPosChange},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			var card *Card
			switch t := ad.(type) {
			case *ActionMove:
				card = t.card
			case *ActionFaceUpCard:
				card = t.card
			case *ActionSilence:
				card = t.card
			case *ActionDataBase:
				for _, v := range params {
					cd, ok := v.(*Card)
					if ok && cd.skillId == ss.GetSkillId() {
						card = cd
					}
				}
			case *ActionActAttack:
				card = t.srcCard
			default:
				return
			}
			if card == nil {
				return
			}
			if card.isBack {
				return
			}
			if !card.HasSkill(ss.GetSkillId()) {
				return
			}
			card.owner.game.PostActData(ss)
			ss.PostActStream(func() {
				StartGetBuff(card, ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
			})

		},
	}
	return []TriggerHandler{th}
}
