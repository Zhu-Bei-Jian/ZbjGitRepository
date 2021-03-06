package appframe

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"sync/atomic"

	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe/request/protoreq"
	"sanguosha.com/baselib/framework/netcluster"
	"sanguosha.com/baselib/ioservice"
	"sanguosha.com/baselib/log"
	"sanguosha.com/baselib/network/msgprocessor"
	"sanguosha.com/baselib/framework/netframe"
)

// DisableApplicationInitGlobalLogrus 创建 Application 时, 禁止初始化全局的 logrus 设置
// 在有多个 Application 实例时, 该选项应该设置为 true. (例如 All in One 模式)
var DisableApplicationInitGlobalLogrus bool

// Application 应用程序类.
type Application struct {
	name             string
	id               uint32
	worker           ioservice.IOService
	slave            *netcluster.Slave
	exitCh           chan int
	services         map[ServerType]Service
	reqc             *protoreq.Client
	respMsgMap       map[reflect.Type]bool
	msgWorker        Worker
	onRun            []func()
	onExit           []func()
	onFinis          []func()
	exit             int32
	closeLogger      func()
	userIDHandlerMap map[reflect.Type]interface{}
}

// NewApplication 创建 Application
// name 参数为 netconfigFile 中对应的 name 值
func NewApplication(netconfigFile string, name string) (*Application, error) {
	netconfig, err := netcluster.ParseClusterConfigFile(netconfigFile)
	if err != nil {
		return nil, err
	}
	a := new(Application)
	err = a.init(netconfig, name)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Application) init(netconfig *netcluster.ClusterConf, name string) error {
	cfg, ok := netconfig.Slaves[name]
	if !ok {
		return fmt.Errorf("can not find app config by name %s", name)
	}

	a.name = name
	a.id = cfg.ServerID
	a.services = make(map[ServerType]Service)

	a.worker = ioservice.NewIOService(fmt.Sprintf("app_%s_main", name), 102400)
	a.worker.Init()
	a.slave = netcluster.NewSlave(netconfig, name, a.worker)
	if a.slave == nil {
		return errors.New("new slave")
	}
	a.slave.Init()
	a.exitCh = make(chan int, 1)

	a.reqc = protoreq.NewClient(a.Post)
	a.reqc.OnNotFind = func(seqid int64, resp interface{}, err error) {
		log := logrus.WithField("seqid", seqid)
		if err != nil {
			log.WithError(err).Warn("Resp not find ", a.id)
		} else {
			log.WithField("respType", reflect.TypeOf(resp)).Warn("Resp not find ", a.id)
		}
	}
	a.respMsgMap = make(map[reflect.Type]bool)
	a.userIDHandlerMap = make(map[reflect.Type]interface{})

	a.RegisterResponse((*protoreq.ErrCode)(nil))

	if !DisableApplicationInitGlobalLogrus {
		close, err := log.InitLogrus(&cfg.Log)
		if err != nil {
			return err
		}
		a.closeLogger = close
	}
	return nil
}

// Name 名称
func (a *Application) Name() string {
	return a.name
}

// ID 唯一标识.
func (a *Application) ID() uint32 {
	return a.id
}

// ServerAddr 服务器IP.
func (a *Application) ServerAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		logrus.WithError(err).Error("Application.ServerAddr InterfaceAddrs error")
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		ipNet, ok := address.(*net.IPNet)
		if !ok {
			continue
		}
		if ipNet.IP.IsLoopback() {
			continue
		}
		if ipNet.IP.To4() == nil {
			continue
		}
		return ipNet.IP.String()
	}

	return ""
}

//主协程工作长度
func (a *Application) WorkerLen() int {
	return a.worker.WorkerLen()
}

// Run 启动 Application, 阻塞调用.
// 退出使用信号 os.Interrupt 或 os.Kill.
func (a *Application) Run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	a.worker.Run()
	a.slave.Run()

	// 在 app 的 worker 中调用这些启动时需要触发的回调函数.
	onruns := a.onRun
	a.worker.Post(func() {
		for _, f := range onruns {
			f()
		}
	})
	a.onRun = nil

	log := logrus.WithField("id", a.ID())
	log.Infof("App (%s) running...", a.name)

	select {
	case sig := <-c:
		log.Infof("App (%s) exiting... signal:(%v)", a.name, sig)
	case <-a.exitCh:
		log.Infof("App (%s) exiting...", a.name)
	}

	atomic.StoreInt32(&a.exit, 1)

	// 通知业务程序将要退出
	onexits := a.onExit

	var wg sync.WaitGroup
	wg.Add(1)
	a.worker.Post(func() {
		defer wg.Done()
		for _, f := range onexits {
			f()
		}
	})
	a.onExit = nil

	//等待close执行完毕
	wg.Wait()

	// 等待所有消息处理完毕
	if a.msgWorker != nil {
		a.msgWorker.WaitDone()
	}

	// TODO 这可能造成程序退出时间过长, 但如果使用超时机制同样也比较尴尬, 见 WaitAllDone 方法内的注释.
	a.reqc.WaitAllDone()

	// 最后的收尾工作, 倒序调用.
	for i := len(a.onFinis) - 1; i >= 0; i-- {
		a.onFinis[i]()
	}
	a.onFinis = nil

	log.Infof("App Exit (%s)", a.name)

	a.slave.Fini()
	a.worker.Fini()

	if a.closeLogger != nil {
		a.closeLogger()
	}
}

