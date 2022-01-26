package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe/master"
)

var help = flag.Bool("h", false, "help")
var netconfigFile = flag.String("netconfig", "netconfig.json", "netconfig file")
var appName = flag.String("name", "master", "netconfig key")
var configFile = flag.String("config", "app.yaml", "app config file")

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	m, err := master.New(*netconfigFile, *appName)
	if err != nil {
		logrus.WithField("name", *appName).WithError(err).Panic("New master fail")
	}
	m.Run()
}
