package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
)

// 64	左慈	同命	被动技：消灭主动攻击此武将的敌方武将。

type SkillTongMing struct {
	HeroSkill
}

func (ss *SkillTongMing) PreUseSkill() {
	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})

	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})
}

func (ss *SkillTongMing) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillTongMing",
		triggerTypes: []TriggerType{TriggerType_BeAttackCard},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionAttackCard)
			if !ok {
				return
			}

			if ac.targetCard.skillId == -1 {
				return
			}
			if ac.targetCard.skillCfg.SkillID != ss.GetSkillId() {
				return
			}
			g.GetCurrentPlayer().Log(fmt.Sprintf("触发被动技：%v", ss.skillCfg.Name))
			g.PostActData(ss)
			ss.PostActStream(func() {
				ss.SetCard(ac.targetCard)
			})
			ss.PostActStream(func() {
				ss.PreUseSkill()
			})
			ss.PostActStream(func() {
				oldHP := ac.srcCard.GetHP()
				ac.srcCard.SubHP(INF, false)
				SyncChangeHP(ac.srcCard, oldHP, ac.srcCard.GetHP(), ac.targetCard, ss.GetSkillId())
			})
			ss.PostActStream(func() {
				StartSetDeadAndNotify(ac.srcCard, ac.targetCard)
			})

		},
	}
	return []TriggerHandler{th}
}
