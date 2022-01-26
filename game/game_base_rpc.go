package game

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"reflect"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/gameshared/entityservice"
	"time"
)

//TODO 每局游戏一个协程，这样可以使用RequestCall阻塞式api开发，不需要下面这种复杂的方式控制流程
//cb为func(resp *smsg.XXXX,err error){}
func (g *GameBase) ReqEntity(msg entityservice.MessageWithUserid, cb interface{}, duration time.Duration) {
	timer := g.WaitingParallel(duration+200*time.Millisecond, func() {
		reqSugarHandle(msg.String(), nil, errors.New("ErrTimeOut"), cb)
	})

	EntityService.Request(msg, func(resp proto.Message, err error) {
		g.DoNow(func() {
			if ok := timer.Stop(); ok {
				reqSugarHandle(msg.String(), resp, err, cb)
			}
		})
	}, duration)
}

//cb为func(resp *smsg.XXXX,err error){}
func (g *GameBase) ReqServerType(typ appframe.ServerType, msg proto.Message, cb interface{}, duration time.Duration) {
	timer := g.WaitingParallel(duration+200*time.Millisecond, func() {
		reqSugarHandle(msg.String(), nil, errors.New("ErrTimeOut"), cb)
	})

	App.GetService(typ).Request(msg, func(resp proto.Message, err error) {
		g.DoNow(func() {
			if ok := timer.Stop(); ok {
				reqSugarHandle(msg.String(), resp, err, cb)
			}
		})
	}, duration)
}

func reqSugarHandle(reqMsgName string, resp proto.Message, err error, cb interface{}) {
	log := logrus.WithField("msg", reqMsgName)
	cbType := reflect.TypeOf(cb)
	if cbType.Kind() != reflect.Func {
		log.Error("ReqSugar cb not a func")
		return
	}
	cbValue := reflect.ValueOf(cb)
	numArgs := cbType.NumIn()
	if numArgs != 2 {
		log.Error("ReqSugar cb param num args !=2")
		return
	}
	args0 := cbType.In(0)
	if args0.Kind() != reflect.Ptr {
		log.Error("ReqSugar cb param args0 not ptr")
		return
	}

	//TODO 严格检查参数类型
	args1 := cbType.In(1)
	oV := make([]reflect.Value, 2)

	if err == nil {
		if reflect.TypeOf(resp) != args0 {
			oV[0] = reflect.New(args0).Elem()
			oV[1] = reflect.ValueOf(fmt.Errorf("respType:%v realType:%v", args0, reflect.TypeOf(resp)))
		} else {
			oV[0] = reflect.ValueOf(resp)
			oV[1] = reflect.New(args1).Elem()
		}
	} else {
		oV[0] = reflect.New(args0).Elem()
		oV[1] = reflect.ValueOf(err)
	}
	cbValue.Call(oV)
}
