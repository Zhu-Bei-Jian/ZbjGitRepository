package game

import "fmt"

type ActionAttackCamp struct {
	ActionDataBase

	srcCard           *Card   //攻击发动者
	campPlayer        *Player //被攻击的 玩家
	retDamage         int32
	extraDamageToCamp int32
	multi             int32 //伤害倍数 默认一倍
}

func NewActionAttackCamp(g *GameBase, campPlayer *Player, srcCard *Card) *ActionAttackCamp {
	ac := &ActionAttackCamp{}
	ac.game = g
	ac.srcCard = srcCard
	ac.campPlayer = campPlayer
	ac.multi = 1
	return ac
}

func (ac *ActionAttackCamp) DoAttack() {
	ac.game.GetCurrentPlayer().Log(fmt.Sprintf("%v对%v的军营发动攻击", ac.srcCard.GetOwnInfo(), ac.campPlayer.user.userBrief.Nickname))

	srcCell := ac.srcCard.cell
	ac.game.PostActData(ac)

	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_AttackCamp, ac)
	})
	ac.PostActStream(func() {
		ac.srcCard.attackCntInTurn++
		ac.srcCard.owner.attackCount++
	})

	//军营受伤
	ac.PostActStream(func() {
		damage := ac.srcCard.attack*ac.multi + ac.extraDamageToCamp
		actDamage := NewActionDamageCamp(ac.game, ac.campPlayer, ac.srcCard, damage)
		actDamage.DoDamage()
	})

	ac.PostActStream(func() {
		switch srcCell.Row {
		case 1:
			ac.retDamage = 2
		default:

		}
	})

	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_BeRetAttackByCamp, ac)
	})

	//攻击者受伤
	ac.PostActStream(func() {
		damage := ac.retDamage
		if damage <= 0 {
			return
		}
		//军营对 攻击卡牌 造成 反伤
		if ac.srcCard.IsDead() {
			return
		}
		actDamage := NewActionDamageCard(ac.game, ac.srcCard, nil, ac.campPlayer, damage, 0)
		actDamage.DoDamage()
	})
	ac.PostActStream(func() {
		//即使 攻击者 已死亡，仍要触发
		ac.game.GetTriggerMgr().Trigger(TriggerType_AfterAttack, ac)

	})

}
