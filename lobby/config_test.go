package lobby

import (
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	"sanguosha.com/sgs_herox/gameshared/manager"
	"testing"
)

var config1 *conf.GameConfig

func TestConfig(t *testing.T) {

	cfg, err := config.ParseConfigFile("../bin/app.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	netConfigFile := "../bin/netconfig.json"
	app, err := appframe.NewApplication(netConfigFile, "lobby")
	if err != nil {
		t.Error(err)
		return
	}

	confManager := manager.NewConfManager(app, cfg.GameCfgPath, cfg.Develop, func(conf *conf.GameConfig) {
		config1 = conf
	})
	confManager.LoadConf()

}
