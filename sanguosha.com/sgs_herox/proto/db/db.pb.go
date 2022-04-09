// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: db/db.proto

package db

import (
	def "sanguosha.com/sgs_herox/proto/def"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

//角色其它信息
type CharInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//玩家属性
	Icon     int32 `protobuf:"varint,1,opt,name=icon,proto3" json:"icon,omitempty"`         //头像
	IconEdge int32 `protobuf:"varint,2,opt,name=iconEdge,proto3" json:"iconEdge,omitempty"` //头像框
	//统计数据
	OnlineSecToday int32 `protobuf:"varint,20,opt,name=onlineSecToday,proto3" json:"onlineSecToday,omitempty"` //当天在线时间
	OnlineSecTotal int64 `protobuf:"varint,21,opt,name=onlineSecTotal,proto3" json:"onlineSecTotal,omitempty"` //总的在线时间
	LoginDayCnt    int32 `protobuf:"varint,22,opt,name=loginDayCnt,proto3" json:"loginDayCnt,omitempty"`       //总的登录天数
	C11LoginDayCnt int32 `protobuf:"varint,23,opt,name=c11LoginDayCnt,proto3" json:"c11LoginDayCnt,omitempty"` //连接登录的天数
	LastLoginTime  int64 `protobuf:"varint,24,opt,name=lastLoginTime,proto3" json:"lastLoginTime,omitempty"`   //上次登入时间
	LastLogoutTime int64 `protobuf:"varint,25,opt,name=lastLogoutTime,proto3" json:"lastLogoutTime,omitempty"` //上次登出时间
	LastResetTime  int64 `protobuf:"varint,30,opt,name=lastResetTime,proto3" json:"lastResetTime,omitempty"`   //上次重置的时间点
	Init           bool  `protobuf:"varint,51,opt,name=init,proto3" json:"init,omitempty"`
	//当日胜负数据统计
	CurDayWinCountText   int32 `protobuf:"varint,60,opt,name=curDayWinCountText,proto3" json:"curDayWinCountText,omitempty"`
	CurDayLoseCountText  int32 `protobuf:"varint,61,opt,name=curDayLoseCountText,proto3" json:"curDayLoseCountText,omitempty"`
	CurDayWinCountVoice  int32 `protobuf:"varint,62,opt,name=curDayWinCountVoice,proto3" json:"curDayWinCountVoice,omitempty"`
	CurDayLoseCountVoice int32 `protobuf:"varint,63,opt,name=curDayLoseCountVoice,proto3" json:"curDayLoseCountVoice,omitempty"`
	//总体胜负数据统计
	WinCountText   int32            `protobuf:"varint,100,opt,name=winCountText,proto3" json:"winCountText,omitempty"`
	LoseCountText  int32            `protobuf:"varint,101,opt,name=loseCountText,proto3" json:"loseCountText,omitempty"`
	WinCountVoice  int32            `protobuf:"varint,110,opt,name=winCountVoice,proto3" json:"winCountVoice,omitempty"`
	LoseCountVoice int32            `protobuf:"varint,111,opt,name=loseCountVoice,proto3" json:"loseCountVoice,omitempty"`
	Score          []*Score         `protobuf:"bytes,120,rep,name=score,proto3" json:"score,omitempty"`
	CardGroups     []*def.CardGroup `protobuf:"bytes,121,rep,name=cardGroups,proto3" json:"cardGroups,omitempty"`
	NowUseId       int32            `protobuf:"varint,122,opt,name=nowUseId,proto3" json:"nowUseId,omitempty"`
}

func (x *CharInfo) Reset() {
	*x = CharInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_db_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CharInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CharInfo) ProtoMessage() {}

