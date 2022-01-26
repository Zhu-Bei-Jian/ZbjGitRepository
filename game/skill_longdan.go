package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

type SkillLongDang struct {
	HeroSkill
}

func (ss *SkillLongDang) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillLongDang",
		triggerTypes: []TriggerType{TriggerType_AfterAttack},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			var card *Card
			switch ad.(type) {
			case *ActionAttackCamp:
				//ac, ok := ad.(*ActionAttackCamp)
				//if !ok {
				//	return
				//}
				//card = ac.srcCard
				return
			case *ActionAttackCard:
				ac, ok := ad.(*ActionAttackCard)
				if !ok {
					return
				}
				card = ac.srcCard
			default:
				return
			}

			if card.isBack {
				return
			}
			if card.GetSkillId() != ss.GetSkillId() {
				return
			}
			g.PostActData(ss)
			ss.PostActStream(func() {
				//if card.HasBuff(ss.GetSkillId()) {
				//	return
				//}
				StartGetBuff(card, ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
			})

			ss.PostActStream(func() {
				oldAt := card.attack
				card.attack *= 2
				SyncChangeAttack(card, oldAt, card.attack, card)
				g.GetCurrentPlayer().Log(fmt.Sprintf("%v触发被动技 %v", card.GetOwnInfo(), card.skillCfg.Name))
			})

		},
	}
	return []TriggerHandler{th}
}
