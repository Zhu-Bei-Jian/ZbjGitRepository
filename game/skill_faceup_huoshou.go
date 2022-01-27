package game

import (
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//30	孟获	祸首	翻牌技：我方随机一名武将攻击力-2。
type SkillHuoShou struct {
	HeroSkill
}

func (s *SkillHuoShou) OnFaceUp(card *Card) {
	myCards := FindAllMyOwnCards(card)
	if len(myCards) == 0 {
		return
	}
	card.owner.game.PostActData(s)
	randIndex := gameutil.Intn(len(myCards))
	s.PostActStream(func() {
		StartGetBuff(myCards[randIndex], s.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
	})

}
