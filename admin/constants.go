package admin

import (
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox"
)

const loginHtmlPath = "/temp/#/login"
const homeHtmlPath = "/temp/#/dashboard"

const keyToken = "token"
const keyLogin = "login"
const keyRole = "role"
const keyAccount = "account"
const keyAuthority = "authority"

const resCode = "code"
const resToken = "token"
const resMsg = "msg"
const resHref = "href"
const resItems = "items"
const resGuildUsers = "guild_users"
const resGuildApplies = "guild_applies"
const resItem = "item"
const resNum = "num"
const resExp = "exp"
const resLevel = "level"

const adminAccount = "username"
const adminPsw = "pwd"

const msgSuccess = "成功"
const msgErrAccount = "账号密码错误"
const msgLoginAgain = "请刷新页面重新登录"
const msgLogoutSuccess = "注销成功"
const msgEmailAttach = "道具不存在"
const msgInvalid = "数据格式不正确"
const msgUnknown = "未知错误"
const msgServers = "服务器不存在"
const msgConfig = "配置错误"
const msgUserNotExist = "用户不存在"
const msgErrTaskID = "任务ID不存在"
const msgUserNoGuild = "未加入公会"
const msgInvalidLockTyp = "封禁类型不正确"
const msgLockTimeInvalid = "封禁时间不合法"
const msgParam = "参数错误"

const codeSuccess = 0
const codeUnknown = 1
const codeErrData = 2
const codeUserNotExist = 3
const codeConfig = 4
const codeServers = 5
const codeErrAccount = 201
const codeLoginAgain = 202
const codeEmailAttach = 301
const codeTaskID = 401
const maxLiftTime float64 =3600

var reloadServers = []appframe.ServerType{
	sgs_herox.SvrTypeLobby,
	sgs_herox.SvrTypeEntity,
	sgs_herox.SvrTypeGame,
}
