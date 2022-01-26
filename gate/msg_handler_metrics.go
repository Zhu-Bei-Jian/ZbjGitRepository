package gate

import (
	"sanguosha.com/sgs_herox/proto/smsg"

	"sanguosha.com/baselib/appframe"
)

func initMetricsMsgHandler(app *appframe.GateApplication) {
	appframe.ListenRequestSugar(app, onAdAllReqMetrics)
}

func onAdAllReqMetrics(sender appframe.Requester, req *smsg.AdAllReqMetrics) {
	resp := &smsg.AdAllRespMetrics{}
	defer sender.Resp(resp)

	resp.Metrics = append(resp.Metrics, &smsg.AdAllRespMetrics_Metrics{
		Key:   smsg.AdAllRespMetrics_OnlineCount,
		Value: int32(len(SessionMgrInstance.sid2session)),
	})
}