func (x *CharInfo) ProtoReflect() protoreflect.Message {
	mi := &file_db_db_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CharInfo.ProtoReflect.Descriptor instead.
func (*CharInfo) Descriptor() ([]byte, []int) {
	return file_db_db_proto_rawDescGZIP(), []int{0}
}

func (x *CharInfo) GetIcon() int32 {
	if x != nil {
		return x.Icon
	}
	return 0
}

func (x *CharInfo) GetIconEdge() int32 {
	if x != nil {
		return x.IconEdge
	}
	return 0
}

func (x *CharInfo) GetOnlineSecToday() int32 {
	if x != nil {
		return x.OnlineSecToday
	}
	return 0
}

func (x *CharInfo) GetOnlineSecTotal() int64 {
	if x != nil {
		return x.OnlineSecTotal
	}
	return 0
}

func (x *CharInfo) GetLoginDayCnt() int32 {
	if x != nil {
		return x.LoginDayCnt
	}
	return 0
}

func (x *CharInfo) GetC11LoginDayCnt() int32 {
	if x != nil {
		return x.C11LoginDayCnt
	}
	return 0
}

func (x *CharInfo) GetLastLoginTime() int64 {
	if x != nil {
		return x.LastLoginTime
	}
	return 0
}

func (x *CharInfo) GetLastLogoutTime() int64 {
	if x != nil {
		return x.LastLogoutTime
	}
	return 0
}

func (x *CharInfo) GetLastResetTime() int64 {
	if x != nil {
		return x.LastResetTime
	}
	return 0
}

func (x *CharInfo) GetInit() bool {
	if x != nil {
		return x.Init
	}
	return false
}

func (x *CharInfo) GetCurDayWinCountText() int32 {
	if x != nil {
		return x.CurDayWinCountText
	}
	return 0
}

func (x *CharInfo) GetCurDayLoseCountText() int32 {
	if x != nil {
		return x.CurDayLoseCountText
	}
	return 0
}

func (x *CharInfo) GetCurDayWinCountVoice() int32 {
	if x != nil {
		return x.CurDayWinCountVoice
	}
	return 0
}

func (x *CharInfo) GetCurDayLoseCountVoice() int32 {
	if x != nil {
		return x.CurDayLoseCountVoice
	}
	return 0
}

func (x *CharInfo) GetWinCountText() int32 {
	if x != nil {
		return x.WinCountText
	}
	return 0
}

func (x *CharInfo) GetLoseCountText() int32 {
	if x != nil {
		return x.LoseCountText
	}
	return 0
}

func (x *CharInfo) GetWinCountVoice() int32 {
	if x != nil {
		return x.WinCountVoice
	}
	return 0
}

func (x *CharInfo) GetLoseCountVoice() int32 {
	if x != nil {
		return x.LoseCountVoice
	}
	return 0
}

func (x *CharInfo) GetScore() []*Score {
	if x != nil {
		return x.Score
	}
	return nil
}

func (x *CharInfo) GetCardGroups() []*def.CardGroup {
	if x != nil {
		return x.CardGroups
	}
	return nil
}

func (x *CharInfo) GetNowUseId() int32 {
	if x != nil {
		return x.NowUseId
	}
	return 0
}

//道具
type DBProp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Props    []*Int64KV `protobuf:"bytes,1,rep,name=props,proto3" json:"props,omitempty"`
	Consumes []*Int64KV `protobuf:"bytes,2,rep,name=consumes,proto3" json:"consumes,omitempty"` //道具消耗记录，目前支付货币和消耗品
}

func (x *DBProp) Reset() {
	*x = DBProp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_db_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBProp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBProp) ProtoMessage() {}

func (x *DBProp) ProtoReflect() protoreflect.Message {
	mi := &file_db_db_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBProp.ProtoReflect.Descriptor instead.
func (*DBProp) Descriptor() ([]byte, []int) {
	return file_db_db_proto_rawDescGZIP(), []int{1}
}

func (x *DBProp) GetProps() []*Int64KV {
	if x != nil {
		return x.Props
	}
	return nil
}

func (x *DBProp) GetConsumes() []*Int64KV {
	if x != nil {
		return x.Consumes
	}
	return nil
}

//隐藏分
type Score struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ScoreType int32      `protobuf:"varint,1,opt,name=ScoreType,proto3" json:"ScoreType,omitempty"`
	Score     []*Int32KV `protobuf:"bytes,2,rep,name=score,proto3" json:"score,omitempty"`
}

