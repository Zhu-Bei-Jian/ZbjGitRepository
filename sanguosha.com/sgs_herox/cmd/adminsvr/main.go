package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sanguosha.com/sgs_herox/admin"
	"time"

	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
)

var help = flag.Bool("h", false, "help")
var netconfigFile = flag.String("netconfig", "netconfig.json", "netconfig file")
var appName = flag.String("name", "admin", "netconfig key")
var configFile = flag.String("config", "app.yaml", "app config file")
var port = flag.Int("port", 8080, "listen port")
var webRoot = flag.String("webroot","","web root")

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
		logrus.WithField("name", *appName).WithError(err).Panic("New application fail")
	}
	err = admin.InitAdmin(app, *configFile, fmt.Sprintf("0.0.0.0:%d", *port),*webRoot)
	if err != nil {
		logrus.WithField("name", *appName).WithError(err).Panic("Init application fail")
	}
	app.Run()
}
