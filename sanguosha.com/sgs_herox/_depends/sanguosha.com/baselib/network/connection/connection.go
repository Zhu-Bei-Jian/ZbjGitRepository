package connection

import (
	"net"
)

// Connection 连接
type Connection interface {
	Name() string

	WriteMsg(ext interface{}, msg interface{}) error
	WriteBytes(ext interface{}, msgid uint32, bytes []byte) error

	Close()
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}
