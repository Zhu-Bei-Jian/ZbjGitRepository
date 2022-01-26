package game

import "sanguosha.com/baselib/appframe"

type PlayerManager struct {
	session2Player map[appframe.SessionID]*Player
	userId2Player  map[uint64]*Player
}

func newPlayerManager() *PlayerManager {
	return &PlayerManager{
		session2Player: make(map[appframe.SessionID]*Player),
		userId2Player:  make(map[uint64]*Player),
	}
}

func (pm *PlayerManager) add(p *Player) {
	session := p.GetUser().SessionID()
	if session.SvrID != 0 {
		pm.session2Player[session] = p
	}
	pm.userId2Player[p.user.userId] = p
}

func (pm *PlayerManager) del(p *Player) {
	session := p.GetUser().SessionID()
	if session.SvrID != 0 {
		pm.delSession(session)
	}
	pm.delUserId(p.user.userId)
}

func (pm *PlayerManager) addSession(session appframe.SessionID, p *Player) {
	pm.session2Player[session] = p
}

func (pm *PlayerManager) delSession(session appframe.SessionID) {
	delete(pm.session2Player, session)
}
func (pm *PlayerManager) delUserId(userId uint64) {
	delete(pm.userId2Player, userId)
}

func (pm *PlayerManager) findPlayerBySessionID(sessionID appframe.SessionID) (*Player, bool) {
	p, exist := pm.session2Player[sessionID]
	return p, exist
}

func (pm *PlayerManager) findPlayerByUserId(userId uint64) (*Player, bool) {
	p, exist := pm.userId2Player[userId]
	return p, exist
}
