package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/proto/cmsg"
)

type ActionDamageCamp struct {
	ActionDataBase

	pos     *Position
	player  *Player //受害的玩家
	srcCard *Card   //伤害来源
	skillId int32

	damage            int32
	backDamage        int32
	extraDamageToCamp int32
}

func NewActionDamageCamp(game *GameBase, player *Player, srcCard *Card, damage int32) *ActionDamageCamp {
	ac := &ActionDamageCamp{}
	ac.game = game
	ac.player = player
	ac.srcCard = srcCard
	ac.skillId = 0
	ac.damage = damage
	ac.backDamage = 0
	return ac
}

func (ac *ActionDamageCamp) DoDamage() {

	ac.game.PostActData(ac)

	ac.PostActStream(func() {
		ac.game.GetTriggerMgr().Trigger(TriggerType_MakeDamageCamp, ac)
	})

	//大本营受伤
	ac.PostActStream(func() {
		if ac.damage+ac.extraDamageToCamp <= 0 {
			return
		}
		var src string
		if ac.srcCard != nil {
			src = ac.srcCard.GetOwnInfo()
		} else {
			src = "无"
		}
		ac.game.GetCurrentPlayer().Log(fmt.Sprintf("%v 的军营受到%v点伤害，伤害来源为%v", ac.player.user.userBrief.Nickname, ac.damage+ac.extraDamageToCamp, src))
		old := ac.player.GetHP()
		ac.player.SubHP(ac.damage + ac.extraDamageToCamp)
		SyncCampChangeHP(ac.player, ac.srcCard, old, ac.damage+ac.extraDamageToCamp)
	})

	ac.PostActStream(func() {
		if ac.player.GetHP() > 0 {
			return
		}

		ac.player.SetDead()
		ac.game.Send2All(&cmsg.SSyncPlayerState{
			Seat:         ac.player.GetSeatID(),
			ConnectState: ac.player.connectState,
			Dead:         ac.player.IsDead(),
			ChangeType:   cmsg.SSyncPlayerState_Die,
		})
	})

	ac.PostActStream(func() {
		ac.game.TryEndGame()
		//ac.game.Print()
	})
}
