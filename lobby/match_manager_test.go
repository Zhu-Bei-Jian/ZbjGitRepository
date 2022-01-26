package lobby

import (
	"fmt"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	"sanguosha.com/sgs_herox/gameshared/manager"
	"testing"
	"time"
)

func TestMatchManager(t *testing.T) {

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
		SetConfig(conf)
	})
	confManager.LoadConf()

	matchMgr := &MatchManager{
		Testing: true,
	}
	matchMgr.init()

	funChan := make(chan func(), 100)
	ticker := time.NewTicker(time.Second)
	util.SafeGo(func() {
		for {
			select {
			case <-ticker.C:
				for _, v := range matchMgr.Mode2ModeMatchManager {
					oldNum := v.matchingUserCount()
					now := time.Now()
					v.match()
					if v.matchingUserCount() != oldNum {
						fmt.Printf("match before:%v,after:%v matchTime:%v\n", oldNum, v.matchingUserCount(), time.Since(now))
					}
				}
			case v := <-funChan:
				v()
			}
		}
	})

	var userID uint64
	for j := 0; j < 8; j++ {
		time.Sleep(time.Second * time.Duration(1))
		fmt.Printf("insert time:%v\n", time.Now())
		funChan <- func() {
			userID++
			err := matchMgr.joinMatch(1, userID, 1, 1)
			if err != nil {
				t.Error(err)
			}
		}
	}

	time.Sleep(time.Second * 10)
}
