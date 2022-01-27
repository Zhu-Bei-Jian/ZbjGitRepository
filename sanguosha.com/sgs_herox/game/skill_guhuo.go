package game

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/game/core"
)

//蛊惑：被动技，回合结束结束，恢复5点生命值（不超过生命值上限）

type SkillGuHuo struct {
	HeroSkill
}

func (ss *SkillGuHuo) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillGuHuo",
		triggerTypes: []TriggerType{TriggerType_PhaseBegin},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {

			for _, cells := range g.board.cells {
				for _, cell := range cells {
					if !cell.HasCard() || cell.Card.isBack || cell.Card.owner.seatId != g.GetCurrentPlayer().seatId {
						continue
					}
					if cell.Card.HasSkill(ss.skillCfg.SkillID) {
						cd := cell.Card
						oldHp := cd.GetHP()
						cd.AddHP(ss.GetValue(1))
						if cd.GetHP() == oldHp {
							logrus.Info("-蛊惑-：于吉已满血，此次触发将不回复血量 ")
							return
						}
						SyncChangeHP(cd, oldHp, cd.GetHP(), cd, ss.GetSkillId())
						g.GetCurrentPlayer().Log(fmt.Sprintf("%v触发被动技 -蛊惑- ，于吉回复 %v 点hp", cd.GetOwnInfo(), cd.GetHP()-oldHp))
					}

				}
			}

		},
	}
	return []TriggerHandler{th}
}
