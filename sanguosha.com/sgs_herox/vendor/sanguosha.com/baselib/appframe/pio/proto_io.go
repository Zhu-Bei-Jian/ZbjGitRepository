package pio

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"sanguosha.com/baselib/network/crypto"
)

var (
	// ErrMsg2ID 无法找到 msg 到 id 的映射
	ErrMsg2ID = errors.New("ErrMsg2ID")
	// ErrID2Msg 无法找到 id 到 msg 的映射
	ErrID2Msg = errors.New("ErrID2Msg")
	// ErrTimeout 超时错误
	ErrTimeout = errors.New("ErrTimeout")
)

// ProtoIO 使用 proto.Message 进行 io 通信的接口抽象.
type ProtoIO interface {
	Read() (proto.Message, error)
	ReadTimeout(time.Duration) (proto.Message, error)
	Write(proto.Message) error
	WriteTimeout(proto.Message, time.Duration) error
	Close() error
}

type protoIO struct {
	conn net.Conn
	wmtx sync.Mutex

	byteOrder binary.ByteOrder
	crypto    crypto.Crypto

	msg2id Msg2IDMapper
	id2msg ID2MsgMapper
}

// Msg2IDMapper msg to id mapper
type Msg2IDMapper func(msg proto.Message) (id uint32, exist bool)

// ID2MsgMapper id to msg mapper
type ID2MsgMapper func(id uint32) (msg proto.Message, exist bool)

// New 创建一个 proto io 对象
func New(conn net.Conn, msg2id Msg2IDMapper, id2msg ID2MsgMapper, byteOrder binary.ByteOrder, crypto crypto.Crypto) ProtoIO {
	return &protoIO{
		conn:      conn,
		byteOrder: byteOrder,
		crypto:    crypto,
		msg2id:    msg2id,
		id2msg:    id2msg,
	}
}

// NewTCP 使用 tcp 模式建立连接
func NewTCP(addr string, msg2id Msg2IDMapper, id2msg ID2MsgMapper) (ProtoIO, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return New(conn, msg2id, id2msg, binary.LittleEndian, crypto.NewAesCrypto()), nil
}

func (pio *protoIO) Read() (proto.Message, error) {
	return pio.read(time.Time{})
}

func (pio *protoIO) ReadTimeout(d time.Duration) (proto.Message, error) {
	return pio.read(time.Now().Add(d))
}

func (pio *protoIO) read(deadline time.Time) (proto.Message, error) {
	err := pio.conn.SetReadDeadline(deadline)
	if err != nil {
		return nil, err
	}

	head := [2]byte{}
	_, err = io.ReadFull(pio.conn, head[:])
	if err != nil {
		return nil, err
	}

	bodyLen := pio.byteOrder.Uint16(head[0:2]) + 4

	body := make([]byte, bodyLen)

	n, err := io.ReadFull(pio.conn, body)
	if err != nil {
		return nil, err
	} else if n != int(bodyLen) {
		return nil, errors.New("Read Data Error")
	}

	if pio.crypto != nil {
		err = pio.crypto.Decrypt(body, body)
		if err != nil {
			return nil, err
		}
	}

	msgid := pio.byteOrder.Uint32(body[:4])

	msg, ok := pio.id2msg(msgid)
	if !ok {
		return nil, ErrID2Msg
	}

	err = proto.Unmarshal(body[4:], msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (pio *protoIO) Write(msg proto.Message) error {
	return pio.write(msg, time.Time{})
}

func (pio *protoIO) WriteTimeout(msg proto.Message, d time.Duration) error {
	return pio.write(msg, time.Now().Add(d))
}

func (pio *protoIO) write(msg proto.Message, deadline time.Time) error {
	msgid, ok := pio.msg2id(msg)
	if !ok {
		return ErrMsg2ID
	}

	msgdata, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	data := make([]byte, 6+len(msgdata))
	pio.byteOrder.PutUint16(data[:2], uint16(len(msgdata)))

	body := data[2:]
	pio.byteOrder.PutUint32(body[:4], msgid)
	copy(body[4:], msgdata)

	pio.wmtx.Lock()
	defer pio.wmtx.Unlock()

	if pio.crypto != nil {
		pio.crypto.Encrypt(body, body)
	}

	err = pio.conn.SetWriteDeadline(deadline)
	if err != nil {
		return err
	}

	n, err := pio.conn.Write(data)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return ErrTimeout
		}
		return err
	} else if n != len(data) {
		return errors.New("write data error")
	}

	return nil
}

func (pio *protoIO) Close() error {
	return pio.conn.Close()
}