func (x *Score) Reset() {
	*x = Score{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_db_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Score) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Score) ProtoMessage() {}

func (x *Score) ProtoReflect() protoreflect.Message {
	mi := &file_db_db_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Score.ProtoReflect.Descriptor instead.
func (*Score) Descriptor() ([]byte, []int) {
	return file_db_db_proto_rawDescGZIP(), []int{2}
}

func (x *Score) GetScoreType() int32 {
	if x != nil {
		return x.ScoreType
	}
	return 0
}

func (x *Score) GetScore() []*Int32KV {
	if x != nil {
		return x.Score
	}
	return nil
}

type Int32KV struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key int32 `protobuf:"varint,1,opt,name=key,proto3" json:"key,omitempty"`
	V   int32 `protobuf:"varint,2,opt,name=v,proto3" json:"v,omitempty"`
}

func (x *Int32KV) Reset() {
	*x = Int32KV{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_db_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Int32KV) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Int32KV) ProtoMessage() {}

func (x *Int32KV) ProtoReflect() protoreflect.Message {
	mi := &file_db_db_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Int32KV.ProtoReflect.Descriptor instead.
func (*Int32KV) Descriptor() ([]byte, []int) {
	return file_db_db_proto_rawDescGZIP(), []int{3}
}

func (x *Int32KV) GetKey() int32 {
	if x != nil {
		return x.Key
	}
	return 0
}

func (x *Int32KV) GetV() int32 {
	if x != nil {
		return x.V
	}
	return 0
}

type Int64KV struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key int64 `protobuf:"varint,1,opt,name=key,proto3" json:"key,omitempty"`
	V   int64 `protobuf:"varint,2,opt,name=v,proto3" json:"v,omitempty"`
}

func (x *Int64KV) Reset() {
	*x = Int64KV{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_db_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Int64KV) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Int64KV) ProtoMessage() {}

