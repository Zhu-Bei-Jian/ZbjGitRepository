package gate

import (
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/cmsg"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"time"

	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox"
)

type session struct {
	appframe.GateSession
	ticket  string
	userid  uint64
	account string

	state sessionState

	machine  sessionStateMachine
	timer    *time.Timer
	cancel   func()
	bindSvrs map[appframe.ServerType]appframe.Server

	authInfo *gamedef.AuthInfo
	ip       string
}

func newSession() *session {
	s := new(session)
	s.bindSvrs = make(map[appframe.ServerType]appframe.Server)
	return s
}

func (s *session) isAuthed() bool {
	if s.state >= stateLogining {
		return true
	}
	return false
}
func (s *session) isLogined() bool {
	if s.state >= stateLogined {
		return true
	}
	return false
}
func (s *session) setClosed() {
	if s.state == stateClosed {
		return
	}
	s.state = stateClosed
	s.machine = nil
	if s.timer != nil {
		s.timer.Stop()
	}
	if s.cancel != nil {
		s.cancel()
	}
}
func (s *session) event(e interface{}) {
	if s.state == stateClosed {
		return
	}
	m := s.machine(s, e)
	if m != nil {
		s.machine = m
	}
}

const calculateOnlineTime = 20 * time.Second
const calculateOnlineTime1s = 1 * time.Second

type sessionManager struct {
	sid2session  map[uint64]*session
	uid2session  map[uint64]*session
	stateMachine sessionStateMachine
	tick         *time.Ticker
	onlineCount  int32
	offlineCount int32
}

func (sm *sessionManager) addSession(gateSession appframe.GateSession) *session {
	s := newSession()
	s.GateSession = gateSession
	s.state = stateUnAuth
	s.machine = sm.stateMachine
	sm.sid2session[s.ID()] = s

	ChannelMgInstance.AddUser(gameconf.ChatChannelTyp_ChatCTLobby, s)
	// 超时未发起认证请求则关闭连接.
	s.timer = gameutil.SafeCallAfter(waitAuthTimeout, s.Close)
	return s
}
func (sm *sessionManager) removeSession(sid uint64) {
	s, ok := sm.sid2session[sid]
	if ok {
		s.setClosed()
		us, exist := sm.getSessionByUserID(s.userid)
		if exist && s.ID() == sid {
			delete(sm.uid2session, us.userid)
		}
		delete(sm.sid2session, sid)
		ChannelMgInstance.DelUserFromAllChannels(s.ID())
	}
}

func (sm *sessionManager) execByEveryUser(f func(uid uint64, s *session)) {
	for uid, session := range sm.uid2session {
		f(uid, session)
	}
}

func (sm *sessionManager) getSessionByUserID(userID uint64) (*session, bool) {
	if s, ok := sm.uid2session[userID]; ok {
		return s, true
	}
	return nil, false
}

func (p *sessionManager) addUserOnLoginSuccess(sid uint64) {
	if s, ok := p.sid2session[sid]; ok {
		s.state = stateLogined
		p.uid2session[s.userid] = s
	}
}

func (sm *sessionManager) getSession(sid uint64) (*session, bool) {
	s, ok := sm.sid2session[sid]
	return s, ok
}

func (p *sessionManager) execByEverySession(f func(*session)) {
	for _, s := range p.uid2session {
		f(s)
	}
}

func (p *sessionManager) close() {
	if p.tick != nil {
		p.tick.Stop()
	}
}

func initSessionManager(app *appframe.GateApplication) *sessionManager {
	sm := new(sessionManager)
	sm.sid2session = make(map[uint64]*session)
	sm.uid2session = make(map[uint64]*session)
	sm.stateMachine = initSessionStateMachine(app, sm)

	// session 的心跳已经由底层维护了.

	// 监听 session 的开启和关闭事件.
	app.ListenSessionEvent(func(sid uint64) {
		sm.onlineCount++
		sm.addSession(app.GetSession(sid))
	}, func(sid uint64) {
		sm.offlineCount++
		if s, ok := sm.getSession(sid); ok {
			if s.isAuthed() {
				// 通知 lobby 玩家断开连接.
				app.GetService(sgs_herox.SvrTypeLobby).SendMsg(&smsg.NoticeSessionClosed{
					Session: sid,
				})
				s.setClosed()
			}
			sm.removeSession(sid)
		}
	})

	// 关闭 session 指令, 发生在顶号时.
	appframe.ListenMsgSugar(app, func(sender appframe.Server, msg *smsg.LoGaNtfCloseSession) {
		if s, ok := sm.getSession(msg.Session); ok {
			logrus.WithFields(logrus.Fields{
				"session": s.ID(),
				"userid":  s.userid,
			}).Info("CMD close user session")
			s.SendMsg(&cmsg.SNoticeLogout{Reason: msg.Reason, Msg: msg.Msg})
			s.Close()
		}
	})

	// 绑定 session 到 具体服务节点的, 用于消息路由
	appframe.ListenMsgSugar(app, func(sender appframe.Server, msg *smsg.BindSessionToServer) {
		if s, ok := sm.getSession(msg.Session); ok {
			if msg.Svrid != 0 {
				logrus.WithFields(logrus.Fields{
					"session": s.ID(),
					"userid":  s.userid,
					"svrtype": msg.SvrType,
					"svrid":   msg.Svrid,
				}).Debug("Bind session to server")
				s.bindSvrs[appframe.ServerType(msg.SvrType)] = app.GetServer(msg.Svrid)
				if msg.SvrType == uint32(sgs_herox.SvrTypeGame) {
					ChannelMgInstance.DelUser(gameconf.ChatChannelTyp_ChatCTLobby, msg.Session)
				}
			} else {
				delete(s.bindSvrs, appframe.ServerType(msg.SvrType))
				if msg.SvrType == uint32(sgs_herox.SvrTypeGame) {
					ChannelMgInstance.AddUser(gameconf.ChatChannelTyp_ChatCTLobby, s)
				}
				logrus.WithFields(logrus.Fields{
					"session": s.ID(),
					"userid":  s.userid,
					"svrtype": msg.SvrType,
				}).Debug("Unbind session to server")
			}
		}
	})

	return sm
}
