package netcluster

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"sanguosha.com/baselib/network/msgprocessor"

	"github.com/sirupsen/logrus"

	"strings"
	"sync"

	"sanguosha.com/baselib/framework/netframe"
	"sanguosha.com/baselib/ioservice"
	"sanguosha.com/baselib/util"
)

const (
	shiftMasterStartTimeMin = 5
	shiftMasterStartTimeMax = 10
	shiftMsgQueueCap        = 1000
	shiftMasterWaitTime     = time.Second * 2
)


type MsgSrc string

const (
	//被动接收的消息 来源：客户端
	MsgSrcIn_Client MsgSrc = "in_c"
	//主动向客户端发消息
	MsgSrcOut_Client MsgSrc = "out_c"
	//被动接收的消息 来源：其它服务器
	MsgSrcIn_Server MsgSrc = "in_s"
	//主动向其它服务发消息
	MsgSrcOut_Server MsgSrc = "out_s"
)

//type MsgHead struct{
//	ServerId uint32
//	SessId uint64
//	SeqId int64
//	UserId uint64
//	ExtendId int64
//}

// ConnMasterInfo master服务的数据
type ConnMasterInfo struct {
	config *MasterConf
	loadLv uint32
}

// RouterStatus 通路状态
type RouterStatus struct {
	masterID    uint32
	isConnected bool // 是否已经连接上
	isWorking   bool // 是否已经在工作
}

// SubServerInfo 关注服务信息
type SubServerInfo struct {
	config     *SlaveConf
	fMasterID  uint32                   //0时未初始化，后就不再是0
	mid2Status map[uint32]*RouterStatus //关注的slave连的m和我的m交集
	//
	isInitOk       bool     //是否已经初始化连接成功
	shiftMsgQueue  msgQueue //切换消息缓冲队列
	shiftingMaster bool     //是否真正切换
	shiftMsgNewPri int64    //消息优先级
	//tmpShiftTargetM uint32					//临时切换目标 检测是否稳定使用
}

// SubServerGroup 关注服务组（按type分类）
type SubServerGroup struct {
	serverType    uint32
	id2SubServers map[uint32]*SubServerInfo
	// 服务状态回调
	handler SvrEventHandler
}

// Slave ...
type Slave struct {
	NetIO                ioservice.IOService
	ClusterConfig        *ClusterConf
	SlaveConfig          *SlaveConf
	DebugPrintMessage    bool
	PrintLoadLevelStatus bool

	svrmutex sync.RWMutex
	// master
	masterType uint32
	id2Masters map[uint32]*ConnMasterInfo //我连接的m
	// 我关注的服务 type->group
	mySubscribers map[uint32]*SubServerGroup
	// 所有关注服务用id索引
	id2Subscribers map[uint32]*SubServerInfo
	// 网络
	net netframe.MetaNet
	// fix master

	// master 回调
	masterNetConnect netframe.OnNetConnect
	masterNetClose   netframe.OnNetClose
	masterBytes      netframe.OnNetBytes

	// client 回调
	clientNetConnect netframe.OnNetConnect
	clientNetClose   netframe.OnNetClose
	clientBytes      netframe.OnNetBytes

	interceptFunc func(source MsgSrc,connId uint32,msgId uint32,msgData []byte,extend *netframe.Server_Extend)
}

// NewSlave ...
func NewSlave(config *ClusterConf, key string, io ioservice.IOService) *Slave {
	keyConfig, ok := config.Slaves[key]
	if !ok {
		logger.WithField("server", key).Panic("NewSlave can not find server")
	}

	slave := &Slave{}

	slave.NetIO = io
	slave.ClusterConfig = config
	slave.SlaveConfig = keyConfig
	slave.DebugPrintMessage = keyConfig.DebugPrintMessage
	slave.PrintLoadLevelStatus = keyConfig.PrintLoadLevelStatus

	return slave
}

