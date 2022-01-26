package lobby

import (
	"sanguosha.com/sgs_herox/proto/smsg"

	"sanguosha.com/baselib/appframe"
)

func initMetricsMsgHandler(app *appframe.Application) {
	appframe.ListenRequestSugar(app, onAdAllReqMetrics)
}

func onAdAllReqMetrics(sender appframe.Requester, req *smsg.AdAllReqMetrics) {
	resp := &smsg.AdAllRespMetrics{}
	defer sender.Resp(resp)

	resp.Metrics = append(resp.Metrics, &smsg.AdAllRespMetrics_Metrics{
		Key:   smsg.AdAllRespMetrics_OnlineCount,
		Value: int32(len(userMgr.ss2user)),
	}, &smsg.AdAllRespMetrics_Metrics{
		Key:   smsg.AdAllRespMetrics_GameCount,
		Value: int32(len(gameMgr.games)),
	}, &smsg.AdAllRespMetrics_Metrics{
		Key:   smsg.AdAllRespMetrics_RoomCount,
		Value: int32(len(roomMgr.roomId2room)),
	})
}
