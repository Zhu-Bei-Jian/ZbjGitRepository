package game

import "sanguosha.com/sgs_herox/game/core"

type BuffSongCi struct {
	HeroBuff
}

func (ss *BuffSongCi) TriggerHandler() []TriggerHandler {
	th := TriggerHandler{
		name:         "BuffSongCi",
		triggerTypes: []TriggerType{TriggerType_PhaseEnd},
		handle: func(g *GameBase, ad core.IActionData, params ...interface{}) {
			for _, rows := range g.board.cells {
				for _, v := range rows {
					if !v.HasCard() || !v.HasBuff(ss.GetBuffId()) || v.owner.seatId != g.GetCurrentPlayer().seatId {
						continue
					}
					oldAt := v.GetAttack()
					var value int32 = 0
					for _, bf := range v.Card.buffs {
						if bf.buffCfg.BuffID == ss.GetBuffId() {
							value += 1 * bf.buffCount
						}
					}
					v.AddAttack(value)
					SyncChangeAttack(v.Card, oldAt, v.GetAttack(), nil)

				}
			}

		},
	}

	return []TriggerHandler{th}
}
