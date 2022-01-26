package protoreq

//import (
//	"reflect"
//
//	"github.com/golang/protobuf/proto"
//)

//// IsValidMsg 判断是否为有效的请求响应类型消息 (消息必须有字段 Seqid int64)
//func IsValidMsg(msg proto.Message) bool {
//	t := reflect.TypeOf(msg)
//	if t.Kind() == reflect.Ptr {
//		t = t.Elem()
//	}
//	if t.Kind() == reflect.Struct {
//		f, ok := t.FieldByName("Seqid")
//		if ok && f.Type.Kind() == reflect.Int64 {
//			return true
//		}
//	}
//
//	return false
//}
//
//// GetSeqid 获取消息的 Seqid
//func GetSeqid(msg proto.Message) (int64, bool) {
//	if i, ok := msg.(interface {
//		GetSeqid() int64
//	}); ok {
//		return i.GetSeqid(), true
//	}
//
//	v := reflect.ValueOf(msg)
//	if v.Kind() == reflect.Ptr {
//		v = v.Elem()
//	}
//	if v.Kind() == reflect.Struct {
//		v := v.FieldByName("Seqid")
//		if v.Kind() == reflect.Int64 {
//			return v.Int(), true
//		}
//	}
//
//	return 0, false
//}

//// SetSeqid 设置响应消息的 Seqid
//func SetSeqid(msg proto.Message, seqid int64) {
//	v := reflect.ValueOf(msg)
//	if v.Kind() == reflect.Ptr {
//		v = v.Elem()
//	}
//	v.FieldByName("Seqid").SetInt(seqid)
//}
