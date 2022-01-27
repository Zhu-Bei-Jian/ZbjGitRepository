package pio

import (
	"encoding/binary"
	"errors"
	"io"
	"time"

	"sanguosha.com/baselib/network/crypto"

	"github.com/gorilla/websocket"
)

type webSocketConn struct {
	*websocket.Conn
	r io.Reader
}

func (c *webSocketConn) Read(b []byte) (n int, err error) {
	if c.r == nil {
		typ, r, err := c.NextReader()
		if err != nil {
			return 0, err
		}
		if typ != websocket.BinaryMessage {
			return 0, errors.New("frame type is not websocket.BinaryMessage")
		}
		c.r = r
	}
	n, err = c.r.Read(b)
	if err != nil && err != io.EOF {
		c.r = nil
		return n, err
	}
	if n < len(b) {
		c.r = nil
		n2, err := c.Read(b[n:])
		return n + n2, err
	}
	return n, err
}

func (c *webSocketConn) Write(b []byte) (n int, err error) {
	err = c.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *webSocketConn) SetDeadline(t time.Time) error {
	err := c.SetReadDeadline(t)
	if err != nil {
		return err
	}
	err = c.SetWriteDeadline(t)
	if err != nil {
		return err
	}
	return nil
}

// NewWebsocket 使用 websocket 协议创建 ProtoIO
func NewWebsocket(addr string, msg2id Msg2IDMapper, id2msg ID2MsgMapper) (ProtoIO, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}
	return New(&webSocketConn{Conn: conn}, msg2id, id2msg, binary.LittleEndian, crypto.NewAesCrypto()), nil
}
