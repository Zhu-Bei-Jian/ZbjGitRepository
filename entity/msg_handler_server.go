package entity

import (
	"github.com/sirupsen/logrus"
	"sanguosha.com/baselib/appframe"
	"sanguosha.com/sgs_herox/entity/userdb"
	"sanguosha.com/sgs_herox/proto/smsg"
)

func onAllEnReqUserInfo(sender appframe.Requester, req *smsg.AllEnReqUserInfo) {
	UserDBInstance.GetUser(req.Userid, func(u *userdb.User, err error) {
		userID := req.Userid
		resp := smsg.AllEnRespUserInfo{
			Userid: userID,
			Data:   &smsg.UserDataContent{},
		}
		defer sender.Resp(&resp)

		if err != nil {
			resp.ErrCode = smsg.AllEnRespUserInfo_GetUserError
			logrus.WithField("userID", req.Userid).WithError(err).Error("onAsEnReqUserInfo error")
			return
		}

		loginTime := u.LoginTime()
		logoutTime := u.LogoutTime()

		for _, dataType := range req.DataTypes {
			switch dataType {
			case smsg.AllEnReqUserInfo_UDTUser:
				resp.Data.User = &smsg.User{
					UserBrief:           u.Brief(),
					IsOnline:            u.IsOnline(),
					Status:              u.Status(),
					LoginIP:             u.LoginIP(),
					ThirdAccountId:      u.ThirdAccountId(),
					CreateTime:          u.CreatedTime().Unix(),
					RegisterTime:        u.CreatedTime().Unix(),
					RegisterIP:          u.CreatedIP(),
					LoginTime:           loginTime.Unix(),
					LogoutTime:          logoutTime.Unix(),
					BornTime:            u.BornTime(),
					ContinueLoginDayCnt: u.GetCharInfo().C11LoginDayCnt,
					LoginDayCnt:         u.GetCharInfo().LoginDayCnt,
					OnlineSec:           u.GetOnlineSecTotal(),
					OnlineSecToday:      u.GetOnlineSecToday(),
				}
			}
		}
		return
	})
}
