package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

// 74	曹丕	行殇	被动技：每当你的军营受到伤害时，随机对一个敌方武将造成5点伤害
type SkillXingShang struct {
	HeroSkill
}

func (ss *SkillXingShang) TriggerHandler() []TriggerHandler {
	th1 := TriggerHandler{
		name:         "SkillXingShang",
		triggerTypes: []TriggerType{TriggerType_MakeDamageCamp},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionDamageCamp)
			if !ok {
				return
			}
			has := false
			var caoPi *Card
			for _, rows := range g.board.cells {
				for _, v := range rows {
					if !v.HasCard() {
						continue
					}
					if v.owner.seatId == ac.player.seatId && v.GetSkillId() == ss.GetSkillId() { // 受伤玩家 场上有 曹丕
						has = true
						caoPi = v.Card
						break
					}
				}
				if has {
					break
				}
			}
			if has {
				enemy := FindCardsByType(ac.srcCard, gamedef.SelectCardType_MyOwn) // 伤害来源卡牌 的队友
				if len(enemy) == 0 {
					return
				}
				randIndex := gameutil.Intn(len(enemy))
				g.GetCurrentPlayer().Log(fmt.Sprintf("%v 行殇：每当你的军营受到伤害时，随机对一个敌方武将(%v)造成5点伤害", caoPi.GetOwnInfo(), enemy[randIndex].GetOwnInfo()))

				actionData := &ActionDataBase{}
				g.PostActData(actionData)
				actionData.PostActStream(func() {
					g.Send2All(&cmsg.SSyncUseSkill{
						Seat:        caoPi.owner.seatId,
						SkillId:     caoPi.skillId,
						TargetCards: cardIds(enemy[randIndex]),
					})
				})
				actionData.PostActStream(func() {
					g.StartWaitingNoneFloat(1.0, nil)
				})

				actionData.PostActStream(func() {
					NewActionDamageCard(g, enemy[randIndex], caoPi, nil, 5, ss.GetSkillId()).DoDamage()
				})

			}

		},
	}

	th2 := TriggerHandler{
		name:         "SkillXingShang",
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
		name:         "SkillXingShang",
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

	return []TriggerHandler{th1, th2, th3}
}