func (x *Int64KV) ProtoReflect() protoreflect.Message {
	mi := &file_db_db_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Int64KV.ProtoReflect.Descriptor instead.
func (*Int64KV) Descriptor() ([]byte, []int) {
	return file_db_db_proto_rawDescGZIP(), []int{4}
}

func (x *Int64KV) GetKey() int64 {
	if x != nil {
		return x.Key
	}
	return 0
}

func (x *Int64KV) GetV() int64 {
	if x != nil {
		return x.V
	}
	return 0
}

var File_db_db_proto protoreflect.FileDescriptor

var file_db_db_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x64, 0x62, 0x2f, 0x64, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x64,
	0x62, 0x1a, 0x2b, 0x73, 0x61, 0x6e, 0x67, 0x75, 0x6f, 0x73, 0x68, 0x61, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x73, 0x67, 0x73, 0x5f, 0x68, 0x65, 0x72, 0x6f, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x64, 0x65, 0x66, 0x2f, 0x64, 0x65, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xad,
	0x06, 0x0a, 0x08, 0x43, 0x68, 0x61, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x69,
	0x63, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x12,
	0x1a, 0x0a, 0x08, 0x69, 0x63, 0x6f, 0x6e, 0x45, 0x64, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x08, 0x69, 0x63, 0x6f, 0x6e, 0x45, 0x64, 0x67, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x6f,
	0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x53, 0x65, 0x63, 0x54, 0x6f, 0x64, 0x61, 0x79, 0x18, 0x14, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0e, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x53, 0x65, 0x63, 0x54, 0x6f,
	0x64, 0x61, 0x79, 0x12, 0x26, 0x0a, 0x0e, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x53, 0x65, 0x63,
	0x54, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x15, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0e, 0x6f, 0x6e, 0x6c,
	0x69, 0x6e, 0x65, 0x53, 0x65, 0x63, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x20, 0x0a, 0x0b, 0x6c,
	0x6f, 0x67, 0x69, 0x6e, 0x44, 0x61, 0x79, 0x43, 0x6e, 0x74, 0x18, 0x16, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0b, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x44, 0x61, 0x79, 0x43, 0x6e, 0x74, 0x12, 0x26, 0x0a,
	0x0e, 0x63, 0x31, 0x31, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x44, 0x61, 0x79, 0x43, 0x6e, 0x74, 0x18,
	0x17, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x63, 0x31, 0x31, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x44,
	0x61, 0x79, 0x43, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67,
	0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x18, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x6c, 0x61,
	0x73, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x6c,
	0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x19, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x52, 0x65, 0x73, 0x65, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x6c, 0x61, 0x73, 0x74,
	0x52, 0x65, 0x73, 0x65, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x6e, 0x69,
	0x74, 0x18, 0x33, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x69, 0x6e, 0x69, 0x74, 0x12, 0x2e, 0x0a,
	0x12, 0x63, 0x75, 0x72, 0x44, 0x61, 0x79, 0x57, 0x69, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x54,
	0x65, 0x78, 0x74, 0x18, 0x3c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x12, 0x63, 0x75, 0x72, 0x44, 0x61,
	0x79, 0x57, 0x69, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x65, 0x78, 0x74, 0x12, 0x30, 0x0a,
	0x13, 0x63, 0x75, 0x72, 0x44, 0x61, 0x79, 0x4c, 0x6f, 0x73, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x54, 0x65, 0x78, 0x74, 0x18, 0x3d, 0x20, 0x01, 0x28, 0x05, 0x52, 0x13, 0x63, 0x75, 0x72, 0x44,
	0x61, 0x79, 0x4c, 0x6f, 0x73, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x65, 0x78, 0x74, 0x12,
	0x30, 0x0a, 0x13, 0x63, 0x75, 0x72, 0x44, 0x61, 0x79, 0x57, 0x69, 0x6e, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x56, 0x6f, 0x69, 0x63, 0x65, 0x18, 0x3e, 0x20, 0x01, 0x28, 0x05, 0x52, 0x13, 0x63, 0x75,
	0x72, 0x44, 0x61, 0x79, 0x57, 0x69, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x6f, 0x69, 0x63,
	0x65, 0x12, 0x32, 0x0a, 0x14, 0x63, 0x75, 0x72, 0x44, 0x61, 0x79, 0x4c, 0x6f, 0x73, 0x65, 0x43,
	0x6f, 0x75, 0x6e, 0x74, 0x56, 0x6f, 0x69, 0x63, 0x65, 0x18, 0x3f, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x14, 0x63, 0x75, 0x72, 0x44, 0x61, 0x79, 0x4c, 0x6f, 0x73, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x56, 0x6f, 0x69, 0x63, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x77, 0x69, 0x6e, 0x43, 0x6f, 0x75, 0x6e,
	0x74, 0x54, 0x65, 0x78, 0x74, 0x18, 0x64, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x77, 0x69, 0x6e,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x65, 0x78, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x6c, 0x6f, 0x73,
	0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x65, 0x78, 0x74, 0x18, 0x65, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0d, 0x6c, 0x6f, 0x73, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x54, 0x65, 0x78, 0x74, 0x12,
	0x24, 0x0a, 0x0d, 0x77, 0x69, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x6f, 0x69, 0x63, 0x65,
	0x18, 0x6e, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x77, 0x69, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x56, 0x6f, 0x69, 0x63, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x6c, 0x6f, 0x73, 0x65, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x56, 0x6f, 0x69, 0x63, 0x65, 0x18, 0x6f, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x6c,
	0x6f, 0x73, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x56, 0x6f, 0x69, 0x63, 0x65, 0x12, 0x1f, 0x0a,
	0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x78, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x64,
	0x62, 0x2e, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x32,
	0x0a, 0x0a, 0x63, 0x61, 0x72, 0x64, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18, 0x79, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x64, 0x65, 0x66, 0x2e, 0x43, 0x61, 0x72,
	0x64, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x0a, 0x63, 0x61, 0x72, 0x64, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x6e, 0x6f, 0x77, 0x55, 0x73, 0x65, 0x49, 0x64, 0x18, 0x7a,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6e, 0x6f, 0x77, 0x55, 0x73, 0x65, 0x49, 0x64, 0x22, 0x54,
	0x0a, 0x06, 0x44, 0x42, 0x50, 0x72, 0x6f, 0x70, 0x12, 0x21, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x70,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x64, 0x62, 0x2e, 0x49, 0x6e, 0x74,
	0x36, 0x34, 0x4b, 0x56, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x70, 0x73, 0x12, 0x27, 0x0a, 0x08, 0x63,
	0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e,
	0x64, 0x62, 0x2e, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x4b, 0x56, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x73,
	0x75, 0x6d, 0x65, 0x73, 0x22, 0x48, 0x0a, 0x05, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x09, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x05, 0x73,
	0x63, 0x6f, 0x72, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x64, 0x62, 0x2e,
	0x49, 0x6e, 0x74, 0x33, 0x32, 0x4b, 0x56, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x22, 0x29,
	0x0a, 0x07, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x4b, 0x56, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x0c, 0x0a, 0x01, 0x76,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x01, 0x76, 0x22, 0x29, 0x0a, 0x07, 0x49, 0x6e, 0x74,
	0x36, 0x34, 0x4b, 0x56, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x0c, 0x0a, 0x01, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x01, 0x76, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x64, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_db_db_proto_rawDescOnce sync.Once
	file_db_db_proto_rawDescData = file_db_db_proto_rawDesc
)

