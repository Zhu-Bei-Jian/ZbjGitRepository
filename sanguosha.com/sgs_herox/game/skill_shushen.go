package game

import (
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

// 67	甘夫人	淑慎	被动技：回合结束，随机使一名友方武将生命值+2   Q: +2 是否改变血上限
type SkillShuShen struct {
	HeroSkill
}

func (ss *SkillShuShen) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillShuShen",
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
			if len(myCards) == 0 {
				return
			}
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
