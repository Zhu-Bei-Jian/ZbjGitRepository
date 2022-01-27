package game

import (
	"sanguosha.com/sgs_herox/game/core"
)

//共通的buff功能 回合结束时检查过期buff删除
type BuffCommon struct {
	HeroBuff
}

func (ss *BuffCommon) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "BuffCommon",
		triggerTypes: []TriggerType{TriggerType_PhaseEnd},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			roundCount := g.roundCount
			for _, p := range g.players {
				remove := p.Camp.BuffManager.RemoveIfRoundExpire(roundCount)
				if len(remove) > 0 {
					SyncCampChangeBuff(p, nil)
				}
			}
			g.PostActData(ss)

			//cards := g.board.GetCardsBySeatId(p.GetSeatID())
			var cards []*Card
			for _, rows := range g.board.cells {
				for _, cell := range rows {
					if cell.Card == nil {
						continue
					}
					cards = append(cards, cell.Card)
				}
			}
			for _, v := range cards {
				card := v
				expireBuffIds, hasRoundBuff := v.BuffManager.ExpireBuffAtRound(roundCount)
				ss.PostActStream(func() {
					StartLoseBuffs(card, expireBuffIds)
				})
				ss.PostActStream(func() {
					if hasRoundBuff {
						SyncChangeBuff(v, nil, 0)
					}
				})

			}

		},
	}

	return []TriggerHandler{th}
}
