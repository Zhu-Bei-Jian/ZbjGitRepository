package game

import (
	"math/rand"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//资援：翻牌技：敌方随机一名武将+1/+1

type SkillZiYuan struct {
	HeroSkill
}

func (s *SkillZiYuan) OnFaceUp(card *Card) {
	enemies := FindAllEnemyCards(card)
	if len(enemies) == 0 {
		return
	}

	randIndex := rand.Int() % len(enemies)
	card.owner.game.PostActData(s)
	s.PostActStream(func() {
		StartGetBuff(enemies[randIndex], s.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
	})

}
