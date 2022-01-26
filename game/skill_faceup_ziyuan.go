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

	s.PostActStream(func() {
		//可以改变生命上限
		t := enemies[randIndex]
		oldHP := t.GetHP()
		t.AddHpMax(1)
		t.AddHP(1)
		SyncChangeHP(t, oldHP, t.GetHP(), card, s.GetSkillId())

		oldAtk := t.GetAttack()
		t.AddAttack(1)
		SyncChangeAttack(t, oldAtk, t.GetAttack(), card)
	})

}
