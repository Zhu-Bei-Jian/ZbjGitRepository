package gameshared

import (
	"strconv"
	"sync"
)

//server_info初始化值
const (
	//版本相关
	//c-渠道号 版本号
	//c-%d-ignore 不检查更新的版本号
	//s-game
	//s-ai

	ServerOpenTime  = "2019-07-25 11:55:00"
	ServerCloseTime = "2099-07-25 11:55:00"

	InnerIp      = "183.129.243.134;60.190.230.178"
	WhiteChannel = ""

	NameRenamePre   = "请改名"
	NameCheckFilter = "游卡;三国杀;手杀"

	NameCheckKeep   = "管理;管理员;admin;系统;system;gm;gamemaster;game_master;yoka;robot;小杀;小怒;??????;未知;unknown"
	NameCheckRegexp = "(?m)[^a-zA-Z0-9\u4e00-\u9fa5]" //{}\u0020-\u002F\u003a\u003d\u003f\u0040\u007e\u005a-\u005e  \uff5e分割号 横线 逗号\u002c\u002d\u00b7 省略号\u2026 \uFF01-\uFF1F \u3001-\u3020
	TokenTime       = 600

	LoadAI       = "0,2550,5,2500,50,2400,200,200"
	LoadGame     = "0,5250,5,5200,50,5000,200,1000"
	LoadGate     = 10000
	LoadGateAuth = 1000
	LoadFriend   = 10000

	AheadOpen = 0
)

type ServerInfos struct {
	infos map[string]string

	lock sync.Mutex
}

func NewServerInfos() *ServerInfos {
	return &ServerInfos{infos: make(map[string]string)}
}

func (si *ServerInfos) Reset(m map[string]string) {
	si.lock.Lock()
	defer si.lock.Unlock()

	si.infos = make(map[string]string)
	for k, v := range m {
		si.infos[k] = v
	}
	//logrus.Info("ServerInfos Reset ", si.infos)
}

func (si *ServerInfos) Set(k string, v string) {
	si.lock.Lock()
	defer si.lock.Unlock()

	si.infos[k] = v
	//logrus.Info("ServerInfos Set ", k, " ", v)
}

func (si *ServerInfos) SetAll(keys []string, values []string) {
	if len(keys) != len(values) {
		return
	}
	si.lock.Lock()
	defer si.lock.Unlock()

	si.infos = make(map[string]string)
	for i, k := range keys {
		si.infos[k] = values[i]
	}
	//logrus.Info("ServerInfos SetAll ", si.infos)
}

func (si *ServerInfos) Get(k string) string {
	si.lock.Lock()
	defer si.lock.Unlock()

	v, ok := si.infos[k]
	if ok {
		return v
	}
	return ""
}

func (si *ServerInfos) GetInt64(k string, defaultValue int64) int64 {
	si.lock.Lock()
	defer si.lock.Unlock()

	v, ok := si.infos[k]
	if ok {
		i, e := strconv.Atoi(v)
		if e == nil {
			return int64(i)
		}
		return defaultValue
	}
	return defaultValue
}
