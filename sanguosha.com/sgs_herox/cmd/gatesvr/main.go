package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/gate"
)

var help = flag.Bool("h", false, "help")
var netconfigFile = flag.String("netconfig", "netconfig.json", "netconfig file")
var appName = flag.String("name", "gate", "netconfig key")
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

	app, err := appframe.NewGateApplication(*netconfigFile, *appName)
	if err != nil {
		logrus.WithField("name", *appName).WithError(err).Error("New application fail")
		return
	}
	err = gate.InitGateSvr(app, *configFile)
	if err != nil {
		logrus.WithField("name", *appName).WithError(err).Error("Init application fail")
		return
	}
	app.Run()
}