// Exit 发送退出指令, 用于主动结束程序.
func (a *Application) Exit() {
	a.exitCh <- 0
}

// IsExit 正在退出或已退出.
func (a *Application) IsExit() bool {
	return atomic.LoadInt32(&a.exit) == 1
}

// OnRunHandler 在 app 启动时触发, 在 Run 之前设置, 可设置多个, 按顺序调用, 在 app 的 worker 中执行.
func (a *Application) OnRunHandler(f func()) {
	a.onRun = append(a.onRun, f)
}

// OnExitHandler 在 app 将要退出时触发, 时间点在 OnFiniHandler 之前.
// 在 Run 之前设置, 可设置多个, 按顺序调用, 在 app 的 worker 中执行.
func (a *Application) OnExitHandler(f func()) {
	a.onExit = append(a.onExit, f)
}

// OnFiniHandler app 退出时的最后处理函数, 在所有消息都已处理完毕后触发.
// 在 Run 之前设置, 可设置多个, 倒序调用.
func (a *Application) OnFiniHandler(f func()) {
	a.onFinis = append(a.onFinis, f)
}

// Post 将工作任务派发到 Application 的 worker 中执行.
func (a *Application) Post(f func()) {
	a.worker.Post(f)
}

// AfterPost 将工作任务派发到 Application 的 worker 中执行.
func (a *Application) AfterPost(d time.Duration, f func()) (cancel func() bool) {
	return a.worker.AfterPost(d, f)
}

// ListenMsg 监听来自服务的消息
// 如果没有调用 SetMsgHandlerWorker, 则回调函数将在 app 的 worker 中执行.
func (a *Application) ListenMsg(msg proto.Message, handler MsgHandler) {
	isResponseMessage := a.isResponseMessage(msg)
	a.slave.ListenServerMessage(0, msg, func(_ uint32, _ uint32,_ uint32,_ []byte,imsg interface{},extend netframe.Server_Extend) {
		// 在 app 退出的过程中, 非响应类型消息不再接受处理.
		if a.IsExit() {
			if !isResponseMessage {
				return
			}
		}
		msg, ok := imsg.(proto.Message)
		if !ok {
			logrus.Error("Msg from server is not a proto.Message")
			return
		}
		if a.msgWorker != nil { // 在业务自定义的消息函数执行者中执行回调函数.
			a.msgWorker.Post(func() {
				handler(a.GetServer(extend.ServerId),extend.SeqId, msg)
			})
		} else { // 默认在 app 的 worker 中执行回调函数.
			handler(a.GetServer(extend.ServerId),extend.SeqId, msg)
		}
	})
}

// ListenSessionMsg 监听来自用户会话的消息
// 如果没有调用 SetMsgHandlerWorker, 则回调函数将在 app 的 worker 中执行.
func (a *Application) ListenSessionMsg(msg proto.Message, handler SessionMsgHandler) {
	isResponseMessage := a.isResponseMessage(msg)
	a.slave.ListenServerMessage(0, msg, func(_ uint32, _ uint32,_ uint32,_ []byte,imsg interface{},extend netframe.Server_Extend) {
		// 在 app 退出的过程中, 非响应类型消息不再接受处理.
		if a.IsExit() {
			if !isResponseMessage {
				return
			}
		}
		msg, ok := imsg.(proto.Message)
		if !ok {
			logrus.Error("Msg from server is not a proto.Message")
			return
		}
		sid:=SessionID{
			SvrID: extend.ServerId,
			ID:    extend.SessionId,
		}
		sidExt := SessionIDExtend{SessionID:sid,UserID: extend.UserId,ExtParam: extend.ExtParam}

		if a.msgWorker != nil { // 在业务自定义的消息函数执行者中执行回调函数.
			a.msgWorker.Post(func() {
				handler(a.GetSessionExtend(sidExt), msg)
			})
		} else { // 默认在 app 的 worker 中执行回调函数.
			handler(a.GetSessionExtend(sidExt), msg)
		}
	})
}

// SetMsgHandlerWorker 设置如何执行 ListenMsg 时设置的回调函数, 这是一个可选项.
// 默认所有监听消息的回调函数是在 app 的 worker 中执行的, 这里提供一个选项, 让业务可以自己定义如何执行这些消息的回调函数.
// 例如, 在无状态的应用中, 可以设置为每一个消息在独立的 goroutine 中执行回调函数, 然后业务层就可以使用同步阻塞风格的方式编程了.
func (a *Application) SetMsgHandlerWorker(worker Worker) {
	a.msgWorker = worker
}

func (a *Application) MsgWorkerGCount() int32 {
	if a.msgWorker == nil {
		return 0
	}

	return a.msgWorker.GroutineCount()
}

