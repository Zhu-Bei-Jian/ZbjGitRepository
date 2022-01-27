package game

import (
	"github.com/golang/protobuf/proto"
	"sanguosha.com/sgs_herox/proto/cmsg"
)

type ActionSelectCamp struct {
	ActionDataBase
	callback func(seatId int32) bool
}

func NewActionSelectCamp(game *GameBase, callback func(seatId int32) bool) *ActionSelectCamp {
	ac := &ActionSelectCamp{}
	ac.game = game
	ac.callback = callback
	return ac
}

func (ad *ActionSelectCamp) OnMessage(p *Player, msg proto.Message) bool {
	switch msg.(type) {
	case *cmsg.CReqOpt:
		resp := &cmsg.SRespOpt{}
		req := msg.(*cmsg.CReqOpt)
		if req.OpType != cmsg.CReqOpt_SelectCamp {
			return false
		}
		defer func() {
			if resp.ErrCode != 0 {
				p.SendMsg(resp)
			}
		}()
		seatId := req.SeatId

		if ad.callback != nil {
			ok := ad.callback(seatId)
			if !ok {
				return false
			}
		}

		p.game.StopWaiting()
		return true
	default:
		return false
	}
}

//func NoticeOpSelectCamp(opPlayer *Player, spellCard *Card, callback func(seatId int32) bool) {
//
//	game := opPlayer.GetGame()
//
//	now := time.Now().Unix()
//	opEndTime := now + opTime
//	msg := &cmsg.SNoticeOp{
//		OpType:    cmsg.SNoticeOp_SelectCamp,
//		OpSeatId:  opPlayer.GetSeatID(),
//		TargetPos: nil,
//		OpEndTime: opEndTime,
//		SpellCard: spellCard.ToDef(-1),
//		Data: &cmsg.SNoticeOp_Data{
//			OptCount:      1,
//			SelectCardTyp: 0,
//		},
//	}
//
//	game.Send2All(msg)
//
//	g := opPlayer.game
//	ac := NewActionSelectCamp(g, callback)
//	g.PostActData(ac)
//
//	ac.PostActStream(func() {
//		g.StartWaiting([]core.Player{opPlayer}, opTime, msg, GetOpMsgList(cmsg.SNoticeOp_SelectCamp), func(timeout bool) {
//			if timeout {
//				if ac.callback != nil {
//					logrus.Debug("军营选择时间 已到，默认seatId=0")
//					ac.callback(0)
//				}
//			}
//		})
//	})
//}
