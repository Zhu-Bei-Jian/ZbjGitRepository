package netframe

import (
	"reflect"
	"sync/atomic"
	"time"

	"sanguosha.com/baselib/ioservice"
	"sanguosha.com/baselib/network/connection"
	"sanguosha.com/baselib/network/crypto"
	"sanguosha.com/baselib/network/msgprocessor"
	"sanguosha.com/baselib/network/tcp"
	"sanguosha.com/baselib/network/websocket"
)

// ServerConnectHandler ...
// 如果是GT的话 ServerId/GateConnId = 0, ServerType = 0
type ServerConnectHandler func(conn connection.Connection, ServerID uint32, ServerType uint32, ServerStartTime int64)

// WorkerExitHandler ...
type WorkerExitHandler func(conn connection.Connection, ServerStartTime int64)

var indexGateConnID uint32 = gateConnIDBegin

func genGateConnID() uint32 {
	newid := atomic.AddUint32(&indexGateConnID, 1)
	if newid <= gateConnIDBegin {
		newid = atomic.AddUint32(&indexGateConnID, 1)
	}

	return newid
}

// Server ...
type Server struct {
	closeServer func()

	appConfig    *AppConfig
	serverConfig *ServerConfig

	handlerIO ioservice.IOService
	// 初始化之后, 该map只在主mainworker线程中访问
	heartCheckTimers map[connection.Connection]*time.Timer

	msgHandlers                      msgprocessor.GetMsgHandler
	msgClientReqHeartBeatHandler     msgprocessor.MsgHandler
	msgServerReqHeartBeatHandler     msgprocessor.MsgHandler
	msgServerReqRegisterHandler      msgprocessor.MsgHandler
	msgServerReportUnRegisterHandler msgprocessor.MsgHandler
}

// NewGateServer 创建监听服务器
func NewGateServer(appConfig *AppConfig, serverConfig *ServerConfig, io ioservice.IOService,
	connectHandler ServerConnectHandler, closeHandler msgprocessor.CloseHandler,
	bytesHandler msgprocessor.BytesHandler, msgHandlers msgprocessor.GetMsgHandler) *Server {
	s := new(Server)

	if io == nil {
		logger.Error("NewGateServer io == nil")
		return nil
	}

	s.appConfig = appConfig
	s.serverConfig = serverConfig
	s.handlerIO = io
	s.heartCheckTimers = make(map[connection.Connection]*time.Timer)

	fconnect := func(conn connection.Connection) {
		connid := genGateConnID()

		//心跳时间设置
		s.heartCheckTimers[conn] = time.AfterFunc(heartBeatWaitTime2, func() {
			conn.Close()
			logger.WithField("serveraddr", conn.RemoteAddr()).Warning("gate server heartbeat timeout")
		})

		if connectHandler != nil {
			connectHandler(conn, connid, 0, 0)
		}
	}

	fclose := func(conn connection.Connection) {
		if closeHandler != nil {
			closeHandler(conn)
		}

		s.heartCheckTimers[conn].Stop()
		delete(s.heartCheckTimers, conn)
	}

	gateProcessor := msgprocessor.NewMetaProcessor(nil, io)
	gateProcessor.ConnectHandler = fconnect
	gateProcessor.CloseHandler = fclose
	gateProcessor.MsgHandlers = s
	gateProcessor.BytesHandler = bytesHandler

	s.msgClientReqHeartBeatHandler = func(conn connection.Connection, ext interface{},_ uint32,_ []byte, msg interface{}) {
		req := msg.(*Client_ReqHeartBeat)

		resp := &Client_RespHeartBeat{
			TimeStamp: req.TimeStamp,
		}

		s.SendGateMsg(conn, resp)

		// 心跳记录重置.
		if timer, exist := s.heartCheckTimers[conn]; exist {
			timer.Reset(heartBeatWaitTime2)
		}
	}

	s.msgHandlers = msgHandlers

	newCrypto := func() crypto.Crypto {
		return crypto.NewAesCrypto()
	}
	if s.serverConfig.DisableCrypto {
		newCrypto = func() crypto.Crypto {
			return nil
		}
	}

	//只支持单个类型
	if s.serverConfig.UseWebsocket {
		svr, err := websocket.NewServer(
			s.serverConfig.Name,
			s.serverConfig.ListenAddr,
			s.serverConfig.MaxConnCnt,
			tcpGateWriteChanLen,
			true,
			gatePackager,
			gateProcessor,
			newCrypto,
			s.serverConfig.DisableWSCheckOrigin,
			s.serverConfig.OpenTLS,
			s.serverConfig.CertFile,
			s.serverConfig.KeyFile,
		)

		if err != nil {
			return nil
		}
		s.closeServer = svr.Close
	} else {
		svr := tcp.NewTCPServer(
			s.serverConfig.Name,
			s.serverConfig.ListenAddr,
			s.serverConfig.MaxConnCnt,
			tcpGateWriteChanLen,
			true,
			gatePackager,
			gateProcessor,
			newCrypto)
		if svr == nil {
			return nil
		}
		s.closeServer = svr.Close
	}

	return s
}

