package manager

import (
	"github.com/robfig/cron"
	"sanguosha.com/baselib/appframe"
)

const (
	AM0spec         = "0 0 0 * * ?"
	AM5spec         = "0 0 5 * * ?"
	FM6spec         = "0 0 18 * * ?"
	EveryHour       = "0 0 */1 * * ?"
	EveryMinuteSpec = "0 */1 * * * ?"
	EveryWeekSpec   = "0 0 0 * * 1"
	Every10sec      = "*/10 * * * * ?"
)

type CronManager struct {
	cron *cron.Cron
}

func NewCronMgr() *CronManager {
	p := &CronManager{}
	p.cron = cron.New()
	return p
}

func (p *CronManager) Start() {
	p.cron.Start()
}

func (p *CronManager) Stop() {
	p.cron.Stop()
}

func (p *CronManager) AddAppStatusPushJob(app *appframe.Application, statusPushMgr *reportManager) {
	//p.Add10SecJob(func() {
	//	//statusPushMgr.PushNatsServerLoad(config.AppWorkerLen, int32(app.WorkerLen()))
	//	//statusPushMgr.PushNatsServerLoad(config.MsgWorkerGoutineLen, app.MsgWorkerGCount())
	//	start := time.Now().UnixNano()
	//	app.Post(func() {
	//		//elapse := (time.Now().UnixNano() - start) / 1e6
	//		//statusPushMgr.PushNatsServerLoad(config.ResponseMS, int32(elapse))
	//	})
	//})

	p.AddMinuteJob(func() {
		statusPushMgr.PushServerHeartBeat()
	})
}

func (p *CronManager) AddFunc(spec string, cmd func()) error {
	return p.cron.AddFunc(spec, cmd)
}

func (p *CronManager) AddWeekJob(cmd func()) error {
	return p.cron.AddFunc(EveryWeekSpec, cmd)
}

func (p *CronManager) AddDayJob(cmd func()) error {
	return p.cron.AddFunc(AM0spec, cmd)
}

func (p *CronManager) AddMinuteJob(cmd func()) error {
	return p.cron.AddFunc(EveryMinuteSpec, cmd)
}

func (p *CronManager) Add10SecJob(cmd func()) error {
	return p.cron.AddFunc(Every10sec, cmd)
}
