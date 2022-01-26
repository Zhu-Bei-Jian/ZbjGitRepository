package client

import (
	"github.com/golang/protobuf/proto"
	"reflect"
	"sanguosha.com/sgs_herox/proto/cmsg"
	"time"
)

type UserData struct {
	myData *cmsg.SRespMyData
	//friendData *cmsg.SNoticeFriendData

	//queryUser []*cmsg.HeadInfo
	userID uint64

	waitting map[reflect.Type]int64

	gameMatch bool
	gameUuid  string
	seatid    int32
}

func (p *UserData) SetWaitMsg(msg proto.Message) {
	if p.waitting == nil {
		p.waitting = make(map[reflect.Type]int64)
	}
	p.waitting[reflect.TypeOf(msg)] = time.Now().UnixNano()
}

func (p *UserData) IsWaitMsg(msg proto.Message) bool {
	if p.waitting == nil {
		return false
	}
	_, ok := p.waitting[reflect.TypeOf(msg)]
	return ok
}

func (p *UserData) CheckWaitMsg(msg proto.Message) {
	if msg == nil || p.waitting == nil {
		return
	}
	refType := reflect.TypeOf(msg)
	v, ok := p.waitting[refType]
	if !ok {
		return
	}
	delete(p.waitting, refType)
	diff := (time.Now().UnixNano() - v)
	if diff >= 0 {
		pushCol("resp", refType, diff)
	}
}