func file_db_db_proto_rawDescGZIP() []byte {
	file_db_db_proto_rawDescOnce.Do(func() {
		file_db_db_proto_rawDescData = protoimpl.X.CompressGZIP(file_db_db_proto_rawDescData)
	})
	return file_db_db_proto_rawDescData
}

var file_db_db_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_db_db_proto_goTypes = []interface{}{
	(*CharInfo)(nil),      // 0: db.CharInfo
	(*DBProp)(nil),        // 1: db.DBProp
	(*Score)(nil),         // 2: db.Score
	(*Int32KV)(nil),       // 3: db.Int32KV
	(*Int64KV)(nil),       // 4: db.Int64KV
	(*def.CardGroup)(nil), // 5: gamedef.CardGroup
}
var file_db_db_proto_depIdxs = []int32{
	2, // 0: db.CharInfo.score:type_name -> db.Score
	5, // 1: db.CharInfo.cardGroups:type_name -> gamedef.CardGroup
	4, // 2: db.DBProp.props:type_name -> db.Int64KV
	4, // 3: db.DBProp.consumes:type_name -> db.Int64KV
	3, // 4: db.Score.score:type_name -> db.Int32KV
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_db_db_proto_init() }
func file_db_db_proto_init() {
	if File_db_db_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_db_db_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CharInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_db_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBProp); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_db_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Score); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_db_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Int32KV); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_db_db_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Int64KV); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_db_db_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_db_db_proto_goTypes,
		DependencyIndexes: file_db_db_proto_depIdxs,
		MessageInfos:      file_db_db_proto_msgTypes,
	}.Build()
	File_db_db_proto = out.File
	file_db_db_proto_rawDesc = nil
	file_db_db_proto_goTypes = nil
	file_db_db_proto_depIdxs = nil
}
