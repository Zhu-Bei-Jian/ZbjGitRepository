package logproducer

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/mq"
	"sync"
	"sync/atomic"
	"time"
)

type LogTyp int32

const (
	LogTyp_LogTInvalid LogTyp = 0
	// 登入
	LogTyp_LTLogin LogTyp = 101
	// 登出
	LogTyp_LTLogout LogTyp = 102
	// 创角
	LogTyp_LTCreate LogTyp = 103
	// 注册
	LogTyp_LTRegist LogTyp = 104
	// 获得物品
	LogTyp_LTGoodsGet LogTyp = 200
	// 购买物品
	LogTyp_LTGoodsBuy LogTyp = 201
	// 货币获得
	LogTyp_LTCurrencyGet LogTyp = 206
	// 货币消耗
	LogTyp_LTCurrencyConsume LogTyp = 207
	// 游戏记录
	LogTyp_LTGameRecord LogTyp = 301
	// 匹配日志
	LogTyp_LTMatch LogTyp = 401
	// 充值日志
	LogTyp_LTCharge LogTyp = 501
	// 任务接收
	LogTyp_LTTaskReceive LogTyp = 600
	// 任务完成
	LogTyp_LTTaskFinish LogTyp = 601
	// 任务奖励接收
	LogTyp_LTTaskRewardReceive LogTyp = 602
	// 武将经验变动
	LogTyp_LTGeneralExp LogTyp = 701
	// 武将星级变动
	LogTyp_LTGeneralLevel LogTyp = 702
)

type LogMessage struct {
	GameId      int32                  `json:"game_id,omitempty"`
	AreaId      int32                  `json:"area_id,omitempty"`
	ServerId    int32                  `json:"server_id,omitempty"`
	LoginFrom   int32                  `json:"login_from,omitempty"`
	LogTime     int64                  `json:"log_time,omitempty"`
	UserAccount string                 `json:"user_account,omitempty"`
	UserLevel   int32                  `json:"user_level,omitempty"`
	OpType      LogTyp                 `json:"op_type,omitempty"`
	Param1      int64                  `json:"param1,omitempty"`
	Param2      int64                  `json:"param2,omitempty"`
	LogInfo     map[string]interface{} `json:"log_info,omitempty"`
	Opmark      int32                  `json:"opmark,omitempty"`
}

// 获取用户属性接口
type LogUser struct {
	Account string
	//角色等级
	CharLevel int32
	//设备类型
	DeviceType int32
	//参数1
	Param1 int64
	//参数2
	Param2 int64
}

type LogManager struct {
	stopped  int32
	wg       sync.WaitGroup
	gameId   int32
	serverID uint32
	areaID   int32

	producer mq.Producer
	topic    string
}

func New(cfg *config.AppConfig, serverId uint32) (*LogManager, error) {
	mqCfg, exist := cfg.MQNodes[config.MQNode_Log]
	if !exist {
		return nil, errors.New("no MQNode_Log")
	}

	producer, err := mq.NewProducer(mq.Config{
		Open:    mqCfg.Open,
		Type:    mqCfg.Type,
		Address: mqCfg.Address,
	})

	if err != nil {
		return nil, err
	}

	r := &LogManager{
		stopped:  0,
		wg:       sync.WaitGroup{},
		gameId:   config.GameID,
		serverID: serverId,
		areaID:   cfg.GameArea,
		producer: producer,
		topic:    config.TopicLog,
	}
	return r, nil
}

func (p *LogManager) Close() error {
	p.Stop()
	p.producer.Close()
	return nil
}

func (p *LogManager) Stop() {
	atomic.StoreInt32(&p.stopped, 1)
	p.wg.Wait()
}

func (p *LogManager) pub(logData LogMessage) {
	payload, err := json.Marshal(logData)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"data": logData,
		}).WithError(err).Error("pub json.Marshal")
		return
	}

	p.producer.PublishAsync(&mq.Msg{
		Topic: p.topic,
		Data:  payload,
	})
}

func (p *LogManager) isClosed() bool {
	return atomic.LoadInt32(&p.stopped) != 0
}

//注册日志
func (p *LogManager) AddLogRegist(account string, detail map[string]interface{}, deviceType int32) {
	logData := LogMessage{
		GameId:      p.gameId,
		AreaId:      p.areaID,
		LoginFrom:   int32(deviceType),
		LogTime:     time.Now().Unix(),
		UserAccount: account,
		UserLevel:   0,
		OpType:      LogTyp_LTRegist,
		Param1:      0,
		Param2:      0,
		LogInfo:     detail,
		Opmark:      0,
		ServerId:    int32(p.serverID),
	}
	p.pub(logData)
}

