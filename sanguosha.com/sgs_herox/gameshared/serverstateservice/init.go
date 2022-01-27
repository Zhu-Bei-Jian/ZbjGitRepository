package serverstateservice

import (
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameshared/balance"
	"sync"
)

type ServerStateService interface {
	GetServerInfos() *gameshared.ServerInfos
	ResetServerInfo(m map[string]string)
	OnServerInfoUpdate(callback func(serverInfos *gameshared.ServerInfos))

	WatchServer(svrType appframe.ServerType)
	NewStateBalance(svrType appframe.ServerType) *balance.StateBalance

	SetLoad(svrID uint32, load int64)
	ModLoad(svrID uint32, load int64)
	SetServerDelay(svrID uint32, svrType appframe.ServerType, delay int64)
	GetLoadInfo() string

	GetLoadableServer(svrType appframe.ServerType) (appframe.Server, error)
}

var serviceLock sync.Mutex
var serviceMap map[*appframe.Application]*serverStateService

func Init(app *appframe.Application) ServerStateService {
	return Get(app)
}
func Get(app *appframe.Application) ServerStateService {
	serviceLock.Lock()
	defer serviceLock.Unlock()

	if serviceMap == nil {
		serviceMap = make(map[*appframe.Application]*serverStateService)
	}
	service, ok := serviceMap[app]
	if ok {
		return service
	}
	service = newServerStateService(app)
	appframe.ListenMsgSugar(app, service.onSyncServerInfo)

	serviceMap[app] = service
	return service
}
