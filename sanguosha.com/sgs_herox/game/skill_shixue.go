package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/cmsg"
)

//被动技：你造成的伤害可以给军营恢复等量的血量
//造成伤害：1.对card造成伤害  2.对camp造成伤害
type SkillShiXue struct {
	HeroSkill
}

func (ss *SkillShiXue) TriggerHandler() []TriggerHandler {
	th1 := TriggerHandler{
		name:         "SkillShiXueDamageCard",
		triggerTypes: []TriggerType{TriggerType_MakeDamage},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionDamageCard)
			if !ok {
				return
			}
			if ac.srcCard == nil {
				return
			}
			if ac.srcCard.isBack {
				return
			}

			if ac.srcCard.GetSkillId() != ss.GetSkillId() {
				return
			}
			// srcCard 为 魏延  ac.damage为魏延实际造成的伤害
			if ac.srcCard.owner.hp == g.config.CampHP {
				g.GetCurrentPlayer().Log(fmt.Sprintf("%v技能触发，即将为己方军营回复血量，但军营已经满血，取消加血", ac.srcCard.GetOwnInfo()))
				return
			}
			oldhp := ac.srcCard.owner.hp
			ac.srcCard.owner.hp += ac.damage //给自己的军营回复
			if ac.srcCard.owner.hp > g.config.CampHP {
				ac.srcCard.owner.hp = g.config.CampHP
			}
			g.Send2All(&cmsg.SSyncUseSkill{
				Seat:        ac.srcCard.owner.seatId,
				SkillId:     ss.GetSkillId(),
				TargetCards: []int32{ac.card.ID()},
			})

			g.StartWaitingNone(1, nil)

			SyncCampChangeHP(ac.srcCard.owner, ac.srcCard, ac.srcCard.owner.hp-ac.damage, ac.damage)

			g.GetCurrentPlayer().Log(fmt.Sprintf("%v技能触发,造成的伤害可以给军营恢复等量的血量,军营实际恢复%v血", ac.srcCard.GetOwnInfo(), ac.srcCard.owner.hp-oldhp))
		},
	}
	th2 := TriggerHandler{
		name:         "SkillShiXueDamageCamp",
		triggerTypes: []TriggerType{TriggerType_MakeDamageCamp},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionDamageCamp)
			if !ok {
				return
			}

			if ac.srcCard != nil && ac.srcCard.isBack {
				return
			}

			if ac.srcCard.GetSkillId() != ss.GetSkillId() {
				return
			}

			if ac.srcCard.owner.hp == g.config.CampHP {
				g.GetCurrentPlayer().Log(fmt.Sprintf("%v 攻击敌方军营.嗜血触发：即将为己方军营回复血量，但军营已经满血，取消加血", ac.srcCard.GetOwnInfo()))
				return
			}

			oldhp := ac.srcCard.owner.hp
			ac.srcCard.owner.AddHP(ac.damage + ac.extraDamageToCamp)
			if ac.srcCard.owner.hp > g.config.CampHP {
				ac.srcCard.owner.hp = g.config.CampHP
			}

			g.Send2All(&cmsg.SSyncUseSkill{
				Seat:         ac.srcCard.owner.seatId,
				SkillId:      ss.GetSkillId(),
				TargetSeatId: ac.player.seatId,
			})

			g.StartWaitingNone(1, nil)

			SyncCampChangeHP(ac.srcCard.owner, ac.srcCard, oldhp, ac.srcCard.owner.GetHP()-oldhp)
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v技能触发,造成的伤害可以给军营恢复等量的血量,军营实际恢复%v血", ac.srcCard.GetOwnInfo(), ac.srcCard.owner.GetHP()-oldhp))

		},
	}

	return []TriggerHandler{th1, th2}
}
