package game

import (
	"fmt"
)

type ActionAttackCard struct {
	ActionDataBase

	srcCard    *Card
	targetCard *Card

	targetCell *Cell //为解决 目标卡牌死亡后，card于cell解绑后， 无法通过card访问cell的问题

	extraRetDamage    int32 //额外反击伤害
	targetExtraDamage int32 // >0 被攻击者的 额外伤害 ; <0 可以为负值 ，代表减伤

}

func NewActionAttackCard(g *GameBase, opCard *Card, targetCard *Card) *ActionAttackCard {
	ac := &ActionAttackCard{}
	ac.game = g
	ac.srcCard = opCard
	ac.targetCard = targetCard
	ac.targetCell = ac.targetCard.cell
	return ac
}

func (ac *ActionAttackCard) DoAttack() {
	ac.game.PostActData(ac)

	opCard := ac.srcCard
	targetCard := ac.targetCard

	ac.PostActStream(func() {
		ac.game.GetCurrentPlayer().Log(fmt.Sprintf("%v对%v 发动攻击", ac.srcCard.GetOwnInfo(), ac.targetCard.GetOwnInfo()))
		ac.game.GetTriggerMgr().Trigger(TriggerType_AttackCard, ac)
	})

	//因需求改变（现只要受到伤害就会被动翻面），所以原本的被动翻面流程 移到 dodamage

	// 触发 被攻击者的 beAttacked   触发器
	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_BeAttackCard, ac)
	})
	//被攻击方受伤
	ac.PostActStream(func() {

		damage := opCard.attack + ac.targetExtraDamage
		if damage <= 0 {
			return
		}
		actDamage := NewActionDamageCard(ac.game, ac.targetCard, ac.srcCard, nil, damage, 0)
		actDamage.DoDamage()
	})

	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_BeRetAttack, ac)
	})

	//攻击方受伤   //即 被反击  BeRetAttack
	ac.PostActStream(func() {
		if opCard.IsDead() {
			return
		}
		damage := targetCard.attack + ac.extraRetDamage
		if damage <= 0 {
			return
		}
		actDamage := NewActionDamageCard(ac.game, ac.srcCard, ac.targetCard, nil, damage, 0)
		actDamage.DoDamage()
	})

	//攻击后触发
	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_AfterAttack, ac)
	})

	ac.PostActStream(func() { //卡牌受伤后检测，无法判断平局， 改为攻击后检测
		ac.game.TryEndGame()
	})
}
