package netcluster

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"reflect"
	"sort"
	"strings"

	"sanguosha.com/baselib/log"
)

// ServerConf 服务器节点通用配置
type ServerConf struct {
	ServerName string `json:"-"`
	ServerID   uint32 `json:"id"`
	ServerType uint32 `json:"type"`

	Log log.Config `json:"log"`

	HTTPPProfHost string `json:"profhost,omitempty"`
	HTTPPProfPort int    `json:"profport,omitempty"`

	DebugPrintMessage    bool `json:"printMsg,omitempty"`
	PrintLoadLevelStatus bool `json:"printLoadLevelStatus,omitempty"`
}

// MasterConf Master 服务器配置
type MasterConf struct {
	ServerConf
	ListenAddr   string            `json:"addr"`
	MaxConnCnt   int               `json:"maxconn,omitempty"`
	SlaveWeights map[uint32]uint32 `json:"slweights"`
	Shiftload    uint32            `json:"shiftload"`
}

// SlaveConf Slave 服务器配置
type SlaveConf struct {
	ServerConf
	SubscribedTypes []uint32 `json:"scbtype"`

	// for gate
	UseWebsocket bool   `json:"useWebsocket"`
	OpenTLS      bool   `json:"openTLS"`  //是否开启TLS，目前只支持websocket
	CertFile     string `json:"certFile"` //证书文件路径
	KeyFile      string `json:"keyFile"`  //key文件路径

	ListenAddr           string   `json:"addr"`
	MaxConnCnt           int      `json:"maxconn,omitempty"`
	MasterIDs            []uint32 `json:"masters"`
	DisableCrypto        bool     `json:"disableCrypto"`
	DisableWSCheckOrigin bool     `json:"disableWSCheckOrigin"`
}

// IsGate 是否是 Gate 服务.
func (c *SlaveConf) IsGate() bool {
	return c.ListenAddr != ""
}

// ClusterConf 集群配置
type ClusterConf struct {
	// config file name
	FileName string                 `json:"-"`
	FileMd5  string                 `json:"-"`
	Masters  map[string]*MasterConf `json:"masters"`
	Slaves   map[string]*SlaveConf  `json:"slaves"`
}

// LoadNewCfgFile ...
func (c *ClusterConf) LoadNewCfgFile() (*ClusterConf, error) {
	newConfig, err := ParseClusterConfigFile(c.FileName)
	if err != nil {
		logger.WithError(err).Error("Reload config File load failed")
		return nil, err
	}

	if strings.Compare(newConfig.FileMd5, c.FileMd5) == 0 {
		logger.Info("Reload config File no change")
		return newConfig, nil
	}

	return newConfig, nil
}

// IsSameMasterCfg ...
func (c *ClusterConf) IsSameMasterCfg(cfg *MasterConf) bool {
	for _, m := range c.Masters {
		if m.ServerID == cfg.ServerID {
			if reflect.DeepEqual(cfg, m) {
				return true
			}
		}
	}
	return false
}

// IsSameSlaveCfg ...
func (c *ClusterConf) IsSameSlaveCfg(cfg *SlaveConf) bool {
	for _, s := range c.Slaves {
		if s.ServerID == cfg.ServerID {
			if reflect.DeepEqual(cfg, s) {
				return true
			}
		}
	}
	return false
}