// Init 初始化
func (s *Slave) Init() {
	if s.SlaveConfig.HTTPPProfPort > 0 {
		util.HTTPPProf(s.SlaveConfig.HTTPPProfHost, s.SlaveConfig.HTTPPProfPort)
	}
	// 初始化变量
	s.net = netframe.NewMetaNet()

	s.svrmutex.Lock()
	defer s.svrmutex.Unlock()
	s.id2Masters = make(map[uint32]*ConnMasterInfo)
	s.mySubscribers = make(map[uint32]*SubServerGroup)
	s.id2Subscribers = make(map[uint32]*SubServerInfo)

	for _, mid := range s.SlaveConfig.MasterIDs {
		for _, mconf := range s.ClusterConfig.Masters {
			if mid == mconf.ServerID {
				s.masterType = mconf.ServerType
				s.id2Masters[mconf.ServerID] = &ConnMasterInfo{
					config: mconf,
					loadLv: 0,
				}
				break
			}
		}
	}

	// 我关注的服务
	for _, serverType := range s.SlaveConfig.SubscribedTypes {
		group := &SubServerGroup{}

		group.serverType = serverType
		group.id2SubServers = make(map[uint32]*SubServerInfo)

		s.mySubscribers[serverType] = group
	}

	for _, sconf := range s.ClusterConfig.Slaves {
		if sconf == s.SlaveConfig {
			continue
		}

		if group, ok := s.mySubscribers[sconf.ServerType]; ok {
			sitem := &SubServerInfo{config: sconf, fMasterID: 0, isInitOk: false}
			sitem.mid2Status = make(map[uint32]*RouterStatus)

			for _, mid := range sconf.MasterIDs {
				if _, ok := s.id2Masters[mid]; ok {
					sitem.mid2Status[mid] = &RouterStatus{masterID: mid}
				}
			}

			group.id2SubServers[sitem.config.ServerID] = sitem
		}
	}

	for _, g := range s.mySubscribers {
		for _, item := range g.id2SubServers {
			s.id2Subscribers[item.config.ServerID] = item
		}
	}

	// 启动网络
	s.net.Init(&netframe.AppConfig{
		ServerID:   s.SlaveConfig.ServerID,
		ServerType: s.SlaveConfig.ServerType,
		StartTime:  time.Now().Unix(),
	}, s.NetIO)

	// 启动监听
	s.net.ListenConnect(s.OnConnect)
	s.net.ListenClose(s.OnClose)
	s.net.ListenMessage((*Master_ReqVerifyConfigFile)(nil), s.onReqVerifyConfigFile)
	s.net.ListenMessage((*Master_PublishServerStatus)(nil), s.onPublishServerStatus)
	s.net.ListenMessage((*Master_PublishLoadLevel)(nil), s.onPublishLoadLevel)
	s.net.ListenMessage((*Slave_ReqShiftFixMaster)(nil), s.onShiftMasterReq)
	s.net.ListenMessage((*Slave_RepShiftFixMaster)(nil), s.onShiftMasterRep)
	s.net.ListenMessage((*SS_CmdPrepareCloseServer)(nil), s.OnSlaveReqPreCloseServer)
	s.net.ListenBytes(s.OnBytes)
}

// Run ...
func (s *Slave) Run() {
	// 连mids，没有则全连
	for _, mscfg := range s.id2Masters {
		cc := &netframe.ClientConfig{
			Name:        mscfg.config.ServerName,
			ConnectAddr: mscfg.config.ListenAddr,
		}
		s.net.Connect(cc)
	}

	if len(s.SlaveConfig.ListenAddr) > 0 {
		sc := &netframe.ServerConfig{
			Name:                 s.SlaveConfig.ServerName,
			UseWebsocket:         s.SlaveConfig.UseWebsocket,
			OpenTLS:              s.SlaveConfig.OpenTLS,
			CertFile:             s.SlaveConfig.CertFile,
			KeyFile:              s.SlaveConfig.KeyFile,
			ListenAddr:           s.SlaveConfig.ListenAddr,
			MaxConnCnt:           s.SlaveConfig.MaxConnCnt,
			DisableCrypto:        s.SlaveConfig.DisableCrypto,
			DisableWSCheckOrigin: s.SlaveConfig.DisableWSCheckOrigin,
		}
		s.net.Listen(sc, s.IsGate())
	}
}

// Fini ...
func (s *Slave) Fini() {
	s.net.Fini()
}

// Post ...
func (s *Slave) Post(f func()) {
	s.net.Post(f)
}

// LocalAddr ...
func (s *Slave) LocalAddr(ID uint32) net.Addr {
	return s.net.LocalAddr(ID)
}

// RemoteAddr ...
func (s *Slave) RemoteAddr(ID uint32) net.Addr {
	return s.net.RemoteAddr(ID)
}

// SendClientMsg ...
func (s *Slave) SendClientMsg(ID uint32,msg interface{},extend netframe.Server_Extend) error {
	if s.DebugPrintMessage {
		logger.WithFields(logrus.Fields{
			"session": ID,
			"msg":     fmt.Sprintf("%#v", msg),
		}).Debug("Send msg to client")
	}
	//直接序列化
	msgid, data, err1 := msgprocessor.OnMarshal(msg)
	if err1 != nil {
		return err1
	}
	//return s.SendClientBytes(ID,msgid,data)
	return s.sendByteWithIntercepter(ID, msgid, data,&extend)
	//return s.net.SendMsg(ID, 0, 0, msg)
}

// SendClientBytes ...
func (s *Slave) SendClientBytes(ID uint32, msgid uint32, bytes []byte, extend netframe.Server_Extend,) error {
	return s.sendByteWithIntercepter(ID, msgid, bytes, &extend)
}


