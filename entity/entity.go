package entity

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"net/http"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/entity/instance"
	"sanguosha.com/sgs_herox/entity/userdb"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	"sanguosha.com/sgs_herox/gameshared/manager"
	"sanguosha.com/sgs_herox/gameshared/notifier"
	"sanguosha.com/sgs_herox/gameshared/serverstateservice"
	"sanguosha.com/sgs_herox/gameutil"
	"time"

	_ "net/http/pprof"
)

var UserDBInstance *userdb.DB
var onlineMgr *onlineManager
var AppCfg *config.AppConfig

var AppInstance *appframe.Application

var serverStateService serverstateservice.ServerStateService
var serverInfos *gameshared.ServerInfos
var useTokenTime int64

var workerMgr *manager.WorkerManager
var idMgr *manager.IDManager
var confManager *manager.ConfManager
var tpNotifier *notifier.TableParkNotifier

// InitEntitySvr 初始化
func InitEntitySvr(app *appframe.Application, cfgFile string) error {
	//远程获取pprof数据
	util.SafeGo(func() {
		err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 20000+app.ID()), nil)
		fmt.Println(err)
	})

	AppInstance = app

	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return err
	}

	idMgr = &manager.IDManager{}
	idMgr.Init(app.ID())

	AppCfg = cfg
	nodeCfg, ok := cfg.GetEntityNodeConfig(app.ID())
	if !ok {
		return fmt.Errorf("can not find entity node (%d) config", app.ID())
	}

	dbGameSource := cfg.DBGame
	logrus.WithField("source", dbGameSource).Info("Connect to DB ...")
	dbGame, err := gameutil.OpenDBWithMaxMinConn(dbGameSource, nodeCfg.ShardCnt, nodeCfg.ShardCnt/10)
	if err != nil {
		logrus.WithField("source", dbGameSource).WithError(err).Error("Connect to DB failed")
		return err
	}
	logrus.WithField("source", dbGameSource).Info("Connect to DB succ")
	app.OnFiniHandler(func() { dbGame.Close() })

	err = gameutil.CheckProcAddColumn(dbGame)
	if err != nil {
		logrus.WithField("source", dbGameSource).WithError(err).Error("CheckProcAddColumn")
		return err
	}

	onlineMgr = newUserOnlineManager()
	initUserMsgHandler(app)

	serverStateService = serverstateservice.Init(app)
	serverInfos = serverStateService.GetServerInfos()
	useTokenTime = gameshared.TokenTime
	serverStateService.OnServerInfoUpdate(func(si *gameshared.ServerInfos) {
		useTokenTime = si.GetInt64("TokenTime", gameshared.TokenTime)
	})

	msgToService := func(typ appframe.ServerType, msg proto.Message) {
		app.GetService(typ).SendMsg(msg)
	}

	msgRequest := func(typ appframe.ServerType, msg proto.Message, cbk func(resp proto.Message, err error), timeout time.Duration) {
		app.GetService(typ).Request(msg, cbk, timeout)
	}

	useridToShard := func(userid uint64) (int, error) {
		svrid, shardid := cfg.GetUserEntityIDAndShard(userid)
		if svrid != app.ID() {
			logrus.WithFields(logrus.Fields{
				"userid":        userid,
				"shouldToSvrid": svrid,
				"currentSvrid":  app.ID(),
			}).Error("Route to wrong entity node")
			return 0, errors.New("Route to wrong entity node")
		}
		return shardid, nil
	}

	UserDBInstance = userdb.NewDB(dbGame, nodeCfg.ShardCnt, nodeCfg.ChanLen, useridToShard, msgToService, msgRequest, AppInstance, cfg.Develop)
	app.OnFiniHandler(func() {
		logrus.Info("Sync data to DB...")
		err := UserDBInstance.Close()
		if err != nil {
			logrus.WithError(err).Error("Close userdb done with error")
		} else {
			logrus.Info("Close userdb done")
		}
	})
	UserDBInstance.AreaID = cfg.GameArea

	//加载配置
	confManager = manager.NewConfManager(app, cfg.GameCfgPath, cfg.Develop, func(conf *conf.GameConfig) {
		UserDBInstance.EveryShardUpdateConf(conf)
	})
	confManager.LoadConf()

	tpNotifier = notifier.NewTableParkNotifier(config.GameID, cfg.TablePark)

	initEntityMsgHandler(app)

	// 注册服务.
	app.RegisterService(sgs_herox.SvrTypeAuth, appframe.WithLoadBalanceRandom(app, sgs_herox.SvrTypeAuth))
	app.RegisterService(sgs_herox.SvrTypeLobby, appframe.WithLoadBalanceSingleton(app, sgs_herox.SvrTypeLobby))

	workerMgr = &manager.WorkerManager{}
	workerMgr.Init(app)

	err = instance.InitLogManagerInstance(app, cfg)
	if err != nil {
		return err
	}

	reportMgr, err := manager.NewReportManager(cfg, app)
	app.OnFiniHandler(func() {
		reportMgr.Close()
	})

	gameshared.RegisterCommonCommand(app)
	gameshared.RegisterCommonServerStatus(app, reportMgr.PushServerStatus)

	app.OnFiniHandler(Finish)

	c := manager.NewCronMgr()
	c.AddAppStatusPushJob(app, reportMgr)
	c.Add10SecJob(func() {
		//reportMgr.PushNatsServerLoad(config.EntityWorkerLen, int32(UserDBInstance.ShareWorkerLen()))
	})
	c.AddDayJob(AM0Fresh)
	c.AddMinuteJob(EveryMinuteFresh)
	c.Start()

	app.OnFiniHandler(func() {
		c.Stop()
	})

	Run()
	return nil
}

func Run() {
	workerMgr.Run()
}

func Finish() {
	workerMgr.Close()
}
