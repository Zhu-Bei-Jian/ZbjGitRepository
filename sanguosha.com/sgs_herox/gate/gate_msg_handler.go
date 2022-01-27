package gate

import (
	"sanguosha.com/baselib/framework/netframe"
	"sanguosha.com/sgs_herox/proto/cmsg"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"

	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/network/msgprocessor"
	"sanguosha.com/sgs_herox"
)

func RegisterBroadcastMsg() {
	msgprocessor.RegisterMessage((*cmsg.SNoticeChatMsg)(nil))
}

func initGateMsgHandler(app *appframe.GateApplication) {

	logrus.Info("initGateMsgHandler !!!!!!!!!")

	appframe.ListenGateSessionMsgSugar(app, CReqChat)

	appframe.ListenMsgSugar(app, onChannelBroadcast)
	appframe.ListenMsgSugar(app, onKickUserOut)

	appframe.ListenRequestSugar(app, OnResponseTest)
	appframe.ListenMsgSugar(app, onServerNotice)
}

func initGateMsgRoute(app *appframe.GateApplication) {
	logrus.Info("initGateMsgRoute !!!!!!!!!")

	//entity
	RouteSessionRowMsgToEntity(app, (*cmsg.CReqMyData)(nil))
	RouteSessionRowMsgToEntity(app, (*cmsg.CReqGMCommand)(nil))
	RouteSessionRowMsgToEntity(app, (*cmsg.CReqCardGroups)(nil))
	RouteSessionRowMsgToEntity(app, (*cmsg.CCardGroupOpt)(nil))

	//game
	RouteSessionRowMsgToGame(app, (*cmsg.CReqGameScene)(nil))
	RouteSessionRowMsgToGame(app, (*cmsg.CReqGameExit)(nil))
	RouteSessionRowMsgToGame(app, (*cmsg.CReqAct)(nil))
	RouteSessionRowMsgToGame(app, (*cmsg.CReqOpt)(nil))
	RouteSessionRowMsgToGame(app, (*cmsg.CReqCancelCurOpt)(nil))

	//lobby
	lobby := app.GetService(sgs_herox.SvrTypeLobby)
	RouteSessionRowMsg(app, (*cmsg.CReqMyRoom)(nil), lobby)
	RouteSessionRowMsg(app, (*cmsg.CReqRoomCreate)(nil), lobby)
	RouteSessionRowMsg(app, (*cmsg.CReqRoomQuickJoin)(nil), lobby)
	RouteSessionRowMsg(app, (*cmsg.CReqRoomSettingChange)(nil), lobby)
	RouteSessionRowMsg(app, (*cmsg.CReqRoomJoin)(nil), lobby)
	RouteSessionRowMsg(app, (*cmsg.CReqRoomSeatChange)(nil), lobby)
	RouteSessionRowMsg(app, (*cmsg.CReqRoomLookerInfo)(nil), lobby)
	RouteSessionRowMsg(app, (*cmsg.CReqRoomReady)(nil), lobby)
	RouteSessionRowMsg(app, (*cmsg.CReqRoomLeave)(nil), lobby)
	RouteSessionRowMsg(app, (*cmsg.CReqRoomKick)(nil), lobby)
	RouteSessionRowMsg(app, (*cmsg.CReqRoomStartGame)(nil), lobby)
	RouteSessionRowMsg(app, (*cmsg.CReqMatch)(nil), lobby)
	RouteSessionRowMsg(app, (*cmsg.CCancelMatch)(nil), lobby)
	RouteSessionRowMsg(app, (*cmsg.CReqPing)(nil), lobby)
}

// 辅助函数, 只转发已登录的用户消息给需要的服务器.
func GetRawMsgRouterWhenLogined(f func(s *session, msgid uint32, data []byte, extParam int64)) appframe.GateSessionRawMsgRouter {
	return func(sid uint64, msgid uint32, data []byte, extParam int64) {
		s, ok := SessionMgrInstance.getSession(sid)
		if ok && s.isLogined() {
			f(s, msgid, data, extParam)
		}
	}
}

func RouteSessionRowMsg(app *appframe.GateApplication, msg proto.Message, service appframe.Service) {
	app.RouteSessionRawMsg(msg, GetRawMsgRouterWhenLogined(func(s *session, mssgid uint32, data []byte, extParam int64) {
		service.ForwardRawMsgFromSession(mssgid, data, netframe.Server_Extend{
			SessionId: s.ID(),
			UserId:    s.userid,
			ExtParam:  extParam,
		})
	}))
}

func RouteSessionRowMsgToEntity(app *appframe.GateApplication, msg proto.Message) {
	app.RouteSessionRawMsg(msg, GetRawMsgRouterWhenLogined(func(s *session, mssgid uint32, data []byte, extParam int64) {
		GetEntity(s.userid).ForwardRawMsgFromSession(mssgid, data, netframe.Server_Extend{
			SessionId: s.ID(),
			UserId:    s.userid,
			ExtParam:  extParam,
		})
	}))
}

