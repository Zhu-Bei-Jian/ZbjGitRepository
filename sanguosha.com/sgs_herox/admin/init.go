package admin

import (
	"database/sql"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/framework/netcluster"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameshared/accountservice"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/entityservice"
	"sanguosha.com/sgs_herox/gameshared/manager"
	"sanguosha.com/sgs_herox/gameshared/manager/session"
	_ "sanguosha.com/sgs_herox/gameshared/manager/session/providers/memory"
	"sanguosha.com/sgs_herox/gameshared/serverstateservice"
	"sanguosha.com/sgs_herox/proto/smsg"
	"strings"
	"time"
)

var app *appframe.Application
var dbMgr *database
var gmMgr *gmManager
var entityMgr entityservice.EntityService
var server *httpServer
var cacheAppErr chan *nats.Msg

var sessionMgr *session.Manager
var accountService *accountservice.Service
var SessionRoot string
var WebRoot string
var AreaId int32

var serverStateService serverstateservice.ServerStateService
var appConfig *config.AppConfig
var reporterMgr *ReporterManager

var slaveDB *sql.DB
var msgMonitor *MsgMonitor

var Develop bool

func InitAdmin(application *appframe.Application, cfgFile string, listenAddress string,webRoot string) error {
	//远程获取pprof数据
	util.SafeGo(func() {
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 20000+application.ID()), nil)
	})

	WebRoot = webRoot
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return err
	}
	appConfig = cfg
	AreaId = cfg.GameArea
	SessionRoot = cfg.SessionRoot
	Develop = cfg.Develop
	app = application
	//WebRoot = "/home/yoka/goProjects/sanguosha.com/sgs_herox"
	gmMgr = NewGMManager()
	dbMgr, err = openDB(cfg.DBGame)
	if err != nil {
		logrus.WithField("source", cfg.DBGame).WithError(err).Error("Connect to DB failed")
		return err
	}
	InitActionVerify(app)

	//dbMgr.CheckServerInfoTable()
	//dbMgr.CheckAuthAdminTable()

	//TODO slave为从库，暂时用主库，上线时用从库
	slaveDB = dbMgr.db

	serverStateService = serverstateservice.Init(app)
	err = CheckServerInfoInit()
	if err != nil {
		logrus.WithError(err).Error("CheckServerInfoInit")
		return err
	}

	cacheAppErr = make(chan *nats.Msg, 102400)

	app.RegisterService(sgs_herox.SvrTypeGate, appframe.WithLoadBalanceRandom(app, sgs_herox.SvrTypeGate))
	app.RegisterService(sgs_herox.SvrTypeAuth, appframe.WithLoadBalanceRandom(app, sgs_herox.SvrTypeAuth))
	app.RegisterService(sgs_herox.SvrTypeLobby, appframe.WithLoadBalanceSingleton(app, sgs_herox.SvrTypeLobby))
	app.RegisterService(sgs_herox.SvrTypeEntity, appframe.WithLoadBalanceRandom(app, sgs_herox.SvrTypeEntity))
	app.RegisterService(sgs_herox.SvrTypeGame, appframe.WithLoadBalanceRandom(app, sgs_herox.SvrTypeGame))

	entityMgr = entityservice.NewEntityService(app, cfg)
	logrus.WithField("source", cfg.DBGame).Info("Connect to DB succ")
	accountService = accountservice.New(app)

	server = &httpServer{}
	//server.authKey = cfg.Admin.AuthKey
	server.init(listenAddress)

	var dingdingTokens []string
	if m, err := dbMgr.GetServerInfo(); err == nil {
		if len(strings.TrimSpace(m["DingDingTokens"])) > 0 {
			dingdingTokens = strings.Split(m["DingDingTokens"], ";")
		}
	}
	reporterMgr, err = newReporterManager(dingdingTokens)
	if err != nil {
		return err
	}

	sessionMgr, err = session.NewManager("memory", "gosessionid", 3600)
	if err != nil {
		logrus.WithError(err).Error("sessionMgr.NewManager error")
		return err
	}
	sessionMgr.RunGC()

	reportMgr, err := manager.NewReportManager(cfg, app)
	if err != nil {
		return err
	}
	app.OnFiniHandler(func() {
		reportMgr.Close()
	})

	msgMonitor, err = NewMsgMonitor(cfg.MQNodes)
	if err != nil {
		return err
	}

	c := manager.NewCronMgr()
	c.AddAppStatusPushJob(app, reportMgr)
	c.Start()
	app.OnFiniHandler(func() {
		c.Stop()
	})

	gameshared.RegisterCommonCommand(app)
	gameshared.RegisterCommonServerStatus(app, reportMgr.PushServerStatus)

	//startRecordMetrics()
	listenMsg(app)

	evName := make(map[netcluster.SvrEvent]string)
	evName[netcluster.SvrEventStart] = "SvrEventStart"
	evName[netcluster.SvrEventReconnect] = "SvrEventReconnect"
	evName[netcluster.SvrEventQuit] = "SvrEventQuit"
	evName[netcluster.SvrEventDisconnect] = "SvrEventDisconnect"
	CommonServerEvent := func(svrType appframe.ServerType, svrID uint32, event netcluster.SvrEvent) {
		switch event {
		case netcluster.SvrEventStart, netcluster.SvrEventReconnect:
			CheckServerInfoOnConnected(svrID, svrType)
		}

		s, _ := evName[event]
		reporterMgr.Send(fmt.Sprintf("服务器【%v】%v %v", AreaId, svrID, s))
	}

	for svrType := sgs_herox.SvrTypeGate; svrType < sgs_herox.SvrTypeEnd; svrType++ {
		if svrType == sgs_herox.SvrTypeAdmin {
			continue
		}
		curType := svrType
		app.ListenServerEvent(curType, func(svrID uint32, event netcluster.SvrEvent) {
			CommonServerEvent(curType, svrID, event)
		})
	}

	//注册关闭调用
	app.OnExitHandler(Close)
	app.OnFiniHandler(Finish)

	util.SafeGo(func() {
		server.run()
	})

	return nil
}

func Close() {
	close(cacheAppErr)
}
func Finish() {
	time.Sleep(5 * time.Second)
	server.close()
	dbMgr.close()
}

func listenMsg(app *appframe.Application) {
	app.RegisterResponse((*smsg.AdEnRespChangeUserData)(nil))

	app.RegisterResponse((*smsg.AdAllRespGMCommand)(nil))
	app.RegisterResponse((*smsg.RespUserSyncDB)(nil))
	app.RegisterResponse((*smsg.RespQueryBackup)(nil))

	app.RegisterResponse((*smsg.AdAllRespReload)(nil))
	app.RegisterResponse((*smsg.AdAllRespCloseServer)(nil))
	app.RegisterResponse((*smsg.AdAllRespPingServer)(nil))

	app.RegisterResponse((*smsg.AllEnRespUserInfo)(nil))
	app.RegisterResponse((*smsg.AdAllRespMetrics)(nil))
	app.RegisterResponse((*smsg.AsAllRespMSGMonitor)(nil))
}
