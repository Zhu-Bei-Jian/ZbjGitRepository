package userdb

import (
	"sanguosha.com/sgs_herox/entity/instance"
	"sanguosha.com/sgs_herox/gameshared/logproducer"
	"time"
)

func (u *User) LogUser() logproducer.LogUser {
	return logproducer.LogUser{
		Account:    u.Account(),
		CharLevel:  u.CharLevel(),
		DeviceType: 0,
		Param1:     0,
		Param2:     0,
	}
}

//添加登录日志
func (u *User) AddLogLogin() {
	logMgr := instance.LogMgr()
	detail := make(map[string]interface{})
	detail["ip"] = u.LoginIP()

	logMgr.AddLogLogin(u.LogUser(), detail)
}

//添加登出日志
func (u *User) AddLogLogout() {
	logMgr := instance.LogMgr()
	detail := make(map[string]interface{})

	detail["ip"] = u.LoginIP()
	detail["duration"] = time.Now().Unix() - u.loginTime.Unix()

	logMgr.AddLogLogout(u.LogUser(), detail)
}

func (u *User) AddLogGame(detail map[string]interface{}) {
	logMgr := instance.LogMgr()
	logMgr.AddLogGameRecord(u.LogUser(), detail)
}
