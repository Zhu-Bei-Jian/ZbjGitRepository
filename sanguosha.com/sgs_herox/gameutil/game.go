package gameutil

import (
	gamedef "sanguosha.com/sgs_herox/proto/def"
)

func UserBrief2HeadInfo(msg *gamedef.UserBrief) *gamedef.HeadInfo {
	return &gamedef.HeadInfo{
		UserID:     msg.UserID,
		ShowID:     0,
		Name:       msg.Nickname,
		Level:      msg.Level,
		Icon:       msg.Icon,
		IconEdge:   msg.IconEdge,
		Star:       0,
		Sex:        0,
		HeadImgUrl: "",
	}
}
