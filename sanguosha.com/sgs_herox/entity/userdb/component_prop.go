package userdb

import (
	"github.com/golang/protobuf/proto"
	"math"
	"sanguosha.com/sgs_herox/proto/db"
)

const (
	LimitItemNum int32 = math.MaxInt32
)

type prop struct {
	BaseComp

	user  *User
	props map[int64]int64 //物品
}

func (p *prop) dbFieldName() string {
	return "prop"
}

func (p *prop) init(user *User, data []byte) error {
	p.user = user

	dbProp := db.DBProp{}
	err := proto.Unmarshal(data, &dbProp)
	if err != nil {
		return err
	}

	p.props = make(map[int64]int64, len(dbProp.Props))
	for _, v := range dbProp.Props {
		p.props[v.Key] = v.V
	}

	return nil
}

func (p *prop) toProtoMessage() proto.Message {
	dbProp := db.DBProp{}
	for k, v := range p.props {
		dbProp.Props = append(dbProp.Props, &db.Int64KV{
			Key: k,
			V:   v,
		})
	}
	return &dbProp
}

////返回最新的值
//func (p *prop) add(propID int32, add int32) (int32, int32) { //now realGet
//	v, _ := p.props[propID]
//
//	if add <= 0 {
//		return v, 0
//	}
//
//	if v >= LimitItemNum {
//		return v, 0
//	}
//
//	maxAdd := LimitItemNum - v
//	if add > maxAdd {
//		add = maxAdd
//	}
//
//	now := v + add
//	p.props[propID] = now
//	p.setDirty()
//	return now, add
//}
//
////返回错误，最新的值
//func (p *prop) sub(propID int32, sub int32) (error, int32) {
//	v, exist := p.props[propID]
//	if !exist {
//		return ErrHaveNoThisProp, 0
//	}
//
//	if v < sub {
//		delete(p.props, propID)
//		p.setDirty()
//		return ErrPropNumSubNotEnough, 0
//	}
//
//	new := v - sub
//	p.props[propID] = new
//
//	if new == 0 {
//		delete(p.props, propID)
//	}
//	p.setDirty()
//	return nil, new
//}
//
//func (p *prop) count(propID int32) int32 {
//	v, _ := p.props[propID]
//	return v
//}
//
//func (p *prop) All() map[int32]int32 {
//	m := make(map[int32]int32, len(p.props))
//	for k, v := range p.props {
//		m[k] = v
//	}
//	return m
//}
