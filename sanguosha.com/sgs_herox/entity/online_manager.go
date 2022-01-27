package entity

import (
	"errors"
	"sanguosha.com/baselib/appframe"
	"sync"

	"sanguosha.com/sgs_herox/entity/userdb"
)

var (
	errNoSession = errors.New("errNoSession")
)

type onlineManager struct {
	userIds        map[uint64]*userdb.User
	MapSessionUser map[appframe.SessionID]uint64
	rw             sync.RWMutex
}

func newUserOnlineManager() *onlineManager {
	m := new(onlineManager)
	m.userIds = make(map[uint64]*userdb.User)
	m.MapSessionUser = make(map[appframe.SessionID]uint64)
	return m
}

func (m *onlineManager) addUser(u *userdb.User, sid appframe.SessionID) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.userIds[u.ID()] = u
	m.MapSessionUser[sid] = u.ID()
}

func (m *onlineManager) removeUser(userId uint64, sid appframe.SessionID) {
	m.rw.Lock()
	defer m.rw.Unlock()
	delete(m.userIds, userId)
	delete(m.MapSessionUser, sid)
}

func (m *onlineManager) AllUserIds() (userIds []uint64) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	for userId, _ := range m.userIds {
		userIds = append(userIds, userId)
	}
	return
}

func (m *onlineManager) execByEveryUser(f func(u *userdb.User, err error)) {
	ids := m.AllUserIds()
	for _, userid := range ids {
		UserDBInstance.GetUser(userid, f)
	}
}
