package game

import (
	"fmt"
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
			ss.PostActStream(func() {

				//天命 改变 血上限
				card.AddHpMax(1)
				card.AddHP(1)
				SyncChangeHP(card, card.GetHP()-1, card.GetHP(), card, ss.GetSkillId())

				card.attack++
				SyncChangeAttack(card, card.attack-1, card.attack, card)

				g.GetCurrentPlayer().Log(fmt.Sprintf("触发被动技：%v", ss.skillCfg.Name))
			})

		},
	}
	return []TriggerHandler{th}
}
