package game

import "sanguosha.com/sgs_herox/proto/gameconf"

// 29	马谡	失守	翻牌技：3回合内，己方军营每次受到伤害额外增加3点
type SkillShiShou struct {
	HeroSkill
}

func (s *SkillShiShou) OnFaceUp(card *Card) {
	g := card.owner.game
	expireV := g.GetBuffCfg(s.GetBuffId0()).GetExpireV() * 2
	AddCampBuff(card, s.GetBuffId0(), gameconf.ExpireTyp_ETRound, expireV)
}
