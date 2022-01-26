package game

import (
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//机巧：翻牌技，令敌方1名随机武将的攻击力-2，若在交战区，令敌方两名随机武将攻击力-2

type SkillJiQiao struct {
	HeroSkill
}

func (ss *SkillJiQiao) OnFaceUp(card *Card) {

	enemies := FindAllEnemyCards(card)
	numEnemy := len(enemies)
	if numEnemy == 0 {
		return
	}
	card.owner.game.PostActData(ss)
	if !IsInWarZone(card) {
		randIndex := gameutil.Intn(numEnemy) // 线程安全的随机数
		oldAttack := enemies[randIndex].attack
		enemies[randIndex].SubAttack(2)
		ss.PostActStream(func() {
			StartGetBuff(enemies[randIndex], ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
		})
		ss.PostActStream(func() {
			SyncChangeAttack(enemies[randIndex], oldAttack, enemies[randIndex].attack, card)
		})

	} else {
		if numEnemy == 1 {
			randIndex := gameutil.Intn(numEnemy) // 线程安全的随机数
			oldAttack := enemies[randIndex].attack
			enemies[randIndex].SubAttack(2)
			ss.PostActStream(func() {
				StartGetBuff(enemies[randIndex], ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
			})
			ss.PostActStream(func() {
				SyncChangeAttack(enemies[randIndex], oldAttack, enemies[randIndex].attack, card)
			})
		} else if numEnemy >= 2 { //处于交战区，且敌方武将数量>=2
			randIndex1 := gameutil.Intn(numEnemy)
			oldAttack := enemies[randIndex1].attack
			enemies[randIndex1].SubAttack(2)
			ss.PostActStream(func() {
				StartGetBuff(enemies[randIndex1], ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
			})
			ss.PostActStream(func() {
				SyncChangeAttack(enemies[randIndex1], oldAttack, enemies[randIndex1].attack, card)
			})

			randIndex2 := gameutil.Intn(numEnemy)
			if randIndex1 == randIndex2 { //保证第二个武将与第一个武将不是同一个
				randIndex2 = (randIndex2 + 1) % numEnemy
			}
			oldAttack = enemies[randIndex2].attack
			enemies[randIndex2].SubAttack(2)
			ss.PostActStream(func() {
				StartGetBuff(enemies[randIndex2], ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)

			})
			ss.PostActStream(func() {
				SyncChangeAttack(enemies[randIndex2], oldAttack, enemies[randIndex2].attack, card)
			})
		}
	}

}
