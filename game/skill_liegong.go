package game

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/game/core"
)

// 烈弓：被动技，你无攻击距离限制，你不会受到距离大于1的武将的反击，不能攻击军营。
type SkillLieGong struct {
	HeroSkill
}

func (ss *SkillLieGong) PreUseSkill() {
	g := ss.card.owner.game
	g.PostActData(ss)

	ss.PostActStream(func() {
		ss.HeroSkill.PreUseSkill()
	})

	ss.PostActStream(func() {
		g.StartWaitingNone(1, nil)
	})
}

func (ss *SkillLieGong) TriggerHandler() []TriggerHandler {
	th1 := TriggerHandler{
		name:         "SkillLieGong",
		triggerTypes: []TriggerType{TriggerType_CheckDistance},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac := ad.(*ActionCheck)
			if ac.card.isBack {
				return
			}
			if !ac.card.HasSkill(ss.GetSkillId()) {
				return
			}

			setAttackDistanceAndNotify(ac.card, INF)
			ac.canAttackCamp = false
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v烈弓生效 : 无攻击距离限制", ac.card.GetOwnInfo()))

			if ac.card != nil && ac.targetCard != nil && calDistanceByPosition(ac.card.cell.Position, ac.targetCard.cell.Position) > 1 {
				card := ac.card
				skill, ok := NewSkill(card.skillCfg)
				if !ok {
					logrus.Errorf("new skill %d fail", card.skillCfg.SkillID)
					return
				}

				skill.SetCard(card)
				skill.SetTargets([]*Card{ac.targetCard})
				skill.SetDataToClient([]int32{ac.card.heroCfg.HeroID})
				ac.AddEffectSkill(skill)
			}

		},
	}
	th2 := TriggerHandler{
		name:         "SkillLieGong",
		triggerTypes: []TriggerType{TriggerType_BeRetAttack},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionAttackCard)
			if !ok {
				return
			}
			if ac.srcCard.isBack {
				return
			}
			if ac.srcCard.skillId != ss.GetSkillId() {
				return
			}
			p1 := ac.srcCard.cell.Position
			p2 := ac.targetCell.Position
			dis := Abs(p1.Row-p2.Row) + Abs(p2.Col-p1.Col)
			//fmt.Println(ac.srcCard.heroCfg.Name, ac.targetCard.heroCfg.Name, dis)
			if dis > 1 {
				ac.extraRetDamage -= INF
				g.GetCurrentPlayer().Log(fmt.Sprintf("%v烈弓生效：(黄忠此次与攻击目标距离为%v),不会受到距离大于1的武将的反击. ", ac.srcCard.GetOwnInfo(), dis))
			} else {
				g.GetCurrentPlayer().Log(fmt.Sprintf("%v此次与攻击目标距离为%v,烈弓的减伤效果 未触发 ", ac.srcCard.GetOwnInfo(), dis))
			}

		},
	}

	return []TriggerHandler{th1, th2}
}
