package game

import "sanguosha.com/sgs_herox/proto/gameconf"

type SkillWuSheng struct {
	HeroSkill
}

//关羽	武圣	翻牌技，获得自身攻击力的生命值。
func (s *SkillWuSheng) OnFaceUp(card *Card) {
	//武圣 会改变 血量上限
	card.owner.game.PostActData(s)
	s.PostActStream(func() {
		StartGetBuff(card, s.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
	})
	s.PostActStream(func() {
		oldHp := card.GetHP()
		card.AddHpMax(card.attack)
		card.AddHP(card.attack)
		SyncChangeHP(card, oldHp, card.GetHP(), card, s.GetSkillId())
	})

}
