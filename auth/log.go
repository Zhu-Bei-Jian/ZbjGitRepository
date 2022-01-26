package auth

//增加注册日志
func addLogRegist(account string, userId uint64, ip string) {
	detail := make(map[string]interface{})
	detail["ip"] = ip
	detail["userId"] = userId
	logMgr.AddLogRegist(account, detail, 0)
}
