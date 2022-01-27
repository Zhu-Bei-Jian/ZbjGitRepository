package game

//陷阵：翻牌技：对交战区整排武将造成4点伤害
type SkillXianZhen struct {
	HeroSkill
}

func (ss *SkillXianZhen) PreUseSkill() {
	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})

	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})
}

func (ss *SkillXianZhen) OnFaceUp(card *Card) {
	var warRow = 1 //交战区
	g := card.owner.game
	g.PostActData(ss)
	ss.PostActStream(func() {
		for i := 0; i < 3; i++ {
			ti := i
			if cd, ok := card.owner.game.board.cells[warRow][ti].GetCard(); ok {
				ss.PostActStream(func() {
					NewActionDamageCard(cd.owner.game, cd, card, nil, ss.GetValue(1), ss.GetSkillId()).DoDamage()
				})
			}
		}
	})
	ss.PostActStream(func() {
		card.owner.game.TryEndGame()
	})

}
