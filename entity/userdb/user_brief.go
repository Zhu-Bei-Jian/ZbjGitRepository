package userdb

import (
	"sanguosha.com/sgs_herox/proto/cmsg"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sanguosha.com/sgs_herox/proto/gameconf"
	"sanguosha.com/sgs_herox/proto/smsg"
)

func (u *User) SyncBrief() {
	u.SendToClient(&cmsg.SSyncUserBrief{
		UserBrief: u.Brief(),
	})

	u.SendToLobby(&smsg.EsAllNtfUserBrief{
		Userid:    u.ID(),
		UserBrief: u.Brief(),
	})
}

func (u *User) Brief() *gamedef.UserBrief {
	charInfo := u.GetCharInfo()
	return &gamedef.UserBrief{
		UserID:          u.ID(),
		CreateTime:      u.CreatedTime().Unix(),
		LoginIP:         u.loginIP,
		Nickname:        u.nickname,
		Level:           u.level,
		Account:         u.account,
		AccountType:     gameconf.AccountLoginTyp(u.accountType),
		ThirdAccountId:  u.thirdAccountId,
		Icon:            charInfo.Icon,
		IconEdge:        charInfo.IconEdge,
		Exp:             u.exp,
		Sex:             u.Sex(),
		HeadImgUrl:      u.headImgUrl,
		HeadFrameImgUrl: u.headFrameImgUrl,
		WinCount:        charInfo.WinCountText + charInfo.WinCountVoice,
		LoseCount:       charInfo.LoseCountText + charInfo.LoseCountVoice,
	}
}
