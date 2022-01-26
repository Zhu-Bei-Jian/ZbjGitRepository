package appframe

import (
	"time"

	"github.com/golang/protobuf/proto"
	"sanguosha.com/baselib/framework/netframe"
)

// ServerType 服务节点的类型
type ServerType uint32

// Server 一个具体的服务器节点, 有唯一的 ID 号标识.
type Server interface {
	// 服务唯一标识符.
	ID() uint32
	// 继承自 Service 接口
	Service
}

// MsgHandler 消息处理函数声明
type MsgHandler func(sender Server, seqId int64,msg proto.Message)

type server struct {
	id  uint32
	app *Application
}

func (s *server) ID() uint32 {
	return s.id
}
func (s *server) Type() ServerType {
	typ, _ := s.app.slave.GetServerType(s.id)
	return ServerType(typ)
}
func (s *server) Available() bool {
	return s.app.slave.IsServerAvailable(s.id)
}
func (s *server) SendMsg(msg proto.Message) error {
	return s.app.slave.SendServerMsg(msg,netframe.Server_Extend{ServerId:s.id})
}

func (s *server) SendMsgWithSeqID(msg proto.Message,seqId int64) error {
	return s.app.slave.SendServerMsg(msg,netframe.Server_Extend{ServerId:s.id,SeqId: seqId})
}

func (s *server) Request(msg proto.Message, cbk func(proto.Message, error), timeout time.Duration) (cancel func()) {
	return s.app.reqc.Req(s.SendMsgWithSeqID, msg, cbk, timeout)
}
func (s *server) RequestCall(msg proto.Message, timeout time.Duration) (proto.Message, error) {
	return s.app.reqc.Call(s.SendMsgWithSeqID, msg, timeout)
}
func (s *server) ForwardMsgFromSession(msg proto.Message,extend netframe.Server_Extend) error {
	extend.ServerId= s.id
	return s.app.slave.SendServerMsg(msg,extend)
}
func (s *server) ForwardRawMsgFromSession(msgid uint32, data []byte,extend netframe.Server_Extend) error {
	extend.ServerId = s.id
	return s.app.slave.SendServerBytes(msgid, data,extend)
}

func (s *server)ReqSugar(msg proto.Message,cbk interface{},duration time.Duration) (cancel func()){
	return s.app.reqc.ReqSugar(s.SendMsgWithSeqID,msg,cbk,duration)
}

func (s *server)CallSugar(msg proto.Message,resp proto.Message,duration time.Duration) (error){
	return s.app.reqc.CallSugar(s.SendMsgWithSeqID,msg,resp,duration)
}
