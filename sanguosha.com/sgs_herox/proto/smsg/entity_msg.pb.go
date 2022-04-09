// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: smsg/entity_msg.proto

package smsg

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	def "sanguosha.com/sgs_herox/proto/def"
	_ "sanguosha.com/sgs_herox/proto/gameconf"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AllEnReqUserInfo_UserDataTyp int32

const (
	AllEnReqUserInfo_UDTInvalid AllEnReqUserInfo_UserDataTyp = 0
	AllEnReqUserInfo_UDTUser    AllEnReqUserInfo_UserDataTyp = 1 //角色
	AllEnReqUserInfo_UDTProp    AllEnReqUserInfo_UserDataTyp = 2 //道具列表
)

// Enum value maps for AllEnReqUserInfo_UserDataTyp.
var (
	AllEnReqUserInfo_UserDataTyp_name = map[int32]string{
		0: "UDTInvalid",
		1: "UDTUser",
		2: "UDTProp",
	}
	AllEnReqUserInfo_UserDataTyp_value = map[string]int32{
		"UDTInvalid": 0,
		"UDTUser":    1,
		"UDTProp":    2,
	}
)

func (x AllEnReqUserInfo_UserDataTyp) Enum() *AllEnReqUserInfo_UserDataTyp {
	p := new(AllEnReqUserInfo_UserDataTyp)
	*p = x
	return p
}

func (x AllEnReqUserInfo_UserDataTyp) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AllEnReqUserInfo_UserDataTyp) Descriptor() protoreflect.EnumDescriptor {
	return file_smsg_entity_msg_proto_enumTypes[0].Descriptor()
}

func (AllEnReqUserInfo_UserDataTyp) Type() protoreflect.EnumType {
	return &file_smsg_entity_msg_proto_enumTypes[0]
}

func (x AllEnReqUserInfo_UserDataTyp) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AllEnReqUserInfo_UserDataTyp.Descriptor instead.
func (AllEnReqUserInfo_UserDataTyp) EnumDescriptor() ([]byte, []int) {
	return file_smsg_entity_msg_proto_rawDescGZIP(), []int{0, 0}
}

type AllEnRespUserInfo_ErrCode int32

const (
	AllEnRespUserInfo_Invalid      AllEnRespUserInfo_ErrCode = 0
	AllEnRespUserInfo_GetUserError AllEnRespUserInfo_ErrCode = 1
)

// Enum value maps for AllEnRespUserInfo_ErrCode.
var (
	AllEnRespUserInfo_ErrCode_name = map[int32]string{
		0: "Invalid",
		1: "GetUserError",
	}
	AllEnRespUserInfo_ErrCode_value = map[string]int32{
		"Invalid":      0,
		"GetUserError": 1,
	}
)

func (x AllEnRespUserInfo_ErrCode) Enum() *AllEnRespUserInfo_ErrCode {
	p := new(AllEnRespUserInfo_ErrCode)
	*p = x
	return p
}

func (x AllEnRespUserInfo_ErrCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AllEnRespUserInfo_ErrCode) Descriptor() protoreflect.EnumDescriptor {
	return file_smsg_entity_msg_proto_enumTypes[1].Descriptor()
}

func (AllEnRespUserInfo_ErrCode) Type() protoreflect.EnumType {
	return &file_smsg_entity_msg_proto_enumTypes[1]
}

func (x AllEnRespUserInfo_ErrCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AllEnRespUserInfo_ErrCode.Descriptor instead.
func (AllEnRespUserInfo_ErrCode) EnumDescriptor() ([]byte, []int) {
	return file_smsg_entity_msg_proto_rawDescGZIP(), []int{1, 0}
}

//请求玩家信息
type AllEnReqUserInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Userid    uint64                         `protobuf:"varint,2,opt,name=userid,proto3" json:"userid,omitempty"`
	DataTypes []AllEnReqUserInfo_UserDataTyp `protobuf:"varint,3,rep,packed,name=dataTypes,proto3,enum=smsg.AllEnReqUserInfo_UserDataTyp" json:"dataTypes,omitempty"`
}

func (x *AllEnReqUserInfo) Reset() {
	*x = AllEnReqUserInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smsg_entity_msg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllEnReqUserInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllEnReqUserInfo) ProtoMessage() {}