// NewCommonServer ...
func NewCommonServer(appConfig *AppConfig, serverConfig *ServerConfig, io ioservice.IOService,
	connectHandler ServerConnectHandler, closeHandler msgprocessor.CloseHandler,
	bytesHandler msgprocessor.BytesHandler, msgHandlers msgprocessor.GetMsgHandler, exitHandler WorkerExitHandler) *Server {
	s := new(Server)

	if io == nil {
		logger.Error("NewCommonServer io == nil")
		return nil
	}

	s.appConfig = appConfig
	s.serverConfig = serverConfig
	s.handlerIO = io
	s.heartCheckTimers = make(map[connection.Connection]*time.Timer)

	fconnect := func(conn connection.Connection) {
		//心跳时间设置
		s.heartCheckTimers[conn] = time.AfterFunc(heartBeatWaitTime, func() {
			conn.Close()
			logger.WithField("clientaddr", conn.RemoteAddr()).Warning("common server heartbeat timeout")
		})
	}

	fclose := func(conn connection.Connection) {
		if closeHandler != nil {
			closeHandler(conn)
		}

		s.heartCheckTimers[conn].Stop()
		delete(s.heartCheckTimers, conn)
	}

	commonProcessor := msgprocessor.NewMetaProcessor((*Server_Extend)(nil), io)
	commonProcessor.ConnectHandler = fconnect
	commonProcessor.CloseHandler = fclose
	commonProcessor.MsgHandlers = s
	commonProcessor.BytesHandler = bytesHandler

	s.msgServerReqHeartBeatHandler = func(conn connection.Connection, ext interface{},_ uint32,_ []byte, msg interface{}) {
		req := msg.(*Server_ReqHeartBeat)

		resp := &Server_RespHeartBeat{
			TimeStamp: req.TimeStamp,
		}

		s.SendCommonServerMsg(conn, appConfig.ServerID, 0,0, resp)

		// 心跳记录重置.
		if timer, exist := s.heartCheckTimers[conn]; exist {
			timer.Reset(heartBeatWaitTime)
		}
	}
	s.msgServerReqRegisterHandler = func(conn connection.Connection, ext interface{},_ uint32,_ []byte, msg interface{}) {
		req := msg.(*Server_ReqRegister)

		resp := &Server_RespRegister{
			ServerType:      s.appConfig.ServerType,
			ServerId:        s.appConfig.ServerID,
			ServerStartTime: s.appConfig.StartTime,
		}
		s.SendCommonServerMsg(conn, appConfig.ServerID, 0,0, resp)

		if connectHandler != nil {
			connectHandler(conn, req.ClientId, req.ClientType, req.ClientStartTime)
		}
	}
	s.msgServerReportUnRegisterHandler = func(conn connection.Connection, ext interface{},_ uint32,_ []byte, msg interface{}) {
		rpt := msg.(*Server_ReportUnRegister)

		if exitHandler != nil {
			exitHandler(conn, rpt.ServerStartTime)
		}
	}

	s.msgHandlers = msgHandlers

	if s.serverConfig.UseWebsocket {
		svr, err := websocket.NewServer(
			s.serverConfig.Name,
			s.serverConfig.ListenAddr,
			s.serverConfig.MaxConnCnt,
			tcpCommonWriteChanLen,
			true,
			commonPackager,
			commonProcessor,
			nil,
			s.serverConfig.DisableWSCheckOrigin,
			s.serverConfig.OpenTLS,
			s.serverConfig.CertFile,
			s.serverConfig.KeyFile,
		)
		if err != nil {
			return nil
		}
		s.closeServer = svr.Close
	} else {
		svr := tcp.NewTCPServer(
			s.serverConfig.Name,
			s.serverConfig.ListenAddr,
			s.serverConfig.MaxConnCnt,
			tcpCommonWriteChanLen,
			true,
			commonPackager,
			commonProcessor,
			nil,
		)
		if svr == nil {
			return nil
		}
		s.closeServer = svr.Close
	}

	return s
}

// Close ...
func (s *Server) Close() {
	s.closeServer()
}

// SendGateMsg 发送消息
func (s *Server) SendGateMsg(conn connection.Connection, msg interface{}) {
	if conn != nil {
		conn.WriteMsg(nil, msg)
	}
}

// SendCommonServerMsg ...
func (s *Server) SendCommonServerMsg(conn connection.Connection, serverId uint32,sessId uint64, userId uint64, msg interface{}) {
	if conn != nil {
		conn.WriteMsg(&Server_Extend{ServerId: serverId, SessionId:sessId,UserId: userId}, msg)
	}
}

var (
	msgTypeClientReqHeartBeat     = reflect.TypeOf((*Client_ReqHeartBeat)(nil))
	msgTypeServerReqHeartBeat     = reflect.TypeOf((*Server_ReqHeartBeat)(nil))
	msgTypeServerReqRegister      = reflect.TypeOf((*Server_ReqRegister)(nil))
	msgTypeServerReportUnRegister = reflect.TypeOf((*Server_ReportUnRegister)(nil))
)

// GetMsgHandler ...
func (s *Server) GetMsgHandler(typ reflect.Type) (msgprocessor.MsgHandler, bool) {
	switch typ {
	case msgTypeClientReqHeartBeat:
		if s.msgClientReqHeartBeatHandler == nil {
			return nil, false
		}
		return s.msgClientReqHeartBeatHandler, true
	case msgTypeServerReqHeartBeat:
		if s.msgServerReqHeartBeatHandler == nil {
			return nil, false
		}
		return s.msgServerReqHeartBeatHandler, true
	case msgTypeServerReqRegister:
		if s.msgServerReqRegisterHandler == nil {
			return nil, false
		}
		return s.msgServerReqRegisterHandler, true
	case msgTypeServerReportUnRegister:
		if s.msgServerReportUnRegisterHandler == nil {
			return nil, false
		}
		return s.msgServerReportUnRegisterHandler, true
	default:
		return s.msgHandlers.GetMsgHandler(typ)
	}
}
