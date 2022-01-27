package game

import (
	"github.com/golang/protobuf/proto"
	"sanguosha.com/baselib/appframe"
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

type User struct {
	userId uint64
	roomId uint32

	seatId int32

	session appframe.Session

	isRobot   bool
	userBrief *gamedef.UserBrief

	disconnected bool
}

func newUser(userId uint64, seatId int32, userBrief *gamedef.UserBrief, session appframe.Session) *User {
	return &User{
		userId:    userId,
		seatId:    seatId,
		session:   session,
		isRobot:   false,
		userBrief: userBrief,
	}
}

func (p *User) isConnected() bool {
	return !p.disconnected
}

func (p *User) SendMsg(msg proto.Message) {
	if !p.isConnected() {
		return
	}
	p.session.SendMsg(msg)
}

func (p *User) OnReConnect(session appframe.Session) {
	p.disconnected = false
	p.session = session
}

func (p *User) OnDisconnect() {
	p.disconnected = true
	p.session = nil
}

func (p *User) SessionID() appframe.SessionID {
	if p.session == nil {
		return appframe.SessionID{}
	}
	return p.session.ID()
}
