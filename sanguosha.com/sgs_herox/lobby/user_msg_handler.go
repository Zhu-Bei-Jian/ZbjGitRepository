package lobby

import (
	"github.com/golang/protobuf/proto"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/network/msgprocessor"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/smsg"
	"time"
)

func initUserMsgHandler(app *appframe.Application) {
	app.RegisterResponse((*smsg.RespGameRemove)(nil))
	app.RegisterResponse((*smsg.LsEsRespLogin)(nil))
	app.RegisterResponse((*smsg.LoGaRespNewPVPGame)(nil))

	appframe.ListenRequestSugar(app, OnReqQueryBindGame)
	appframe.ListenRequestSugar(app, onGtLsReqLogin)

	appframe.ListenMsgSugar(app, SSNoticeUser)

	appframe.ListenMsgSugar(app, onNoticeSessionClosed)
	appframe.ListenMsgSugar(app, onSyncUserBrief)
	appframe.ListenSessionMsgSugar(app, onReqMatch)
	appframe.ListenSessionMsgSugar(app, onCancelMatch)
}

//登录
func onGtLsReqLogin(sender appframe.Requester, req *smsg.GtLsReqLogin) {
	userId := req.Userid

	sid := appframe.SessionID{SvrID: sender.From().ID(), ID: req.Session}
	u, ok := userMgr.findUser(userId)
	if !ok {
		u = userMgr.addUser(userId)
	}

	if u.loginState == StateLogining {
		sender.Resp(&smsg.GtLsRespLogin{
			ErrCode: smsg.GtLsRespLogin_ErrSystem,
		})
		return
	}

	u.onConnect(sid)

	u.loginState = StateLogining
	EntityInstance.Request(&smsg.LsEsReqLogin{
		Userid:   userId,
		Time:     time.Now().Unix(),
		IP:       req.IP,
		GateId:   sender.From().ID(),
		Session:  req.Session,
		AuthInfo: req.AuthInfo,
	}, func(msg proto.Message, err error) {
		resp := &smsg.GtLsRespLogin{}
		defer sender.Resp(resp)

		if err != nil {
			u.onDisconnect()
			resp.ErrCode = smsg.GtLsRespLogin_ErrSystem
			return
		}

		respCall, _ := msg.(*smsg.LsEsRespLogin)
		if respCall.ErrCode != 0 {
			u.onDisconnect()
			resp.ErrCode = smsg.GtLsRespLogin_ErrSystem
			return
		}

		userBrief := respCall.UserBrief

		resp.GameInfo = u.gameInfo()
		resp.UserBrief = userBrief
		resp.ServerCfg = &gamedef.ServerConfig{
			MinPlayerCount: gameConfig.Global.RoomSettingMinPlayerCount,
			MaxPlayerCount: gameConfig.Global.RoomSettingMaxPlayerCount,
		}
		u.loginState = StateLogin
		u.setUserBrief(userBrief)
	}, time.Second*30)
}

//断线事件
func onNoticeSessionClosed(sender appframe.Server, notice *smsg.NoticeSessionClosed) {
	sid := appframe.SessionID{SvrID: sender.ID(), ID: notice.Session}
	u, ok := userMgr.findUserBySessionID(sid)
	if !ok {
		return
	}
	u.onDisconnect()
}

func onSyncUserBrief(sender appframe.Server, msg *smsg.EsAllNtfUserBrief) {
	u, ok := userMgr.findUser(msg.Userid)
	if !ok {
		return
	}
	u.setUserBrief(msg.UserBrief)
}

func OnReqQueryBindGame(sender appframe.Requester, req *smsg.ReqUserGameInfo) {
	if u, ok := userMgr.findUser(req.Userid); ok && u.isInGame() {
		sender.Resp(&smsg.RespUserGameInfo{
			Userid:   u.userid,
			Svrid:    u.game.node.ID(),
			GameUUID: string(u.game.gameId),
			GameMode: u.game.gameMode,
		})
		return
	}

	sender.Resp(&smsg.RespUserGameInfo{Userid: req.Userid})
}

func SSNoticeUser(sender appframe.Server, req *smsg.SSNoticeUser) {
	u, ok := userMgr.findUser(req.UserId)
	if !ok {
		return
	}

	if msg, e := msgprocessor.OnUnmarshal(req.MsgId, req.LogicMsg); e == nil {
		u.SendMsg(msg.(proto.Message))
	}
}

func onReqMatch(sender appframe.Session, req *cmsg.CReqMatch) {
	u, ok := userMgr.findUserBySessionID(sender.ID())
	if !ok {
		return
	}
	resp := &cmsg.SRespMatch{
		Model: int32(req.Mode),
	}
	defer sender.SendMsg(resp)

	req.Mode = 1

	err := matchMgr.joinMatch(MatchQueueTyp_MQTWaitQueue, u.userid, int32(req.Mode), 1)
	if err != nil {
		resp.ErrCode = cmsg.SRespMatch_Failed
		return
	}
	resp.ErrCode = cmsg.SRespMatch_Success
}

func onCancelMatch(sender appframe.Session, req *cmsg.CCancelMatch) {
	u, ok := userMgr.findUserBySessionID(sender.ID())
	if !ok {
		return
	}
	resp := &cmsg.SRespCancelMatch{
		ErrCode: cmsg.SRespCancelMatch_Success,
	}
	defer sender.SendMsg(resp)

	ok, _ = matchMgr.quitMatch(u)
	if !ok {
		resp.ErrCode = cmsg.SRespCancelMatch_Failed
		return
	}
}
