package auth

import (
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"
	"strings"
)

func initAuthMsgHandler(app *appframe.Application) {
	appframe.ListenRequestSugar(app, Auth)
}

// 处理账户验证的业务逻辑.
func Auth(sender appframe.Requester, req *smsg.ReqAuth) {
	token := req.Ticket
	loginType := req.LoginType

	log := logrus.WithFields(logrus.Fields{
		"token":     token,
		"loginType": loginType,
	})

	//兼容处理
	if loginType == gameconf.AccountLoginTyp_ALTInvalid {
		loginType = gameconf.AccountLoginTyp_ALTTest
	}

	resp := &smsg.RespAuth{}
	defer sender.Resp(resp)

	//版本检查
	if ok := checkVersion(req.Version); !ok {
		resp.ErrCode = smsg.RespAuth_ErrVersionNotMatch
		return
	}

	//开服时间检查
	if ok, _ := checkServerTime(); ok {

	} else {
		//不在开服时间内部ip可登陆
		if ok := checkInnerIp(req.Ip); !ok {
			resp.ErrCode = smsg.RespAuth_ErrNotOpenTime
			return
		}
	}

	auth, err := NewAuthenticator(loginType, token)
	if err != nil {
		log.WithError(err).Error("NewAuthenticator")
		resp.ErrCode = smsg.RespAuth_ErrUnknow
		return
	}

	authInfo, originData, err := auth.Auth()
	if err != nil {
		switch err {
		case ErrTicketInvalid:
			resp.ErrCode = smsg.RespAuth_ErrTicketInvalid
		default:
			resp.ErrCode = smsg.RespAuth_ErrUnknow
		}
		log.WithError(err).Error("auth.GetAccountInfo")
		return
	}

	account := makeAccount(loginType, authInfo.ThirdAccountId)

	userId, err := dbDao.getAccountUserID(account)
	if err == sql.ErrNoRows {
		userId, err = dbDao.createAccount(account, int(loginType), req.Ip, authInfo.ThirdAccountId, authInfo.Nickname, authInfo.HeadImgUrl, authInfo.HeadFrameImgUrl, authInfo.Sex, originData)
		if err != nil {
			resp.ErrCode = smsg.RespAuth_ErrCreateAccountFailed
			return
		}
		addLogRegist(account, userId, req.Ip)
	} else if err != nil {
		resp.ErrCode = smsg.RespAuth_ErrLoadAccountFailed
		return
	}

	dbDao.updateLoginInfo(userId, req.Ip)

	resp.Userid = userId
	resp.Account = account
	resp.AuthInfo = authInfo
}

func makeAccount(loginType gameconf.AccountLoginTyp, uuid string) string {
	var prefix string
	switch loginType {
	case gameconf.AccountLoginTyp_ALTTest:
		prefix = "test"
	case gameconf.AccountLoginTyp_ALTTablePark:
		prefix = "tablepark"
	default:
		prefix = "default"
	}

	return prefix + "_" + uuid
}

func onServerInfoUpdate(serverInfos *gameshared.ServerInfos) {
	//logrus.Debug("Auth onServerInfoUpdate")

	v := serverInfos.Get("ServerOpenTime")
	ServerOpenTime, _ = gameutil.ParseDatetime2Timestamp(v)

	v = serverInfos.Get("ServerCloseTime")
	ServerCloseTime, _ = gameutil.ParseDatetime2Timestamp(v)

	v = serverInfos.Get("InnerIp")
	InnerIp = strings.Split(v, ";")

	v = serverInfos.Get("ApkVersionLimit")
	ApkVersion = v
}

func checkServerTime() (bool, error) {
	timeNow := gameutil.GetCurrentTimestamp()
	if timeNow < ServerOpenTime {
		return false, fmt.Errorf("server not open")
	} else if timeNow > ServerCloseTime {
		return false, fmt.Errorf("server not open")
	}
	return true, nil
}

func checkInnerIp(ip string) bool {
	if len(ip) == 0 {
		return false
	}
	for _, v := range InnerIp {
		if ip == v {
			return true
		}
	}
	return false
}

func checkVersion(ver string) bool {
	if ApkVersion == "" {
		return true
	}
	if ver != ApkVersion {
		return false
	}
	return true
}
