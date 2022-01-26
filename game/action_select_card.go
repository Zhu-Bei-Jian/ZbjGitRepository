package game

import (
	"github.com/golang/protobuf/proto"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

type ActionSelectCard struct {
	ActionDataBase
	callback    func(cards []*Card) bool
	selectCount int32
	selectType  gamedef.SelectCardType
	card        *Card //技能发起者
}

func NewActionSelectCard(game *GameBase, selectCount int32, selectType gamedef.SelectCardType, card *Card, callback func(cards []*Card) bool) *ActionSelectCard {
	ac := &ActionSelectCard{}
	ac.game = game
	ac.selectType = selectType
	ac.selectCount = selectCount
	ac.callback = callback
	ac.card = card
	return ac
}

func (ad *ActionSelectCard) OnMessage(p *Player, msg proto.Message) bool {
	switch msg.(type) {
	case *cmsg.CReqOpt:
		resp := &cmsg.SRespOpt{}
		req := msg.(*cmsg.CReqOpt)
		if req.OpType != cmsg.CReqOpt_SelectCard {
			return false
		}
		defer func() {
			if resp.ErrCode != 0 {
				p.SendMsg(resp)
			}
		}()
		//positon 去重 ，保留唯一position

		var tPos []*gamedef.Position
		for _, pos1 := range req.TargetPos {
			hasSamePos := false
			for _, pos2 := range tPos {
				if pos1.Row == pos2.Row && pos1.Col == pos2.Col {
					hasSamePos = true
					break
				}
			}
			if !hasSamePos {
				tPos = append(tPos, pos1)
			}
		}
		//检查 slectCount 与传回的position去重后的个数是否一致
		if ad.selectCount != int32(len(tPos)) {
			resp.ErrCode = cmsg.SRespOpt_ErrSelectCountNotSame
			return false
		}
		resp.ErrCode = cmsg.SRespOpt_ErrSelectType //假定选牌类型异常 ，检查之后若无问题，再将Errcode归零
		var cards []*Card
		for _, pos := range tPos {
			card, exist := ad.game.board.GetCardByPos(pos)
			if !exist {
				return false
			}
			cards = append(cards, card)
		}

		// 根据selectCardType检查牌合法性
		switch ad.selectType {
		case gamedef.SelectCardType_Enemy:
			for _, card := range cards {
				if card.owner == ad.card.owner {
					return false
				}
			}
		case gamedef.SelectCardType_MyOwn:
			for _, card := range cards {
				if card.owner != ad.card.owner {
					return false
				}
			}
		case gamedef.SelectCardType_OtherMyOwnFaceUp:
			for _, card := range cards {
				if card.owner != ad.card.owner || (card.cell.Position.Row == ad.card.cell.Row && card.cell.Position.Col == ad.card.cell.Col) || card.isBack {
					return false
				}
			}
		case gamedef.SelectCardType_OtherMyOwn:
			for _, card := range cards {
				if card.owner != ad.card.owner || (card.cell.Position.Row == ad.card.cell.Row && card.cell.Position.Col == ad.card.cell.Col) {
					return false
				}
			}
		case gamedef.SelectCardType_EnemyBack:
			for _, card := range cards {
				if card.owner == ad.card.owner || (!card.isBack) {
					return false
				}
			}
		case gamedef.SelectCardType_EnemyFaceUp:
			for _, card := range cards {
				if card.owner == ad.card.owner || (card.isBack) {
					return false
				}
			}
		case gamedef.SelectCardType_MyOwnBack:
			for _, card := range cards {
				if card.owner != ad.card.owner || (!card.isBack) {
					return false
				}
			}
		case gamedef.SelectCardType_MyOwnFaceUp:
			for _, card := range cards {
				if card.owner != ad.card.owner || (card.isBack) {
					return false
				}
			}
		case gamedef.SelectCardType_OtherMyOwnBack:
			for _, card := range cards {
				if card.owner != ad.card.owner || (card.cell.Position.Row == ad.card.cell.Row && card.cell.Position.Col == ad.card.cell.Col) || (!card.isBack) {
					return false
				}
			}
		case gamedef.SelectCardType_OtherEnemyBack:
			for _, card := range cards {
				if card.owner == ad.card.owner || (card.cell.Position.Row == ad.card.cell.Row && card.cell.Position.Col == ad.card.cell.Col) || (!card.isBack) {
					return false
				}
			}
		case gamedef.SelectCardType_NotHeavy: //非重装明将
			for _, card := range cards {
				if card.isBack || (!card.IsHeavyCard()) {
					return false
				}
			}
		case gamedef.SelectCardType_OneOtherMyOwnAndOneEnemy: //选择一名己方武将 和一名 敌方武将
			if len(cards) != 2 {
				return false
			}
			if cards[0].owner == cards[1].owner {
				return false
			}

		}
		if ad.callback != nil {
			ok := ad.callback(cards)
			if !ok {
				return false
			}
		}
		resp.ErrCode = 0 //检测类型无异常，错误码置零

		p.game.StopWaiting()

		return true
	default:
		return false
	}
}

//
//func NoticeOpSelectCard(opPlayer *Player, spellCard *Card, selectCount int32, selectType gamedef.SelectCardType, callback func(cards []*Card) bool) {
//
//	game := opPlayer.GetGame()
//
//	now := time.Now().Unix()
//	opEndTime := now + opTime
//	msg := &cmsg.SNoticeOp{
//		OpType:    cmsg.SNoticeOp_SelectCard,
//		OpSeatId:  opPlayer.GetSeatID(),
//		TargetPos: nil,
//		OpEndTime: opEndTime,
//		SpellCard: spellCard.ToDef(-1),
//		Data: &cmsg.SNoticeOp_Data{
//			OptCount:      selectCount,
//			SelectCardTyp: selectType,
//		},
//	}
//
//	game.Send2All(msg)
//
//	g := opPlayer.game
//	ac := NewActionSelectCard(g, selectCount, selectType, spellCard, callback)
//	g.PostActData(ac)
//
//	ac.PostActStream(func() {
//		g.StartWaiting([]core.Player{opPlayer}, opTime, msg, GetOpMsgList(cmsg.SNoticeOp_SelectCard), func(timeout bool) {
//			if timeout {
//				if ac.callback != nil {
//					ac.callback(nil)
//				}
//			}
//		})
//	})
//}
