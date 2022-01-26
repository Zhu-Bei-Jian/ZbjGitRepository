package gate

import (
	"github.com/golang/protobuf/proto"
	"sanguosha.com/sgs_herox/proto/gameconf"
)

type Channel struct {
	sessions map[uint64]*session
}

func NewChannel() *Channel {
	return &Channel{
		sessions: make(map[uint64]*session),
	}
}

type ChannelMg struct {
	channels map[gameconf.ChatChannelTyp]*Channel
}

func NewChannelMg() *ChannelMg {
	//return &ChannelMg{
	//	channels: map[gameconf.ChatChannelTyp]*Channel{},
	//}
	cm := &ChannelMg{}
	cm.channels = make(map[gameconf.ChatChannelTyp]*Channel)
	cm.channels[gameconf.ChatChannelTyp_ChatCTLobby] = NewChannel()
	return cm
}

func (cm *ChannelMg) SendMsg(ctype gameconf.ChatChannelTyp, message proto.Message, ignoreUser uint64, versionGE string) {
	if c, ok := cm.channels[ctype]; ok {
		for _, gateSession := range c.sessions {
			if gateSession.userid == ignoreUser || !gateSession.isLogined() {
				continue
			}
			gateSession.SendMsg(message)
			//if gateSession, ok := SessionMgrInstance.getSession(s); ok {
			//	gateSession.SendMsg(message)
			//}
		}
	}
}

func (cm *ChannelMg) AddUser(ctype gameconf.ChatChannelTyp, s *session) {
	if c, ok := cm.channels[ctype]; ok {
		c.sessions[s.ID()] = s
	}
}

func (cm *ChannelMg) DelUser(ctype gameconf.ChatChannelTyp, session uint64) {
	if c, ok := cm.channels[ctype]; ok {
		delete(c.sessions, session)
	}
}

func (cm *ChannelMg) DelUserFromAllChannels(session uint64) {
	for t := range cm.channels {
		cm.DelUser(t, session)
	}
}
