package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"sanguosha.com/baselib/util"
	"sanguosha.com/sgs_herox/gameutil"
	"time"
)

const DingDingURLBase = "https://oapi.dingtalk.com/robot/send?access_token="

func newDingDingReporter(token string) *DingDingReporter {
	return &DingDingReporter{
		token: token,
	}
}

type DingDingReporter struct {
	token string
}

type DingDingReq struct {
	Type string            `json:"msgtype"`
	Text map[string]string `json:"text"`
}

type DingDingResp struct {
	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (p *DingDingReporter) Send(data string) {
	content := fmt.Sprintf("%s\n时间:%s", data, gameutil.ParseTimestamp2String(time.Now().Unix()))
	util.SafeGo(func() {
		err := p.sendSync(content)
		if err != nil {
			logrus.WithError(err).Error("DingDingReporter Send")
		}
	})
}

func (p *DingDingReporter) sendSync(content string) error {
	req := &DingDingReq{Type: "text"}
	req.Text = make(map[string]string)
	req.Text["content"] = content
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	resp, err := http.Post(DingDingURLBase+p.token, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var ret DingDingResp
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return err
	}

	if ret.ErrCode != 0 {
		return fmt.Errorf("resp errcode:%d errmsg:%s", ret.ErrCode, ret.ErrMsg)
	}

	return nil
}
