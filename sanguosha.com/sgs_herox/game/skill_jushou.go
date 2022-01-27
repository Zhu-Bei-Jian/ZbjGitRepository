package game

import (
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

type SkillJuShou struct {
	HeroSkill
}

func (ss *SkillJuShou) TriggerHandler() []TriggerHandler {
	th1 := TriggerHandler{
		name:         "SkillJuShou",
		triggerTypes: []TriggerType{TriggerType_CheckAP_PlaceCard},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionCheckPlaceCard)
			if !ok {
				return
			}

			if ac.skillCfg.SkillID != ss.GetSkillId() {
				return
			}
			ac.ap = g.config.HeavyCost
			bfConfig, ok := g.config.Get(ss.GetBuffId0())
			if !ok {
				logrus.Fatal("???????buff呢？？？？？？")
				return
			}
			ac.buffs = append(ac.buffs, HeavyHaloBuff{
				buff: buff{
					buffCfg:    bfConfig,
					ExpireType: gameconf.ExpireTyp_ETInvalid,
				},
				effectType:   EffectType_camp,
				targetSeatID: g.GetCurrentPlayer().seatId,
			})
		},
	}

	th2 := TriggerHandler{
		name:         "SkillJuShou",
		triggerTypes: []TriggerType{TriggerType_CheckHeavyAP},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionCheck)
			if !ok {
				return
			}

			if ac.card.skillId != ss.GetSkillId() {
				return
			}
			ac.ap[MOVE_CARD] = g.config.HeavyCost
			ac.ap[ATTACK_CAMP] = g.config.HeavyCost
			ac.ap[ATTACK_CARD] = g.config.HeavyCost

		},
	}
	th3 := TriggerHandler{
		name:         "Skilljushou",
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

	th4 := TriggerHandler{
		name:         "Skilljushou",
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

	th5 := TriggerHandler{
		name:         "skillJushou",
		triggerTypes: []TriggerType{TriggerType_MakeDamageCamp},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			ac, ok := ad.(*ActionDamageCamp)
			if !ok {
				return
			}
			if ac.player.HasBuff(ss.GetBuffId0()) {
				ac.extraDamageToCamp -= INF
			}

		},
	}
	return []TriggerHandler{th1, th2, th3, th4, th5}
}
