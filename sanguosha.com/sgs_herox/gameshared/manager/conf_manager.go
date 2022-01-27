package manager

import (
	"fmt"
	"sanguosha.com/sgs_herox/proto/smsg"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	"sanguosha.com/sgs_herox/gameutil"
)

// ConfManager ...
type ConfManager struct {
	paths   *conf.GameConfigPathNode
	gameCfg *conf.GameConfig

	develop bool

	reloadTime   time.Time
	mutex        sync.RWMutex
	loadCallback func(config *conf.GameConfig)
}

func NewConfManager(app *appframe.Application, paths *conf.GameConfigPathNode, dev bool, callback func(config *conf.GameConfig)) *ConfManager {
	manager := &ConfManager{
		paths:        paths,
		loadCallback: callback,
		develop:      dev,
	}
	manager.initReload(app)
	return manager
}

func (p *ConfManager) LoadConf() *conf.GameConfig {
	p.gameCfg = &conf.GameConfig{}
	p.gameCfg.Develop = p.develop
	err := p.gameCfg.Init(p.paths)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"paths": p.paths,
		}).WithError(err).Error()
		panic(fmt.Errorf("load game config error, error: %v", err))
	}
	p.gameCfg.SetVersion(p.generateVersion())
	p.loadCallback(p.gameCfg)
	return p.gameCfg
}

// Init ...
func (p *ConfManager) InitTest(app *appframe.Application, paths *conf.GameConfigPathNode, dev bool) {
	p.paths = paths
	//p.reloadCallback = reloadCallback
	p.gameCfg = &conf.GameConfig{}
	err := p.gameCfg.Init(paths)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"paths": paths,
		}).WithError(err).Error()
		panic(fmt.Errorf("load game config error, error: %v", err))
	}
	p.gameCfg.Develop = dev
	p.gameCfg.SetVersion(p.generateVersion())

	// 注册配置重载方法
	//p.initReload(app)

	logrus.Debug("init game config success")
}

func (p *ConfManager) generateVersion() string {
	return fmt.Sprintf("%d", gameutil.GetCurrentTimestamp())
}

//func generateVersion() string {
//	return fmt.Sprintf("%d", gameutil.GetCurrentTimestamp())
//}

// Reload ...
//func (p *ConfManager) Reload() error {
//	//p.mutex.Lock()
//	//defer p.mutex.Unlock()
//	// 缓存旧配置(注意拷贝对象的时候先*在&)
//	tmp := *p.gameCfg
//	oldGameCfg := &tmp
//	// 重载
//	err := p.gameCfg.Reload()
//	if err != nil {
//		logrus.WithError(err).Error("reload game config, reload game config error")
//		// 还原配置
//		p.gameCfg = oldGameCfg
//		return err
//	}
//
//	p.gameCfg.SetVersion(p.generateVersion())
//	logrus.WithFields(logrus.Fields{
//		"time": gameutil.GetCurrentTime().String(),
//	}).Debug("reload game config success")
//	return nil
//}

// GetConfig ...
func (p *ConfManager) GetConfig() *conf.GameConfig {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.gameCfg
}

// GetVersion ...
func (p *ConfManager) GetVersion() string {
	return p.gameCfg.GetVersion()
}

func (p *ConfManager) initReload(app *appframe.Application) {
	// 重载配置方法注册
	if app == nil {
		return
	}
	appframe.ListenRequestSugar(app, p.onAsAllReqReload)
}

func (p *ConfManager) onAsAllReqReload(sender appframe.Requester, req *smsg.AdAllReqReload) {
	resp := &smsg.AdAllRespReload{}
	gameCfg := &conf.GameConfig{}
	err := gameCfg.Init(p.paths)
	if err != nil {
		resp.ErrCode = smsg.AdAllRespReload_ReloadErr
		resp.ErrMsg = err.Error()
		sender.Resp(resp)
		logrus.WithFields(logrus.Fields{
			"paths": p.paths,
		}).WithError(err).Error()
		return
	}
	gameCfg.Develop = p.develop
	p.loadCallback(gameCfg)
	sender.Resp(resp)
}
