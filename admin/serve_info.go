package admin

import (
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/gameshared"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/smsg"
	"strings"
	"time"
)

func CheckServerInfoInit() error {
	infoMap, err := dbMgr.GetServerInfo()
	if err != nil {
		return err
	}
	add := make(map[string]string)
	CheckOne := func(k, defaultValue string) {
		if _, ok := infoMap[k]; !ok {
			infoMap[k] = defaultValue
			add[k] = infoMap[k]
		}
	}

	CheckOne("ApkVersionLimit", "")
	CheckOne("ApkVersionInfo", "版本过低，请更新客户端")
	CheckOne("ServerOpenTime", gameshared.ServerOpenTime)
	CheckOne("ServerCloseTime", gameshared.ServerCloseTime)
	CheckOne("InnerIp", gameshared.InnerIp)

	CheckOne("DingDingTokens", "")

	CheckOne("UnloadGame", "")
	CheckOne("UnloadSvr", "")

	for k, v := range add {
		err = dbMgr.UpdateServerInfo(k, v)
		if err != nil {
			return err
		}
	}

	serverStateService.ResetServerInfo(infoMap)

	time.AfterFunc(3*time.Second, func() {
		msg := &smsg.SyncServerInfo{}
		for k, v := range infoMap {
			msg.Keys = append(msg.Keys, k)
			msg.Values = append(msg.Values, v)
		}
		SyncServerInfoChange(msg)
	})

	return nil
}

func CheckServerInfo(infoMap map[string]string, key string, value string) (dirty bool, old string, err error) {
	if strings.Index(key, "open-") == 0 || key == "ServerCloseTime" || key == "ServerOpenTime" {
		if value == "" {
			return
		}
		_, e := gameutil.ParseDatetime2Timestamp(value)
		if e != nil {
			err = e
			return
		}
	}
	if key == "" {
		return
	}

	if v, ok := infoMap[key]; ok {
		old = v
		if v != value {
			dirty = true
		}
	} else {
		dirty = true
	}

	if dirty {
		err = dbMgr.UpdateServerInfo(key, value)
		if err != nil {
			return
		} else {
			infoMap[key] = value
		}
	}

	serverStateService.ResetServerInfo(infoMap)

	msg := &smsg.SyncServerInfo{}
	for k, v := range infoMap {
		msg.Keys = append(msg.Keys, k)
		msg.Values = append(msg.Values, v)
	}
	SyncServerInfoChange(msg)
	return
}

func SyncServerInfoChange(msg *smsg.SyncServerInfo) {
	for svrType := sgs_herox.SvrTypeGate; svrType < sgs_herox.SvrTypeEnd; svrType++ {
		switch svrType {
		case sgs_herox.SvrTypeAdmin:
			break
		case sgs_herox.SvrTypeAI:
			break
		default:
			svrIds := app.GetAvailableServerIDs(svrType)
			for _, serverId := range svrIds {
				app.GetServer(serverId).SendMsg(msg)
			}
		}
	}
}

func CheckServerInfoOnConnected(svrID uint32, svrType appframe.ServerType) {
	sendServerInfo := func(svrID uint32) {
		infoMap, e := dbMgr.GetServerInfo()
		if e == nil {
			msg := &smsg.SyncServerInfo{}
			for k, v := range infoMap {
				msg.Keys = append(msg.Keys, k)
				msg.Values = append(msg.Values, v)
			}
			app.GetServer(svrID).SendMsg(msg)
		}
	}
	switch svrType {
	case sgs_herox.SvrTypeAdmin:
		break
	case sgs_herox.SvrTypeAI:
		break
	default:
		sendServerInfo(svrID)
	}
}
