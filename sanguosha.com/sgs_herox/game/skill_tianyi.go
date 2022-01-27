package game

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/game/core"
)

//天义: 被动技 你的攻击会从指定目标开始，顺时针影响所有相邻的敌方目标：每个武将对你造成反击伤害-2
type SkillTianYi struct {
	HeroSkill
}

func (ss *SkillTianYi) TriggerHandler() []TriggerHandler {
	th1 := TriggerHandler{
		name:         "SkillTianYi",
		triggerTypes: []TriggerType{TriggerType_BeRetAttack},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionAttackCard)
			if !ok {
				return
			}
			taiShiCi := ac.srcCard //太史慈发动的攻击  ，即将受到敌方反击伤害
			if taiShiCi.isBack {
				return
			}
			if !taiShiCi.HasSkill(ss.GetSkillId()) {
				return
			}
			ac.extraRetDamage -= ss.GetValue(1)

			g.GetCurrentPlayer().Log(fmt.Sprintf("%v 触发天义, 受到%v的反击伤害-%v", taiShiCi.GetOwnInfo(), ac.targetCard.GetOwnInfo(), ss.GetValue(1)))
		},
	}

	th2 := TriggerHandler{
		name:         "SkillTianYi",
		triggerTypes: []TriggerType{TriggerType_BeforeActAttack},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionActAttack)
			if !ok {
				return
			}
			card := ac.srcCard
			if card.isBack {
				return
			}
			if !card.HasSkill(ss.GetSkillId()) {
				return
			}

			if len(ac.targetCards) != 1 || ac.spellCard.heroCfg.Name == "貂蝉" {
				return
			}

			cards := g.board.GetClockWiseEnemyFromTarget(card, ac.targetCards[0])
			ac.targetCards = cards
			g.GetCurrentPlayer().Log("太史慈 触发被动技：天义。攻击会从指定目标开始，顺时针影响所有相邻的敌方目标。目标按顺序如下：")
			for id, v := range cards {
				logrus.Info(id+1, ".", v.GetOwnInfo())
			}
		},
	}
	return []TriggerHandler{th1, th2}
}

func (b *Board) GetClockWiseEnemyFromTarget(src *Card, target *Card) []*Card {

	row := src.cell.Position.Row
	col := src.cell.Position.Col
	var pos = []*Position{{row - 1, col},
		{row, col - 1},
		{row + 1, col},
		{row, col + 1},
	}
	var cards []*Card
	for _, p := range pos {
		if !IsInsideBoard(p.Row, p.Col) {
			continue
		}
		if !b.cells[p.Row][p.Col].HasCard() {
			continue
		}
		if b.cells[p.Row][p.Col].owner == src.owner {
			continue
		}
		cards = append(cards, b.cells[p.Row][p.Col].Card)
	}
	for id, v := range cards {
		if v.cell.Position.Row == target.cell.Position.Row && v.cell.Position.Col == target.cell.Position.Col {
			cards = append(cards[id:], cards[:id]...)
			break
		}
	}
	return cards
}