//创角日志
func (p *LogManager) AddLogCreate(account string, detail map[string]interface{}, deviceType int32) {
	logData := LogMessage{
		GameId:      p.gameId,
		AreaId:      p.areaID,
		LoginFrom:   int32(deviceType),
		LogTime:     time.Now().Unix(),
		UserAccount: account,
		UserLevel:   0,
		OpType:      LogTyp_LTCreate,
		Param1:      0,
		Param2:      0,
		LogInfo:     detail,
		Opmark:      0,
		ServerId:    int32(p.serverID),
	}

	p.pub(logData)
}

//登入日志
func (p *LogManager) AddLogLogin(u LogUser, detail map[string]interface{}) {
	logData := LogMessage{
		GameId:      p.gameId,
		AreaId:      p.areaID,
		LoginFrom:   u.DeviceType,
		LogTime:     time.Now().Unix(),
		UserAccount: u.Account,
		UserLevel:   u.CharLevel,
		OpType:      LogTyp_LTLogin,
		Param1:      u.Param1,
		Param2:      u.Param2,
		LogInfo:     detail,
		Opmark:      0,
		ServerId:    int32(p.serverID),
	}

	p.pub(logData)
}

//登出日志
func (p *LogManager) AddLogLogout(u LogUser, detail map[string]interface{}) {
	logData := LogMessage{
		GameId:      p.gameId,
		AreaId:      p.areaID,
		LoginFrom:   u.DeviceType,
		LogTime:     time.Now().Unix(),
		UserAccount: u.Account,
		UserLevel:   u.CharLevel,
		OpType:      LogTyp_LTLogout,
		Param1:      u.Param1,
		Param2:      u.Param2,
		LogInfo:     detail,
		Opmark:      0,
		ServerId:    int32(p.serverID),
	}

	p.pub(logData)
}

//物品日志
func (p *LogManager) AddLogGoods(typ LogTyp, u LogUser, detail map[string]interface{}) {
	logData := LogMessage{
		GameId:      p.gameId,
		AreaId:      p.areaID,
		LoginFrom:   u.DeviceType,
		LogTime:     time.Now().Unix(),
		UserAccount: u.Account,
		UserLevel:   u.CharLevel,
		OpType:      typ,
		Param1:      u.Param1,
		Param2:      u.Param2,
		LogInfo:     detail,
		Opmark:      0,
		ServerId:    int32(p.serverID),
	}

	p.pub(logData)
}

//任务日志
func (p *LogManager) AddLogTask(u LogUser, taskID int32, logTyp LogTyp, detail map[string]interface{}) {
	logData := LogMessage{
		GameId:      p.gameId,
		AreaId:      p.areaID,
		LoginFrom:   u.DeviceType,
		LogTime:     time.Now().Unix(),
		UserAccount: u.Account,
		UserLevel:   u.CharLevel,
		OpType:      logTyp,
		Param1:      u.Param1,
		Param2:      u.Param2,
		LogInfo:     detail,
		Opmark:      0,
		ServerId:    int32(p.serverID),
	}
	p.pub(logData)
}

func (p *LogManager) AddLogGameRecord(u LogUser, detail map[string]interface{}) {
	logData := LogMessage{
		GameId:      p.gameId,
		AreaId:      p.areaID,
		LoginFrom:   u.DeviceType,
		LogTime:     time.Now().Unix(),
		UserAccount: u.Account,
		UserLevel:   u.CharLevel,
		OpType:      LogTyp_LTGameRecord,
		Param1:      u.Param1,
		Param2:      u.Param2,
		LogInfo:     detail,
		Opmark:      0,
		ServerId:    int32(p.serverID),
	}
	p.pub(logData)
}

func (p *LogManager) AddLogCharge(u LogUser, detail map[string]interface{}) {
	logData := LogMessage{
		GameId:      p.gameId,
		AreaId:      p.areaID,
		LoginFrom:   u.DeviceType,
		LogTime:     time.Now().Unix(),
		UserAccount: u.Account,
		UserLevel:   u.CharLevel,
		OpType:      LogTyp_LTCharge,
		Param1:      u.Param1,
		Param2:      u.Param2,
		LogInfo:     detail,
		Opmark:      0,
		ServerId:    int32(p.serverID),
	}
	p.pub(logData)
}
