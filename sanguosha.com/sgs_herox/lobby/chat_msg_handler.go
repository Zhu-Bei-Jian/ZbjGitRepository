package lobby

import (
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/proto/cmsg"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

func initChatMsgHandler(app *appframe.Application) {
	appframe.ListenSessionMsgSugar(app, onReqChat)
	appframe.ListenSessionMsgSugar(app, onReqPing)
}

func onReqChat(sender appframe.Session, req *cmsg.CReqChat) {
	u, ok := userMgr.findUserBySessionID(sender.ID())
	if !ok {
		return
	}

	resp := &cmsg.SRespChat{}
	defer sender.SendMsg(resp)

	chatMsg := &cmsg.SNoticeChatMsg{
		Channel:  req.Channel,
		MsgType:  req.MsgType,
		FromUser: u.userBrief,
		Msg:      req.Msg,
		MsgId:    req.MsgId,
	}
	switch req.Channel {
	case gameconf.ChatChannelTyp_ChatCTPrivate:
		targetUser, exist := userMgr.findUser(req.TargetUserId)
		if !exist {
			resp.ErrCode = cmsg.SRespChat_ErrTargetNotOnline
			return
		}
		targetUser.SendMsg(chatMsg)
	case gameconf.ChatChannelTyp_ChatCTRoom:
		if !u.isInRoom() {
			resp.ErrCode = cmsg.SRespChat_ErrNotInRoom
			return
		}

		if seatId, ok := u.room.findSeatId(u.userid); ok {
			chatMsg.FromSeat = seatId
		}

		u.room.notifyMessage(chatMsg, nil)
	default:
		resp.ErrCode = cmsg.SRespChat_ErrNotSupport
		return
	}
}

func onReqPing(sender appframe.Session, req *cmsg.CReqPing) {
	_, ok := userMgr.findUserBySessionID(sender.ID())
	if !ok {
		return
	}

	resp := &cmsg.SRespPing{
		TimeTag: req.TimeTag,
	}
	sender.SendMsg(resp)
}
