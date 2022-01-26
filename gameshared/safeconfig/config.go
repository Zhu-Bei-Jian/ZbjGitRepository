package safeconfig

import (
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	"sync"
)

var gameConfigLock sync.Mutex
var gameConfig *conf.GameConfig

func SetConfig(c *conf.GameConfig) {
	gameConfigLock.Lock()
	defer gameConfigLock.Unlock()
	gameConfig = c
}

func GetConfig() *conf.GameConfig {
	gameConfigLock.Lock()
	defer gameConfigLock.Unlock()
	return gameConfig
}