// SendServerMsg serverID:目标服务器
func (s *Slave) SendServerMsg(msg interface{},extend netframe.Server_Extend) error {
	//if s.DebugPrintMessage {
	//	if extend.SessionId != 0 {
	//		logger.WithFields(logrus.Fields{
	//			"svrid":   extend.ServerId,
	//			"session": extend.SessionId,
	//			"msg":     fmt.Sprintf("%#v", msg),
	//		}).Debug("Send msg to session")
	//	} else {
	//		logger.WithFields(logrus.Fields{
	//			"svrid": extend.ServerId,
	//			"msg":   fmt.Sprintf("%#v", msg),
	//		}).Debug("Send msg to server")
	//	}
	//}

	//直接序列化
	msgid, data, err1 := msgprocessor.OnMarshal(msg)
	if err1 != nil {
		logger.WithFields(logrus.Fields{
			"serverID": extend.ServerId,
			"error":    err1,
			"msgid":    msgid,
			"msg":      fmt.Sprintf("%#v", msg),
		}).Error("Marshal message error")
		return fmt.Errorf("Slave SendServerMsg, OnMarshal error! serverId:%d, error:%s", extend.ServerId, err1)
	}
	return s.SendServerBytes(msgid, data,extend)
}

// SendServerBytes ...
func (s *Slave) SendServerBytes(msgid uint32, bytes []byte,extend netframe.Server_Extend) error {
	s.svrmutex.Lock()
	defer s.svrmutex.Unlock()
	if ssinfo, ok := s.id2Subscribers[extend.ServerId]; ok {
		if ssinfo.shiftingMaster || len(ssinfo.shiftMsgQueue) != 0 {
			ssinfo.shiftMsgNewPri++
			ssinfo.shiftMsgQueue.Push(&wMessage{bytes: bytes, pri: ssinfo.shiftMsgNewPri, extend:extend,msgID: msgid})
			return nil
		}

		return s.sendByteWithIntercepter(ssinfo.fMasterID, msgid, bytes,&extend)
	}

	return fmt.Errorf("Slave SendServerBytes, no serverinfo. id:%d", extend.ServerId)
}

func(s *Slave) sendByteWithIntercepter(ID uint32,msgid uint32, bytes []byte,extend *netframe.Server_Extend)error{
	if s.interceptFunc!=nil{
		msgSrc:=MsgSrcOut_Client
		if netframe.IsServerID(ID){
			msgSrc = MsgSrcOut_Server
		}
		s.interceptFunc(msgSrc,ID,msgid,bytes,extend)
	}
	return s.net.SendBytes(ID, msgid, bytes,extend)
}

// Close ...
func (s *Slave) Close(ID uint32) {
	s.net.Close(ID)
}

// IsGate 是否为网关Slave
func (s *Slave) IsGate() bool {
	return s.SlaveConfig.IsGate()
}

// GetServerRandomly 获取连接ok且正在作用的working
func (s *Slave) GetServerRandomly(serverType uint32) (serverID uint32, err error) {
	s.svrmutex.RLock()
	defer s.svrmutex.RUnlock()
	if group, okg := s.mySubscribers[serverType]; okg {
		servers := make([]uint32, 0)

		for _, sitem := range group.id2SubServers {
			if rs, ok := sitem.mid2Status[sitem.fMasterID]; ok {
				if rs.isConnected && rs.isWorking {
					servers = append(servers, sitem.config.ServerID)
				}
			}
		}

		serverCnt := len(servers)
		if serverCnt <= 0 {
			return 0, fmt.Errorf("ServerType:%d has not connect, Check Config or Wait for A Minutue", serverType)
		}

		i := rand.Int31n((int32)(serverCnt))

		return servers[i], nil
	}

	return 0, fmt.Errorf("Have Not Focus ServerType:%d, Check Config", serverType)
}

// GetServerAllAvailable ..
func (s *Slave) GetServerAllAvailable(serverType uint32) []uint32 {
	s.svrmutex.RLock()
	defer s.svrmutex.RUnlock()
	if group, okg := s.mySubscribers[serverType]; okg {
		var servers []uint32
		for _, sitem := range group.id2SubServers {
			if rs, ok := sitem.mid2Status[sitem.fMasterID]; ok {
				if rs.isConnected && rs.isWorking {
					servers = append(servers, sitem.config.ServerID)
				}
			}
		}
		return servers
	}

	return nil
}

// IsServerAvailable ...
func (s *Slave) IsServerAvailable(svrid uint32) bool {
	s.svrmutex.RLock()
	defer s.svrmutex.RUnlock()
	for _, group := range s.mySubscribers {
		for id, sitem := range group.id2SubServers {
			if id == svrid {
				if rs, ok := sitem.mid2Status[sitem.fMasterID]; ok {
					if rs.isConnected && rs.isWorking {
						return true
					}
				}
				return false
			}
		}
	}
	return false
}

