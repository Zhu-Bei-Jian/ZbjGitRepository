package gate

import (
	"sanguosha.com/sgs_herox/proto/cmsg"
	"sanguosha.com/sgs_herox/proto/smsg"
	"time"

	"sanguosha.com/sgs_herox"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"net"
	"sanguosha.com/baselib/appframe"
)

const (
	// 等待用户发起验证请求的最大时限.
	waitAuthTimeout = 30 * time.Second
	reqAuthTimeout  = 15 * time.Second
	reqLoginTimeout = 60 * time.Second
)

type sessionState int

const (
	stateUnAuth sessionState = iota
	stateAuthing
	stateLogining
	stateLogined
	stateClosed sessionState = -1
)

type sessionStateMachine func(s *session, event interface{}) sessionStateMachine

func initSessionStateMachine(app *appframe.GateApplication, sm *sessionManager) sessionStateMachine {
	app.RegisterResponse((*smsg.RespAuth)(nil))
	//app.RegisterResponse((*smsg.RespGetUserID)(nil))

	lobby := app.GetService(sgs_herox.SvrTypeLobby)
	app.RegisterResponse((*smsg.GtLsRespLogin)(nil))

	// 用户的登录请求.
	appframe.ListenGateSessionMsgSugar(app, func(sender appframe.GateSession, req *cmsg.CReqLogin) {
		s, ok := sm.getSession(sender.ID())
		if ok {
			s.event(req)
		}
	})

	// 用户登陆状态机.
	var sUnAuth, sAuthing, sLogining, sLogined sessionStateMachine

	// 未认证状态.
	sUnAuth = func(s *session, e interface{}) sessionStateMachine {
		switch e := e.(type) {
		case (*cmsg.CReqLogin):
			logrus.WithFields(logrus.Fields{
				"session": s.ID(),
				"ticket":  e.Ticket,
			}).Debug("User req login")

			//msg := &cmsg.SRespLogin{
			//	ErrCode:    0,
			//	ErrData:    nil,
			//	UserId:     1000,
			//	Account:    "",
			//	UserBrief:  nil,
			//	ServerTime: 0,
			//	IsInGame:   false,
			//}
			//s.SendMsg(msg)

			addr := s.Addr()
			ip, _, _ := net.SplitHostPort(addr.String())

			s.ticket = e.Ticket

			if s.timer != nil {
				s.timer.Stop()
				s.timer = nil
			}

			s.ip = ip
			// 发起认证请求.
			auth := app.GetService(sgs_herox.SvrTypeAuth)
			s.cancel = auth.Request(&smsg.ReqAuth{
				Ticket:    e.Ticket,
				LoginType: e.LoginType,
				Ip:        ip,
				Version:   e.Version,
			}, func(resp proto.Message, err error) {
				s.cancel = nil
				if err != nil {
					s.event(err)
				} else {
					s.event(resp)
				}
			}, reqAuthTimeout)

			s.state = stateAuthing
			return sAuthing
		}
		return nil
	}

	// 正在认证状态.
	sAuthing = func(s *session, e interface{}) sessionStateMachine {
		switch e := e.(type) {
		case error:
			logrus.Error(e)
			resp := &cmsg.SRespLogin{}
			resp.ErrCode = cmsg.SRespLogin_ErrSystem
			s.SendMsg(resp)
			s.state = stateUnAuth
			// 认证失败则关闭连接.
			s.Close()
			return sUnAuth
		case (*smsg.RespAuth):
			if e.ErrCode != 0 {
				resp := &cmsg.SRespLogin{}
				resp.ErrCode = cmsg.SRespLogin_ErrCode(e.ErrCode + 1) //后续新增保证errcode的对应关系
				s.SendMsg(resp)
				s.state = stateUnAuth
				// 认证失败则关闭连接.
				s.Close()
				return sUnAuth
			}

			s.userid = e.Userid
			s.account = e.Account
			s.authInfo = e.AuthInfo

			s.cancel = lobby.Request(&smsg.GtLsReqLogin{
				Userid:   s.userid,
				Session:  s.ID(),
				IP:       s.ip,
				AuthInfo: s.authInfo,
			}, func(resp proto.Message, err error) {
				s.cancel = nil
				if err != nil {
					s.event(err)
				} else {
					s.event(resp)
				}
			}, reqLoginTimeout)

			s.state = stateLogining
			return sLogining
		}

		return nil
	}

	// 已通过认证, 正在登录状态.
	sLogining = func(s *session, e interface{}) sessionStateMachine {
		log := logrus.WithFields(logrus.Fields{
			"session": s.ID(),
			"userid":  s.userid,
		})

		switch e := e.(type) {
		case error:
			var logAsWarn bool
			resp := &cmsg.SRespLogin{ErrCode: cmsg.SRespLogin_ErrSystem}
			s.SendMsg(resp)

			if logAsWarn {
				log.WithError(e).Warn("User login fail")
			} else {
				log.WithError(e).Error("User login fail")
			}

			s.state = stateUnAuth
			s.Close()
			return sUnAuth

		case (*smsg.GtLsRespLogin):
			if e.ErrCode != 0 {
				resp := &cmsg.SRespLogin{ErrCode: cmsg.SRespLogin_ErrSystem}
				s.SendMsg(resp)
				s.state = stateUnAuth
				s.Close()
				return sUnAuth
			}

			log.Debug("User login succ")

			SessionMgrInstance.addUserOnLoginSuccess(s.ID())

			now := time.Now().Unix()
			resp := &cmsg.SRespLogin{
				ErrCode:    0,
				ErrData:    nil,
				UserId:     s.userid,
				Account:    s.account,
				UserBrief:  e.UserBrief,
				IsInGame:   e.GameInfo != nil,
				ServerTime: now,
				ServerCfg:  e.ServerCfg,
			}
			userMsgMonitor.onUserLogin(s.userid, s.Session32ID(), s.ID())
			s.SendMsg(resp)
			s.SendMsg(&cmsg.SSyncServerTime{ServerTime: now})
			return sLogined
		}
		return nil
	}

	// 已登录状态.
	sLogined = func(s *session, e interface{}) sessionStateMachine {
		return nil
	}

	return sUnAuth
}
