package lobby

import (
	"sanguosha.com/sgs_herox/gameshared/config/conf"
	"sync"
)

var gameConfig *conf.GameConfig
var mutex sync.RWMutex

//多线程访问j时，为防止配置重载中，在一个函数上下文情况下取到的配置不同，在一个上下文情况下要只取出一份用
func GetConfig() *conf.GameConfig {
	mutex.RLock()
	defer mutex.RUnlock()
	return gameConfig
}

func SetConfig(config *conf.GameConfig) {
	mutex.Lock()
	defer mutex.Unlock()
	gameConfig = config
}