// GetServerType ...
func (s *Slave) GetServerType(svrid uint32) (uint32, bool) {
	s.svrmutex.RLock()
	defer s.svrmutex.RUnlock()
	for typ, group := range s.mySubscribers {
		for id := range group.id2SubServers {
			if id == svrid {
				return typ, true
			}
		}
	}
	return 0, false
}

// ListenMasterNetEvent ...
func (s *Slave) ListenMasterNetEvent(con netframe.OnNetConnect, discon netframe.OnNetClose) {
	s.masterNetConnect = con
	s.masterNetClose = discon
}

// ListenClientNetEvent ...
func (s *Slave) ListenClientNetEvent(con netframe.OnNetConnect, discon netframe.OnNetClose) {
	s.clientNetConnect = con
	s.clientNetClose = discon
}

// ListenClientBytes ...
func (s *Slave) ListenClientBytes(onBytes netframe.OnNetBytes) {
	s.clientBytes = func(ID uint32, serverType uint32,  msgid uint32, bytes []byte,extend netframe.Server_Extend){
		if s.interceptFunc!=nil{
			s.interceptFunc(MsgSrcIn_Client,ID,msgid,bytes,&extend)
		}
		onBytes(ID,serverType,msgid,bytes,extend)
	}
}

func(s *Slave) RegisterIntercepter(f func(source MsgSrc,connId uint32,msgId uint32,msgData []byte,extend *netframe.Server_Extend)){
	s.interceptFunc = f
}

// ListenClientMessage ...
func (s *Slave) ListenClientMessage(msg interface{}, message netframe.OnNetMessage) {
	if message == nil {
		return
	}

	s.net.ListenMessage(msg, func(cID uint32, cServerType uint32,msgId uint32,msgData []byte, cmsg interface{},extend netframe.Server_Extend) {
		if !netframe.IsServerID(cID) {
			if s.DebugPrintMessage {
				logger.WithFields(logrus.Fields{
					"session": cID,
					"msg":     fmt.Sprintf("%#v", cmsg),
				}).Debug("Msg from client")
			}
			if s.interceptFunc!=nil{
				s.interceptFunc(MsgSrcIn_Client,cID,msgId,msgData,&extend)
			}
			message(cID, cServerType,msgId,msgData,cmsg,extend)
		}
	})
}

// ListenServerStatus ...
func (s *Slave) ListenServerStatus(serverType uint32, handler SvrEventHandler) {
	s.svrmutex.RLock()
	defer s.svrmutex.RUnlock()
	if group, ok := s.mySubscribers[serverType]; ok {
		preHandler := group.handler
		group.handler = func(svrid uint32, event SvrEvent) {
			if preHandler != nil {
				preHandler(svrid, event)
			}
			handler(svrid, event)
		}
	}
}

// ListenServerBytes ...
func (s *Slave) ListenServerBytes(onBytes netframe.OnNetBytes) {
	//s.masterBytes = onBytes
	s.masterBytes = func(ID uint32, serverType uint32,  msgid uint32, bytes []byte,extend netframe.Server_Extend){
		if s.interceptFunc!=nil {
			s.interceptFunc(MsgSrcIn_Server,ID, msgid,bytes,&extend)
		}
		onBytes(ID,serverType, msgid,bytes,extend)
	}
}

// ListenServerMessage ...
func (s *Slave) ListenServerMessage(serverType uint32, msg interface{}, message netframe.OnNetMessage) {
	if message == nil {
		return
	}

	s.net.ListenMessage(msg, func(cID uint32, cServerType uint32,msgId uint32,msgData []byte,cmsg interface{},extend netframe.Server_Extend) {
		s.svrmutex.RLock()

		if s.interceptFunc!=nil{
			s.interceptFunc(MsgSrcIn_Server,cID, msgId,msgData,&extend)
		}

		if s.masterType == cServerType {
			if _, ok := s.id2Subscribers[extend.ServerId]; ok {
				s.svrmutex.RUnlock()
				if s.DebugPrintMessage {
					if extend.SessionId != 0 {
						//logger.WithFields(logrus.Fields{
						//	"svrid":   extend.ServerId,
						//	"session": extend.SessionId,
						//	"msg":     fmt.Sprintf("%#v", cmsg),
						//}).Debug("Msg from session")
					} else {
						//logger.WithFields(logrus.Fields{
						//	"svrid": extend.ServerId,
						//	"msg":   fmt.Sprintf("%#v", cmsg),
						//}).Debug("Msg from server")
					}
				}
				message(cID, cServerType,msgId,msgData, cmsg,extend)
				return
			}
		}
		s.svrmutex.RUnlock()
	})
}

