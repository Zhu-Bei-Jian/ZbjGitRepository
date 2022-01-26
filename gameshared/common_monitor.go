package gameshared

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/framework/netcluster"
	"sanguosha.com/baselib/framework/netframe"
	"sanguosha.com/baselib/network/msgprocessor"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/smsg"
)

func RegisterMSGMonitorCommand(app *appframe.Application, pub func(data []byte) error) {
	app.RegisterService(sgs_herox.SvrTypeAdmin, appframe.WithLoadBalanceSingleton(app, sgs_herox.SvrTypeAdmin))

	appframe.ListenRequestSugar(app, func(sender appframe.Requester, req *smsg.AsAllReqMSGMonitor) {
		resp := &smsg.AsAllRespMSGMonitor{}
		defer sender.Resp(resp)

		switch req.Type {
		case smsg.AsAllReqMSGMonitor_All_PrintUserMsgDetail:
			monitorUserId := req.UserId
			var monitorSessId uint64
			app.RegisterIntercepter(func(msgSrc netcluster.MsgSrc, connId uint32, msgId uint32, msgData []byte, extend *netframe.Server_Extend) {
				if !(extend.UserId == monitorUserId || (monitorSessId != 0 && extend.SessionId == monitorSessId)) {
					return
				}

				if extend.SessionId > 0 {
					monitorSessId = extend.SessionId
				}

				data, err := MarshalMonitorMsg("trace", app.ID(), msgSrc, msgId, msgData)
				if err != nil {
					logrus.WithError(err).Error("MarshalMonitorMsg")
					return
				}

				pub(data)
			})
		case smsg.AsAllReqMSGMonitor_All_StopPrint:
			app.RegisterIntercepter(nil)
		default:

		}
	})
}

func MsgNameContent(msgId uint32, msgData []byte) (msgName string, i interface{}) {
	msg, err := msgprocessor.OnUnmarshal(msgId, msgData)
	if err != nil {
		return fmt.Sprintf("msgId:%d", msgId), string(msgData)
	}
	return string(proto.MessageReflect(msg.(proto.Message)).Descriptor().FullName()), msg
}

func MarshalMonitorMsg(id string, serverId uint32, msgSrc netcluster.MsgSrc, msgId uint32, msgData []byte) ([]byte, error) {
	type Data struct {
		Id         string            `json:"trace"`
		ServerId   uint32            `json:"serverId"`
		MsgSrc     netcluster.MsgSrc `json:"msgSrc"`
		MsgName    string            `json:"msgName"`
		Time       string            `json:"time"`
		MsgLen     string            `json:"msgLen"`
		MsgContent interface{}       `json:"msgContent"`
	}

	msgName, content := MsgNameContent(msgId, msgData)
	d := Data{
		Id:         id,
		ServerId:   serverId,
		MsgSrc:     msgSrc,
		MsgName:    msgName,
		Time:       gameutil.GetCurrentDateTime(),
		MsgLen:     Human_ByteLen(msgData),
		MsgContent: content,
	}

	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func MsgName(msgId uint32) (msgName string) {
	return msgprocessor.MessageName(msgId)
}

func Human_ByteLen(data []byte) string {
	l := len(data)
	if l < 1024 {
		return fmt.Sprintf("%d", l)
	}

	KB := l / 1024
	if KB < 1024 {
		return fmt.Sprintf("%dKB", KB)
	}

	MB := KB / 1024
	return fmt.Sprintf("%dMB", MB)
}
