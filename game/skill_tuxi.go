package game

import (
	"fmt"
	"sanguosha.com/sgs_herox/game/core"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//突袭：被动技，当你处于敌方后排时，攻击力+5
//持续性的技能
type SkillTuXi struct {
	HeroSkill
}

func (ss *SkillTuXi) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "SkillTuXi",
		triggerTypes: []TriggerType{TriggerType_AfterPosChange, TriggerType_GetSkill},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			var card *Card
			switch t := ad.(type) {
			case *ActionMove:
				card = t.card
			case *ActionFaceUpCard:
				card = t.card
			case *ActionSilence:
				card = t.card
			case *ActionDataBase:
				for _, v := range params {
					cd, ok := v.(*Card)
					if ok && cd.skillId == ss.GetSkillId() {
						card = cd
					}
				}
			default:
				return
			}
			if card == nil {
				return
			}
			card.owner.game.PostActData(ss)
			canActivate := IsInEnemyZone(card) && card.HasSkill(ss.GetSkillId())
			if canActivate {
				if card.HasBuff(ss.GetBuffId0()) {
					return
				}
				StartGetBuff(card, ss.GetBuffId0(), gameconf.ExpireTyp_ETInvalid, 0, card)

			} else {
				if !card.HasBuff(ss.GetBuffId0()) {
					return
				}
				StartLoseBuff(card, ss.GetBuffId0())
			}

		},
	}

	return []TriggerHandler{th}
}

func StartGetBuff(card *Card, buffId int32, expireType gameconf.ExpireTyp, expireV int32, srcCard *Card) {

	g := card.owner.game

	buffCfg, exist := g.config.Buff.Get(buffId)
	if !exist {
		return
	}

	ad := &ActionDataBase{}
	g.PostActData(ad)

	ad.PostActStream(func() {
		//if card.BuffManager.HasBuff(buffId) {
		//	ad.Stop()
		//	return
		//}

		card.BuffManager.Add(buffCfg, g.roundCount, expireType, expireV)

		g.GetCurrentPlayer().Log(fmt.Sprintf("%v 获得 buff %v %v", card.GetOwnInfo(), buffCfg.GetName(), buffCfg.GetDesc()))
		SyncChangeBuff(card, srcCard, buffId)
	})

	ad.PostActStream(func() {
		buff, ok := newBuff(buffCfg)
		if !ok {
			return
		}
		buff.OnEnable(card)
	})
}

func StartLoseBuff(card *Card, buffId int32) {
	g := card.owner.game

	buffCfg, exist := g.config.Buff.Get(buffId)
	if !exist {
		return
	}

	ad := &ActionDataBase{}
	g.PostActData(ad)

	ad.PostActStream(func() {
		if !card.BuffManager.HasBuff(buffId) {
			ad.Stop()
			return
		}
		card.BuffManager.Remove(buffId)
		SyncChangeBuff(card, nil, 0)
	})

	ad.PostActStream(func() {
		buff, ok := newBuff(buffCfg)
		if !ok {
			return
		}
		buff.OnDisable(card)
		bfCfg, _ := g.config.Buff.Get(buffId)
		g.GetCurrentPlayer().Log(fmt.Sprintf("%v 失去buff %v", card.GetOwnInfo(), bfCfg.GetName()))
	})
}

func StartLoseBuffs(card *Card, buffIds []int32) {
	g := card.owner.game

	ad := &ActionDataBase{}
	g.PostActData(ad)

	for _, v := range buffIds {
		buffId := v
		ad.PostActStream(func() {
			StartLoseBuff(card, buffId)
		})
	}
}
