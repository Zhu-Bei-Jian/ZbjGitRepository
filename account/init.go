package account

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/entityservice"
	"sanguosha.com/sgs_herox/gameshared/manager"
)

var app *appframe.Application

//var configMgr *ConfigMgr
var entityMgr entityservice.EntityService
var userCacheMgr *UserCacheManager
var dbDao *database

//目前仅提供玩家信息查询服务
func InitSvr(application *appframe.Application, cfgFile string) error {
	//远程获取pprof数据
	util.SafeGo(func() {
		http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 20000+application.ID()), nil)
	})

	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return err
	}

	app = application
	entityMgr = entityservice.NewEntityService(app, cfg)

	//configMgr = NewConfigMgr()
	//confManager := manager.NewConfManager(app, cfg.GameCfgPath, cfg.Develop, func(conf *conf.GameConfig) {
	//	configMgr.SetConfig(conf)
	//})
	//confManager.LoadConf()

	dbDao, err = newDatabase(cfg.DBGame)
	if err != nil {
		logrus.WithField("source", cfg.DBGame).WithError(err).Error("Connect to DB failed")
		return err
	}

	userCacheMgr, err = newUserCacheManager(cfg.RedisNodes[config.RedisUserCache])
	if err != nil {
		return fmt.Errorf("newUserCacheManager %w", err)
	}

	reportMgr, err := manager.NewReportManager(cfg, app)
	if err != nil {
		return err
	}

	app.OnFiniHandler(func() {
		reportMgr.Close()
	})

	msgWorker := appframe.NewParallelWorker(0, true)
	app.SetMsgHandlerWorker(msgWorker)

	c := manager.NewCronMgr()
	c.AddAppStatusPushJob(app, reportMgr)
	c.Start()
	app.OnFiniHandler(func() {
		c.Stop()
	})

	gameshared.RegisterCommonCommand(app)
	gameshared.RegisterCommonServerStatus(app, reportMgr.PushServerStatus)

	initMsgHandler(app)
	return nil
}
