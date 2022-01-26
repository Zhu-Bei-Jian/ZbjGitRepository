package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sanguosha.com/sgs_herox/gameutil"
	gamedef "sanguosha.com/sgs_herox/proto/def"
	"strconv"
)

type TableParkAuthenticator struct {
	Authenticator
	token string
}

func (p *TableParkAuthenticator) Auth() (*gamedef.AuthInfo, string, error) {
	if p.token == "hello" {
		return p.GetTestUser(), "", nil
	}
	req, err := http.NewRequest("GET", tableParkCfg.UrlAuth, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", p.token)

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	//str := `{"code":98460336,"data":{"avatar":"occaecat eu ut est","birthDay":"dolor consectetur quis","registerTime":"2009-04-23T19:44:12.003Z","sex":40238038,"userId":-9629820,"userNick":"laborum ad"},"msg":"adipisicing in commodo dolor occaecat","status":true}`
	//body := []byte(str)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	fmt.Println(string(body))

	type Result struct {
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Status bool   `json:"status"`
		Data   struct {
			Avatar       string `json:"avatar"`
			AvatarFrame  string `json:"avatarFrame"`
			BirthDay     string `json:"birthDay"`
			RegisterTime string `json:"registerTime"`
			Sex          int    `json:"sex"`
			UserId       int64  `json:"userId"`
			UserNick     string `json:"userNick"`
		} `json:"data"`
	}

	var ret Result
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, "", err
	}

	if !ret.Status {
		if ret.Code == 401 {
			return nil, "", ErrTicketInvalid
		}
		return nil, "", fmt.Errorf("code:%d msg:%s", ret.Code, ret.Msg)
	}

	registerTime, _ := gameutil.ParseDatetime2Timestamp(ret.Data.RegisterTime)

	u := &gamedef.AuthInfo{
		UserId:          uint64(ret.Data.UserId),
		Nickname:        ret.Data.UserNick,
		Sex:             int32(ret.Data.Sex),
		HeadImgUrl:      ret.Data.Avatar,
		HeadFrameImgUrl: ret.Data.AvatarFrame,
		Birthday:        ret.Data.BirthDay,
		RegisterTime:    registerTime,
		ThirdAccountId:  strconv.FormatInt(ret.Data.UserId, 10),
	}

	return u, string(body), nil
}

func (p *TableParkAuthenticator) GetTestUser() *gamedef.AuthInfo {
	u := &gamedef.AuthInfo{
		UserId:          100,
		Nickname:        "hello",
		Sex:             0,
		HeadImgUrl:      "https://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIB4iapibkZyGtwrqoHIfA4gpyUNwQ63SzT0BfVFDsDsgLIG3JO1VTrQqqQBmPHTBcxCcNXoa1VBI0Q/132",
		HeadFrameImgUrl: "https://jqsj-oss-online.oss-cn-hangzhou.aliyuncs.com/md2/avatarFrame/火焰球球.png",
		Birthday:        "",
		RegisterTime:    0,
		ThirdAccountId:  "hello",
	}

	return u
}
