package balance

import (
	"github.com/golang/protobuf/proto"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/framework/netcluster"
)

//根据配置的负载限制、负载开关来选择服务器 也实现了appframe.LoadBalance的接口
type StateBalance struct {
	ls  LoadableService
	app *appframe.Application
	typ appframe.ServerType
}

func NewStateBalance(ls LoadableService, app *appframe.Application, typ appframe.ServerType) *StateBalance {
	return &StateBalance{ls, app, typ}
}

func (b *StateBalance) Available() bool {
	_, e := b.GetLoadableServer()
	if e != nil {
		return false
	}
	return true
}
func (b *StateBalance) GetServerID(msg proto.Message) (uint32, bool) {
	s, e := b.GetLoadableServer()
	if e != nil {
		return 0, false
	}
	return s.ID(), true
}
func (b *StateBalance) GetServerIDForwardMsgFromSession(session uint64, msg proto.Message) (uint32, bool) {
	return b.GetServerID(msg)
}
func (b *StateBalance) GetServerIDForwardRawMsgFromSession(session uint64, msgid uint32) (uint32, bool) {
	return b.GetServerID(nil)
}
func (b *StateBalance) OnServerEvent(svrid uint32, event netcluster.SvrEvent) {

}
func (b *StateBalance) GetLoadableServer() (appframe.Server, error) {
	return b.ls.GetLoadableServer(b.typ)
}