func (x *AllEnReqUserInfo) ProtoReflect() protoreflect.Message {
	mi := &file_smsg_entity_msg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllEnReqUserInfo.ProtoReflect.Descriptor instead.
func (*AllEnReqUserInfo) Descriptor() ([]byte, []int) {
	return file_smsg_entity_msg_proto_rawDescGZIP(), []int{0}
}

func (x *AllEnReqUserInfo) GetUserid() uint64 {
	if x != nil {
		return x.Userid
	}
	return 0
}

func (x *AllEnReqUserInfo) GetDataTypes() []AllEnReqUserInfo_UserDataTyp {
	if x != nil {
		return x.DataTypes
	}
	return nil
}

type AllEnRespUserInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Seqid   int64                     `protobuf:"varint,1,opt,name=seqid,proto3" json:"seqid,omitempty"`
	Userid  uint64                    `protobuf:"varint,2,opt,name=userid,proto3" json:"userid,omitempty"`
	ErrCode AllEnRespUserInfo_ErrCode `protobuf:"varint,3,opt,name=errCode,proto3,enum=smsg.AllEnRespUserInfo_ErrCode" json:"errCode,omitempty"`
	Data    *UserDataContent          `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *AllEnRespUserInfo) Reset() {
	*x = AllEnRespUserInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smsg_entity_msg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllEnRespUserInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllEnRespUserInfo) ProtoMessage() {}

func (x *AllEnRespUserInfo) ProtoReflect() protoreflect.Message {
	mi := &file_smsg_entity_msg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllEnRespUserInfo.ProtoReflect.Descriptor instead.
func (*AllEnRespUserInfo) Descriptor() ([]byte, []int) {
	return file_smsg_entity_msg_proto_rawDescGZIP(), []int{1}
}

func (x *AllEnRespUserInfo) GetSeqid() int64 {
	if x != nil {
		return x.Seqid
	}
	return 0
}

func (x *AllEnRespUserInfo) GetUserid() uint64 {
	if x != nil {
		return x.Userid
	}
	return 0
}

func (x *AllEnRespUserInfo) GetErrCode() AllEnRespUserInfo_ErrCode {
	if x != nil {
		return x.ErrCode
	}
	return AllEnRespUserInfo_Invalid
}

func (x *AllEnRespUserInfo) GetData() *UserDataContent {
	if x != nil {
		return x.Data
	}
	return nil
}

type UserDataContent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *UserDataContent) Reset() {
	*x = UserDataContent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smsg_entity_msg_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserDataContent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserDataContent) ProtoMessage() {}

func (x *UserDataContent) ProtoReflect() protoreflect.Message {
	mi := &file_smsg_entity_msg_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserDataContent.ProtoReflect.Descriptor instead.
func (*UserDataContent) Descriptor() ([]byte, []int) {
	return file_smsg_entity_msg_proto_rawDescGZIP(), []int{2}
}

func (x *UserDataContent) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserBrief           *def.UserBrief `protobuf:"bytes,1,opt,name=userBrief,proto3" json:"userBrief,omitempty"`
	IsOnline            bool           `protobuf:"varint,11,opt,name=isOnline,proto3" json:"isOnline,omitempty"`                       //是否在线
	Status              int32          `protobuf:"varint,12,opt,name=status,proto3" json:"status,omitempty"`                           //状态
	LoginIP             string         `protobuf:"bytes,15,opt,name=loginIP,proto3" json:"loginIP,omitempty"`                          //登录的ip
	ThirdAccountId      string         `protobuf:"bytes,16,opt,name=thirdAccountId,proto3" json:"thirdAccountId,omitempty"`            //登录渠道的账号Id
	CreateTime          int64          `protobuf:"varint,17,opt,name=createTime,proto3" json:"createTime,omitempty"`                   //账号创建日期
	RegisterTime        int64          `protobuf:"varint,20,opt,name=registerTime,proto3" json:"registerTime,omitempty"`               //注册时间
	RegisterIP          string         `protobuf:"bytes,21,opt,name=registerIP,proto3" json:"registerIP,omitempty"`                    //注册IP
	LoginTime           int64          `protobuf:"varint,22,opt,name=loginTime,proto3" json:"loginTime,omitempty"`                     //上次登录时间
	LogoutTime          int64          `protobuf:"varint,24,opt,name=logoutTime,proto3" json:"logoutTime,omitempty"`                   //上次登出时间
	BornTime            int64          `protobuf:"varint,25,opt,name=bornTime,proto3" json:"bornTime,omitempty"`                       //出生时间
	ContinueLoginDayCnt int32          `protobuf:"varint,33,opt,name=continueLoginDayCnt,proto3" json:"continueLoginDayCnt,omitempty"` //连续登录天数
	LoginDayCnt         int32          `protobuf:"varint,34,opt,name=loginDayCnt,proto3" json:"loginDayCnt,omitempty"`                 //总登录天数
	OnlineSec           int64          `protobuf:"varint,36,opt,name=onlineSec,proto3" json:"onlineSec,omitempty"`                     //在线时间
	OnlineSecToday      int32          `protobuf:"varint,37,opt,name=onlineSecToday,proto3" json:"onlineSecToday,omitempty"`           //今日在线时间
	CardGroup           *def.CardGroup `protobuf:"bytes,38,opt,name=cardGroup,proto3" json:"cardGroup,omitempty"`                      //玩家当前所使用的卡组
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_smsg_entity_msg_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_smsg_entity_msg_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_smsg_entity_msg_proto_rawDescGZIP(), []int{3}
}

func (x *User) GetUserBrief() *def.UserBrief {
	if x != nil {
		return x.UserBrief
	}
	return nil
}

func (x *User) GetIsOnline() bool {
	if x != nil {
		return x.IsOnline
	}
	return false
}

func (x *User) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *User) GetLoginIP() string {
	if x != nil {
		return x.LoginIP
	}
	return ""
}

func (x *User) GetThirdAccountId() string {
	if x != nil {
		return x.ThirdAccountId
	}
	return ""
}

func (x *User) GetCreateTime() int64 {
	if x != nil {
		return x.CreateTime
	}
	return 0
}

func (x *User) GetRegisterTime() int64 {
	if x != nil {
		return x.RegisterTime
	}
	return 0
}

func (x *User) GetRegisterIP() string {
	if x != nil {
		return x.RegisterIP
	}
	return ""
}

func (x *User) GetLoginTime() int64 {
	if x != nil {
		return x.LoginTime
	}
	return 0
}

func (x *User) GetLogoutTime() int64 {
	if x != nil {
		return x.LogoutTime
	}
	return 0
}

func (x *User) GetBornTime() int64 {
	if x != nil {
		return x.BornTime
	}
	return 0
}

func (x *User) GetContinueLoginDayCnt() int32 {
	if x != nil {
		return x.ContinueLoginDayCnt
	}
	return 0
}

func (x *User) GetLoginDayCnt() int32 {
	if x != nil {
		return x.LoginDayCnt
	}
	return 0
}

func (x *User) GetOnlineSec() int64 {
	if x != nil {
		return x.OnlineSec
	}
	return 0
}

func (x *User) GetOnlineSecToday() int32 {
	if x != nil {
		return x.OnlineSecToday
	}
	return 0
}

func (x *User) GetCardGroup() *def.CardGroup {
	if x != nil {
		return x.CardGroup
	}
	return nil
}

var File_smsg_entity_msg_proto protoreflect.FileDescriptor

var file_smsg_entity_msg_proto_rawDesc = []byte{
	0x0a, 0x15, 0x73, 0x6d, 0x73, 0x67, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x6d, 0x73,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x73, 0x6d, 0x73, 0x67, 0x1a, 0x36, 0x73,
	0x61, 0x6e, 0x67, 0x75, 0x6f, 0x73, 0x68, 0x61, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x67, 0x73,
	0x5f, 0x68, 0x65, 0x72, 0x6f, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x61, 0x6d,
	0x65, 0x63, 0x6f, 0x6e, 0x66, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2b, 0x73, 0x61, 0x6e, 0x67, 0x75, 0x6f, 0x73, 0x68, 0x61,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x67, 0x73, 0x5f, 0x68, 0x65, 0x72, 0x6f, 0x78, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x64, 0x65, 0x66, 0x2f, 0x64, 0x65, 0x66, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xa5, 0x01, 0x0a, 0x10, 0x41, 0x6c, 0x6c, 0x45, 0x6e, 0x52, 0x65, 0x71, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x69, 0x64, 0x12,
	0x40, 0x0a, 0x09, 0x64, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0e, 0x32, 0x22, 0x2e, 0x73, 0x6d, 0x73, 0x67, 0x2e, 0x41, 0x6c, 0x6c, 0x45, 0x6e, 0x52,
	0x65, 0x71, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x44,
	0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x52, 0x09, 0x64, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65,
	0x73, 0x22, 0x37, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70,
	0x12, 0x0e, 0x0a, 0x0a, 0x55, 0x44, 0x54, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x10, 0x00,
	0x12, 0x0b, 0x0a, 0x07, 0x55, 0x44, 0x54, 0x55, 0x73, 0x65, 0x72, 0x10, 0x01, 0x12, 0x0b, 0x0a,
	0x07, 0x55, 0x44, 0x54, 0x50, 0x72, 0x6f, 0x70, 0x10, 0x02, 0x22, 0xd1, 0x01, 0x0a, 0x11, 0x41,
	0x6c, 0x6c, 0x45, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x14, 0x0a, 0x05, 0x73, 0x65, 0x71, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x73, 0x65, 0x71, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x69, 0x64, 0x12, 0x39,
	0x0a, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1f, 0x2e, 0x73, 0x6d, 0x73, 0x67, 0x2e, 0x41, 0x6c, 0x6c, 0x45, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65,
	0x52, 0x07, 0x65, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x29, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x73, 0x6d, 0x73, 0x67, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x28, 0x0a, 0x07, 0x45, 0x72, 0x72, 0x43, 0x6f, 0x64, 0x65, 0x12,
	0x0b, 0x0a, 0x07, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c,
	0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x10, 0x01, 0x22, 0x31,
	0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x12, 0x1e, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0a, 0x2e, 0x73, 0x6d, 0x73, 0x67, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65,
	0x72, 0x22, 0xb8, 0x04, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x30, 0x0a, 0x09, 0x75, 0x73,
	0x65, 0x72, 0x42, 0x72, 0x69, 0x65, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x67, 0x61, 0x6d, 0x65, 0x64, 0x65, 0x66, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x42, 0x72, 0x69, 0x65,
	0x66, 0x52, 0x09, 0x75, 0x73, 0x65, 0x72, 0x42, 0x72, 0x69, 0x65, 0x66, 0x12, 0x1a, 0x0a, 0x08,
	0x69, 0x73, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08,
	0x69, 0x73, 0x4f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x18, 0x0a, 0x07, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x49, 0x50, 0x18, 0x0f, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x49, 0x50, 0x12, 0x26, 0x0a, 0x0e, 0x74, 0x68,
	0x69, 0x72, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x10, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x74, 0x68, 0x69, 0x72, 0x64, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65,
	0x18, 0x11, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x54, 0x69,
	0x6d, 0x65, 0x18, 0x14, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x49, 0x50, 0x18, 0x15, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x49, 0x50, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x54,
	0x69, 0x6d, 0x65, 0x18, 0x16, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x6c, 0x6f, 0x67, 0x69, 0x6e,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x6c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x54, 0x69,
	0x6d, 0x65, 0x18, 0x18, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x6c, 0x6f, 0x67, 0x6f, 0x75, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x6f, 0x72, 0x6e, 0x54, 0x69, 0x6d, 0x65,
	0x18, 0x19, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x62, 0x6f, 0x72, 0x6e, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x30, 0x0a, 0x13, 0x63, 0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x44, 0x61, 0x79, 0x43, 0x6e, 0x74, 0x18, 0x21, 0x20, 0x01, 0x28, 0x05, 0x52, 0x13, 0x63,
	0x6f, 0x6e, 0x74, 0x69, 0x6e, 0x75, 0x65, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x44, 0x61, 0x79, 0x43,
	0x6e, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x44, 0x61, 0x79, 0x43, 0x6e,
	0x74, 0x18, 0x22, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x44, 0x61,
	0x79, 0x43, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x53, 0x65,
	0x63, 0x18, 0x24, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x53,
	0x65, 0x63, 0x12, 0x26, 0x0a, 0x0e, 0x6f, 0x6e, 0x6c, 0x69, 0x6e, 0x65, 0x53, 0x65, 0x63, 0x54,
	0x6f, 0x64, 0x61, 0x79, 0x18, 0x25, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x6f, 0x6e, 0x6c, 0x69,
	0x6e, 0x65, 0x53, 0x65, 0x63, 0x54, 0x6f, 0x64, 0x61, 0x79, 0x12, 0x30, 0x0a, 0x09, 0x63, 0x61,
	0x72, 0x64, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x26, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x67, 0x61, 0x6d, 0x65, 0x64, 0x65, 0x66, 0x2e, 0x43, 0x61, 0x72, 0x64, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x52, 0x09, 0x63, 0x61, 0x72, 0x64, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x42, 0x08, 0x5a, 0x06,
	0x2e, 0x2f, 0x73, 0x6d, 0x73, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_smsg_entity_msg_proto_rawDescOnce sync.Once
	file_smsg_entity_msg_proto_rawDescData = file_smsg_entity_msg_proto_rawDesc
)

func file_smsg_entity_msg_proto_rawDescGZIP() []byte {
	file_smsg_entity_msg_proto_rawDescOnce.Do(func() {
		file_smsg_entity_msg_proto_rawDescData = protoimpl.X.CompressGZIP(file_smsg_entity_msg_proto_rawDescData)
	})
	return file_smsg_entity_msg_proto_rawDescData
}

var file_smsg_entity_msg_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_smsg_entity_msg_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_smsg_entity_msg_proto_goTypes = []interface{}{
	(AllEnReqUserInfo_UserDataTyp)(0), // 0: smsg.AllEnReqUserInfo.UserDataTyp
	(AllEnRespUserInfo_ErrCode)(0),    // 1: smsg.AllEnRespUserInfo.ErrCode
	(*AllEnReqUserInfo)(nil),          // 2: smsg.AllEnReqUserInfo
	(*AllEnRespUserInfo)(nil),         // 3: smsg.AllEnRespUserInfo
	(*UserDataContent)(nil),           // 4: smsg.UserDataContent
	(*User)(nil),                      // 5: smsg.User
	(*def.UserBrief)(nil),             // 6: gamedef.UserBrief
	(*def.CardGroup)(nil),             // 7: gamedef.CardGroup
}
var file_smsg_entity_msg_proto_depIdxs = []int32{
	0, // 0: smsg.AllEnReqUserInfo.dataTypes:type_name -> smsg.AllEnReqUserInfo.UserDataTyp
	1, // 1: smsg.AllEnRespUserInfo.errCode:type_name -> smsg.AllEnRespUserInfo.ErrCode
	4, // 2: smsg.AllEnRespUserInfo.data:type_name -> smsg.UserDataContent
	5, // 3: smsg.UserDataContent.user:type_name -> smsg.User
	6, // 4: smsg.User.userBrief:type_name -> gamedef.UserBrief
	7, // 5: smsg.User.cardGroup:type_name -> gamedef.CardGroup
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_smsg_entity_msg_proto_init() }
func file_smsg_entity_msg_proto_init() {
	if File_smsg_entity_msg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_smsg_entity_msg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllEnReqUserInfo); i {
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
		file_smsg_entity_msg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllEnRespUserInfo); i {
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
		file_smsg_entity_msg_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserDataContent); i {
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
		file_smsg_entity_msg_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
			RawDescriptor: file_smsg_entity_msg_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_smsg_entity_msg_proto_goTypes,
		DependencyIndexes: file_smsg_entity_msg_proto_depIdxs,
		EnumInfos:         file_smsg_entity_msg_proto_enumTypes,
		MessageInfos:      file_smsg_entity_msg_proto_msgTypes,
	}.Build()
	File_smsg_entity_msg_proto = out.File
	file_smsg_entity_msg_proto_rawDesc = nil
	file_smsg_entity_msg_proto_goTypes = nil
	file_smsg_entity_msg_proto_depIdxs = nil
}
