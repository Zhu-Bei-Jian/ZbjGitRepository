package msgpackager

import (
	"encoding/binary"
	"io"

	"sanguosha.com/baselib/network/crypto"
)

//
// msg struct/msg packet
// ----------------------------------------
// | extlen | msglen | id | ext | msg |
// ----------------------------------------
// |          head        |    body   |
// ----------------------------------------
// head 是 可根据不同的protocol来设定
//

// MsgPackager 管理协议的组织
type MsgPackager interface {
	ReadMsg(reader io.Reader, crypto crypto.Crypto) (id uint32, extdata []byte, msgdata []byte, error error)
	WriteMsg(writer io.Writer, id uint32, extdata []byte, msgdata []byte, crypto crypto.Crypto) error
}

var (
	// BigEndian ...
	BigEndian = binary.ByteOrder(binary.BigEndian)
	// LittleEndian ...
	LittleEndian = binary.ByteOrder(binary.LittleEndian)
)

const (
	// MessageIDSize = 4个字节长度, 不能改
	MessageIDSize = 4
	// MessageLenSize 消息头中表示消息长度的字节的大小
	MessageLenSize = 2
	// MessageMaxLen 消息最大长度
	MessageMaxLen = 102400
)
