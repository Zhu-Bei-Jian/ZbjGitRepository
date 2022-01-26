package account

import (
	"github.com/golang/protobuf/proto"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox/proto/smsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sync"
	"time"
)

const reqCallTimeout = time.Second * 30

func initMsgHandler(app *appframe.Application) {
	app.RegisterResponse((*smsg.AcEnRespCacheUserSummary)(nil))

	appframe.ListenSessionMsgSugar(app, func(sender appframe.Session, req *smsg.RouteMessageWithUserID) {
		app.HandleMsgWithUserID(sender, req.MsgID, req.UserID, req.Data)
	})
	//只有上面注册了onGateShopMessage才能使用ListenWithUserIDMsg
	appframe.ListenRequestSugar(app, OnPuAcReqUserSummary)
	appframe.ListenRequestSugar(app, onPuAsReqUserIDByShowID)
}

func OnPuAcReqUserSummary(sender appframe.Requester, req *smsg.PuAcReqUserSummary) {
	loadUserSummary(req.Userids, func(m map[uint64]*gamedef.UserSummary, e error) {
		resp := &smsg.PuAcRespUserSummary{}
		defer sender.Resp(resp)

		if e != nil {
			resp.ErrCode = smsg.PuAcRespUserSummary_ErrSystem
			return
		}

		resp.Summaries = m
	})
}

func loadUserSummary(userIds []uint64, callback func(map[uint64]*gamedef.UserSummary, error)) {
	summaries, unfindUserIds, err := userCacheMgr.GetUserSummaries(userIds)
	if err != nil {
		callback(summaries, err)
		return
	}

	if len(unfindUserIds) == 0 {
		callback(summaries, nil)
		return
	}
	var waitGroup sync.WaitGroup
	summaryChan := make(chan *gamedef.UserSummary, len(unfindUserIds))
	for _, userID := range unfindUserIds {
		waitGroup.Add(1)
		entityMgr.Request(&smsg.AcEnReqCacheUserSummary{
			Userid: userID,
		}, func(msg proto.Message, err error) {
			defer waitGroup.Done()
			if err != nil {
				return
			}
			resp, _ := msg.(*smsg.AcEnRespCacheUserSummary)
			if resp.ErrCode != 0 {
				return
			}
			if resp.Summary == nil {
				return
			}
			summaryChan <- resp.Summary
		}, 30*time.Second)
	}

	util.SafeGo(func() {
		waitGroup.Wait()
		close(summaryChan)
		app.Post(func() {
			for {
				summary, ok := <-summaryChan
				if !ok {
					break
				}
				summaries[summary.UserId] = summary
			}

			callback(summaries, nil)
		})
	})
}

func onPuAsReqUserIDByShowID(sender appframe.Requester, req *smsg.PuAcReqQueryUserID) {
	resp := &smsg.PuAcRespQueryUserID{}
	defer sender.Resp(resp)
	switch req.QueryType {
	case smsg.PuAcReqQueryUserID_ByAccount:
		account := req.ParamStr
		userId, err := dbDao.FindUserIdByAccount(account)
		if err != nil {
			if err == ErrCanNotFind {
				return
			}
			resp.ErrCode = smsg.PuAcRespQueryUserID_ErrSytem
			return
		}
		resp.UserIds = append(resp.UserIds, userId)
	case smsg.PuAcReqQueryUserID_ByNickName:
		nickName := req.ParamStr
		userIds, err := dbDao.FindUserIdsByNickName(nickName)
		if err != nil {
			resp.ErrCode = smsg.PuAcRespQueryUserID_ErrSytem
			return
		}
		resp.UserIds = userIds
	case smsg.PuAcReqQueryUserID_ByUnionId:
		unionId := req.ParamStr
		userId, err := dbDao.FindUserIdByUnionId(unionId, req.AccountType)
		if err != nil {
			if err == ErrCanNotFind {
				return
			}
			resp.ErrCode = smsg.PuAcRespQueryUserID_ErrSytem
			return
		}
		resp.UserIds = append(resp.UserIds, userId)
	case smsg.PuAcReqQueryUserID_LikeNickName:
		nickName := req.ParamStr
		userIds, err := dbDao.FindUserIdsLikeNickName(nickName)
		if err != nil {
			resp.ErrCode = smsg.PuAcRespQueryUserID_ErrSytem
			return
		}
		resp.UserIds = userIds
	default:
		resp.ErrCode = smsg.PuAcRespQueryUserID_ErrQueryType
		return
	}
}
