package lobby

import (
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/cmsg"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"time"

	"github.com/golang/protobuf/proto"
)

type LoginState int8

type Action int32

const (
	Action_Invalid Action = iota
	Action_Team
	Action_TeamCancel
	Action_Match
	Action_MatchCancel
	Action_EnterReadyRoom
	Action_EnterReadyTableCancel
	Action_CreateGameFail
	Action_Game
	Action_GameQuit
)

const (
	// 未登录状态
	StateLogout LoginState = 0
	// 登录中
	StateLogining LoginState = 1
	// 已登录
	StateLogin LoginState = 2
)

type user struct {
	userid         uint64
	session        appframe.SessionID
	loginState     LoginState
	connect        bool
	disconnectTime int64

	mgr       *userManager
	userBrief *gamedef.UserBrief

	isRobot     bool
	IsMatching  bool
	MatchModeID int32
	GameStatus  gameconf.UserGameStatusTyp

	room *Room
	game *Game
}

func (p *user) setRoom(room *Room) {
	p.room = room
}

func (p *user) clearRoom() {
	p.setRoom(nil)
}

func (p *user) isInRoom() bool {
	return p.room != nil
}

func (p *user) headInfo() *gamedef.HeadInfo {
	return gameutil.UserBrief2HeadInfo(p.userBrief)
}

func (p *user) setUserBrief(brief *gamedef.UserBrief) {
	p.userBrief = brief
}

func (p *user) setUserGameStatus(status gameconf.UserGameStatusTyp) {
	p.GameStatus = status
}

func (u *user) isConnected() bool {
	return u.connect
}

func (u *user) quitRoomIfOffline() {
	if !u.isInRoom() {
		return
	}

	if u.isConnected() {
		return
	}

	u.room.quit(u, 0)
}

func (u *user) quitRoomIfMatch() {
	if !u.isInRoom() {
		return
	}
	roomId := u.room.roomId
	u.room.quit(u, 0)
	u.SendMsg(&cmsg.SNoticeRoomKick{
		KickerUserId: u.userid,
		RoomId:       roomId,
		KickType:     cmsg.SNoticeRoomKick_MatchGameEnd,
	})

}

func (p *user) SetMatchModeID(v int32) {
	p.MatchModeID = v
}

func (p *user) SetMatch(match bool, matchModeID int32) {
	p.IsMatching = match
	if match {
		p.SetMatchModeID(matchModeID)
	} else {
		p.SetMatchModeID(0)
	}
}

func (u *user) QuitMatch() {
	if u.IsMatching {
		_, err := matchMgr.quitMatch(u)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err":    err,
				"userID": u.userid,
				"modeID": u.MatchModeID,
			}).Error("退出匹配失败")
		}
	}
}

func (u *user) bindGame(g *Game) {
	u.game = g
	u.updateGateGameServer(g.node.ID())
}

func (u *user) updateGateGameServer(gameSvrId uint32) {
	AppInstance.GetServer(u.session.SvrID).SendMsg(&smsg.BindSessionToServer{
		Session: u.session.ID,
		SvrType: uint32(sgs_herox.SvrTypeGame),
		Svrid:   gameSvrId,
	})
}

func (u *user) clearGame() {
	if u.game == nil {
		return
	}
	u.game = nil
	u.updateGateGameServer(0)
}

func (u *user) isInGame() bool {
	return u.game != nil
}

func (u *user) SendMsg(msg proto.Message) {
	if !u.isConnected() {
		return
	}
	AppInstance.GetSession(u.session).SendMsg(msg)
}

func (u *user) IsConnect() bool {
	return u.connect
}

func (u *user) onDisconnect() {
	if !u.connect {
		return
	}

	u.loginState = StateLogout
	u.connect = false
	u.disconnectTime = time.Now().Unix()
	u.QuitMatch()

	userMgr.removeSession(u.session)
	u.session = appframe.SessionID{}

	u.mgr.entity.SendMsg(&smsg.LsEsLogout{
		Userid: u.userid,
		Time:   time.Now().Unix(),
	})

	if u.room != nil {
		u.room.onUserDisconnect(u)
	}
}

func (u *user) onConnect(sid appframe.SessionID) {
	if u.connect {
		if u.session == sid {
			logrus.WithFields(logrus.Fields{
				"userid":  u.userid,
				"session": sid,
			}).Warn("User reconnect with same session")
			return
		}
		u.mgr.app.GetServer(u.session.SvrID).SendMsg(&smsg.LoGaNtfCloseSession{
			Session: u.session.ID,
			Reason:  gameconf.KickUserOutReason_KUORelogin,
			Msg:     "",
		})
		u.onDisconnect()
	}

	u.connect = true
	u.session = sid
	u.mgr.addSession(sid, u)

	//玩家上线，游戏中通知玩家上线了，此时玩家可能还未进游戏场景，需要玩家请求游戏重连才算真正返回游戏
	if u.isInGame() {
		u.updateGateGameServer(u.game.node.ID())
		u.game.node.SendMsg(&smsg.UserConnect{
			Userid:  u.userid,
			Gateid:  u.session.SvrID,
			Session: u.session.ID,
		})
	}
}

func (u *user) gameInfo() *gamedef.GameInfo {
	if u.game == nil {
		return nil
	}

	return &gamedef.GameInfo{
		GameMode: u.game.gameMode,
		GameUUID: u.game.gameId,
	}
}
