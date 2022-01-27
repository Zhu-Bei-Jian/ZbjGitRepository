package game

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/cmsg"
)

//诛害：被动技，每回合第一次己方武将的攻击不花费行动点。
type SkillZhuHai struct {
	HeroSkill
}

func (ss *SkillZhuHai) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillZhuHai",
		triggerTypes: []TriggerType{TriggerType_CheckAP},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionCheck)
			if !ok {
				return
			}
			if ac.card.isBack {
				return
			}
			//扫描场上的友方武将 是否存在 这个技能
			has := false
			for _, rows := range g.board.cells {
				for _, cell := range rows {
					if !cell.HasCard() {
						continue
					}
					if cell.Card.owner != ac.card.owner {
						continue
					}
					if cell.Card.HasSkill(ss.GetSkillId()) {
						has = true
						break
					}
				}
				if has {
					break
				}
			}
			if !has {
				return
			}

			if g.GetCurrentPlayer().attackCount == 0 {
				//当前玩家还未攻击过
				//下一次攻击 即为本回合的第一次攻击
				//下次攻击（卡牌或者军营） 耗费 0 ap

				ac.ap[ATTACK_CARD] = 0
				ac.ap[ATTACK_CAMP] = 0
				logrus.Infof("skill_zhuhai.go Line:53, 卡牌攻击需要消耗的AP %v  roomId:%v", ac.ap[ATTACK_CAMP], g.roomId)
			}
			g.GetCurrentPlayer().Log(fmt.Sprintf("触发被动技：%v", ss.skillCfg.Name))
		},
	}

	th2 := TriggerHandler{
		name:         "SkillZhuHai",
		triggerTypes: []TriggerType{TriggerType_PhaseBegin}, //通知播放动画
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {

			//扫描场上的友方武将 是否存在 这个技能
			var xuShu *Card
			has := false
			for _, rows := range g.board.cells {
				for _, cell := range rows {
					if !cell.HasCard() {
						continue
					}
					if cell.Card.owner != g.GetCurrentPlayer() {
						continue
					}
					if cell.Card.HasSkill(ss.GetSkillId()) {
						has = true
						xuShu = cell.Card
						break
					}
				}
				if has {
					break
				}
			}
			if !has {
				return
			}
			g.PostActData(ss)
			ss.PostActStream(func() {
				g.StartWaitingNone(1, nil)
			})
			ss.PostActStream(func() {
				g.Send2All(&cmsg.SSyncUseSkill{
					SkillId: ss.GetSkillId(),
					Card:    xuShu.ID(),
				})
			})

		},
	}

	th3 := TriggerHandler{
		name:         "SkillZhuHai",
		triggerTypes: []TriggerType{TriggerType_GetSkill}, //通知播放动画
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			if len(params) != 1 {
				return
			}
			card, ok := params[0].(*Card)
			if !ok {
				return
			}
			if card.owner != g.GetCurrentPlayer() || card.GetSkillId() != ss.GetSkillId() {
				return
			}

			g.PostActData(ss)
			ss.PostActStream(func() {
				g.StartWaitingNoneFloat(0.6, nil)
			})
			ss.PostActStream(func() {
				g.Send2All(&cmsg.SSyncUseSkill{
					SkillId: ss.GetSkillId(),
					Card:    card.ID(),
				})
			})
		},
	}

	return []TriggerHandler{th, th2, th3}
}
