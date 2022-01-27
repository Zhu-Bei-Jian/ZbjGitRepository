package lobby

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	"sanguosha.com/sgs_herox/gameshared/entityservice"
	"sanguosha.com/sgs_herox/gameshared/logproducer"
	"sanguosha.com/sgs_herox/gameshared/manager"
	"sanguosha.com/sgs_herox/gameshared/serverstateservice"
	"sanguosha.com/sgs_herox/proto/smsg"
)

var userMgr *userManager

var EntityInstance entityservice.EntityService

var gameMgr *GameManager
var AppInstance *appframe.Application

var logMgr *logproducer.LogManager

var workerMgr *manager.WorkerManager

var AppCfg *config.AppConfig

var msgOpenInfo *smsg.SyncServerOpenInfo

var roomMgr *RoomManager
var serverStateService serverstateservice.ServerStateService

var matchMgr *MatchManager
var configGlobal *conf.GameConfig

// InitLobbySvr 初始化 lobbysvr.
func InitLobbySvr(app *appframe.Application, cfgFile string) error {
	//远程获取pprof数据
	util.SafeGo(func() {
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 20000+app.ID()), nil)
	})

	AppInstance = app

	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return err
	}
	AppCfg = cfg

	//初始化游戏配置文件读取
	confManager := manager.NewConfManager(app, cfg.GameCfgPath, cfg.Develop, func(conf *conf.GameConfig) {
		SetConfig(conf)
	})
	confManager.LoadConf()

	logMgr, err = logproducer.New(cfg, app.ID())
	if err != nil {
		return err
	}

	workerMgr = &manager.WorkerManager{}
	workerMgr.Init(app)

	roomMgr = NewRoomManager()

	reportMgr, err := manager.NewReportManager(cfg, app)
	if err != nil {
		return err
	}
	app.OnFiniHandler(func() {
		reportMgr.Close()
	})

	gameshared.RegisterCommonCommand(app)
	gameshared.RegisterCommonServerStatus(app, reportMgr.PushServerStatus)

	EntityInstance = entityservice.NewEntityService(app, cfg)

	userMgr = initUserManager(app, EntityInstance)

	initUserMsgHandler(app)

	serverStateService = serverstateservice.Init(app)
	serverStateService.WatchServer(sgs_herox.SvrTypeGame)

	gameMgr = newGameManager(app)
	initGameMsgHandler(app)
	initRoomMsgHandler(app)
	initChatMsgHandler(app)
	initMetricsMsgHandler(app)

	matchMgr = &MatchManager{}
	matchMgr.init()

	c := manager.NewCronMgr()
	c.AddAppStatusPushJob(app, reportMgr)
	c.AddMinuteJob(func() {
		app.Post(func() {
			count := len(userMgr.ss2user)
			reportMgr.PushOnline(int32(count))
			//reportMgr.PushNatsServerLoad(config.OnlineCount, int32(count))
		})
	})

	c.Start()

	app.OnFiniHandler(func() {
		c.Stop()
	})

	app.RegisterService(sgs_herox.SvrTypeAuth, appframe.WithLoadBalanceRandom(app, sgs_herox.SvrTypeAuth))
	app.RegisterService(sgs_herox.SvrTypeAI, appframe.WithLoadBalanceRandom(app, sgs_herox.SvrTypeAI))
	app.RegisterService(sgs_herox.SvrTypeGame, appframe.WithLoadBalanceRandom(app, sgs_herox.SvrTypeGame))

	ListenServerEvent(app)

	app.OnExitHandler(Close)
	app.OnFiniHandler(Finish)

	Run()

	return nil
}

func Run() {
	workerMgr.Run()
	matchMgr.run()
}

func Close() {

}

func Finish() {
	workerMgr.Close()
	logMgr.Close()
}
