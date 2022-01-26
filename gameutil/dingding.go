package gameutil

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sanguosha.com/baselib/util"
	"strings"
)

type DingDingMsg struct {
	Type string            `json:"msgtype"`
	Text map[string]string `json:"text"`
}

var (
	dingdingURLBase = "https://oapi.dingtalk.com/robot/send?access_token="
)

func SendToDingDing(tokens string, sep string, content string) {
	util.SafeGo(func() {
		msg := &DingDingMsg{Type: "text"}
		msg.Text = make(map[string]string)
		msg.Text["content"] = content
		jsonValue, _ := json.Marshal(msg)

		list := strings.Split(tokens, sep)
		for _, t := range list {
			if len(t) == 0 {
				continue
			}
			resp, err := http.Post(dingdingURLBase+t, "application/json", bytes.NewBuffer(jsonValue))
			if err != nil {
				logrus.Debug("SendToDingDing ", err)
			} else {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				logrus.Debug("SendToDingDing ", string(body), " ", err)
			}
		}
	})
}
