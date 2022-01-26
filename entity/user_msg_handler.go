package entity

import (
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/entity/userdb"
	"sanguosha.com/sgs_herox/proto/smsg"
)

func initUserMsgHandler(app *appframe.Application) {
	appframe.ListenRequestSugar(app, onLsEsReqLogin)
	appframe.ListenMsgSugar(app, onLsEsReqLogout)
}

func onLsEsReqLogin(sender appframe.Requester, req *smsg.LsEsReqLogin) {
	UserDBInstance.GetUser(req.Userid, func(u *userdb.User, err error) {
		if err != nil {
			resp := &smsg.LsEsRespLogin{
				ErrCode: smsg.LsEsRespLogin_ErrUnknow,
				Userid:  req.Userid,
			}
			sender.Resp(resp)
			return
		}

		if u.IsOnline() {
			logrus.WithField("userId", u.ID()).Error("onLsEsReqLogin when user online")
			doLogout(u)
		}

		sid := appframe.SessionID{
			SvrID: req.GateId,
			ID:    req.Session,
		}

		u.Login(req.IP, req.AuthInfo)
		u.BindSession(sid)
		onlineMgr.addUser(u, sid)
		u.AddLogLogin()

		resp := &smsg.LsEsRespLogin{
			Userid:    u.ID(),
			UserBrief: u.Brief(),
		}
		sender.Resp(resp)
	})
}

func onLsEsReqLogout(sender appframe.Server, req *smsg.LsEsLogout) {
	UserDBInstance.GetUser(req.Userid, func(u *userdb.User, err error) {
		if err != nil {
			//logrus.WithField("userid", req.Userid).WithError(err).Error("Login entity failed")
			//rerr := appframe.NewErrorResponse()
			//rerr.Code = int32(smsg.EsRespLogin_ErrUnknow)
			//rerr.Msg = err.Error()
			//sender.Resp(rerr)
			return
		}

		if !u.IsOnline() {
			logrus.WithField("userId", u.ID()).Error("onLsEsReqLogout when user offline")
			return
		}

		doLogout(u)
	})
}

func doLogout(u *userdb.User) {
	userId := u.ID()
	sid := u.GetSessionID()
	u.UnbindSession()
	u.Logout()
	u.AddLogLogout()
	onlineMgr.removeUser(userId, sid)
}
