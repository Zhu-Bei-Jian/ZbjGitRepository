package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/cmsg"
)

// 73	周泰	不屈	被动技：阵亡时，对敌方军营造成4点伤害

type SkillBuQu struct {
	HeroSkill
}

func (ss *SkillBuQu) TriggerHandler() []TriggerHandler {
	th1 := TriggerHandler{
		name:         "SkillBuQu",
		triggerTypes: []TriggerType{TriggerType_Die},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			if len(params) != 1 {
				return
			}
			card, ok := params[0].(*Card)
			if !ok {
				return
			}
			if !card.HasSkill(ss.GetSkillId()) {
				return
			}
			enemySeatId := (card.owner.seatId + 1) % 2
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v死亡，触发不屈：阵亡时，对敌方军营造成4点伤害", card.GetOwnInfo()))

			actionData := &ActionDataBase{}
			g.PostActData(actionData)
			actionData.PostActStream(func() {
				g.Send2All(&cmsg.SSyncUseSkill{
					Seat:         card.owner.seatId,
					SkillId:      card.skillId,
					Card:         card.ID(),
					TargetSeatId: enemySeatId,
				})
			})
			actionData.PostActStream(func() {
				g.StartWaitingNoneFloat(1.0, nil)
			})

			actionData.PostActStream(func() {
				NewActionDamageCamp(g, g.players[enemySeatId], card, 4).DoDamage()
			})

		},
	}

	return []TriggerHandler{th1}
}
