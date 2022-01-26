package game

import "sanguosha.com/sgs_herox/proto/gameconf"

// 29	马谡	失守	翻牌技：3回合内，己方军营每次受到伤害额外增加3点
type SkillShiShou struct {
	HeroSkill
}

func (s *SkillShiShou) OnFaceUp(card *Card) {
	var expireV int32 = 3 * 2
	if card.owner.game.roundCount%2 == 0 {
		expireV--
	}
	AddCampBuff(card, s.GetBuffId0(), gameconf.ExpireTyp_ETRound, expireV)
}
