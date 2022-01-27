package config

const GameID = 100

const (
	// LogTopic 日志频道.
	TopicLog     = "log"
	TopicMonitor = "monitor"
	TopicReport  = "report"
)

const (
	MQNode_Log     = "public"
	MQNode_Report  = "public"
	MQNode_Monitor = "public"
)

const (
	RedisUserCache = "user"
	RedisSgsPublic = "sgs_public"
)

type ServerLoadType int32

const (
	AppWorkerLen ServerLoadType = iota
	EntityWorkerLen
	EntityCacheWorkerLen
	AuthDBWorkerLen
	ResponseMS
	ChanEmailLen
	PayDBWorkerLen
	GateOnLineCount
	ShopMsgWorkerCount
	ShopDBWorkerLen
	OnlineCount
	MsgWorkerGoutineLen
)
