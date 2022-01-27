package game

import "fmt"

//苦肉：翻牌技，对自己造成12点伤害。
type SkillKuRou struct {
	HeroSkill
}

func (s *SkillKuRou) OnFaceUp(card *Card) {
	g := card.GetPlayer().game
	g.GetCurrentPlayer().Log(fmt.Sprintf("%v 发动 %v", card.GetOwnInfo(), s.skillCfg.GetDesc()))
	g.PostActData(s)
	s.PostActStream(func() {
		actDamage := NewActionDamageCard(g, card, card, nil, s.GetValue(1), s.GetSkillId())
		actDamage.DoDamage()
	})
}
