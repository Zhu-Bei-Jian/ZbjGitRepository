package game

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

func SyncChangePlace(card *Card) {
	card.owner.game.SendSeatMsg(func(seatId int32) proto.Message {
		return &cmsg.SSyncCard{
			SeatId:    card.GetPlayer().GetSeatID(),
			Card:      card.ToDef(seatId),
			SpellCard: 0,
			Changes: []*cmsg.SSyncCard_Change{{
				ChangeType: cmsg.SSyncCard_Place,
				NewPos:     card.cell.Position.ToDef(),
			}},
		}
	})
}

func SyncChangePos(oldPos *Position, g *GameBase, card *Card) {
	var oldPosDef *gamedef.Position
	if oldPos == nil {
		//oldPos = &card.cell.Position
	}
	if oldPos != nil {
		oldPosDef = oldPos.ToDef()
	}
	g.SendSeatMsg(func(seatId int32) proto.Message {
		return &cmsg.SSyncCard{
			SeatId: card.GetPlayer().GetSeatID(),
			Card:   card.ToDef(seatId),
			Changes: []*cmsg.SSyncCard_Change{{
				ChangeType: cmsg.SSyncCard_Pos,
				OldPos:     oldPosDef,
				NewPos:     card.cell.Position.ToDef(),
			}},
		}
	})
}

func SyncChangeAttack(card *Card, old, new int32, spellCard *Card) {
	if old == new {
		return
	}
	if card.isBack {
		// 需求变化：不能看到敌方背面卡牌的属性变化
		//如果是背面卡牌，只同步给卡牌拥有者
		card.owner.SendMsg(&cmsg.SSyncCard{
			SeatId:    card.GetPlayer().GetSeatID(),
			Card:      card.ToDef(card.owner.seatId),
			SpellCard: spellCard.ID(),
			Changes: []*cmsg.SSyncCard_Change{{
				ChangeType: cmsg.SSyncCard_Attack,
				Old:        old,
				New:        new,
			}},
		})
	} else {
		card.owner.game.SendSeatMsg(func(seatId int32) proto.Message {
			return &cmsg.SSyncCard{
				SeatId:    card.GetPlayer().GetSeatID(),
				Card:      card.ToDef(seatId),
				SpellCard: spellCard.ID(),
				Changes: []*cmsg.SSyncCard_Change{{
					ChangeType: cmsg.SSyncCard_Attack,
					Old:        old,
					New:        new,
				}},
			}
		})
	}

}

func SyncChangeLeftAtkCnt(card *Card, old, new int32, spellCard *Card) {
	if old == new {
		return
	}
	card.owner.game.SendSeatMsg(func(seatId int32) proto.Message {
		return &cmsg.SSyncCard{
			SeatId:    card.GetPlayer().GetSeatID(),
			Card:      card.ToDef(seatId),
			SpellCard: spellCard.ID(),
			Changes: []*cmsg.SSyncCard_Change{{
				ChangeType: cmsg.SSyncCard_AttackCnt,
				Old:        old,
				New:        new,
				Change:     new - old,
			}},
		}
	})
}

func SyncChangeAttackDistance(card *Card, old, new int32, spellCard *Card) {
	card.owner.game.SendSeatMsg(func(seatId int32) proto.Message {
		return &cmsg.SSyncCard{
			SeatId:    card.GetPlayer().GetSeatID(),
			Card:      card.ToDef(seatId),
			SpellCard: spellCard.ID(),
			Changes: []*cmsg.SSyncCard_Change{{
				ChangeType: cmsg.SSyncCard_AttackDistance,
				Old:        old,
				New:        new,
			}},
		}
	})
}

func bool2Int32(v bool) int32 {
	if v {
		return 1
	} else {
		return 0
	}
}

func SyncChangeCanTurnUp(card *Card, old, new bool, spellCard *Card) {
	card.owner.game.SendSeatMsg(func(seatId int32) proto.Message {
		return &cmsg.SSyncCard{
			SeatId:    card.GetPlayer().GetSeatID(),
			Card:      card.ToDef(seatId),
			SpellCard: spellCard.ID(),
			Changes: []*cmsg.SSyncCard_Change{{
				ChangeType: cmsg.SSyncCard_CanFaceUp,
				Old:        bool2Int32(old),
				New:        bool2Int32(new),
			}},
		}
	})
}

