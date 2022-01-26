package game

import (
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//50	于禁	毅重	翻牌技：该武将在我方置牌区和交战区翻开，攻/血变成1/7
type SkillYiZhong struct {
	HeroSkill
}

func (ss *SkillYiZhong) OnFaceUp(card *Card) {

	if IsInMyZone(card) || IsInWarZone(card) {
		card.owner.game.PostActData(ss)
		ss.PostActStream(func() {
			StartGetBuff(card, ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
		})
		ss.PostActStream(func() {
			oldAt := card.attack
			oldHp := card.GetHP()

			if oldHp > 7 { //血量上限减少
				card.SubHpMax(oldHp - 7)
			} else {
				card.AddHpMax(7 - oldHp)
				card.AddHP(7 - oldHp)
			}
			card.attack = 1

			SyncChangeHP(card, oldHp, card.GetHP(), card, ss.GetSkillId())
			SyncChangeAttack(card, oldAt, card.attack, card)
		})
	}
}
