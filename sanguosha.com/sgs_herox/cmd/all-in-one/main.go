package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"sync"

	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/appframe/master"
	"sanguosha.com/baselib/framework/netcluster"
	"sanguosha.com/baselib/log"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox/admin"
	"sanguosha.com/sgs_herox/auth"
	"sanguosha.com/sgs_herox/entity"
	"sanguosha.com/sgs_herox/game"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/gate"
	"sanguosha.com/sgs_herox/lobby"
)

var (
	help          = flag.Bool("h", false, "help")
	netconfigFile = flag.String("netconfig", "netconfig.json", "netconfig file")
	appconfigFile = flag.String("config", "app.yaml", "app config")
	noMaster      = flag.Bool("noMaster", false, "ignore Master server")
	noGame        = flag.Bool("noGame", false, "ignore game server")
	noGate        = flag.Bool("noGate", false, "ignore gate server")
	noEntity      = flag.Bool("noEntity", false, "ignore entity server")
	WebRoot       = flag.String("webroot", "", "web root")
)

var (
	masterName   = "master"
	gateName     = "gate"
	authName     = "auth"
	lobbyName    = "lobby"
	entityName   = "entity"
	gameName     = "game"
	emailName    = "email"
	friendName   = "friend"
	adminName    = "admin"
	adminPort    = 9090
	shopName     = "shop"
	payName      = "pay"
	jobName      = "job"
	rankName     = "rank"
	accountName  = "account"
	activityName = "activity"
	wechatName   = "wechat"
)

func init() {
	master.DisableMasterInitGlobalLogrus = true
	appframe.DisableApplicationInitGlobalLogrus = true
}

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	log.InitLogrus(&log.Config{
		Name:  "sgs_wx",
		Level: int(logrus.DebugLevel),
		Outputs: map[string]map[string]interface{}{
			"file": map[string]interface{}{
				"path":   "./logs",
				"rotate": true,
			},
		},
	})

	util.SafeGo(func() {
		fmt.Println(gameutil.Pprof("0.0.0.0:8888", nil))
	})

	//f, err := os.Create("pprof_report/cpu.prof")
	//if err != nil {
	//	logrus.WithError(err).Panic("create pprof error")
	//}
	//pprof.StartCPUProfile(f)
	//
	//time.AfterFunc(time.Minute, func() {
	//	defer pprof.StopCPUProfile()
	//})

	//{
	//	cfg, err := config.LoadConfig(*appconfigFile)
	//	if err != nil {
	//		logrus.WithError(err).Panic("config.LoadConfig")
	//		return
	//	}
	//	//err = gameshared.Init(cfg)
	//	//if err != nil {
	//	//	logrus.WithError(err).Panic("gameshared.Init")
	//	//	return
	//	//}
	//}

	netconfig, err := netcluster.ParseClusterConfigFile(*netconfigFile)
	if err != nil {
		logrus.WithError(err).Panic("netconfigFile.load", err)
		return
	}

	findOneNode := func(base string) string {
		list := []string{base, base + "1"}
		for _, target := range list {
			for nodeName, _ := range netconfig.Masters {
				if nodeName == target {
					return target
				}
			}
			for nodeName, _ := range netconfig.Slaves {
				if nodeName == target {
					return target
				}
			}
		}
		return base
	}

	wg := sync.WaitGroup{}

	// master
	if !*noMaster {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m, err := master.New(*netconfigFile, findOneNode(masterName))
			if err != nil {
				logrus.WithField("name", masterName).WithError(err).Panic("New master fail")
			}
			m.Run()
		}()
	}

	if !*noGate {
		// gate
		wg.Add(1)
		go func() {
			defer wg.Done()
			app, err := appframe.NewGateApplication(*netconfigFile, findOneNode(gateName))
			if err != nil {
				logrus.WithField("name", gateName).WithError(err).Panic("New gate app fail")
			}
			err = gate.InitGateSvr(app, *appconfigFile)
			if err != nil {
				logrus.WithField("name", gateName).WithError(err).Panic("Init gatesvr fail")
			}
			app.Run()
		}()
	}
	// auth
	wg.Add(1)
	go func() {
		defer wg.Done()
		app, err := appframe.NewApplication(*netconfigFile, findOneNode(authName))
		if err != nil {
			logrus.WithField("name", authName).WithError(err).Panic("New auth app fail")
		}
		err = auth.InitAuthSvr(app, *appconfigFile)
		if err != nil {
			logrus.WithField("name", authName).WithError(err).Panic("Init authsvr fail")
		}
		app.Run()
	}()

	//lobby
	wg.Add(1)
	go func() {
		defer wg.Done()
		app, err := appframe.NewApplication(*netconfigFile, findOneNode(lobbyName))
		if err != nil {
			logrus.WithField("name", lobbyName).WithError(err).Panic("New lobby app fail")
		}
		err = lobby.InitLobbySvr(app, *appconfigFile)
		if err != nil {
			logrus.WithField("name", lobbyName).WithError(err).Panic("Init lobbysvr fail")
		}
		app.Run()
	}()

	if !*noEntity {
		// entity
		wg.Add(1)
		go func() {
			defer wg.Done()
			app, err := appframe.NewApplication(*netconfigFile, findOneNode(entityName))
			if err != nil {
				logrus.WithField("name", entityName).WithError(err).Panic("New entity app fail")
			}
			err = entity.InitEntitySvr(app, *appconfigFile)
			if err != nil {
				logrus.WithField("name", entityName).WithError(err).Panic("Init entitysvr fail")
			}
			app.Run()
		}()
	}

	if !*noGame {
		// game
		wg.Add(1)
		go func() {
			defer wg.Done()
			app, err := appframe.NewApplication(*netconfigFile, findOneNode(gameName))
			if err != nil {
				logrus.WithField("name", gameName).WithError(err).Panic("New game app fail")
			}
			err = game.InitGameSvr(app, *appconfigFile)
			if err != nil {
				logrus.WithField("name", gameName).WithError(err).Panic("Init gamesvr fail")
			}
			app.Run()
		}()
	}

	// admin
	wg.Add(1)
	go func() {
		defer wg.Done()
		app, err := appframe.NewApplication(*netconfigFile, findOneNode(adminName))
		if err != nil {
			logrus.WithField("name", adminName).WithError(err).Panic("New admin app fail")
		}
		err = admin.InitAdmin(app, *appconfigFile, fmt.Sprintf("0.0.0.0:%d", adminPort), *WebRoot)
		if err != nil {
			logrus.WithField("name", adminName).WithError(err).Panic("Init admin server fail")
		}
		app.Run()
	}()
	//if !game.IsZbj() {
	//	game.DingSendMsgTestToTakeMyLord(fmt.Sprintf("Hello，我是机器人小钉。%v,测试服更新，现已完成重启", time.Now().UTC().Local()))
	//}

	wg.Wait()
}
