package balance

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/framework/netcluster"
)

var (
	ErrNoAvailableAI     = errors.New("ErrNoAvailableAI")
	ErrNoAvailableGame   = errors.New("ErrNoAvailableGame")
	ErrNoAvailableServer = appframe.ErrNoAvailableServer //errors.New("ErrNoAvailableServer")
	ErrServiceLoadFull   = errors.New("ErrServiceLoadFull")
)

type LoadableBalance interface {
	Available() bool
	GetServerID(msg proto.Message) (uint32, bool)
	GetServerIDForwardMsgFromSession(session uint64, msg proto.Message) (uint32, bool)
	GetServerIDForwardRawMsgFromSession(session uint64, msgid uint32) (uint32, bool)
	OnServerEvent(svrid uint32, event netcluster.SvrEvent)
	GetLoadableServer() (appframe.Server, error)
}

type LoadableService interface {
	GetLoadableServer(typ appframe.ServerType) (appframe.Server, error)
}
