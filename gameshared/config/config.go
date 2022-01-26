package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox/gameshared/config/conf"
)

type RedisConfig struct {
	Max      int    `yaml:"max"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	MaxIdle  int    `yaml:"max_idle"`
}

// Config 共用配置项
type AppConfig struct {
	Develop            bool                `yaml:"develop"`
	RunMode            int32               `yaml:"runmode"` //0正常模式，1，白名单模式
	AppID              int64               `yaml:"appid"`
	GameArea           int32               `yaml:"game_area"`
	DBGame             string              `yaml:"db_game"`
	DBAccount          string              `yaml:"db_account"`
	DBGameUserShardCnt int                 `yaml:"db_game_user_shard_count"`
	EntityNodes        []*EntityNodeConfig `yaml:"entity_nodes"`
	WordsFilter        string              `yaml:"words_filter"`

	//Admin              Admin               `yaml:"admin"`
	SessionRoot string `yaml:"session_root"`
	WebRoot     string `yaml:"web_root"`
	//KafkaNode *KafkaNode `yaml:"kafka"`

	MQNodes map[string]*MQNode `yaml:"mq"`

	RedisNodes map[string]*RedisConfig `yaml:"redis"`

	entityTotalShardCnt int

	TablePark TablePark `yaml:"tablepark"`

	GameCfgPath *conf.GameConfigPathNode `yaml:"game_config_path"`
}

// EntityNodeConfig entity 节点配置
type EntityNodeConfig struct {
	SvrID    uint32 `yaml:"svrid"`
	ShardCnt int    `yaml:"shard_cnt"`
	ChanLen  int    `yaml:"chan_len"`

	shardMaxIndex int
}

type Admin struct {
	AuthKey string `yaml:"auth_key"`
}

type MQNode struct {
	Open    bool     `yaml:"open"`
	Type    string   `yaml:"type"`
	Address []string `yaml:"address"`
}

type TablePark struct {
	UrlAuth             string `yaml:"url_auth"`
	NoticeOpen          bool   `yaml:"notice_open"`
	UrlNoticeGameStatus string `yaml:"url_notice_gamestatus"`
	UrlNoticeGameData   string `yaml:"url_notice_gamedata"`
}

//加载节点配置文件并检查有效性
func LoadConfig(fileName string) (*AppConfig, error) {
	cfg, err := ParseConfigFile(fileName)

	if err != nil {
		return nil, err
	}

	err = CheckConfig(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func parseConfigData(data []byte) (*AppConfig, error) {
	var cfg AppConfig
	if err := yaml.Unmarshal([]byte(data), &cfg); err != nil {
		return nil, err
	}

	if cfg.GameCfgPath == nil {
		return nil, errors.New("no game file path config")
	}

	return &cfg, nil
}

// ParseConfigFile ...
func ParseConfigFile(fileName string) (*AppConfig, error) {
	abs, err := filepath.Abs(fileName)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(abs)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(abs)
	if err != nil {
		return nil, err
	}

	cfg, err := parseConfigData(data)
	if err != nil {
		return nil, err
	}

	cfg.GameCfgPath.BaseConfigPath = filepath.Join(dir, cfg.GameCfgPath.BaseConfigPath)
	cfg.GameCfgPath.StableConfigPath = filepath.Join(dir, cfg.GameCfgPath.StableConfigPath)
	return cfg, nil
}

// ParseConfigFileUseExecDir 以可执行文件所在目录为基准.
func ParseConfigFileUseExecDir(filename string) (*AppConfig, error) {
	dir, err := util.GetCurrExecDir()
	if err == nil {
		filename = filepath.Join(dir, filename)
	}
	return ParseConfigFile(filename)
}

// GetConfig ...
//func GetConfig() (*Config, error) {
//	if appCfg == nil {
//		return nil, errors.New("init server config first")
//	}
//	return appCfg, nil
//}

// CheckConfig 检查配置有效性
func CheckConfig(c *AppConfig) error {
	//if c.DBGame == "" {
	//	return errors.New("no db game")
	//}
	if c.DBGameUserShardCnt < 0 {
		return errors.New("db game user shard count invaild")
	} else if c.DBGameUserShardCnt == 0 {
		c.DBGameUserShardCnt = 1
	}
	if len(c.EntityNodes) > 0 {
		m := make(map[uint32]bool, len(c.EntityNodes))
		var shardIndex int
		for _, e := range c.EntityNodes {
			// check duplicate
			if m[e.SvrID] {
				return fmt.Errorf("entity nodes svrid (%d) duplicate", e.SvrID)
			}
			if e.ShardCnt <= 0 {
				return fmt.Errorf("entity node (%d) shard invalid", e.SvrID)
			}
			if e.ChanLen <= 0 {
				return fmt.Errorf("entity node (%d) chan len invalid", e.SvrID)
			}
			shardIndex += e.ShardCnt
			e.shardMaxIndex = shardIndex
			c.entityTotalShardCnt += e.ShardCnt
			m[e.SvrID] = true
		}
	} else {
		return errors.New("no entity nodes")
	}
	if c.WordsFilter == "" {
		return errors.New("no words filter")
	}
	//if c.AntiAddiction.Enable {
	//	if c.AntiAddiction.RedisAddr == "" {
	//		return errors.New("no anti addiction redis config")
	//	}
	//}
	return nil
}

// GetUserEntityID 用户对应的 entity 节点
func (c *AppConfig) GetUserEntityID(userid uint64) uint32 {
	svrid, _ := c.GetUserEntityIDAndShard(userid)
	return svrid
}

// GetUserEntityIDAndShard 用户对应的 entity 节点 id 与 shard index
func (c *AppConfig) GetUserEntityIDAndShard(userid uint64) (svrid uint32, shard int) {
	// hash userid, 要保证每一个分片的请求能够均匀落到每一张分表上
	// base on [https://stackoverflow.com/a/12996028]
	x := userid
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	x = x ^ (x >> 31)
	// 计算对应的分片index.
	i := int(x % uint64(c.entityTotalShardCnt))
	for _, e := range c.EntityNodes {
		if e.shardMaxIndex > i {
			return e.SvrID, e.ShardCnt - (e.shardMaxIndex - i)
		}
	}
	panic("entity shard config error")
}

// IsWordsFilterNull 是否配置为不进行过滤
func (c *AppConfig) IsWordsFilterNull() bool {
	return c.WordsFilter == "null"
}

func (c *AppConfig) GetEntityNodeConfig(svrid uint32) (*EntityNodeConfig, bool) {
	for _, v := range c.EntityNodes {
		if v.SvrID == svrid {
			return v, true
		}
	}
	return nil, false
}
