package auth

import (
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"sync"
)

var testUsersGlobal sync.Map

type InsideTestAuthenticator struct {
	Authenticator
	token string
}

func (p *InsideTestAuthenticator) Auth() (*gamedef.AuthInfo, string, error) {
	token := p.token
	v, ok := testUsersGlobal.Load(token)
	if ok {
		return v.(*gamedef.AuthInfo), "", nil
	}

	account := "hi_" + token
	user := &gamedef.AuthInfo{
		UserId:          0,
		Nickname:        account,
		Sex:             0,
		HeadImgUrl:      "https://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIB4iapibkZyGtwrqoHIfA4gpyUNwQ63SzT0BfVFDsDsgLIG3JO1VTrQqqQBmPHTBcxCcNXoa1VBI0Q/132",
		HeadFrameImgUrl: "https://jqsj-oss-online.oss-cn-hangzhou.aliyuncs.com/md2/avatarFrame/火焰球球.png",
		Birthday:        "",
		RegisterTime:    0,
		ThirdAccountId:  token,
	}

	testUsersGlobal.Store(token, user)
	return user, "", nil
}
