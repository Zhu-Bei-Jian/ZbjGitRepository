package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox/cmd/test/client"
	"sanguosha.com/sgs_herox/gameutil"
	"time"
)

var (
	help    = flag.Bool("h", false, "help")
	cfgFile = flag.String("config", "test.yaml", "app config")
	cfg     = &client.TestConf
)

func main() {
	rand.Seed(time.Now().Unix())
	fmt.Println("开始时间:", time.Now())
	client.Init(true)

	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	abs, err := filepath.Abs(*cfgFile)
	if err != nil {
		return
	}
	data, err := ioutil.ReadFile(abs)
	if err != nil {
		return
	}
	if err := yaml.Unmarshal([]byte(data), &client.TestConf); err != nil {
		return
	}

	var interval time.Duration
	if cfg.LoginPerSec > 0 {
		interval = time.Duration(gameutil.SafeDivFloat64(float64(time.Second), float64(cfg.LoginPerSec)))
	} else {
		interval = 0
	}
	for i := cfg.Base; i < (cfg.Base + cfg.Limit); i++ {
		createAccount(int32(i))
		time.Sleep(interval)
	}
	client.GoShowCol()
	for {
		time.Sleep(10 * time.Hour)
	}
}

func createAccounts(min, max int) {
	for i := min; i < max; i++ {
		index := i
		util.SafeGo(func() {
			account := fmt.Sprintf("%s%d", cfg.Prefix, index)
			addrLen := int32(len(cfg.Addr))
			addr := cfg.Addr[gameutil.RandNum(addrLen)]
			c := client.New(addr, account, cfg.Prefix, cfg.Base+cfg.Limit)

			c.Run()
		})
	}
}

func createAccount(index int32) {
	util.SafeGo(func() {
		account := fmt.Sprintf("%s%d", cfg.Prefix, index)
		addrLen := int32(len(cfg.Addr))
		addr := cfg.Addr[gameutil.RandNum(addrLen)]
		c := client.New(addr, account, cfg.Prefix, cfg.Base+cfg.Limit)

		c.Run()
	})
}
