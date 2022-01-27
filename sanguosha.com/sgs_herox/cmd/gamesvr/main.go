package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"sanguosha.com/baselib/util"
	"time"

	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	game "sanguosha.com/sgs_herox/game"
)

var help = flag.Bool("h", false, "help")
var netconfigFile = flag.String("netconfig", "netconfig.json", "netconfig file")
var appName = flag.String("name", "game", "netconfig key")
var configFile = flag.String("config", "app.yaml", "app config file")

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	rand.Seed(time.Now().Unix())

	//{
	//	cfg, err := config.LoadConfig(*configFile)
	//	if err != nil {
	//		logrus.WithError(err).Panic("config.LoadConfig")
	//		return
	//	}
	//	err = gameshared.Init(cfg)
	//	if err != nil {
	//		logrus.WithError(err).Panic("gameshared.Init")
	//		return
	//	}
	//}

	app, err := appframe.NewApplication(*netconfigFile, *appName)
	if err != nil {
		logrus.WithField("name", *appName).WithError(err).Error("New application fail")
		return
	}

	appID := app.ID()
	if appID < 10000 {
		util.SafeGo(func() {
			fmt.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 20000+appID), nil))
		})
	}

	err = game.InitGameSvr(app, *configFile)
	if err != nil {
		logrus.WithField("name", *appName).WithError(err).Error("Init application fail")
		return
	}
	app.Run()
}
