package admin

import (
	"errors"
	"sanguosha.com/sgs_herox/gameutil"
	"sanguosha.com/sgs_herox/proto/smsg"
)

func reqUserSMSGData(userId uint64, dataTypes ...smsg.AllEnReqUserInfo_UserDataTyp) (*smsg.UserDataContent, error) {
	resp, err := entityMgr.RequestCall(&smsg.AllEnReqUserInfo{
		Userid:    userId,
		DataTypes: dataTypes,
	}, reqCallTimeout)

	if err != nil {
		return nil, errors.New("request userInfo error")
	}
	msg := resp.(*smsg.AllEnRespUserInfo)

	if msg.ErrCode != smsg.AllEnRespUserInfo_Invalid {
		return nil, errors.New("request userInfo error")
	}

	return msg.Data, nil
}

func SMSG2UserInfo(u *smsg.User) *UserInfo {
	return &UserInfo{
		UserID:         u.UserBrief.UserID,
		Account:        u.UserBrief.Account,
		Nickname:       u.UserBrief.Nickname,
		Level:          u.UserBrief.Level,
		ThirdAccountId: u.UserBrief.ThirdAccountId,
		HeadImgUrl:     u.UserBrief.HeadImgUrl,
		WinCount:       u.UserBrief.WinCount,
		LoseCount:      u.UserBrief.LoseCount,
		//Score:          u.UserBrief.Score,
		IsOnline:   u.IsOnline,
		CreateTime: gameutil.ParseTimestamp2String(u.CreateTime),
		LoginTime:  gameutil.ParseTimestamp2String(u.LoginTime),
		LogoutTime: gameutil.ParseTimestamp2String(u.LogoutTime),
		LoginIP:    u.UserBrief.LoginIP,
	}
}

type UserInfo struct {
	UserID   uint64 `json:"userID"`
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
	Level    int32  `json:"level"`

	ThirdAccountId string `json:"thirdAccountId"`
	HeadImgUrl     string `json:"headImgUrl"`
	WinCount       int32  `json:"winCount"`
	LoseCount      int32  `json:"loseCount"`
	Score          int32  `json:"score"`

	IsOnline   bool   `json:"isOnline"`
	CreateTime string `json:"createTime"`
	LoginTime  string `json:"loginTime"`
	LogoutTime string `json:"logoutTime"`
	LoginIP    string `json:"loginIP"`
}
