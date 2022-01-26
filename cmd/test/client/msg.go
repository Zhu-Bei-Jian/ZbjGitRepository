package client

import (
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"

	"sanguosha.com/baselib/util"
)

var regMsg, msg2id, type2name, id2msg = func() (
	func(proto.Message),
	func(proto.Message) (uint32, bool),
	func(reflect.Type) string,
	func(uint32) (proto.Message, bool),
) {

	var (
		mtx       sync.RWMutex
		id2type   = map[uint32]reflect.Type{}
		type2id   = map[reflect.Type]uint32{}
		type2name = map[reflect.Type]string{}
	)

	regMsg := func(msg proto.Message) {
		typ := reflect.TypeOf(msg)
		if typ == nil {
			panic("invalid msg type")
		}

		// quick path
		mtx.RLock()
		if _, exist := type2id[typ]; exist {
			mtx.RUnlock()
			return
		}
		mtx.RUnlock()

		mtx.Lock()
		defer mtx.Unlock()

		// double check
		if _, exist := type2id[typ]; exist {
			return
		}

		msgName := proto.MessageName(msg)
		id := util.StringHash(proto.MessageName(msg))

		id2type[id] = typ
		type2id[typ] = id
		type2name[typ] = msgName
	}

	msg2id := func(msg proto.Message) (uint32, bool) {
		typ := reflect.TypeOf(msg)
		mtx.RLock()
		id, ok := type2id[typ]
		mtx.RUnlock()
		return id, ok
	}

	_type2name := func(typ reflect.Type) string {
		mtx.RLock()
		name, _ := type2name[typ]
		mtx.RUnlock()
		return name
	}

	id2msg := func(id uint32) (proto.Message, bool) {
		mtx.RLock()
		typ, ok := id2type[id]
		mtx.RUnlock()
		if !ok {
			return nil, false
		}
		return reflect.New(typ.Elem()).Interface().(proto.Message), true
	}

	return regMsg, msg2id, _type2name, id2msg
}()
