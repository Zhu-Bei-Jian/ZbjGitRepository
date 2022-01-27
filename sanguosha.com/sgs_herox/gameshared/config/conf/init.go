package conf

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"strings"
	"time"
)

// GameConfig ...
type GameConfig struct {
	cfgPath *GameConfigPathNode
	cfgNode *cfgConfigNode

	Develop bool
	version string

	Global
	Hero
	Skill
	Buff
}

type GameConfigPathNode struct {
	BaseConfigPath   string `yaml:"base_config_path"`
	StableConfigPath string `yaml:"stable_config_path"`
}

type cfgConfigNode struct {
	baseCfg   *gameconf.GameBaseConfig
	stableCfg *gameconf.GameStableConfig
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Init ...
func (p *GameConfig) Init(path *GameConfigPathNode) error {
	var err error
	p.cfgPath = path
	//从配置文件中加载数据到结构体中
	err = p.parseGameConfig()
	if err != nil {
		return err
	}
	//从结构体中整理数据到各个模块中
	err = p.load()
	if err != nil {
		return err
	}
	return nil
}

func (p *GameConfig) Reload() error {
	//从配置文件中加载数据到结构体中
	err := p.parseGameConfig()
	if err != nil {
		return err
	}
	//从结构体中整理数据到各个模块中
	err = p.load()
	if err != nil {
		return err
	}
	return nil
}

// SetVersion ...
func (p *GameConfig) SetVersion(version string) {
	p.version = version
}

// GetVersion ...
func (p *GameConfig) GetVersion() string {
	return p.version
}

func (p *GameConfig) loadBaseCfg() error {
	err := p.Global.loadConf(p.cfgNode.baseCfg)
	if err != nil {
		return err
	}
	err = p.Hero.loadConf(p.cfgNode.baseCfg)
	if err != nil {
		return err
	}
	err = p.Skill.loadConf(p.cfgNode.baseCfg)
	if err != nil {
		return err
	}
	err = p.Buff.loadConf(p.cfgNode.baseCfg)
	if err != nil {
		return err
	}

	return nil
}

func (p *GameConfig) loadStableCfg() error {
	return nil
}

func (p *GameConfig) parseGameConfig() error {
	p.cfgNode = &cfgConfigNode{}
	//基础配置
	baseContent, err := p.getContent(p.cfgPath.BaseConfigPath)
	if err != nil {
		return err
	}
	p.cfgNode.baseCfg = &gameconf.GameBaseConfig{}
	if strings.Index(p.cfgPath.BaseConfigPath, ".json") != -1 {
		err = json.Unmarshal(baseContent, p.cfgNode.baseCfg)
	} else {
		err = proto.UnmarshalText(string(baseContent), p.cfgNode.baseCfg)
	}
	if err != nil {
		return err
	}

	// 武将
	heroContent, err := p.getContent(p.cfgPath.StableConfigPath)
	if err != nil {
		return err
	}
	p.cfgNode.stableCfg = &gameconf.GameStableConfig{}
	if strings.Index(p.cfgPath.StableConfigPath, ".json") != -1 {
		err = json.Unmarshal(heroContent, p.cfgNode.stableCfg)
	} else {
		err = proto.UnmarshalText(string(heroContent), p.cfgNode.stableCfg)
	}
	if err != nil {
		return err
	}

	return nil
}

func (p *GameConfig) load() error {

	err := p.loadBaseCfg()
	if err != nil {
		logrus.WithError(err).Error("initBaseCfg error")
		return err
	}
	err = p.loadStableCfg()
	if err != nil {
		logrus.WithError(err).Error("initStableCfg error")
		return err
	}
	return nil
}

func (p *GameConfig) getContent(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, err
}