// RegisterResponse 注册响应消息, 响应类型消息中, 必须有成员字段 Seqid int64.
// 在 app 退出的过程中, 响应类型的消息将能够被正常处理, 非响应消息将会被拒绝处理.
func (a *Application) RegisterResponse(msg proto.Message) {
	//if !protoreq.IsValidMsg(msg) {
	//	logrus.Panic("RegisterResponse not a valid resp msg type")
	//}
	typ := reflect.TypeOf(msg)
	if !a.respMsgMap[typ] {
		a.respMsgMap[typ] = true
		a.ListenMsg(msg, func(sender Server,seqId int64, msg proto.Message) {
			a.reqc.OnResp(msg,seqId)
		})
	}
}

func (a *Application) isResponseMessage(msg proto.Message) bool {
	return a.respMsgMap[reflect.TypeOf(msg)]
}

// DebugPrintMessage 是否打印收到和发送的消息.
func (a *Application) DebugPrintMessage(able bool) {
	a.slave.DebugPrintMessage = able
}

//
//
// 服务相关接口.

// GetServer 获取指定的服务节点.
// 使用者不用对返回结果判空, 对无效 server 的错误处理, 都会延迟到调用 Server 的方法时发生, 这样做的目的就是为了简化错误处理.
func (a *Application) GetServer(id uint32) Server {
	return &server{
		id:  id,
		app: a,
	}
}

// GetAvailableServerIDs 根据类型获取所有可用的服务节点id, 断线的会被排除.
func (a *Application) GetAvailableServerIDs(typ ServerType) []uint32 {
	return a.slave.GetServerAllAvailable(uint32(typ))
}

// ListenServerEvent 监听服务事件, 回调函数将在 app 的 worker 中执行.
func (a *Application) ListenServerEvent(typ ServerType, handler netcluster.SvrEventHandler) {
	a.slave.ListenServerStatus(uint32(typ), handler)
}

// RegisterService 注册一类型的服务, 提供负载均衡策略.
func (a *Application) RegisterService(typ ServerType, loadBalance LoadBalance) Service {
	if _, exist := a.services[typ]; exist {
		// 如果已存在, 提示警告, 并进行替换.
		logrus.WithField("svrtype", typ).Warn("Service already register")
	}
	s := &commonService{
		app:        a,
		typ:        typ,
		loadBlance: loadBalance,
	}
	a.services[typ] = s
	a.ListenServerEvent(typ, func(svrid uint32, event netcluster.SvrEvent) {
		loadBalance.OnServerEvent(svrid, event)
	})
	return s
}

// GetService 获得一种服务, 只有注册过的 RegisterService 的服务才能够通过 GetService 接口获取, 否则将得到一个 unregisteredService 的虚拟服务.
func (a *Application) GetService(typ ServerType) Service {
	result, ok := a.services[typ]
	if !ok {
		return unregisteredService{app: a}
	}
	return result
}

// GetSession 获取用户会话对象
func (a *Application) GetSession(id SessionID) Session {
	return &svrSession{
		id:  id,
		app: a,
	}
}

func(a *Application)GetSessionExtend(id SessionIDExtend)Session{
	return &svrSessionExtend{
		id:  id,
		app: a,
	}
}

func (a *Application) HandleMsgWithUserID(sender Session, msgID uint32, userID uint64, data []byte) {
	if msg, err := msgprocessor.OnUnmarshal(msgID, data); err == nil {
		if handler, exist := a.userIDHandlerMap[reflect.TypeOf(msg)]; exist {
			v := reflect.ValueOf(handler)
			v.Call([]reflect.Value{reflect.ValueOf(sender), reflect.ValueOf(userID), reflect.ValueOf(msg)})
		} else {
			logrus.Error("HandleMsgWithUserID req.MsgID not exist error")
		}
	}
}

//只有注册了上面的HandleMsgWithUserID，这个Listen才有用
func (a *Application) ListenMsgWithUserID(handler interface{}) {
	v := reflect.ValueOf(handler)
	if v.Type().NumIn() != 3 {
		logrus.Panic("ListenWithUserIDMsg handler params num wrong")
	}
	if v.Type().In(0) != reflect.TypeOf((*Session)(nil)).Elem() {
		logrus.Panic("ListenWithUserIDMsg handler num in 0 is not appframe.Session")
	}
	if v.Type().In(1).Kind() != reflect.Uint64 {
		logrus.Panic("ListenWithUserIDMsg handler num in 2 is not uint64")
	}
	msg := reflect.New(v.Type().In(2)).Elem().Interface().(proto.Message)

	msgTyp := reflect.TypeOf(msg)
	//id := util.StringHash(proto.MessageName(msg))
	a.userIDHandlerMap[msgTyp] = handler
	msgprocessor.RegisterMessage(msg)
}

func(a *Application) RegisterIntercepter(f func(msgSrc netcluster.MsgSrc,connId uint32,msgId uint32,msgData []byte,extend *netframe.Server_Extend)){
	a.slave.RegisterIntercepter(f)
}
