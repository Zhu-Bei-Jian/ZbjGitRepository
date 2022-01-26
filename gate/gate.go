package gate

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/framework/netcluster"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/entityservice"
	"sanguosha.com/sgs_herox/gameshared/manager"
)

var SessionMgrInstance *sessionManager
var wordFilterMgr *manager.WordFilterManager
var AppInstance *appframe.GateApplication
var AppCfg *config.AppConfig
var ChannelMgInstance *ChannelMg
var EntityInstance entityservice.EntityService
var userMsgMonitor *UserMsgMonitor

// InitGateSvr 初始化 gatesvr.
func InitGateSvr(app *appframe.GateApplication, cfgFile string) error {
	//远程获取pprof数据
	util.SafeGo(func() {
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 20000+app.ID()), nil)
	})

	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return err
	}
	AppCfg = cfg
	AppInstance = app

	EntityInstance = entityservice.NewEntityService(&app.Application, cfg)

	// 注册 auth 服务.
	app.RegisterService(sgs_herox.SvrTypeAuth, appframe.WithLoadBalanceRandom(app, sgs_herox.SvrTypeAuth))
	// 注册 lobby 服务.
	app.RegisterService(sgs_herox.SvrTypeLobby, appframe.WithLoadBalanceSingleton(app, sgs_herox.SvrTypeLobby))
	// 注册 gate 服务.
	app.RegisterService(sgs_herox.SvrTypeGate, appframe.WithLoadBalanceRandom(app, sgs_herox.SvrTypeGate))

	app.ListenServerEvent(sgs_herox.SvrTypeLobby, onLobbyServerDisconnect)

	SessionMgrInstance = initSessionManager(app)

	userMsgMonitor = &UserMsgMonitor{}
	userMsgMonitor.init()

	ChannelMgInstance = NewChannelMg()

	wordFilterMgr = manager.NewWordFilterManager(cfg.WordsFilter)
	app.OnFiniHandler(wordFilterMgr.Close)

	RegisterBroadcastMsg()
	initGateMsgRoute(app)
	initGateMsgHandler(app)
	initMetricsMsgHandler(app)

	reportMgr, err := manager.NewReportManager(cfg, &app.Application)
	if err != nil {
		return err
	}

	app.OnFiniHandler(func() {
		reportMgr.Close()
	})

	c := manager.NewCronMgr()
	c.AddAppStatusPushJob(&app.Application, reportMgr)
	c.AddMinuteJob(func() {
		app.Post(func() {
			num := int32(len(SessionMgrInstance.sid2session))
			reportMgr.PushGateOnline(num, AppInstance.ServerAddr())
			//reportMgr.PushNatsServerLoad(config.GateOnLineCount, num)
		})
	})
	c.Start()
	app.OnFiniHandler(func() {
		c.Stop()
	})

	gameshared.RegisterCommonCommand(&app.Application)
	gameshared.RegisterCommonServerStatus(&app.Application, reportMgr.PushServerStatus)

	err = RegisterMSGMonitorCommand(cfg, &app.Application, SessionMgrInstance)
	if err != nil {
		return err
	}

	app.OnExitHandler(Close)
	return nil
}

// 获取用户的 entity server, 根据 userid hash
func GetEntity(userid uint64) appframe.Server {
	return AppInstance.GetServer(AppCfg.GetUserEntityID(userid))
}

func Close() {
	SessionMgrInstance.close()
}

func onLobbyServerDisconnect(svrID uint32, event netcluster.SvrEvent) {
	switch event {
	case netcluster.SvrEventQuit, netcluster.SvrEventDisconnect:
		logrus.WithFields(logrus.Fields{
			"svrID": svrID,
			"event": event,
		}).Error("lobby server disconnect")

		SessionMgrInstance.execByEverySession(func(s *session) {
			s.Close()
		})
	}
}
