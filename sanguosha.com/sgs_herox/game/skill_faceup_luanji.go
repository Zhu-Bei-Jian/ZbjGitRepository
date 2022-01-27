package game

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

//乱击：翻牌技，对其他武将造成2点伤害，若自身在交战区，对其他武将造成4点伤害。
type SkillLuanJi struct {
	HeroSkill
}

func (ss *SkillLuanJi) PreUseSkill() {
	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})

	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})
}
func (s *SkillLuanJi) OnFaceUp(card *Card) {
	p := card.GetPlayer()
	g := p.game
	damageValue := s.GetValue(1)
	if IsInWarZone(card) {
		damageValue = s.GetValue(2)
	}
	g.PostActData(s)
	card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v开始发动乱击,对其他武将造成%d点伤害", card.GetOwnInfo(), damageValue))
	for i := int32(0); i < 3; i++ {
		ti := i
		for j := int32(0); j < 3; j++ {
			tj := j
			if IsPositionEqual(card.cell.Position, Position{Row: i, Col: j}) {
				continue
			}
			if !g.board.cells[i][j].HasCard() {
				continue
			}
			s.PostActStream(func() {
				NewActionDamageCard(g, g.board.cells[ti][tj].Card, card, nil, damageValue, s.GetSkillId()).DoDamage()
				logrus.Printf("%v受到乱击 伤害%v", g.board.cells[ti][tj].GetOwnInfo(), damageValue)

			})
		}
	}

}

func IsPositionEqual(p1 Position, p2 Position) bool {
	return p1.Row == p2.Row && p1.Col == p2.Col
}
