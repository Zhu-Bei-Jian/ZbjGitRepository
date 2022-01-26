package admin

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"sanguosha.com/sgs_herox"
	"sanguosha.com/sgs_herox/proto/smsg"
	"strconv"
	"strings"
	"time"
)

func gmSendHorseLamp(cmd string, callback func(s string)) {
	cmdList := strings.Split(cmd, " ")
	lenCmdList := len(cmdList)
	if lenCmdList < 2 {
		callback("没有指定广播消息")
		return
	}
	msg := &smsg.ServerNotice{
		AppID:    app.ID(),
		NoticeID: 1,
		Msg:      cmdList[1],
	}
	gateids := app.GetAvailableServerIDs(sgs_herox.SvrTypeGate)
	for _, gateid := range gateids {
		app.GetServer(gateid).SendMsg(msg)
	}
	callback("\n命令执行完成")
}

func gmCloseServer(cmd string, callback func(s string)) {
	cmdList := strings.Split(cmd, " ")
	lenCmdList := len(cmdList)
	if lenCmdList < 2 {
		callback("没有指定目标服务器id")
		return
	}
	serverid, _ := strconv.Atoi(cmdList[1])
	resp, err := app.GetServer(uint32(serverid)).RequestCall(&smsg.AdAllReqCloseServer{}, 15*time.Second)

	func(resp proto.Message, err error) {
		if err != nil {
			callback(fmt.Sprintf("CloseServer %d %v", serverid, err))
			return
		}
		rMsg := resp.(*smsg.AdAllRespCloseServer)
		callback(fmt.Sprintf("CloseServer %d %v", serverid, rMsg.ErrCode))
	}(resp, err)

	callback("\n命令执行完成")
}

func gmGetServerStatus(callback func(s string)) {
	for svrType := sgs_herox.SvrTypeGate; svrType < sgs_herox.SvrTypeEnd; svrType++ {
		if svrType == sgs_herox.SvrTypeAdmin {
			continue
		}
		ids := app.GetAvailableServerIDs(svrType)
		if len(ids) != 0 {
			callback(fmt.Sprintf("\n%v %v", svrType, ids))
		}
	}
	callback("\n命令执行完成")
}

func gmGetUserInfo(cmd string, callback func(s string)) {
	cmdList := strings.Split(cmd, " ")
	lenCmdList := len(cmdList)
	if lenCmdList < 2 {
		callback("参数不对")
		return
	}
	name := cmdList[1]
	userIds, err := accountService.GetUserIdsByNickNameSync(name, true)
	if err != nil {
		callback(err.Error())
		return
	}

	for _, userId := range userIds {
		resp, err := reqUserSMSGData(userId, smsg.AllEnReqUserInfo_UDTUser)
		if err != nil {
			callback(err.Error())
			return
		}
		data := SMSG2UserInfo(resp.User)
		s, e := json.MarshalIndent(data, "", "\t")
		if e != nil {
			continue
		}
		callback(string(s))
	}
}
