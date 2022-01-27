package gameshared

import (
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/framework/netcluster"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/proto/smsg"
)

func RegisterCommonCommand(app *appframe.Application) {
	appframe.ListenRequestSugar(app, func(sender appframe.Requester, req *smsg.AdAllReqCloseServer) {
		resp := &smsg.AdAllRespCloseServer{}
		sender.Resp(resp)
		app.Exit()
	})

	appframe.ListenRequestSugar(app, func(sender appframe.Requester, req *smsg.AdAllReqPingServer) {
		resp := &smsg.AdAllRespPingServer{}
		sender.Resp(resp)
	})
}

func RegisterCommonServerStatus(app *appframe.Application, cbk func(errType int32, targetID uint32, isDisconnect bool)) {
	handler := commonServerStatusHandler(cbk)
	for serverType := sgs_herox.SvrTypeGate; serverType < sgs_herox.SvrTypeEnd; serverType++ {
		app.ListenServerEvent(serverType, handler)
	}
}

func commonServerStatusHandler(cbk func(int32, uint32, bool)) func(svrid uint32, event netcluster.SvrEvent) {
	return func(svrid uint32, event netcluster.SvrEvent) {
		var errType int32
		isDisconnect := false
		switch event {
		case netcluster.SvrEventQuit:
			return
		case netcluster.SvrEventDisconnect:
			isDisconnect = true
			errType = 1
		}
		cbk(errType, svrid, isDisconnect)
	}
}
