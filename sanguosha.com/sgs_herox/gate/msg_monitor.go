package gate

import (
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/reflect/protoregistry"
	"reflect"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/framework/netcluster"
	"sanguosha.com/baselib/framework/netframe"
	"sanguosha.com/baselib/network/msgprocessor"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/mq"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/smsg"
	"time"
)

type UserMsgMonitor struct {
	userId    uint64
	connId    uint32
	sessionId uint64
}

func (p *UserMsgMonitor) init() {
	p.registerClientMSG()
}

func (p *UserMsgMonitor) registerClientMSG() {
	fd, err := protoregistry.GlobalFiles.FindFileByPath("cmsg/client_msg.proto")
	if err != nil {
		return
	}

	mds := fd.Messages()
	for i := mds.Len() - 1; i >= 0; i-- {
		x := mds.Get(i)
		fullName := x.FullName()
		mt, err := protoregistry.GlobalTypes.FindMessageByName(fullName)
		if err != nil {
			return
		}

		msgprocessor.RegisterMessageNameType(string(fullName), reflect.TypeOf(proto.MessageV1(mt.Zero().Interface())))
	}
}

func (p *UserMsgMonitor) RegisterUser(userId uint64, sessMgr *sessionManager) {

	p.userId = userId

	if sess, ok := sessMgr.getSessionByUserID(userId); ok {
		p.connId = sess.Session32ID()
		p.sessionId = sess.ID()
	} else {
		p.connId = 0
		p.sessionId = 0
	}
}

func (p *UserMsgMonitor) onUserLogin(userId uint64, connId uint32, sessId uint64) {
	if p.userId != userId {
		return
	}

	p.connId = connId
	p.sessionId = sessId
}

func (p *UserMsgMonitor) IsUser(connId uint32, extend *netframe.Server_Extend) bool {
	if (extend.UserId == p.userId) || (p.connId != 0 && connId == p.connId) || (p.sessionId != 0 && p.sessionId == extend.SessionId) {
		return true
	}

	return false
}

func RegisterMSGMonitorCommand(cfg *config.AppConfig, app *appframe.Application, sessMgr *sessionManager) error {
	mqNode, exist := cfg.MQNodes[config.MQNode_Monitor]
	if !exist {
		return errors.New("MQNode_Monitor not exist")
	}

	producer, err := mq.NewProducer(mq.Config{
		Open:    mqNode.Open,
		Type:    mqNode.Type,
		Address: mqNode.Address,
	})

	if err != nil {
		return err
	}
	pub := func(data []byte) {
		producer.PublishAsync(&mq.Msg{
			Topic: config.TopicMonitor,
			Data:  data,
		})
	}

	app.RegisterService(sgs_herox.SvrTypeAdmin, appframe.WithLoadBalanceSingleton(app, sgs_herox.SvrTypeAdmin))

	var gateMsgCounts map[uint32]int64
	var recordTime int64
	var startUserCount int32

	appframe.ListenRequestSugar(app, func(sender appframe.Requester, req *smsg.AsAllReqMSGMonitor) {
		resp := &smsg.AsAllRespMSGMonitor{}
		defer sender.Resp(resp)

		switch req.Type {
		case smsg.AsAllReqMSGMonitor_Gate_StartRECReq:
			gateMsgCounts = make(map[uint32]int64)
			recordTime = time.Now().Unix()
			startUserCount = int32(len(sessMgr.uid2session))
			app.RegisterIntercepter(func(msgSrc netcluster.MsgSrc, connId uint32, msgId uint32, msgData []byte, extend *netframe.Server_Extend) {
				if msgSrc != netcluster.MsgSrcIn_Client {
					return
				}
				gateMsgCounts[msgId]++
			})
		case smsg.AsAllReqMSGMonitor_Gate_StopGetRECReq:
			app.RegisterIntercepter(nil)
			userCount := int32(len(sessMgr.uid2session))

			m := make(map[string]int64, len(gateMsgCounts))
			for msgId, count := range gateMsgCounts {
				name := gameshared.MsgName(msgId)
				m[name] = count
			}

			resp.RECResult = &smsg.AsAllRespMSGMonitor_RECReqResult{
				GateServerId:   int32(app.ID()),
				Msgs:           m,
				StartUserCount: startUserCount,
				EndUserCount:   userCount,
				StartTime:      recordTime,
				EndTime:        time.Now().Unix(),
			}
			gateMsgCounts = nil
		case smsg.AsAllReqMSGMonitor_All_PrintUserMsgDetail:
			userMsgMonitor.RegisterUser(req.UserId, sessMgr)

			app.RegisterIntercepter(func(msgSrc netcluster.MsgSrc, connId uint32, msgId uint32, msgData []byte, extend *netframe.Server_Extend) {
				if !userMsgMonitor.IsUser(connId, extend) {
					return
				}

				data, err := gameshared.MarshalMonitorMsg("trace", app.ID(), msgSrc, msgId, msgData)
				if err != nil {
					logrus.WithError(err).Error("MarshalMonitorMsg")
					return
				}

				pub(data)
			})
		case smsg.AsAllReqMSGMonitor_Gate_PrintUserMsgInOut:
			userMsgMonitor.RegisterUser(req.UserId, sessMgr)

			app.RegisterIntercepter(func(msgSrc netcluster.MsgSrc, connId uint32, msgId uint32, msgData []byte, extend *netframe.Server_Extend) {
				if msgSrc != netcluster.MsgSrcIn_Client && msgSrc != netcluster.MsgSrcOut_Client {
					return
				}

				if !userMsgMonitor.IsUser(connId, extend) {
					return
				}

				data, err := gameshared.MarshalMonitorMsg("inout", app.ID(), msgSrc, msgId, msgData)
				if err != nil {
					logrus.WithError(err).Error("MarshalMonitorMsg")
					return
				}

				pub(data)
			})
		case smsg.AsAllReqMSGMonitor_Gate_PrintSlowResponse:
			slowerThan := req.SlowerThan

			app.RegisterIntercepter(func(msgSrc netcluster.MsgSrc, connId uint32, msgId uint32, msgData []byte, extend *netframe.Server_Extend) {
				switch msgSrc {
				case netcluster.MsgSrcIn_Client:
					extend.ExtParam = time.Now().UnixNano()
				case netcluster.MsgSrcOut_Client:
					lastNano := extend.ExtParam
					if lastNano == 0 {
						return
					}

					elapse := int32((time.Now().UnixNano() - lastNano) / 1e6)
					if elapse < slowerThan {
						return
					}

					msgName := gameshared.MsgName(msgId)

					detail := make(map[string]interface{})
					detail["id"] = "slowresponse"
					detail["time"] = gameutil.GetCurrentDateTime()
					detail["msgName"] = msgName
					detail["elapse"] = elapse
					detail["userId"] = extend.UserId

					data, err := json.Marshal(detail)
					if err != nil {
						return
					}

					pub(data)
				default:

				}
			})
		case smsg.AsAllReqMSGMonitor_All_StopPrint:
			app.RegisterIntercepter(nil)
		default:

		}
	})
	return nil
}
