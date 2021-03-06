// Generated by github.com/davyxu/tabtoy
// Version: 2.9.1
// Modify by nuyan
// DO NOT EDIT!!

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: gameconf/game_stable_config.proto

package gameconf

import (
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

// Defined in table: GameStableConfig
type GameStableConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DoNotUseThis []*DoNotUseThisDefine `protobuf:"bytes,1,rep,name=DoNotUseThis,proto3" json:"DoNotUseThis,omitempty"` // DoNotUseThis
	UsernameConf []*UsernameConfDefine `protobuf:"bytes,2,rep,name=UsernameConf,proto3" json:"UsernameConf,omitempty"` // UsernameConf
}

func (x *GameStableConfig) Reset() {
	*x = GameStableConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gameconf_game_stable_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GameStableConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GameStableConfig) ProtoMessage() {}

func (x *GameStableConfig) ProtoReflect() protoreflect.Message {
	mi := &file_gameconf_game_stable_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GameStableConfig.ProtoReflect.Descriptor instead.
func (*GameStableConfig) Descriptor() ([]byte, []int) {
	return file_gameconf_game_stable_config_proto_rawDescGZIP(), []int{0}
}

func (x *GameStableConfig) GetDoNotUseThis() []*DoNotUseThisDefine {
	if x != nil {
		return x.DoNotUseThis
	}
	return nil
}

func (x *GameStableConfig) GetUsernameConf() []*UsernameConfDefine {
	if x != nil {
		return x.UsernameConf
	}
	return nil
}

// Defined in table: UsernameConf
type UsernameConfDefine struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID   int32  `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`    // ID
	Name string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"` // ??????
}

func (x *UsernameConfDefine) Reset() {
	*x = UsernameConfDefine{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gameconf_game_stable_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UsernameConfDefine) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UsernameConfDefine) ProtoMessage() {}

func (x *UsernameConfDefine) ProtoReflect() protoreflect.Message {
	mi := &file_gameconf_game_stable_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UsernameConfDefine.ProtoReflect.Descriptor instead.
func (*UsernameConfDefine) Descriptor() ([]byte, []int) {
	return file_gameconf_game_stable_config_proto_rawDescGZIP(), []int{1}
}

func (x *UsernameConfDefine) GetID() int32 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *UsernameConfDefine) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_gameconf_game_stable_config_proto protoreflect.FileDescriptor

var file_gameconf_game_stable_config_proto_rawDesc = []byte{
	0x0a, 0x21, 0x67, 0x61, 0x6d, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f,
	0x73, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x08, 0x67, 0x61, 0x6d, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x1a, 0x18, 0x67,
	0x61, 0x6d, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x96, 0x01, 0x0a, 0x10, 0x47, 0x61, 0x6d, 0x65,
	0x53, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x40, 0x0a, 0x0c,
	0x44, 0x6f, 0x4e, 0x6f, 0x74, 0x55, 0x73, 0x65, 0x54, 0x68, 0x69, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x2e, 0x44, 0x6f,
	0x4e, 0x6f, 0x74, 0x55, 0x73, 0x65, 0x54, 0x68, 0x69, 0x73, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65,
	0x52, 0x0c, 0x44, 0x6f, 0x4e, 0x6f, 0x74, 0x55, 0x73, 0x65, 0x54, 0x68, 0x69, 0x73, 0x12, 0x40,
	0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x2e,
	0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x44, 0x65, 0x66, 0x69,
	0x6e, 0x65, 0x52, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x43, 0x6f, 0x6e, 0x66,
	0x22, 0x38, 0x0a, 0x12, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x43, 0x6f, 0x6e, 0x66,
	0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f,
	0x67, 0x61, 0x6d, 0x65, 0x63, 0x6f, 0x6e, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gameconf_game_stable_config_proto_rawDescOnce sync.Once
	file_gameconf_game_stable_config_proto_rawDescData = file_gameconf_game_stable_config_proto_rawDesc
)

func file_gameconf_game_stable_config_proto_rawDescGZIP() []byte {
	file_gameconf_game_stable_config_proto_rawDescOnce.Do(func() {
		file_gameconf_game_stable_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_gameconf_game_stable_config_proto_rawDescData)
	})
	return file_gameconf_game_stable_config_proto_rawDescData
}

var file_gameconf_game_stable_config_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_gameconf_game_stable_config_proto_goTypes = []interface{}{
	(*GameStableConfig)(nil),   // 0: gameconf.GameStableConfig
	(*UsernameConfDefine)(nil), // 1: gameconf.UsernameConfDefine
	(*DoNotUseThisDefine)(nil), // 2: gameconf.DoNotUseThisDefine
}
var file_gameconf_game_stable_config_proto_depIdxs = []int32{
	2, // 0: gameconf.GameStableConfig.DoNotUseThis:type_name -> gameconf.DoNotUseThisDefine
	1, // 1: gameconf.GameStableConfig.UsernameConf:type_name -> gameconf.UsernameConfDefine
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_gameconf_game_stable_config_proto_init() }
func file_gameconf_game_stable_config_proto_init() {
	if File_gameconf_game_stable_config_proto != nil {
		return
	}
	file_gameconf_game_type_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_gameconf_game_stable_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GameStableConfig); i {
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
		file_gameconf_game_stable_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UsernameConfDefine); i {
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
			RawDescriptor: file_gameconf_game_stable_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gameconf_game_stable_config_proto_goTypes,
		DependencyIndexes: file_gameconf_game_stable_config_proto_depIdxs,
		MessageInfos:      file_gameconf_game_stable_config_proto_msgTypes,
	}.Build()
	File_gameconf_game_stable_config_proto = out.File
	file_gameconf_game_stable_config_proto_rawDesc = nil
	file_gameconf_game_stable_config_proto_goTypes = nil
	file_gameconf_game_stable_config_proto_depIdxs = nil
}
