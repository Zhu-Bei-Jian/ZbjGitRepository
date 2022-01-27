package game

import (
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/cmsg"
)

//55	张角	引雷	翻牌技：对其他武将随机造成5次3点伤害
type SkillYinLei struct {
	HeroSkill
}

func (ss *SkillYinLei) PreUseSkill() {

}
func (ss *SkillYinLei) OnFaceUp(card *Card) {
	g := card.owner.game
	g.PostActData(ss)
	ss.PostActStream(func() {

		for i := int32(0); i < ss.GetValue(1); i++ {

			ss.PostActStream(func() {
				var otherCards []*Card
				for _, rows := range g.board.cells {
					for _, cell := range rows {
						if !cell.HasCard() {
							continue
						}
						if cell == card.cell {
							continue
						}
						otherCards = append(otherCards, cell.Card)
					}
				}
				if len(otherCards) == 0 {
					return
				}

				t := otherCards[gameutil.Intn(len(otherCards))]
				ss.PostActStream(func() {
					g.Send2All(&cmsg.SSyncUseSkill{
						SkillId:     ss.GetSkillId(),
						TargetCards: []int32{t.ID()},
					})
				})
				ss.PostActStream(func() {
					g.StartWaitingNone(1, nil)
				})
				ss.PostActStream(func() {
					NewActionDamageCard(g, t, card, nil, ss.GetValue(2), ss.GetSkillId()).DoDamage()
				})
				ss.PostActStream(func() {
					g.StartWaitingNone(1, nil)
				})
			})

		}
	})

	ss.PostActStream(func() {
		g.TryEndGame()
	})

}
