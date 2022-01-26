package protoreq

import (
	"errors"
	"time"
	"reflect"
	"fmt"

	"github.com/golang/protobuf/proto"
	"sanguosha.com/baselib/appframe/request"
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/util"
)

// Client 用于发起消息请求, 接受响应
type Client struct {
	*request.Client
}

// NewClient 创建 Client
func NewClient(worker func(func())) *Client {
	return &Client{
		Client: request.NewClient(worker),
	}
}

// Req 发起请求, 异步回调, 请求消息必须有 Seqid int64 字段
func (c *Client) Req(sendMsg func(proto.Message,int64) error, msg proto.Message, cbk func(resp proto.Message, err error), timeout time.Duration) (cancel func()) {
	return c.Client.Req(func(seqid int64) error {
		//SetSeqid(msg, seqid)
		return sendMsg(msg,seqid)
	}, func(resp interface{}, err error) {
		if err == nil {
			msg, ok := resp.(proto.Message)
			if ok {
				cbk(msg, err)
			} else {
				cbk(nil, errors.New("resp not a proto.Message"))
			}
		} else {
			cbk(nil, err)
		}
	}, timeout)
}

// Call 阻塞调用, 消息必须有 Seqid int64 字段
func (c *Client) Call(sendMsg func(proto.Message,int64) error, msg proto.Message, timeout time.Duration) (proto.Message, error) {
	resp, err := c.Client.Call(func(seqid int64) error {
		//SetSeqid(msg, seqid)
		return sendMsg(msg,seqid)
	}, timeout)
	if err != nil {
		return nil, err
	}
	msg, ok := resp.(proto.Message)
	if !ok {
		return nil, errors.New("resp not a proto.Message")
	}
	return msg, nil
}

// OnResp 响应消息, 消息必须有 Seqid int64 字段
func (c *Client) OnResp(msg proto.Message,seqid int64) {
	//seqid, ok := GetSeqid(msg)
	//if !ok {
	//	return
	//}
	if err, ok := msg.(*ErrCode); ok {
		c.Client.OnErr(seqid, err)
	} else {
		c.Client.OnResp(seqid, msg)
	}
}

// Req 发起请求, 异步回调, 请求消息必须有 Seqid int64 字段
func (c *Client) ReqSugar(sendMsg func(proto.Message,int64) error, msg proto.Message, cb interface{}, timeout time.Duration) (cancel func()) {
	log:=logrus.WithField("msg",msg.String())

	cbType := reflect.TypeOf(cb)
	if cbType.Kind() != reflect.Func {
		log.Error("ReqSugar cb not a func")
		return nil
	}
	cbValue := reflect.ValueOf(cb)
	numArgs := cbType.NumIn()
	if numArgs != 2 {
		log.Error("ReqSugar cb param num args !=2")
		return nil
	}
	args0 := cbType.In(0)
	if args0.Kind() != reflect.Ptr {
		log.Error("ReqSugar cb param args0 not ptr")
		return nil
	}
	//TODO 严格检查参数类型
	args1 := cbType.In(1)

	return c.Client.Req(func(seqid int64) error {
		//SetSeqid(msg, seqid)
		return sendMsg(msg,seqid)
	}, func(resp interface{}, err error) {
		oV := make([]reflect.Value, 2)

		if err == nil {
			if reflect.TypeOf(resp)!=args0{
				oV[0] =reflect.New(args0).Elem()
				oV[1] = reflect.ValueOf(fmt.Errorf("respType:%v realType:%v",args0,reflect.TypeOf(resp)))
			}else{
				oV[0] = reflect.ValueOf(resp)
				oV[1] = reflect.New(args1).Elem()
			}
		} else {
			oV[0] =reflect.New(args0).Elem()
			oV[1] = reflect.ValueOf(err)
		}
		cbValue.Call(oV)
	}, timeout)
}

func (c *Client) CallSugar(sendMsg func(proto.Message,int64) error, msg proto.Message, resp proto.Message,timeout time.Duration) (error) {
	respCall, err := c.Client.Call(func(seqid int64) error {
		return sendMsg(msg,seqid)
	}, timeout)
	if err != nil {
		return err
	}

	msg, ok := respCall.(proto.Message)
	if !ok {
		return errors.New("resp not a proto.Message")
	}

	err = util.DeepCopyUseProtobuf(resp,msg)
	return  err
}
