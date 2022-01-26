// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gameconf/game_stable_config.proto

package gameconf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Defined in table: GameStableConfig
type GameStableConfig struct {
	DoNotUseThis []*DoNotUseThisDefine `protobuf:"bytes,1,rep,name=DoNotUseThis" json:"DoNotUseThis,omitempty"`
	UsernameConf []*UsernameConfDefine `protobuf:"bytes,2,rep,name=UsernameConf" json:"UsernameConf,omitempty"`
}

func (m *GameStableConfig) Reset()                    { *m = GameStableConfig{} }
func (m *GameStableConfig) String() string            { return proto.CompactTextString(m) }
func (*GameStableConfig) ProtoMessage()               {}
func (*GameStableConfig) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *GameStableConfig) GetDoNotUseThis() []*DoNotUseThisDefine {
	if m != nil {
		return m.DoNotUseThis
	}
	return nil
}

func (m *GameStableConfig) GetUsernameConf() []*UsernameConfDefine {
	if m != nil {
		return m.UsernameConf
	}
	return nil
}

// Defined in table: UsernameConf
type UsernameConfDefine struct {
	ID   int32  `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=Name" json:"Name,omitempty"`
}

func (m *UsernameConfDefine) Reset()                    { *m = UsernameConfDefine{} }
func (m *UsernameConfDefine) String() string            { return proto.CompactTextString(m) }
func (*UsernameConfDefine) ProtoMessage()               {}
func (*UsernameConfDefine) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *UsernameConfDefine) GetID() int32 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *UsernameConfDefine) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*GameStableConfig)(nil), "gameconf.GameStableConfig")
	proto.RegisterType((*UsernameConfDefine)(nil), "gameconf.UsernameConfDefine")
}

func init() { proto.RegisterFile("gameconf/game_stable_config.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 189 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4c, 0x4f, 0xcc, 0x4d,
	0x4d, 0xce, 0xcf, 0x4b, 0xd3, 0x07, 0x31, 0xe2, 0x8b, 0x4b, 0x12, 0x93, 0x72, 0x52, 0xe3, 0x41,
	0x02, 0x99, 0xe9, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0x1c, 0x30, 0x25, 0x52, 0x12, 0xa8,
	0x8a, 0x4b, 0x2a, 0x0b, 0x52, 0x21, 0x6a, 0x94, 0xa6, 0x31, 0x72, 0x09, 0xb8, 0x27, 0xe6, 0xa6,
	0x06, 0x83, 0xf5, 0x3b, 0x83, 0xb5, 0x0b, 0x39, 0x70, 0xf1, 0xb8, 0xe4, 0xfb, 0xe5, 0x97, 0x84,
	0x16, 0xa7, 0x86, 0x64, 0x64, 0x16, 0x4b, 0x30, 0x2a, 0x30, 0x6b, 0x70, 0x1b, 0xc9, 0xe8, 0xc1,
	0x4c, 0xd1, 0x43, 0x96, 0x75, 0x49, 0x4d, 0xcb, 0xcc, 0x4b, 0x0d, 0x42, 0xd1, 0x01, 0x32, 0x21,
	0xb4, 0x38, 0xb5, 0x28, 0x2f, 0x31, 0x17, 0x6c, 0xa6, 0x04, 0x13, 0xba, 0x09, 0xc8, 0xb2, 0x30,
	0x13, 0x90, 0xc5, 0x94, 0x2c, 0xb8, 0x84, 0x30, 0xd5, 0x08, 0xf1, 0x71, 0x31, 0x79, 0xba, 0x48,
	0x30, 0x2a, 0x30, 0x6a, 0xb0, 0x06, 0x31, 0x79, 0xba, 0x08, 0x09, 0x71, 0xb1, 0xf8, 0x25, 0xe6,
	0xa6, 0x4a, 0x30, 0x29, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x49, 0x6c, 0x60, 0x9f, 0x19, 0x03,
	0x02, 0x00, 0x00, 0xff, 0xff, 0x35, 0x2d, 0x96, 0x4d, 0x22, 0x01, 0x00, 0x00,
}
