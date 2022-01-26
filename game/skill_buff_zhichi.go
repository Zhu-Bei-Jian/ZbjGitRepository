package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
)

//buff智迟  每次攻击后 恢复 3点生命值
type BuffZhiChi struct {
	HeroBuff
}

func (ss *BuffZhiChi) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "BuffZhiChi",
		triggerTypes: []TriggerType{TriggerType_AfterAttack},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			switch ad.(type) {
			case *ActionAttackCamp:
				ac, ok := ad.(*ActionAttackCamp)
				if !ok {
					return
				}
				if !ac.srcCard.HasBuff(ss.GetBuffId()) {
					return
				}
				var recoverValue int32 = 0
				for _, v := range ac.srcCard.buffs {
					if v.buffCfg.BuffID == ss.GetBuffId() {
						recoverValue += 3
					}
				}
				oldHP := ac.srcCard.GetHP()
				ac.srcCard.AddHP(recoverValue)
				SyncChangeHP(ac.srcCard, oldHP, ac.srcCard.GetHP(), nil, ss.GetBuffId())
				g.GetCurrentPlayer().Log(fmt.Sprintf("%v 触发智迟,实际恢复%v血", ac.srcCard.GetOwnInfo(), ac.srcCard.GetHP()-oldHP))
			case *ActionAttackCard:
				ac, ok := ad.(*ActionAttackCard)
				if !ok {
					return
				}

				if !ac.srcCard.HasBuff(ss.GetBuffId()) {
					return
				}

				//ac.srcCard.hp += 3
				oldHP := ac.srcCard.GetHP()
				ac.srcCard.AddHP(3)
				SyncChangeHP(ac.srcCard, oldHP, ac.srcCard.GetHP(), nil, ss.GetBuffId())
				g.GetCurrentPlayer().Log(fmt.Sprintf("%v 触发智迟,实际恢复%v血", ac.srcCard.GetOwnInfo(), ac.srcCard.GetHP()-oldHP))

			default:
				return
			}

		},
	}

	return []TriggerHandler{th}
}