func SyncChangeHP(card *Card, old, new int32, spellCard *Card, spellSKillId int32) {
	if old == new {
		return
	}

	if card.isBack {
		// 需求变化：不能看到敌方背面卡牌的属性变化
		//如果是背面卡牌，只同步给卡牌拥有者
		card.owner.SendMsg(&cmsg.SSyncCard{
			SeatId:       card.GetPlayer().GetSeatID(),
			Card:         card.ToDef(card.owner.seatId),
			SpellCard:    spellCard.ID(),
			SpellSkillId: spellSKillId,
			Changes: []*cmsg.SSyncCard_Change{{
				ChangeType: cmsg.SSyncCard_HP,
				Old:        old,
				New:        new,
			}},
		})

	} else {
		card.owner.game.SendSeatMsg(func(seatId int32) proto.Message {
			return &cmsg.SSyncCard{
				SeatId:    card.GetPlayer().GetSeatID(),
				Card:      card.ToDef(seatId),
				SpellCard: spellCard.ID(),
				Changes: []*cmsg.SSyncCard_Change{{
					ChangeType: cmsg.SSyncCard_HP,
					Old:        old,
					New:        new,
				}},
			}
		})
	}

}

func SyncChangeFaceUp(card *Card, spellCard *Card) {
	card.owner.game.SendSeatMsg(func(seatId int32) proto.Message {
		return &cmsg.SSyncCard{
			SeatId:    card.GetPlayer().GetSeatID(),
			Card:      card.ToDef(-1),
			SpellCard: spellCard.ID(),
			Changes: []*cmsg.SSyncCard_Change{{
				ChangeType: cmsg.SSyncCard_FaceUp,
			}},
		}
	})
}

func SyncChangeDie(card *Card, spellCard *Card) {
	if spellCard == nil {
		card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("玩家%v的%v 死亡 ", card.GetUserName(), card.heroCfg.Name))

	} else {
		card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v 死亡，伤害来源为%v", card.GetOwnInfo(), spellCard.GetOwnInfo()))
	}

	card.owner.game.SendSeatMsg(func(seatId int32) proto.Message {
		return &cmsg.SSyncCard{
			SeatId:    card.GetPlayer().GetSeatID(),
			Card:      card.ToDef(seatId),
			SpellCard: spellCard.ID(),
			Changes: []*cmsg.SSyncCard_Change{{
				ChangeType: cmsg.SSyncCard_Die,
			}},
		}
	})
}
func StartSetDeadAndNotify(card *Card, srcCard *Card) {

	g := card.owner.game
	ad := &ActionDataBase{}
	g.PostActData(ad)
	ad.PostActStream(func() {
		if card.skillId > 0 {
			g.GetTriggerMgr().Trigger(TriggerType_LoseSkill, nil, card)
		}

	})
	ad.PostActStream(func() {
		g.GetTriggerMgr().Trigger(TriggerType_Die, nil, card)
	})
	ad.PostActStream(func() {
		SyncChangeDie(card, srcCard)
		card.cell.RemoveCard()
	})

}

func SyncChangeBuff(card *Card, spellCard *Card, buffId int32) {

	card.owner.game.SendSeatMsg(func(seatId int32) proto.Message {
		ret := &cmsg.SSyncCard{
			SeatId:    card.GetPlayer().GetSeatID(),
			Card:      card.ToDef(seatId),
			SpellCard: spellCard.ID(),
			Changes: []*cmsg.SSyncCard_Change{{
				ChangeType: cmsg.SSyncCard_Buff,
				NewBuffId:  buffId,
			}},
		}

		return ret
	})
}

func StartExchangePos(card1 *Card, card2 *Card) {
	if card1 == nil || card2 == nil {
		return
	}

	ad := &ActionDataBase{}
	g := card1.owner.game
	g.PostActData(ad)

	ad.PostActStream(func() {
		cell1 := card1.cell
		cell2 := card2.cell

		card1.cell = cell2
		card2.cell = cell1

		cell2.Card = card1
		cell1.Card = card2

		SyncChangePos(&Position{card2.cell.Position.Row, card2.cell.Position.Col}, g, card1)
		SyncChangePos(&Position{card1.cell.Position.Row, card1.cell.Position.Col}, g, card2)
	})

	ad.PostActStream(func() {
		g.triggerMgr.Trigger(TriggerType_AfterPosChange, ad, card1, card2)
	})

}

func StartMoveToCell(card *Card, cell *Cell) {
	if card == nil {
		return
	}

	ad := &ActionDataBase{}
	g := card.owner.game
	g.PostActData(ad)

	ad.PostActStream(func() {
		oldRow := card.cell.Row
		oldCol := card.cell.Col
		card.cell = cell
		cell.Card = card
		g.board.cells[oldRow][oldCol].Card = nil
		SyncChangePos(&Position{oldRow, oldCol}, g, card)
	})

	ad.PostActStream(func() {
		g.triggerMgr.Trigger(TriggerType_AfterPosChange, ad, card)
	})

}
