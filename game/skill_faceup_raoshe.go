package game

//饶舌：翻牌技：对其余己方明将造成6点伤害
type SkillRaoShe struct {
	HeroSkill
}

func (ss *SkillRaoShe) PreUseSkill() {
	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})

	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})
}
func (s *SkillRaoShe) OnFaceUp(card *Card) {
	g := card.owner.game
	g.PostActData(s)

	myCards := FindAllMyOwnCards(card)
	for _, v := range myCards {
		if v.isBack {
			continue
		}
		if card.cell.Row == v.cell.Row && card.cell.Col == v.cell.Col {
			continue
		}
		cd := v
		s.PostActStream(func() {
			NewActionDamageCard(g, cd, card, nil, 6, s.GetSkillId()).DoDamage()
		})

	}

}
