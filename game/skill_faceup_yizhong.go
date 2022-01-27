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
		g := card.owner.game
		g.PostActData(ss)
		cfg := g.GetBuffCfg(ss.GetBuffId0())
		ss.PostActStream(func() {
			StartGetBuff(card, ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
		})
		ss.PostActStream(func() {
			card.SetBodyAndNotify(cfg.GetBuffAttack(), cfg.GetBuffHP())
		})
	}
}

func (p *Card) SetBodyAndNotify(attack int32, hp int32) {
	oldAt := p.attack
	oldHp := p.GetHP()

	if oldHp > hp { //血量上限减少
		p.SubHpMax(oldHp - hp)
	} else {
		p.AddHpMax(hp - oldHp)
		p.AddHP(hp - oldHp)
	}
	p.attack = attack

	SyncChangeHP(p, oldHp, p.GetHP(), p, p.GetSkillId())
	SyncChangeAttack(p, oldAt, p.attack, p)
}
