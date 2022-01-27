package game

import (
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//英魂：被动技，在你的回合结束时，随机令一名己方武将获得+1+1
type SkillYingHun struct {
	HeroSkill
}

func (ss *SkillYingHun) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillYingHun",
		triggerTypes: []TriggerType{TriggerType_PhaseEnd},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			var cards []*Card
			for _, rows := range g.board.cells {
				for _, cell := range rows {
					if cell.HasCard() && cell.Card.skillId == ss.GetSkillId() && cell.Card.owner.seatId == g.GetCurrentPlayer().seatId {
						cards = append(cards, cell.Card)
					}
				}
			}

			if cards == nil {
				return
			}
			myCards := FindAllMyOwnCards(cards[0])
			randIndex := gameutil.Intn(len(myCards))
			card := myCards[randIndex]

			g.PostActData(ss)
			ss.PostActStream(func() {
				StartGetBuff(card, ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, cards[0])
			})

		},
	}
	return []TriggerHandler{th}
}
