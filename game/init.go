package game

import (
	"fmt"
	"math/rand"
	"sanguosha.com/sgs_herox/gameshared/logproducer"
	"sanguosha.com/sgs_herox/gameshared/notifier"
	"sanguosha.com/sgs_herox/proto/smsg"
	"time"

	"github.com/sirupsen/logrus"

	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/framework/netcluster"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	"sanguosha.com/sgs_herox/gameshared/entityservice"
	"sanguosha.com/sgs_herox/gameshared/manager"
)

var App *appframe.Application

var AppCfg *config.AppConfig

var workerMgr *manager.WorkerManager

var lastTick = time.Now()

var logMgr *logproducer.LogManager
var AreaID int32
var gWorkerMgr *manager.WorkerManager
var EntityService entityservice.EntityService

var gameMgr *GameManager
var playerMgr *PlayerManager

var idMgr *manager.IDManager
var wordFilterMgr *manager.WordFilterManager

var gCfg *conf.GameConfig
var tpNotifier *notifier.TableParkNotifier

var Develop bool

// InitGameSvr 初始化 gamesvr
func InitGameSvr(app *appframe.Application, cfgFile string) error {
	rand.Seed(time.Now().UnixNano())
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return err
	}
	App = app
	AreaID = cfg.GameArea

	//gamelogic.GameManagerInstance = gamelogic.NewGameManager()

	confManager := manager.NewConfManager(app, cfg.GameCfgPath, cfg.Develop, func(conf *conf.GameConfig) {
		gCfg = conf
	})
	confManager.LoadConf()

	workerMgr = &manager.WorkerManager{}
	workerMgr.Init(app)

	logMgr, err = logproducer.New(cfg, app.ID())
	if err != nil {
		return err
	}

	EntityService = entityservice.NewEntityService(app, cfg)
	InitGameMsgHandler(app)
	initMetricsMsgHandler(app)
	Develop = cfg.Develop

	App = app
	gWorkerMgr = workerMgr

	idMgr = &manager.IDManager{}
	idMgr.Init(app.ID())

	tpNotifier = notifier.NewTableParkNotifier(config.GameID, cfg.TablePark)

	gameMgr = newGameManager()
	playerMgr = newPlayerManager()

	wordFilterMgr = manager.NewWordFilterManager(cfg.WordsFilter)
	app.OnFiniHandler(wordFilterMgr.Close)

	app.RegisterService(sgs_herox.SvrTypeLobby, appframe.WithLoadBalanceSingleton(app, sgs_herox.SvrTypeLobby))

	app.ListenServerEvent(sgs_herox.SvrTypeLobby, func(svrid uint32, e netcluster.SvrEvent) {
		switch e {
		case netcluster.SvrEventStart, netcluster.SvrEventReconnect:
			app.GetServer(svrid).SendMsg(&smsg.SyncServerVersion{Id: app.ID(), Ver: 0})
		case netcluster.SvrEventQuit, netcluster.SvrEventDisconnect:
			logrus.Error(fmt.Sprintf("SvrTypeLobby %d Disconnect or Quit", svrid))
			for _, g := range gameMgr.allGames() {
				if g.IsOver() {
					continue
				}
				g.ForceOver()
			}
		}
	})
	//ai服崩溃后，所有相关游戏桌子都异常结束掉--->改为仅取消ai关联
	app.ListenServerEvent(sgs_herox.SvrTypeAI, func(svrid uint32, e netcluster.SvrEvent) {
		switch e {
		case netcluster.SvrEventQuit, netcluster.SvrEventDisconnect:
			logrus.Error(fmt.Sprintf("SvrTypeAI %d Disconnect or Quit", svrid))
		}
	})

	reportMgr, err := manager.NewReportManager(cfg, app)
	if err != nil {
		return err
	}
	app.OnFiniHandler(func() {
		reportMgr.Close()
	})

	app.OnFiniHandler(func() {
		reportMgr.Close()
	})

	//启动定时任务
	c := manager.NewCronMgr()
	c.AddAppStatusPushJob(app, reportMgr)
	c.Start()
	app.OnFiniHandler(func() {
		c.Stop()
	})

	gameshared.RegisterCommonCommand(app)
	gameshared.RegisterCommonServerStatus(app, reportMgr.PushServerStatus)

	app.OnExitHandler(Close)
	app.OnFiniHandler(Finish)

	Run()

	return nil
}

func Run() {
	workerMgr.Run()
}

func Close() {

}

func Finish() {
	workerMgr.Close()
}
