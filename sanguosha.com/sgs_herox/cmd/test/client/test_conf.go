package client

type TestConfig struct {
	Addr []string `yaml:"addr"`

	Prefix      string `yaml:"prefix"`
	Base        int    `yaml:"base"`
	Limit       int    `yaml:"limit"`
	LoginPerSec int    `yaml:"loginPerSec"`

	ReqPerSec int `yaml:"reqPerSec"`
	PingDelay int `yaml:"pingDelay"`

	TestEntity bool `yaml:"testEntity"`
	TestGame   bool `yaml:"testGame"`
}

var TestConf TestConfig
