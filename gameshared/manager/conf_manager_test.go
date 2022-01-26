package manager

import (
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameutil"
	"testing"
)

func TestConfManager(t *testing.T) {
	cfg, err := config.ParseConfigFile("../../bin/app.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	netConfigFile := "../../bin/netconfig.json"
	app, err := appframe.NewApplication(netConfigFile, "game")
	if err != nil {
		t.Error(err)
		return
	}

	cfgMgr := &ConfManager{}
	// 初始化配置
	cfgMgr.InitTest(app, cfg.GameCfgPath, cfg.Develop)
	// 获取配置
	t.Logf("init game config success, version: %v, time: %v", cfgMgr.GetVersion(), gameutil.GetCurrentTime().Local().String())

	//item,_ :=cfgMgr.GetConfig().GetItemByItemID(1100001)
	//item, ok := cfgMgr.GetConfig().ActiveCenter.GetModeRankReward(110001, 999)
	//t.Log("rank reward:", item, ok)
	//
	//// 重载配置
	////time.Sleep(time.Second * 1)
	////err = cfgMgr.Reload()
	//if err != nil {
	//	t.Errorf("reload game config error, err: %v", err)
	//	return
	//}

	t.Logf("reload game config success, version: %v, time: %v", cfgMgr.GetVersion(), gameutil.GetCurrentTime().Local().String())
}
