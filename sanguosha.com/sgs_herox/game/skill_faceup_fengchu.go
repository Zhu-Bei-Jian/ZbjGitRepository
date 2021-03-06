package game

import (
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//FengChu：翻牌技，令己方军营免疫2回合伤害。
type SkillFengChu struct {
	HeroSkill
}

func (s *SkillFengChu) OnFaceUp(card *Card) {
	g := card.owner.game
	expireV := g.GetBuffCfg(s.GetBuffId0()).GetExpireV() * 2
	AddCampBuff(card, s.GetBuffId0(), gameconf.ExpireTyp_ETRound, expireV)
}

func AddCampBuff(card *Card, buffId int32, expireTyp gameconf.ExpireTyp, expireV int32) {
	p := card.owner
	g := p.game
	buffCfg, exist := g.config.Buff.Get(buffId)
	if !exist {
		return
	}
	p.Camp.BuffManager.Add(buffCfg, g.roundCount, expireTyp, expireV)
	SyncCampChangeBuff(p, card)
}
func DelCampBuff(card *Card, buffId int32) {
	p := card.owner
	p.Camp.BuffManager.Remove(buffId)
	SyncCampChangeBuff(p, card)
}
