package game

import "fmt"

//焚城：翻牌技：对己方军营造成5点伤害
type SkillFengCheng struct {
	HeroSkill
}

func (ss *SkillFengCheng) PreUseSkill() {
	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})

	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})
}

func (s *SkillFengCheng) OnFaceUp(card *Card) {
	card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v发动%v ,对己方军营造成5点伤害 ", s.card.GetOwnInfo(), s.card.skillCfg.Name))
	g := card.GetPlayer().game
	g.PostActData(s)
	s.PostActStream(func() {
		card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf(" %v发动 焚城：翻牌技：对己方军营造成%v点伤害", card.GetOwnInfo(), s.GetValue(1)))
		actDamage := NewActionDamageCamp(g, card.owner, card, s.GetValue(1))
		actDamage.DoDamage()
	})
}
