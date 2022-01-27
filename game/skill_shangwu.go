package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//尚武：被动技，己方武将攻击军营造成的伤害+2。
type SkillShangWu struct {
	HeroSkill
}

func (ss *SkillShangWu) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillShangWu",
		triggerTypes: []TriggerType{TriggerType_AttackCamp},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			//adi := g.ActDataPeek()
			ac, ok := ad.(*ActionAttackCamp)
			if !ok {
				return
			}
			if ac.srcCard.isBack {
				return
			}
			cards := g.skillCards(g.GetCurrentPlayer(), ss.GetSkillId())

			if len(cards) == 1 {
				ac.extraDamageToCamp += g.GetBuffCfg(ss.GetBuffId0()).GetBuffDamage()
				g.GetCurrentPlayer().Log(fmt.Sprintf("尚武：被动技，己方武将(%v)攻击军营造成的伤害+2。", ac.srcCard))
			}

		},
	}
	th2 := TriggerHandler{
		name:         "SkillShangWu",
		triggerTypes: []TriggerType{TriggerType_GetSkill},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			//adi := g.ActDataPeek()
			if len(params) != 1 {
				return
			}
			card, ok := params[0].(*Card)
			if !ok {
				return
			}
			if card.skillCfg.SkillID != ss.GetSkillId() {
				return
			}
			myCards := FindAllMyOwnCards(card)
			for _, v := range myCards {
				StartGetBuff(v, ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)
			}

		},
	}

	th3 := TriggerHandler{
		name:         "SkillShangWu",
		triggerTypes: []TriggerType{TriggerType_LoseSkill},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			//adi := g.ActDataPeek()
			if len(params) != 1 {
				return
			}
			card, ok := params[0].(*Card)
			if !ok {
				return
			}
			if card.skillCfg.SkillID != ss.GetSkillId() {
				return
			}
			myCards := FindAllMyOwnCards(card)
			for _, v := range myCards {
				StartLoseBuff(v, ss.GetBuffId0())
			}

		},
	}
	return []TriggerHandler{th, th2, th3}
}

// 光环buff 的同步时机点
// 1 姜维 翻至正面 且 未被沉默(TriggerType get skill) ，己方所有武将 start Get   Buff尚武
// 2 姜维 死亡 且 此前未被沉默(TriggerType lose skill) ，己方所有武将 start Lose  Buff尚武
// 3 姜维 在正面状态且未被沉默，此时被沉默，(TriggerType lose skill)己方所有武将 start Lose  Buff尚武
// 4 姜维 在正面状态且未被沉默，此时被孙权制衡翻至背面，(TriggerType lose skill)己方所有武将 start Lose  Buff尚武
