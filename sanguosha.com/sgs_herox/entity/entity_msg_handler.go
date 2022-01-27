package entity

import (
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/entity/userdb"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"
)

func initEntityMsgHandler(app *appframe.Application) {
	logrus.Info("initEntityMsgHandler !!!!!!!!!!!")

	// 注册主动请求的返回消息
	app.RegisterResponse((*smsg.RespMatch)(nil))

	// 注册其他服务的请求消息回调
	appframe.ListenRequestSugar(app, GetUserData)
	appframe.ListenSessionMsgSugar(app, GetSelfData)

	//游戏服发来的消息
	appframe.ListenMsgSugar(app, onGameOver)

	//all
	appframe.ListenRequestSugar(app, onAllEnReqUserInfo)

	//玩家卡组信息 相关 (编辑牌库、选择牌库)
	appframe.ListenSessionMsgSugar(app, onReqCardGroup)
	appframe.ListenSessionMsgSugar(app, onCardGroupOpt)

}
func onReqCardGroup(sender appframe.Session, msg *cmsg.CReqCardGroups) { //客户端 请求玩家已编辑的卡组信息 以及最后一次使用的卡组
	resp := &cmsg.SRespCardGroups{}
	defer sender.SendMsg(resp)
	user := GetUserBySessionID(sender.ID())

	resp.CardGroups = user.GetCharInfo().CardGroups
	resp.NowUseId = user.GetCharInfo().NowUseId

}

func onCardGroupOpt(sender appframe.Session, msg *cmsg.CCardGroupOpt) { //客户端 通知 有新建卡组  或者修改已有卡组
	resp := &cmsg.SRespCardGroupOpt{}
	defer sender.SendMsg(resp)
	u := GetUserBySessionID(sender.ID())

	if msg.OptType == cmsg.CCardGroupOpt_OptCreat {
		u.Char.AddCardGroup(msg.CardsGroup)
		if len(msg.CardsGroup.HeroIds) != 12 {
			resp.ErrCode = cmsg.SRespCardGroupOpt_ErrCardsCount
			return
		}
		if !IsUniqueCardGroup(msg.CardsGroup.HeroIds) {
			resp.ErrCode = cmsg.SRespCardGroupOpt_ErrCardSame
			return
		}
		resp.SCode = cmsg.SRespCardGroupOpt_SuccessCreat
	}

	if msg.OptType == cmsg.CCardGroupOpt_OptModify {
		var modifyId int
		for id, v := range u.GetCharInfo().CardGroups {
			if v.GroupId == msg.CardsGroup.GroupId {
				modifyId = id
				break
			}
		}
		u.Char.ModifyCardGroup(modifyId, msg.CardsGroup)
		resp.SCode = cmsg.SRespCardGroupOpt_SuccessModify
	}
	if msg.OptType == cmsg.CCardGroupOpt_OptSelect {
		u.Char.SelectCardGroup(msg.CardsGroup.GroupId)
		resp.SCode = cmsg.SRespCardGroupOpt_SuccessSelect
	}

}

func GetUserBySessionID(sid appframe.SessionID) *userdb.User {
	uid := onlineMgr.MapSessionUser[sid]
	return onlineMgr.userIds[uid]
}

// 服务请求用户数据
func GetUserData(sender appframe.Requester, req *smsg.ReqUserData) {
	UserDBInstance.GetUser(req.Userid, func(u *userdb.User, err error) {
		if err != nil {
			rerr := appframe.NewErrorResponse()
			rerr.Code = int32(smsg.RespUserData_ErrUnknow)
			rerr.Msg = err.Error()
			sender.Resp(rerr)
			return
		}
		resp := new(smsg.RespUserData)
		resp.Nickname = u.Nickname()
		resp.Level = u.Level()
		resp.Icon = u.GetCharInfo().Icon
		resp.IconEdge = u.GetCharInfo().IconEdge
		resp.Sex = u.Sex()
		nowId := -1
		for id, v := range u.GetCharInfo().CardGroups {
			if v.GroupId == u.GetCharInfo().NowUseId {
				nowId = id
				break
			}
		}
		if u.GetCharInfo().CardGroups != nil && len(u.GetCharInfo().CardGroups) > nowId && nowId != -1 {
			resp.CardGroup = u.GetCharInfo().CardGroups[nowId]
		} else {
			resp.CardGroup = nil
		}

		sender.Resp(resp)
	})
}

// 请求我的数据
func GetSelfData(sender appframe.Session, msg *smsg.SyncReqMyData) {
	UserDBInstance.GetUser(msg.Userid, func(u *userdb.User, err error) {
		resp := &cmsg.SRespMyData{}
		if err != nil {
			resp.ErrCode = cmsg.SRespMyData_ErrUnknow
			resp.ErrMsg = err.Error()
			//sender.SendMsg(resp)
			AppInstance.GetSession(appframe.SessionID{SvrID: msg.Gateid, ID: msg.Session}).SendMsg(resp)
			return
		}

		brief := u.Brief()

		resp.UserBase = &gamedef.UserBase{
			UserID:     u.ID(),
			CreateTime: brief.CreateTime,
			Nickname:   brief.Nickname,
			Level:      brief.Level,
			Account:    brief.Account,
			IconEdge:   brief.IconEdge,
			Exp:        brief.Exp,
			Sex:        brief.Sex,
			HeadImgUrl: brief.HeadImgUrl,
		}

		//回包前，设置发送状态已准备好，这样消息才能发出去
		u.SetSendReady(true)

		u.SendToClient(resp)
	})
}

func onGameOver(sender appframe.Server, msg *smsg.GsEsGameOver) {
	UserDBInstance.GetUser(msg.Userid, func(u *userdb.User, err error) {
		if err != nil {
			return
		}

		if msg.IsException {

		} else {
			u.Char.AddUpWinLoseCount(msg.GameModeType, msg.WinLoseType)
			u.Char.AddScore(msg.Score)
			u.SetExp(msg.Exp)
			u.SetLevel(msg.Level)
			u.SyncBrief()
		}

		//同步桌上学园
		if u.AccountType() == gameconf.AccountLoginTyp_ALTTablePark {
			daySuccCount := u.GetCharInfo().CurDayWinCountText + u.GetCharInfo().CurDayWinCountVoice
			dayLoseCount := u.GetCharInfo().CurDayLoseCountText + u.GetCharInfo().CurDayLoseCountVoice
			tpNotifier.NotifyGameDataAsync(u.ThirdAccountId(), daySuccCount, daySuccCount+dayLoseCount)
		}

		//日志记录
		detail := make(map[string]interface{})
		detail["gameMode"] = msg.GameModeType
		detail["gameUUID"] = msg.GameId
		detail["startTime"] = gameutil.ParseTimestamp2String(msg.GameStartTime)
		detail["endTime"] = gameutil.ParseTimestamp2String(msg.GameEndTime)
		detail["duration"] = msg.GameEndTime - msg.GameStartTime
		detail["isException"] = msg.IsException
		detail["roundCount"] = msg.RoundCount

		detail["seat"] = msg.SeatId
		detail["roleType"] = msg.RoleType
		detail["winLoseType"] = msg.WinLoseType
		detail["isEscape"] = msg.IsEscape
		detail["wordID"] = msg.WordId
		u.AddLogGame(detail)
	})
}
