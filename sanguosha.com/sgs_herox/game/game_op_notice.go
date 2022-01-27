package game

import (
	"github.com/golang/protobuf/proto"
	"sanguosha.com/sgs_herox/proto/cmsg"
)

func GetOpMsgList(opType cmsg.SNoticeOp_OpType) []proto.Message {
	var msgList []proto.Message
	switch opType {
	case cmsg.SNoticeOp_ActionStart:
		// 空闲出牌阶段
		msgList = []proto.Message{(*cmsg.CReqAct)(nil), (*cmsg.CReqCancelCurOpt)(nil)}
	case cmsg.SNoticeOp_SelectCard:
		msgList = []proto.Message{(*cmsg.CReqCancelCurOpt)(nil), (*cmsg.CReqOpt)(nil)}
	case cmsg.SNoticeOp_SelectCamp:
		msgList = []proto.Message{(*cmsg.CReqCancelCurOpt)(nil), (*cmsg.CReqOpt)(nil)}
	}
	return msgList
}

func NoticeOpActionStart(p *Player, cb func(timeout bool)) {
	action := cmsg.SNoticeOp_ActionStart
	waitSec := p.game.phaseLeftSec()
	p.GetGame().NoticeOpCommon(p, action, waitSec, GetOpMsgList(action), cb)
}