// OnConnect ...
func (s *Slave) OnConnect(ID uint32, serverType uint32) {
	if serverType == s.masterType {
		logger.WithFields(logrus.Fields{
			"svrid":   ID,
			"svrtype": serverType,
		}).Info("[Slave] Connect Master Succeed.")

		// 加载新配置连接master，避免断线重连后master更新配置文件
		if newConfig, err := s.ClusterConfig.LoadNewCfgFile(); err == nil {
			if strings.Compare(newConfig.FileMd5, s.ClusterConfig.FileMd5) != 0 && s.canLoadNewConfig(newConfig) {
				s.ClusterConfig = newConfig
			}
		}
		// 检查配置请求
		req := &Slave_ReqVerifyConfigFile{FileMd5: s.ClusterConfig.FileMd5}
		s.net.SendMsg(ID,req,&netframe.Server_Extend{ServerId: ID})

		if s.masterNetConnect != nil {
			s.masterNetConnect(ID, serverType)
		}
	} else if !netframe.IsServerID(ID) {
		if s.clientNetConnect != nil {
			s.clientNetConnect(ID, serverType)
		}
	}
}

// OnClose ...
func (s *Slave) OnClose(ID uint32, serverType uint32) {
	if serverType == s.masterType {
		logger.WithFields(logrus.Fields{
			"svrid":   ID,
			"svrtype": serverType,
		}).Info("[Slave] Close Master Succeed.")
		//清理
		s.svrmutex.Lock()
		for sid, sinfo := range s.id2Subscribers {
			if msta, ok := sinfo.mid2Status[ID]; ok {
				msta.isConnected = false
				msta.isWorking = false
			}
			if sinfo.fMasterID == ID {
				sinfo.fMasterID = 0
				s.fixMaster2SubServer(sid, 1)
			}
		}
		s.svrmutex.Unlock()

		if s.masterNetClose != nil {
			s.masterNetClose(ID, serverType)
		}
	} else if !netframe.IsServerID(ID) {
		if s.clientNetClose != nil {
			s.clientNetClose(ID, serverType)
		}
	}
}

// OnBytes ...
func (s *Slave) OnBytes(ID uint32, serverType uint32,msgid uint32, bytes []byte,extend netframe.Server_Extend) {
	if serverType == s.masterType {
		if s.masterBytes != nil {
			s.masterBytes(ID, serverType,  msgid, bytes,extend)
		}
	} else if !netframe.IsServerID(ID) {
		if s.clientBytes != nil {
			s.clientBytes(ID, serverType, msgid, bytes,extend)
		}
	}
}

func (s *Slave) onReqVerifyConfigFile(ID uint32, serverType uint32, _ uint32,_ []byte, msg interface{},extend netframe.Server_Extend) {
	if serverType != s.masterType {
		return
	}
	req := msg.(*Master_ReqVerifyConfigFile)
	isCfgOk := false

	if strings.Compare(s.ClusterConfig.FileMd5, req.FileMd5) == 0 {
		isCfgOk = true
	} else if newConfig, err := s.ClusterConfig.LoadNewCfgFile(); err == nil {
		if strings.Compare(req.FileMd5, newConfig.FileMd5) == 0 && s.canLoadNewConfig(newConfig) {
			s.ClusterConfig = newConfig
			isCfgOk = true

			// 更新
			for _, mitem := range s.id2Masters {
				s.net.SendMsg(mitem.config.ServerID,&Slave_UptConfigMd5{FileMd5: s.ClusterConfig.FileMd5},&netframe.Server_Extend{ServerId: mitem.config.ServerID})
			}
		}
	}

	s.net.SendMsg(ID, &Slave_RepVerifyConfigFile{IsSucc: isCfgOk, Time: req.Time, ReqServerId: req.ReqServerId, ReqServerType: req.ReqServerType},&netframe.Server_Extend{ServerId: serverType})
}

