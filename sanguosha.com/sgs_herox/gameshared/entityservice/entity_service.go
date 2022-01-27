package entityservice

import (
	"time"

	"github.com/golang/protobuf/proto"
	"sanguosha.com/baselib/appframe"
	"sync"
)

// MessageWithUserid ...
type MessageWithUserid interface {
	proto.Message
	GetUserid() uint64
}

// EntityService ...
type EntityService interface {
	SendMsg(msg MessageWithUserid) error
	Request(msg MessageWithUserid, cbk func(resp proto.Message, err error), timeout time.Duration) (cancel func())
	RequestCall(msg MessageWithUserid, timeout time.Duration) (proto.Message, error)
	IsLocalRequest(msg MessageWithUserid) bool
	GetServerForUser(userid uint64) appframe.Server
}

// GetServer ...
type GetServer interface {
	GetServer(svrid uint32) appframe.Server
}

// GetUserEntityID ...
type GetUserEntityID interface {
	GetUserEntityID(userid uint64) uint32
}

type entityService struct {
	app           *appframe.Application
	getServer     GetServer
	getUserEntity GetUserEntityID
}

func (s *entityService) GetServerForUser(userid uint64) appframe.Server {
	return s.getServer.GetServer(s.getUserEntity.GetUserEntityID(userid))
}
func (s *entityService) SendMsg(msg MessageWithUserid) error {
	return s.GetServerForUser(msg.GetUserid()).SendMsg(msg)
}
func (s *entityService) Request(msg MessageWithUserid, cbk func(resp proto.Message, err error), timeout time.Duration) (cancel func()) {
	return s.GetServerForUser(msg.GetUserid()).Request(msg, cbk, timeout)
}
func (s *entityService) RequestCall(msg MessageWithUserid, timeout time.Duration) (proto.Message, error) {
	return s.GetServerForUser(msg.GetUserid()).RequestCall(msg, timeout)
}
func (s *entityService) IsLocalRequest(msg MessageWithUserid) bool {
	return s.getServer.GetServer(s.getUserEntity.GetUserEntityID(msg.GetUserid())).ID() == s.app.ID()
}

var entityServiceLock sync.Mutex
var entityServiceMap map[*appframe.Application]*entityService

//var LocalRequestCacheUser func(uint64, func(*smsg.RespCacheUser, error))
//var LocalRequestCacheUserSummary func(uint64, func(*smsg.RespCacheUserSummary, error))

// NewEntityService ...
func NewEntityService(app *appframe.Application, getUserEntity GetUserEntityID) EntityService { //getServer GetServer
	entityServiceLock.Lock()
	defer entityServiceLock.Unlock()

	if entityServiceMap == nil {
		entityServiceMap = make(map[*appframe.Application]*entityService)
	}
	es, ok := entityServiceMap[app]
	if ok {
		return es
	}
	es = &entityService{
		app:           app,
		getServer:     app, //getServer,
		getUserEntity: getUserEntity,
	}
	entityServiceMap[app] = es
	return es
}
