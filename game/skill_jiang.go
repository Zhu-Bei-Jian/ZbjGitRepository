package game

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/game/core"
)

// 35	孙策	激昂	被动技：主动攻击时，对目标后方武将造成3点伤害
type SkillJiAng struct {
	HeroSkill
}

func (ss *SkillJiAng) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillJiAng",
		triggerTypes: []TriggerType{TriggerType_AttackCard},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionAttackCard)
			if !ok {
				return
			}
			if ac.srcCard.isBack {
				return
			}
			if !ac.srcCard.HasSkill(ss.GetSkillId()) {
				return
			}

			targets := g.board.pierceThrough(ac.srcCard, ac.targetCard)
			if len(targets) == 0 {
				return
			}
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v：激昂触发,对%v后方的卡牌造成3点伤害", ac.srcCard.GetOwnInfo(), ac.targetCard.GetOwnInfo()))
			t := targets[0]
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v 受到来自孙策的3点伤害", t.GetOwnInfo()))

			g.PostActData(ss)

			ss.PostActStream(func() {
				card := ac.srcCard
				skill, ok := NewSkill(card.skillCfg)
				if !ok {
					logrus.Errorf("new skill %d fail", card.skillCfg.SkillID)
					return
				}
				skill.SetCard(card)
				skill.SetTargets([]*Card{ac.targetCard})
				skill.SetDataToClient(cardIds(targets...))
				skill.PreUseSkill()
			})
			ss.PostActStream(func() {
				g.StartWaitingNone(1, nil)
			})
			ss.PostActStream(func() {
				NewActionDamageCard(ac.srcCard.owner.game, t, ac.srcCard, ac.srcCard.owner, 3, ss.GetSkillId()).DoDamage()
			})

		},
	}
	return []TriggerHandler{th}
}