func RouteSessionRowMsgToGame(app *appframe.GateApplication, msg proto.Message) {
	app.RouteSessionRawMsg(msg, GetRawMsgRouterWhenLogined(func(s *session, mssgid uint32, data []byte, extParam int64) {
		game, ok := s.bindSvrs[sgs_herox.SvrTypeGame]
		if !ok {
			if AppCfg.Develop {
				logrus.WithFields(logrus.Fields{
					"userid":  s.userid,
					"session": s.ID(),
				}).Warn("Route msg to game failed")
			}
			return
		}
		game.ForwardRawMsgFromSession(mssgid, data, netframe.Server_Extend{
			SessionId: s.ID(),
			UserId:    s.userid,
			ExtParam:  extParam,
		})
	}))
}

func RouteSessionRowMsgWithUserID(app *appframe.GateApplication, msg proto.Message, serverType appframe.ServerType) {
	app.RouteSessionRawMsg(msg, GetRawMsgRouterWhenLogined(func(s *session, mssgid uint32, data []byte, extParam int64) {
		shopPayMsg := &smsg.RouteMessageWithUserID{}
		shopPayMsg.UserID = s.userid
		shopPayMsg.MsgID = mssgid
		shopPayMsg.Data = data
		if newMsgID, newData, err := msgprocessor.OnMarshal(shopPayMsg); err == nil {
			err := app.GetService(serverType).ForwardRawMsgFromSession(newMsgID, newData, netframe.Server_Extend{
				SessionId: s.ID(),
				UserId:    s.userid,
				ExtParam:  extParam,
			})
			if err != nil {
				logrus.WithField("userid", s.userid).WithError(err).Error("RouteSessionRowMsgWithUserID error")
			}
		}
	}))
}

func CReqChat(sender appframe.GateSession, req *cmsg.CReqChat) {
	s, ok := SessionMgrInstance.getSession(sender.ID())
	if !ok {
		return
	}

	if !s.isLogined() {
		return
	}

	//玩家在游戏中，聊天信息会转到游戏中
	game, ok := s.bindSvrs[sgs_herox.SvrTypeGame]
	if ok && (req.Channel == gameconf.ChatChannelTyp_ChatCTRoom || req.Channel == gameconf.ChatChannelTyp_ChatCTGame) {
		game.ForwardMsgFromSession(req, netframe.Server_Extend{
			SessionId: sender.ID(),
			UserId:    s.userid,
		})
		return
	}

	if len(req.Msg) == 0 {
		AppInstance.GetService(sgs_herox.SvrTypeLobby).ForwardMsgFromSession(req, netframe.Server_Extend{
			SessionId: sender.ID(),
			UserId:    s.userid,
		})
		return
	}

	wordFilterMgr.FilterAsync(req.Msg, func(content string, err error) {
		AppInstance.Post(func() {
			if err != nil {
				sender.SendMsg(&cmsg.SRespChat{
					ErrCode:      cmsg.SRespChat_ErrWordFilterError,
					Channel:      req.Channel,
					MsgType:      req.MsgType,
					TargetUserId: req.TargetUserId,
				})
				return
			}

			req.Msg = content
			AppInstance.GetService(sgs_herox.SvrTypeLobby).ForwardMsgFromSession(req, netframe.Server_Extend{
				SessionId: sender.ID(),
				UserId:    s.userid,
			})
		})
	})
}

func sendMsgGameLogic(sender appframe.GateSession, msg proto.Message) {
	s, ok := SessionMgrInstance.getSession(sender.ID())
	if ok && s.isLogined() {
		if gamesvr, ok := s.bindSvrs[sgs_herox.SvrTypeGame]; ok {
			if gamesvr != nil {
				gamesvr.ForwardMsgFromSession(msg, netframe.Server_Extend{
					SessionId: sender.ID(),
					UserId:    s.userid,
				})
			}
		}
	}
}

func onChannelBroadcast(sender appframe.Server, req *smsg.SSChannelBroadcast) {
	if msg, e := msgprocessor.OnUnmarshal(req.MsgId, req.LogicMsg); e == nil {
		ChannelMgInstance.SendMsg(req.Channel, msg.(proto.Message), req.IgnoreUser, req.VersionGE)
	}
}

func onKickUserOut(sender appframe.Server, req *smsg.AllGaNtfKickUserOut) {
	kickAll := req.KickAll
	userIds := req.UserIds

	msg := &cmsg.SNoticeLogout{
		Reason: req.Reason,
		Msg:    "",
	}

	if kickAll {
		SessionMgrInstance.execByEveryUser(func(uid uint64, s *session) {
			if uid == 0 || s == nil {
				return
			}
			s.SendMsg(msg)
			s.Close()
		})
	} else {
		for _, userId := range userIds {
			s, ok := SessionMgrInstance.getSessionByUserID(userId)
			if !ok {
				continue
			}
			s.SendMsg(msg)
			s.Close()
		}
	}
}

func OnResponseTest(sender appframe.Requester, msg *smsg.ReqResponseTest) {
	resp := &smsg.RespResponseTest{}
	sender.Resp(resp)
}

func onServerNotice(sender appframe.Server, req *smsg.ServerNotice) {
	SessionMgrInstance.execByEveryUser(func(uid uint64, s *session) {
		if uid == 0 || s == nil {
			return
		}
		msg := &cmsg.ServerNotice{
			Msg: req.Msg,
		}
		s.SendMsg(msg)
	})
}
