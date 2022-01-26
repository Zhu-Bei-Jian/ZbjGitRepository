package manager

import (
	"encoding/json"
	"errors"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/mq"
	"sanguosha.com/sgs_herox/gameutil"
	"time"
)

const (
	Timestamp = "timestamp" // 时间戳
	ReportID  = "id"
)

const (
	// ReportIDMemoryStatus 内存监控.
	ReportIDMemoryStatus = 101
	// ReportIDHeartBeat 心跳.
	ReportIDHeartBeat = 107
	// ReportIDServerDisconnect 服务断开连接.
	ReportIDServerDisconnect = 201
	// ReportIDServerConnect 服务连接.
	ReportIDServerConnect = 209

	// ReportIDReportGameStart 游戏开始消息.
	ReportIDReportGameStart = 1004
	// ReportIDReportGameOver 游戏结束消息.
	ReportIDReportGameOver = 1005
	// ReportIDReportScoreChanged 分数变更消息.
	ReportIDReportScoreChanged = 1006

	ReportID_Online = 1001
)

type reportMsg struct {
	AppId     int32                  `json:"app_id"`
	AreaId    int32                  `json:"area_id"`
	ServerId  int32                  `json:"server_id"`
	EventInfo map[string]interface{} `json:"event_info"`
}

//状态推送管理器
type reportManager struct {
	producer mq.Producer
	topic    string
	appId    int64
	gameArea int32
	app      *appframe.Application
}

func NewReportManager(cfg *config.AppConfig, app *appframe.Application) (*reportManager, error) {
	mqCfg, exist := cfg.MQNodes[config.MQNode_Report]
	if !exist {
		return nil, errors.New("no MQNode_Report")
	}

	producer, err := mq.NewProducer(mq.Config{
		Open:    mqCfg.Open,
		Type:    mqCfg.Type,
		Address: mqCfg.Address,
	})

	if err != nil {
		return nil, err
	}

	r := &reportManager{
		producer: producer,
		topic:    config.TopicReport,
		appId:    cfg.AppID,
		gameArea: cfg.GameArea,
		app:      app,
	}
	return r, nil
}

func (p *reportManager) Close() {
	p.producer.Close()
}

func (p *reportManager) pub(detail map[string]interface{}) error {
	kafkaMsg := &reportMsg{
		AppId:     int32(p.appId),
		AreaId:    int32(p.gameArea),
		ServerId:  int32(p.app.ID()),
		EventInfo: detail,
	}

	data, err := json.Marshal(kafkaMsg)
	if err != nil {
		return err
	}

	p.producer.PublishAsync(&mq.Msg{
		Topic: p.topic,
		Data:  data,
	})

	return nil
}

func (p *reportManager) PushGateOnline(num int32, ip string) {
	detail := make(map[string]interface{})
	detail["online"] = num
	detail[Timestamp] = gameutil.GetCurrentTimestamp()
	detail[ReportID] = ReportID_Online
	detail["ip"] = ip
	p.pub(detail)
}

func (p *reportManager) PushOnline(num int32) {
	detail := make(map[string]interface{})
	detail["online"] = num
	detail[Timestamp] = gameutil.GetCurrentTimestamp()
	detail[ReportID] = ReportID_Online
	p.pub(detail)
}

func (p *reportManager) PushServerStatus(errType int32, targetID uint32, isDisconnect bool) {
	detail := make(map[string]interface{})
	detail[Timestamp] = gameutil.GetCurrentTimestamp()
	detail[ReportID] = ReportIDServerConnect
	if isDisconnect {
		detail[ReportID] = ReportIDServerDisconnect
		detail["error_type"] = errType
	}

	detail["my_server_id"] = p.app.ID()
	detail["other_server_id"] = targetID
	detail["server_name"] = p.app.Name()
	p.pub(detail)
}

func (p *reportManager) PushServerHeartBeat() {
	timeNow := time.Now()
	detail := make(map[string]interface{})
	detail[Timestamp] = timeNow.Unix()
	detail[ReportID] = ReportIDHeartBeat
	detail["server_name"] = p.app.Name()
	detail["time"] = timeNow.Format("2006-01-02 15:04:05")
	p.pub(detail)
}