func (s *Slave) onPublishServerStatus(ID uint32, serverType uint32,_ uint32,_ []byte,msg interface{},extend netframe.Server_Extend) {
	if serverType != s.masterType {
		return
	}

	req := msg.(*Master_PublishServerStatus)

	s.svrmutex.Lock()
	if !s.addNewSalveCfg(req.ServerId, req.ServerType) {
		s.svrmutex.Unlock()
		logger.WithFields(logrus.Fields{
			"svrid":     req.ServerId,
			"svrtype":   req.ServerType,
			"connected": req.IsConnected,
			"working":   req.IsWorking,
		}).Error("[Slave] Other Slave Server Status 加载配置失败!")
		return
	}

	// 关注连接
	if group, okg := s.mySubscribers[req.ServerType]; okg {
		if server, oks := group.id2SubServers[req.ServerId]; oks {
			if rs, ok := server.mid2Status[ID]; ok {
				var event SvrEvent
				if rs.isWorking != req.IsWorking {
					if req.IsWorking {
						event = SvrEventStart
					} else {
						event = SvrEventQuit
					}
				} else if rs.isConnected != req.IsConnected {
					if req.IsConnected {
						event = SvrEventReconnect
					} else {
						event = SvrEventDisconnect
					}
				}

				rs.isConnected = req.IsConnected
				rs.isWorking = req.IsWorking

				//都连上了，选择fixmaster
				if !server.isInitOk {
					server.isInitOk = true
					for _, sta := range server.mid2Status {
						//if !sta.isConnected || !sta.isWorking {
						if !sta.isConnected {
							server.isInitOk = false
							break
						}
					}
					if server.isInitOk {
						if fmid := s.fixMaster2SubServer(req.ServerId, 2); fmid == 0 {
							logger.WithField("svrid", req.ServerId).Panic("OnPublishServerStatus init fix master failed!")
						}
					} else {
						event = 0
					}
				} else {
					if server.fMasterID == 0 {
						//if rs.isConnected && rs.isWorking {
						if rs.isConnected {
							s.fixMaster2SubServer(req.ServerId, 3)
						}
					}
					if server.fMasterID == ID {
						//if !rs.isWorking {
						//	server.fMasterID = 0
						//} else if !rs.isConnected {
						if !rs.isConnected {
							//清理fixmaster
							server.fMasterID = 0
							s.fixMaster2SubServer(req.ServerId, 4)
							if server.fMasterID != 0 {
								event = 0
							}
						}
					} else {
						//if rs.isConnected && rs.isWorking {
						if rs.isConnected {
							s.checkShiftMaster(ID, req.ServerId)
						}
						// TODO: 重复的退出事件
						// 如果服务在断线后重启, 由于先前的断线已经造成 fMasterID 的值归为0, 所以无法判断以哪个master为准,
						// 但退出事件又必须通知给上层, 这是就可能造成上层连续收到多次同一服务的退出事件, 业务层重复的退出事件一般并不会影响逻辑, 这里还是列为TODO.
						if !(event == SvrEventQuit && server.fMasterID == 0) {
							event = 0
						}
					}
				}

				logger.WithFields(logrus.Fields{
					"svrid":     req.ServerId,
					"svrtype":   req.ServerType,
					"connected": req.IsConnected,
					"working":   req.IsWorking,
				}).Info("[Slave] Other Slave Server Status.")
				if event != 0 && group.handler != nil {
					s.svrmutex.Unlock()
					group.handler(server.config.ServerID, event)
					return
				}

				s.svrmutex.Unlock()
				return
			}
		}
	}
	s.svrmutex.Unlock()

	logger.WithFields(logrus.Fields{
		"svrid":     req.ServerId,
		"svrtype":   req.ServerType,
		"connected": req.IsConnected,
		"working":   req.IsWorking,
	}).Error("[Slave] Other Slave Server Status 失败!")
}

// fixMaster2SubServer 基于masterinfo上loadlv设置fix master
func (s *Slave) fixMaster2SubServer(serverID uint32, idx uint32) uint32 {
	if sinfo, ok := s.id2Subscribers[serverID]; ok {
		var tminfo *ConnMasterInfo

		//旧fmasterid是否可用
		if sinfo.fMasterID != 0 {
			if rs, ok := sinfo.mid2Status[sinfo.fMasterID]; ok {
				//if rs.isConnected && rs.isWorking {
				if rs.isConnected {
					return sinfo.fMasterID
				}
			}
		}

		//查找新master
		for mid, mv := range sinfo.mid2Status {
			//if !mv.isConnected || !mv.isWorking {
			if !mv.isConnected {
				continue
			}

			minfo, have := s.id2Masters[mid]
			if !have {
				continue
			}

			if tminfo == nil {
				tminfo = minfo
			} else if minfo.loadLv < tminfo.loadLv {
				tminfo = minfo
			}
		}

		if tminfo != nil {
			tmpw := uint32(5)
			if wv, ok := tminfo.config.SlaveWeights[s.SlaveConfig.ServerType]; ok {
				tmpw = wv
			}
			sinfo.fMasterID = tminfo.config.ServerID
			tminfo.loadLv += tmpw

			//上报master其压力值变化,然后master会同步其所有slave
			msg := &Slave_ReportLoadLevel{
				IsFix:     true,
				AServerID: s.SlaveConfig.ServerID,
				BServerID: serverID,
			}
			s.net.SendMsg(sinfo.fMasterID, msg,&netframe.Server_Extend{ServerId:sinfo.fMasterID})

			//检查是否切换master
			if sinfo.shiftingMaster {
				s.sendMsgQueue(serverID)
			}

			logger.Infof("[Fix]连接. %s,%d ---%d---> %s,%d", s.SlaveConfig.ServerName, s.SlaveConfig.ServerID, sinfo.fMasterID, sinfo.config.ServerName, sinfo.config.ServerID)
			return tminfo.config.ServerID
		}

		logger.WithFields(logrus.Fields{
			"svrid": serverID,
			"idx":   idx,
		}).Warning("fixMaster2SubServer not find fix master for server")
	}

	return 0
}

