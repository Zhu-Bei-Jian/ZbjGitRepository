package game

import (
	"sanguosha.com/sgs_herox/proto/cmsg"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

//请求游戏场景,客户端发出此消息代表客户端场景已准备好
func (g *GameBase) onReqChat(p *Player, req *cmsg.CReqChat) {
	resp := &cmsg.SRespChat{}
	defer p.SendMsg(resp)

	msg := req.Msg

	//非旁观者要进行爆词检测
	if !p.IsLooker() {
		//msg = filterPlayerWord(req.Msg, p.word)
		//if msg != req.Msg {
		//	p.SendMsg(&cmsg.SNoticeChatMsg{
		//		Channel: gameconf.ChatChannelTyp_ChatCTGame,
		//		MsgType: gameconf.ChatMsgTyp_CMTSysTip,
		//		MsgId:   "baoChi",
		//	})
		//}
	}

	chatMsg := &cmsg.SNoticeChatMsg{
		Channel:  req.Channel,
		MsgType:  req.MsgType,
		FromUser: p.user.userBrief,
		FromSeat: p.seatId,
		Msg:      msg,
		MsgId:    req.MsgId,
	}

	switch req.Channel {
	case gameconf.ChatChannelTyp_ChatCTPrivate:
	case gameconf.ChatChannelTyp_ChatCTRoom:
		g.Send2All(chatMsg)
	default:
		resp.ErrCode = cmsg.SRespChat_ErrNotSupport
		return
	}
}
