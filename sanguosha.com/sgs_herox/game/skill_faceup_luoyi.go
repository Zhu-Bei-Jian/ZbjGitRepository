package game

import "sanguosha.com/sgs_herox/proto/gameconf"

//裸衣：翻牌技，攻击力增加2，若此时自身在交战区或敌方置牌区，攻击力增加4。
type SkillLuoYi struct {
	HeroSkill
}

func (ss *SkillLuoYi) OnFaceUp(card *Card) {
	if IsInWarZone(card) || IsInEnemyZone(card) {
		StartGetBuff(card, ss.skillCfg.Buffs[1], gameconf.ExpireTyp_ETInvalid, 0, card)
	} else {
		StartGetBuff(card, ss.skillCfg.Buffs[0], gameconf.ExpireTyp_ETInvalid, 0, card)
	}
}
