package game

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

//蛮夷：翻牌技：弃掉一张随机手牌
type SkillManYi struct {
	HeroSkill
}

func (ss *SkillManYi) OnFaceUp(card *Card) {
	p := card.owner
	g := card.owner.game
	if p.HandCard.Count() == 0 {
		return
	}

	//oldHandCard := p.HandCard.Clone()
	trashCard, ok := p.HandCard.RemoveRandom()
	if !ok {
		return
	}

	card.owner.game.GetCurrentPlayer().Log(fmt.Sprintf("%v 发动 蛮夷：翻牌技：弃掉一张随机手牌.此次丢弃的为%v", card.GetOwnInfo(), trashCard))

	g.SendSeatMsg(func(seatId int32) proto.Message {
		show := p.GetSeatID() == seatId
		return &cmsg.SSyncHandCard{
			ChangeTypes: cmsg.SSyncHandCard_Skill,
			SeatId:      p.GetSeatID(),
			GetCards:    []*gamedef.PoolCard{trashCard},
			HandCards:   p.HandCard.Cards(show),
			SpellCard:   card.ID(),
		}
	})
}
