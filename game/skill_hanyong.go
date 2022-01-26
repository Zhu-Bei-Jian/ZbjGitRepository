package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

// 72	兀突骨	悍勇	被动技：主动攻击后，随机对一名敌方武将造成2点伤害

type SkillHanYong struct {
	HeroSkill
}

func (ss *SkillHanYong) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillHanYong",
		triggerTypes: []TriggerType{TriggerType_AfterAttack, TriggerType_AfterAttackCamp},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			var card *Card
			switch ad.(type) {
			case *ActionAttackCard:
				card = ad.(*ActionAttackCard).srcCard

			case *ActionAttackCamp:
				card = ad.(*ActionAttackCamp).srcCard

			default:
				return
			}
			if card == nil || card.GetSkillId() != ss.GetSkillId() {
				return
			}

			enemy := FindCardsByType(card, gamedef.SelectCardType_Enemy)
			if len(enemy) == 0 {
				return
			}
			randIndex := gameutil.Intn(len(enemy))
			t := enemy[randIndex]
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v 触发悍勇：主动攻击后，随机对一名敌方武将（%v）造成2点伤害", card.GetOwnInfo(), t.GetOwnInfo()))
			actionData := &ActionDataBase{}
			g.PostActData(actionData)

			actionData.PostActStream(func() {
				g.Send2All(&cmsg.SSyncUseSkill{
					Seat:         card.owner.seatId,
					SkillId:      card.skillId,
					Card:         card.ID(),
					TargetCards:  cardIds(t),
					Data:         nil,
					TargetSeatId: 0,
				})
			})
			actionData.PostActStream(func() {
				g.StartWaitingNoneFloat(1.0, nil)
			})

			actionData.PostActStream(func() {
				NewActionDamageCard(g, t, card, nil, 2, ss.GetSkillId()).DoDamage()
			})

		},
	}
	return []TriggerHandler{th}
}