// ParseClusterConfigFile ...
func ParseClusterConfigFile(fileName string) (*ClusterConf, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("ParseClusterConfigFile read %v failed, err:%v", fileName, err)
	}

	config, err := ParseClusterConfData(data)
	if err != nil {
		return nil, fmt.Errorf("ParseClusterConfigFile prase %v failed, err:%v", fileName, err)
	}

	err = CheckClusterConf(config)
	if err != nil {
		return nil, err
	}

	config.FileName = fileName
	config.FileMd5 = hex.EncodeToString(GetCfgFileMd5(fileName))
	for k, v := range config.Masters {
		v.ServerName = k
		if v.Log.Name == "" {
			v.Log.Name = k
		}
	}

	for _, v := range config.Slaves {
		if len(v.MasterIDs) == 0 {
			for _, m := range config.Masters {
				v.MasterIDs = append(v.MasterIDs, m.ServerID)
			}
			sort.Slice(v.MasterIDs, func(i, j int) bool {
				return v.MasterIDs[i] < v.MasterIDs[j]
			})
		}
	}

	var allServerTypes = map[uint32]bool{}
	for k, v := range config.Slaves {
		v.ServerName = k
		allServerTypes[v.ServerType] = true
		if v.Log.Name == "" {
			v.Log.Name = k
		}
	}
	for _, v := range config.Slaves {
		if len(v.SubscribedTypes) == 0 {
			for t := range allServerTypes {
				if t != v.ServerType {
					v.SubscribedTypes = append(v.SubscribedTypes, t)
				}
			}
			sort.Slice(v.SubscribedTypes, func(i, j int) bool {
				return v.SubscribedTypes[i] < v.SubscribedTypes[j]
			})
		}
	}

	for _, v := range config.Masters {
		if v.HTTPPProfHost == "" {
			v.HTTPPProfHost = "localhost"
		}
	}
	for _, v := range config.Slaves {
		if v.HTTPPProfHost == "" {
			v.HTTPPProfHost = "localhost"
		}
	}

	return config, nil
}

// ParseClusterConfData ...
func ParseClusterConfData(data []byte) (*ClusterConf, error) {
	var cfg ClusterConf
	if err := json.Unmarshal([]byte(data), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// CheckClusterConf 检查配置是否合法
func CheckClusterConf(config *ClusterConf) error {
	if config == nil || len(config.Masters) == 0 {
		return fmt.Errorf("cluster config is nil or master is nil")
	}

	ids := make(map[uint32]struct{})
	mids := make(map[uint32]struct{})
	addrs := make(map[string]struct{})

	var masterSvrType uint32
	for sn, s := range config.Masters {
		if s == nil {
			return fmt.Errorf("cluster config master.%s == nil", sn)
		}

		if masterSvrType == 0 {
			masterSvrType = s.ServerType
		} else if s.ServerType != masterSvrType {
			return fmt.Errorf("cluster config master.%s diff servertype %d", sn, s.ServerType)
		}

		if _, ok := ids[s.ServerID]; ok {
			return fmt.Errorf("cluster config master.%s.id has used", sn)
		}
		ids[s.ServerID] = struct{}{}
		mids[s.ServerID] = struct{}{}

		if len(s.ListenAddr) > 0 {
			if _, err := net.ResolveTCPAddr("tcp", s.ListenAddr); err != nil {
				return fmt.Errorf("cluster config master:%d listen on illegal adress:%s", s.ServerID, s.ListenAddr)
			}

			if _, ok := addrs[s.ListenAddr]; ok {
				return fmt.Errorf("cluster config master.%s.addr has used", sn)
			}
			addrs[s.ListenAddr] = struct{}{}
		}
	}

	for sn, s := range config.Slaves {
		if s == nil {
			return fmt.Errorf("cluster config servers.%s == nil", sn)
		}

		if s.ServerType == masterSvrType {
			return fmt.Errorf("cluster config master.type == %s.type", sn)
		}

		if _, ok := ids[s.ServerID]; ok {
			return fmt.Errorf("cluster config %s.id has used", sn)
		}
		ids[s.ServerID] = struct{}{}

		for _, mid := range s.MasterIDs {
			if _, ok := mids[mid]; !ok {
				return fmt.Errorf("cluster config %s.mids not a master id, id=%d", sn, mid)
			}
		}

		if len(s.ListenAddr) > 0 {
			if _, err := net.ResolveTCPAddr("tcp", s.ListenAddr); err != nil {
				return fmt.Errorf("cluster config server:%d listen on illegal adress:%s", s.ServerID, s.ListenAddr)
			}

			if _, ok := addrs[s.ListenAddr]; ok {
				return fmt.Errorf("cluster config %s.addr has used", sn)
			}
			addrs[s.ListenAddr] = struct{}{}
		}
	}

	return nil
}

// GetCfgFileMd5 获取配置文件md5值
func GetCfgFileMd5(fileName string) []byte {
	f, err := os.Open(fileName)
	if err != nil {
		return nil
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		fmt.Println(err)
	}
	return h.Sum(nil)
}
