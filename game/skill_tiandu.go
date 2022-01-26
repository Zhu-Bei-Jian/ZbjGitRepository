package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
)

type SkillTianDu struct {
	HeroSkill
}

func (ss *SkillTianDu) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillTianDu",
		triggerTypes: []TriggerType{TriggerType_PhaseEnd},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {

			for _, rows := range g.board.cells {
				for _, cell := range rows {
					if !cell.HasCard() {
						continue
					}
					if g.GetCurrentPlayer() != cell.Card.owner {
						continue
					}
					if cell.Card.skillId == ss.GetSkillId() {
						//cell.Card.hp -= 5
						oldHP := cell.Card.GetHP()
						cell.Card.SubHP(5, false)
						SyncChangeHP(cell.Card, oldHP, cell.Card.GetHP(), nil, ss.GetSkillId())

						if cell.Card.GetHP() <= 0 {
							StartSetDeadAndNotify(cell.Card, nil)

						}
						g.GetCurrentPlayer().Log(fmt.Sprintf("触发被动技：%v", ss.skillCfg.Name))
					}
				}
			}

		},
	}
	return []TriggerHandler{th}
}