// sendMsgQueue 发送消息队列
func (s *Slave) sendMsgQueue(sID uint32) {
	if sinfo, ok := s.id2Subscribers[sID]; ok {
		sinfo.shiftingMaster = false
		sinfo.shiftMsgNewPri = 0
		//sinfo.tmpShiftTargetM = 0
		//log.BDebug("[shift sendMsgQueue] -++++++++++++len of msgq:%d", len(sinfo.shiftMsgQueue))
		for {
			wmsg := sinfo.shiftMsgQueue.Pop()
			if wmsg == nil {
				break
			}
			if wmsg.msg != nil {
				extend:=wmsg.extend
				extend.ServerId = sID
				s.net.SendMsg(sinfo.fMasterID,wmsg.msg,&extend)
			} else {
				extend:=wmsg.extend
				extend.ServerId = sID
				s.net.SendBytes(sinfo.fMasterID, wmsg.msgID, wmsg.bytes,&extend)
			}
		}
	}
}

func (s *Slave) onPublishLoadLevel(ID uint32, serverType uint32,_ uint32,_ []byte, msg interface{},extend netframe.Server_Extend) {
	if serverType != s.masterType {
		return
	}

	req := msg.(*Master_PublishLoadLevel)

	if mv, ok := s.id2Masters[ID]; ok {
		mv.loadLv = req.LoadLevel
		if s.PrintLoadLevelStatus {
			logger.WithFields(logrus.Fields{
				"master":    ID,
				"loadLevel": req.LoadLevel,
			}).Info("OnPublishLoadLevel")
		}
	}
}

//检查动态切换master
func (s *Slave) checkShiftMaster(mID uint32, sID uint32) {
	//新master
	minfo, mok := s.id2Masters[mID]
	if !mok {
		return
	}
	//目标slave
	sinfo, sok := s.id2Subscribers[sID]
	if !sok {
		return
	}
	if !sinfo.isInitOk {
		return
	}
	if sinfo.fMasterID == mID {
		return
	}
	//fix master
	fminfo, have := s.id2Masters[sinfo.fMasterID]
	if !have {
		sinfo.fMasterID = 0
		s.fixMaster2SubServer(sID, 5)
		return
	}
	// 检查load值
	if fminfo.loadLv < minfo.loadLv+fminfo.config.Shiftload {
		return
	}

	// 随机时间后切换
	wtime := time.Second * time.Duration(rand.Intn(shiftMasterStartTimeMax-shiftMasterStartTimeMin)+shiftMasterStartTimeMin)
	//log.BDebug("checkShiftMaster-----------------------1  %d",sID)
	s.NetIO.AfterPost(wtime, func() {
		s.svrmutex.Lock()
		defer s.svrmutex.Unlock()
		if sinfo, ok := s.id2Subscribers[sID]; ok {
			var tminfo *ConnMasterInfo
			fminfo, fhave := s.id2Masters[sinfo.fMasterID]
			if !fhave {
				sinfo.fMasterID = 0
				s.fixMaster2SubServer(sID, 6)
				return
			}

			for mid, mv := range sinfo.mid2Status {
				//if !mv.isConnected || !mv.isWorking {
				if !mv.isConnected {
					continue
				}

				if minfo, have := s.id2Masters[mid]; have {
					if tminfo == nil || minfo.loadLv < tminfo.loadLv {
						tminfo = minfo
					}
				}
			}

			if tminfo != nil && fminfo != nil {
				//log.BDebug("checkShiftMaster-----------------------2  %d",sID)
				if fminfo.loadLv > tminfo.loadLv+fminfo.config.Shiftload {
					sinfo.shiftingMaster = true
					sinfo.shiftMsgQueue = newMsgQueue(shiftMsgQueueCap)
					smsg := &Slave_ReqShiftFixMaster{
						MasterID:  tminfo.config.ServerID,
						AServerID: s.SlaveConfig.ServerID,
						BServerID: sID,
					}
					s.net.SendMsg(sinfo.fMasterID, smsg,&netframe.Server_Extend{ServerId: sID})
					//超时
					s.NetIO.AfterPost(shiftMasterWaitTime, func() {
						s.svrmutex.Lock()
						defer s.svrmutex.Unlock()
						if sinfo, ok := s.id2Subscribers[sID]; ok {
							if sinfo.shiftingMaster {
								s.sendMsgQueue(sID)
								//log.BDebug("checkShiftMaster 请求切换超时 A:%d,M:%d,B:%d",s.SlaveConfig.ServerID,tminfo.config.ServerID,sID)
							}
						}
					})
					//log.BDebug("checkShiftMaster 发送切换请求 A:%d,M:%d,B:%d",s.SlaveConfig.ServerID,tminfo.config.ServerID,sID)
				}
			}
		}
	})
}

