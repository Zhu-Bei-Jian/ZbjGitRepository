package netframe

import (
	"net"
	"time"

	"github.com/sirupsen/logrus"

	"math"

	"sanguosha.com/baselib/ioservice"
	"sanguosha.com/baselib/network/msgpackager"
	"sanguosha.com/baselib/network/msgprocessor"
)

var logger = logrus.WithField("module", "netframe")

// private -------------------------------------
const (
	// 服务id范围 [serverIDBegin, serverIDEnd]
	serverIDBegin = 0
	serverIDEnd   = 100000
	// 客户端连接id范围 (gateConnIDBegin, gateConnIDEnd]
	gateConnIDBegin = serverIDEnd
	gateConnIDEnd   = math.MaxUint32
)

var (
	gatePackager   = msgpackager.NewMetaPackager(0, 2, msgpackager.LittleEndian)
	commonPackager = msgpackager.NewMetaPackager(1, 2, msgpackager.LittleEndian)

	heartBeatInterval  = time.Second * 40      // 心跳发送间隔时间
	heartBeatWaitTime  = heartBeatInterval * 3 // 心跳超时间隔时间
	heartBeatWaitTime2 = heartBeatInterval * 5 // gate客户端心跳超时等待时间

	autoReconnetInterval = time.Second * 5 // 自动重连时间

	tcpGateWriteChanLen   = 10240
	tcpCommonWriteChanLen = 40960
)

// private -------------------------------------

// public --------------------------------------

// IsServerID 检查是否是服务ID
func IsServerID(ID uint32) bool {
	if ID >= serverIDBegin && ID <= serverIDEnd {
		return true
	}
	return false
}

// AppConfig ...
type AppConfig struct {
	ServerID   uint32 `json:"ServerID"`
	ServerType uint32 `json:"ServerType"`
	StartTime  int64  `json:"StartTime"`
}

// ClientConfig 客户端配置
type ClientConfig struct {
	Name        string `json:"Name"`
	ConnectAddr string `json:"ConnectAddr"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Name         string `json:"Name"`
	UseWebsocket bool
	OpenTLS      bool   //是否开启wss
	CertFile     string //证书文件路径
	KeyFile      string //key文件路径

	ListenAddr    string `json:"ListenAddr"`
	MaxConnCnt    int    `json:"MaxConnCnt"`
	DisableCrypto bool

	DisableWSCheckOrigin bool
}

// OnNetConnect ...
type OnNetConnect func(ID uint32, serverType uint32)

// OnNetClose ...
type OnNetClose func(ID uint32, serverType uint32)

// OnNetBytes ...
type OnNetBytes func(ID uint32, serverType uint32,  msgid uint32, bytes []byte,extend Server_Extend)

// OnNetMessage ...
type OnNetMessage func(ID uint32, serverType uint32,msgId uint32,bytes []byte, msg interface{},extend Server_Extend)

// OnWorkerExit ...
type OnWorkerExit func(connID uint32, serverType uint32)

// MetaNet ...
type MetaNet interface {
	Init(config *AppConfig, io ioservice.IOService)
	Connect(config *ClientConfig)
	Listen(config *ServerConfig, isGate bool)
	Fini()

	Post(f func())

	// Listen系列 必须在Connect和Listen之前调用
	// 监听网络连接成功
	ListenConnect(fconnect OnNetConnect)
	// 监听网络断开
	ListenClose(fclose OnNetClose)
	// 监听收到byte流
	ListenBytes(fbytes OnNetBytes)
	// 监听收到消息，一个消息能被多个func监听，调用次序按注册次序
	ListenMessage(msg interface{}, fmessage OnNetMessage)
	// 监听工作端强退，即server的client，对应集群中的slave
	ListenWorkerExit(fexit OnWorkerExit)

	LocalAddr(ID uint32) net.Addr
	RemoteAddr(ID uint32) net.Addr

	SendGateMsg(ID uint32, msg interface{}) error
	SendGateBytes(ID uint32, msgid uint32, bytes []byte) error
	SendMsg(ID uint32, msg interface{},extend *Server_Extend) error
	SendBytes(ID uint32, msgid uint32, bytes []byte,extend *Server_Extend) error

	Close(ID uint32)
}

// public --------------------------------------

func initMessage() {
	msgprocessor.RegisterMessage((*Server_ReqHeartBeat)(nil))
	msgprocessor.RegisterMessage((*Server_RespHeartBeat)(nil))

	msgprocessor.RegisterMessage((*Server_ReqRegister)(nil))
	msgprocessor.RegisterMessage((*Server_RespRegister)(nil))
	msgprocessor.RegisterMessage((*Server_ReportUnRegister)(nil))

	msgprocessor.RegisterMessage((*Client_ReqHeartBeat)(nil))
	msgprocessor.RegisterMessage((*Client_RespHeartBeat)(nil))
}
