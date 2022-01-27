package game

import (
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

func SubAPAndSync(p *Player, ap int32) {
	p.SubAP(ap)
	SyncPlayerState(p, cmsg.SSyncPlayerState_AP)
}

func SetAPAndSync(p *Player, ap int32) {
	p.SetAP(ap)
	SyncPlayerState(p, cmsg.SSyncPlayerState_AP)
}

//func AddHandCardPoolAndSync(p *Player, cardIds ...int32) {
//	p.HandCardPool.AddCard(cardIds...)
//	SyncPlayerState(p, cmsg.SSyncPlayerState_CardPoolCount)
//
//}

func SyncPlayerState(player *Player, state cmsg.SSyncPlayerState_ChangeType) {
	g := player.game

	var new int32
	switch state {
	case cmsg.SSyncPlayerState_CardPoolCount:
		new = int32(player.HandCardPool.Count())
	case cmsg.SSyncPlayerState_AP:
		new = player.GetAP()
	case cmsg.SSyncPlayerState_State:
	default:

	}

	g.Send2All(&cmsg.SSyncPlayerState{
		ChangeType:   state,
		Seat:         player.seatId,
		ConnectState: player.connectState,
		Dead:         player.IsDead(),
		New:          new,
	})

}

func SyncCampChangeBuff(p *Player, spellCard *Card) {
	p.game.Send2All(&cmsg.SSyncCamp{
		Changes: []*cmsg.SSyncCamp_Change{{
			ChangeType: cmsg.SSyncCamp_Buff,
			NewBuffs:   p.BuffManager.ToDef(),
		}},
		SeatId:    p.GetSeatID(),
		SpellCard: spellCard.ID(),
	})
}

func SyncCampChangeHP(p *Player, spellCard *Card, old int32, change int32) {
	if change == 0 {
		return
	}
	new := p.GetHP()
	p.game.Send2All(&cmsg.SSyncPlayerState{
		Seat:         p.GetSeatID(),
		ConnectState: p.connectState,
		Dead:         p.IsDead(),
		HpChange: &gamedef.Change{
			Old:    old,
			New:    new,
			Change: change,
		},
		ChangeType: cmsg.SSyncPlayerState_HP,
	})

	p.game.Send2All(&cmsg.SSyncCamp{
		SpellCard: spellCard.ID(),
		Changes: []*cmsg.SSyncCamp_Change{{
			ChangeType: cmsg.SSyncCamp_HP,
			Change:     0,
			Old:        old,
			New:        p.Camp.GetHP(),
		}},
		SeatId: p.GetSeatID(),
	})
}
