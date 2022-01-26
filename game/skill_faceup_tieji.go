package game

import "sanguosha.com/sgs_herox/proto/gameconf"

//铁骑：翻牌技，本回合攻击力+10
type SkillTieJi struct {
	HeroSkill
}

func (s *SkillTieJi) OnFaceUp(card *Card) {
	StartGetBuff(card, s.GetBuffId0(), gameconf.ExpireTyp_ETRound, 1, card)
}
