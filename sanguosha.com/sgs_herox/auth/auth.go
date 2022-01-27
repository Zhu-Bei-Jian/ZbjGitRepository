package auth

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"sanguosha.com/baselib/appframe"
	wordsfilter "sanguosha.com/baselib/sgslib/filter"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	"sanguosha.com/sgs_herox/gameshared/entityservice"
	"sanguosha.com/sgs_herox/gameshared/logproducer"
	"sanguosha.com/sgs_herox/gameshared/manager"
	"sanguosha.com/sgs_herox/gameshared/serverstateservice"

	_ "net/http/pprof"
)

var wordsFilter wordsfilter.Filter

var authApp *appframe.Application
var logMgr *logproducer.LogManager
var workerMgr *manager.WorkerManager

var gameConfig *conf.GameConfig
var entityMgr entityservice.EntityService
var dbDao *database

var serverStateService serverstateservice.ServerStateService
var serverInfos *gameshared.ServerInfos

var ServerOpenTime int64
var ServerCloseTime int64
var InnerIp []string
var ApkVersion string

var tableParkCfg config.TablePark

// InitAuthSvr 初始化 authsvr 服务
func InitAuthSvr(app *appframe.Application, cfgFile string) error {
	//远程获取pprof数据
	util.SafeGo(func() {
		err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 20000+app.ID()), nil)
		if err != nil {
			logrus.WithError(err).Error("pprof listen error")
		}
	})

	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return err
	}

	tableParkCfg = cfg.TablePark

	dbDao, err = NewDatabase(cfg.DBGame)
	if err != nil {
		return fmt.Errorf("NewDataBase error:%w", err)
	}

	app.SetMsgHandlerWorker(appframe.NewParallelWorker(0, true))
	workerMgr = &manager.WorkerManager{}
	workerMgr.Init(app)

	reportMgr, err := manager.NewReportManager(cfg, app)
	if err != nil {
		return err
	}

	app.OnFiniHandler(func() {
		reportMgr.Close()
	})

	//服务器状态信息
	serverStateService = serverstateservice.Init(app)
	serverInfos = serverStateService.GetServerInfos()
	serverStateService.OnServerInfoUpdate(onServerInfoUpdate)

	logMgr, err = logproducer.New(cfg, app.ID())
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

	initAuthMsgHandler(app)
	app.OnFiniHandler(Finish)

	workerMgr.Run()
	return nil
}

func Finish() {
	workerMgr.Close()
}
