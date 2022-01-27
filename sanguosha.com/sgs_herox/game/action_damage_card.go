package game

import (
	"fmt"
)

type ActionDamageCard struct {
	ActionDataBase

	card *Card //受伤的卡牌
	//卡牌受伤 ，伤害来源分为 卡牌和军营
	srcCard   *Card   //伤害来源为卡牌
	srcPlayer *Player //伤害来源为军营
	skillId   int32   //造成伤害的技能
	damage    int32
}

func NewActionDamageCard(game *GameBase, card *Card, srcCard *Card, srcPlayer *Player, damage int32, skillId int32) *ActionDamageCard {
	ac := &ActionDamageCard{}
	ac.game = game
	ac.srcPlayer = srcPlayer
	ac.srcCard = srcCard
	ac.card = card
	ac.damage = damage
	ac.skillId = skillId
	return ac
}

func (ac *ActionDamageCard) DoDamage() {
	ac.game.PostActData(ac)

	ac.PostActStream(func() { //需求：一旦受到伤害 就会被动翻面
		if ac.card.isBack {
			ad := NewActionFaceUpCard(ac.game, ac.card, true, ac.card, nil)
			ad.DoFaceUp(nil)
		}
	})
	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_MakeDamage, ac)
	})

	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_OnDamage, ac)
	})

	ac.PostActStream(func() {
		if ac.damage <= 0 {
			return
		}

		card := ac.card

		old := card.GetHP()
		card.SubHP(ac.damage, true)
		new := card.GetHP()
		ac.game.GetCurrentPlayer().Log(fmt.Sprintf("%v 受到%v伤害，血量从%v变为%v ", ac.card.GetOwnInfo(), ac.damage, old, new))

		SyncChangeHP(card, old, new, ac.srcCard, ac.skillId)

	})
	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_AfterBeDamaged, ac)
	})
	ac.PostActStream(func() {
		if ac.card.GetHP() > 0 {
			ac.Stop()
		}
	})

	ac.PostActStream(func() {

		StartSetDeadAndNotify(ac.card, ac.srcCard)

	})

}
