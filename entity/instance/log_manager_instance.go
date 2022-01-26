package instance

import (
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/gameshared/config"
	"sanguosha.com/sgs_herox/gameshared/logproducer"
)

var logMgr *logproducer.LogManager

func LogMgr() *logproducer.LogManager {
	return logMgr
}

func InitLogManagerInstance(app *appframe.Application, cfg *config.AppConfig) error {
	ret, err := logproducer.New(cfg, app.ID())
	if err != nil {
		return err
	}
	logMgr = ret
	return nil
}
