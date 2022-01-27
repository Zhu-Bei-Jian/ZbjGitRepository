package appframe

import (
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe/request/protoreq"
)

// Requester 请求者
type Requester interface {
	From() Server
	Resp(proto.Message) error
}

// ListenRequest 监听请求消息, 请求必须有 Seqid int64 字段
func ListenRequest(app iListenMsg, msg proto.Message, handler func(sender Requester, req proto.Message)) {
	// 检查 msg 是否为合格的请求消息
	//if !protoreq.IsValidMsg(msg) {
	//	logrus.Panic("ListenRequest message do not have a Seqid int64 feild")
	//}
	// 监听消息
	app.ListenMsg(msg, func(sender Server,seqid int64, msg proto.Message) {
		//seqid, ok := protoreq.GetSeqid(msg)
		//if !ok {
		//	// msg 类型已经在前面检查过了, 不应该再有获取不到 seqid 的情况了
		//	logrus.Panic("Can not get request msg seqid")
		//}
		handler(&requester{
			s:     sender,
			seqid: seqid,
		}, msg)
	})
}

type requester struct {
	s     Server
	seqid int64
}

func (r *requester) From() Server {
	return r.s
}

// 响应消息必须有 Seqid int64 字段
func (r *requester) Resp(msg proto.Message) error {
	//protoreq.SetSeqid(msg, r.seqid)
	//return r.s.SendMsg()
	return r.s.SendMsgWithSeqID(msg,r.seqid)
}

// NewErrorResponse 创建一个通用的错误响应消息, 该消息会以 error 的形式返回给请求者的回调函数, 帮助简化错误处理
func NewErrorResponse() *protoreq.ErrCode {
	return new(protoreq.ErrCode)
}

// CheckErrorRespone 检查错误是否为 ErrorResponse
func CheckErrorRespone(err error) (*protoreq.ErrCode, bool) {
	ec, ok := err.(*protoreq.ErrCode)
	return ec, ok
}

// ListenRequestSugar 为消息监听处理提供便利.
// app 参数可以是 *Application 或 *GateApplication 对象.
// reqHandler 必须是 func(sender Requester, msg *FooMsg) 形式的函数.
func ListenRequestSugar(app iListenMsg, reqHandler interface{}) {
	v := reflect.ValueOf(reqHandler)

	// type check.
	if v.Type().NumIn() != 2 {
		logrus.Panic("ListenRequestSugar handler params num wrong")
	}
	var tempSender Requester
	if v.Type().In(0) != reflect.TypeOf(&tempSender).Elem() {
		logrus.Panic("ListenRequestSugar handler num in 0 is not Requester")
	}

	msg := reflect.New(v.Type().In(1)).Elem().Interface().(proto.Message)
	ListenRequest(app, msg, func(sender Requester, msg proto.Message) {
		v.Call([]reflect.Value{reflect.ValueOf(sender), reflect.ValueOf(msg)})
	})
}
