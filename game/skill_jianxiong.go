package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//奸雄：被动技：该武将受到伤害但未阵亡时，全体友方武将+1攻击力

type SkillJianXiong struct {
	HeroSkill
}

func (ss *SkillJianXiong) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillJianXiong",
		triggerTypes: []TriggerType{TriggerType_AfterBeDamaged},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionDamageCard)
			if !ok {
				return
			}
			if ac.card.isBack {
				return
			}

			if ac.card.skillId != ss.GetSkillId() {
				return
			}
			if ac.card.IsDead() {
				return
			}
			ac.card.owner.game.PostActData(ss)
			g.GetCurrentPlayer().Log(fmt.Sprintf("触发被动技：%v", ss.skillCfg.Name))

			for _, cells := range g.board.cells {
				tCells := cells
				for _, cell := range tCells {
					tCell := cell
					if !tCell.HasCard() {
						continue
					}
					if tCell.Card.owner != ac.card.owner {
						continue
					}
					tCell.Card.attack++
					ss.PostActStream(func() {
						StartGetBuff(tCell.Card, ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, ac.card)
					})
					ss.PostActStream(func() {
						SyncChangeAttack(tCell.Card, tCell.Card.attack-1, tCell.Card.attack, ac.card)
						g.GetCurrentPlayer().Log(fmt.Sprintf("%v受到奸雄加成，攻击力+1", tCell.heroCfg.Name))
					})

				}
			}

		},
	}
	return []TriggerHandler{th}
}
