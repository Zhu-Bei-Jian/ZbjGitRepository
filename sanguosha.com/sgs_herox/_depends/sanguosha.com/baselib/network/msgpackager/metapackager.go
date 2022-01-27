package msgpackager

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"sanguosha.com/baselib/network/crypto"
)

// msg struct
// ----------------------------------------
// | extlen | msglen | id | ext | msg |
// ----------------------------------------
// |          head        |    body   |
// ----------------------------------------

// 服务器内部协议
// Each packet has a fix length packet header to present packet length.
type metaPackager struct {
	MsgPackager

	headLen   int    //headlenSize+datalenSize+MessageIDSize占用的字节数
	extMaxLen uint32 // 扩展数据最大长度
	msgMaxLen uint32 // 数据最大长度

	// byte流
	encodeHead func([]byte, uint32, uint32, uint32)
	// byte流 -> headlen, datalen, msgid
	decodeHeadLen   func([]byte) (uint32, uint32)
	decodeHeadMsgID func([]byte) uint32
}

// NewMetaPackager Create a {| headlen | lendata | id | head | msg |} msg.
// The extLenSize 是 extlen 的字节数. extLenSize must is 0、1、2、4
// The msgLenSize 是 msglen 的字节数. msgLenSize must is 1、2、4
func NewMetaPackager(extLenSize int, msgLenSize int, byteOrder binary.ByteOrder) MsgPackager {
	packager := &metaPackager{
		headLen: (extLenSize + msgLenSize + MessageIDSize),
	}

	if extLenSize == 0 {
		packager.extMaxLen = 0
	} else if extLenSize == 1 {
		packager.extMaxLen = math.MaxUint8
	} else if extLenSize == 2 {
		packager.extMaxLen = math.MaxUint16
	} else if extLenSize == 4 {
		packager.extMaxLen = math.MaxUint32
	} else {
		panic("unsupported packet ext len size")
	}
	if packager.extMaxLen > MessageMaxLen {
		packager.extMaxLen = MessageMaxLen
	}

	if msgLenSize == 1 {
		packager.msgMaxLen = math.MaxUint8
	} else if msgLenSize == 2 {
		packager.msgMaxLen = math.MaxUint16
	} else if msgLenSize == 4 {
		packager.msgMaxLen = math.MaxUint32
	} else {
		panic("unsupported packet msg len size")
	}
	if packager.msgMaxLen > MessageMaxLen {
		packager.msgMaxLen = MessageMaxLen
	}

	packager.encodeHead = func(buffer []byte, el uint32, ml uint32, id uint32) {
		var pos int

		if extLenSize == 1 {
			buffer[pos] = byte(el)
			pos = pos + 1
		} else if extLenSize == 2 {
			byteOrder.PutUint16(buffer[pos:], uint16(el))
			pos = pos + 2
		} else if extLenSize == 4 {
			byteOrder.PutUint32(buffer[pos:], uint32(el))
			pos = pos + 4
		}

		if msgLenSize == 1 {
			buffer[pos] = byte(ml)
			pos = pos + 1
		} else if msgLenSize == 2 {
			byteOrder.PutUint16(buffer[pos:], uint16(ml))
			pos = pos + 2
		} else if msgLenSize == 4 {
			byteOrder.PutUint32(buffer[pos:], uint32(ml))
			pos = pos + 4
		}

		byteOrder.PutUint32(buffer[pos:], id)
		pos = pos + MessageIDSize
	}

	packager.decodeHeadLen = func(buffer []byte) (el uint32, ml uint32) {
		var pos int

		if extLenSize == 1 {
			el = uint32(buffer[pos])
			pos = pos + 1
		} else if extLenSize == 2 {
			el = uint32(byteOrder.Uint16(buffer[pos:]))
			pos = pos + 2
		} else if extLenSize == 4 {
			el = uint32(byteOrder.Uint32(buffer[pos:]))
			pos = pos + 4
		}

		if msgLenSize == 1 {
			ml = uint32(buffer[pos])
			pos = pos + 1
		} else if msgLenSize == 2 {
			ml = uint32(byteOrder.Uint16(buffer[pos:]))
			pos = pos + 2
		} else if msgLenSize == 4 {
			ml = uint32(byteOrder.Uint32(buffer[pos:]))
			pos = pos + 4
		}

		return el, ml
	}

	packager.decodeHeadMsgID = func(buffer []byte) uint32 {
		return byteOrder.Uint32(buffer)
	}

	return packager
}

// ReadMsg ...
func (p *metaPackager) ReadMsg(reader io.Reader, crypto crypto.Crypto) (uint32, []byte, []byte, error) {
	// head with msgid
	msgHeadWithoutMsgID := make([]byte, p.headLen-MessageIDSize)

	// read head
	if n, err := io.ReadFull(reader, msgHeadWithoutMsgID); err != nil {
		if !(err == io.EOF && n == len(msgHeadWithoutMsgID)) {
			return 0, nil, nil, err
		}
	}

	extLen, msgLen := p.decodeHeadLen(msgHeadWithoutMsgID)

	if extLen > p.extMaxLen {
		return 0, nil, nil, errors.New("read ext too max")
	}

	if msgLen > p.msgMaxLen {
		return 0, nil, nil, errors.New("read msg too max")
	}

	// body
	msgBodyWithMsgID := make([]byte, MessageIDSize+extLen+msgLen)

	if n, err := io.ReadFull(reader, msgBodyWithMsgID); err != nil {
		if !(err == io.EOF && n == len(msgHeadWithoutMsgID)) {
			return 0, nil, nil, err
		}
	}

	if crypto != nil {
		err := crypto.Decrypt(msgBodyWithMsgID, msgBodyWithMsgID)
		if err != nil {
			return 0, nil, nil, errors.New("read msg decrypt error")
		}
	}

	msgid := p.decodeHeadMsgID(msgBodyWithMsgID)

	body := msgBodyWithMsgID[MessageIDSize:]

	// ext
	if p.extMaxLen != 0 && extLen > 0 {
		//fmt.Println("--------------------read, type:", p.extType, ", msg:", msgBody)
		return msgid, body[:extLen], body[extLen:], nil
	}

	return msgid, nil, body, nil
}

// WriteMsg ...
func (p *metaPackager) WriteMsg(writer io.Writer, id uint32, extdata []byte, msgdata []byte, crypto crypto.Crypto) error {
	msgLen := uint32(len(msgdata))

	if msgLen > p.msgMaxLen {
		return fmt.Errorf("write msgdata too max msgid: %d, len: %d", id, msgLen)
	}

	var extLen uint32
	if p.extMaxLen != 0 && extdata != nil {
		extLen = uint32(len(extdata))
	}

	// new buffer
	buffer := make([]byte, uint32(p.headLen)+extLen+msgLen)

	// write head
	p.encodeHead(buffer, extLen, msgLen, id)
	pos := (uint32)(p.headLen)

	// write ext
	if extLen > 0 {
		copy(buffer[pos:], extdata)
		pos = pos + extLen
	}

	// write msg
	copy(buffer[pos:], msgdata)

	if crypto != nil {
		startPos := p.headLen - MessageIDSize
		crypto.Encrypt(buffer[startPos:], buffer[startPos:])
	}

	// write to io
	if _, err := writer.Write(buffer); err != nil {
		return err
	}

	//fmt.Println("--------------------send, type:", p.extType, ", msg:", buffer)

	return nil
}