func (s *Slave) onShiftMasterReq(ID uint32, serverType uint32,_ uint32,_ []byte,msg interface{},extend netframe.Server_Extend) {
	req := msg.(*Slave_ReqShiftFixMaster)
	repmsg := &Slave_RepShiftFixMaster{
		MasterID:  req.MasterID,
		AServerID: req.AServerID,
		BServerID: req.BServerID,
	}

	if s.net.SendMsg(req.MasterID,repmsg,&extend) != nil {
		logger.WithFields(logrus.Fields{
			"asvrid": req.AServerID,
			"bsvrid": req.BServerID,
			"master": req.MasterID,
		}).Warning("OnShiftMasterReq req shift master failed!")
	}
}

func (s *Slave) onShiftMasterRep(ID uint32, serverType uint32,_ uint32,_ []byte, msg interface{},extend netframe.Server_Extend) {
	rep := msg.(*Slave_RepShiftFixMaster)

	s.svrmutex.Lock()
	defer s.svrmutex.Unlock()
	serverId:=extend.ServerId
	if ssinfo, ok := s.id2Subscribers[serverId]; ok {
		if ssinfo.shiftingMaster {
			//同步
			msg := &Slave_ReportLoadLevel{
				IsFix:     false,
				AServerID: s.SlaveConfig.ServerID,
				BServerID: serverId,
			}
			s.net.SendMsg(ssinfo.fMasterID, msg,&netframe.Server_Extend{ServerId: ssinfo.fMasterID,})

			ssinfo.fMasterID = 0
			s.fixMaster2SubServer(serverId, 7)
			//log.BDebug("OnShiftMasterRep A:%d req shift master:%d to B%d successed!", rep.AServerID, rep.MasterID, rep.BServerID)
			return
		}
	}

	logger.WithFields(logrus.Fields{
		"asvrid": rep.AServerID,
		"bsvrid": rep.BServerID,
		"master": rep.MasterID,
	}).Debug("OnShiftMasterRep req shift master failed!")
}

func (s *Slave) addNewSalveCfg(newID uint32, newType uint32) bool {

	//ClusterConfig *ClusterConf
	//SlaveConfig   *ServerConf
	//// 我关注的服务 type->group
	//mySubscribers  map[uint32]*SubServerGroup
	//// 所有关注服务用id索引
	//id2Subscribers map[uint32]*SubServerInfo

	if newID == s.SlaveConfig.ServerID {
		return false
	}

	if _, ok := s.id2Subscribers[newID]; ok {
		return true
	}

	for _, sconf := range s.ClusterConfig.Slaves {
		if sconf.ServerID == newID && sconf.ServerType == newType {
			var group *SubServerGroup
			if g, ok := s.mySubscribers[sconf.ServerType]; ok {
				group = g
			} else {
				group := &SubServerGroup{}

				group.serverType = sconf.ServerType
				group.id2SubServers = make(map[uint32]*SubServerInfo)

				s.mySubscribers[sconf.ServerType] = group
			}

			sitem := &SubServerInfo{config: sconf, fMasterID: 0, isInitOk: false}
			sitem.mid2Status = make(map[uint32]*RouterStatus)
			for _, mid := range sconf.MasterIDs {
				if _, ok := s.id2Masters[mid]; ok {
					sitem.mid2Status[mid] = &RouterStatus{masterID: mid}
				}
			}
			group.id2SubServers[sitem.config.ServerID] = sitem
			s.id2Subscribers[sitem.config.ServerID] = sitem

			return true
		}
	}

	return false
}

func (s *Slave) canLoadNewConfig(newConfig *ClusterConf) bool {
	s.svrmutex.RLock()
	defer s.svrmutex.RUnlock()
	if newConfig == nil {
		return false
	}

	if !newConfig.IsSameSlaveCfg(s.SlaveConfig) {
		return false
	}

	for _, mitem := range s.id2Masters {
		if !newConfig.IsSameMasterCfg(mitem.config) {
			return false
		}
	}

	for _, sitem := range s.id2Subscribers {
		if !newConfig.IsSameSlaveCfg(sitem.config) {
			return false
		}
	}

	return true
}

// OnSlaveReqPreCloseServer 准备关闭服务
func (s *Slave) OnSlaveReqPreCloseServer(ID uint32, serverType uint32,_ uint32,_ []byte,msg interface{},extend netframe.Server_Extend) {
	//req := msg.(*SS_CmdPrepareCloseServer)

	for _, mitem := range s.id2Masters {
		s.net.SendMsg(mitem.config.ServerID,&SM_ReqPrepareCloseMyself{},&netframe.Server_Extend{
			ServerId: mitem.config.ServerID,
		})
	}
}
