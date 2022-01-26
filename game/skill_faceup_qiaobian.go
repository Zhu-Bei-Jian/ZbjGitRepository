package game

import "sanguosha.com/sgs_herox/proto/gameconf"

type SkillQiaoBian struct {
	HeroSkill
}

func (ss *SkillQiaoBian) PreUseSkill() {
	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})

	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})
}

//54	张郃	巧变	翻牌技：当前回合，你无攻击距离限制（军营除外）
func (ss *SkillQiaoBian) OnFaceUp(card *Card) {
	StartGetBuff(card, ss.GetBuffId0(), gameconf.ExpireTyp_ETRound, 1, card)
}
