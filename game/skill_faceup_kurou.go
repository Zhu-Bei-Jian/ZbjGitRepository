package game

import "fmt"

//苦肉：翻牌技，对自己造成12点伤害。
type SkillKuRou struct {
	HeroSkill
}

func (s *SkillKuRou) OnFaceUp(card *Card) {
	//card := cell.Card
	card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v 发动苦肉：翻牌技，对自己造成12点伤害。", card.GetOwnInfo()))
	g := card.GetPlayer().game
	g.PostActData(s)
	s.PostActStream(func() {
		actDamage := NewActionDamageCard(g, card, card, nil, 12, s.GetSkillId())
		actDamage.DoDamage()
	})
}
