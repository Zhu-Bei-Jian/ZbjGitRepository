package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

// 刘表	自守	被动技：你的军营反伤增加2点
type SkillZiShou struct {
	HeroSkill
}

func (ss *SkillZiShou) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillZiSHou",
		triggerTypes: []TriggerType{TriggerType_OnDamage},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionDamageCard)
			if !ok {
				return
			}
			if ac.srcCard != nil || ac.srcPlayer == nil { //确认是军营反伤
				return
			}
			//扫描场上的友方武将 是否存在 这个技能
			has := false
			var liuBiao *Card
			for _, rows := range g.board.cells {
				for _, cell := range rows {
					if !cell.HasCard() {
						continue
					}
					//造成反伤的玩家是否有 自守
					if cell.Card.owner != ac.srcPlayer {
						continue
					}
					if cell.Card.HasSkill(ss.GetSkillId()) {
						has = true
						liuBiao = cell.Card
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
			g.GetCurrentPlayer().Log(fmt.Sprintf("%v 自守：你的军营反伤增加2点. %v 受到额外2点的反伤", liuBiao.GetOwnInfo(), ac.card.GetOwnInfo()))
			ac.damage += 2

		},
	}

	th2 := TriggerHandler{
		name:         "Skillzishou",
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
			AddCampBuff(card, ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0)

		},
	}

	th3 := TriggerHandler{
		name:         "Skillzishou",
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
			DelCampBuff(card, ss.GetBuffId0())

		},
	}

	return []TriggerHandler{th, th2, th3}
}
